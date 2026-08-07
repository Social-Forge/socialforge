import { redirect } from '@sveltejs/kit';
import { PUBLIC_API_URL } from '$env/static/public';
import { NODE_ENV } from '$env/static/private';
import { isSocialAuthProvider } from '$lib/auth/social-providers';
import { localizePath } from '$lib/utils/localize-path';

const OAUTH_REDIRECT_COOKIE = 'oauth_redirect_target';
const OAUTH_MODE_COOKIE = 'oauth_auth_mode';
const OAUTH_SESSION_COOKIE = 'oauth_backend_session';

export const GET = async ({ params, url, cookies, locals, fetch }) => {
        const provider = params.provider;

        if (!provider || !isSocialAuthProvider(provider)) {
                throw redirect(302, localizePath('/signin', locals.lang));
        }

        const redirectTarget = sanitizeRedirectTarget(url.searchParams.get('redirect'));
        const mode = url.searchParams.get('mode') === 'signup' ? 'signup' : 'signin';

        cookies.set(OAUTH_REDIRECT_COOKIE, redirectTarget, {
                path: '/',
                httpOnly: true,
                secure: NODE_ENV === 'production',
                sameSite: 'lax',
                maxAge: 60 * 10
        });

        cookies.set(OAUTH_MODE_COOKIE, mode, {
                path: '/',
                httpOnly: true,
                secure: NODE_ENV === 'production',
                sameSite: 'lax',
                maxAge: 60 * 10
        });

        const backendResponse = await fetch(`${PUBLIC_API_URL}/auth/oauth/${provider}`, {
                method: 'GET',
                redirect: 'manual'
        });
        const redirectLocation = backendResponse.headers.get('location');
        const backendCookieHeader = serializeCookieHeader(readSetCookies(backendResponse));

        if (!redirectLocation) {
                throw redirect(
                        302,
                        localizePath(`/${mode}?oauth_error=${encodeURIComponent('OAuth provider redirect failed')}`, locals.lang)
                );
        }

        if (backendCookieHeader) {
                cookies.set(OAUTH_SESSION_COOKIE, backendCookieHeader, {
                        path: '/',
                        httpOnly: true,
                        secure: NODE_ENV === 'production',
                        sameSite: 'lax',
                        maxAge: 60 * 10
                });
        }

        throw redirect(302, redirectLocation);
};

function sanitizeRedirectTarget(target: string | null): string {
        if (!target || !target.startsWith('/') || target.startsWith('//')) {
                return '/app/chats';
        }

        return target;
}

function readSetCookies(response: Response): string[] {
        const headers = response.headers as Headers & {
                getSetCookie?: () => string[];
                raw?: () => Record<string, string[]>;
        };

        if (typeof headers.getSetCookie === 'function') {
                return headers.getSetCookie();
        }

        if (typeof headers.raw === 'function') {
                return headers.raw()['set-cookie'] ?? [];
        }

        const headerValue = response.headers.get('set-cookie');
        return headerValue ? [headerValue] : [];
}

function serializeCookieHeader(setCookieHeaders: string[]): string {
        return setCookieHeaders
                .map((value) => value.split(';', 1)[0]?.trim())
                .filter((value): value is string => Boolean(value))
                .join('; ');
}
