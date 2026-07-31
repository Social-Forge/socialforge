import { superValidate } from 'sveltekit-superforms';
import { verifyTwoFactorSchema } from '$lib/utils/validators';
import { fail, redirect } from '@sveltejs/kit';
import { zod4 } from 'sveltekit-superforms/adapters';
import { definePageMetaTags } from 'svelte-meta-tags';
import { localizeHref } from '$lib/paraglide/runtime.js';

export const load = async ({ url, locals, parent }) => {
	const defaultOrigin = await parent().then((data) => data.canonicalUrl || '');
	const alternates = await parent().then((data) => data.alternates || []);

	const token = url.searchParams.get('token');
	if (!token) {
		throw redirect(302, localizeHref('/signin'));
	}

	const pageMetaTags = definePageMetaTags({
		title: 'Verify Two-Factor',
		robots: 'index, follow',
		canonical: defaultOrigin,
		twitter: {
			cardType: 'summary_large_image',
			site: '@socialforge',
			image: '/logo.png',
			title: 'Verify Two-Factor'
		}
	});

	const form = await superValidate(zod4(verifyTwoFactorSchema));
	return {
		...pageMetaTags,
		form,
		token
	};
};
export const actions = {
	default: async ({ request, locals }) => {
		const form = await superValidate(request, zod4(verifyTwoFactorSchema));
		if (!form.valid) {
			return fail(400, {
				form,
				success: false,
				message: Object.values(form.errors)
					.map((error) => error)
					.join(', ')
			});
		}

		const response = await locals.helper.auth.verifyTwoFactor(form.data);
		if (!response.success) {
			return fail(response.status, {
				form,
				success: false,
				message: response.message
			});
		}
		if (!response.data?.access_token || !response.data?.refresh_token) {
			return fail(500, {
				form,
				success: false,
				message: 'Invalid response'
			});
		}
		locals.helper.session.setAuthCookies(
			{
				accessToken: response.data.access_token,
				refreshToken: response.data.refresh_token
			},
			response.data?.expires_in || 60 * 60 * 24,
			response.data?.expires_refresh_in || 60 * 60 * 24 * 7
		);

		throw redirect(302, localizeHref('/app/chats'));
	}
};
