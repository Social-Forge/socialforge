import type { Handle, RequestEvent } from '@sveltejs/kit';
import { sequence } from '@sveltejs/kit/hooks';
import { redirect } from '@sveltejs/kit';
import {
	getTextDirection,
	type Locale,
	locales as SUPPORTED_LOCALES
} from '$lib/paraglide/runtime';
import { paraglideMiddleware } from '$lib/paraglide/server';
import { ApiHandler } from '$lib/server/api';
import { ServiceHelper } from '$lib/server/index';
import { PUBLIC_API_URL } from '$env/static/public';
import { authMiddleware, authenticatedAppMiddleware } from '$lib/middleware/auth';
import { tenantMiddleware } from '$lib/middleware/tenant';
import { adminMiddleware } from '$lib/middleware/admin';
import { apiMiddleware } from '$lib/middleware/api';
import { NODE_ENV } from '$env/static/private';
import { BASE_LOCALE, getLocaleFromPath } from '$lib/utils/localize-path';
import {
        authRoutes,
        matchesAnyRoute,
        matchesRoutePrefix,
        protectedAppRoutes,
        restrictedSuperAdminRoutes,
        restrictedTenantOwnerRoutes
} from '$lib/middleware/rules';

const LOCALE_COOKIE_NAME = 'PARAGLIDE_LOCALE';

const handleParaglideWithAutoDetectedLocale: Handle = ({ event, resolve }) => {
	const { request } = event;
	const pathname = event.url.pathname;

        if (shouldBypassLocaleHandling(pathname)) {
		return resolve(event);
	}

	const ua = request.headers.get('user-agent');
	const isBot = !!ua && /bot|crawl|spider|facebookexternalhit|twitterbot/i.test(ua);
	const pathLocale = getLocaleFromPath(pathname);
        const cookieLocale = getCookieLocale(event);
        const detectedLocale = resolvePreferredLocale({
                event,
                pathLocale,
                cookieLocale
        });

	if (request.method !== 'GET' && request.method !== 'HEAD') {
                return applyParaglideLocale({ event, resolve, locale: detectedLocale });
	}

	if (isBot) {
                return applyParaglideLocale({
                        event,
                        resolve,
                        locale: pathLocale ?? BASE_LOCALE
                });
	}

        if (pathname === `/${BASE_LOCALE}` || pathname.startsWith(`/${BASE_LOCALE}/`)) {
                const stripped = pathname === `/${BASE_LOCALE}` ? '/' : pathname.slice(BASE_LOCALE.length + 1);
		throw redirect(302, `${stripped}${event.url.search}`);
	}

	if (!pathLocale) {
                const targetLocale = cookieLocale ?? detectedLocale;

                if (targetLocale !== BASE_LOCALE) {
			throw redirect(302, `/${targetLocale}${pathname === '/' ? '' : pathname}${event.url.search}`);
		}

                return applyParaglideLocale({ event, resolve, locale: targetLocale });
	}

        return applyParaglideLocale({ event, resolve, locale: pathLocale });
};

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

        const isApiRoute = matchesRoutePrefix(pathname, '/api');
        const isAppRoute = matchesAnyRoute(pathname, protectedAppRoutes);
        const isTenantRoute = matchesAnyRoute(pathname, restrictedTenantOwnerRoutes);
        const isAdminRoute = matchesAnyRoute(pathname, restrictedSuperAdminRoutes);
        const isAuthRoute = matchesAnyRoute(pathname, authRoutes);

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
                                userRoleLevel: userRoleLevel ?? null,
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
                if (isAppRoute) {
                        return await authenticatedAppMiddleware({
                                event,
                                resolve,
                                isAuthenticated,
                                hasTenant,
                                method,
                                pathname
                        });
                }
	} catch (error: any) {
                if (typeof error?.status === 'number' && error.status >= 300 && error.status < 400) {
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
export const handle: Handle = sequence(handleParaglideWithAutoDetectedLocale, initServer, auth);

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

function getCookieLocale(event: RequestEvent): Locale | null {
        const locale = event.cookies.get(LOCALE_COOKIE_NAME);
        return locale && SUPPORTED_LOCALES.includes(locale as Locale) ? (locale as Locale) : null;
}

function detectLocaleFromBrowser(event: RequestEvent): Locale {
	const accept = event.request.headers.get('accept-language');
	const l = accept?.split(',')[0].split('-')[0] as Locale;
	const supported = SUPPORTED_LOCALES.includes(l);
        return supported ? l : BASE_LOCALE;
}

function resolvePreferredLocale({
        event,
        pathLocale,
        cookieLocale
}: {
        event: RequestEvent;
        pathLocale: Locale | null;
        cookieLocale: Locale | null;
}): Locale {
        return pathLocale ?? cookieLocale ?? detectLocaleFromBrowser(event);
}

function shouldBypassLocaleHandling(pathname: string): boolean {
        return (
                pathname.startsWith('/api') ||
                pathname.startsWith('/_app') ||
                pathname === '/health' ||
                pathname.includes('.')
        );
}

function applyParaglideLocale({
        event,
        resolve,
        locale
}: {
        event: RequestEvent;
        resolve: Parameters<Handle>[0]['resolve'];
        locale: Locale;
}) {
        event.locals.lang = locale;
        setCookie(event, locale);

        return paraglideMiddleware(event.request, ({ request: localizedRequest }) => {
                event.request = localizedRequest;

                try {
                        (event as any).url = new URL(localizedRequest.url);
                } catch {
                        try {
                                Object.defineProperty(event, 'url', { value: new URL(localizedRequest.url) });
                        } catch {
                                // ignore
                        }
                }

                return resolve(event, {
                        transformPageChunk: ({ html }) =>
                                html
                                        .replace('%paraglide.lang%', locale)
                                        .replace('%paraglide.dir%', getTextDirection(locale))
                });
        });
}

function setCookie(event: RequestEvent, locale: Locale) {
	event.cookies.set('PARAGLIDE_LOCALE', locale, {
                // Locale preference is also updated from the client via Paraglide's setLocale().
                // Keep this cookie readable/writable in the browser so client-side switches
                // don't get overwritten by a stale server-only value on the next request.
                httpOnly: false,
		secure: NODE_ENV === 'production',
		sameSite: 'lax',
		path: '/',
		maxAge: 60 * 60 * 24 * 7
	});
}
