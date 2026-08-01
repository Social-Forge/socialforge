<script lang="ts">
	import { AppAvatar } from '$lib/components/app';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu/index.js';
	import { chatUiState } from '$lib/hooks/chat-ui.svelte';
	import { buttonVariants } from '$lib/components/ui/button/index.js';
	import { ArrowLeft, EllipsisVertical, Info } from '@lucide/svelte';
	import type { ChatSummary } from './types';
	import { cn } from '$lib/utils';

	let { chat }: { chat: ChatSummary } = $props();

	const subtitle = $derived(
		chat.isGroup
			? `${chat.memberCount} member${chat.memberCount === 1 ? '' : 's'}, ${chat.onlineCount} online`
			: chat.presence === 'online'
				? 'Online'
				: 'Offline'
	);

	function handleKeyDown(event: KeyboardEvent) {
		if (event.key === 'Enter' || event.key === ' ') {
			event.preventDefault(); // Prevents page scrolling on Spacebar
			chatUiState.toggleInfo();
		}
	}

	function handleMenuItemClick(action: string) {
		console.log(`Action: ${action} for chat ${chat.id}`);
	}
</script>

<header class="flex h-18.25 shrink-0 items-center justify-between border-b border-border px-4">
	<div class="flex items-center gap-3">
		<button
			class="grid h-9 w-9 place-items-center rounded-md hover:bg-accent md:hidden"
			onclick={() => chatUiState.backToList()}
		>
			<ArrowLeft class="h-5 w-5" />
		</button>
		<AppAvatar
			src={chat.avatarUrl}
			fallback={chat.name.slice(0, 1)}
			size="md"
			channel={chat.channel}
		/>
		<div
			role="button"
			tabindex="0"
			onkeydown={handleKeyDown}
			onclick={() => chatUiState.toggleInfo()}
		>
			<div class="text-sm font-semibold text-foreground">{chat.name}</div>
			<!-- <div class="text-xs text-muted-foreground">{subtitle}</div> -->
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
	</div>

	<div class="flex items-center gap-1 text-muted-foreground">
		<button
			class="grid h-9 w-9 place-items-center rounded-md hover:bg-accent {chatUiState.infoOpen
				? 'bg-accent text-foreground'
				: ''}"
			title="Chat info"
			onclick={() => chatUiState.toggleInfo()}
		>
			<Info class="h-4.5 w-4.5" />
		</button>
		<DropdownMenu.Root>
			<DropdownMenu.Trigger class={buttonVariants({ variant: 'ghost', size: 'icon-sm' })}>
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
					<DropdownMenu.Item onSelect={() => handleMenuItemClick('leave')}>
						Leave Chat
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
	</div>
</header>
