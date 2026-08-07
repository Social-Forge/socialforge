import { redirect, error } from '@sveltejs/kit';
import { buildSignInRedirect, canManageTenant, isSuperAdmin } from '$lib/middleware/rules';
import { localizePath } from '$lib/utils/localize-path';

export const tenantMiddleware = async (handler: RequestHandlerParams) => {
        const { event, resolve, isAuthenticated, hasTenant } = handler;

	if (!isAuthenticated) {
                throw redirect(302, buildSignInRedirect(`${event.url.pathname}${event.url.search}`, event.locals.lang));
	}

        const userRoleLevel = event.locals.user?.role?.level ?? null;

        if (!canManageTenant(userRoleLevel)) {
		throw error(403, {
			message: 'You do not have permission to access this resource',
			code: 'FORBIDDEN'
		});
	}

        if (!hasTenant && !isSuperAdmin(userRoleLevel)) {
                throw redirect(302, localizePath('/signup', event.locals.lang));
	}

	return resolve(event);
};
