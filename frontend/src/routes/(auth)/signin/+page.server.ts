import { superValidate } from 'sveltekit-superforms';
import { loginSchema } from '$lib/utils/validators';
import { fail, redirect } from '@sveltejs/kit';
import { zod4 } from 'sveltekit-superforms/adapters';
import { definePageMetaTags } from 'svelte-meta-tags';
import { localizeHref } from '$lib/paraglide/runtime.js';

export const load = async ({ url, locals, parent }) => {
	const defaultOrigin = await parent().then((data) => data.canonicalUrl || '');
	const alternates = await parent().then((data) => data.alternates || []);

	const pageMetaTags = definePageMetaTags({
		title: 'Sign In',
		robots: 'index, follow',
		canonical: defaultOrigin,
		twitter: {
			cardType: 'summary_large_image',
			site: '@socialforge',
			image: '/logo.png',
			title: 'Sign In'
		}
	});

	const form = await superValidate(zod4(loginSchema));
	return {
		...pageMetaTags,
		form
	};
};
export const actions = {
	default: async ({ request, locals }) => {
		const form = await superValidate(request, zod4(loginSchema));

		if (!form.valid) {
			return fail(400, {
				form,
				success: false,
				message: Object.values(form.errors)
					.map((error) => error)
					.join(', ')
			});
		}
		const response = await locals.helper.auth.login(form.data);

		if (!response.success) {
			return fail(400, {
				form,
				success: false,
				message: response.message || 'Invalid input'
			});
		}

		if (response.data?.two_fa_token && response.status === 202) {
			throw redirect(302, localizeHref(`/two-factor?token=${response.data?.two_fa_token || ''}`));
		}
		locals.helper.session.setAuthCookies(
			{
				accessToken: response.data?.access_token || '',
				refreshToken: response.data?.refresh_token || ''
			},
			response.data?.expires_in || 60 * 60 * 24,
			response.data?.expires_refresh_in || 60 * 60 * 24 * 7
		);
		throw redirect(302, localizeHref('/app/chats'));
	}
};
