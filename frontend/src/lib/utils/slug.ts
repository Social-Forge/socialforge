/** Slugs that can never be used by campaigns (they collide with app routes). */
export const RESERVED_SLUGS = new Set([
	'admin',
	'api',
	'about',
	'contact',
	'dashboard',
	'dmca',
	'favicon.ico',
	'forgot-password',
	'go',
	'login',
	'logout',
	'payout-rates',
	'privacy',
	'register',
	'report',
	'reset-password',
	'robots.txt',
	's',
	'sitemap.xml',
	'terms'
]);

const ALPHABET = 'abcdefghijkmnpqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ23456789'; // no 0/O/1/l/I

/** Cryptographically random, URL-safe, unambiguous slug. */
export function generateSlug(length = 7): string {
	const bytes = new Uint8Array(length);
	crypto.getRandomValues(bytes);
	let out = '';
	for (const b of bytes) out += ALPHABET[b % ALPHABET.length];
	return out;
}

export const SLUG_PATTERN = /^[a-zA-Z0-9-_]{3,64}$/;

export function isValidCustomSlug(slug: string): boolean {
	return SLUG_PATTERN.test(slug) && !RESERVED_SLUGS.has(slug.toLowerCase());
}
