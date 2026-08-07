<script lang="ts" module>
	import type { SocialAuthProvider } from '$lib/auth/social-providers';

	export type SocialAuthButtonsProps = {
		mode?: 'signin' | 'signup';
		redirectTo?: string;
	};
</script>

<script lang="ts">
	import { localizeHref } from '$lib/paraglide/runtime.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import * as Field from '$lib/components/ui/field/index.js';
	import * as i18n from '$lib/paraglide/messages.js';

	let { mode = 'signin', redirectTo = '/app/chats' }: SocialAuthButtonsProps = $props();

	const providers = [
		{
			code: 'google',
			label: i18n.auth_button_continue_with_google(),
			icon: 'google'
		},
		{
			code: 'facebook',
			label: i18n.auth_button_continue_with_facebook(),
			icon: 'facebook'
		},
		{
			code: 'github',
			label: i18n.auth_button_continue_with_github(),
			icon: 'github'
		}
	] as const satisfies ReadonlyArray<{
		code: SocialAuthProvider;
		label: string;
		icon: 'google' | 'facebook' | 'github';
	}>;

	function buildProviderHref(provider: SocialAuthProvider): string {
		const params = new URLSearchParams({
			mode,
			redirect: redirectTo
		});

		return localizeHref(`/auth/${provider}?${params.toString()}`);
	}
</script>

<Field.Field>
	{#each providers as provider (provider.code)}
		<Button variant="outline" type="button" href={buildProviderHref(provider.code)}>
			{#if provider.icon === 'google'}
				<svg xmlns="http://www.w3.org/2000/svg" width="32" height="32" viewBox="0 0 16 16">
					<g fill="none" fill-rule="evenodd" clip-rule="evenodd">
						<path
							fill="#f44336"
							d="M7.209 1.061c.725-.081 1.154-.081 1.933 0a6.57 6.57 0 0 1 3.65 1.82a100 100 0 0 0-1.986 1.93q-1.876-1.59-4.188-.734q-1.696.78-2.362 2.528a78 78 0 0 1-2.148-1.658a.26.26 0 0 0-.16-.027q1.683-3.245 5.26-3.86"
							opacity=".987"
						/>
						<path
							fill="#ffc107"
							d="M1.946 4.92q.085-.013.161.027a78 78 0 0 0 2.148 1.658A7.6 7.6 0 0 0 4.04 7.99q.037.678.215 1.331L2 11.116Q.527 8.038 1.946 4.92"
							opacity=".997"
						/>
						<path
							fill="#448aff"
							d="M12.685 13.29a26 26 0 0 0-2.202-1.74q1.15-.812 1.396-2.228H8.122V6.713q3.25-.027 6.497.055q.616 3.345-1.423 6.032a7 7 0 0 1-.51.49"
							opacity=".999"
						/>
						<path
							fill="#43a047"
							d="M4.255 9.322q1.23 3.057 4.51 2.854a3.94 3.94 0 0 0 1.718-.626q1.148.812 2.202 1.74a6.62 6.62 0 0 1-4.027 1.684a6.4 6.4 0 0 1-1.02 0Q3.82 14.524 2 11.116z"
							opacity=".993"
						/>
					</g>
				</svg>
			{:else if provider.icon === 'facebook'}
				<svg xmlns="http://www.w3.org/2000/svg" width="32" height="32" viewBox="0 0 256 256">
					<path
						fill="#1877f2"
						d="M256 128C256 57.308 198.692 0 128 0S0 57.308 0 128c0 63.888 46.808 116.843 108 126.445V165H75.5v-37H108V99.8c0-32.08 19.11-49.8 48.348-49.8C170.352 50 185 52.5 185 52.5V84h-16.14C152.959 84 148 93.867 148 103.99V128h35.5l-5.675 37H148v89.445c61.192-9.602 108-62.556 108-126.445"
					/>
					<path
						fill="#fff"
						d="m177.825 165l5.675-37H148v-24.01C148 93.866 152.959 84 168.86 84H185V52.5S170.352 50 156.347 50C127.11 50 108 67.72 108 99.8V128H75.5v37H108v89.445A129 129 0 0 0 128 256a129 129 0 0 0 20-1.555V165z"
					/>
				</svg>
			{:else}
				<svg xmlns="http://www.w3.org/2000/svg" width="32" height="32" viewBox="0 0 24 24">
					<path
						fill="currentColor"
						d="M12 2A10 10 0 0 0 2 12c0 4.42 2.87 8.17 6.84 9.5c.5.08.66-.23.66-.5v-1.69c-2.77.6-3.36-1.34-3.36-1.34c-.46-1.16-1.11-1.47-1.11-1.47c-.91-.62.07-.6.07-.6c1 .07 1.53 1.03 1.53 1.03c.87 1.52 2.34 1.07 2.91.83c.09-.65.35-1.09.63-1.34c-2.22-.25-4.55-1.11-4.55-4.92c0-1.11.38-2 1.03-2.71c-.1-.25-.45-1.29.1-2.64c0 0 .84-.27 2.75 1.02c.79-.22 1.65-.33 2.5-.33s1.71.11 2.5.33c1.91-1.29 2.75-1.02 2.75-1.02c.55 1.35.2 2.39.1 2.64c.65.71 1.03 1.6 1.03 2.71c0 3.82-2.34 4.66-4.57 4.91c.36.31.69.92.69 1.85V21c0 .27.16.59.67.5C19.14 20.16 22 16.42 22 12A10 10 0 0 0 12 2"
					/>
				</svg>
			{/if}
			{provider.label}
		</Button>
	{/each}
</Field.Field>
