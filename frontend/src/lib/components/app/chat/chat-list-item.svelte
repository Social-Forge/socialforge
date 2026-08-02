<script lang="ts">
	import { AppAvatar, AppChipLabel } from '$lib/components/app';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu/index.js';
	import { buttonVariants } from '$lib/components/ui/button/index.js';
	import { formatTime } from '$lib/utils/time';
	import { Check, CheckCheck, EllipsisVertical } from '@lucide/svelte';
	import type { ChatSummary } from './types';
	import { cn } from '$lib/utils';
	import { IsMobile } from '$lib/hooks/is-mobile.svelte';
	import { longpress } from '$lib/hooks/use-long-press.svelte';

	let {
		chat,
		active = false,
		onclick
	}: { chat: ChatSummary; active?: boolean; onclick: () => void } = $props();

	let isHovered = $state(false);
	let isDropdownOpen = $state(false);
	let isMobile = $derived(new IsMobile().current);

	function handleLongPress() {
		if (isMobile) {
			isDropdownOpen = true;
			isHovered = true;
			if (navigator.vibrate) {
				navigator.vibrate(50);
			}
		}
	}

	function handleShortPress() {
		if (!isDropdownOpen) {
			onclick?.();
		}
	}

	function handleDropdownOpenChange(open: boolean) {
		isDropdownOpen = open;
		if (!open) {
			isHovered = false;
		}
		// if (!open && !isHovered) {
		// 	isHovered = false;
		// }
	}

	function handleMenuItemClick(action: string) {
		console.log(`Action: ${action} for chat ${chat.id}`);
		isDropdownOpen = false;
		isHovered = false;
	}
</script>

<div
	role="button"
	tabindex="0"
	use:longpress={500}
	onlongpress={handleLongPress}
	onshortpress={handleShortPress}
	class={cn(
		'conversation-item-wrapper relative w-full rounded-lg transition-colors hover:bg-accent',
		active && 'bg-accent'
	)}
	onmouseenter={() => {
		if (!isMobile) isHovered = true;
	}}
	onmouseleave={() => {
		if (!isMobile && !isDropdownOpen) isHovered = false;
	}}
>
	<button class="flex w-full items-center gap-3 rounded-lg px-3 py-2.5 text-left">
		<AppAvatar
			src={chat.avatarUrl}
			fallback={chat.name.slice(0, 1)}
			channel={chat.channel}
			size="md"
		/>
		<div class="min-w-0 flex-1">
			<div class="relative flex items-center justify-between gap-2">
				<div class="flex min-w-0 items-center gap-0.5">
					<span class="truncate text-sm font-semibold text-foreground">{chat.name}</span>
					{#if chat.labels?.length}
						<AppChipLabel labels={chat.labels} maxDisplay={2} size="sm" />
					{/if}
				</div>
				{#if isHovered || isDropdownOpen}
					<DropdownMenu.Root open={isDropdownOpen} onOpenChange={handleDropdownOpenChange}>
						<DropdownMenu.Trigger
							class={buttonVariants({ variant: 'ghost', size: 'icon-sm' })}
							onclick={(e) => {
								e.stopPropagation();
								isDropdownOpen = !isDropdownOpen;
							}}
						>
							<EllipsisVertical class="h-4 w-4" />
						</DropdownMenu.Trigger>
						<DropdownMenu.Content class="w-56" align="start">
							<DropdownMenu.Group>
								<DropdownMenu.Item onSelect={() => handleMenuItemClick('pinned')}>
									Pin Chat
								</DropdownMenu.Item>
								<DropdownMenu.Item onSelect={() => handleMenuItemClick('archive')}>
									Archive
								</DropdownMenu.Item>
								<DropdownMenu.Item onSelect={() => handleMenuItemClick('complete')}>
									Mark As Completed
								</DropdownMenu.Item>
								<DropdownMenu.Separator />
								<DropdownMenu.Item
									onSelect={() => handleMenuItemClick('delete')}
									class="text-destructive focus:text-destructive"
								>
									Delete
								</DropdownMenu.Item>
							</DropdownMenu.Group>
						</DropdownMenu.Content>
					</DropdownMenu.Root>
				{:else}
					<div class="flex items-center gap-1">
						<span class="shrink-0 text-xs text-muted-foreground">
							{formatTime(chat.lastMessageAt)}
						</span>
					</div>
				{/if}
			</div>
			<div class="mt-0.5 flex items-center gap-1">
				{#if chat.lastMessageStatus === 'read'}
					<CheckCheck class="h-3.5 w-3.5 shrink-0 text-primary" />
				{:else if chat.lastMessageStatus === 'delivered'}
					<CheckCheck class="h-3.5 w-3.5 shrink-0 text-muted-foreground" />
				{:else if chat.lastMessageStatus === 'sent'}
					<Check class="h-3.5 w-3.5 shrink-0 text-muted-foreground" />
				{/if}
				<span class="truncate text-xs text-muted-foreground">{chat.lastMessagePreview}</span>
				{#if !isHovered}
					<div class="ml-auto flex shrink-0 items-center justify-end gap-1.5">
						{#if chat.unreadCount}
							<span
								class=" grid h-5 min-w-5 place-items-center rounded-full bg-primary px-1 text-[10px] font-semibold text-primary-foreground"
							>
								{chat.unreadCount}
							</span>
						{/if}
						<span
							class={cn(
								'rounded-lg border px-1.5 py-0.5 text-[10px]',
								chat.agentName
									? 'dark:textborder-emerald-400 border-emerald-600 text-emerald-600 dark:border-emerald-400'
									: 'border-amber-400 text-amber-400'
							)}
						>
							{chat.agentName ? chat.agentName : 'Unassigned'}
						</span>
					</div>
				{/if}
			</div>
		</div>
	</button>
</div>
