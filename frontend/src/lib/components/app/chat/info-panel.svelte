<script lang="ts">
	import { AppAvatar } from '$lib/components/app';
	import { activeContact, documents, lastMedia } from './data';
	import { chatUiState } from '$lib/hooks/chat-ui.svelte';
	import { formatBytes } from '$lib/utils/formatter';
	import { X, Info, Image, FileText, UserPlus, Heart, Ban, MoreVertical } from '@lucide/svelte';
</script>

<aside
	class="flex h-full w-full flex-col border-l border-border bg-background md:w-[320px] lg:w-85"
>
	<div class="flex items-center justify-between px-4 pt-4">
		<h2 class="text-sm font-semibold text-foreground">Profile Details</h2>
		<button
			class="grid h-8 w-8 place-items-center rounded-md text-muted-foreground hover:bg-accent"
			onclick={() => chatUiState.toggleInfo()}
		>
			<X class="h-4 w-4" />
		</button>
	</div>

	<div class="scroll-thin flex-1 overflow-y-auto px-4 pb-6">
		<div class="flex flex-col items-center gap-3 py-6">
			<AppAvatar
				src={activeContact.avatarUrl}
				fallback={activeContact.name.slice(0, 1)}
				size="xl"
			/>
			<div class="text-center">
				<div class="text-base font-semibold text-foreground">{activeContact.name}</div>
				<div class="text-sm text-muted-foreground">{activeContact.location}</div>
			</div>
			<div class="flex items-center gap-3">
				<button
					class="grid h-9 w-9 place-items-center rounded-full border border-border text-muted-foreground hover:bg-accent"
					title="Add contact"
				>
					<UserPlus class="h-4 w-4" />
				</button>
				<button
					class="grid h-9 w-9 place-items-center rounded-full bg-destructive/10 text-destructive hover:bg-destructive/15"
					title="Favorite"
				>
					<Heart class="h-4 w-4" />
				</button>
				<button
					class="grid h-9 w-9 place-items-center rounded-full border border-border text-muted-foreground hover:bg-accent"
					title="Block"
				>
					<Ban class="h-4 w-4" />
				</button>
			</div>
		</div>

		<section class="border-t border-border pt-4">
			<div class="mb-3 flex items-center justify-between">
				<h3 class="text-sm font-semibold text-foreground">User Information</h3>
				<Info class="h-4 w-4 text-muted-foreground" />
			</div>
			<div class="flex flex-col gap-3 text-sm">
				<div>
					<div class="text-xs text-muted-foreground">Phone</div>
					<div class="text-foreground">{activeContact.phone}</div>
				</div>
				<div>
					<div class="text-xs text-muted-foreground">Email</div>
					<div class="break-all text-foreground">{activeContact.email}</div>
				</div>
				<div>
					<div class="text-xs text-muted-foreground">Address</div>
					<div class="text-foreground">{activeContact.address}</div>
				</div>
			</div>
		</section>

		<section class="mt-4 border-t border-border pt-4">
			<div class="mb-3 flex items-center justify-between">
				<h3 class="text-sm font-semibold text-foreground">Last Media</h3>
				<Image class="h-4 w-4 text-muted-foreground" />
			</div>
			<div class="grid grid-cols-3 gap-2">
				{#each lastMedia as src, i (i)}
					<div class="aspect-square overflow-hidden rounded-md bg-muted">
						<img {src} alt="" class="h-full w-full object-cover" />
					</div>
				{/each}
			</div>
		</section>

		<section class="mt-4 border-t border-border pt-4">
			<div class="mb-3 flex items-center justify-between">
				<h3 class="text-sm font-semibold text-foreground">Document</h3>
				<FileText class="h-4 w-4 text-muted-foreground" />
			</div>
			<div class="flex flex-col gap-1">
				{#each documents as doc, di (di)}
					<div class="flex items-center gap-3 rounded-lg px-1.5 py-2 hover:bg-accent">
						<span
							class="grid h-9 w-9 shrink-0 place-items-center rounded-md bg-primary/15 text-primary"
						>
							<FileText class="h-4.5 w-4.5" />
						</span>
						<div class="min-w-0 flex-1">
							<div class="truncate text-sm font-medium text-foreground">{doc.name}</div>
							<div class="text-xs text-muted-foreground">
								{formatBytes(doc.sizeBytes)} · {doc.kind}
							</div>
						</div>
						<button
							class="grid h-7 w-7 shrink-0 place-items-center rounded-md text-muted-foreground hover:bg-accent"
						>
							<MoreVertical class="h-4 w-4" />
						</button>
					</div>
				{/each}
			</div>
		</section>
	</div>
</aside>
