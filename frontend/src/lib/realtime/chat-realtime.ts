import { chatsStore } from '$lib/stores/chats';
import type { ChatMessage } from '$lib/components/app/chat/types';
import { chatUiState } from '$lib/hooks/chat-ui.svelte';

/**
 * Minimal Centrifugo bidirectional JSON-protocol client over a raw WebSocket
 * (avoids the `centrifuge` npm dependency). Connects with the backend-issued
 * connection token, subscribes to the tenant channel (list updates) and to
 * per-conversation channels (new messages), routing publications into the
 * chats store. The `conversation`/`tenant` namespaces have
 * allow_subscribe_for_client=true, so no per-channel subscription tokens.
 */
let ws: WebSocket | null = null;
let cmdId = 0;
let connected = false;
let cfg: { token: string; wsUrl: string; tenantId: string } | null = null;
const subscribed = new Set<string>();
const pending = new Set<string>();
let reconnectTimer: ReturnType<typeof setTimeout> | null = null;

function nextId() {
	return ++cmdId;
}

function sendCommand(obj: Record<string, unknown>) {
	if (ws && ws.readyState === WebSocket.OPEN) {
		ws.send(JSON.stringify(obj));
	}
}

function subscribeChannel(channel: string) {
	if (subscribed.has(channel)) return;
	if (!connected) {
		pending.add(channel);
		return;
	}
	subscribed.add(channel);
	sendCommand({ subscribe: { channel }, id: nextId() });
}

function mapRealtimeMessage(d: Record<string, unknown>): ChatMessage {
	const direction = d.direction === 'in' ? 'in' : 'out';
	return {
		id: String(d.id),
		chatId: String(d.conversation_id),
		direction,
		timestamp: String(d.created_at ?? new Date().toISOString()),
		status: direction === 'out' ? ((d.status as ChatMessage['status']) ?? 'sent') : undefined,
		text: (d.body as string) || undefined
	};
}

function handlePublication(channel: string, data: Record<string, unknown>) {
	if (channel.startsWith('tenant:')) {
		if (data.type === 'conversation_updated') {
			chatsStore.applyConversationUpdate({
				conversation_id: String(data.conversation_id),
				preview: data.preview as string | undefined,
				last_message_at: data.last_message_at as string | undefined
			});
		}
		return;
	}
	if (channel.startsWith('conversation:')) {
		// New-message publications are wrapped: { type: 'new_message', message: {...} }.
		const m = (data.type === 'new_message' ? data.message : data) as Record<string, unknown>;
		if (m?.id && m?.conversation_id) {
			const msg = mapRealtimeMessage(m);
			const isActive = chatUiState.activeChatId === msg.chatId;
			chatsStore.appendMessage(msg, { markUnread: !isActive && msg.direction === 'in' });
		}
	}
}

function handleFrame(raw: string) {
	// Centrifugo may batch multiple newline-separated JSON commands.
	for (const line of raw.split('\n')) {
		if (!line.trim()) continue;
		let msg: any;
		try {
			msg = JSON.parse(line);
		} catch {
			continue;
		}
		// Server ping is an empty object; reply with an empty command.
		if (msg && Object.keys(msg).length === 0) {
			sendCommand({});
			continue;
		}
		if (msg.push?.pub?.data) {
			handlePublication(msg.push.channel, msg.push.pub.data);
			continue;
		}
		// Reply to our connect command (id 1) -> flush queued subscriptions.
		if (msg.connect || msg.id === 1) {
			connected = true;
			for (const ch of pending) {
				pending.delete(ch);
				subscribeChannel(ch);
			}
		}
	}
}

export function connectRealtime(opts: { token: string; wsUrl: string; tenantId: string }) {
	if (typeof window === 'undefined') return;
	cfg = opts;
	openSocket();
	// Tenant channel: conversation-list updates.
	subscribeChannel(`tenant:${opts.tenantId}`);
}

function openSocket() {
	if (!cfg || (ws && ws.readyState <= WebSocket.OPEN)) return;
	try {
		ws = new WebSocket(cfg.wsUrl);
	} catch {
		scheduleReconnect();
		return;
	}
	ws.onopen = () => {
		connected = false;
		cmdId = 0;
		sendCommand({ connect: { token: cfg!.token }, id: nextId() });
	};
	ws.onmessage = (e) => handleFrame(typeof e.data === 'string' ? e.data : '');
	ws.onclose = () => {
		connected = false;
		subscribed.clear();
		scheduleReconnect();
	};
	ws.onerror = () => ws?.close();
}

function scheduleReconnect() {
	if (reconnectTimer || !cfg) return;
	reconnectTimer = setTimeout(() => {
		reconnectTimer = null;
		// Re-queue known channels so they resubscribe after reconnect.
		openSocket();
	}, 3000);
}

/** Subscribe to a conversation's channel for realtime new messages. */
export function subscribeConversation(conversationId: string) {
	if (typeof window === 'undefined' || !cfg) return;
	subscribeChannel(`conversation:${conversationId}`);
}

export function disconnectRealtime() {
	if (reconnectTimer) clearTimeout(reconnectTimer);
	reconnectTimer = null;
	cfg = null;
	connected = false;
	subscribed.clear();
	pending.clear();
	ws?.close();
	ws = null;
}
