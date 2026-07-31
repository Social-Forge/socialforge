<script lang="ts">
	import * as Sheet from '$lib/components/ui/sheet';
	import * as Accordion from '$lib/components/ui/accordion';
	import { Menu, X } from '@lucide/svelte';
	import { Button } from '$lib/components/ui/button';
	import { LightSwitch } from '$lib/components/extras/light-switch';
	import { LanguageSwitcher } from '$lib/components/extras/language-switcher';
	import { mode } from 'mode-watcher';
	import { cn } from '$lib/utils';
	import { NAV_LINKS } from '$lib/components/landing';
	import { type Locale, locales as availableLocales, localizeHref } from '$lib/paraglide/runtime';
	import { LanguageLabels } from '$lib/utils/localize-path.js';

	let { lang }: { lang: string } = $props();

	let isOpen = $state(false);
	// svelte-ignore state_referenced_locally
	let currentLang = $state(lang);

	const languages = availableLocales.map((code) => ({
		code,
		label: LanguageLabels[code] ?? code.toUpperCase()
	}));
</script>

<div class="flex items-center justify-end lg:hidden">
	<Sheet.Root bind:open={isOpen} onOpenChange={(open) => (isOpen = open)}>
		<Sheet.Trigger
			class={cn('flex cursor-pointer items-center justify-center active:scale-95', {
				hidden: isOpen
			})}
		>
			<Menu class="h-5 w-5" />
		</Sheet.Trigger>
		<Sheet.Content side="right" class="w-full py-10">
			<Sheet.Close class="bg-background text-foreground">
				<div class="absolute top-3 right-5 z-20 flex items-center justify-center bg-background">
					<X class="h-5 w-5" />
				</div>
			</Sheet.Close>
			<div class="flex max-w-3xl flex-col items-start py-2">
				<div class="flex w-full items-center justify-center gap-4">
					<Button href={localizeHref('/signin')} variant="outline" class="w-[43%]">Sign In</Button>
					<Button href={localizeHref('/signup')} class="w-[43%]">Sign Up</Button>
				</div>
				<ul class="mt-6 flex w-full flex-col items-start">
					<Accordion.Root type="single" class="w-full">
						{#each NAV_LINKS as link (link.title)}
							<Accordion.Item value={link.title} class="px-5 last:border-none!">
								{#if link.menu}
									<Accordion.Trigger class="w-full justify-start py-4">
										{link.title}
									</Accordion.Trigger>
									<Accordion.Content class="w-full">
										<div
											role="button"
											tabindex="0"
											aria-label={`${link.title} menu`}
											class={cn('w-full cursor-pointer')}
											onclick={() => (isOpen = false)}
											onkeydown={(e) => {
												if (e.key === 'Enter') {
													isOpen = false;
												}
											}}
										>
											{#each link.menu as item (item.title)}
												<a
													href={localizeHref(item.href)}
													title={item.title}
													class={cn(
														'block space-y-1 rounded-lg p-3 leading-none no-underline transition-colors outline-none select-none hover:bg-accent hover:text-accent-foreground focus:bg-accent focus:text-accent-foreground'
													)}
												>
													<div class="flex items-center space-x-2 text-foreground">
														<item.icon class="h-5 w-5" />
														<h6 class="text-sm leading-none!">{item.title}</h6>
													</div>
													<p
														title={item.tagline}
														class="line-clamp-1 text-sm leading-snug text-muted-foreground"
													>
														{item.tagline}
													</p>
												</a>
											{/each}
										</div>
									</Accordion.Content>
								{:else}
									<a
										href={localizeHref(link.href)}
										class="block space-y-1 rounded-lg py-4 leading-none no-underline transition-colors outline-none select-none hover:bg-accent hover:text-accent-foreground focus:bg-accent focus:text-accent-foreground"
										onclick={() => (isOpen = false)}
										onkeydown={(e) => {
											if (e.key === 'Enter') {
												isOpen = false;
											}
										}}
									>
										<span class="block w-full text-left">{link.title}</span>
									</a>
								{/if}
							</Accordion.Item>
						{/each}
						<Accordion.Item class="px-5 last:border-none!">
							<div
								class="flex items-center justify-between space-y-1 rounded-lg py-4 leading-none no-underline transition-colors outline-none select-none hover:bg-accent hover:text-accent-foreground focus:bg-accent focus:text-accent-foreground"
							>
								<span class="block w-full text-left">
									{mode.current === 'light' ? 'Dark' : 'Light'} Mode
								</span>
								<LightSwitch />
							</div>
						</Accordion.Item>
						<Accordion.Item class="px-5 last:border-none!">
							<div
								class="flex items-center justify-between space-y-1 rounded-lg py-4 leading-none no-underline transition-colors outline-none select-none hover:bg-accent hover:text-accent-foreground focus:bg-accent focus:text-accent-foreground"
							>
								<span class="block w-full text-left">
									{languages.find((lang) => lang.code === currentLang)?.label || 'Select Language'}
								</span>
								<LanguageSwitcher {languages} bind:value={currentLang} variant="outline" />
							</div>
						</Accordion.Item>
					</Accordion.Root>
				</ul>
			</div>
		</Sheet.Content>
	</Sheet.Root>
</div>
