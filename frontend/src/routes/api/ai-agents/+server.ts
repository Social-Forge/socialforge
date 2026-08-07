import { json, type RequestHandler } from '@sveltejs/kit';

/** POST /api/ai-agents — create an agent. */
export const POST: RequestHandler = async ({ request, locals }) => {
	const payload = await request.json().catch(() => null);
	if (!payload) return json({ success: false, message: 'Invalid body' }, { status: 400 });
	const res = await locals.helper.aiAgent.create(payload);
	return json(res, { status: res.success ? 201 : 400 });
};
