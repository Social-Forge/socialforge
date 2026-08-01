<script lang="ts">
	import { DropdownMenu as DM } from 'bits-ui';
	import { cn } from '$lib/utils';

	let {
		items,
		align = 'end',
		trigger
	}: {
		items: {
			label: string;
			onSelect: () => void;
			destructive?: boolean;
			icon?: import('svelte').Snippet;
		}[];
		align?: 'start' | 'end' | 'center';
		trigger: import('svelte').Snippet;
	} = $props();
</script>

<DM.Root>
	<DM.Trigger>
		{@render trigger()}
	</DM.Trigger>
	<DM.Portal>
		<DM.Content
			{align}
			sideOffset={4}
			class="z-50 min-w-40 overflow-hidden rounded-md border border-border bg-popover p-1 text-popover-foreground shadow-md"
		>
			{#each items as item, i (i)}
				<DM.Item
					onSelect={item.onSelect}
					class={cn(
						'flex cursor-pointer items-center gap-2 rounded-sm px-2 py-1.5 text-sm outline-none select-none hover:bg-accent hover:text-accent-foreground',
						item.destructive && 'text-destructive hover:text-destructive'
					)}
				>
					{item.label}
				</DM.Item>
			{/each}
		</DM.Content>
	</DM.Portal>
</DM.Root>

