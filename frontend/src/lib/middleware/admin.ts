import { redirect, error } from '@sveltejs/kit';
import { buildSignInRedirect, isSuperAdmin } from '$lib/middleware/rules';

export const adminMiddleware = async (handler: RequestHandlerParams) => {
        const { event, resolve, isAuthenticated } = handler;

	if (!isAuthenticated) {
                throw redirect(302, buildSignInRedirect(`${event.url.pathname}${event.url.search}`, event.locals.lang));
	}

        if (!isSuperAdmin(event.locals.user?.role?.level ?? null)) {
		throw error(403, {
			message: 'You do not have permission to access this resource',
			code: 'FORBIDDEN'
		});
	}

	return resolve(event);
};
