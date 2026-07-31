<script lang="ts" module>
	interface ReviewProps {
		image: string;
		name: string;
		userName: string;
		comment: string;
		rating: number;
	}
</script>

<script lang="ts">
	import * as Card from '$lib/components/ui/card';
	import * as Avatar from '$lib/components/ui/avatar';
	import * as Carousel from '$lib/components/ui/carousel';
	import { MaxWidthWrapper, HomeAnimationContainer } from '$lib/components/landing';
	import { Star } from '@lucide/svelte';
	import type { CarouselAPI } from '$lib/components/ui/carousel/context';
	import Autoplay from 'embla-carousel-autoplay';

	const reviewList: ReviewProps[] = [
		{
			image: 'https://github.com/shadcn.png',
			name: 'Sarah Chen',
			userName: 'Frontend Developer',
			comment:
				'This SvelteKit landing page template by Memet Zx is exactly what I needed! The conversion from Vue to Svelte is seamless and the components are well-organized.',
			rating: 5.0
		},
		{
			image: 'https://github.com/zxce3.png',
			name: 'Memet Zx',
			userName: 'Creator & Developer',
			comment:
				'I created this template to help developers quickly build beautiful landing pages with SvelteKit and Shadcn. Hope you find it useful!',
			rating: 5.0
		},
		{
			image: 'https://github.com/shadcn.png',
			name: 'Alex Rivera',
			userName: 'Full Stack Developer',
			comment:
				"Zxce3's implementation of Shadcn components in SvelteKit is brilliant. The dark mode feature and responsive design work flawlessly.",
			rating: 4.9
		},
		{
			image: 'https://github.com/shadcn.png',
			name: 'Emily Watson',
			userName: 'UI/UX Designer',
			comment:
				'The attention to detail in this template is impressive. Memet has done an excellent job maintaining the design aesthetics while converting to SvelteKit.',
			rating: 5.0
		},
		{
			image: 'https://github.com/shadcn.png',
			name: 'David Kim',
			userName: 'Web Developer',
			comment:
				'This is now my go-to template for SvelteKit projects. The documentation is clear and the implementation by Zxce3 is top-notch.',
			rating: 4.9
		},
		{
			image: 'https://github.com/shadcn.png',
			name: 'Lisa Chen',
			userName: 'Software Engineer',
			comment:
				"Thanks to Memet's template, I was able to launch my landing page in record time. The TypeScript integration is particularly well done.",
			rating: 5.0
		}
	];

	let api = $state<CarouselAPI>();
	const plugin = Autoplay({ delay: 4000, stopOnInteraction: true });
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
		<section id="testimonials" class="container py-24 sm:py-32">
			<div class="mb-8 text-center">
				<h2 class="mb-2 text-center text-lg tracking-wider text-primary">Testimonials</h2>

				<h2 class="mb-4 text-center text-3xl font-bold md:text-4xl">
					Hear What Our 1000+ Clients Say
				</h2>
			</div>

			<Carousel.Root
				opts={{
					align: 'start',
					loop: true
				}}
				plugins={[plugin]}
				class="relative mx-auto w-[80%] sm:w-[90%] lg:max-w-7xl"
				setApi={(emblaApi) => (api = emblaApi)}
				onmouseenter={plugin.stop}
				onmouseleave={plugin.reset}
			>
				<Carousel.Content class="-ml-4">
					{#each reviewList as review (review.name)}
						<Carousel.Item class="pl-4 md:basis-1/2 lg:basis-1/3">
							<Card.Root class="h-full bg-muted/50 dark:bg-card">
								<Card.Content class="flex h-full flex-col p-6">
									<div class="flex-1">
										<div class="mb-6 flex gap-1">
											{#each Array(5) as _, i (i)}
												<Star class="size-4 fill-primary text-primary" />
											{/each}
										</div>

										<p class="mb-6">"{review.comment}"</p>
									</div>

									<div class="flex items-center gap-4 border-t pt-6">
										<Avatar.Root>
											<Avatar.Image src={review.image} alt={`Avatar of ${review.name}`} />
											<Avatar.Fallback>
												{review.name
													.split(' ')
													.map((n) => n[0])
													.join('')}
											</Avatar.Fallback>
										</Avatar.Root>

										<div class="flex flex-col">
											<Card.Title class="text-lg">{review.name}</Card.Title>
											<Card.Description>{review.userName}</Card.Description>
										</div>
									</div>
								</Card.Content>
							</Card.Root>
						</Carousel.Item>
					{/each}
				</Carousel.Content>
				<Carousel.Previous />
				<Carousel.Next />
			</Carousel.Root>
		</section>
	</HomeAnimationContainer>
</MaxWidthWrapper>
