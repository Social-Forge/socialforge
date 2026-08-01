import { localizeHref } from '$lib/paraglide/runtime';
import { redirect } from '@sveltejs/kit';

export const load = async ({ locals }) => {
	throw redirect(307, localizeHref('/app/chats'));
};
