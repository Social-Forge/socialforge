import { redirect } from '@sveltejs/kit';
import {
	allowedWhenAuthenticated,
	authRoutes,
	buildSignInRedirect,
	isSuperAdmin,
	matchesAnyRoute,
	sanitizeRedirectTarget
} from '$lib/middleware/rules';
import { localizePath } from '$lib/utils/localize-path';

export const authMiddleware = async (handler: RequestHandlerParams) => {
	const { event, resolve, isAuthenticated, method, pathname } = handler;

	if (isAuthenticated) {
		const isAuthPage =
			matchesAnyRoute(pathname, authRoutes) && !matchesAnyRoute(pathname, allowedWhenAuthenticated);

		if (isAuthPage) {
			const redirectTo = sanitizeRedirectTarget(event.url.searchParams.get('redirect'));
			throw redirect(302, localizePath(redirectTo, event.locals.lang));
		}
	}

	if (method !== 'GET' && method !== 'POST') {
		return new Response('Method not allowed', { status: 405 });
	}

	return resolve(event);
};

export const authenticatedAppMiddleware = async (handler: RequestHandlerParams) => {
	const { event, resolve, isAuthenticated, hasTenant } = handler;

	if (!isAuthenticated) {
		throw redirect(
			302,
			buildSignInRedirect(`${event.url.pathname}${event.url.search}`, event.locals.lang)
		);
	}

	const userRoleLevel = event.locals.user?.role?.level ?? null;

	if (!hasTenant && !isSuperAdmin(userRoleLevel)) {
		throw redirect(302, localizePath('/signup', event.locals.lang));
	}

	return resolve(event);
};

