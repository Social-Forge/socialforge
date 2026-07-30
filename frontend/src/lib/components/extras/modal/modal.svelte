<script lang="ts">
	import * as Dialog from '$lib/components/extras/dialog/index.js';
	import * as Drawer from '$lib/components/extras/drawer/index.js';
	import type { DialogRootProps } from 'bits-ui';
	import { useModal } from '$lib/components/extras/modal/modal.svelte.js';

	let { open = $bindable(false), children, ...rest }: DialogRootProps = $props();

	const modal = useModal();
</script>

{#if modal.view === 'desktop'}
	<Dialog.Root bind:open {...rest}>
		{@render children?.()}
	</Dialog.Root>
{:else}
	<Drawer.Root bind:open {...rest}>
		{@render children?.()}
	</Drawer.Root>
{/if}
