<script lang="ts">
	import { goto, invalidateAll } from '$app/navigation';
	import { superForm } from 'sveltekit-superforms';
	import * as Card from '$lib/components/ui/card/index.js';
	import * as Field from '$lib/components/ui/field/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Spinner } from '$lib/components/ui/spinner/index.js';
	import Icon from '@iconify/svelte';
	import * as Alert from '$lib/components/ui/alert/index.js';
	import * as InputOTP from '$lib/components/ui/input-otp/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { REGEXP_ONLY_DIGITS_AND_CHARS } from 'bits-ui';
	import { localizeHref } from '$lib/paraglide/runtime.js';
	import * as i18n from '$lib/paraglide/messages.js';

	let { data } = $props();
	let metaTags = $derived(data.pageMetaTags);

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
			await goto(localizeHref('/app/chats'));
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
			<Card.Title class="text-xl">{i18n.auth_two_factor_authentication()}</Card.Title>
			<Card.Description>{i18n.auth_two_factor_authentication_description()}</Card.Description>
		</Card.Header>
		<Card.Content>
			<form method="POST" class="w-full" use:enhance>
				<Field.Group>
					<Input
						bind:value={$form.token}
						name="token"
						type="hidden"
						class="ps-10"
						placeholder="Enter your token"
						aria-invalid={!!$errors.token}
						autocomplete="on"
					/>
					{#if errorMessage}
						<Alert.Root variant="destructive">
							<Icon icon="mingcute:warning-line" class="size-4" />
							<Alert.Title>Error</Alert.Title>
							<Alert.Description>{errorMessage}</Alert.Description>
						</Alert.Root>
					{/if}
					<Field.Field>
						<Field.Label for="otp">
							One-Time Password <span class="text-red-500 dark:text-red-400">*</span>
						</Field.Label>
						<InputOTP.Root
							maxlength={6}
							pattern={REGEXP_ONLY_DIGITS_AND_CHARS}
							bind:value={$form.otp}
							name="otp"
						>
							{#snippet children({ cells })}
								<InputOTP.Group>
									{#each cells as cell (cell)}
										<InputOTP.Slot {cell} />
									{/each}
								</InputOTP.Group>
							{/snippet}
						</InputOTP.Root>
						{#if $errors.otp}
							<Field.Error>{$errors.otp}</Field.Error>
						{/if}
					</Field.Field>
					<Field.Field>
						<Button type="submit" disabled={$submitting}>
							{#if $submitting}
								<Spinner />
							{/if}
							{$submitting ? i18n.auth_verifing_two_factor() : i18n.auth_verify_two_factor()}
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
