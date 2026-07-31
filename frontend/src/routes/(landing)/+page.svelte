<script lang="ts">
	import { MetaTags } from 'svelte-meta-tags';
	import {
		HomeNavbar,
		MaxWidthWrapper,
		HomeAnimationContainer,
		COMPANIES
	} from '$lib/components/landing';
	import { ArrowRight } from '@lucide/svelte';
	import { Button } from '$lib/components/ui/button/index.js';
	import BorderBeam from '$lib/components/border-beam.svelte';
	import MagicBadge from '$lib/components/magic-badge.svelte';
	import { setMode, mode } from 'mode-watcher';
	// import { COMPANIES } from '@/constants/misc.js';

	let { data } = $props();
	let metaTags = $derived(data.pageMetaTags);

	let isDark = $derived(mode.current === 'dark');
</script>

<div class="size-full overflow-x-hidden">
	<MaxWidthWrapper>
		<div
			class="flex w-full flex-col items-center justify-center bg-linear-to-t from-background text-center"
		>
			<HomeAnimationContainer
				variant="fade-up"
				delay={0.2}
				once={false}
				duration={1}
				distance={60}
				class="flex w-full flex-col items-center justify-center text-center"
			>
				<button
					class="group relative grid overflow-hidden rounded-full px-4 py-1 shadow-[0_1000px_0_0_hsl(0_0%_20%)_inset] transition-colors duration-200"
				>
					<span>
						<span
							class="spark mask-gradient absolute inset-0 h-full w-full animate-flip overflow-hidden rounded-full [mask:linear-gradient(white,transparent_50%)] before:absolute before:inset-[0_auto_auto_50%] before:aspect-square before:w-[200%] before:[translate:-50%_-15%] before:-rotate-90 before:animate-rotate before:bg-[conic-gradient(from_0deg,transparent_0_340deg,white_360deg)] before:content-['']"
						></span>
					</span>
					<span
						class="backdrop absolute inset-px rounded-full bg-neutral-950 transition-colors duration-200 group-hover:bg-neutral-900"
					></span>
					<span
						class="absolute inset-x-0 bottom-0 h-full w-full bg-linear-to-tr from-primary/20 blur-md"
					></span>
					<span class="z-10 flex items-center justify-center gap-1 py-0.5 text-sm text-neutral-100">
						✨ Smart CRM and Web Builder
						<ArrowRight
							class="ml-1 size-3 transition-transform duration-300 ease-in-out group-hover:translate-x-0.5"
						/>
					</span>
				</button>
				<h1
					class="font-heading w-full py-6 text-center text-5xl leading-[1.15]! font-medium tracking-normal text-balance text-foreground sm:text-6xl md:text-7xl lg:text-8xl"
				>
					AI-powered <span
						class="inline-bloc bg-linear-to-r from-violet-500 to-fuchsia-500 bg-clip-text text-transparent"
					>
						Omnichannel CRM and Web Builder
					</span>
				</h1>
				<p class="mb-12 text-lg tracking-tight text-balance text-muted-foreground md:text-xl">
					Helps businesses to connect with customers across multiple channels, including email,
					chat, social media and web builder.
					<br class="hidden md:block" />
					<span class="hidden md:block"
						>Boosting marketing, customer service, and sales growth by up to 5x.</span
					>
				</p>
				<div class="z-50 flex items-center justify-center gap-4 whitespace-nowrap">
					<Button href="/auth/sign-in" class="flex items-center">
						Book a Demo
						<ArrowRight class="ml-2 h-4 w-4" />
					</Button>
				</div>
			</HomeAnimationContainer>
			<HomeAnimationContainer
				variant="fade-up"
				delay={0.2}
				once={false}
				duration={1}
				distance={60}
				class="relative w-full bg-transparent px-2 pt-20 pb-20 md:py-32"
			>
				<div
					class="mask-gradient absolute inset-0 left-1/2 h-1/4 w-3/4 -translate-x-1/2 animate-image-glow blur-[5rem] md:top-[10%] md:h-1/3"
				></div>
				<div
					class="bg-opacity-50 -m-2 rounded-xl p-2 ring-1 ring-foreground/20 backdrop-blur-3xl ring-inset lg:-m-4 lg:rounded-2xl"
				>
					<BorderBeam size={250} duration={12} delay={9} />
					{#if isDark}
						<img
							src="/images/dark-laptop.png"
							alt="Dashboard"
							width={1200}
							height={1200}
							loading="eager"
							decoding="async"
							class="rounded-md bg-foreground/10 ring-1 ring-border lg:rounded-xl"
						/>
					{:else}
						<img
							src="/images/light-laptop.png"
							alt="Dashboard"
							width={1200}
							height={1200}
							loading="eager"
							decoding="async"
							class="rounded-md bg-foreground/10 ring-1 ring-border lg:rounded-xl"
						/>
					{/if}
					<div
						class="absolute inset-x-0 -bottom-4 z-40 h-1/2 w-full bg-linear-to-t from-background"
					></div>
					<div
						class="absolute inset-x-0 bottom-0 z-50 h-1/4 w-full bg-linear-to-t from-background md:-bottom-8"
					></div>
				</div>
			</HomeAnimationContainer>
		</div>
	</MaxWidthWrapper>
	<MaxWidthWrapper>
		<HomeAnimationContainer
			variant="fade-up"
			delay={0.2}
			once={false}
			duration={1}
			distance={60}
			class="relative"
		>
			<div class="container mx-auto w-full max-w-6xl px-4 py-14 md:px-8">
				<div class=" relative flex flex-col space-y-8 rounded-2xl px-6 py-10 md:px-10 md:py-14">
					<h2
						class="font-heading text-center text-base font-semibold tracking-wider text-neutral-800 uppercase dark:text-neutral-200"
					>
						Trusted by the best in the industry
					</h2>
					<div class="relative w-full overflow-hidden">
						<div
							class="pointer-events-none absolute top-0 left-0 z-10 h-full w-20 bg-linear-to-r from-white to-transparent dark:from-neutral-950 dark:to-transparent"
						></div>
						<div
							class="pointer-events-none absolute top-0 right-0 z-10 h-full w-20 bg-linear-to-l from-white to-transparent dark:from-neutral-950 dark:to-transparent"
						></div>
						<ul
							class="flex animate-marquee flex-wrap items-center justify-center gap-x-6 gap-y-6 md:gap-x-16"
						>
							{#each COMPANIES as company (company.name)}
								<li>
									<img
										src={company.logo}
										alt={company.name}
										width={80}
										height={80}
										loading="eager"
										decoding="async"
										class="h-10 w-auto transition-opacity hover:opacity-75 md:w-28"
									/>
								</li>
							{/each}
						</ul>
					</div>
				</div>
			</div>
		</HomeAnimationContainer>
	</MaxWidthWrapper>
	<MaxWidthWrapper class="pt-10">
		<HomeAnimationContainer variant="fade-down" threshold={0.1} once={false}>
			<div class="flex w-full flex-col items-center justify-center py-8 lg:items-center">
				<MagicBadge title="Features" />
				<h2
					class="font-heading mt-6 text-center text-3xl leading-[1.1]! font-medium text-foreground md:text-5xl lg:text-center"
				>
					Manage Links Like a Pro
				</h2>
				<p class="mt-4 max-w-lg text-center text-lg text-muted-foreground lg:text-center">
					Linkify is a powerful link management tool that helps you shorten, track, and organize all
					your links in one place.
				</p>
			</div>
		</HomeAnimationContainer>
		<HomeAnimationContainer
			variant="fade-up"
			threshold={0.2}
			once={false}
			duration={1}
			distance={30}
		></HomeAnimationContainer>
	</MaxWidthWrapper>
</div>

<style scoped>
	.mask-gradient {
		background: conic-gradient(
			from 230.29deg at 51.63% 52.16%,
			rgb(36, 0, 255) 0deg,
			rgb(0, 135, 255) 67.5deg,
			rgb(108, 39, 157) 198.75deg,
			rgb(24, 38, 163) 251.25deg,
			rgb(54, 103, 196) 301.88deg,
			rgb(105, 30, 255) 360deg
		);
	}
</style>
