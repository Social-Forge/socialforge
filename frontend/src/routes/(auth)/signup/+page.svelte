<script lang="ts">
	import { goto, invalidateAll } from '$app/navigation';
	import * as Card from '$lib/components/ui/card/index.js';
  import { Button } from '$lib/components/ui/button/index.js';
  import { SocialAuthButtons } from '$lib/components/auth';
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
  const oauthError = $derived(data.oauthError ?? null);

	let phoneInput = $state<string | undefined>('');
  let errorMessage = $state<string | null>(null);

  $effect(() => {
          if (oauthError) {
                  errorMessage = oauthError;
          }
  });

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
                                   <SocialAuthButtons mode="signup" redirectTo={data.redirectTarget} />
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
