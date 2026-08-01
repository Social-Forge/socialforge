<script lang="ts">
	import { MessageSquareMore } from '@lucide/svelte';
</script>

<div
	class="flex h-full flex-1 flex-col items-center justify-center gap-3 bg-background text-center"
>
	<span class="grid h-16 w-16 place-items-center rounded-full bg-secondary text-muted-foreground">
		<MessageSquareMore class="h-8 w-8" />
	</span>
	<div>
		<p class="font-medium text-foreground">Select a chat to start messaging</p>
		<p class="text-sm text-muted-foreground">Choose from your existing conversations on the left</p>
	</div>
</div>
