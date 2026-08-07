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
</script>

<script lang="ts">
	import {
		Plus,
		Smile,
		Send,
		Paperclip,
		Image,
		Video,
		FileText,
		MapPinXInside,
		XIcon
	} from '@lucide/svelte';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu/index.js';
	import { buttonVariants, Button } from '$lib/components/ui/button';
	import { EmojiPicker } from '$lib/components/ui/emoji-picker/index.js';

	let { onsend }: { onsend?: (text: string) => void } = $props();

	let value = $state('');
	let showEmojiPicker = $state(false);
	let selectedEmoji = $state('');

	function handleSend() {
		if (!value.trim()) return;
		onsend?.(value);
		value = '';
	}

	function handleKeydown(e: KeyboardEvent) {
		if (e.key === 'Enter' && !e.shiftKey) {
			e.preventDefault();
			handleSend();
		}
	}

	const onEmojiClick = (emoji: Emoji) => {
		const emojiChar = emoji.emoji || emoji.unicode;
		value += emojiChar;
	};
</script>

<div class="border-t bg-background px-2 py-4">
	{#if showEmojiPicker}
		<div class="mb-2 rounded-lg border bg-card p-2">
			<div class="flex w-full justify-end">
				<Button variant="outline" size="icon-xs" onclick={() => (showEmojiPicker = false)}>
					<XIcon class="h-5 w-5" />
				</Button>
			</div>
			<EmojiPicker bind:value={selectedEmoji} class="shrink-0" onselect={onEmojiClick} />
		</div>
	{/if}
	<div class="flex items-center gap-1.5">
		<DropdownMenu.Root>
			<DropdownMenu.Trigger class={buttonVariants({ variant: 'ghost', size: 'icon' })}>
				<Plus class="h-5 w-5" />
			</DropdownMenu.Trigger>
			<DropdownMenu.Content class="w-56" align="start">
				<DropdownMenu.Group>
					<DropdownMenu.Item>
						<Image />
						Image
					</DropdownMenu.Item>
					<DropdownMenu.Item>
						<Video />
						Video
					</DropdownMenu.Item>
					<DropdownMenu.Item>
						<FileText />
						File
					</DropdownMenu.Item>
					<DropdownMenu.Item>
						<MapPinXInside />
						Location
					</DropdownMenu.Item>
				</DropdownMenu.Group>
			</DropdownMenu.Content>
		</DropdownMenu.Root>
		<button
			class="hidden h-9 w-9 shrink-0 place-items-center rounded-full text-muted-foreground hover:bg-accent sm:grid"
			title="Attach file"
		>
			<Paperclip class="h-5 w-5" />
		</button>

		<div
			class="flex flex-1 items-center gap-2 rounded-full border border-border bg-secondary px-4 py-2"
		>
			<textarea
				bind:value
				onkeydown={handleKeydown}
				rows="1"
				placeholder="Type your message here..."
				class="flex-1 rounded-md bg-transparent text-sm text-foreground placeholder:text-muted-foreground focus:outline-none"
			></textarea>
			<button
				class="text-muted-foreground hover:text-foreground"
				title="Emoji"
				onclick={() => (showEmojiPicker = !showEmojiPicker)}
			>
				<Smile class="h-5 w-5" />
			</button>
		</div>

		<button
			onclick={handleSend}
			class="grid h-10 w-10 shrink-0 place-items-center rounded-full bg-primary text-primary-foreground transition-transform hover:scale-105 active:scale-95"
			title="Send"
		>
			<Send class="h-4.5 w-4.5" />
		</button>
	</div>
</div>
