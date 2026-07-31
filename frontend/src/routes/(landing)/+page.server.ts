import { definePageMetaTags } from 'svelte-meta-tags';

export const load = async ({ locals }) => {
	const { user, lang } = locals;
	const pageMetaTags = definePageMetaTags({
		title: 'Social Forge - Multi-agent Customer Service and Omnichannel CRM',
		robots: 'index, follow',
		twitter: {
			cardType: 'summary_large_image',
			site: '@socialforge',
			image: '/logo.png',
			title: 'Social Forge - Multi-agent Customer Service and Omnichannel CRM'
		}
	});

	return {
		...pageMetaTags,
		user,
		lang
	};
};
