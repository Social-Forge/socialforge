<script lang="ts">
	import { invalidateAll } from '$app/navigation';
	import { superForm } from 'sveltekit-superforms';
	import { Input } from '$lib/components/ui/input/index.js';
	import * as Card from '$lib/components/ui/card/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import * as Field from '$lib/components/ui/field/index.js';
	import { Spinner } from '$lib/components/ui/spinner/index.js';
	import Icon from '@iconify/svelte';
	import * as Alert from '$lib/components/ui/alert/index.js';
	import { localizeHref } from '$lib/paraglide/runtime.js';
	import * as i18n from '$lib/paraglide/messages.js';

	let { data } = $props();
	let metaTags = $derived(data.pageMetaTags);

	let errorMessage = $state<string | null>(null);
	let successMessage = $state<string | null>(null);

	// svelte-ignore state_referenced_locally
	const { form, enhance, errors, submitting } = superForm(data.form, {
		async onSubmit(input) {
			errorMessage = null;
			successMessage = null;
		},
		async onUpdate(event) {
			if (event.result.type === 'failure') {
				errorMessage = event.result.data.message;
				return;
			}
			successMessage = event.result.data.message;
			await invalidateAll();
		}
	});
</script>

<div class="flex flex-col gap-6">
	<Card.Root>
		<Card.Header class="text-center">
			<Card.Title class="text-xl">{i18n.auth_forgot_password()}</Card.Title>
			<Card.Description>{i18n.auth_forgot_password_description()}</Card.Description>
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
					{#if successMessage}
						<Alert.Root variant="default">
							<Icon icon="mingcute:check-line" class="size-4" />
							<Alert.Title>Success</Alert.Title>
							<Alert.Description>{successMessage}</Alert.Description>
						</Alert.Root>
					{/if}
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
						<Button type="submit" disabled={$submitting}>
							{#if $submitting}
								<Spinner />
							{/if}
							{$submitting ? i18n.auth_sending() : i18n.auth_button_resest_link()}
						</Button>
						<Field.Description class="text-center">
							{i18n.auth_forgot_remember_password()}
							{i18n.auth_button_back_to_login()}
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
