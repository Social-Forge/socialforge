<script lang="ts">
	import MessageBubble from './message-bubble.svelte';
	import type { ChatMessage } from './types';

	let { messages }: { messages: ChatMessage[] } = $props();
</script>

<div class="scroll-thin flex-1 overflow-y-auto py-4">
	<div class="mb-4 flex justify-center">
		<span class="rounded-full bg-secondary px-3 py-1 text-xs font-medium text-muted-foreground">
			Today
		</span>
	</div>

	<div class="flex flex-col gap-3">
		{#each messages as message, i (message.id)}
			{@const prev = messages[i - 1]}
			{@const showSender = message.direction === 'in' && prev?.senderName !== message.senderName}
			<MessageBubble {message} {showSender} />
		{/each}
	</div>
</div>
