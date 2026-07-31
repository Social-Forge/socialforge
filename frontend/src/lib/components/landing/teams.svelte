<script lang="ts" module>
	interface TeamProps {
		imageUrl: string;
		firstName: string;
		lastName: string;
		positions: string[];
		socialNetworks: SocialNetworkProps[];
	}

	interface SocialNetworkProps {
		name: 'LinkedIn' | 'Github' | 'X';
		url: string;
	}
</script>

<script lang="ts">
	import * as Card from '$lib/components/ui/card';
	import { MaxWidthWrapper, HomeAnimationContainer } from '$lib/components/landing';
	import { LinkedinIcon, GithubIcon, XIcon } from '$lib/components/icons';

	const teamList: TeamProps[] = [
		{
			imageUrl: 'https://i.pravatar.cc/250?img=58',
			firstName: 'Leo',
			lastName: 'Miranda',
			positions: ['Vue Fronted Developer', 'Creator Of This Website'],
			socialNetworks: [
				{
					name: 'LinkedIn',
					url: 'https://www.linkedin.com/in/leopoldo-miranda/'
				},
				{
					name: 'Github',
					url: 'https://github.com/leoMirandaa'
				},
				{
					name: 'X',
					url: 'https://x.com/leo_mirand4'
				}
			]
		},
		{
			imageUrl:
				'https://images.unsplash.com/photo-1534528741775-53994a69daeb?q=80&w=1528&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D',
			firstName: 'Elizabeth',
			lastName: 'Moore',
			positions: ['UI/UX Designer'],
			socialNetworks: [
				{
					name: 'LinkedIn',
					url: 'https://www.linkedin.com/in/leopoldo-miranda/'
				},
				{
					name: 'X',
					url: 'https://x.com/leo_mirand4'
				}
			]
		},
		{
			imageUrl:
				'https://images.unsplash.com/photo-1527980965255-d3b416303d12?q=80&w=1760&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D',
			firstName: 'David',
			lastName: 'Diaz',
			positions: ['Machine Learning Engineer', 'TensorFlow Tinkerer'],
			socialNetworks: [
				{
					name: 'LinkedIn',
					url: 'https://www.linkedin.com/in/leopoldo-miranda/'
				},
				{
					name: 'Github',
					url: 'https://github.com/leoMirandaa'
				}
			]
		},
		{
			imageUrl:
				'https://images.unsplash.com/photo-1573497161161-c3e73707e25c?q=80&w=1587&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D',
			firstName: 'Sarah',
			lastName: 'Robinson',
			positions: ['Cloud Native Developer', ' Kubernetes Orchestrator'],
			socialNetworks: [
				{
					name: 'LinkedIn',
					url: 'https://www.linkedin.com/in/leopoldo-miranda/'
				},
				{
					name: 'Github',
					url: 'https://github.com/leoMirandaa'
				},
				{
					name: 'X',
					url: 'https://x.com/leo_mirand4'
				}
			]
		},
		{
			imageUrl:
				'https://images.unsplash.com/photo-1616805765352-beedbad46b2a?q=80&w=1887&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D',
			firstName: 'Michael',
			lastName: 'Holland',
			positions: ['DevOps Engineer', 'CI/CD Pipeline Mastermind'],
			socialNetworks: [
				{
					name: 'LinkedIn',
					url: 'https://www.linkedin.com/in/leopoldo-miranda/'
				}
			]
		},
		{
			imageUrl:
				'https://images.unsplash.com/photo-1573497019940-1c28c88b4f3e?q=80&w=1587&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D',
			firstName: 'Zoe',
			lastName: 'Garcia',
			positions: ['JavaScript Evangelist', 'Deno Champion'],
			socialNetworks: [
				{
					name: 'LinkedIn',
					url: 'https://www.linkedin.com/in/leopoldo-miranda/'
				},
				{
					name: 'Github',
					url: 'https://github.com/leoMirandaa'
				}
			]
		},
		{
			imageUrl:
				'https://images.unsplash.com/photo-1633332755192-727a05c4013d?q=80&w=1480&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D',
			firstName: 'Evan',
			lastName: 'James',
			positions: ['Backend Developer'],
			socialNetworks: [
				{
					name: 'LinkedIn',
					url: 'https://www.linkedin.com/in/leopoldo-miranda/'
				},
				{
					name: 'Github',
					url: 'https://github.com/leoMirandaa'
				},
				{
					name: 'X',
					url: 'https://x.com/leo_mirand4'
				}
			]
		},
		{
			imageUrl:
				'https://images.unsplash.com/photo-1573497019236-17f8177b81e8?q=80&w=1740&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3Dhttps://images.unsplash.com/photo-1573497019236-17f8177b81e8?q=80&w=1740&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D',
			firstName: 'Pam',
			lastName: 'Taylor',
			positions: ['Fullstack developer', 'UX Researcher'],
			socialNetworks: [
				{
					name: 'X',
					url: 'https://x.com/leo_mirand4'
				}
			]
		}
	];

	const socialIconMap = {
		LinkedIn: LinkedinIcon,
		Github: GithubIcon,
		X: XIcon
	};
</script>

<MaxWidthWrapper>
	<section id="team" class="container mx-auto max-w-7xl px-4 py-16">
		<div class="mb-12 text-center">
			<h2 class="mb-3 text-center text-lg tracking-wider text-primary">Team</h2>

			<h2 class="text-center text-3xl font-bold md:text-4xl">The Company Dream Team</h2>
		</div>

		<div
			class="grid grid-cols-1 items-start justify-center gap-6 sm:grid-cols-2 md:grid-cols-3 md:gap-8 xl:grid-cols-4"
		>
			{#each teamList as { imageUrl, firstName, lastName, positions, socialNetworks }, pi (pi)}
				<Card.Root
					class="group/hoverimg flex h-full flex-col overflow-hidden bg-muted/60 dark:bg-card"
				>
					<Card.Header class="gap-0 p-0">
						<div class="h-full overflow-hidden">
							<img
								src={imageUrl}
								alt={`${firstName} ${lastName}`}
								class="aspect-square w-full object-cover saturate-0 transition-all duration-300 ease-in-out group-hover/hoverimg:scale-105 group-hover/hoverimg:saturate-100"
							/>
						</div>
						<Card.Title class="p-4 sm:p-6">
							{firstName}
							<span class="text-primary">{lastName}</span>
						</Card.Title>
					</Card.Header>

					<div class="grow px-4 sm:px-6">
						{#each positions as position, si (si)}
							<Card.Content class="p-0 pb-2 text-muted-foreground">
								{position}{#if si < positions.length - 1},{/if}
							</Card.Content>
						{/each}
					</div>

					<Card.Footer class="space-x-4 px-4 py-4 sm:px-6">
						{#each socialNetworks as { name, url }, sni (sni)}
							<a
								href={url}
								target="_blank"
								class="transition-opacity hover:opacity-80"
								aria-label="Visit our {name} page"
							>
								{#if socialIconMap[name]}
									{@const Icon = socialIconMap[name]}
									<Icon />
								{/if}
							</a>
						{/each}
					</Card.Footer>
				</Card.Root>
			{/each}
		</div>
	</section>
</MaxWidthWrapper>
