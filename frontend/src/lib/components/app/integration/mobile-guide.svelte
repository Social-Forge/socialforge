<script lang="ts" module>
	type ChannelType = 'whatsapp_waha' | 'whatsapp_meta' | 'messenger' | 'instagram' | 'telegram';
</script>

<script lang="ts">
	import * as Sheet from '$lib/components/ui/sheet/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { buttonVariants } from '$lib/components/ui/button/index.js';
	import * as Accordion from '$lib/components/ui/accordion/index.js';

	let { channel, type }: { channel: string; type: string } = $props();

	const ChannelMap = {
		whatsapp_waha: 'WhatsApp',
		whatsapp_meta: 'WhatsApp API',
		messenger: 'Messenger',
		instagram: 'Instagram',
		telegram: 'Telegram'
	};

	let currentChannelName = $derived(() => ChannelMap[type as keyof typeof ChannelMap] || type);

	const GuidesItems = {
		linkchat: {
			title: 'Link Chat',
			guides: ["Click Add Chat Link'", 'Enter the chat link name', 'Click Save', 'Done']
		},
		webchat: {
			title: 'Web Chat',
			guides: [
				'Click Add Web Chat',
				'Enter the Web Chat name and website URL',
				'Click Save',
				'Done.'
			]
		},
		whatsapp_waha: {
			title: 'WhatsApp',
			guides: [
				'Click Add WhatsApp',
				'Enter your name and WhatsApp number (e.g. +6281234567890)',
				'Click Next',
				'Tap the three dots on your phone',
				'Tap Linked Devices',
				'Scan the QR code that appears',
				'Done.'
			]
		},
		whatsapp_meta: {
			title: 'WhatsApp API',
			guides: [
				'Click Add WhatsApp',
				'Select Business',
				'Select or Create WhatsApp Account',
				'Select WhatsApp Number',
				'Add payment method',
				'Done'
			]
		},
		messenger: {
			title: 'Messenger',
			guides: ['Click Add Messenger', 'Approve permissions from Facebook', 'Select a Page', 'Done']
		},
		instagram: {
			title: 'Instagram',
			guides: [
				'Click Add Instagram',
				'Approve permissions from Facebook',
				'Select Instagram',
				'Done'
			]
		},
		telegram: {
			title: 'Telegram',
			guides: [
				'Log in to the Telegram app',
				'Search for the user @BotFather',
				'Type /newbot',
				'Enter the name and username',
				'Copy the Access Token',
				'Paste the Access Token into the Telegram integration in Socialchat',
				'DDone'
			]
		}
	};

	const currentGuide = $derived(GuidesItems[type as keyof typeof GuidesItems]);
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
					class="h-auto w-full overflow-y-auto rounded-xl border-t border-border  px-5 pt-5 pb-20"
				>
					<Accordion.Trigger>Bagaimana cara menambahkan {currentChannelName()}?</Accordion.Trigger>
					<Accordion.Content>
						<div class="space-y-4">
							{#each currentGuide.guides as guide, index (index)}
								<div class="flex w-full items-center gap-2">
									<div
										class="mt-1 flex h-4 w-4 shrink-0 items-center justify-center rounded-full bg-primary select-none"
									>
										<div class="text-[11px] text-white">{index + 1}</div>
									</div>
									<div class="text-[11px]">{guide}</div>
								</div>
							{/each}
						</div>
					</Accordion.Content>
				</Accordion.Item>
			</Accordion.Root>
		</div>
	</Sheet.Content>
</Sheet.Root>
