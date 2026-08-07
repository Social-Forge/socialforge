import { json, type RequestHandler } from '@sveltejs/kit';

/** POST /api/contacts/:id — block/unblock ({ blocked: boolean }). */
export const POST: RequestHandler = async ({ params, request, locals }) => {
	const body = await request.json().catch(() => ({}));
	const blocked = !!body?.blocked;
	const res = await locals.helper.contact.setBlocked(params.id!, blocked);
	return json({ success: res.success, message: res.message }, { status: res.success ? 200 : 400 });
};

/** DELETE /api/contacts/:id — delete a contact. */
export const DELETE: RequestHandler = async ({ params, locals }) => {
	const res = await locals.helper.contact.remove(params.id!);
	return json({ success: res.success, message: res.message }, { status: res.success ? 200 : 400 });
};
