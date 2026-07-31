<script lang="ts">
	import { LanguageSwitcher } from '$lib/components/extras/language-switcher';
	import { LightSwitch } from '$lib/components/extras/light-switch';
	import { locales as availableLocales, localizeHref } from '$lib/paraglide/runtime';
	import { LanguageLabels } from '$lib/utils/localize-path.js';

	let { data, children } = $props();

	let currentLang = $derived(data.lang);

	const languages = availableLocales.map((code) => ({
		code,
		label: LanguageLabels[code] ?? code.toUpperCase()
	}));
</script>

<div
	class="relative flex min-h-svh flex-col items-center justify-center gap-6 bg-muted bg-linear-to-br from-primary/20 via-primary/10 to-primary/30 p-6 md:p-10"
>
	<div class="fixed top-4 right-4 z-10">
		<div class="flex items-center justify-center gap-2">
			<LightSwitch />
			<LanguageSwitcher {languages} bind:value={currentLang} variant="outline" />
		</div>
	</div>
	<div class="flex w-full max-w-md flex-col gap-6">
		<a href={localizeHref('/')} class="flex items-center justify-center gap-2">
			<div class="flex items-center justify-center">
				<img src="/logo.png" alt="Social Forge" class="h-10 w-auto object-cover" />
			</div>
			<span class="text-2xl font-bold">Social Forge</span>
		</a>
		{@render children?.()}
	</div>
</div>
