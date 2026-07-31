<script lang="ts" module>
	const ProService = {
		YES: 1,
		NO: 0
	} as const;

	type ProServiceType = (typeof ProService)[keyof typeof ProService];

	interface ServiceProps {
		title: string;
		pro: ProServiceType;
		description: string;
	}
</script>

<script lang="ts">
	import * as Card from '$lib/components/ui/card';
	import { Badge } from '$lib/components/ui/badge';
	import { MaxWidthWrapper, HomeAnimationContainer } from '$lib/components/landing';

	const serviceList: ServiceProps[] = [
		{
			title: 'Custom Domain Integration',
			description: 'Lorem ipsum dolor sit, amet consectetur adipisicing elit adipisicing.',
			pro: 0
		},
		{
			title: 'Social Media Integrations',
			description: 'Lorem ipsum dolor sit, amet consectetur adipisicing elit. Molestiae, dicta.',
			pro: 0
		},
		{
			title: 'Email Marketing Integrations',
			description: 'Lorem dolor sit amet adipisicing.',
			pro: 0
		},
		{
			title: 'SEO Optimization',
			description: 'Lorem ipsum dolor sit amet consectetur.',
			pro: 1
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
		<section id="services" class="container py-24 sm:py-32">
			<h2 class="mb-2 text-center text-lg tracking-wider text-primary">Services</h2>

			<h2 class="mb-4 text-center text-3xl font-bold md:text-4xl">Grow Your Business</h2>
			<h3 class="mx-auto mb-8 text-center text-xl text-muted-foreground md:w-1/2">
				From marketing and sales to operations and strategy, we have the expertise to help you
				achieve your goals.
			</h3>

			<div class="mx-auto grid w-full gap-4 sm:grid-cols-2 lg:w-[60%] lg:grid-cols-2">
				{#each serviceList as { title, description, pro }, index (index)}
					<Card.Root class="relative h-full bg-muted/60 pb-6 dark:bg-card">
						<Card.Header>
							<Card.Title>{title}</Card.Title>
							<Card.Description>{description}</Card.Description>
						</Card.Header>
						{#if pro === ProService.YES}
							<Badge variant="secondary" class="absolute -top-2 -right-3">PRO</Badge>
						{/if}
					</Card.Root>
				{/each}
			</div>
		</section>
	</HomeAnimationContainer>
</MaxWidthWrapper>
