<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import * as Card from '$lib/components/ui/card/index.js';
	import { Progress } from '$lib/components/ui/progress/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import * as Alert from '$lib/components/ui/alert/index.js';
	import { Spinner } from '$lib/components/ui/spinner/index.js';
	import Icon from '@iconify/svelte';
	import { localizeHref } from '$lib/paraglide/runtime.js';
	import * as i18n from '$lib/paraglide/messages.js';

	let { data } = $props();

	let errorMessage = $state<string | undefined>(undefined);
	let successMessage = $state<string | undefined>(undefined);
	let progressValue = $state<number>(0);

	const validateToken = async () => {
		progressValue = 20;
		try {
			const response = await fetch('/api/auth/email-validate', {
				method: 'POST',
				body: JSON.stringify({ token: data.token })
			});
			const json = await response.json();
			if (!response.ok) {
				errorMessage = json.message;
			} else {
				successMessage = json.message;
				// await goto('/auth/sign-in');
			}
		} catch (error) {
			errorMessage = error instanceof Error ? error.message : 'Internal server error';
		} finally {
			progressValue = 100;
		}
	};

	onMount(async () => {
		await validateToken();
	});
</script>

<div class="flex flex-col gap-6">
	<Card.Root>
		<Card.Header class="text-center">
			<Card.Title class="text-xl">{i18n.auth_confirm_email()}</Card.Title>
			<Card.Description>{i18n.auth_confirm_email_description()}</Card.Description>
		</Card.Header>
		<Card.Content class="space-y-4">
			{#if errorMessage}
				<Alert.Root variant="destructive">
					<Icon icon="mingcute:warning-line" class="size-4" />
					<Alert.Title>Error</Alert.Title>
					<Alert.Description>{errorMessage}</Alert.Description>
				</Alert.Root>
			{:else if successMessage}
				<Alert.Root variant="default">
					<Icon icon="mingcute:check-line" class="size-4" />
					<Alert.Title>Success</Alert.Title>
					<Alert.Description>{successMessage}</Alert.Description>
				</Alert.Root>
				<Button type="button" variant="default" href={localizeHref('/signin')}
					>{i18n.auth_button_back_to_login()}</Button
				>
			{:else}
				<Progress value={progressValue} class="w-full" />
				<div class="flex items-center justify-center gap-2">
					<Spinner class="size-4" />
					<p class="text-sm opacity-70">{i18n.auth_confirm_email_verification()}</p>
				</div>
			{/if}
		</Card.Content>
	</Card.Root>
</div>
