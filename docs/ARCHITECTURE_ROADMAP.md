# Social Forge — Architecture Roadmap

> **Produk:** _Social Forge - Multi-agent Customer Service and Omnichannel CRM_
> **Status dokumen:** Disepakati (v1) — acuan eksekusi coding bertahap.
> **Terakhir diperbarui:** 2026-07-31

Dokumen ini adalah **sumber kebenaran (source of truth)** untuk arsitektur & roadmap eksekusi.
Semua keputusan di sini sudah disepakati bersama. Jika ada perubahan arah, update dokumen ini
dulu sebelum mengubah kode.

---

## 1. Ringkasan Produk

Platform CRM omnichannel untuk customer service multi-agent dengan integrasi banyak channel
social media, ditenagai AI Agent untuk meningkatkan **rate closing**. Dipakai untuk kebutuhan
online shop pribadi, lalu dikomersilkan (SaaS multi-tenant) dengan payment gateway.

**Channel yang didukung:**

| Channel | Library / API | Prioritas Build |
|---|---|---|
| WhatsApp (unofficial) | `devlikeapro/waha` (WEBJS, NOWEB, GOWS) | Fase 2 (pertama) |
| Telegram | `telegram-bot-api` | Fase 2 |
| WhatsApp Business | Meta Cloud API (WABA) | Fase 4 |
| Messenger | Meta Graph API | Fase 4 |
| Instagram | Meta Graph API | Fase 4 |
| Webchat | Widget script embed sendiri | Fase 3–5 |
| Linkchat | Share-link (mirip wa.me) | Fase 3 |

**Target aplikasi:**
- **Web** — Landing page + Admin Panel + Chat Portal (SvelteKit v5).
- **Mobile** — React Native (chat only), dibangun setelah API + Web stabil.

---

## 2. Keputusan Arsitektur (DIPATENKAN)

| # | Keputusan | Pilihan Final | Alasan |
|---|---|---|---|
| 1 | **Isolasi multi-tenant** | **Shared-schema + Postgres RLS** | Skala ke ribuan tenant, migrasi sekali. Hapus sisa eksperimen schema-per-tenant. |
| 2 | **RAG / knowledge store AI** | **pgvector** (embedding) + **Typesense** (keyword/typo-tolerant) | pgvector untuk semantic AI, Typesense untuk full-text chat/contact. |
| 3 | **LLM provider** | **OpenRouter** (multi-model) via `AIClient` yang **provider-agnostic** | Fleksibel ganti model per tenant/plan, kontrol biaya untuk komersil. |
| 4 | **Realtime** | **Centrifugo** | Sudah ter-wire; presence + channel per-tenant/conversation/agent. |
| 5 | **Reliability pesan** | **Outbox (outbound) + Inbox/dedup (inbound)** | Tabel `message_outbox` sudah ada; tambah dedup by `provider_message_id`. |
| 6 | **Message broker** | **RabbitMQ** | Antar-worker ingestion, AI job, notifikasi, billing job. |
| 7 | **Object storage** | **MinIO** | Media semua channel + AI asset. |
| 8 | **Auth** | JWT access (±15m) + refresh token + session + OAuth (goth) | Superadmin via middleware role. |
| 9 | **Search** | **Typesense** | Contact, conversation, message full-text; sync via outbox. |
| 10 | **Deployment** | Dev: `docker-compose`. Prod: Ubuntu + Nginx + TLS. | Observability: Prometheus + Loki + Promtail + Grafana. |

---

## 3. Tech Stack (Final)

### Backend
- **Go 1.26** + **Fiber v3**
- **PostgreSQL** (+ **pgvector** extension) + **Typesense**
- **Redis** (cache, session, rate-limit, refresh token, pub/sub tenant refresh)
- **RabbitMQ** (message broker antar-worker)
- **MinIO** (object storage)
- **Centrifugo** (realtime websocket)
- **pgx v5** + **scany** (data access), **goth** (OAuth), **zap** (logging), **golang-jwt/v5**
- **Prometheus client** (metrics)

### Frontend (Web)
- **SvelteKit v5**
- **shadcn-svelte** + **shadcn-svelte-extras** + **TailwindCSS** (Dark/Light mode)
- **i18n** (EN & ID) — sudah ada `project.inlang` + `messages/`
- **Centrifugo JS client** (realtime), Payment checkout page (Midtrans/Xendit/PayPal)

### Mobile
- **React Native** + **NativeWind** + **TailwindCSS**
- Dark/Light mode, **i18n (EN & ID)**, **FCM** push notification
- Chat-only scope

### Infra & Observability
- **Docker** (dev compose + prod), **Nginx** reverse proxy + TLS
- **Grafana + Prometheus + Loki + Promtail**
- Exporters: `postgres-exporter`, `redis-exporter`

### Tambahan yang direkomendasikan (diadopsi)
- **pgvector** — semantic search RAG (ganti kolom `embedding JSONB` → `vector`).
- **OpenRouter** provider di `ai-client` (OpenAI-compatible).
- **FCM** untuk push mobile.
- **CDC ringan** ke Typesense via outbox (bukan trigger berat).

---

## 4. Arsitektur High-Level

```
                         ┌─────────────────────────────────────────────┐
   Social Channels       │                 SOCIAL FORGE                 │
 ┌──────────────┐        │                                             │
 │ WAHA (3 eng) │──webhook──▶ ┌──────────────┐   ┌─────────────────┐   │
 │ Telegram     │──webhook──▶ │  Ingestion   │──▶│  RabbitMQ       │   │
 │ Meta WABA    │──webhook──▶ │  (verify +   │   │  (queues)       │   │
 │ Messenger    │──webhook──▶ │  normalize + │   └────────┬────────┘   │
 │ Instagram    │──webhook──▶ │  dedup)      │            │            │
 │ Webchat      │──ws/http──▶ └──────────────┘            ▼            │
 └──────────────┘        │                        ┌──────────────┐     │
        ▲                │                        │   Workers    │     │
        │                │                        │ - ingest     │     │
        │ outbound       │                        │ - ai-agent   │     │
        │ (message_      │   ┌──────────────┐     │ - auto-assign│     │
        └────outbox)─────┼───│  API (Fiber) │◀───▶│ - billing    │     │
                         │   └──────┬───────┘     │ - sub-expiry │     │
                         │          │             │ - search-sync│     │
                         │   ┌──────▼───────┐     └──────┬───────┘     │
                         │   │  PostgreSQL  │◀───────────┘             │
                         │   │  (RLS +      │     ┌──────────────┐     │
                         │   │   pgvector)  │     │  Centrifugo  │────▶ Web / Mobile
                         │   └──────────────┘     └──────────────┘     │
                         │   Redis · Typesense · MinIO · RabbitMQ      │
                         └─────────────────────────────────────────────┘
```

---

## 5. Multi-Tenancy (Shared-schema + RLS)

**Prinsip:**
- Semua tabel tenant-scoped punya kolom `tenant_id UUID NOT NULL`.
- **RLS aktif** (migrasi `enable_*_rls` sudah ada) — policy pakai session var `app.current_tenant_id`.
- Middleware `tenant.go` set `SET LOCAL app.current_tenant_id = $tenant` **per-transaksi** (bukan per-koneksi, karena pool pgx shared).
- **Superadmin** bypass RLS lewat role Postgres khusus / policy `BYPASSRLS`, di-guard middleware role.

**Yang harus dibereskan di Fase 0:**
- Renumber migrasi supaya konsisten (`users` sebelum `tenant`? tentukan urутan FK yang benar).
- Hapus sisa `user_repository_schema.go` / logika schema-per-tenant.
- Pastikan setiap repository melewati helper yang meng-`SET LOCAL app.current_tenant_id`.
- Verifikasi semua tabel tenant-scoped punya RLS policy (SELECT/INSERT/UPDATE/DELETE).

---

## 6. Ingestion Pipeline & Unified Message Envelope

**Alur inbound (semua channel seragam):**
1. Provider kirim webhook → endpoint `POST /api/webhooks/{provider}/{channel_id}`.
2. **Verify signature** (Meta `X-Hub-Signature-256`, Telegram secret token, WAHA HMAC).
3. **Normalize** ke *Unified Message Envelope* (struct tunggal).
4. **Dedup** by `(channel_id, provider_message_id)` — idempotent (Redis SETNX + unique index).
5. Publish ke RabbitMQ `ingest.inbound`.
6. Worker `ingest`: resolve/kreasi contact + conversation → persist message → update search → publish Centrifugo.

**Unified Message Envelope (konsep):**
```
Envelope {
  provider        // waha | telegram | meta_wa | messenger | instagram | webchat
  channel_id
  provider_msg_id
  direction       // inbound | outbound
  contact         { external_id, name, avatar, phone/username }
  type            // text | image | video | audio | document | location | template | sticker | call | system
  text
  media[]         { mime, url(minio), size, thumb, caption }
  reply_to        // provider_msg_id yang di-reply
  timestamp
  raw             // payload asli (audit/debug)
}
```

**Alur outbound:**
Agent/AI kirim → tulis ke `message_outbox` (status `pending`) → worker `dispatch` ambil →
kirim ke provider (rate-limited per channel) → update status `sent`/`failed` (retry backoff) →
publish status ke Centrifugo (double-check icon).

---

## 7. Realtime (Centrifugo)

**Channel naming:**
- `tenant:{tenant_id}` — event global tenant (badge, notif).
- `conversation:{conversation_id}` — pesan & status di 1 room.
- `agent:{user_id}` — assign, notif personal.
- `presence:tenant:{tenant_id}` — agent online/offline (untuk auto-assign & working hours).

**Auth:** JWT connection token dari backend; subscribe permission dicek per role/assignment.

---

## 8. AI Agent (RAG + Persona + Guardrails)

**Komponen (entity sudah ada):**
- `ai_agents` — identitas: nama, soul, system prompt, tone, gender, gaya, safety instructions.
- `ai_playbooks` — aturan "kalau pelanggan tanya X / ada kata Y → jawab Z, dalam lingkup ...".
- `ai_knowledge` — knowledge base, **embedding → migrasi ke `vector` (pgvector)**.
- `ai_assets` — referensi file gambar/video/testimoni (di MinIO) untuk dikirim ke customer.
- `ai_credit_ledger` — metering token (debit/credit) per tenant.

**Alur jawab AI:**
1. Pesan masuk ke channel yang AI-nya aktif → override auto-first-reply.
2. Retrieve: playbook match (keyword/rule) + top-K knowledge (pgvector cosine).
3. Compose prompt: persona + system + guardrails + retrieved context + history.
4. Call **OpenRouter** (model per plan/tenant) via `AIClient` (provider-agnostic).
5. Post-process: guardrails filter, attach `ai_assets` bila playbook mensyaratkan.
6. Kirim balasan lewat outbox + **debit `ai_credit_ledger`** sesuai token usage.
7. Bila topik sensitif / low-confidence → handoff ke human agent (unassign AI → antre agent).

**Provider-agnostic client:** `openai.go`, `anthropic.go`, `gemini.go` sudah ada; **tambah `openrouter.go`** (OpenAI-compatible) + `embeddings.go` (embedding model untuk RAG).

---

## 9. Channels — Detail

### 9.1 WAHA (WhatsApp unofficial)
- **3 engine** (compose sudah ada): **GOWS** (default, ringan/stabil) → **NOWEB** (no-browser, hemat RAM) → **WEBJS** (fitur terlengkap, fallback).
- Session di **Postgres**, media auto-download ke **MinIO**.
- **Auto-reject call**: subscribe event `call.received` → reject + kirim auto-response (bila enabled).
- Routing session→engine disimpan di `channel` (metadata engine).

### 9.2 Telegram
- Bot API webhook, secret token verification. Paling cepat diintegrasi (Fase 2).

### 9.3 Meta (WABA / Messenger / Instagram)
- Graph API + webhook verify (`hub.challenge` + `X-Hub-Signature-256`).
- Token page/app disimpan **terenkripsi** (`secret_helper` sudah ada).
- WABA: template message (HSM) untuk pesan di luar 24h window.

### 9.4 Webchat
- Widget floating bubble → script embed → koneksi via Centrifugo/WS → jadi channel `webchat`.

### 9.5 Linkchat
- Short-link share (mirip wa.me) menuju group/division; landing → mulai conversation.

---

## 10. Billing & Subscription

**Entity sudah ada:** `plans`, `subscriptions`, `invoices`, `payment_events`, `subscription_addons`.

- **Plan enforcement** via **quota middleware**: cek limit (channel/agent/quick-reply/AI token)
  sebelum operasi create. Limit ditarik dari plan + addon aktif.
- **Gateway:** Midtrans (ID card/ewallet), Xendit (VA/QRIS), PayPal (internasional).
- **Alur:** create invoice → redirect/checkout → webhook gateway → verify signature →
  `payment_events` → aktivasi/perpanjang subscription/addon → notif realtime.
- **Worker sub-expiry:** cron cek subscription kadaluarsa → downgrade/suspend + notif.
- **AI token:** default 0, top-up via addon → tambah credit ke `ai_credit_ledger`.

**Default limit plan (dari spec):**

| Fitur | Free | Pro |
|---|---|---|
| WhatsApp WAHA | 0 | 1 |
| WhatsApp Business (Meta) | 0 | 1 |
| Messenger | 1 | 10 |
| Instagram | 1 | 10 |
| Telegram | 1 | 10 |
| Agent CS (termasuk supervisor) | 1 | 10 |
| Quick reply | (default kecil) | 100 |
| Linkchat & Webchat | terbatas | Unlimited |
| AI Agent integration | 0 | 1 |
| AI token | 0 (top-up addon) | 0 (top-up addon) |

---

## 11. Model Data (existing — 34 migrations)

Sudah ada & dipakai sebagai fondasi (tidak dibangun ulang, hanya disempurnakan):

`users`, `roles`, `tenants`, `user_tenant`, `divisions`, `division_members`,
`oauth_providers`, `ai_agents`, `channel`, `contact`, `conversation`, `message`,
`message_outbox`, `conversation_events`, `ai_credit_ledger`, `ai_knowledge`,
`labels`, `conversation_labels`, `quick_replies`, `plans`, `subscriptions`,
`invoices`, `payment_events`, `subscription_addons`, `ai_playbooks`, `ai_assets`,
`audit_logs` + RLS untuk tenant/messaging/ai/label/quick_reply/billing/ai_advanced.

### Tabel BARU yang perlu ditambah (roadmap)
- `conversation_metrics` — SLA (first response time, resolution time) & **CSAT**.
- `agent_working_hours` — jam kerja per agent/division.
- `auto_responses` — first-reply per channel (text/media/hybrid).
- `auto_assign_rules` — tipe (round-robin, percentage, ...) per division/channel.
- `webhook_events` — log & dedup inbound (idempotency).
- `push_tokens` — FCM token device (mobile).
- (opsi) `blocked_contacts` bila belum tercakup di `contact`.

---

## 12. Struktur Direktori (konvensi yang dijaga)

```
backend/internal/
  app/            # bootstrap fiber + route registration
  dependencies/   # DI container
  factory/        # wiring handler↔service↔repo per domain
  handlers/       # HTTP handlers (thin)
  services/       # business logic
  infra/
    repository/   # data access (pgx + scany), RLS-aware
    channels/     # waha | meta | telegram | webchat | linkchat  ← DIISI
    ai-client/    # openai|anthropic|gemini|openrouter + embeddings ← DIISI
    centrifugo/ minio-client/ redis-client/ typesense-client/ oauth/
  middlewares/    # auth, tenant, permission, rate_limiter, csrf, flag, ...
  entity/  dto/  helpers/  utils/
cmd/
  api/            # HTTP server
  worker/         # RabbitMQ consumers + cron  ← DIISI
```

Pola per domain baru: `entity → migration → repository → service → handler → factory → route`.

---

## 13. Roadmap Eksekusi (8 Fase)

Tiap fase = increment yang **jalan & bisa dites**. Selesai tiap fase → commit + verifikasi build.

### Fase 0 — Cleanup & Fondasi ⬅️ **SEDANG DIKERJAKAN**
- [x] Migrasi identity chain konsisten (urutan FK benar: users→roles→tenant→user_tenants→divisions→division_members).
- [x] Fix typo migrasi `000030` (`subssubscription_addonscriptions` → `subscription_addons`).
- [x] Kunci multi-tenant RLS — GUC `app.current_tenant` distandarkan; **fondasi RLS dibuat**: `repository.WithTenantID` + `RunInTenantTx` + `q(ctx)` (`set_config` transaction-scoped). `TenantGuard` inject tenant ke context. Tak ada sisa schema-per-tenant.
- [x] Normalisasi identity: `is_active`→`status`, `is_verified`→`email_verified_at` (derived methods `User.IsActive()/IsVerified()`). Hapus RBAC permission (level/name-based saja).
- [x] Perbaiki rantai identity: `UserTenantWithDetails` + `UserResponse` aggregate, `GetUserTenantWithDetailsByUserID`, **Register atomik + buat `user_tenants` membership** (sebelumnya putus). Hapus 3 method repo dead+broken (`u.is_active`).
- [x] `division_repository` diperbaiki total (INSERT arg, `owner_id` hantu, `$?` mismatch, scany, RLS-aware). `tenant_helper` diperbaiki.
- [x] **Build hijau (`go build ./...`) + `go vet` clean.**
- [ ] `docker-compose` sanity: service `backend` + `worker` + `frontend` + `nginx` lengkap. _(berikutnya)_
- [ ] `migrate up/down` mulus di DB nyata _(perlu dijalankan user via tooling migrate)_.

> **Dipindah ke fase pemakainya** (hindari skema spekulatif tanpa konsumen):
> - `ai_knowledge.embedding` JSONB → **pgvector `vector(N)`** → **Fase 5** (dimensi N mengikuti model embedding OpenRouter yang dipilih; butuh dependency `pgvector-go`).
> - Tabel `webhook_events` (dedup) → **Fase 2**; `auto_responses` → **Fase 3**; `auto_assign_rules` → **Fase 3**.
> - RLS untuk `contacts` (PII) → **Fase 2** (saat contact dibangun).

### Fase 1 — Auth & Tenant Lengkap
- [ ] Signup, Signin, Verify Email, Forgot/Reset, Confirm + Social OAuth (signin/signup).
- [ ] Refresh token rotation + session revoke; rate-limit + anti-abuse di auth.
- [ ] Superadmin panel (middleware role) + tenant provisioning (default Free plan).
- [ ] Division CRUD, member (Supervisor/Agent), channel-assign, quota middleware dasar.

### Fase 2 — Channel & Ingestion (WAHA + Telegram dulu)
- [ ] Unified Message Envelope + normalizer.
- [ ] Webhook endpoints + signature verify + dedup (`webhook_events`).
- [ ] WAHA adapter (3 engine, session Postgres, media→MinIO, auto-reject call).
- [ ] Telegram adapter.
- [ ] Worker `ingest` + `dispatch` (outbox) + RabbitMQ queues.
- [ ] Contact auto-create + label channel.

### Fase 3 — Conversation & Chat Core
- [ ] Conversation room, assign engine (round-robin/percentage), manual transfer/unassign.
- [ ] Label management (unique per tenant), quick-reply (type text/media/hybrid + upload).
- [ ] Realtime via Centrifugo (pesan, status, unread badge, presence).
- [ ] Archived, pinned conversation/message, filters, search (Typesense).
- [ ] Auto-response first-reply per channel + working hours.
- [ ] Linkchat + Webchat widget dasar.

### Fase 4 — Meta Channels
- [ ] WABA (template/HSM), Messenger, Instagram adapters + verify + token terenkripsi.

### Fase 5 — AI Agent
- [ ] `openrouter.go` + `embeddings.go` (provider-agnostic).
- [ ] Playbook engine + knowledge RAG (pgvector) + persona + guardrails.
- [ ] Override auto-first-reply, handoff ke human, metering `ai_credit_ledger`.
- [ ] Webchat AI chatbot.

### Fase 6 — Billing & Subscription
- [ ] Plan enforcement penuh (semua quota), addon top-up (slot & AI token).
- [ ] Midtrans + Xendit + PayPal (checkout + webhook verify + `payment_events`).
- [ ] Worker sub-expiry + invoice + notif realtime.

### Fase 7 — Frontend & Mobile
- [ ] SvelteKit: Landing (i18n EN/ID), Auth pages, App shell (sidebar WhatsApp-style),
      Chat portal, Analytics, Contact, AI Agent, Integrations, Settings, Superadmin.
- [ ] Analytics: General, Agent performance, AI analysis, SLA & CSAT, Contact.
- [ ] Payment checkout invoice page (realtime webhook status).
- [ ] React Native chat app (drag gestures, i18n, dark/light, FCM).

---

## 14. Cross-Cutting Concerns (berlaku semua fase)

- **Idempotency** di endpoint mutasi + inbound webhook (dedup).
- **Signature verification** tiap provider.
- **Per-channel outbound rate-limit** (hindari ban WA/Meta).
- **Audit log** (`audit_logs`) untuk aksi sensitif.
- **Feature flags** (`flag.go`) untuk rollout bertahap.
- **Observability**: metrics Prometheus, log terstruktur Loki, dashboard Grafana.
- **Secret encryption** (`secret_helper`) untuk token channel/gateway.
- **Graceful shutdown** (sudah ada di container).

---

## 15. Risiko & Catatan

- **RLS + pooling**: wajib `SET LOCAL` di dalam transaksi, jangan `SET` di level koneksi.
- **WAHA stabilitas**: engine unofficial rawan; sediakan health-check + reconnect + failover engine.
- **Meta review**: WABA/IG butuh app review & template approval — siapkan lebih awal untuk komersil.
- **Biaya AI**: metering ketat via ledger; set hard-limit per plan agar tidak boncos.
- **Migrasi**: setelah Fase 0, migrasi bersifat **append-only** (jangan renumber lagi).

---

## 16. Definition of Done (per fase)

Sebuah fase dianggap selesai bila: (1) build hijau, (2) migrasi `up/down` mulus,
(3) endpoint/worker terkait berjalan & bisa dites manual/otomatis, (4) tidak ada mismatch
entity↔migrasi baru, (5) dokumen ini diperbarui bila ada perubahan arsitektur.
