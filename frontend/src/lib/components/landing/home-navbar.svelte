<script lang="ts" module>
	type ListItemMenuProps = {
		className?: string;
		title: string;
		href: string;
		content: string;
	};
</script>

<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import { browser } from '$app/environment';
	import { cn } from '$lib/utils';
	import { HomeAnimationContainer, MaxWidthWrapper, HomeMobileNavbar } from './index';
	import { LightSwitch } from '$lib/components/extras/light-switch';
	import { LanguageSwitcher } from '$lib/components/extras/language-switcher';
	import * as NavigationMenu from '$lib/components/ui/navigation-menu/index.js';
	import { ZapIcon } from '@lucide/svelte';
	import { Button } from '$lib/components/ui/button';
	import { NAV_LINKS } from '$lib/components/landing';
	import { locales as availableLocales, localizeHref } from '$lib/paraglide/runtime';
	import { LanguageLabels } from '$lib/utils/localize-path.js';

	let { lang }: { lang: string } = $props();

	let isScroll = $state(false);

	// svelte-ignore state_referenced_locally
	let currentLang = $state(lang);

	const languages = availableLocales.map((code) => ({
		code,
		label: LanguageLabels[code] ?? code.toUpperCase()
	}));

	const handleScroll = () => {
		isScroll = window.scrollY > 0;
	};

	onMount(() => {
		if (browser) {
			window.addEventListener('scroll', handleScroll);
		}
		return () => {
			window.removeEventListener('scroll', handleScroll);
		};
	});

	onDestroy(() => {
		if (browser) {
			window.removeEventListener('scroll', handleScroll);
		}
	});
</script>

<header
	class={cn(
		'sticky inset-x-0 top-0 z-99999 h-14 w-full border-b border-transparent select-none',
		isScroll ? 'border-background/80 bg-background/40 backdrop-blur-md' : ''
	)}
>
	<HomeAnimationContainer delay={0.1} class="size-full py-3">
		<MaxWidthWrapper class="flex items-center justify-between">
			<div class="flex items-center space-x-12">
				<a href={localizeHref('/')} class="flex items-center justify-center gap-2">
					<img src="/logo.png" alt="logo" class="h-7 w-auto" />
					<span class="text-xl font-medium text-foreground">SocialForge</span>
				</a>
				<NavigationMenu.Root class=" hidden lg:flex">
					<NavigationMenu.List>
						{#each NAV_LINKS as link (link.title)}
							<NavigationMenu.Item>
								{#if link.menu}
									<NavigationMenu.Trigger
										class="bg-transparent hover:bg-transparent active:bg-transparent"
									>
										{link.title}
									</NavigationMenu.Trigger>
									<NavigationMenu.Content>
										<ul
											class={cn(
												'grid gap-1 rounded-xl p-4 md:w-100 lg:w-125',
												link.title === 'Features' ? 'lg:grid-cols-[.75fr_1fr]' : 'lg:grid-cols-2'
											)}
										>
											{#if link.title === 'Features'}
												<li class="relative row-span-4 overflow-hidden rounded-lg pr-2">
													<div
														class="absolute inset-0 z-10! h-full w-[calc(100%-10px)] bg-[linear-gradient(to_right,rgb(38,38,38,0.5)_1px,transparent_1px),linear-gradient(to_bottom,rgb(38,38,38,0.5)_1px,transparent_1px)] bg-size-[1rem_1rem]"
													></div>
													<NavigationMenu.Link class="relative z-20 h-full">
														{#snippet children()}
															<!-- svelte-ignore a11y_invalid_attribute -->
															<a
																href="#"
																class="flex h-full w-full flex-col justify-end rounded-lg bg-linear-to-b from-muted/50 to-muted p-4 no-underline outline-none select-none focus:shadow-md"
															>
																<h6 class="mt-4 mb-2 text-lg font-medium">All Features</h6>
																<p class="text-sm leading-tight text-muted-foreground">
																	Manage links, track performance, and more.
																</p>
															</a>
														{/snippet}
													</NavigationMenu.Link>
												</li>
											{/if}
											{#each link.menu as subItem (subItem.title)}
												<li>
													<NavigationMenu.Link
														href={localizeHref(subItem.href)}
														class={cn(
															'block space-y-1 rounded-lg p-3 leading-none no-underline transition-all duration-100 ease-out outline-none select-none hover:bg-accent hover:text-accent-foreground focus:bg-accent focus:text-accent-foreground'
														)}
													>
														<div
															class="flex items-center space-x-2 text-neutral-600 dark:text-neutral-300"
														>
															<subItem.icon class="h-4 w-4" />
															<h6 class="text-sm leading-none! font-medium">
																{subItem.title}
															</h6>
														</div>
														<p
															title={subItem.tagline}
															class="line-clamp-1 text-sm leading-snug text-muted-foreground"
														>
															{subItem.tagline}
														</p>
													</NavigationMenu.Link>
												</li>
											{/each}
										</ul>
									</NavigationMenu.Content>
								{:else}
									<NavigationMenu.Link href={link.href}>
										{link.title}
									</NavigationMenu.Link>
								{/if}
							</NavigationMenu.Item>
						{/each}
					</NavigationMenu.List>
				</NavigationMenu.Root>
			</div>
			<div class="hidden items-center lg:flex">
				<div class="flex items-center gap-x-4">
					<Button href={localizeHref('/signin')} size="sm" variant="ghost">Sign In</Button>
					<Button href={localizeHref('/signup')} size="sm">
						Get Started
						<ZapIcon class="ml-1.5 size-3.5 fill-orange-500 text-orange-500" />
					</Button>
					<LightSwitch />
					<LanguageSwitcher {languages} bind:value={currentLang} variant="outline" />
				</div>
			</div>
			<HomeMobileNavbar {lang} />
		</MaxWidthWrapper>
	</HomeAnimationContainer>
</header>

<style scoped>
	/* Your styles here */
</style>
