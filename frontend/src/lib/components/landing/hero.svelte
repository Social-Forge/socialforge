<script lang="ts">
	import { MaxWidthWrapper, HomeAnimationContainer } from '$lib/components/landing';
	import { Button, buttonVariants } from '$lib/components/ui/button/index.js';
	import { localizeHref } from '$lib/paraglide/runtime.js';
	import { mode } from 'mode-watcher';
	import { ArrowRight } from '@lucide/svelte';

	let isDark = $derived(mode.current === 'dark');
</script>

<MaxWidthWrapper>
	<HomeAnimationContainer
		variant="fade-up"
		delay={0.2}
		once={false}
		duration={1}
		distance={30}
		class="flex w-full flex-col items-center justify-center text-center"
	>
		<section class="container">
			<div class="mx-auto grid place-items-center gap-8 pb-20 md:pb-32 lg:max-w-7xl">
				<div class="flex w-full flex-col items-center justify-center gap-8 text-center">
					<button
						class="group relative grid overflow-hidden rounded-full px-4 py-1 shadow-[0_1000px_0_0_hsl(0_0%_20%)_inset] transition-colors duration-200"
					>
						<span>
							<span
								class="spark mask-gradient absolute inset-0 h-full w-full animate-flip overflow-hidden rounded-full [mask:linear-gradient(white,transparent_50%)] before:absolute before:inset-[0_auto_auto_50%] before:aspect-square before:w-[200%] before:[translate:-50%_-15%] before:-rotate-90 before:animate-rotate before:bg-[conic-gradient(from_0deg,transparent_0_340deg,white_360deg)] before:content-['']"
							></span>
						</span>
						<span
							class="backdrop absolute inset-px rounded-full bg-white/95 transition-colors duration-200 group-hover:bg-white/90 dark:bg-neutral-950 dark:group-hover:bg-neutral-900"
						></span>
						<span
							class="absolute inset-x-0 bottom-0 h-full w-full bg-linear-to-tr from-primary/20 blur-md"
						></span>
						<span
							class="z-10 flex items-center justify-center gap-1 py-0.5 text-sm text-neutral-950 dark:text-neutral-100"
						>
							✨ Smart CRM and Web Builder
							<ArrowRight
								class="ml-1 size-3 transition-transform duration-300 ease-in-out group-hover:translate-x-0.5"
							/>
						</span>
					</button>

					<div class="mx-auto max-w-3xl text-center text-5xl font-bold md:text-6xl">
						<h1
							class="font-heading w-full py-6 text-center leading-[1.15]! font-medium tracking-normal text-balance text-foreground"
						>
							<span
								class="inline-bloc bg-linear-to-r from-neutral-800 to-neutral-600 bg-clip-text text-transparent dark:from-neutral-200 dark:to-neutral-50"
							>
								AI-powered Omnichannel</span
							>
							<span
								class="inline-bloc bg-linear-to-r from-violet-500 to-fuchsia-500 bg-clip-text text-transparent"
							>
								CRM and Multi Agent
							</span>
						</h1>
					</div>

					<p class="mx-auto max-w-screen-sm text-xl text-muted-foreground">
						Helps businesses to connect with customers across multiple channels, including email,
						chat, social media and web builder.
						<br class="hidden md:block" />
						<span class="hidden md:block"
							>Boosting marketing, customer service, and sales growth by up to 5x.</span
						>
					</p>

					<div class="space-y-4 md:space-y-0 md:space-x-4">
						<Button
							href={localizeHref('/signup')}
							variant="default"
							class="group/arrow w-40 font-bold"
						>
							Get Started
							<ArrowRight
								class="ml-2 size-5 transition-transform group-hover/arrow:translate-x-1"
							/>
						</Button>
						<a
							href="#pricing"
							class={buttonVariants({ size: 'lg', variant: 'outline', class: 'w-40 font-bold' })}
						>
							Pricing
						</a>
					</div>
				</div>

				<div class="group relative mt-14">
					<div
						class="img-shadow-animation absolute -top-6 right-12 h-12 w-[90%] rounded-full bg-primary/50 blur-3xl lg:h-[80%]"
					></div>

					{#if isDark}
						<img
							src="/images/hero-image-dark.jpg"
							alt="Dashboard"
							loading="eager"
							decoding="async"
							class="rouded-lg img-border-animation relative mx-auto flex w-full items-center rounded-lg border border-t-2 border-t-primary/30 leading-none md:w-300"
						/>
					{:else}
						<img
							src="/images/hero-image-light.jpg"
							alt="Dashboard"
							loading="eager"
							decoding="async"
							class="rouded-lg img-border-animation relative mx-auto flex w-full items-center rounded-lg border border-t-2 border-t-primary/30 leading-none md:w-300"
						/>
					{/if}

					<div
						class="absolute bottom-0 left-0 h-20 w-full rounded-lg bg-linear-to-b from-background/0 via-background/50 to-background md:h-28"
					></div>
				</div>
			</div>
		</section>
	</HomeAnimationContainer>
</MaxWidthWrapper>

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
