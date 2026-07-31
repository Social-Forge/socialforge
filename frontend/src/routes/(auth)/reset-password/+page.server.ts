import { superValidate } from 'sveltekit-superforms';
import { resetPasswordSchema } from '$lib/utils/validators';
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
		title: 'Reset Password',
		robots: 'index, follow',
		canonical: defaultOrigin,
		twitter: {
			cardType: 'summary_large_image',
			site: '@socialforge',
			image: '/logo.png',
			title: 'Reset Password'
		}
	});

	const form = await superValidate(
		{
			token,
			new_password: '',
			confirm_password: ''
		},
		zod4(resetPasswordSchema)
	);
	return {
		...pageMetaTags,
		form,
		token
	};
};
export const actions = {
	default: async ({ request, locals }) => {
		const form = await superValidate(request, zod4(resetPasswordSchema));
		if (!form.valid) {
			return fail(400, {
				form,
				succcess: false,
				message: Object.values(form.errors)
					.map((error) => error)
					.join(', ')
			});
		}

		const response = await locals.helper.auth.resetPassword(form.data);
		if (!response.success) {
			return fail(400, {
				form,
				succcess: false,
				message: response.message
			});
		}
		throw redirect(302, localizeHref('/signin'));
	}
};
