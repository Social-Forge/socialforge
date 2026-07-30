import type { Handle, RequestEvent } from '@sveltejs/kit';
import { sequence } from '@sveltejs/kit/hooks';
import { redirect } from '@sveltejs/kit';
import { getTextDirection } from '$lib/paraglide/runtime';
import { paraglideMiddleware } from '$lib/paraglide/server';
import { ApiHandler } from '$lib/server/api';
import { ServiceHelper } from '$lib/server/index';
import { PUBLIC_API_URL } from '$env/static/public';
import { authMiddleware, authRoutes } from '$lib/middleware/auth';
import { tenantMiddleware, restrictedTenantOwnerRoutes } from '$lib/middleware/tenant';
import { adminMiddleware, restrictedSuperAdminRoutes } from '$lib/middleware/admin';
import {
	apiMiddleware,
	protectedApiRoutes,
	adminApiRoutes,
	tenantApiRoutes
} from '$lib/middleware/api';

const handleParaglide: Handle = ({ event, resolve }) =>
	paraglideMiddleware(event.request, ({ request, locale }) => {
		event.request = request;

		return resolve(event, {
			transformPageChunk: ({ html }) =>
				html
					.replace('%paraglide.lang%', locale)
					.replace('%paraglide.dir%', getTextDirection(locale))
		});
	});

const initServer: Handle = async ({ event, resolve }) => {
	event.locals.origin = PUBLIC_API_URL;
	event.locals.api = new ApiHandler(event);
	event.locals.helper = new ServiceHelper(event);
	event.locals.safeGetUser = async (): Promise<UserResponse | null> => {
		try {
			const shouldRefresh = await handleAutoRefresh(event);
			if (shouldRefresh) {
				console.log('🔄 Token refreshed automatically');
			}

			const user = await event.locals.helper.user.currentUser();
			if (!user) {
				return null;
			}

			return user;
		} catch (error) {
			console.error('Error fetching user:', error);
			return null;
		}
	};
	const response = await resolve(event);

	if (response.status === 404) {
		throw redirect(303, '/');
	}
	if (response.status === 403) {
		throw redirect(307, '/');
	}
	return response;
};
const auth: Handle = async ({ event, resolve }) => {
	const { url, request } = event;
	const pathname = url.pathname;
	const method = request.method;

	const isApiRoute = pathname.startsWith('/api');
	const isTenantRoute = restrictedTenantOwnerRoutes.some((route) => pathname.startsWith(route));
	const isAdminRoute = restrictedSuperAdminRoutes.some((route) => pathname.startsWith(route));
	const isAuthRoute = authRoutes.some((route) => pathname.startsWith(route));

	try {
		if (!isApiRoute) {
			await handleAutoRefresh(event);
		}

		const user = await event.locals.safeGetUser();
		event.locals.user = user;

		const isAuthenticated = user !== null;
		const hasTenant = user?.tenant !== null && user?.user_tenant !== null;
		const userRoleLevel = user?.role?.level;

		if (isApiRoute) {
			await apiMiddleware({
				event,
				method,
				pathname,
				isAuthenticated,
				userRoleLevel: userRoleLevel || 0,
				hasTenant
			});

			return resolve(event);
		}
		if (isAuthRoute) {
			return await authMiddleware({
				event,
				resolve,
				isAuthenticated,
				hasTenant,
				method,
				pathname
			});
		}
		if (isAdminRoute) {
			return await adminMiddleware({
				event,
				resolve,
				isAuthenticated,
				hasTenant,
				method,
				pathname
			});
		}
		if (isTenantRoute) {
			return await tenantMiddleware({
				event,
				resolve,
				isAuthenticated,
				hasTenant,
				method,
				pathname
			});
		}
	} catch (error: any) {
		if (error?.status === 302 || error?.status === 301) {
			throw error;
		}
		if (isApiRoute) {
			return new Response(
				JSON.stringify({
					success: false,
					error: {
						code: error?.code || 'INTERNAL_ERROR',
						message: error?.message || 'An unexpected error occurred',
						status: error?.status || 500
					}
				}),
				{
					status: error?.status || 500,
					headers: {
						'Content-Type': 'application/json'
					}
				}
			);
		}
		console.error('Auth middleware error:', error);
	}

	return resolve(event);
};
export const handle: Handle = sequence(
	handleParaglide,
	initServer
	// auth
);

async function handleAutoRefresh(event: RequestEvent): Promise<boolean> {
	try {
		const accessToken = event.locals.helper.session.getAccessToken();
		const refreshToken = event.locals.helper.session.getRefreshToken();

		// Jika ada access token yang valid, tidak perlu refresh
		if (accessToken && !event.locals.helper.session.isTokenExpired(accessToken)) {
			return false;
		}

		// Jika ada refresh token yang masih valid, lakukan refresh
		if (refreshToken && !event.locals.helper.session.isTokenExpired(refreshToken)) {
			console.log('🔄 Attempting token refresh...');
			const refreshResult = await handleRefreshSession(event);

			if (refreshResult) {
				console.log('✅ Token refresh successful');
				return true;
			} else {
				console.log('❌ Token refresh failed');
				event.locals.helper.session.setAuthCookies(null, 0, 0);
				return false;
			}
		}

		return false;
	} catch (error) {
		console.error('Auto refresh error:', error);
		return false;
	}
}
async function handleRefreshSession(event: RequestEvent) {
	try {
		const refreshToken = event.locals.helper.session.getRefreshToken();

		if (!refreshToken) {
			return null;
		}

		const response = await event.locals.helper.auth.refreshToken(refreshToken);

		if (!response.success) {
			console.warn('Refresh token failed:', response.message);
			return null;
		}

		if (!response.data?.access_token || !response.data?.refresh_token) {
			console.warn('Refresh token response missing tokens');
			return null;
		}
		event.locals.helper.session.setAuthCookies(
			{
				accessToken: response.data?.access_token || '',
				refreshToken: response.data?.refresh_token || ''
			},
			response.data?.expires_in || 60 * 60 * 24, // 24 jam default
			response.data?.expires_refresh_in || 60 * 60 * 24 * 7 // 7 hari default
		);

		return response;
	} catch (error) {
		console.error('Refresh session error:', error);
		return null;
	}
}
