import { definePageMetaTags } from 'svelte-meta-tags';

export const load = async ({ locals, url }) => {
	const { user, lang } = locals;
	const pageMetaTags = definePageMetaTags({
		title: 'Contacts',
		robots: 'noindex, nofollow',
		twitter: {
			cardType: 'summary_large_image',
			site: '@socialforge',
			image: '/logo.png',
			title: 'Contacts'
		}
	});

	const page = Number(url.searchParams.get('page') ?? '1') || 1;
	const perPage = Number(url.searchParams.get('limit') ?? '20') || 20;
	const search = url.searchParams.get('search') ?? '';
	const { contacts, meta } = await locals.helper.contact
		.list({ page, perPage, search })
		.catch(() => ({ contacts: [], meta: { page: 1, per_page: 20, total: 0, total_pages: 0, has_more: false } }));

	return {
		...pageMetaTags,
		user,
		lang,
		contacts,
		meta,
		search
	};
};
