<script lang="ts">
	import ChatHeader from './chat-header.svelte';
	import MessageList from './message-list.svelte';
	import MessageInput from './message-input.svelte';
	import { chatsStore } from '$lib/stores/chats';
	import { subscribeConversation } from '$lib/realtime/chat-realtime';
	import type { ChatMessage, ChatSummary } from './types';

	let { chat }: { chat: ChatSummary } = $props();

	const messages = $derived(chatsStore.messagesFor(chat.id));

	// Lazily load a conversation's messages on first open, then subscribe for
	// realtime updates on its channel.
	$effect(() => {
		const id = chat.id;
		subscribeConversation(id);
		if (chatsStore.hasMessages(id)) return;
		chatsStore.setMessages(id, []); // mark as loading to avoid refetch
		fetch(`/api/chats/${id}/messages`)
			.then((r) => (r.ok ? r.json() : { data: [], meta: {} }))
			.then((res) => chatsStore.setMessages(id, (res.data as ChatMessage[]) ?? []))
			.catch(() => chatsStore.setMessages(id, []));
	});

	async function handleSend(text: string) {
		const trimmed = text.trim();
		if (!trimmed) return;
		const res = await fetch(`/api/chats/${chat.id}/messages`, {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify({ text: trimmed })
		});
		const json = await res.json().catch(() => null);
		const sent = json?.data as ChatMessage | undefined;
		if (sent) chatsStore.appendMessage(sent);
	}
</script>

<div class="flex h-full min-w-0 flex-1 flex-col bg-background">
	<ChatHeader {chat} />
	<MessageList {messages} />
	<MessageInput onsend={handleSend} />
</div>
