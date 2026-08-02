import { definePageMetaTags } from 'svelte-meta-tags';
import { redirect } from '@sveltejs/kit';
import { localizeHref } from '$lib/paraglide/runtime.js';

export const load = async ({ locals, params }) => {
	const { user, lang } = locals;

	const slug = params.slug;
	if (!slug) {
		throw redirect(302, localizeHref('/app/integrations'));
	}

	const pageMetaTags = definePageMetaTags({
		title: `Integrations - ${slug.charAt(0).toUpperCase() + slug.slice(1)}`,
		robots: 'noindex, nofollow',
		twitter: {
			cardType: 'summary_large_image',
			site: '@socialforge',
			image: '/logo.png',
			title: `Integrations - ${slug.charAt(0).toUpperCase() + slug.slice(1)}`
		}
	});

	return {
		...pageMetaTags,
		user,
		lang,
		slug
	};
};
