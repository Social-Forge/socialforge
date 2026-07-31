<script lang="ts">
	import { goto, invalidateAll } from '$app/navigation';
	import { superForm } from 'sveltekit-superforms';
	import * as Card from '$lib/components/ui/card/index.js';
	import * as Field from '$lib/components/ui/field/index.js';
	import * as Password from '$lib/components/extras/password';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Spinner } from '$lib/components/ui/spinner/index.js';
	import type { ZxcvbnResult } from '@zxcvbn-ts/core';
	import Icon from '@iconify/svelte';
	import * as Alert from '$lib/components/ui/alert/index.js';
	import { localizeHref } from '$lib/paraglide/runtime.js';
	import * as i18n from '$lib/paraglide/messages.js';

	let { data } = $props();
	let metaTags = $derived(data.pageMetaTags);

	let showConfirmPassword = $state(false);
	let errorMessage = $state<string | undefined>(undefined);
	let passwordInput = $state<string | undefined>('');
	const SCORE_NAMING = ['Poor', 'Weak', 'Average', 'Strong', 'Secure'];
	let strength = $state<ZxcvbnResult>();

	// svelte-ignore state_referenced_locally
	const { form, enhance, errors, submitting } = superForm(data.form, {
		async onSubmit(input) {
			errorMessage = undefined;
		},
		async onUpdate(event) {
			if (event.result.type === 'failure') {
				errorMessage = event.result.data.message;
				return;
			}
			await invalidateAll();
			await goto(localizeHref('/signin'));
		}
	});

	$effect(() => {
		if (data.token.trim().length === 0) {
			errorMessage = 'Token is required';
			return;
		}
		$form.token = data.token;
	});
</script>

<div class="flex flex-col gap-6">
	<Card.Root>
		<Card.Header class="text-center">
			<Card.Title class="text-xl">{i18n.auth_reset_password()}</Card.Title>
			<Card.Description>{i18n.auth_reset_password_description()}</Card.Description>
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
					<Input
						bind:value={$form.token}
						name="token"
						type="hidden"
						aria-invalid={!!$errors.token}
						autocomplete="on"
					/>
					<Field.Field>
						<Field.Label for="new_password">{i18n.field_password()}</Field.Label>
						<Password.Root>
							<Password.Input
								name="new_password"
								placeholder="Password"
								required
								autocomplete="new-password"
								disabled={$submitting}
								oninput={(e) => ($form.new_password = (e.target as HTMLInputElement)?.value || '')}
							>
								<Password.ToggleVisibility />
							</Password.Input>
							<Password.Strength />
						</Password.Root>
						{#if $errors.new_password}
							<Field.Error>{$errors.new_password}</Field.Error>
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
							{$submitting ? i18n.auth_resetting_password() : i18n.auth_reset_password()}
						</Button>
					</Field.Field>
				</Field.Group>
			</form>
		</Card.Content>
	</Card.Root>
	<Field.Description class="px-6 text-center">
		{i18n.auth_terms_conditions()}
	</Field.Description>
</div>
