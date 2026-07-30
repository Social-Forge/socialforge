import { error, redirect } from '@sveltejs/kit';
import type { RequestEvent } from '@sveltejs/kit';

export const protectedApiRoutes = [
	'/api/users',
	'/api/chats',
	'/api/messages',
	'/api/tenants',
	'/api/divisions',
	'/api/agents',
	'/api/channels',
	'/api/settings',
	'/api/analytics',
	'/api/profile'
];

export const adminApiRoutes = [
	'/api/admin/users',
	'/api/admin/tenants',
	'/api/admin/system',
	'/api/admin/logs',
	'/api/admin/backup'
];

export const tenantApiRoutes = [
	'/api/tenants',
	'/api/divisions',
	'/api/agents',
	'/api/channels',
	'/api/settings'
];

const RoleLevelMap = {
	SUPER_ADMIN: 0,
	TENANT_OWNER: 1,
	SUPERVISOR: 2,
	AGENT: 3
} as const;

export const apiMiddleware = async ({
	event,
	method,
	pathname,
	isAuthenticated,
	userRoleLevel,
	hasTenant
}: ApiMiddlewareParams) => {
	if (!pathname.startsWith('/api')) {
		return { allowed: true };
	}

	const requiresAuth = protectedApiRoutes.some((route) => pathname.startsWith(route));

	if (requiresAuth && !isAuthenticated) {
		throw error(401, {
			message: 'Authentication required',
			code: 'UNAUTHORIZED'
		});
	}

	const isAdminRoute = adminApiRoutes.some((route) => pathname.startsWith(route));
	if (isAdminRoute) {
		if (userRoleLevel !== RoleLevelMap.SUPER_ADMIN) {
			throw error(403, {
				message: 'Admin access required',
				code: 'FORBIDDEN'
			});
		}
	}

	const requiresTenant = tenantApiRoutes.some((route) => pathname.startsWith(route));
	if (requiresTenant && !hasTenant) {
		throw error(403, {
			message: 'Tenant required',
			code: 'TENANT_REQUIRED'
		});
	}

	const restrictedMethods = ['POST', 'PUT', 'PATCH', 'DELETE'];
	const isRestrictedMethod = restrictedMethods.includes(method);

	if (requiresTenant && isRestrictedMethod) {
		const allowedRoles = [RoleLevelMap.SUPER_ADMIN, RoleLevelMap.TENANT_OWNER] as const as number[];

		if (userRoleLevel == null || !allowedRoles.includes(userRoleLevel)) {
			throw error(403, {
				message: 'Insufficient permissions for this action',
				code: 'FORBIDDEN'
			});
		}
	}

	// 6. Rate limiting untuk API (opsional)

	return { allowed: true };
};

export const validateApiRequest = async (event: RequestEvent, schema?: any) => {
	if (event.request.method === 'GET' || event.request.method === 'DELETE') {
		return { valid: true };
	}

	try {
		const body = await event.request.json();

		if (schema) {
			// const result = schema.safeParse(body);
			// if (!result.success) {
			// 	throw error(400, {
			// 		message: 'Invalid request body',
			// 		code: 'VALIDATION_ERROR',
			// 		details: result.error.errors
			// 	});
			// }
			// return { valid: true, data: result.data };
		}

		return { valid: true, data: body };
	} catch (err) {
		throw error(400, {
			message: 'Invalid JSON body',
			code: 'INVALID_JSON'
		});
	}
};
