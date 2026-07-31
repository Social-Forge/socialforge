<script lang="ts" module>
	interface SponsorsProps {
		icon: 'crown' | 'vegan' | 'ghost' | 'puzzle' | 'squirrel' | 'cookie' | 'drama';
		name: string;
	}
</script>

<script lang="ts">
	import { MaxWidthWrapper, HomeAnimationContainer } from '$lib/components/landing';
	import { Crown, Vegan, Ghost, Puzzle, Squirrel, Cookie, Drama } from '@lucide/svelte';

	const sponsors: SponsorsProps[] = [
		{ icon: 'crown', name: 'Acmebrand' },
		{ icon: 'vegan', name: 'Acmelogo' },
		{ icon: 'ghost', name: 'Acmesponsor' },
		{ icon: 'puzzle', name: 'Acmeipsum' },
		{ icon: 'squirrel', name: 'Acme' },
		{ icon: 'cookie', name: 'Accmee' },
		{ icon: 'drama', name: 'Acmetech' }
	];

	const iconMap = {
		crown: Crown,
		vegan: Vegan,
		ghost: Ghost,
		puzzle: Puzzle,
		squirrel: Squirrel,
		cookie: Cookie,
		drama: Drama
	};

	let duplicatedCompanies = [...sponsors, ...sponsors];
</script>

<MaxWidthWrapper>
	<HomeAnimationContainer
		variant="fade-up"
		delay={0.2}
		once={false}
		duration={1}
		distance={60}
		class="relative w-full bg-transparent px-2"
	>
		<section id="sponsors" class="mx-auto pb-24 sm:pb-32">
			<h2 class="mb-6 text-center text-lg md:text-xl">Trusted by the best in the industry</h2>

			<div class="marquee-wrapper mx-auto overflow-hidden">
				<div class="marquee-track flex gap-12">
					{#each duplicatedCompanies as company, index (company.name + index)}
						<div class="flex items-center text-xl font-medium whitespace-nowrap md:text-2xl">
							{#if iconMap[company.icon]}
								{@const Icon = iconMap[company.icon]}
								<Icon class="mr-2" strokeWidth={3} />
							{/if}
							{company.name}
						</div>
					{/each}
				</div>
			</div>
		</section>
	</HomeAnimationContainer>
</MaxWidthWrapper>

<style>
	.marquee-wrapper {
		position: relative;
		width: 100%;
		overflow: hidden;
	}

	.marquee-track {
		display: flex;
		gap: 3rem; /* 12 = 3rem */
		width: max-content;
		animation: marquee-scroll 30s linear infinite;
	}

	.marquee-wrapper:hover .marquee-track {
		animation-play-state: paused;
	}

	@keyframes marquee-scroll {
		0% {
			transform: translateX(0);
		}
		100% {
			transform: translateX(-50%);
		}
	}
</style>
