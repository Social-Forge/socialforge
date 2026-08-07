import { definePageMetaTags } from 'svelte-meta-tags';

export const load = async ({ locals }) => {
	const { user, lang } = locals;
	const pageMetaTags = definePageMetaTags({
		title: 'Integrations',
		robots: 'noindex, nofollow',
		twitter: {
			cardType: 'summary_large_image',
			site: '@socialforge',
			image: '/logo.png',
			title: 'Blog'
		}
	});

	const channelCounts = await locals.helper.channel.countsByType().catch(() => ({}));

	return {
		...pageMetaTags,
		user,
		lang,
		channelCounts
	};
};
