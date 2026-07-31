<script lang="ts">
	import * as Card from '$lib/components/ui/card';
	import { MaxWidthWrapper, HomeAnimationContainer } from '$lib/components/landing';
	import { Button } from '$lib/components/ui/button';
	import { DiscordIcon } from '$lib/components/icons';
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
		<section id="community" class="container mx-auto max-w-max px-4 py-16">
			<Card.Root
				class="flex flex-col items-center justify-center border-0 bg-background text-center shadow-none"
			>
				<Card.Header class="flex w-full flex-col items-center justify-center gap-4">
					<Card.Title class="text-4xl font-bold md:text-5xl">
						<DiscordIcon class="m-auto mb-4 h-20 w-20" />
						Ready to join this
						<span class="bg-linear-to-r from-[#D247BF] to-primary bg-clip-text text-transparent">
							Community?
						</span>
					</Card.Title>
				</Card.Header>
				<Card.Content class="text-xl text-muted-foreground lg:w-[80%]">
					Join our vibrant Discord community! Connect, share, and grow with like-minded enthusiasts.
					Click to dive in! 🚀
				</Card.Content>

				<Card.Footer>
					<Button>
						<a href="https://discord.com/" target="_blank" rel="noreferrer" class="h-full w-full">
							Join Discord
						</a>
					</Button>
				</Card.Footer>
			</Card.Root>
		</section>
	</HomeAnimationContainer>
</MaxWidthWrapper>
