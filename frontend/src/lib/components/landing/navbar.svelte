<script lang="ts" module>
	interface RouteProps {
		href: string;
		label: string;
	}

	interface FeatureProps {
		title: string;
		description: string;
	}
</script>

<script lang="ts">
	import { LightSwitch } from '$lib/components/extras/light-switch';
	import { LanguageSwitcher } from '$lib/components/extras/language-switcher';
	import {
		DropdownMenu,
		DropdownMenuContent,
		DropdownMenuItem,
		DropdownMenuTrigger
	} from '$lib/components/ui/dropdown-menu';
	import {
		Sheet,
		SheetContent,
		SheetFooter,
		SheetHeader,
		SheetTitle,
		SheetTrigger
	} from '$lib/components/ui/sheet';
	import { Button, buttonVariants } from '$lib/components/ui/button';
	import { Separator } from '$lib/components/ui/separator';
	import { Menu } from '@lucide/svelte';
	import { locales as availableLocales, localizeHref } from '$lib/paraglide/runtime';
	import { LanguageLabels } from '$lib/utils/localize-path.js';

	let { lang }: { lang: string } = $props();

	let isOpen = $state(false);
	// svelte-ignore state_referenced_locally
	let currentLang = $state(lang);

	const languages = availableLocales.map((code) => ({
		code,
		label: LanguageLabels[code] ?? code.toUpperCase()
	}));

	const routeList: RouteProps[] = [
		{ href: '#pricing', label: 'Pricing' },
		{ href: localizeHref('/about'), label: 'About' },
		{ href: localizeHref('/contact'), label: 'Contact' },
		{ href: localizeHref('/faq'), label: 'FAQ' }
	];

	const featureList: FeatureProps[] = [
		{
			title: 'Showcase Your Value ',
			description: 'Highlight how your product solves user problems.'
		},
		{
			title: 'Build Trust',
			description: 'Leverages social proof elements to establish trust and credibility.'
		},
		{
			title: 'Capture Leads',
			description: 'Make your lead capture form visually appealing and strategically.'
		}
	];
</script>

<header
	class="dark:shadow-dark shadow-light sticky top-5 z-40 mx-auto flex w-[90%] items-center justify-between rounded-2xl border bg-card p-2 px-4 shadow-md md:w-[70%] lg:w-[75%] lg:max-w-7xl"
>
	<a href={localizeHref('/')} class="flex items-center justify-center gap-2">
		<img src="/logo.png" alt="logo" class="h-7 w-auto" />
		<span class="text-xl font-medium text-foreground">SocialForge</span>
	</a>
	<div class="flex items-center lg:hidden">
		<Sheet bind:open={isOpen}>
			<SheetTrigger>
				<Menu onclick={() => (isOpen = true)} class="cursor-pointer" />
			</SheetTrigger>

			<SheetContent
				side="left"
				class="flex flex-col justify-between rounded-tr-2xl rounded-br-2xl bg-card"
			>
				<div>
					<SheetHeader class="mb-4 ml-4">
						<SheetTitle class="flex items-center">
							<a href={localizeHref('/')} class="flex items-center justify-center gap-2">
								<img src="/logo.png" alt="logo" class="h-7 w-auto" />
								<span class="text-xl font-medium text-foreground">SocialForge</span>
							</a>
						</SheetTitle>
					</SheetHeader>

					<div class="flex flex-col gap-2 px-3">
						{#each routeList as { href, label } (href)}
							<a onclick={() => (isOpen = false)} {href}>
								<Button variant="ghost" class="w-full justify-start text-base">
									{label}
								</Button>
							</a>
						{/each}
					</div>
					<Separator class="mb-2" />
					<div class="flex items-center justify-center gap-4">
						<Button href={localizeHref('/signin')} size="sm" variant="outline">Sign In</Button>
						<Button href={localizeHref('/signup')} size="sm">Get Started</Button>
					</div>
				</div>

				<SheetFooter class="flex-col items-center justify-center sm:flex-col">
					<Separator class="mb-2" />
					<div class="flex items-center justify-center gap-4">
						<LightSwitch />
						<LanguageSwitcher {languages} bind:value={currentLang} variant="outline" />
					</div>
				</SheetFooter>
			</SheetContent>
		</Sheet>
	</div>

	<div class="hidden items-center gap-1 lg:flex">
		<DropdownMenu>
			<DropdownMenuTrigger
				class={`${buttonVariants({ variant: 'ghost', size: 'default' })} text-base`}
			>
				Features
			</DropdownMenuTrigger>
			<DropdownMenuContent class="w-150">
				<div class="grid grid-cols-2 gap-5 p-4">
					<img
						src="https://github.com/sveltejs.png"
						alt="Beach"
						class="h-full w-full rounded-md object-cover"
					/>
					<ul class="flex flex-col gap-2">
						{#each featureList as { title, description } (title)}
							<DropdownMenuItem class="cursor-pointer rounded-md p-3 text-sm">
								<div>
									<p class="mb-1 leading-none font-semibold text-foreground">
										{title}
									</p>
									<p class="line-clamp-2 text-muted-foreground">
										{description}
									</p>
								</div>
							</DropdownMenuItem>
						{/each}
					</ul>
				</div>
			</DropdownMenuContent>
		</DropdownMenu>

		<!-- Navigation Links -->
		{#each routeList as { href, label } (label)}
			<a {href} class={buttonVariants({ variant: 'ghost', size: 'default' })}>
				{label}
			</a>
		{/each}
	</div>

	<div class="hidden lg:flex">
		<div class="flex items-center gap-x-4">
			<Button href={localizeHref('/signin')} size="sm" variant="outline">Sign In</Button>
			<Button href={localizeHref('/signup')} size="sm">Get Started</Button>
			<LightSwitch />
			<LanguageSwitcher {languages} bind:value={currentLang} variant="outline" />
		</div>
	</div>
</header>
