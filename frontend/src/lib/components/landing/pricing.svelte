<script lang="ts" module>
	interface PlanProps {
		title: string;
		popular: boolean;
		price: number;
		description: string;
		buttonText: string;
		benefitList: string[];
	}
</script>

<script lang="ts">
	import * as Card from '$lib/components/ui/card';
	import { Button } from '$lib/components/ui/button';
	import { MaxWidthWrapper, HomeAnimationContainer } from '$lib/components/landing';
	import { Check } from '@lucide/svelte';

	const plans: PlanProps[] = [
		{
			title: 'Free',
			popular: false,
			price: 0,
			description: 'Lorem ipsum dolor sit, amet ipsum consectetur adipisicing elit.',
			buttonText: 'Start Free Trial',
			benefitList: [
				'1 team member',
				'1 GB storage',
				'Upto 2 pages',
				'Community support',
				'AI assistance'
			]
		},
		{
			title: 'Premium',
			popular: true,
			price: 45,
			description: 'Lorem ipsum dolor sit, amet ipsum consectetur adipisicing elit.',
			buttonText: 'Get Started',
			benefitList: [
				'4 team member',
				'8 GB storage',
				'Upto 6 pages',
				'Priority support',
				'AI assistance'
			]
		},
		{
			title: 'Enterprise',
			popular: false,
			price: 120,
			description: 'Lorem ipsum dolor sit, amet ipsum consectetur adipisicing elit.',
			buttonText: 'Contact Us',
			benefitList: [
				'10 team member',
				'20 GB storage',
				'Upto 10 pages',
				'Phone & email support',
				'AI assistance'
			]
		}
	];
</script>

<MaxWidthWrapper>
	<HomeAnimationContainer
		variant="fade-up"
		delay={0.2}
		once={false}
		duration={1}
		distance={30}
		class="relative w-full bg-transparent px-2"
	>
		<section class="container py-16 sm:py-32">
			<h2 class="mb-2 text-center text-lg tracking-wider text-primary">Pricing</h2>

			<h2 class="mb-4 text-center text-3xl font-bold md:text-4xl">Get unlimited access</h2>

			<h3 class="mx-auto pb-14 text-center text-xl text-muted-foreground md:w-1/2">
				Lorem ipsum dolor sit amet consectetur adipisicing reiciendis.
			</h3>

			<!-- Mobile First Layout -->
			<div class="block space-y-6 lg:hidden">
				{#each plans as { title, popular, price, description, buttonText, benefitList }, index (index)}
					<div class="relative">
						<Card.Root class={popular ? 'popular-mobile' : ''}>
							{#if popular}
								<div class="absolute -top-4 left-1/2 -translate-x-1/2">
									<span
										class="rounded-full bg-primary px-3 py-1 text-xs font-medium whitespace-nowrap text-primary-foreground sm:px-4 sm:text-sm"
									>
										Most Popular
									</span>
								</div>
							{/if}
							<Card.Header>
								<Card.Title class="pb-2 text-lg sm:text-xl">
									{title}
								</Card.Title>

								<Card.Description class="pb-4 text-sm sm:text-base">{description}</Card.Description>

								<div>
									<span class="text-2xl font-bold sm:text-3xl">${price}</span>
									<span class="text-sm text-muted-foreground sm:text-base"> /month</span>
								</div>
							</Card.Header>

							<Card.Content>
								<div class="space-y-3 sm:space-y-4">
									{#each benefitList as benefit, bi (bi)}
										<span class="flex items-center text-sm sm:text-base">
											<Check class="mr-2 h-4 w-4 text-primary sm:h-5 sm:w-5" />
											<span>{benefit}</span>
										</span>
									{/each}
								</div>
							</Card.Content>

							<Card.Footer>
								<Button variant={popular ? 'default' : 'secondary'} class="w-full">
									{buttonText}
								</Button>
							</Card.Footer>
						</Card.Root>
					</div>
				{/each}
			</div>

			<!-- Desktop Layout -->
			<div class="relative hidden gap-8 lg:grid lg:grid-cols-3">
				{#each plans as { title, popular, price, description, buttonText, benefitList }, index (index)}
					<div
						class="relative {popular
							? 'z-10 scale-110 hover:scale-112'
							: 'hover:scale-105'} transition-all duration-200"
					>
						<Card.Root class={popular ? 'popular-desktop' : 'hover:border-primary/50'}>
							{#if popular}
								<div class="absolute -top-4 left-1/2 -translate-x-1/2">
									<span
										class="rounded-full bg-primary px-4 py-1 text-sm font-medium whitespace-nowrap text-primary-foreground shadow-lg"
									>
										Most Popular
									</span>
								</div>
							{/if}
							<Card.Header>
								<Card.Title class="pb-2 text-lg sm:text-xl">
									{title}
								</Card.Title>

								<Card.Description class="pb-4 text-sm sm:text-base">{description}</Card.Description>

								<div>
									<span class="text-2xl font-bold sm:text-3xl">${price}</span>
									<span class="text-sm text-muted-foreground sm:text-base"> /month</span>
								</div>
							</Card.Header>

							<Card.Content>
								<div class="space-y-3 sm:space-y-4">
									{#each benefitList as benefit, bi (bi)}
										<span class="flex items-center text-sm sm:text-base">
											<Check class="mr-2 h-4 w-4 text-primary sm:h-5 sm:w-5" />
											<span>{benefit}</span>
										</span>
									{/each}
								</div>
							</Card.Content>

							<Card.Footer>
								<Button variant={popular ? 'default' : 'secondary'} class="w-full">
									{buttonText}
								</Button>
							</Card.Footer>
						</Card.Root>
					</div>
				{/each}
			</div>
		</section>
	</HomeAnimationContainer>
</MaxWidthWrapper>

<!-- svelte-ignore css_unused_selector-->
<style scoped>
	@reference "tailwindcss";
	.popular-mobile {
		@apply relative my-8 scale-[1.02] border-2 border-primary bg-card shadow-md ring-2 ring-primary/50 ring-offset-2 ring-offset-background;
	}

	.popular-desktop {
		@apply relative transform-gpu border-2 border-primary bg-card shadow-2xl ring-2 ring-primary/50 ring-offset-4 ring-offset-background;
	}

	:global(.dark) .popular-desktop {
		@apply shadow-primary/20 ring-primary/40;
	}
</style>
