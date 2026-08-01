<script lang="ts">
	import { AppAvatar } from '$lib/components/app';
	import { cn } from '$lib/utils';
	import { formatBytes } from '$lib/utils/formatter';
	import { formatTime } from '$lib/utils/time';
	import { AppDropdownMenu } from '$lib/components/app';
	import {
		Check,
		CheckCheck,
		Clock,
		AlertCircle,
		MoreVertical,
		FileText,
		Play,
		Download,
		ExternalLink
	} from '@lucide/svelte';
	import type { ChatMessage } from './types';

	let { message, showSender = false }: { message: ChatMessage; showSender?: boolean } = $props();

	const isOut = $derived(message.direction === 'out');

	const menuItems = [
		{ label: 'Reply', onSelect: () => {} },
		{ label: 'Forward', onSelect: () => {} },
		{ label: 'Copy', onSelect: () => {} },
		{ label: 'Delete', onSelect: () => {}, destructive: true }
	];
</script>

{#snippet statusIcon()}
	{#if message.status === 'sending'}
		<Clock class="h-3.5 w-3.5 text-bubble-out-foreground/60" />
	{:else if message.status === 'sent'}
		<Check class="h-3.5 w-3.5 text-bubble-out-foreground/60" />
	{:else if message.status === 'delivered'}
		<CheckCheck class="h-3.5 w-3.5 text-bubble-out-foreground/60" />
	{:else if message.status === 'read'}
		<CheckCheck class="h-3.5 w-3.5 text-sky-300" />
	{:else if message.status === 'failed'}
		<AlertCircle class="h-3.5 w-3.5 text-destructive" />
	{/if}
{/snippet}

{#snippet kebabTrigger()}
	<span
		class="grid h-6 w-6 place-items-center rounded-full opacity-0 transition-opacity group-hover:opacity-100 hover:bg-black/10 dark:hover:bg-white/10"
	>
		<MoreVertical class="h-4 w-4" />
	</span>
{/snippet}

<div class={cn('group flex items-end gap-2 px-2', isOut ? 'flex-row-reverse' : 'flex-row')}>
	{#if !isOut}
		<AppAvatar
			src={message.senderAvatarUrl}
			fallback={message.senderName?.slice(0, 1) ?? '?'}
			size="sm"
			class="mb-0.5"
		/>
	{/if}

	<div class={cn('flex max-w-[75%] flex-col md:max-w-[60%]', isOut ? 'items-end' : 'items-start')}>
		{#if showSender && !isOut && message.senderName}
			<span class="mb-1 px-1 text-xs font-semibold text-primary">{message.senderName}</span>
		{/if}

		<div
			class={cn(
				'relative rounded-2xl px-3 py-2 shadow-sm',
				isOut
					? 'rounded-br-sm bg-bubble-out text-bubble-out-foreground'
					: 'rounded-bl-sm bg-bubble-in text-bubble-in-foreground'
			)}
		>
			<div class="absolute -top-1 {isOut ? '-left-7' : '-right-7'}">
				<AppDropdownMenu items={menuItems} align={isOut ? 'end' : 'start'} trigger={kebabTrigger} />
			</div>

			<!-- Template message -->
			{#if message.template}
				<div class="w-64 sm:w-72">
					{#if message.template.headerText}
						<div class="mb-1 text-sm font-bold">{message.template.headerText}</div>
					{/if}
					<p class="text-sm leading-relaxed">{message.template.bodyText}</p>
					{#if message.template.footerText}
						<p class="mt-1 text-xs opacity-70">{message.template.footerText}</p>
					{/if}
					{#if message.template.buttons?.length}
						<div class="mt-2 flex flex-col gap-1 border-t border-current/15 pt-2">
							{#each message.template.buttons as btn, bi (bi)}
								<button
									class="rounded-md py-1.5 text-center text-sm font-medium text-primary hover:bg-black/5 dark:hover:bg-white/5"
								>
									{btn.label}
								</button>
							{/each}
						</div>
					{/if}
				</div>

				<!-- Media (image / video / document / audio), optionally with caption -->
			{:else if message.media?.length}
				<div class="flex flex-col gap-2">
					{#each message.media as item, mi (mi)}
						{#if item.kind === 'image'}
							<img
								src={item.url}
								alt={item.name ?? 'photo'}
								class="max-h-72 w-full rounded-lg object-cover sm:w-72"
							/>
						{:else if item.kind === 'video'}
							<div class="relative w-full overflow-hidden rounded-lg sm:w-72">
								<img
									src={item.thumbnailUrl}
									alt={item.name ?? 'video'}
									class="max-h-72 w-full object-cover"
								/>
								<div class="absolute inset-0 grid place-items-center bg-black/25">
									<span
										class="grid h-11 w-11 place-items-center rounded-full bg-white/90 text-foreground"
									>
										<Play class="h-5 w-5 fill-current" />
									</span>
								</div>
								{#if item.durationSeconds}
									<span
										class="absolute right-1.5 bottom-1.5 rounded bg-black/60 px-1.5 py-0.5 text-[11px] text-white"
									>
										0:{item.durationSeconds.toString().padStart(2, '0')}
									</span>
								{/if}
							</div>
						{:else if item.kind === 'document'}
							<div
								class="flex w-64 items-center gap-3 rounded-lg bg-black/5 p-2.5 sm:w-72 dark:bg-white/5"
							>
								<span
									class="grid h-10 w-10 shrink-0 place-items-center rounded-md bg-primary/15 text-primary"
								>
									<FileText
										class={cn('h-5 w-5 ', isOut ? 'text-white' : 'text-muted-foreground')}
									/>
								</span>
								<div class="min-w-0 flex-1">
									<div class="truncate text-sm font-medium">{item.name}</div>
									<div class="text-xs opacity-70">
										{item.sizeBytes ? formatBytes(item.sizeBytes) : ''}
										{item.mimeType ? `· ${item.mimeType}` : ''}
									</div>
								</div>
								<Download class="h-4 w-4 shrink-0 opacity-70" />
							</div>
						{/if}
					{/each}
					{#if message.text}
						<p class="text-sm leading-relaxed">{message.text}</p>
					{/if}
				</div>

				<!-- Link preview, optionally with leading text -->
			{:else if message.linkPreview}
				<div class="flex w-64 flex-col gap-2 sm:w-72">
					{#if message.text}
						<p class="text-sm leading-relaxed">{message.text}</p>
					{/if}
					<a
						href={message.linkPreview.url}
						class="block overflow-hidden rounded-lg bg-black/5 hover:bg-black/10 dark:bg-white/5 dark:hover:bg-white/10"
					>
						{#if message.linkPreview.imageUrl}
							<img src={message.linkPreview.imageUrl} alt="" class="h-32 w-full object-cover" />
						{/if}
						<div class="p-2.5">
							<div class="flex items-center gap-1 text-[11px] tracking-wide uppercase opacity-60">
								{message.linkPreview.siteName}
								<ExternalLink class="h-3 w-3" />
							</div>
							<div class="mt-0.5 text-sm font-semibold">{message.linkPreview.title}</div>
							{#if message.linkPreview.description}
								<div class="mt-0.5 line-clamp-2 text-xs opacity-70">
									{message.linkPreview.description}
								</div>
							{/if}
						</div>
					</a>
				</div>

				<!-- Plain text -->
			{:else if message.text}
				<p class="text-sm leading-relaxed">{message.text}</p>
			{/if}

			<div
				class={cn(
					'mt-1 flex items-center justify-end gap-1 text-[11px]',
					isOut ? 'text-bubble-out-foreground/70' : 'text-bubble-in-foreground/60'
				)}
			>
				{formatTime(message.timestamp)}
				{#if isOut}
					{@render statusIcon()}
				{/if}
			</div>
		</div>
	</div>
</div>
