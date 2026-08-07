export const SOCIAL_AUTH_PROVIDERS = ['google', 'facebook', 'github'] as const;

export type SocialAuthProvider = (typeof SOCIAL_AUTH_PROVIDERS)[number];

export function isSocialAuthProvider(value: string): value is SocialAuthProvider {
        return SOCIAL_AUTH_PROVIDERS.includes(value as SocialAuthProvider);
}
