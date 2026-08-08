import { json, type RequestHandler } from '@sveltejs/kit';

/** PUT /api/user/password — change the current user's password. */
export const PUT: RequestHandler = async ({ request, locals }) => {
	const body = await request.json().catch(() => null);
	if (!body) return json({ success: false, message: 'Invalid body' }, { status: 400 });
	const res = await locals.helper.user.changePassword(body);
	return json({ success: res.success, message: res.message }, { status: res.success ? 200 : 400 });
};
