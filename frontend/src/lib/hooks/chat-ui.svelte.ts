import { chatsStore } from '$lib/stores/chats';
import type { ChatSummary } from '$lib/components/app/chat/types.js';

/**
 * Central UI state for the chat screen. Plain Svelte 5 runes module —
 * import `chatUiState` anywhere and mutate it directly.
 *
 * Mobile flow: list -> chat (activeChatId set) -> info panel (infoOpen).
 * Desktop shows list + chat + (optional) info panel side by side.
 */
class ChatUiState {
	activeChatId = $state<string | null>(null);
	infoOpen = $state(false);
	mobileView = $state<'list' | 'chat' | 'info'>('list');

	get activeChat(): ChatSummary | undefined {
		return this.activeChatId ? chatsStore.getConversation(this.activeChatId) : undefined;
	}

	openChat(chatId: string) {
		this.activeChatId = chatId;
		this.mobileView = 'chat';
	}

	backToList() {
		this.mobileView = 'list';
		this.infoOpen = false;
	}

	toggleInfo() {
		this.infoOpen = !this.infoOpen;
		this.mobileView = this.infoOpen ? 'info' : 'chat';
	}
}

export const chatUiState = new ChatUiState();
