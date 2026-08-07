import { definePageMetaTags } from 'svelte-meta-tags';

export const load = async ({ locals }) => {
	const { user, lang } = locals;
	const pageMetaTags = definePageMetaTags({
		title: 'Agent Settings',
		robots: 'noindex, nofollow',
		twitter: {
			cardType: 'summary_large_image',
			site: '@socialforge',
			image: '/logo.png',
			title: 'Agent Settings'
		}
	});

	return {
		...pageMetaTags,
		user,
		lang
	};
};
