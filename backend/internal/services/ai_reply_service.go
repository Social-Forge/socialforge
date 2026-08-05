package services

import (
	"context"
	"fmt"
	"github/socialforge/internal/entity"
	aiclient "github/socialforge/internal/infra/ai-client"
	"github/socialforge/internal/infra/contextpool"
	"github/socialforge/internal/infra/repository"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// aiHistoryLimit is how many recent messages to feed the model as context.
const aiHistoryLimit = 20

// AIReplyService generates and sends AI agent replies. It is the headline
// differentiator: when a channel has an active AI agent with auto-reply on, an
// inbound customer message triggers a Claude-generated response persisted and
// dispatched as sender_type=ai, metered against the tenant's credit balance.
type AIReplyService struct {
	aiAgentRepo      repository.AIAgentRepository
	messageRepo      repository.MessageRepository
	conversationRepo repository.ConversationRepository
	knowledgeRepo    repository.AIKnowledgeRepository
	playbookRepo     repository.AIPlaybookRepository
	assetRepo        repository.AIAssetRepository
	creditRepo       repository.AICreditLedgerRepository
	aiClient         *aiclient.AIClient
	outbound         *OutboundService
	logger           *zap.Logger
}

func NewAIReplyService(
	aiAgentRepo repository.AIAgentRepository,
	messageRepo repository.MessageRepository,
	conversationRepo repository.ConversationRepository,
	knowledgeRepo repository.AIKnowledgeRepository,
	playbookRepo repository.AIPlaybookRepository,
	assetRepo repository.AIAssetRepository,
	creditRepo repository.AICreditLedgerRepository,
	aiClient *aiclient.AIClient,
	outbound *OutboundService,
	logger *zap.Logger,
) *AIReplyService {
	return &AIReplyService{
		aiAgentRepo:      aiAgentRepo,
		messageRepo:      messageRepo,
		conversationRepo: conversationRepo,
		knowledgeRepo:    knowledgeRepo,
		playbookRepo:     playbookRepo,
		assetRepo:        assetRepo,
		creditRepo:       creditRepo,
		aiClient:         aiClient,
		outbound:         outbound,
		logger:           logger,
	}
}

// metaAIPaused is the conversation.metadata flag set on human handoff.
const metaAIPaused = "ai_paused"

// GenerateAndReply loads the channel's AI agent, builds a persona-aware prompt
// from the recent conversation history, calls the model, and sends the reply.
// It is best-effort: any failure is logged and swallowed so a model outage never
// blocks message ingestion.
func (s *AIReplyService) GenerateAndReply(ctx context.Context, channel *entity.Channel, conversationID uuid.UUID) {
	if channel.AIAgentID == nil || s.aiClient == nil || s.outbound == nil {
		return
	}
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 45*time.Second)
	defer cancel()

	tctx := repository.WithTenantID(subCtx, channel.TenantID)

	// 1. Load the agent, conversation, history, playbooks, knowledge + balance in
	// one tenant tx.
	var agent *entity.AIAgent
	var conv *entity.Convertation
	var history []*entity.Message
	var playbooks []*entity.AIPlaybook
	var knowledge []*entity.AIKnowledge
	var assets []*entity.AIAsset
	var balance int
	err := s.aiAgentRepo.RunInTenantTx(tctx, func(txCtx context.Context) error {
		var err error
		if agent, err = s.aiAgentRepo.FindByID(txCtx, *channel.AIAgentID); err != nil {
			return err
		}
		if conv, err = s.conversationRepo.FindByID(txCtx, conversationID); err != nil {
			return err
		}
		if history, err = s.messageRepo.ListByConversation(txCtx, conversationID, aiHistoryLimit); err != nil {
			return err
		}
		if playbooks, err = s.playbookRepo.ListByAgent(txCtx, agent.ID); err != nil {
			return err
		}
		if knowledge, err = s.knowledgeRepo.ListByAgent(txCtx, agent.ID); err != nil {
			return err
		}
		if assets, err = s.assetRepo.ListByAgent(txCtx, agent.ID); err != nil {
			return err
		}
		balance, err = s.creditRepo.Balance(txCtx, channel.TenantID)
		return err
	})
	if err != nil {
		s.logger.Warn("ai reply: failed to load agent/history", zap.Error(err))
		return
	}
	if agent == nil || !agent.IsActive || !agent.AutoReplyEnabled {
		return
	}
	// Human already took over this conversation -> AI stays silent.
	if conv != nil && conv.Metadata != nil {
		if paused, _ := (*conv.Metadata)[metaAIPaused].(bool); paused {
			return
		}
	}
	if balance <= 0 {
		s.logger.Info("ai reply skipped: insufficient credits",
			zap.String("tenant_id", channel.TenantID.String()), zap.Int("balance", balance))
		return
	}

	// 2. Build the model conversation (chronological) from history.
	msgs := buildAIMessages(history)
	if len(msgs) == 0 {
		return
	}
	// Only reply if the latest turn is from the customer (avoid replying to our
	// own outbound / looping).
	if last := msgs[len(msgs)-1]; last.Role != "user" {
		return
	}
	lastMsg := strings.ToLower(msgs[len(msgs)-1].Content)

	// 2a. Working-hours gate: outside the agent's hours, do not auto-reply (an
	// off-hours notice is sent once, then a human handles it during hours).
	if !aiWithinWorkingHours(agent.WorkingHours, time.Now()) {
		if note := workingHoursMessage(agent.WorkingHours); note != "" {
			_, _ = s.outbound.SendSystemText(subCtx, channel.TenantID, conversationID, note)
		}
		s.logger.Info("ai reply skipped: outside working hours", zap.String("agent", agent.Name))
		return
	}

	// 2b. Handoff detection: explicit request or a safety topic -> pause AI, hand
	// to a human, and stop.
	if reason := handoffReason(lastMsg, agent.Safety); reason != "" {
		s.handoff(subCtx, tctx, channel, conv, conversationID, reason)
		return
	}

	// 2c. Playbook match: highest-priority active playbook whose keyword appears
	// in the customer's message steers the reply and attaches assets.
	matched := matchPlaybook(playbooks, lastMsg)

	// 2d. Knowledge retrieval: semantic (pgvector) when embeddings are enabled,
	// else lexical term-overlap over the loaded rows.
	knowledgeBlock := s.retrieveKnowledge(subCtx, tctx, agent.ID, msgs[len(msgs)-1].Content, knowledge)

	// 3. Call the model with persona + knowledge + playbook context.
	systemPrompt := buildSystemPrompt(agent) + knowledgeBlock + buildPlaybookBlock(matched, assets)
	resp, err := s.aiClient.Chat(subCtx, msgs, systemPrompt)
	if err != nil {
		s.logger.Warn("ai reply: model call failed", zap.Error(err))
		return
	}
	reply := strings.TrimSpace(resp.Content)
	if reply == "" {
		return
	}

	// 4. Send the reply as sender_type=ai.
	sent, err := s.outbound.SendAIText(subCtx, channel.TenantID, conversationID, reply)
	if err != nil {
		s.logger.Warn("ai reply: failed to send", zap.Error(err))
		return
	}

	// 5. Meter the usage against the tenant credit balance (best-effort).
	tokens := resp.TokensIn + resp.TokensOut
	if tokens <= 0 {
		tokens = len(reply) / 4 // rough fallback estimate
	}
	entry := &entity.AICreditLedger{
		TenantID:       channel.TenantID,
		ConversationID: &conversationID,
		MessageID:      &sent.ID,
		Model:          resp.Model,
		InputTokens:    resp.TokensIn,
		OutputTokens:   resp.TokensOut,
	}
	if err := s.creditRepo.RunInTenantTx(tctx, func(txCtx context.Context) error {
		_, err := s.creditRepo.Debit(txCtx, entry, tokens)
		return err
	}); err != nil {
		s.logger.Warn("ai reply: failed to meter credits", zap.Error(err))
	}

	playbookName := ""
	if matched != nil {
		playbookName = matched.Name
	}
	s.logger.Info("ai reply sent",
		zap.String("agent", agent.Name),
		zap.String("model", resp.Model),
		zap.Int("tokens", tokens),
		zap.Int("balance_after", entry.BalanceAfter),
		zap.String("playbook", playbookName),
		zap.Int("knowledge_available", len(knowledge)),
	)
}

// retrieveKnowledge selects reference knowledge for the query: semantic search
// via pgvector when embeddings are enabled (and rows are embedded), otherwise a
// lexical term-overlap over the already-loaded rows.
func (s *AIReplyService) retrieveKnowledge(ctx, tctx context.Context, agentID uuid.UUID, query string, loaded []*entity.AIKnowledge) string {
	if s.aiClient != nil && s.aiClient.EmbeddingsEnabled() {
		vec, err := s.aiClient.Embed(ctx, query)
		if err != nil {
			s.logger.Warn("query embedding failed, lexical fallback", zap.Error(err))
		} else {
			var hits []*entity.AIKnowledge
			e := s.knowledgeRepo.RunInTenantTx(tctx, func(txCtx context.Context) error {
				var err error
				hits, err = s.knowledgeRepo.SearchByEmbedding(txCtx, agentID, vec, 3)
				return err
			})
			if e == nil && len(hits) > 0 {
				s.logger.Info("knowledge retrieval via pgvector", zap.Int("hits", len(hits)))
				return formatKnowledgeBlock(hits)
			}
			if e != nil {
				s.logger.Warn("vector search failed, lexical fallback", zap.Error(e))
			}
		}
	}
	return buildKnowledgeBlock(loaded, query)
}

// handoff pauses the AI on a conversation and hands it to a human: it flags the
// conversation metadata (ai_paused=true), notifies the customer once, and moves
// the conversation to unassigned so an agent picks it up.
func (s *AIReplyService) handoff(ctx, tctx context.Context, channel *entity.Channel, conv *entity.Convertation, conversationID uuid.UUID, reason string) {
	meta := entity.MetDataConfig{}
	if conv != nil && conv.Metadata != nil {
		meta = *conv.Metadata
	}
	meta[metaAIPaused] = true
	meta["handoff_reason"] = reason
	err := s.conversationRepo.RunInTenantTx(tctx, func(txCtx context.Context) error {
		if err := s.conversationRepo.SetMetadata(txCtx, conversationID, meta); err != nil {
			return err
		}
		// Leave it open for a human; unassign so it re-enters the queue.
		return s.conversationRepo.UpdateStatus(txCtx, conversationID, entity.ConversationStatusUnassigned)
	})
	if err != nil {
		s.logger.Warn("ai reply: handoff failed", zap.Error(err))
		return
	}
	_, _ = s.outbound.SendSystemText(ctx, channel.TenantID, conversationID,
		"Baik, permintaan Anda akan kami teruskan ke agen kami. Mohon tunggu sebentar ya 🙏")
	s.logger.Info("ai reply: handed off to human", zap.String("reason", reason), zap.String("conversation_id", conversationID.String()))
}

// buildAIMessages maps persisted messages (returned newest-first) to the model's
// chronological role/content turns, collapsing consecutive same-role turns is
// left to the model. Empty-body and non-text rows are skipped.
func buildAIMessages(history []*entity.Message) []aiclient.Message {
	// Copy + sort ascending by CreatedAt (repo returns DESC).
	ordered := make([]*entity.Message, len(history))
	copy(ordered, history)
	sort.SliceStable(ordered, func(i, j int) bool {
		return ordered[i].CreatedAt.Before(ordered[j].CreatedAt)
	})

	msgs := make([]aiclient.Message, 0, len(ordered))
	for _, m := range ordered {
		if !m.Body.Valid || strings.TrimSpace(m.Body.String) == "" {
			continue
		}
		role := "assistant"
		if m.Direction == entity.MessageDirectionIn || m.SenderType == entity.SenderTypeContact {
			role = "user"
		}
		msgs = append(msgs, aiclient.Message{Role: role, Content: m.Body.String})
	}
	return msgs
}

// buildSystemPrompt composes the effective system prompt from the agent's base
// prompt plus its persona/guardrails/safety configuration.
func buildSystemPrompt(agent *entity.AIAgent) string {
	var b strings.Builder
	if strings.TrimSpace(agent.SystemPrompt) != "" {
		b.WriteString(agent.SystemPrompt)
		b.WriteString("\n\n")
	}

	if agent.Persona != nil {
		p := map[string]interface{}(*agent.Persona)
		writeField(&b, "Your name", p, "name", "agent_name")
		writeField(&b, "Your role/soul", p, "soul", "role")
		writeField(&b, "Tone of voice", p, "tone")
		writeField(&b, "Gender", p, "gender")
		writeField(&b, "Language", p, "language")
		if greeting := firstString(p, "greeting"); greeting != "" {
			fmt.Fprintf(&b, "Opening greeting to use when appropriate: %q\n", greeting)
		}
	}

	if agent.Guardrails != nil {
		if rules := listValues(map[string]interface{}(*agent.Guardrails)); len(rules) > 0 {
			b.WriteString("\nGuardrails you MUST follow:\n")
			for _, r := range rules {
				fmt.Fprintf(&b, "- %s\n", r)
			}
		}
	}
	if agent.Safety != nil {
		if rules := listValues(map[string]interface{}(*agent.Safety)); len(rules) > 0 {
			b.WriteString("\nSafety rules you MUST NOT violate:\n")
			for _, r := range rules {
				fmt.Fprintf(&b, "- %s\n", r)
			}
		}
	}

	out := strings.TrimSpace(b.String())
	if out == "" {
		out = "You are a helpful, professional customer service assistant. Keep replies concise and friendly."
	}
	return out
}

func writeField(b *strings.Builder, label string, m map[string]interface{}, keys ...string) {
	for _, k := range keys {
		if v := firstString(m, k); v != "" {
			fmt.Fprintf(b, "%s: %s\n", label, v)
			return
		}
	}
}

func firstString(m map[string]interface{}, key string) string {
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok {
			return strings.TrimSpace(s)
		}
	}
	return ""
}

// listValues flattens a config map's values into human-readable strings,
// handling both scalar values and arrays (e.g. guardrails: ["no refunds", ...]).
func listValues(m map[string]interface{}) []string {
	var out []string
	// Deterministic order for stable prompts.
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		switch v := m[k].(type) {
		case string:
			if s := strings.TrimSpace(v); s != "" {
				out = append(out, s)
			}
		case []interface{}:
			for _, item := range v {
				if s, ok := item.(string); ok && strings.TrimSpace(s) != "" {
					out = append(out, strings.TrimSpace(s))
				}
			}
		case bool:
			if v {
				out = append(out, k)
			}
		}
	}
	return out
}

// ============================ Working hours ============================

// aiWithinWorkingHours reports whether `now` falls inside the agent's configured
// hours. Config shape (all optional): {enabled, timezone, start:"HH:MM",
// end:"HH:MM", days:[0-6]}. When disabled/unset the agent is always on.
func aiWithinWorkingHours(wh *entity.WorkingHours, now time.Time) bool {
	if wh == nil {
		return true
	}
	m := map[string]interface{}(*wh)
	if enabled, ok := m["enabled"].(bool); !ok || !enabled {
		return true
	}
	loc := time.UTC
	if tz, _ := m["timezone"].(string); tz != "" {
		if l, err := time.LoadLocation(tz); err == nil {
			loc = l
		}
	} else if l, err := time.LoadLocation("Asia/Jakarta"); err == nil {
		loc = l
	}
	now = now.In(loc)

	// Day-of-week filter (0=Sunday). Absent -> every day allowed.
	if days, ok := m["days"].([]interface{}); ok && len(days) > 0 {
		today := int(now.Weekday())
		allowed := false
		for _, d := range days {
			if n, ok := toInt(d); ok && n == today {
				allowed = true
				break
			}
		}
		if !allowed {
			return false
		}
	}

	start := parseHHMM(firstString(m, "start"), 0)
	end := parseHHMM(firstString(m, "end"), 24*60)
	cur := now.Hour()*60 + now.Minute()
	if start <= end {
		return cur >= start && cur < end
	}
	// Overnight window (e.g. 20:00-06:00).
	return cur >= start || cur < end
}

func workingHoursMessage(wh *entity.WorkingHours) string {
	if wh == nil {
		return ""
	}
	m := map[string]interface{}(*wh)
	return firstString(m, "offhours_message")
}

func toInt(v interface{}) (int, bool) {
	switch n := v.(type) {
	case float64:
		return int(n), true
	case int:
		return n, true
	case int64:
		return int(n), true
	}
	return 0, false
}

// parseHHMM turns "HH:MM" into minutes-since-midnight, or def on failure.
func parseHHMM(s string, def int) int {
	if s == "" {
		return def
	}
	var h, min int
	if _, err := fmt.Sscanf(s, "%d:%d", &h, &min); err != nil {
		return def
	}
	return h*60 + min
}

// ============================ Handoff ============================

// handoffPhrases are explicit customer requests to talk to a human.
var handoffPhrases = []string{
	"agen manusia", "customer service manusia", "cs manusia", "bicara dengan manusia",
	"ngomong sama manusia", "orang asli", "manusia asli", "agen asli", "human agent",
	"talk to a human", "speak to a human", "real person", "operator",
}

// handoffReason returns a non-empty reason when the message should be escalated
// to a human: an explicit request, or a safety topic the agent must not handle.
func handoffReason(lastMsg string, safety *entity.AiSafetyConfig) string {
	for _, p := range handoffPhrases {
		if strings.Contains(lastMsg, p) {
			return "customer_requested_human"
		}
	}
	if safety != nil {
		for _, topic := range listValues(map[string]interface{}(*safety)) {
			t := strings.ToLower(strings.TrimSpace(topic))
			if len(t) >= 4 && strings.Contains(lastMsg, t) {
				return "safety_topic:" + t
			}
		}
	}
	return ""
}

// ============================ Playbook ============================

// matchPlaybook returns the highest-priority active playbook whose keyword
// appears in the customer message (playbooks arrive ordered priority DESC).
func matchPlaybook(playbooks []*entity.AIPlaybook, lastMsg string) *entity.AIPlaybook {
	for _, p := range playbooks {
		if !p.IsActive {
			continue
		}
		for _, kw := range p.Keywords {
			kw = strings.ToLower(strings.TrimSpace(kw))
			if kw != "" && strings.Contains(lastMsg, kw) {
				return p
			}
		}
	}
	return nil
}

// buildPlaybookBlock injects the matched playbook's instruction and any attached
// assets into the system prompt so the model follows the scripted response.
func buildPlaybookBlock(matched *entity.AIPlaybook, assets []*entity.AIAsset) string {
	if matched == nil {
		return ""
	}
	var b strings.Builder
	fmt.Fprintf(&b, "\n\nActive playbook %q — follow this instruction for your reply:\n%s\n", matched.Name, matched.Instruction)

	if len(matched.AssetIDs) > 0 && len(assets) > 0 {
		byID := make(map[string]*entity.AIAsset, len(assets))
		for _, a := range assets {
			byID[a.ID.String()] = a
		}
		var lines []string
		for _, id := range matched.AssetIDs {
			if a, ok := byID[id]; ok {
				desc := a.Name
				if a.Description.Valid && a.Description.String != "" {
					desc += " — " + a.Description.String
				}
				lines = append(lines, fmt.Sprintf("%s (%s)", desc, a.Type))
			}
		}
		if len(lines) > 0 {
			b.WriteString("Offer/mention these materials when relevant:\n")
			for _, l := range lines {
				fmt.Fprintf(&b, "- %s\n", l)
			}
		}
	}
	return b.String()
}

// ============================ Knowledge (RAG-lite) ============================

// buildKnowledgeBlock scores knowledge entries by keyword overlap with the
// customer's message and injects the top matches as reference context. This is a
// lightweight lexical retriever; pgvector embeddings are a deeper refinement.
func buildKnowledgeBlock(knowledge []*entity.AIKnowledge, query string) string {
	if len(knowledge) == 0 {
		return ""
	}
	terms := queryTerms(query)
	type scored struct {
		k     *entity.AIKnowledge
		score int
	}
	var ranked []scored
	for _, k := range knowledge {
		hay := strings.ToLower(k.Title + " " + k.Content)
		score := 0
		for t := range terms {
			if strings.Contains(hay, t) {
				score++
			}
		}
		if score > 0 {
			ranked = append(ranked, scored{k: k, score: score})
		}
	}
	if len(ranked) == 0 {
		return ""
	}
	sort.SliceStable(ranked, func(i, j int) bool { return ranked[i].score > ranked[j].score })

	top := make([]*entity.AIKnowledge, 0, 3)
	for i, r := range ranked {
		if i >= 3 {
			break
		}
		top = append(top, r.k)
	}
	return formatKnowledgeBlock(top)
}

// formatKnowledgeBlock renders already-selected knowledge entries as a reference
// block for the system prompt. Shared by lexical and vector retrieval.
func formatKnowledgeBlock(entries []*entity.AIKnowledge) string {
	if len(entries) == 0 {
		return ""
	}
	var b strings.Builder
	b.WriteString("\n\nReference knowledge (use it to answer accurately; do not invent facts beyond it):\n")
	for _, k := range entries {
		content := k.Content
		if len(content) > 600 {
			content = content[:600] + "…"
		}
		fmt.Fprintf(&b, "- %s: %s\n", k.Title, content)
	}
	return b.String()
}

// queryTerms extracts distinct lowercase tokens (>=4 chars) from the query.
func queryTerms(q string) map[string]struct{} {
	terms := make(map[string]struct{})
	for _, w := range strings.FieldsFunc(strings.ToLower(q), func(r rune) bool {
		return !((r >= 'a' && r <= 'z') || (r >= '0' && r <= '9'))
	}) {
		if len(w) >= 4 {
			terms[w] = struct{}{}
		}
	}
	return terms
}
