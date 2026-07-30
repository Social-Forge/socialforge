<script lang="ts" generics="TData">
	import type { Snippet } from 'svelte';
	import type { Table } from '@tanstack/table-core';
	import { Badge } from '@/components/ui/badge';
	import { Button, buttonVariants } from '@/components/ui/button';
	import { Separator } from '@/components/ui/separator';
	import { X } from '@lucide/svelte';
	import { cn } from '@/utils.js';
	import * as i18n from '@/paraglide/messages.js';
	import * as Tooltip from '@/components/ui/tooltip';

	let {
		table = $bindable<Table<TData>>(),
		entityName,
		children
	}: {
		table: Table<TData>;
		entityName: string;
		children?: Snippet;
	} = $props();

	let selectedRows = $derived(table.getSelectedRowModel());
	let selectedCount = $derived(selectedRows.rows.length);
	let toolbarRef = $state<HTMLDivElement | null>(null);
	let announcement = $state<string | undefined>(undefined);

	const handleClearSelection = () => {
		table.resetRowSelection();
	};

	const handleKeyDown = (event: KeyboardEvent) => {
		// if (!toolbarRef) return;

		const buttons = toolbarRef?.querySelectorAll('button');
		if (!buttons) return;

		const currentIndex = Array.from(buttons).findIndex(
			(button) => button === document.activeElement
		);

		switch (event.key) {
			case 'ArrowLeft':
				event.preventDefault();
				const nextIndex = (currentIndex - 1 + buttons.length) % buttons.length;
				if (nextIndex !== currentIndex) {
					buttons[nextIndex].focus();
				}
				break;
			case 'ArrowRight':
				event.preventDefault();
				const prevIndex = currentIndex === 0 ? buttons.length - 1 : currentIndex - 1;
				if (prevIndex !== currentIndex) {
					buttons[prevIndex].focus();
				}
				break;
			case 'Home':
				event.preventDefault();
				if (currentIndex !== 0) {
					buttons[0].focus();
				}
				break;
			case 'End':
				event.preventDefault();
				if (currentIndex !== buttons.length - 1) {
					buttons[buttons.length - 1].focus();
				}
				break;
			case 'Escape': {
				const target = event.target as HTMLElement;
				const activeElement = document.activeElement as HTMLElement;

				const isFromDropdownTrigger =
					target?.getAttribute('data-slot') === 'dropdown-menu-trigger' ||
					activeElement?.getAttribute('data-slot') === 'dropdown-menu-trigger' ||
					target?.closest('[data-slot="dropdown-menu-trigger"]') ||
					activeElement?.closest('[data-slot="dropdown-menu-trigger"]');

				const isFromDropdownContent =
					activeElement?.closest('[data-slot="dropdown-menu-content"]') ||
					target?.closest('[data-slot="dropdown-menu-content"]');

				if (isFromDropdownTrigger || isFromDropdownContent) {
					// Escape was meant for the dropdown - don't clear selection
					return;
				}

				event.preventDefault();
				handleClearSelection();
			}
		}
	};

	$effect(() => {
		if (selectedCount > 0) {
			const message = `${selectedCount} ${entityName} ${selectedCount > 1 ? 's' : ''} selected`;

			queueMicrotask(() => {
				announcement = message;
			});

			const timer = setTimeout(() => (announcement = undefined), 3000);
			return () => clearTimeout(timer);
		}
	});

	$effect(() => {
		if (selectedCount === 0) {
			announcement = undefined;
		}
	});
</script>

<div aria-live="polite" aria-atomic="true" role="status" class="sr-only">
	{announcement}
</div>
<div
	bind:this={toolbarRef}
	role="toolbar"
	aria-label={`Bulk actions for ${selectedCount} selected ${entityName}${selectedCount > 1 ? 's' : ''}`}
	aria-describedby="bulk-actions-description"
	tabindex={-1}
	onkeydown={handleKeyDown}
	class={cn(
		'fixed bottom-6 left-1/2 z-50 -translate-x-1/2 rounded-xl',
		'transition-all delay-100 duration-300 ease-out hover:scale-105',
		'focus-visible:ring-2 focus-visible:ring-ring/50 focus-visible:outline-none'
	)}
>
	<div
		class={cn(
			'p-2 shadow-xl',
			'rounded-xl border',
			'bg-background/95 backdrop-blur-lg supports-backdrop-filter:bg-background/60',
			'flex items-center gap-x-2'
		)}
	>
		<Tooltip.Provider>
			<Tooltip.Root>
				<Tooltip.Trigger
					class={buttonVariants({
						variant: 'outline',
						size: 'icon',
						class: 'size-6 rounded-full'
					})}
					aria-label="Clear selection"
					title="Clear selection (Escape)"
					onclick={handleClearSelection}
				>
					<X class="size-4" />
					<span class="sr-only">Clear selection</span>
				</Tooltip.Trigger>
				<Tooltip.Content>
					<p>Clear selected items</p>
				</Tooltip.Content>
			</Tooltip.Root>
		</Tooltip.Provider>
		<Separator class="h-5" orientation="vertical" aria-hidden="true" />
		<div class="flex items-center gap-x-1 text-sm" id="bulk-actions-description">
			<Badge variant="default" class="min-w-8 rounded-lg" aria-label={`${selectedCount} selected`}>
				{selectedCount}
			</Badge>{' '}
			<span class="hidden sm:inline">
				{entityName}
				{selectedCount > 1 ? 's' : ''}
			</span>{' '}
			Selected items
		</div>

		<Separator class="h-5" orientation="vertical" aria-hidden="true" />

		{@render children?.()}
	</div>
</div>
