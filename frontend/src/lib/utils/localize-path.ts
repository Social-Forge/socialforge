import { goto } from '$app/navigation';
import { type Locale } from '$lib/paraglide/runtime.js';

export const LanguageLabels: Partial<Record<Locale, string>> = {
	en: 'English',
	id: 'Bahasa Indonesia'
};

export function localizePath(path: string, lang: string) {
	return `/${lang}${path.startsWith('/') ? path : `/${path}`}`;
}

export function gotoLocale(path: string, lang: string) {
	const href = `/${lang}${path.startsWith('/') ? path : `/${path}`}`;
	goto(href);
}

export function removeLocaleFromPath(path: string): string {
	const match = path.match(/^\/(en|id|es|ru|pt|fr|de|zh|hi|ar|ja|tr|vi|th|el|it)(\/|$)/);
	if (match) {
		const remaining = path.slice(match[0].length - (match[2] === '/' ? 1 : 0));
		return remaining || '/';
	}
	return path || '/';
}
