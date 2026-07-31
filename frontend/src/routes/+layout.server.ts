import type { MetaTag, LinkTag } from 'svelte-meta-tags';
import { defineBaseMetaTags } from 'svelte-meta-tags';
import { NODE_ENV } from '$env/static/private';
import { locales as SUPPORTED_LOCALES } from '$lib/paraglide/runtime';

export const load = async ({ locals, platform, url }) => {
	const { user, lang } = locals;
	const defaultOrigin = locals.origin ?? url.origin;

	let canonicalUrl = defaultOrigin;
	if (NODE_ENV && NODE_ENV === 'production' && canonicalUrl.startsWith('http://')) {
		canonicalUrl = canonicalUrl.replace('http://', 'https://');
	}
	const alternates = SUPPORTED_LOCALES.map((lang) => ({
		lang,
		href: `${canonicalUrl}/${lang}`
	}));
	const normalizedAlternates = alternates.map((alt) => ({
		...alt,
		href: normalizeUrl(alt.href)
	}));

	const baseTags = defineBaseMetaTags({
		title: 'Social Forge - Multi-agent Customer Service and Omnichannel CRM',
		titleTemplate: '%s | Social Forge',
		description:
			'Social Forge is the best multi-agent and multi-channel WhatsApp CRM to serve customers on WhatsApp, Instagram, Facebook, and your website.',
		keywords: [
			'crm whatsapp',
			'omnichannel crm',
			'multi-agent customer service',
			'whatsapp blast',
			'integrasi instagram facebook'
		],
		canonical: canonicalUrl + url.pathname,
		additionalMetaTags: [
			{
				name: 'viewport',
				content: 'width=device-width, initial-scale=1.0'
			},
			{
				property: 'dc:creator',
				content: 'Social Forge'
			},
			{
				name: 'application-name',
				content: 'Social Forge'
			},
			{
				httpEquiv: 'x-ua-compatible',
				content: 'IE=edge'
			},
			{
				name: 'description',
				content:
					'Social Forge is the best multi-agent and multi-channel WhatsApp CRM to serve customers on WhatsApp, Instagram, Facebook, and your website.'
			},
			{
				name: 'mobile-web-app-capable',
				content: 'yes'
			},
			{
				name: 'mobile-web-app-status-bar-style',
				content: 'black-translucent'
			},
			{
				name: 'mobile-web-app-title',
				content: 'Social Forge'
			},
			{
				name: 'mobile-web-app-icon',
				content: '/favicon.ico'
			}
		] as MetaTag[],
		additionalLinkTags: [
			{
				rel: 'canonical',
				href: canonicalUrl + url.pathname
			},
			{
				rel: 'alternate',
				hreflang: 'x-default',
				href: canonicalUrl + url.pathname
			},
			{
				rel: 'icon',
				type: 'image/x-icon',
				sizes: '96x96',
				href: '/logo.png'
			},
			{
				rel: 'icon',
				type: 'image/png',
				sizes: '32x32',
				href: '/favicon-32x32.png'
			},
			{
				rel: 'icon',
				type: 'image/png',
				sizes: '16x16',
				href: '/favicon-16x16.png'
			},
			{
				rel: 'icon',
				type: 'image/png',
				sizes: '192x192',
				href: '/android-chrome-192x192.png'
			},
			{
				rel: 'icon',
				type: 'image/png',
				sizes: '512x512',
				href: '/android-chrome-512x512.png'
			},
			{
				rel: 'apple-touch-icon',
				type: 'image/png',
				sizes: '180x180',
				href: '/apple-touch-icon.png'
			}
		] as LinkTag[],
		openGraph: {
			type: 'website',
			url: canonicalUrl + url.pathname,
			locale: 'en_IE',
			title: 'Social Forge',
			description:
				'Social Forge is the best multi-agent and multi-channel WhatsApp CRM to serve customers on WhatsApp, Instagram, Facebook, and your website.',
			siteName: 'Social Forge',
			images: [
				{
					url: '/logo.png',
					width: 800,
					height: 600,
					alt: 'Social Forge Cover Image',
					type: 'image/png'
				},
				{
					url: '/logo.png',
					width: 512,
					height: 512,
					alt: 'Social Forge Android Chrome Icon',
					type: 'image/x-icon'
				}
			],
			profile: {
				firstName: 'Social Forge',
				lastName: 'Social Forge',
				username: 'socialforge'
			}
		}
	});

	return {
		...baseTags,
		user,
		lang,
		canonicalUrl,
		alternates: normalizedAlternates
	};
};

function normalizeUrl(urlString: string): string {
	try {
		const url = new URL(urlString);
		if (NODE_ENV === 'production' && url.protocol === 'http:') {
			url.protocol = 'https:';
		}
		return url.href;
	} catch {
		return urlString;
	}
}
