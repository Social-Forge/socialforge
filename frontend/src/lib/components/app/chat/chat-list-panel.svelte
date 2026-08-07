<script lang="ts">
	import { Button } from '$lib/components/ui/button/index.js';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu/index.js';
	import { Badge } from '$lib/components/ui/badge/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import ChatListItem from './chat-list-item.svelte';
	import { chatUiState } from '$lib/hooks/chat-ui.svelte';
	import { chatsStore } from '$lib/stores/chats';
	import { MoreVertical, Search, ChevronDown, Bookmark, Plus } from '@lucide/svelte';

	const bookmarked = $derived(chatsStore.bookmarked);
	const rest = $derived(chatsStore.rest);
	let selectedChat = $state<string>('all');
	let search = $state<string>('');

	const exampleLabels = [
		{ id: '0', label: 'All Chats', color: '#cbd5e1', textColor: '#000000' },
		{ id: '1', label: 'Pesanan Baru', color: '#00c335', textColor: '#FFFFFF' },
		{ id: '2', label: 'Pembayaran Tertunda', color: '#ae662a', textColor: '#FFFFFF' },
		{ id: '3', label: 'Sedang Diproses', color: '#006dff', textColor: '#FFFFFF' }
	];
</script>

<div class="flex h-full w-full flex-col md:w-[320px] md:border-r md:border-border lg:w-90">
	<div class="flex items-center justify-between px-4 pt-4 pb-3">
		<h1 class="text-xl font-bold text-foreground">Chats</h1>
		<div class="flex items-center gap-1 text-muted-foreground">
			<Button variant="ghost" size="icon" title="More">
				<Plus class="h-4.5 w-4.5" />
			</Button>
			<Button variant="ghost" size="icon" title="More">
				<MoreVertical class="h-4.5 w-4.5" />
			</Button>
		</div>
	</div>

	<div class="flex items-center gap-2 px-4 pb-3">
		<DropdownMenu.Root>
			<DropdownMenu.Trigger>
				{#snippet child({ props })}
					<Button
						{...props}
						variant="outline"
						class="w-full max-w-[120px] items-center justify-start truncate text-start"
					>
						{selectedChat === 'all'
							? 'All Chats'
							: selectedChat === 'pesanan'
								? 'pesanan Baru'
								: selectedChat === 'pembayaran'
									? 'Pembayaran Tertunda'
									: 'Sedang Diproses'}
						<ChevronDown class="h-3.5 w-3.5" />
					</Button>
				{/snippet}
			</DropdownMenu.Trigger>
			<DropdownMenu.Content>
				<DropdownMenu.Group>
					<DropdownMenu.Label>Select Chat</DropdownMenu.Label>
					<DropdownMenu.Separator />
					<DropdownMenu.Item onSelect={() => (selectedChat = 'all')}>All Chats</DropdownMenu.Item>
					<DropdownMenu.Item onSelect={() => (selectedChat = 'pesanan')}
						>Pesanan Baru</DropdownMenu.Item
					>
					<DropdownMenu.Item onSelect={() => (selectedChat = 'pembayaran')}
						>Pembayaran Tertunda</DropdownMenu.Item
					>
					<DropdownMenu.Item onSelect={() => (selectedChat = 'sedang')}
						>Sedang Diproses</DropdownMenu.Item
					>
				</DropdownMenu.Group>
			</DropdownMenu.Content>
		</DropdownMenu.Root>

		<div class="relative flex-1">
			<Search
				class="pointer-events-none absolute top-1/2 left-3 h-4 w-4 -translate-y-1/2 text-muted-foreground"
			/>
			<Input
				type="text"
				placeholder="Search users"
				class="w-full rounded-lg border border-border bg-secondary py-2 pr-3 pl-9 text-sm text-foreground placeholder:text-muted-foreground focus:ring-2 focus:ring-ring focus:outline-none"
			/>
		</div>
	</div>

	<div class="scrollbar-primary mt-2 flex items-center gap-2 overflow-x-auto px-4 pb-3">
		{#each exampleLabels as label (label.id)}
			<Badge
				class="cursor-pointer bg-neutral-800 text-[10px] font-semibold text-white active:scale-95 dark:bg-neutral-50 dark:text-neutral-900"
			>
				{label.label}
			</Badge>
		{/each}
	</div>

	<div class="scroll-thin flex-1 overflow-y-auto px-2 pb-4">
		{#if bookmarked.length > 0}
			<div
				class="flex items-center gap-1.5 px-2 pt-2 pb-1 text-xs font-semibold tracking-wide text-muted-foreground uppercase"
			>
				<Bookmark class="h-3.5 w-3.5" />
				Bookmarked
			</div>
			<div class="flex flex-col gap-0.5">
				{#each bookmarked as chat (chat.id)}
					<ChatListItem
						{chat}
						active={chatUiState.activeChatId === chat.id}
						onclick={() => chatUiState.openChat(chat.id)}
					/>
				{/each}
			</div>
		{/if}

		<div class="px-2 pt-4 pb-1 text-xs font-semibold tracking-wide text-muted-foreground uppercase">
			All Messages
		</div>
		<div class="flex flex-col gap-0.5">
			{#each rest as chat (chat.id)}
				<ChatListItem
					{chat}
					active={chatUiState.activeChatId === chat.id}
					onclick={() => chatUiState.openChat(chat.id)}
				/>
			{/each}
		</div>
	</div>
</div>
