<script lang="ts">
	import * as Tooltip from '$lib/components/ui/tooltip/index.js';
	import { cn } from '$lib/utils.js';
	import type { Component } from 'svelte';
	import type { LucideProps } from '@lucide/svelte';

	let {
		icon: Icon,
		label,
		active = false,
		onClick
	}: {
		icon: Component<LucideProps, {}, ''>;
		label: string;
		active?: boolean;
		onClick?: () => void;
	} = $props();
</script>

<Tooltip.Root>
	<Tooltip.Trigger
		class={cn(
			'grid h-10 w-10 place-items-center rounded-lg transition-colors hover:bg-white/10',
			active && 'bg-primary text-white'
		)}
		onclick={onClick}
	>
		<Icon class="h-5 w-5" strokeWidth={1.75} />
	</Tooltip.Trigger>
	<Tooltip.Content side="right">
		<p>{label}</p>
	</Tooltip.Content>
</Tooltip.Root>
