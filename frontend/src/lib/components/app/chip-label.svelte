<script lang="ts">
	import { cn } from '$lib/utils';

	export type ChipLabel = {
		id: string;
		label: string;
		color: string;
		textColor?: string;
	};

	let {
		labels = [],
		maxDisplay = 2,
		size = 'sm',
		className
	}: {
		labels: ChipLabel[];
		maxDisplay?: number;
		size?: 'sm' | 'md';
		className?: string;
	} = $props();

	let showAll = $state(false);

	function getColorStyle(color: string): string {
		if (color.startsWith('#')) {
			return `background-color: ${color};`;
		}
		return `background-color: hsl(var(--${color.replace('bg-', '')}));`;
	}
</script>

<div class={cn('relative flex flex-wrap items-center', className)}>
	{#each labels.slice(0, showAll ? labels.length : maxDisplay) as label, index ((label.id, index))}
		<span
			class={cn(
				'relative inline-flex items-center overflow-hidden rounded-l-md px-2 py-1 text-xs font-medium ring-1 ring-background',
				size === 'sm' ? 'text-[10px]' : 'text-xs',
				label.color.startsWith('#') ? '' : `bg-${label.color}`
			)}
			style={`
				${label.color.startsWith('#') ? getColorStyle(label.color) : ''}
				clip-path: polygon(0% 0%, 85% 0%, 100% 50%, 85% 100%, 0% 100%);
				
				/* Kunci untuk tumpang tindih: */
				z-index: ${10 + index};
				margin-left: ${index === 0 ? '0px' : '-8px'};
			`}
		>
		</span>
	{/each}
</div>
