<script lang="ts">
	import { page } from '$app/state';
	import { goto } from '$app/navigation';
	import { cn } from '$lib/utils';
	import { BotIcon, MessageSquare, LayoutGrid, CircleUserRound, Settings } from '@lucide/svelte';
	import { localizeHref } from '$lib/paraglide/runtime';

	const items = [
		{ icon: MessageSquare, label: 'Chats', href: '/app/chats' },
		{ icon: BotIcon, label: 'Ai Agent', href: '/app/ai-agents' },
		{ icon: CircleUserRound, label: 'Schedule', href: '/app/contacts' },
		{ icon: LayoutGrid, label: 'Notifications', href: '/app/integrations' },
		{ icon: Settings, label: 'Settings', href: '/app/settings' }
	];

	async function onNavClick(link: string) {
		await goto(localizeHref(link));
	}
</script>

<nav
	class="flex h-16 shrink-0 items-center justify-around bg-sidebar text-sidebar-foreground md:hidden"
>
	{#each items as item, i (i)}
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
</nav>
