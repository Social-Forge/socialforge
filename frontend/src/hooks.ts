import type { Reroute } from '@sveltejs/kit';
import { deLocalizeUrl } from '$lib/paraglide/runtime';

export const reroute: Reroute = ({ url }) => {
	const pathname = url.pathname;
	if (
		pathname.startsWith('/api') ||
		pathname.startsWith('/_app') ||
		pathname.startsWith('/favicon.ico') ||
		pathname.includes('.')
	) {
		return;
	}
	return deLocalizeUrl(url).pathname;
};

