<script lang="ts">
	import type { Snippet } from 'svelte';
	import { useRename } from '$lib/components/extras/rename/rename.svelte.js';

	useRename();

	let { children }: { children: Snippet } = $props();
</script>

{@render children?.()}
