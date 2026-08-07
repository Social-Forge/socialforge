import { superValidate } from 'sveltekit-superforms';
import { registerSchema } from '$lib/utils/validators';
import { fail, redirect } from '@sveltejs/kit';
import { zod4 } from 'sveltekit-superforms/adapters';
import { definePageMetaTags } from 'svelte-meta-tags';
import { localizeHref } from '$lib/paraglide/runtime.js';

export const load = async ({ url, locals, parent }) => {
	const defaultOrigin = await parent().then((data) => data.canonicalUrl || '');
	const alternates = await parent().then((data) => data.alternates || []);
        const oauthError = url.searchParams.get('oauth_error');
        const redirectTarget = url.searchParams.get('redirect') || '/app/chats';

	const pageMetaTags = definePageMetaTags({
		title: 'Sign Up',
		robots: 'index, follow',
		canonical: defaultOrigin,
		twitter: {
			cardType: 'summary_large_image',
			site: '@socialforge',
			image: '/logo.png',
			title: 'Sign Up'
		}
	});

	const form = await superValidate(zod4(registerSchema));
	return {
		...pageMetaTags,
                form,
                oauthError,
                redirectTarget
	};
};
export const actions = {
	default: async ({ request, locals }) => {
		const form = await superValidate(request, zod4(registerSchema));

		if (!form.valid) {
			return fail(400, {
				form,
				success: false,
				message: Object.values(form.errors)
					.map((error) => error)
					.join(', ')
			});
		}
		const response = await locals.helper.auth.register(form.data);

		if (!response.success) {
			return fail(400, {
				form,
				success: false,
				message: response.message || 'Invalid input'
			});
		}

		throw redirect(302, localizeHref(`/verify-email?email=${form.data.email}`));
	}
};
