import { json, type RequestHandler } from '@sveltejs/kit';

/** GET /api/chats/:id/messages — conversation messages (oldest-first). */
export const GET: RequestHandler = async ({ params, locals }) => {
	const messages = await locals.helper.conversation.messages(params.id!);
	return json({ success: true, data: messages });
};

/** POST /api/chats/:id/messages — send an agent text reply. */
export const POST: RequestHandler = async ({ params, request, locals }) => {
	const body = await request.json().catch(() => ({}));
	const text = typeof body?.text === 'string' ? body.text.trim() : '';
	if (!text) {
		return json({ success: false, message: 'text is required' }, { status: 400 });
	}
	const message = await locals.helper.conversation.send(params.id!, text);
	return json({ success: !!message, data: message }, { status: message ? 201 : 400 });
};
