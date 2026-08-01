<script lang="ts">
	import * as Tooltip from '$lib/components/ui/tooltip/index.js';
	import { AppIconRail, AppMobileBottomNav } from '$lib/components/app';
	import { chatUiState } from '$lib/hooks/chat-ui.svelte';

	let { children } = $props();
</script>

<Tooltip.Provider>
	<div class="flex h-dvh w-full overflow-hidden bg-background">
		<AppIconRail />

		{@render children()}

		{#if chatUiState.mobileView === 'list'}
			<div class="fixed right-0 bottom-0 left-0 md:hidden">
				<AppMobileBottomNav />
			</div>
		{/if}
	</div>
</Tooltip.Provider>
