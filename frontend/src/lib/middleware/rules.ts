import { localizePath } from '$lib/utils/localize-path';

export const ROLE_LEVELS = {
        SUPER_ADMIN: 0,
        TENANT_OWNER: 1,
        SUPERVISOR: 2,
        AGENT: 3
} as const;

export const allowedWhenAuthenticated = ['/signout'];
export const authRoutes = [
        '/signin',
        '/signup',
        '/verify-email',
        '/reset-password',
        '/two-factor',
        '/forgot-password',
        '/confirm'
] as const;

export const protectedAppRoutes = ['/app'] as const;
export const restrictedSuperAdminRoutes = ['/app/admin'] as const;
export const restrictedTenantOwnerRoutes = [
        '/app/settings',
        '/app/divisions',
        '/app/agents',
        '/app/channels',
        '/app/analytics'
] as const;

export const publicApiRoutes = ['/api/auth/email-validate'] as const;
export const adminApiRoutes = ['/api/admin'] as const;
export const tenantMemberApiRoutes = ['/api/chats'] as const;
export const tenantOwnerApiRoutes = [
        '/api/tenants',
        '/api/divisions',
        '/api/agents',
        '/api/channels',
        '/api/settings',
        '/api/analytics'
] as const;

export function matchesRoutePrefix(pathname: string, prefix: string): boolean {
        return pathname === prefix || pathname.startsWith(`${prefix}/`);
}

export function matchesAnyRoute(pathname: string, routes: readonly string[]): boolean {
        return routes.some((route) => matchesRoutePrefix(pathname, route));
}

export function isSuperAdmin(roleLevel?: number | null): boolean {
        return roleLevel === ROLE_LEVELS.SUPER_ADMIN;
}

export function canManageTenant(roleLevel?: number | null): boolean {
        return roleLevel === ROLE_LEVELS.SUPER_ADMIN || roleLevel === ROLE_LEVELS.TENANT_OWNER;
}

export function sanitizeRedirectTarget(target: string | null | undefined, fallback = '/app/chats'): string {
        if (!target || !target.startsWith('/') || target.startsWith('//')) {
                return fallback;
        }

        return target;
}

export function buildSignInRedirect(pathname: string, lang: string): string {
        return `${localizePath('/signin', lang)}?redirect=${encodeURIComponent(pathname)}`;
}
