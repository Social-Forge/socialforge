<script lang="ts">
	import ChatHeader from './chat-header.svelte';
	import MessageList from './message-list.svelte';
	import MessageInput from './message-input.svelte';
	import { messages as allMessages } from './data';
	import type { ChatSummary } from './types';

	let { chat }: { chat: ChatSummary } = $props();

	const messages = $derived(allMessages.filter((m) => m.chatId === chat.id));
</script>

<div class="flex h-full min-w-0 flex-1 flex-col bg-background">
	<ChatHeader {chat} />
	<MessageList {messages} />
	<MessageInput />
</div>
