<script lang="ts">
	import { page } from '$app/state';
	import { goto } from '$app/navigation';
	import { cn } from '$lib/utils';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import {
		BotIcon,
		MessageSquare,
		LayoutGrid,
		CircleUserRound,
		Settings,
		ChartLine,
		Headset,
		EllipsisVertical
	} from '@lucide/svelte';
	import { localizeHref } from '$lib/paraglide/runtime';

	const mainItems = [
		{ icon: MessageSquare, label: 'Chats', href: '/app/chats' },
		{ icon: BotIcon, label: 'Ai Agent', href: '/app/ai-agents' },
		{ icon: CircleUserRound, label: 'Contacts', href: '/app/contacts' },
		{ icon: LayoutGrid, label: 'Integrations', href: '/app/integrations' }
	];
	const moreItems = [
		{ icon: Settings, label: 'Settings', href: '/app/settings' },
		{ icon: ChartLine, label: 'Analytics', href: '/app/analytics' },
		{ icon: Headset, label: 'Agent', href: '/app/agents' }
	];

	async function onNavClick(link: string) {
		await goto(localizeHref(link));
	}
</script>

<nav
	class="flex h-16 shrink-0 items-center justify-around bg-sidebar text-sidebar-foreground md:hidden"
>
	{#each mainItems as item, i (i)}
		<button
			class={cn(
				'grid h-10 w-10 place-items-center rounded-lg transition-colors',
				page.url.pathname.includes(item.href)
					? 'bg-primary text-white'
					: 'text-sidebar-foreground/70'
			)}
			title={item.label}
			onclick={() => onNavClick(item.href)}
		>
			<item.icon class="h-5.5 w-5.5" strokeWidth={1.75} />
		</button>
	{/each}
	<DropdownMenu.Root>
		<DropdownMenu.Trigger>
			{#snippet child({ props })}
				<button
					{...props}
					class="grid h-10 w-10 place-items-center rounded-lg transition-colors"
					title="More"
				>
					<EllipsisVertical class="h-5.5 w-5.5" strokeWidth={1.75} />
				</button>
			{/snippet}
		</DropdownMenu.Trigger>
		<DropdownMenu.Content align="center">
			<DropdownMenu.Group>
				{#each moreItems as item, i (i)}
					<DropdownMenu.Item onSelect={() => onNavClick(item.href)}>{item.label}</DropdownMenu.Item>
				{/each}
			</DropdownMenu.Group>
		</DropdownMenu.Content>
	</DropdownMenu.Root>
</nav>
