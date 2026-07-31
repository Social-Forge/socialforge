import { definePageMetaTags } from 'svelte-meta-tags';

export const load = async ({ locals }) => {
	const { user, lang } = locals;
	const pageMetaTags = definePageMetaTags({
		title: 'Career',
		robots: 'index, follow',
		twitter: {
			cardType: 'summary_large_image',
			site: '@socialforge',
			image: '/logo.png',
			title: 'Career'
		}
	});

	return {
		...pageMetaTags,
		user,
		lang
	};
};
