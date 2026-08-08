import { definePageMetaTags } from 'svelte-meta-tags';

export const load = async ({ locals }) => {
	const { user, lang } = locals;
	const pageMetaTags = definePageMetaTags({
		title: 'Billing Settings',
		robots: 'noindex, nofollow',
		twitter: {
			cardType: 'summary_large_image',
			site: '@socialforge',
			image: '/logo.png',
			title: 'Billing Settings'
		}
	});

	const [subscription, invoices, plans] = await Promise.all([
		locals.helper.billing.subscription().catch(() => null),
		locals.helper.billing.invoices().catch(() => []),
		locals.helper.billing.plans().catch(() => [])
	]);

	return {
		...pageMetaTags,
		user,
		lang,
		subscription,
		invoices,
		plans
	};
};
