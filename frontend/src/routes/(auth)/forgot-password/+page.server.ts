import { superValidate } from 'sveltekit-superforms';
import { forgotPasswordSchema } from '$lib/utils/validators';
import { fail } from '@sveltejs/kit';
import { zod4 } from 'sveltekit-superforms/adapters';
import { definePageMetaTags } from 'svelte-meta-tags';

export const load = async ({ url, locals, parent }) => {
	const defaultOrigin = await parent().then((data) => data.canonicalUrl || '');
	const alternates = await parent().then((data) => data.alternates || []);

	const pageMetaTags = definePageMetaTags({
		title: 'Forgot Password',
		robots: 'index, follow',
		canonical: defaultOrigin,
		twitter: {
			cardType: 'summary_large_image',
			site: '@socialforge',
			image: '/logo.png',
			title: 'Forgot Password'
		}
	});

	const form = await superValidate(zod4(forgotPasswordSchema));
	return {
		...pageMetaTags,
		form
	};
};
export const actions = {
	default: async ({ request, locals }) => {
		const form = await superValidate(request, zod4(forgotPasswordSchema));

		if (!form.valid) {
			return fail(400, {
				form,
				success: false,
				message: Object.values(form.errors)
					.map((error) => error)
					.join(', ')
			});
		}

		try {
			const response = await locals.helper.auth.forgot(form.data);
			if (!response.success) {
				return fail(400, {
					form,
					success: false,
					message: response?.message || 'Email not found'
				});
			}
			return {
				form,
				success: true,
				message: response.message || 'Password reset email sent'
			};
		} catch (error) {
			return fail(500, {
				form,
				success: false,
				message: error instanceof Error ? error.message : 'Internal server error'
			});
		}
	}
};
