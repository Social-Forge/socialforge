<script lang="ts">
	import * as Accordion from '$lib/components/ui/accordion/index.js';

	let { channel, type }: { channel: string; type: string } = $props();

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

<div
	class="bg-opacity-40 fixed top-0 left-0 z-60 hidden h-full w-full bg-black lg:relative lg:z-0 lg:flex lg:h-full lg:w-[30%] lg:bg-transparent"
>
	<div
		class="bg-base-100 absolute bottom-0 left-0 z-10 flex h-122.5 w-full flex-col gap-6 rounded-t-3xl px-10 pt-5 lg:h-full lg:bg-transparent"
	>
		<div
			class="absolute top-2 right-0 flex h-15 w-15 cursor-pointer items-center justify-center lg:hidden"
		>
			<i class="bi bi-x-lg text-[20px]"></i>
		</div>
		<div class="w-full">
			<div class="mb-3 text-xl font-medium capitalize">
				Integrasi Channel {channel}
			</div>
			<div class="text-muted-foreground"></div>
		</div>
		<div class="h-auto w-full overflow-y-auto rounded-xl border-t border-border px-5 pt-5 pb-10">
			<Accordion.Root type="single">
				<Accordion.Item value="item-1">
					<Accordion.Trigger>
						How to add a {channel} ?
					</Accordion.Trigger>
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
	</div>
</div>
