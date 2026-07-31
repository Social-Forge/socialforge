import { definePageMetaTags } from 'svelte-meta-tags';

export const load = async ({ locals }) => {
	const { user, lang } = locals;
	const pageMetaTags = definePageMetaTags({
		title: 'Blog',
		robots: 'index, follow',
		twitter: {
			cardType: 'summary_large_image',
			site: '@socialforge',
			image: '/logo.png',
			title: 'Blog'
		}
	});

	return {
		...pageMetaTags,
		user,
		lang
	};
};
