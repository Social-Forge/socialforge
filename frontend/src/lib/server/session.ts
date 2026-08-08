import { NODE_ENV } from '$env/static/private';
import type { RequestEvent } from '@sveltejs/kit';
import { BaseHandler } from './base';
import { redirect, type Cookies } from '@sveltejs/kit';

export class SessionHelper extends BaseHandler {
	constructor(protected readonly event: RequestEvent) {
		super(event);
	}

	setAuthCookies = (
		tokens: {
			accessToken: string;
			refreshToken: string;
		} | null,
		expiresAccIn: number,
		expiresRefreshIn: number
	) => {
		const isProduction = NODE_ENV === 'production';
		if (!tokens) {
			this.event.cookies.delete('access_token', { path: '/' });
			this.event.cookies.delete('refresh_token', { path: '/' });
			return;
		}
		this.event.cookies.set('access_token', tokens.accessToken, {
			path: '/',
			httpOnly: false,
			secure: isProduction,
			sameSite: 'lax',
			maxAge: expiresAccIn,
			expires: new Date(Date.now() + expiresAccIn * 1000)
		});

		this.event.cookies.set('refresh_token', tokens.refreshToken, {
			path: '/',
			httpOnly: false,
			secure: isProduction,
			sameSite: 'lax',
			maxAge: expiresRefreshIn,
			expires: new Date(Date.now() + expiresRefreshIn * 1000)
		});
	};
	validateCSRF = (): boolean => {
		const cookieToken =
			this.event.cookies.get('csrf_token') || this.event.cookies.get('XSRF-TOKEN');
		const headerToken =
			this.event.request.headers.get('X-XSRF-TOKEN') ||
			this.event.request.headers.get('X-Xsrf-Token');
		return !!cookieToken && cookieToken === headerToken;
	};
	setSecurityHeaders = () => {
		this.event.request.headers.set('Cache-Control', 'no-store, max-age=0');
		this.event.request.headers.set('CDN-Cache-Control', 'max-age=60, stale-while-revalidate=300');
	};
	clearAuthCookies = () => {
		this.event.cookies.delete('access_token', { path: '/' });
		this.event.cookies.delete('refresh_token', { path: '/' });
		this.event.cookies.delete('twofa_session_id', { path: '/' });
                this.event.cookies.delete('oauth_redirect_target', { path: '/' });
                this.event.cookies.delete('oauth_auth_mode', { path: '/' });
                this.event.cookies.delete('oauth_backend_session', { path: '/' });
                this.event.cookies.delete('csrf_token', { path: '/' });
                this.event.cookies.delete('XSRF-TOKEN', { path: '/' });
	};
	isAuthenticated = async (): Promise<boolean> => {
		const accessToken = this.getAccessToken();
		const refreshToken = this.getRefreshToken();

		if (accessToken && !this.isTokenExpired(accessToken)) {
			return true;
		}
		if (refreshToken && !this.isTokenExpired(refreshToken)) {
			return true;
		}
		return false;
	};
	handleUnauthorized = (redirectTo = '/signin') => {
		this.clearAuthCookies();
		return redirect(302, `/signin?from=${encodeURIComponent(redirectTo)}`);
	};
	getAccessToken = (): string | undefined => {
		return this.event.cookies.get('access_token');
	};
	getRefreshToken = (): string | undefined => {
		return this.event.cookies.get('refresh_token');
	};
	getTwoSessionToken = (): string | undefined => {
		const twoFaSession =
			this.event.request.headers.get('X-2FA-Session') || this.event.cookies.get('twofa_session_id');
		return twoFaSession;
	};
	isTokenExpired = (token: string): boolean => {
		try {
			const payload = JSON.parse(atob(token.split('.')[1]));
			const exp = payload.exp * 1000;
			const now = Date.now();
			const buffer = 5 * 60 * 1000;
			return now >= exp - buffer;
		} catch {
			return true;
		}
	};
}
