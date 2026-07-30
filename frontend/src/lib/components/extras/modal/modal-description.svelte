<script lang="ts">
	import { useModalSub } from '$lib/components/extras/modal/modal.svelte.js';
	import * as Dialog from '$lib/components/extras/dialog/index.js';
	import * as Drawer from '$lib/components/extras/drawer/index.js';
	import type { DialogDescriptionProps } from 'bits-ui';

	const modal = useModalSub();

	let { ref = $bindable(null), ...rest }: DialogDescriptionProps = $props();
</script>

{#if modal.view === 'desktop'}
	<Dialog.Description bind:ref {...rest} />
{:else}
	<Drawer.Description bind:ref {...rest} />
{/if}
