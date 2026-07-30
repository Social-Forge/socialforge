<script lang="ts">
	import { useModalSub } from '$lib/components/extras/modal/modal.svelte.js';
	import * as Dialog from '$lib/components/extras/dialog/index.js';
	import * as Drawer from '$lib/components/extras/drawer/index.js';
	import type { WithElementRef } from 'bits-ui';
	import type { HTMLAttributes } from 'svelte/elements';

	const modal = useModalSub();

	let {
		ref = $bindable(null),
		children,
		...rest
	}: WithElementRef<HTMLAttributes<HTMLDivElement>> = $props();
</script>

{#if modal.view === 'desktop'}
	<Dialog.Header bind:ref {...rest}>
		{@render children?.()}
	</Dialog.Header>
{:else}
	<Drawer.Header bind:ref {...rest}>
		{@render children?.()}
	</Drawer.Header>
{/if}
