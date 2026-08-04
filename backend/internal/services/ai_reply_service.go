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
	aiAgentRepo repository.AIAgentRepository
	messageRepo repository.MessageRepository
	creditRepo  repository.AICreditLedgerRepository
	aiClient    *aiclient.AIClient
	outbound    *OutboundService
	logger      *zap.Logger
}

func NewAIReplyService(
	aiAgentRepo repository.AIAgentRepository,
	messageRepo repository.MessageRepository,
	creditRepo repository.AICreditLedgerRepository,
	aiClient *aiclient.AIClient,
	outbound *OutboundService,
	logger *zap.Logger,
) *AIReplyService {
	return &AIReplyService{
		aiAgentRepo: aiAgentRepo,
		messageRepo: messageRepo,
		creditRepo:  creditRepo,
		aiClient:    aiClient,
		outbound:    outbound,
		logger:      logger,
	}
}

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

	// 1. Load the agent + its recent conversation history in one tenant tx.
	var agent *entity.AIAgent
	var history []*entity.Message
	var balance int
	err := s.aiAgentRepo.RunInTenantTx(tctx, func(txCtx context.Context) error {
		var err error
		if agent, err = s.aiAgentRepo.FindByID(txCtx, *channel.AIAgentID); err != nil {
			return err
		}
		if history, err = s.messageRepo.ListByConversation(txCtx, conversationID, aiHistoryLimit); err != nil {
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

	// 3. Call the model.
	systemPrompt := buildSystemPrompt(agent)
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

	s.logger.Info("ai reply sent",
		zap.String("agent", agent.Name),
		zap.String("model", resp.Model),
		zap.Int("tokens", tokens),
		zap.Int("balance_after", entry.BalanceAfter),
	)
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
