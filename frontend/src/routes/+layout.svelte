<script lang="ts">
	import type { Pathname } from '$app/types';
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import { MetaTags, deepMerge } from 'svelte-meta-tags';
	import { locales, localizeHref } from '$lib/paraglide/runtime';
	import './layout.css';
	import { ModeWatcher } from 'mode-watcher';
	import { ToastContent } from '$lib/components/toast';
	import { SvelteKitTopLoader } from 'sveltekit-top-loader';

	let { data, children } = $props();
	let metaTags = $derived(deepMerge(data.baseMetaTags, page.data.pageMetaTags));
</script>

<MetaTags {...metaTags} />
<ModeWatcher />
<ToastContent />
<SvelteKitTopLoader color="#007a55" />

<main class="min-h-screen antialiased">
	{@render children?.()}
</main>

<div style="display:none">
	{#each locales as locale (locale)}
		<!-- @ts-expect-error -->
		<a href={resolve(localizeHref(page.url.pathname, { locale }) as any)}>{locale}</a>
	{/each}
</div>
