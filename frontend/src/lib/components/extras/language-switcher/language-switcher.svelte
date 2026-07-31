<script lang="ts" module>
	export type Language = {
		/** Language code (e.g., 'en', 'de') */
		code: string;
		/** Display name (e.g., 'English', 'Deutsch') */
		label: string;
	};

	export type LanguageSwitcherProps = {
		/** List of available languages */
		languages: Language[];

		/** Current selected language code */
		value?: string;

		/** Dropdown alignment */
		align?: 'start' | 'center' | 'end';

		/** Button variant */
		variant?: 'outline' | 'ghost';

		/** Called when the language changes */
		onChange?: (code: string) => void;

		class?: string;
	};
</script>

<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import GlobeIcon from '@lucide/svelte/icons/globe';
	import * as DropdownMenu from '$lib/components/extras/dropdown-menu';
	import { buttonVariants } from '$lib/components/extras/button';
	import { cn } from '$lib/utils.js';
	import { LanguageLabels } from '$lib/utils/localize-path.js';
	import { setLocale, isLocale } from '$lib/paraglide/runtime';

	let {
		languages = [],
		value = $bindable(''),
		align = 'end',
		variant = 'ghost',
		onChange,
		class: className
	}: LanguageSwitcherProps = $props();

	let pathname = $derived(page.url.pathname);

	// svelte-ignore state_referenced_locally
	if (value === '' && languages.length > 0) {
		// svelte-ignore state_referenced_locally
		value = languages[0].code;
	}

	function setFlagUrl(code?: string): string {
		if (!code) return '';
		const flags: Record<string, string> = {
			en: 'us',
			id: 'id',
			ja: 'jp',
			ar: 'sa',
			zh: 'cn',
			hi: 'in',
			el: 'gr',
			ko: 'kr',
			vi: 'vn',
			ms: 'my',
			tl: 'ph',
			sv: 'se',
			pl: 'pl',
			cs: 'cz',
			ro: 'ro',
			hu: 'hu',
			fi: 'fi',
			da: 'dk',
			no: 'no',
			he: 'il',
			bn: 'bd',
			es: 'es',
			ru: 'ru',
			pt: 'pt',
			fr: 'fr',
			de: 'de',
			tr: 'tr',
			it: 'it',
			th: 'th'
		};
		const flag = flags[code] || code;
		return `https://flagicons.lipis.dev/flags/4x3/${flag}.svg`;
	}

	function handleLanguageChange(code: string) {
		if (!isLocale(code)) return;

		setLocale(code);

		value = code;

		if (onChange) {
			onChange(code);
		}

		const pathWithoutLocale = removeLocaleFromPath(pathname);
		const newPath =
			code === 'en'
				? `/${pathWithoutLocale}${pathWithoutLocale === '' ? '' : ''}`
				: `/${code}${pathWithoutLocale}`;

		goto(newPath, { replaceState: true, invalidateAll: true });
	}

	function removeLocaleFromPath(path: string): string {
		const match = path.match(/^\/(en|id|es)(\/|$)/);
		if (match) {
			const remaining = path.slice(match[0].length - (match[2] === '/' ? 1 : 0));
			return remaining || '/';
		}
		return path || '/';
	}
</script>

<DropdownMenu.Root>
	<DropdownMenu.Trigger
		class={cn(buttonVariants({ variant, size: 'icon' }), className)}
		aria-label="Change language"
	>
		{#if value}
			<img
				src={setFlagUrl(value)}
				alt={LanguageLabels[value as keyof typeof LanguageLabels] || value.toUpperCase()}
				class="size-4"
			/>
		{:else}
			<GlobeIcon class="size-4" />
		{/if}
		<span class="sr-only">Change language</span>
	</DropdownMenu.Trigger>
	<DropdownMenu.Content {align} class="w-full">
		<DropdownMenu.RadioGroup bind:value onValueChange={handleLanguageChange}>
			{#each languages.sort((a, b) => a.code.localeCompare(b.code)) as language (language.code)}
				<DropdownMenu.RadioItem value={language.code} class="flex items-center gap-2">
					<img src={setFlagUrl(language.code)} alt={language.label} class="no-copy size-4" />
					{language.label}
				</DropdownMenu.RadioItem>
			{/each}
		</DropdownMenu.RadioGroup>
	</DropdownMenu.Content>
</DropdownMenu.Root>
