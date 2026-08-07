<script lang="ts">
	import { onMount } from 'svelte';
	import { fly } from 'svelte/transition';
	import {
		AppChatListPanel,
		AppChatWindow,
		AppEmptyChatState,
		AppInfoPanel
	} from '$lib/components/app/chat';
	import { chatUiState } from '$lib/hooks/chat-ui.svelte';
	import { chatsStore } from '$lib/stores/chats';
	import { connectRealtime, disconnectRealtime } from '$lib/realtime/chat-realtime';
	import type { ChatSummary } from '$lib/components/app/chat/types';

	let { data } = $props();

	// Seed the store with real conversations from the server load.
	// svelte-ignore state_referenced_locally
	chatsStore.setConversations((data.conversations ?? []) as ChatSummary[]);

	onMount(() => {
		if (data.realtime?.token && data.realtime?.tenantId) {
			connectRealtime({
				token: data.realtime.token,
				wsUrl: data.realtime.wsUrl,
				tenantId: data.realtime.tenantId
			});
		}
		return () => disconnectRealtime();
	});
</script>

<!-- Chat list: always visible on desktop; on mobile only in 'list' view -->
<div
	class={chatUiState.mobileView === 'list'
		? 'flex h-full w-full flex-col pb-16 md:flex md:w-auto md:pb-0'
		: 'hidden md:flex md:w-auto'}
>
	<AppChatListPanel />
</div>

<!-- Chat window: always visible on desktop when a chat is active; on mobile only in 'chat' view -->
<div
	class={chatUiState.mobileView === 'chat'
		? 'flex h-full w-full min-w-0 flex-1 md:flex'
		: 'hidden md:flex md:min-w-0 md:flex-1'}
>
	{#if chatUiState.activeChat}
		<AppChatWindow chat={chatUiState.activeChat} />
	{:else}
		<AppEmptyChatState />
	{/if}
</div>

<!-- Info panel: shown alongside on desktop when toggled; full-screen on mobile in 'info' view -->

{#if chatUiState.infoOpen}
	<div
		transition:fly={{ x: 250, duration: 300 }}
		class={chatUiState.mobileView === 'info'
			? 'flex h-full w-full md:flex md:w-auto'
			: 'hidden md:flex md:w-auto'}
	>
		<AppInfoPanel />
	</div>
{/if}
