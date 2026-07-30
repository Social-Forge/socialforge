import { redirect, error } from '@sveltejs/kit';

export const restrictedSuperAdminRoutes = ['/app/admin'];

export const adminMiddleware = async (handler: RequestHandlerParams) => {
	const { event, resolve, isAuthenticated, hasTenant, method, pathname } = handler;

	if (!isAuthenticated) {
		throw redirect(302, `/signin?redirect=${encodeURIComponent(pathname)}`);
	}

	const userRole = event.locals.user?.role;
	const allowedRoles = [0];

	if (!userRole || !allowedRoles.includes(userRole.level)) {
		throw error(403, {
			message: 'You do not have permission to access this resource',
			code: 'FORBIDDEN'
		});
	}

	if (method === 'POST' || method === 'PUT' || method === 'PATCH' || method === 'DELETE') {
		const isRestricted = restrictedSuperAdminRoutes.some((route) => pathname.startsWith(route));

		if (isRestricted && !allowedRoles.includes(userRole.level)) {
			throw error(403, {
				message: 'You do not have permission to access this action',
				code: 'FORBIDDEN'
			});
		}
	}

	return resolve(event);
};
