import { json, type RequestHandler } from '@sveltejs/kit';

const ALLOWED = new Set(['knowledge', 'playbooks', 'assets']);

/** GET /api/ai-agents/:id/:resource — list nested resource. */
export const GET: RequestHandler = async ({ params, locals }) => {
	if (!ALLOWED.has(params.resource!)) return json({ success: false }, { status: 404 });
	const data = await locals.helper.aiAgent.listResource(params.id!, params.resource!);
	return json({ success: true, data });
};

/** POST /api/ai-agents/:id/:resource — create nested resource. */
export const POST: RequestHandler = async ({ params, request, locals }) => {
	if (!ALLOWED.has(params.resource!)) return json({ success: false }, { status: 404 });
	const payload = await request.json().catch(() => null);
	if (!payload) return json({ success: false, message: 'Invalid body' }, { status: 400 });
	const res = await locals.helper.aiAgent.createResource(params.id!, params.resource!, payload);
	return json(res, { status: res.success ? 201 : 400 });
};
