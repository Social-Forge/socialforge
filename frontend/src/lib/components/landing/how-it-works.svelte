<script lang="ts" module>
	interface HowItWorksProps {
		badgeTitle: string;
		title: string;
		description: string;
		image: string;
	}
</script>

<script lang="ts">
	import * as Card from '$lib/components/ui/card';
	import { Badge } from '$lib/components/ui/badge';
	import { MaxWidthWrapper, HomeAnimationContainer } from '$lib/components/landing';

	const HowItWorksList: HowItWorksProps[] = [
		{
			badgeTitle: 'Contribute',
			title: 'Make a Difference with Every Contribution',
			description:
				'Lorem ipsum dolor, sit amet consectetur adipisicing elit. Voluptatem fugiat, quasi sint reiciendis quidem iure.',
			image: '/images/roboto.png'
		},
		{
			badgeTitle: 'Discover',
			title: 'Create an Impact with Every Choice',
			description:
				'Lorem ipsum dolor, sit amet consectetur adipisicing elit. Voluptatem fugiat, sint reiciendis quidem iure veritatis optio facere tenetur.',
			image: '/images/runner.png'
		},
		{
			badgeTitle: 'Innovate',
			title: 'Make Every Action Count',
			description:
				'Lorem ipsum dolor, sit amet consectetur adipisicing elit. Voluptatem fugiat, odit similique quasi sint.',
			image: '/images/pacheco.png'
		},
		{
			badgeTitle: 'Collaborate',
			title: 'Tackle this challenge together',
			description:
				'Lorem ipsum dolor, sit amet consectetur adipisicing elit. Voluptatem fugiat, quasi sint reiciendis quidem iure.',
			image: '/images/gamestation.png'
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
		<section id="features" class="container py-24">
			<div class="mb-8 text-center">
				<h2 class="mb-2 text-center text-lg tracking-wider text-primary">How It Works</h2>

				<h2 class="text-center text-3xl font-bold md:text-4xl">Step-by-Step Process</h2>
			</div>

			<div class="relative mx-auto lg:w-[80%]">
				{#each HowItWorksList as { badgeTitle, title, description, image }, index (index)}
					<div class="mb-8 flex items-center" class:flex-row-reverse={index % 2 !== 0}>
						<Card.Root class="h-full border-0 bg-transparent shadow-none">
							<Card.Header>
								<div class="pb-4">
									<Badge>{badgeTitle}</Badge>
								</div>

								<Card.Title>
									{title}
								</Card.Title>
							</Card.Header>

							<Card.Content class="w-[80%] text-muted-foreground">
								{description}
							</Card.Content>
						</Card.Root>

						<img
							src={image}
							alt="Image describing {title}"
							class="mx-auto w-37.5 -scale-x-100 md:w-62.5 lg:w-75"
						/>

						<div
							class="absolute right-0 -z-10 h-72 w-44 rounded-full bg-primary/15 blur-3xl lg:h-80 lg:w-64 dark:bg-primary/10"
							class:left-0={index % 2 !== 0}
						></div>
					</div>
				{/each}
			</div>
		</section>
	</HomeAnimationContainer>
</MaxWidthWrapper>
