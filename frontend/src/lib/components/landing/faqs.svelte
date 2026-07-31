<script lang="ts" module>
	interface FAQProps {
		question: string;
		answer: string;
		value: string;
	}
</script>

<script lang="ts">
	import { MaxWidthWrapper, HomeAnimationContainer } from '$lib/components/landing';
	import * as Accordion from '$lib/components/ui/accordion';

	const FAQList: FAQProps[] = [
		{
			question: 'Is this template free?',
			answer: 'Yes. It is a free Shadcn/Svelte template.',
			value: 'item-1'
		},
		{
			question: 'Duis aute irure dolor in reprehenderit in voluptate velit?',
			answer:
				'Lorem ipsum dolor sit amet, consectetur adipisicing elit. Sint labore quidem quam consectetur sapiente, iste rerum reiciendis animi nihil nostrum sit quo, modi quod.',
			value: 'item-2'
		},
		{
			question: 'Lorem ipsum dolor sit amet Consectetur natus dolor minus quibusdam?',
			answer:
				'Lorem ipsum dolor sit amet consectetur, adipisicing elit. Labore qui nostrum reiciendis veritatis.',
			value: 'item-3'
		},
		{
			question: 'Excepteur sint occaecat cupidata non proident sunt?',
			answer: 'Lorem ipsum dolor sit amet consectetur, adipisicing elit.',
			value: 'item-4'
		},
		{
			question: 'Enim ad minim veniam, quis nostrud exercitation ullamco laboris?',
			answer: 'consectetur adipisicing elit. Sint labore.',
			value: 'item-5'
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
		<section id="faq" class="container w-full py-24">
			<div class="mb-8 text-center">
				<h2 class="mb-2 text-center text-lg tracking-wider text-primary">FAQS</h2>
				<h2 class="text-center text-3xl font-bold md:text-4xl">Common Questions</h2>
			</div>

			<div class="relative mx-auto lg:w-[80%]">
				<Accordion.Root type="single" class="w-full">
					{#each FAQList as { question, answer, value }, index (index)}
						<Accordion.Item {value}>
							<Accordion.Trigger class="text-left">
								{question}
							</Accordion.Trigger>
							<Accordion.Content>
								{answer}
							</Accordion.Content>
						</Accordion.Item>
					{/each}
				</Accordion.Root>
			</div>
			<h3 class="mt-4 text-center font-medium">
				Still have questions?
				<a href="#contact" class="text-muted-foreground">
					<span class="underline">Contact us</span>
				</a>
			</h3>
		</section>
	</HomeAnimationContainer>
</MaxWidthWrapper>
