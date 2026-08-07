import { definePageMetaTags } from 'svelte-meta-tags';

export const load = async ({ locals }) => {
	const { user, lang } = locals;
	const pageMetaTags = definePageMetaTags({
		title: 'Chats',
		robots: 'noindex, nofollow',
		twitter: {
			cardType: 'summary_large_image',
			site: '@socialforge',
			image: '/logo.png',
			title: 'Chats'
		}
	});

	// Real conversations + realtime connection token from the backend.
	const [conversations, realtime] = await Promise.all([
		locals.helper.conversation.list().catch(() => []),
		locals.helper.conversation.realtimeToken().catch(() => null)
	]);

	return {
		...pageMetaTags,
		user,
		lang,
		conversations,
		realtime
	};
};
