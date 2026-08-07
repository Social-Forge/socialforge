import { definePageMetaTags } from 'svelte-meta-tags';

export const load = async ({ locals }) => {
	const { user, lang } = locals;
	const pageMetaTags = definePageMetaTags({
		title: 'AI Agents',
		robots: 'noindex, nofollow',
		twitter: {
			cardType: 'summary_large_image',
			site: '@socialforge',
			image: '/logo.png',
			title: 'Blog'
		}
	});

	const agents = await locals.helper.aiAgent.list().catch(() => []);

	return {
		...pageMetaTags,
		user,
		lang,
		agents
	};
};
