<script lang="ts">
	import { goto } from '$app/navigation';
	import { Button } from '$lib/components/ui/button/index.js';
	import * as UnderlineTabs from '$lib/components/extras/underline-tabs/index.js';
	import { Separator } from '$lib/components/ui/separator/index.js';
	import * as Empty from '$lib/components/ui/empty/index.js';
	import { localizeHref } from '$lib/paraglide/runtime.js';
	import { Plus, Lock, User, Rocket, ArrowLeft } from '@lucide/svelte';

	let { data } = $props();

	let selectedTab = $state('channel');

	const dummyChannel = [
		{
			name: 'Linkchat',
			description: 'Start communicating with your customers via Linkchat',
			integrations: 0,
			href: '/app/integrations/linkchat',
			icon: '/images/channel/linkchat.svg'
		},
		{
			name: 'Webchat',
			description: 'Start communicating with your customers via Webchat',
			integrations: 0,
			href: '/app/integrations/webchat',
			icon: '/images/channel/webchat.svg'
		},
		{
			name: 'WhatsApp',
			description: 'Start communicating with your customers via WhatsApp',
			integrations: 0,
			href: '/app/integrations/whatsapp',
			icon: '/images/channel/whatsapp-unofficial.svg'
		},
		{
			name: 'WhatsApp API',
			description: 'Start communicating with your customers via WhatsApp API',
			integrations: 0,
			href: '/app/integrations/whatsapp',
			icon: '/images/channel/whatsapp-official.svg'
		},
		{
			name: 'Messenger',
			description: 'Start communicating with your customers via Messenger',
			integrations: 0,
			href: '/app/integrations/messenger',
			icon: '/images/channel/messenger.svg'
		},
		{
			name: 'Telegram',
			description: 'Start communicating with your customers via Telegram',
			integrations: 0,
			href: '/app/integrations/telegram',
			icon: '/images/channel/telegram.svg'
		},
		{
			name: 'Instagram',
			description: 'Start communicating with your customers via Instagram',
			integrations: 0,
			href: '/app/integrations/instagram',
			icon: '/images/channel/instagram.svg'
		}
	];

	async function onClickIntegration(url: string) {
		// Navigate to the integration page for the selected channel
		await goto(localizeHref(url));
	}
</script>

<div class="h-[calc(100dvh-70px)] w-full overflow-y-auto bg-background lg:h-full">
	<div class="relative flex h-full w-full">
		<div class="scrollbar-primary flex h-full w-full flex-col overflow-y-auto">
			<div class="flex flex-col gap-2 px-6 pt-6">
				<div class="text-xl font-semibold">App Integration and Connection</div>
				<div class="text-sm text-muted-foreground">
					Manage and connect multiple communication channels and third-party applications in one
					place.
				</div>
			</div>
			<div class="mt-6 flex flex-col gap-4 px-6 md:flex-row md:items-center md:justify-between">
				<UnderlineTabs.Root value={selectedTab} onValueChange={(value) => (selectedTab = value)}>
					<UnderlineTabs.List>
						<UnderlineTabs.Trigger value="channel">Integrasi Channel</UnderlineTabs.Trigger>
						<UnderlineTabs.Trigger value="app">Integrasi App</UnderlineTabs.Trigger>
					</UnderlineTabs.List>
				</UnderlineTabs.Root>
				<div class="flex items-center gap-3 text-xl whitespace-nowrap">
					<Button variant="ghost" class="border border-primary text-primary hover:text-primary"
						><Lock /></Button
					>
					<Button variant="ghost" class="border border-primary text-primary hover:text-primary"
						><User /></Button
					>
					<div class="text-sm text-muted-foreground">
						Limit Integrations <span class="ml-2 font-semibold">0</span><span class="font-semibold"
							>/0</span
						>
					</div>
				</div>
			</div>
			<Separator class="my-4" />
			{#if selectedTab === 'channel'}
				<div class="relative mt-3 flex-1 px-6 pb-6">
					<div
						class="mt-6 mb-45 grid w-full grid-cols-1 gap-6 text-[12.8px] sm:grid-cols-2 lg:mb-0 lg:grid-cols-3"
					>
						{#each dummyChannel as channel (channel.name)}
							<div
								class="flex w-full flex-col justify-between gap-6 rounded-lg border border-border bg-card p-6 shadow-md transition-all hover:scale-105 hover:shadow-lg"
							>
								<div class="flex items-start justify-between">
									<div class="relative flex">
										<img
											class="top-[35.5%] left-[19.5%] h-[29%] w-14 object-contain"
											src={channel.icon}
											alt={channel.name}
										/>
									</div>
									<div class="flex items-center gap-3">
										<a href={localizeHref(channel.href)} class="rounded-lg px-4 py-2">
											<div class="flex items-center gap-2">
												<span> See Integrations </span>
												<div>(0)</div>
											</div>
										</a>
										<Button
											variant="ghost"
											size="icon"
											class="border border-primary text-primary hover:text-primary"
											onclick={() => onClickIntegration(channel.href)}
										>
											<Plus />
										</Button>
									</div>
								</div>
								<div class="flex flex-col gap-3 text-left text-[22.72px]">
									<div class="leading-[150%] font-semibold">{channel.name}</div>
									<div class="text-[14.4px] leading-[150%] text-muted-foreground">
										{channel.description}
									</div>
								</div>
							</div>
						{/each}
					</div>
				</div>
			{:else}
				<div
					class="relative mt-3 flex h-full w-full flex-col items-center justify-center px-6 pb-6"
				>
					<div class="mx-auto">
						<Empty.Root class="h-full w-full">
							<Empty.Header>
								<Empty.Media variant="icon" class="mb-10">
									<Rocket class="size-20" />
								</Empty.Media>
								<Empty.Title class="text-2xl font-semibold">Coming Soon</Empty.Title>
								<Empty.Description>We're working on this feature. Stay tuned!</Empty.Description>
							</Empty.Header>
							<Empty.Content>
								<Button onclick={() => (selectedTab = 'channel')} class="mt-4">
									<ArrowLeft />
									Back
								</Button>
							</Empty.Content>
						</Empty.Root>
					</div>
				</div>
			{/if}
		</div>
	</div>
</div>
