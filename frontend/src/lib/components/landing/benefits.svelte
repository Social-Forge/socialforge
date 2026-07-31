<script lang="ts" module>
	interface BenefitsProps {
		icon: 'sparkle' | 'tag' | 'paintbrush' | 'blocks' | 'lineChart' | 'wallet';
		title: string;
		description: string;
	}
</script>

<script lang="ts">
	import * as Card from '$lib/components/ui/card';
	import { MaxWidthWrapper, HomeAnimationContainer } from '$lib/components/landing';
	import { Sparkle, Tag, Paintbrush, Blocks, LineChart, Wallet } from '@lucide/svelte';

	const benefitList: BenefitsProps[] = [
		{
			icon: 'blocks',
			title: 'Build Brand Trust',
			description:
				'Lorem ipsum dolor sit amet consectetur adipisicing elit. A odio velit cum aliquam. Natus consectetur dolores.'
		},
		{
			icon: 'lineChart',
			title: 'More Leads',
			description:
				'Lorem ipsum dolor sit amet consectetur adipisicing elit. A odio velit cum aliquam, natus consectetur.'
		},
		{
			icon: 'wallet',
			title: 'Higher Conversions',
			description:
				'Lorem ipsum dolor sit amet consectetur adipisicing elit. Natus consectetur. A odio velit cum aliquam'
		},
		{
			icon: 'sparkle',
			title: 'Test Marketing Ideas',
			description:
				'Lorem ipsum dolor sit amet consectetur adipisicing elit. A odio velit cum aliquam. Natus consectetur dolores.'
		}
	];

	const iconMap = {
		sparkle: Sparkle,
		tag: Tag,
		paintbrush: Paintbrush,
		blocks: Blocks,
		lineChart: LineChart,
		wallet: Wallet
	};
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
		<section id="benefits" class="container py-24 sm:py-32">
			<div class="grid place-items-center lg:grid-cols-2 lg:gap-24">
				<div>
					<h2 class="mb-2 text-lg tracking-wider text-primary">Benefits</h2>
					<h2 class="mb-4 text-3xl font-bold md:text-4xl">Your Shortcut to Success</h2>
					<p class="mb-8 text-xl text-muted-foreground">
						Lorem ipsum dolor sit amet consectetur, adipisicing elit. Non ducimus reprehenderit
						architecto rerum similique facere odit deleniti necessitatibus quo quae.
					</p>
				</div>

				<div class="grid w-full gap-4 lg:grid-cols-2">
					{#each benefitList as { icon, title, description }, index (index)}
						<Card.Root
							class="group/number bg-muted/50 transition-all delay-75 hover:bg-background dark:bg-card dark:hover:bg-background"
						>
							<Card.Header>
								<div class="flex justify-between">
									{#if iconMap[icon]}
										{@const Icon = iconMap[icon]}
										<Icon class="mb-6 size-8 text-primary" />
									{/if}
									<span
										class="text-5xl font-medium text-muted-foreground/15 transition-all delay-75 group-hover/number:text-muted-foreground/30"
									>
										0{index + 1}
									</span>
								</div>
								<Card.Title>{title}</Card.Title>
							</Card.Header>

							<Card.Content class="text-muted-foreground">
								{description}
							</Card.Content>
						</Card.Root>
					{/each}
				</div>
			</div>
		</section>
	</HomeAnimationContainer>
</MaxWidthWrapper>
