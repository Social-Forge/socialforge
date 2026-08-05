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
import { authMiddleware, authRoutes } from '$lib/middleware/auth';
import { tenantMiddleware, restrictedTenantOwnerRoutes } from '$lib/middleware/tenant';
import { adminMiddleware, restrictedSuperAdminRoutes } from '$lib/middleware/admin';
import {
	apiMiddleware,
	protectedApiRoutes,
	adminApiRoutes,
	tenantApiRoutes
} from '$lib/middleware/api';
import { NODE_ENV } from '$env/static/private';

const LOCALE_COOKIE_NAME = 'PARAGLIDE_LOCALE';

const handleParaglide: Handle = ({ event, resolve }) =>
	paraglideMiddleware(event.request, ({ request, locale }) => {
		event.request = request;
		event.locals.lang = locale;

		return resolve(event, {
			transformPageChunk: ({ html }) =>
				html
					.replace('%paraglide.lang%', locale)
					.replace('%paraglide.dir%', getTextDirection(locale))
		});
	});
const paraglideHandleWithAutoDetectedLocale: Handle = ({ event, resolve }) => {
	const { request } = event;
	const pathname = event.url.pathname;

	if (
		pathname.startsWith('/api') ||
		pathname.startsWith('/_app') ||
		pathname === '/health' ||
		pathname.includes('.')
	) {
		return resolve(event);
	}

	const ua = request.headers.get('user-agent');
	const isBot = !!ua && /bot|crawl|spider|facebookexternalhit|twitterbot/i.test(ua);
	const pathLocale = getLocaleFromPath(pathname);
	const cookieLocale = event.cookies.get(LOCALE_COOKIE_NAME) as Locale | null;
	let detectedLocale: Locale;
	if (pathLocale) {
		detectedLocale = pathLocale;
	} else if (cookieLocale && SUPPORTED_LOCALES.includes(cookieLocale)) {
		detectedLocale = cookieLocale;
	} else {
		detectedLocale = detectLocaleFromBrowser(event);
	}

	const resolveWithParaglide = (locale: Locale) => {
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
	};

	if (request.method !== 'GET' && request.method !== 'HEAD') {
		event.locals.lang = detectedLocale;
		setCookie(event, detectedLocale);
		return resolveWithParaglide(detectedLocale);
	}

	if (isBot) {
		event.locals.lang = (pathLocale ?? 'en') as Locale;
		setCookie(event, event.locals.lang as Locale);
		return resolveWithParaglide(event.locals.lang as Locale);
	}

	if (pathname === '/en' || pathname.startsWith('/en/')) {
		const stripped = pathname === '/en' ? '/' : pathname.slice(3);
		throw redirect(302, `${stripped}${event.url.search}`);
	}

	if (!pathLocale) {
		const cookieLang = event.cookies.get(LOCALE_COOKIE_NAME) as Locale | null;
		const targetLocale =
			cookieLang && SUPPORTED_LOCALES.includes(cookieLang) ? cookieLang : detectedLocale;

		event.locals.lang = targetLocale;
		setCookie(event, targetLocale);

		if (targetLocale !== 'en') {
			throw redirect(302, `/${targetLocale}${pathname === '/' ? '' : pathname}${event.url.search}`);
		}
		return resolveWithParaglide(targetLocale);
	}

	event.locals.lang = pathLocale;
	setCookie(event, pathLocale);
	return resolveWithParaglide(pathLocale);
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
	handleParaglide
	// initServer
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
function hasLocalePrefix(path: string): boolean {
	return SUPPORTED_LOCALES.some((l) => path === `/${l}` || path.startsWith(`/${l}/`));
}

function getLocaleFromPath(pathname: string): Locale | null {
	const match = pathname.match(/^\/(en|id)(\/|$)/);
	return match ? (match[1] as Locale) : null;
}
function detectLocaleFromBrowser(event: RequestEvent): Locale {
	const accept = event.request.headers.get('accept-language');
	const l = accept?.split(',')[0].split('-')[0] as Locale;
	const supported = SUPPORTED_LOCALES.includes(l);
	return supported ? l : 'en';
}

function setCookie(event: RequestEvent, locale: Locale) {
	event.cookies.set('PARAGLIDE_LOCALE', locale, {
		httpOnly: true,
		secure: NODE_ENV === 'production',
		sameSite: 'lax',
		path: '/',
		maxAge: 60 * 60 * 24 * 7
	});
}
