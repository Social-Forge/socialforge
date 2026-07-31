<script lang="ts">
	import { goto, invalidateAll } from '$app/navigation';
	import * as Card from '$lib/components/ui/card/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import * as Field from '$lib/components/ui/field/index.js';
	import { superForm } from 'sveltekit-superforms';
	import { Input } from '$lib/components/ui/input/index.js';
	import * as Password from '$lib/components/extras/password';
	import { Spinner } from '$lib/components/ui/spinner/index.js';
	import Icon from '@iconify/svelte';
	import * as Alert from '$lib/components/ui/alert/index.js';
	import { PhoneInput } from '$lib/components/extras/phone-input/index.js';
	import { localizeHref } from '$lib/paraglide/runtime.js';
	import * as i18n from '$lib/paraglide/messages.js';

	let { data } = $props();

	let phoneInput = $state<string | undefined>('');
	let errorMessage = $state<string | null>(null);

	// svelte-ignore state_referenced_locally
	const { form, enhance, errors, submitting } = superForm(data.form, {
		async onSubmit(input) {
			errorMessage = null;
		},
		async onUpdate(event) {
			if (event.result.type === 'failure') {
				errorMessage = event.result.data.message;
				return;
			}
			await invalidateAll();
			await goto(localizeHref(`/verify-email?email=${$form.email}`));
		}
	});

	const formatPhoneInput = (event: Event) => {
		let value = (event.target as HTMLInputElement).value;
		value = value.replace(/[^\d+]/g, '');

		if (value.startsWith('+')) {
			value = '+' + value.slice(1).replace(/\+/g, '');
		} else {
			value = value.replace(/\+/g, '');
		}

		(event.target as HTMLInputElement).value = value;
		$form.phone = value;
	};

	$effect(() => {
		if (data.form.data.phone && !phoneInput) {
			phoneInput = data.form.data.phone;
		}
		if (phoneInput && typeof phoneInput === 'string' && phoneInput.trim() !== '') {
			const cleanNumber = phoneInput.replace(/[\s-]/g, '');
			$form.phone = cleanNumber;
		} else {
			$form.phone = '';
		}
	});
</script>

<div class="flex flex-col gap-6">
	<Card.Root>
		<Card.Header class="text-center">
			<Card.Title class="text-xl">{i18n.auth_dign_up_title()}</Card.Title>
			<Card.Description>{i18n.auth_signup_description()}</Card.Description>
		</Card.Header>
		<Card.Content>
			<form method="POST" class="w-full" use:enhance>
				<Field.Group>
					{#if errorMessage}
						<Alert.Root variant="destructive">
							<Icon icon="mingcute:warning-line" class="size-4" />
							<Alert.Title>Error</Alert.Title>
							<Alert.Description>{errorMessage}</Alert.Description>
						</Alert.Root>
					{/if}
					<Field.Field>
						<Button variant="outline" type="button">
							<svg xmlns="http://www.w3.org/2000/svg" width="32" height="32" viewBox="0 0 16 16"
								>><g fill="none" fill-rule="evenodd" clip-rule="evenodd"
									><path
										fill="#f44336"
										d="M7.209 1.061c.725-.081 1.154-.081 1.933 0a6.57 6.57 0 0 1 3.65 1.82a100 100 0 0 0-1.986 1.93q-1.876-1.59-4.188-.734q-1.696.78-2.362 2.528a78 78 0 0 1-2.148-1.658a.26.26 0 0 0-.16-.027q1.683-3.245 5.26-3.86"
										opacity=".987"
									/><path
										fill="#ffc107"
										d="M1.946 4.92q.085-.013.161.027a78 78 0 0 0 2.148 1.658A7.6 7.6 0 0 0 4.04 7.99q.037.678.215 1.331L2 11.116Q.527 8.038 1.946 4.92"
										opacity=".997"
									/><path
										fill="#448aff"
										d="M12.685 13.29a26 26 0 0 0-2.202-1.74q1.15-.812 1.396-2.228H8.122V6.713q3.25-.027 6.497.055q.616 3.345-1.423 6.032a7 7 0 0 1-.51.49"
										opacity=".999"
									/><path
										fill="#43a047"
										d="M4.255 9.322q1.23 3.057 4.51 2.854a3.94 3.94 0 0 0 1.718-.626q1.148.812 2.202 1.74a6.62 6.62 0 0 1-4.027 1.684a6.4 6.4 0 0 1-1.02 0Q3.82 14.524 2 11.116z"
										opacity=".993"
									/></g
								></svg
							>
							{i18n.auth_button_continue_with_google()}
						</Button>
						<Button variant="outline" type="button">
							<svg xmlns="http://www.w3.org/2000/svg" width="32" height="32" viewBox="0 0 256 256">
								<path
									fill="#1877f2"
									d="M256 128C256 57.308 198.692 0 128 0S0 57.308 0 128c0 63.888 46.808 116.843 108 126.445V165H75.5v-37H108V99.8c0-32.08 19.11-49.8 48.348-49.8C170.352 50 185 52.5 185 52.5V84h-16.14C152.959 84 148 93.867 148 103.99V128h35.5l-5.675 37H148v89.445c61.192-9.602 108-62.556 108-126.445"
								/><path
									fill="#fff"
									d="m177.825 165l5.675-37H148v-24.01C148 93.866 152.959 84 168.86 84H185V52.5S170.352 50 156.347 50C127.11 50 108 67.72 108 99.8V128H75.5v37H108v89.445A129 129 0 0 0 128 256a129 129 0 0 0 20-1.555V165z"
								/></svg
							>
							{i18n.auth_button_continue_with_facebook()}
						</Button>
						<Button variant="outline" type="button">
							<svg xmlns="http://www.w3.org/2000/svg" width="32" height="32" viewBox="0 0 24 24">
								<path
									fill="currentColor"
									d="M12 2A10 10 0 0 0 2 12c0 4.42 2.87 8.17 6.84 9.5c.5.08.66-.23.66-.5v-1.69c-2.77.6-3.36-1.34-3.36-1.34c-.46-1.16-1.11-1.47-1.11-1.47c-.91-.62.07-.6.07-.6c1 .07 1.53 1.03 1.53 1.03c.87 1.52 2.34 1.07 2.91.83c.09-.65.35-1.09.63-1.34c-2.22-.25-4.55-1.11-4.55-4.92c0-1.11.38-2 1.03-2.71c-.1-.25-.45-1.29.1-2.64c0 0 .84-.27 2.75 1.02c.79-.22 1.65-.33 2.5-.33s1.71.11 2.5.33c1.91-1.29 2.75-1.02 2.75-1.02c.55 1.35.2 2.39.1 2.64c.65.71 1.03 1.6 1.03 2.71c0 3.82-2.34 4.66-4.57 4.91c.36.31.69.92.69 1.85V21c0 .27.16.59.67.5C19.14 20.16 22 16.42 22 12A10 10 0 0 0 12 2"
								/></svg
							>
							{i18n.auth_button_continue_with_github()}
						</Button>
					</Field.Field>
					<Field.Separator class="*:data-[slot=field-separator-content]:bg-card">
						{i18n.auth_or_continue_with()}
					</Field.Separator>
					<Field.Field>
						<Field.Label for="full_name">{i18n.field_full_name()}</Field.Label>
						<Input
							bind:value={$form.full_name}
							name="full_name"
							type="text"
							placeholder="Jane Doe"
							required
							aria-invalid={!!$errors.full_name}
							disabled={$submitting}
							autocomplete="name"
						/>
						{#if $errors.full_name}
							<Field.Error>{$errors.full_name}</Field.Error>
						{/if}
					</Field.Field>
					<Field.Field>
						<Field.Label for="email">{i18n.field_email()}</Field.Label>
						<Input
							bind:value={$form.email}
							name="email"
							type="email"
							placeholder="jane@example.com"
							required
							aria-invalid={!!$errors.email}
							disabled={$submitting}
							autocomplete="email"
						/>
						{#if $errors.email}
							<Field.Error>{$errors.email}</Field.Error>
						{/if}
					</Field.Field>
					<Field.Field>
						<Field.Label for="phone">{i18n.field_phone()}</Field.Label>
						<PhoneInput
							bind:value={phoneInput}
							name="phone"
							placeholder="Enter your phone number"
							required
							disabled={$submitting}
						/>
						{#if $errors.phone}
							<Field.Error>{$errors.phone}</Field.Error>
						{/if}
					</Field.Field>
					<Field.Field>
						<Field.Label for="password">{i18n.field_password()}</Field.Label>
						<Password.Root>
							<Password.Input
								name="password"
								placeholder="Password"
								required
								autocomplete="current-password"
								disabled={$submitting}
								oninput={(e) => ($form.password = (e.target as HTMLInputElement)?.value || '')}
							>
								<Password.ToggleVisibility />
							</Password.Input>
							<Password.Strength />
						</Password.Root>
						{#if $errors.password}
							<Field.Error>{$errors.password}</Field.Error>
						{/if}
					</Field.Field>
					<Field.Field>
						<Field.Label for="confirm_password">{i18n.field_confirm_password()}</Field.Label>
						<Password.Root>
							<Password.Input
								name="confirm_password"
								placeholder="Confirm Password"
								required
								autocomplete="current-password"
								disabled={$submitting}
								oninput={(e) =>
									($form.confirm_password = (e.target as HTMLInputElement)?.value || '')}
							>
								<Password.ToggleVisibility />
							</Password.Input>
						</Password.Root>
						{#if $errors.confirm_password}
							<Field.Error>{$errors.confirm_password}</Field.Error>
						{/if}
					</Field.Field>
					<Field.Field>
						<Button type="submit" disabled={$submitting}>
							{#if $submitting}
								<Spinner />
							{/if}
							{$submitting ? i18n.loading_submitting() : i18n.auth_button_signup()}
						</Button>
						<Field.Description class="text-center">
							{i18n.auth_already_have_account()}
							<a href={localizeHref('/signin')}>{i18n.auth_button_login()}</a>
						</Field.Description>
					</Field.Field>
				</Field.Group>
			</form>
		</Card.Content>
	</Card.Root>
	<Field.Description class="px-6 text-center">
		{i18n.auth_terms_conditions()}
	</Field.Description>
</div>
