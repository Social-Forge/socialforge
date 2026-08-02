<script lang="ts" module>
	type ChannelType = 'whatsapp_waha' | 'whatsapp_meta' | 'messenger' | 'instagram' | 'telegram';
</script>

<script lang="ts">
	import * as Sheet from '$lib/components/ui/sheet/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { buttonVariants } from '$lib/components/ui/button/index.js';
	import * as Accordion from '$lib/components/ui/accordion/index.js';

	let {
		type,
		contents = []
	}: {
		type: string;
		contents?: string[];
	} = $props();

	const ChannelMap = {
		whatsapp_waha: 'WhatsApp',
		whatsapp_meta: 'WhatsApp API',
		messenger: 'Messenger',
		instagram: 'Instagram',
		telegram: 'Telegram'
	};

	let currentChannelName = $derived(() => ChannelMap[type as keyof typeof ChannelMap] || type);
</script>

<Sheet.Root>
	<Sheet.Trigger
		class={buttonVariants({
			variant: 'ghost',
			class: 'font-medium text-primary lg:hidden'
		})}
	>
		Tutorial
	</Sheet.Trigger>
	<Sheet.Content side="bottom" class="w-full">
		<Sheet.Header>
			<Sheet.Title>Integrasi Channel {currentChannelName()}</Sheet.Title>
			<Sheet.Description></Sheet.Description>
		</Sheet.Header>
		<div class="px-4">
			<Accordion.Root type="single">
				<Accordion.Item
					class="h-auto w-full overflow-y-auto rounded-xl border-t border-border bg-background px-5 pt-5 pb-20"
				>
					<Accordion.Trigger>Bagaimana cara menambahkan {currentChannelName()}?</Accordion.Trigger>
					<Accordion.Content>
						{contents.join('\n')}
					</Accordion.Content>
				</Accordion.Item>
			</Accordion.Root>
		</div>
	</Sheet.Content>
</Sheet.Root>
