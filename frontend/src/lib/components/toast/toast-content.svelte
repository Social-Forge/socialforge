<script lang="ts">
	import Toast from './toast-item.svelte';
	import { toast } from '$lib/stores/toast';
</script>

{#each $toast as item (item.id)}
	<Toast
		message={item.message}
		type={item.type}
		duration={item.duration ?? 0}
		onDismiss={() => toast.remove(item.id)}
	/>
{/each}

