import { json, type RequestHandler } from '@sveltejs/kit';

const ALLOWED = new Set(['knowledge', 'playbooks', 'assets']);

/** PUT /api/ai-agents/:id/:resource/:childId — update nested resource. */
export const PUT: RequestHandler = async ({ params, request, locals }) => {
	if (!ALLOWED.has(params.resource!)) return json({ success: false }, { status: 404 });
	const payload = await request.json().catch(() => null);
	if (!payload) return json({ success: false, message: 'Invalid body' }, { status: 400 });
	const res = await locals.helper.aiAgent.updateResource(
		params.id!,
		params.resource!,
		params.childId!,
		payload
	);
	return json(res, { status: res.success ? 200 : 400 });
};

/** DELETE /api/ai-agents/:id/:resource/:childId — delete nested resource. */
export const DELETE: RequestHandler = async ({ params, locals }) => {
	if (!ALLOWED.has(params.resource!)) return json({ success: false }, { status: 404 });
	const res = await locals.helper.aiAgent.deleteResource(params.id!, params.resource!, params.childId!);
	return json(res, { status: res.success ? 200 : 400 });
};
