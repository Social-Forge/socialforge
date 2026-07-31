import { fail, redirect } from '@sveltejs/kit';
import { definePageMetaTags } from 'svelte-meta-tags';
import { localizeHref } from '$lib/paraglide/runtime.js';

export const load = async ({ url, locals, parent }) => {
	const defaultOrigin = await parent().then((data) => data.canonicalUrl || '');
	const alternates = await parent().then((data) => data.alternates || []);

	const email = url.searchParams.get('email');
	if (!email) {
		throw redirect(302, localizeHref('/signin'));
	}

	const pageMetaTags = definePageMetaTags({
		title: 'Verify Email',
		robots: 'index, follow',
		canonical: defaultOrigin,
		twitter: {
			cardType: 'summary_large_image',
			site: '@socialforge',
			image: '/logo.png',
			title: 'Verify Email'
		}
	});

	return {
		...pageMetaTags,
		email
	};
};
