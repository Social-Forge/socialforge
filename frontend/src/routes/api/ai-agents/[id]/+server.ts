import { json, type RequestHandler } from '@sveltejs/kit';

/** PUT /api/ai-agents/:id — update. */
export const PUT: RequestHandler = async ({ params, request, locals }) => {
	const payload = await request.json().catch(() => null);
	if (!payload) return json({ success: false, message: 'Invalid body' }, { status: 400 });
	const res = await locals.helper.aiAgent.update(params.id!, payload);
	return json(res, { status: res.success ? 200 : 400 });
};

/** DELETE /api/ai-agents/:id — delete. */
export const DELETE: RequestHandler = async ({ params, locals }) => {
	const res = await locals.helper.aiAgent.remove(params.id!);
	return json(res, { status: res.success ? 200 : 400 });
};
