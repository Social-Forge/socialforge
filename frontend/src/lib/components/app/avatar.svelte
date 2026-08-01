<script lang="ts">
	import { cn } from '$lib/utils';
	import {
		WhatsAppIcon,
		WhatsAppBusinessIcon,
		MessengerIcon,
		InstagramIcon,
		TelegramIcon,
		GithubIcon
	} from '$lib/components/icons';

	let {
		src,
		alt = '',
		fallback = '',
		size = 'md',
		channel,
		class: className = ''
	}: {
		src?: string;
		alt?: string;
		fallback?: string;
		size?: 'sm' | 'md' | 'lg' | 'xl';
		channel: 'whatsapp_waha' | 'whatsapp_meta' | 'messenger' | 'instagram' | 'telegram';
		class?: string;
	} = $props();

	const sizes = {
		sm: 'h-8 w-8 text-xs',
		md: 'h-10 w-10 text-sm',
		lg: 'h-14 w-14 text-base',
		xl: 'h-24 w-24 text-2xl'
	} as const;

	const dotSizes = {
		sm: 'h-2 w-2',
		md: 'h-3 w-3',
		lg: 'h-4 w-4',
		xl: 'h-5 w-5'
	} as const;

	const ChannelMap = {
		whatsapp_waha: WhatsAppIcon,
		whatsapp_meta: WhatsAppBusinessIcon,
		messenger: MessengerIcon,
		instagram: InstagramIcon,
		telegram: TelegramIcon
	};
</script>

<span class={cn('relative inline-flex shrink-0', className)}>
	<span
		class={cn(
			'flex items-center justify-center overflow-hidden rounded-full bg-muted font-medium text-muted-foreground',
			sizes[size]
		)}
	>
		{#if src}
			<img {src} {alt} class="h-full w-full object-cover" />
		{:else}
			{fallback}
		{/if}
	</span>
	{#if channel}
		{@const Icon = ChannelMap[channel]}
		<Icon
			class={cn('absolute right-0 bottom-0 rounded-full ring-2 ring-background', dotSizes[size])}
		></Icon>
	{/if}
</span>
