import { BaseHandler } from './base';
import type { ChatSummary, ChatMessage, MessageStatus, ChipLabel } from '$lib/components/app/chat/types';

/** Raw conversation list item returned by GET /conversations/protected (enriched). */
interface ConversationListItem {
	id: string;
	channel_id: string;
	contact_id: string;
	assigned_agent_id: string | null;
	status: string;
	is_pinned: boolean;
	is_archived: boolean;
	unread_count: number;
	last_message_at: string | null;
	updated_at: string;
	contact_name: string;
	contact_avatar: string | null;
	contact_external_id: string;
	channel_type: string;
	channel_name: string | null;
	agent_name: string | null;
	last_message_body: string | null;
	last_message_direction: string | null;
	last_message_sender_type: string | null;
	labels: ChipLabel[] | null;
}

/** Raw message row returned by GET /conversations/protected/:id/messages. */
interface BackendMessage {
	id: string;
	conversation_id: string;
	direction: string; // 'in' | 'out'
	sender_type: string; // contact | agent | ai | system
	content_type: string;
	body: string | null;
	media: Record<string, unknown> | null;
	status: string; // pending | sent | delivered | read | failed
	created_at: string;
	reply_to_id?: string | null;
}

const CHANNEL_MAP: Record<string, ChatSummary['channel']> = {
	whatsapp_waha: 'whatsapp_waha',
	whatsapp_meta: 'whatsapp_meta',
	messenger: 'messenger',
	instagram: 'instagram',
	telegram: 'telegram'
};

function mapStatus(s: string | null | undefined): MessageStatus | undefined {
	if (!s) return undefined;
	return s === 'pending' ? 'sending' : (s as MessageStatus);
}

function mapConversation(it: ConversationListItem): ChatSummary {
	return {
		id: it.id,
		name: it.contact_name || it.contact_external_id || 'Unknown',
		avatarUrl: it.contact_avatar || undefined,
		lastMessagePreview: it.last_message_body || '',
		lastMessageAt: it.last_message_at || it.updated_at,
		unreadCount: it.unread_count > 0 ? it.unread_count : undefined,
		pinned: it.is_pinned || undefined,
		channel: CHANNEL_MAP[it.channel_type] ?? (it.channel_type as ChatSummary['channel']),
		agentName: it.agent_name || undefined,
		labels: it.labels && it.labels.length > 0 ? it.labels : undefined
	};
}

function mapMessage(m: BackendMessage): ChatMessage {
	return {
		id: m.id,
		chatId: m.conversation_id,
		direction: (m.direction === 'in' ? 'in' : 'out'),
		timestamp: m.created_at,
		status: m.direction === 'out' ? mapStatus(m.status) : undefined,
		text: m.body || undefined
	};
}

export interface ConversationListFilters {
	status?: string;
	search?: string;
	channel_id?: string;
	archived?: boolean;
}

export class ConversationHandler extends BaseHandler {
	/** List conversations mapped to the chat-list UI shape. */
	async list(filters: ConversationListFilters = {}): Promise<ChatSummary[]> {
		const qs = new URLSearchParams();
		if (filters.status) qs.set('status', filters.status);
		if (filters.search) qs.set('search', filters.search);
		if (filters.channel_id) qs.set('channel_id', filters.channel_id);
		if (filters.archived != null) qs.set('archived', String(filters.archived));
		const suffix = qs.toString() ? `?${qs}` : '';
		const res = await this.api.authRequest<ConversationListItem[]>(
			'GET',
			`/conversations/protected/${suffix}`
		);
		if (!res.success || !res.data) return [];
		return res.data.map(mapConversation);
	}

	/** Messages for a conversation, oldest-first (backend returns newest-first). */
	async messages(conversationId: string): Promise<ChatMessage[]> {
		const res = await this.api.authRequest<BackendMessage[]>(
			'GET',
			`/conversations/protected/${conversationId}/messages`
		);
		if (!res.success || !res.data) return [];
		return res.data.map(mapMessage).reverse();
	}

	/** Send an agent text reply; returns the mapped message on success. */
	async send(conversationId: string, text: string): Promise<ChatMessage | null> {
		const res = await this.api.authRequest<BackendMessage>(
			'POST',
			`/conversations/protected/${conversationId}/messages`,
			{ text }
		);
		if (!res.success || !res.data) return null;
		return mapMessage(res.data);
	}

	/** Mark a conversation read. */
	async markRead(conversationId: string) {
		return this.api.authRequest('POST', `/conversations/protected/${conversationId}/read`);
	}

	/** Fetch a short-lived Centrifugo connection token + ws url for realtime. */
	async realtimeToken(): Promise<{ token: string; wsUrl: string; tenantId: string } | null> {
		const res = await this.api.authRequest<{ token: string; ws_url: string; tenant_id: string }>(
			'GET',
			'/token/centrifugo'
		);
		if (!res.success || !res.data?.token) return null;
		return { token: res.data.token, wsUrl: res.data.ws_url, tenantId: res.data.tenant_id };
	}
}
