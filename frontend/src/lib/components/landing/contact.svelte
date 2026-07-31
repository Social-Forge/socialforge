<script lang="ts" module>
	interface ContactFormProps {
		firstName: string;
		lastName: string;
		email: string;
		subject: string;
		message: string;
	}
</script>

<script lang="ts">
	import * as Card from '$lib/components/ui/card';
	import { Button } from '$lib/components/ui/button';
	import { MaxWidthWrapper, HomeAnimationContainer } from '$lib/components/landing';
	import { Label } from '$lib/components/ui/label';
	import { Input } from '$lib/components/ui/input';
	import { Textarea } from '$lib/components/ui/textarea';
	import * as Alert from '$lib/components/ui/alert';
	import { AlertCircle, Building2, Phone, Mail, Clock } from '@lucide/svelte';
	import * as Select from '$lib/components/ui/select/index.js';

	let contactForm = $state<ContactFormProps>({
		firstName: '',
		lastName: '',
		email: '',
		subject: 'Web Development',
		message: ''
	});

	let invalidInputForm = false;

	function handleSubmit(event: SubmitEvent) {
		event.preventDefault();
		const { firstName, lastName, email, subject, message } = contactForm;
		console.log(contactForm);

		const mailToLink = `mailto:leomirandadev@gmail.com?subject=${subject}&body=Hello I am ${firstName} ${lastName}, my Email is ${email}. %0D%0A${message}`;
		window.location.href = mailToLink;
	}

	const subjects = [
		{ value: 'Web Development', label: 'Web Development' },
		{ value: 'Mobile Development', label: 'Mobile Development' },
		{ value: 'Figma Design', label: 'Figma Design' },
		{ value: 'REST API', label: 'REST API' },
		{ value: 'FullStack Project', label: 'FullStack Project' }
	];

	let triggerContent = $derived(
		subjects.find((s) => s.value === contactForm.subject)?.label ?? 'Select a subject'
	);
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
		<section id="contact" class="container py-24 sm:py-32">
			<section class="grid grid-cols-1 gap-8 md:grid-cols-2">
				<div>
					<div class="mb-4">
						<h2 class="mb-2 text-lg tracking-wider text-primary">Contact</h2>
						<h2 class="text-3xl font-bold md:text-4xl">Connect With Us</h2>
					</div>
					<p class="mb-8 text-muted-foreground lg:w-5/6">
						Lorem ipsum dolor sit amet consectetur adipisicing elit. Voluptatum ipsam sint enim
						exercitationem ex autem corrupti quas tenetur
					</p>

					<div class="flex flex-col gap-4">
						<div>
							<div class="mb-1 flex gap-2">
								<Building2 />
								<div class="font-bold">Find Us</div>
							</div>
							<div>742 Evergreen Terrace, Springfield, IL 62704</div>
						</div>

						<div>
							<div class="mb-1 flex gap-2">
								<Phone />
								<div class="font-bold">Call Us</div>
							</div>
							<div>+1 (619) 123-4567</div>
						</div>

						<div>
							<div class="mb-1 flex gap-2">
								<Mail />
								<div class="font-bold">Mail Us</div>
							</div>
							<div>leomirandadev@gmail.com</div>
						</div>

						<div>
							<div class="mb-1 flex gap-2">
								<Clock />
								<div class="font-bold">Visit Us</div>
							</div>
							<div>
								<div>Monday - Friday</div>
								<div>8AM - 4PM</div>
							</div>
						</div>
					</div>
				</div>

				<!-- Form -->
				<Card.Root class="bg-muted/60 dark:bg-card">
					<Card.Header class="text-2xl text-primary" />
					<Card.Content>
						<form onsubmit={handleSubmit} class="grid gap-4">
							<div class="flex flex-col gap-8 md:flex-row">
								<div class="flex w-full flex-col gap-1.5">
									<Label for="firstName">First Name</Label>
									<Input
										id="firstName"
										type="text"
										placeholder="Leopoldo"
										bind:value={contactForm.firstName}
									/>
								</div>

								<div class="flex w-full flex-col gap-1.5">
									<Label for="lastName">Last Name</Label>
									<Input
										id="lastName"
										type="text"
										placeholder="Miranda"
										bind:value={contactForm.lastName}
									/>
								</div>
							</div>

							<div class="flex flex-col gap-1.5">
								<Label for="contactEmail">Email</Label>
								<Input
									id="contactEmail"
									type="email"
									placeholder="leomirandadev@gmail.com"
									bind:value={contactForm.email}
								/>
							</div>

							<div class="flex flex-col gap-1.5">
								<Label for="contactSubject">Subject</Label>
								<Select.Root type="single" bind:value={contactForm.subject}>
									<Select.Trigger id="contactSubject" class="w-full">
										{triggerContent}
									</Select.Trigger>
									<Select.Content>
										<Select.Group>
											{#each subjects as subject, si (si)}
												<Select.Item value={subject.value} label={subject.label}>
													{subject.label}
												</Select.Item>
											{/each}
										</Select.Group>
									</Select.Content>
								</Select.Root>
							</div>

							<div class="flex flex-col gap-1.5">
								<Label for="contactMessage">Message</Label>
								<Textarea
									id="contactMessage"
									placeholder="Your message..."
									rows={5}
									bind:value={contactForm.message}
								/>
							</div>

							{#if invalidInputForm}
								<Alert.Root variant="destructive">
									<AlertCircle class="h-4 w-4" />
									<Alert.Title>Error</Alert.Title>
									<Alert.Description>
										There is an error in the form. Please check your input.
									</Alert.Description>
								</Alert.Root>
							{/if}

							<Button class="mt-4">Send message</Button>
						</form>
					</Card.Content>
					<Card.Footer />
				</Card.Root>
			</section>
		</section>
	</HomeAnimationContainer>
</MaxWidthWrapper>
