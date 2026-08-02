<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import { LightSwitch } from '$lib/components/extras/light-switch';
	import { AppNavUser, AppNavItem } from '$lib/components/app';
	import { Separator } from '$lib/components/ui/separator';
	import {
		LayoutGrid,
		MessageSquare,
		Settings,
		ChartLine,
		BotIcon,
		CircleUserRound,
		Funnel
	} from '@lucide/svelte';
	import { localizeHref } from '$lib/paraglide/runtime';

	const topItems = [
		{ icon: MessageSquare, label: 'Chats', href: '/app/chats' },
		{ icon: ChartLine, label: 'Analytics', href: '/app/analytics' },
		{ icon: BotIcon, label: 'Ai Agent', href: '/app/ai-agents' },
		{ icon: CircleUserRound, label: 'Contact', href: '/app/contacts' }
	];

	const bottomItems = [
		{ icon: LayoutGrid, label: 'Notifications', href: '/app/integrations' },
		{ icon: Settings, label: 'Settings', href: '/app/settings' }
	];

	async function onNavClick(link: string) {
		await goto(localizeHref(link));
	}
</script>

<nav
	class="hidden h-full w-16 shrink-0 flex-col items-center justify-between bg-sidebar py-4 text-sidebar-foreground md:flex"
>
	<div class="flex flex-col items-center gap-6">
		<a href={localizeHref('/app/chats')} class="flex items-center justify-center">
			<img src="/logo.png" alt="Social Forge" class="h-9 w-9 object-cover" />
		</a>
		<div class="flex flex-col items-center gap-1">
			{#each topItems as item, i (i)}
				<AppNavItem
					icon={item.icon}
					label={item.label}
					active={page.url.pathname.includes(item.href)}
					onClick={() => onNavClick(item.href)}
				/>
			{/each}
			<Separator class="my-4" />
			{#each bottomItems as item, i (i)}
				<AppNavItem
					icon={item.icon}
					label={item.label}
					active={page.url.pathname.includes(item.href)}
					onClick={() => onNavClick(item.href)}
				/>
			{/each}
		</div>
	</div>

	<div class="flex flex-col items-center gap-3">
		<LightSwitch />
		<AppNavUser />
	</div>
</nav>
