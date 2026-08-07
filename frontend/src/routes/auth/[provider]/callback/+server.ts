import { redirect } from '@sveltejs/kit';
import { NODE_ENV } from '$env/static/private';
import { isSocialAuthProvider } from '$lib/auth/social-providers';
import { localizePath } from '$lib/utils/localize-path';

const OAUTH_REDIRECT_COOKIE = 'oauth_redirect_target';
const OAUTH_MODE_COOKIE = 'oauth_auth_mode';
const OAUTH_SESSION_COOKIE = 'oauth_backend_session';

export const GET = async ({ params, url, locals, cookies }) => {
        const provider = params.provider;
        const defaultMode = cookies.get(OAUTH_MODE_COOKIE) === 'signup' ? 'signup' : 'signin';
        const authPage = defaultMode === 'signup' ? '/signup' : '/signin';
        const redirectTarget = sanitizeRedirectTarget(cookies.get(OAUTH_REDIRECT_COOKIE));
        const backendSessionCookie = cookies.get(OAUTH_SESSION_COOKIE);

        if (!provider || !isSocialAuthProvider(provider)) {
                throw redirect(302, localizePath(authPage, locals.lang));
        }

        const response = await locals.helper.auth.oAuthCallback(
                provider,
                url.searchParams,
                backendSessionCookie ? { Cookie: backendSessionCookie } : undefined
        );

        cookies.delete(OAUTH_REDIRECT_COOKIE, { path: '/' });
        cookies.delete(OAUTH_MODE_COOKIE, { path: '/' });
        cookies.delete(OAUTH_SESSION_COOKIE, {
                path: '/',
                secure: NODE_ENV === 'production',
                sameSite: 'lax'
        });

        if (!response.success || !response.data?.access_token || !response.data?.refresh_token) {
                const errorMessage = encodeURIComponent(response.message || 'OAuth sign in failed');
                throw redirect(302, localizePath(`${authPage}?oauth_error=${errorMessage}`, locals.lang));
        }

        locals.helper.session.setAuthCookies(
                {
                        accessToken: response.data.access_token,
                        refreshToken: response.data.refresh_token
                },
                response.data.expires_in || 60 * 60 * 24,
                response.data.expires_refresh_in || 60 * 60 * 24 * 7
        );

        throw redirect(302, localizePath(redirectTarget, locals.lang));
};

function sanitizeRedirectTarget(target: string | undefined): string {
        if (!target || !target.startsWith('/') || target.startsWith('//')) {
                return '/app/chats';
        }

        return target;
}
