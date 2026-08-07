import { json, type RequestHandler } from '@sveltejs/kit';
import type { CheckoutRequest } from '$lib/server/billing';

/** POST /api/billing/checkout — start a checkout, return the gateway URL. */
export const POST: RequestHandler = async ({ request, locals }) => {
	const body = (await request.json().catch(() => ({}))) as CheckoutRequest;
	if (!body?.kind || !body?.provider) {
		return json({ success: false, message: 'kind and provider are required' }, { status: 400 });
	}
	const result = await locals.helper.billing.checkout(body);
	if ('error' in result) {
		return json({ success: false, message: result.error }, { status: 400 });
	}
	return json({ success: true, checkout_url: result.checkoutUrl });
};
