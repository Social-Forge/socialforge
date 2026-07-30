import { redirect, error } from '@sveltejs/kit';

export const restrictedTenantOwnerRoutes = [
	'/app/settings',
	'/app/divisions',
	'/app/agents',
	'/app/channels',
	'/app/analytics'
];

const RoleLevelMap = {
	SUPER_ADMIN: 0,
	TENANT_OWNER: 1,
	SUPERVISOR: 2,
	AGENT: 3
} as const;

export const tenantMiddleware = async (handler: RequestHandlerParams) => {
	const { event, resolve, isAuthenticated, hasTenant, method, pathname } = handler;

	if (!isAuthenticated) {
		throw redirect(302, `/signin?redirect=${encodeURIComponent(pathname)}`);
	}

	const userRole = event.locals.user?.role;
	const allowedLevel = [RoleLevelMap.SUPER_ADMIN, RoleLevelMap.TENANT_OWNER] as number[];

	if (!userRole || !userRole.level) {
		throw error(403, {
			message: 'You do not have permission to access this resource',
			code: 'FORBIDDEN'
		});
	}

	const allowedRole = allowedLevel.includes(userRole.level);

	if (!allowedRole) {
		throw error(403, {
			message: 'You do not have permission to access this resource',
			code: 'FORBIDDEN'
		});
	}

	if (!hasTenant) {
		throw redirect(302, '/signup');
	}

	if (method === 'POST' || method === 'PUT' || method === 'PATCH' || method === 'DELETE') {
		const isRestricted = restrictedTenantOwnerRoutes.some((route) => pathname.startsWith(route));

		if (isRestricted && !allowedRole) {
			throw error(403, {
				message: 'You do not have permission to access this action',
				code: 'FORBIDDEN'
			});
		}
	}

	return resolve(event);
};
