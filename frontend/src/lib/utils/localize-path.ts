import { baseLocale, locales as SUPPORTED_LOCALES, type Locale } from '$lib/paraglide/runtime.js';

export const BASE_LOCALE = baseLocale;
const localePrefixPattern = new RegExp(
        `^/(${SUPPORTED_LOCALES.map(escapeRegExp).join('|')})(?=/|$)`,
        'i'
);

export const LanguageLabels: Partial<Record<Locale, string>> = {
	en: 'English',
	id: 'Bahasa Indonesia'
};

export function getLocaleFromPath(path: string): Locale | null {
        const pathname = getPathname(path);
        const match = pathname.match(localePrefixPattern);
        if (!match) {
                return null;
        }

        return toSupportedLocale(match[1]);
}

export function localizePath(path: string, lang: string): string {
        const locale = toSupportedLocale(lang) ?? BASE_LOCALE;
        const pathWithoutLocale = removeLocaleFromPath(path);
        const pathname = normalizePath(getPathname(pathWithoutLocale));
        const suffix = pathWithoutLocale.slice(getPathname(pathWithoutLocale).length);

        if (locale === BASE_LOCALE) {
                return `${pathname}${suffix}`;
        }

        return pathname === '/' ? `/${locale}${suffix}` : `/${locale}${pathname}${suffix}`;
}

export function removeLocaleFromPath(path: string): string {
        const rawPathname = getPathname(path);
        const pathname = normalizePath(rawPathname);
        const suffix = path.slice(rawPathname.length);
        const match = pathname.match(localePrefixPattern);

        if (!match) {
                return `${pathname}${suffix}`;
        }

        const localizedPrefix = match[0];
        const remaining = pathname.slice(localizedPrefix.length);
        const normalizedRemaining = remaining.startsWith('/') ? remaining : `/${remaining}`;

        return `${normalizedRemaining === '/' ? '/' : normalizedRemaining}${suffix}`;
}

export function hasLocalePrefix(path: string): boolean {
        return getLocaleFromPath(path) !== null;
}

function getPathname(path: string): string {
        const pathnameMatch = path.match(/^[^?#]*/);
        return pathnameMatch?.[0] || '/';
}

function normalizePath(path: string): string {
        if (!path) {
                return '/';
        }

        return path.startsWith('/') ? path : `/${path}`;
}

function toSupportedLocale(value: string): Locale | null {
        const locale = SUPPORTED_LOCALES.find((candidate) => candidate.toLowerCase() === value.toLowerCase());
        return locale ?? null;
}

function escapeRegExp(value: string): string {
        return value.replace(/[.*+?^${}()|[\]\\]/g, '\\$&');
}
