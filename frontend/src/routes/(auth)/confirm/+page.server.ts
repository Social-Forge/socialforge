import { fail, redirect } from '@sveltejs/kit';
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
		title: 'Confirm Email',
		robots: 'index, follow',
		canonical: defaultOrigin,
		twitter: {
			cardType: 'summary_large_image',
			site: '@socialforge',
			image: '/logo.png',
			title: 'Confirm Email'
		}
	});

	return {
		...pageMetaTags,
		token
	};
};
