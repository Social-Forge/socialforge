<script lang="ts" module>
	interface Emoji {
		annotation?: string;
		emoji?: string;
		group?: number;
		hexcode: string;
		order?: number;
		shortcodes?: string[];
		skins?: Emoji[];
		tags?: string[];
		unicode: string;
		version?: number;
	}

	interface Props {
		value?: string;
		placeholder?: string;
		disabled?: boolean;
		locale?: string; // 'en', 'id', etc
		class?: string;
		onselect?: (emoji: Emoji) => void;
	}
</script>

<script lang="ts">
	import { cn } from '$lib/utils.js';
	import { onMount, type Component } from 'svelte';
	import * as Tabs from '$lib/components/ui/tabs/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { ScrollArea } from '$lib/components/ui/scroll-area/index.js';
	import enEmojis from 'emojibase-data/en/compact.json';
	import enMessages from 'emojibase-data/en/messages.json';
	import {
		Clock,
		Smile,
		User,
		Leaf,
		Coffee,
		MapPin,
		Activity,
		Box,
		Hash,
		Flag,
		Circle,
		Inbox,
		SearchX,
		type LucideProps
	} from '@lucide/svelte';

	let {
		value = $bindable(''),
		placeholder = 'Search emoji...',
		disabled = false,
		locale = 'en',
		class: className = '',
		onselect
	}: Props = $props();

	let isOpen = $state(false);
	let searchQuery = $state('');
	let selectedCategory = $state('frequent');
	let recentEmojis = $state<Emoji[]>([]);
	let isLoaded = $state(false);
	let emojiData = $state<Emoji[]>([]);

	const groups = enMessages.groups.map((g: any) => ({
		key: g.key,
		label: g.message,
		order: g.order
	}));

	const initEmojis = async () => {
		emojiData = enEmojis as Emoji[];
		isLoaded = true;
	};
	onMount(() => {
		initEmojis();
	});

	const filteredEmojis = $derived(() => {
		if (!isLoaded) return [];

		let emojis = [...emojiData];

		if (selectedCategory !== 'frequent' && selectedCategory !== 'search') {
			const groupIndex = groups.findIndex((g) => g.key === selectedCategory);
			if (groupIndex !== -1) {
				emojis = emojis.filter((e) => e.group === groupIndex);
			}
		}

		if (searchQuery.trim()) {
			const query = searchQuery.toLowerCase();
			emojis = emojis.filter(
				(e) =>
					e.annotation?.toLowerCase().includes(query) ||
					e.tags?.some((tag) => tag.toLowerCase().includes(query)) ||
					e.shortcodes?.some((code) => code.toLowerCase().includes(query))
			);

			if (searchQuery) {
				selectedCategory = 'search';
				return emojis;
			}
		}

		if (selectedCategory === 'frequent') {
			const recentHexcodes = new Set(recentEmojis.map((e) => e.hexcode));
			return emojiData.filter((e) => recentHexcodes.has(e.hexcode));
		}

		return emojis;
	});

	const selectEmoji = (emoji: Emoji) => {
		const existing = recentEmojis.findIndex((e) => e.hexcode === emoji.hexcode);
		if (existing !== -1) {
			recentEmojis.splice(existing, 1);
		}
		recentEmojis.unshift(emoji);

		if (recentEmojis.length > 30) {
			recentEmojis = recentEmojis.slice(0, 30);
		}

		const emojiChar = emoji.emoji || emoji.unicode;
		value = emojiChar; // Update binding
		onselect?.(emoji);

		isOpen = false;
		searchQuery = '';
	};

	const getCategoryIcon = (key: string): Component<LucideProps, {}, ''> => {
		const icons: Record<string, Component<LucideProps, {}, ''>> = {
			frequent: Clock,
			'smileys-emotion': Smile,
			'people-body': User,
			'animals-nature': Leaf,
			'food-drink': Coffee,
			'travel-places': MapPin,
			activities: Activity,
			objects: Box,
			symbols: Hash,
			flags: Flag
		};
		return icons[key] || Circle;
	};

	const handleKeydown = (e: KeyboardEvent) => {
		if (e.key === 'Escape') {
			isOpen = false;
			searchQuery = '';
		}
	};

	$effect(() => {
		if (recentEmojis.length > 0) {
			try {
				localStorage.setItem('recent-emojis', JSON.stringify(recentEmojis.map((e) => e.hexcode)));
			} catch (e) {
				// ignore
			}
		}
	});

	onMount(() => {
		try {
			const stored = localStorage.getItem('recent-emojis');
			if (stored) {
				const hexcodes = JSON.parse(stored) as string[];
				const emojis = hexcodes
					.map((hex: string) => emojiData.find((e) => e.hexcode === hex))
					.filter((e): e is Emoji => e !== undefined);
				recentEmojis = emojis;
			}
		} catch (e) {
			// ignore
		}
	});
</script>

<div class={cn('w-full p-0', className)}>
	<div class="flex flex-col">
		<div class="border-b p-3">
			<Input
				bind:value={searchQuery}
				{placeholder}
				class="w-full"
				onkeydown={handleKeydown}
				{disabled}
			/>
		</div>

		<!-- Scroll Area -->
		<ScrollArea class="h-75">
			<!-- Tabs (Hanya muncul jika tidak sedang search) -->
			<Tabs.Root
				bind:value={selectedCategory}
				class={cn('w-full', className, searchQuery ? 'hidden' : '')}
			>
				<Tabs.List
					class="h-auto w-full flex-wrap justify-start rounded-none border-b bg-transparent p-0"
				>
					<!-- Recent Tab -->
					<Tabs.Trigger
						value="frequent"
						class="rounded-none border-b-2 border-transparent px-3 py-2 data-[state=active]:border-primary data-[state=active]:bg-transparent"
					>
						<Clock class="h-4 w-4" />
						<span class="sr-only">Recent</span>
					</Tabs.Trigger>

					<!-- Category Tabs -->
					{#each groups as group (group.key)}
						<Tabs.Trigger
							value={group.key}
							class="rounded-none border-b-2 border-transparent px-3 py-2 data-[state=active]:border-primary data-[state=active]:bg-transparent"
						>
							{@const Icon = getCategoryIcon(group.key)}
							<Icon class="h-4 w-4" />
							<span class="sr-only">{group.label}</span>
						</Tabs.Trigger>
					{/each}
				</Tabs.List>

				<!-- Tabs Content -->
				<div class="p-2">
					{#each [{ key: 'frequent', label: 'Recent' }, ...groups] as group (group.key)}
						<Tabs.Content value={group.key} class="mt-0">
							{#if filteredEmojis().length === 0}
								<div class="py-8 text-center text-muted-foreground">
									<Inbox class="mx-auto mb-2 h-8 w-8" />
									<p class="text-sm">No emojis found</p>
								</div>
							{:else}
								<div class="grid grid-cols-8 gap-1">
									{#each filteredEmojis() as emoji (emoji.hexcode)}
										<button
											type="button"
											class="rounded p-1 text-xl transition-colors hover:bg-muted focus:ring-2 focus:ring-primary focus:outline-none"
											onclick={() => selectEmoji(emoji)}
										>
											{emoji.emoji || emoji.unicode}
										</button>
									{/each}
								</div>
							{/if}
						</Tabs.Content>
					{/each}
				</div>
			</Tabs.Root>

			<!-- Tampilan Khusus Search (Jika ada query) -->
			{#if searchQuery}
				<div class="p-2">
					{#if filteredEmojis().length === 0}
						<div class="py-8 text-center text-muted-foreground">
							<SearchX class="mx-auto mb-2 h-8 w-8" />
							<p class="text-sm">No emojis found for "{searchQuery}"</p>
						</div>
					{:else}
						<div class="grid grid-cols-8 gap-1">
							{#each filteredEmojis() as emoji (emoji.hexcode)}
								<button
									type="button"
									class="rounded p-1 text-xl transition-colors hover:bg-muted focus:ring-2 focus:ring-primary focus:outline-none"
									onclick={() => selectEmoji(emoji)}
								>
									{emoji.emoji || emoji.unicode}
								</button>
							{/each}
						</div>
					{/if}
				</div>
			{/if}
		</ScrollArea>
	</div>
</div>
