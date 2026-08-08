import { json, type RequestHandler } from '@sveltejs/kit';

/** PUT /api/user/profile — update the current user's profile. */
export const PUT: RequestHandler = async ({ request, locals }) => {
	const body = await request.json().catch(() => null);
	if (!body) return json({ success: false, message: 'Invalid body' }, { status: 400 });
	const res = await locals.helper.user.updateProfile(body);
	return json({ success: res.success, message: res.message, data: res.data }, { status: res.success ? 200 : 400 });
};
