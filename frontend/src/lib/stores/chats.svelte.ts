import type { ChatMessage, ChatSummary } from '$lib/components/app/chat/types';

/**
 * Central chats data store (Svelte 5 runes — `.svelte.ts` so runes compile).
 * Seeded from server load (`data.conversations`) and kept live by the
 * Centrifugo client (`lib/realtime/chat-realtime`).
 */
class ChatsStore {
	conversations = $state<ChatSummary[]>([]);
	/** Loaded messages per conversation id (oldest-first). */
	messagesByChat = $state<Record<string, ChatMessage[]>>({});

	get bookmarked(): ChatSummary[] {
		return this.conversations.filter((c) => c.pinned);
	}
	get rest(): ChatSummary[] {
		return this.conversations.filter((c) => !c.pinned);
	}

	setConversations(list: ChatSummary[]) {
		this.conversations = list;
	}

	getConversation(id: string): ChatSummary | undefined {
		return this.conversations.find((c) => c.id === id);
	}

	setMessages(chatId: string, msgs: ChatMessage[]) {
		this.messagesByChat[chatId] = msgs;
	}

	messagesFor(chatId: string): ChatMessage[] {
		return this.messagesByChat[chatId] ?? [];
	}

	hasMessages(chatId: string): boolean {
		return this.messagesByChat[chatId] != null;
	}

	/** Append a message (dedup by id) + refresh the conversation preview. */
	appendMessage(msg: ChatMessage, opts: { markUnread?: boolean } = {}) {
		const existing = this.messagesByChat[msg.chatId] ?? [];
		if (existing.some((m) => m.id === msg.id)) return;
		this.messagesByChat[msg.chatId] = [...existing, msg];
		this.touchConversation(msg.chatId, msg.text ?? '', msg.timestamp, opts.markUnread);
	}

	/** Move a conversation to the top and update its preview/time. */
	touchConversation(chatId: string, preview: string, at: string, markUnread = false) {
		const idx = this.conversations.findIndex((c) => c.id === chatId);
		if (idx === -1) return;
		const conv = this.conversations[idx];
		conv.lastMessagePreview = preview || conv.lastMessagePreview;
		conv.lastMessageAt = at || conv.lastMessageAt;
		if (markUnread) conv.unreadCount = (conv.unreadCount ?? 0) + 1;
		this.conversations = [...this.conversations].sort((a, b) => {
			if (!!a.pinned !== !!b.pinned) return a.pinned ? -1 : 1;
			return new Date(b.lastMessageAt).getTime() - new Date(a.lastMessageAt).getTime();
		});
	}

	/** Apply a `tenant:` channel conversation-updated event (list refresh). */
	applyConversationUpdate(evt: {
		conversation_id: string;
		preview?: string;
		last_message_at?: string;
	}) {
		this.touchConversation(
			evt.conversation_id,
			evt.preview ?? '',
			evt.last_message_at ?? new Date().toISOString(),
			true
		);
	}

	clearUnread(chatId: string) {
		const conv = this.getConversation(chatId);
		if (conv) conv.unreadCount = undefined;
	}
}

export const chatsStore = new ChatsStore();
