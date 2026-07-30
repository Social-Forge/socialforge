<script lang="ts">
	import { useModalSub } from '$lib/components/extras/modal/modal.svelte.js';
	import * as Dialog from '$lib/components/extras/dialog/index.js';
	import * as Drawer from '$lib/components/extras/drawer/index.js';
	import type { DialogTitleProps } from 'bits-ui';

	const modal = useModalSub();

	let { ref = $bindable(null), ...rest }: DialogTitleProps = $props();
</script>

{#if modal.view === 'desktop'}
	<Dialog.Title bind:ref {...rest} />
{:else}
	<Drawer.Title bind:ref {...rest} />
{/if}
