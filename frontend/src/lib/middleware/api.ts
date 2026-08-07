import { error } from '@sveltejs/kit';
import type { RequestEvent } from '@sveltejs/kit';
import {
        adminApiRoutes,
        canManageTenant,
        isSuperAdmin,
        matchesAnyRoute,
        matchesRoutePrefix,
        publicApiRoutes,
        tenantMemberApiRoutes,
        tenantOwnerApiRoutes
} from '$lib/middleware/rules';

export const apiMiddleware = async ({
	method,
	pathname,
	isAuthenticated,
	userRoleLevel,
	hasTenant
}: ApiMiddlewareParams) => {
        if (!matchesRoutePrefix(pathname, '/api')) {
		return { allowed: true };
	}

        if (matchesAnyRoute(pathname, publicApiRoutes)) {
                return { allowed: true };
        }

        if (!isAuthenticated) {
		throw error(401, {
			message: 'Authentication required',
			code: 'UNAUTHORIZED'
		});
	}

        if (matchesAnyRoute(pathname, adminApiRoutes) && !isSuperAdmin(userRoleLevel)) {
                throw error(403, {
                        message: 'Admin access required',
                        code: 'FORBIDDEN'
                });
	}

        const requiresTenant =
                matchesAnyRoute(pathname, tenantMemberApiRoutes) ||
                matchesAnyRoute(pathname, tenantOwnerApiRoutes);

	if (requiresTenant && !hasTenant) {
		throw error(403, {
			message: 'Tenant required',
			code: 'TENANT_REQUIRED'
		});
	}

        const restrictedMethods = ['POST', 'PUT', 'PATCH', 'DELETE'];
        const isRestrictedMethod = restrictedMethods.includes(method);
        const requiresTenantOwner = matchesAnyRoute(pathname, tenantOwnerApiRoutes);

        if (requiresTenantOwner && isRestrictedMethod && !canManageTenant(userRoleLevel)) {
                throw error(403, {
                        message: 'Insufficient permissions for this action',
                        code: 'FORBIDDEN'
                });
	}

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
