<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Separator } from '$lib/components/ui/separator/index.js';
	import { MobileGuide, DesktopGuide } from '$lib/components/app/integration/index.js';
	import { localizeHref } from '$lib/paraglide/runtime.js';
	import { Plus, Crown, User, Rocket, ArrowLeft } from '@lucide/svelte';

	let { data } = $props();

	const ChannelMap = {
		whatsapp_waha: 'WhatsApp',
		whatsapp_meta: 'WhatsApp API',
		messenger: 'Messenger',
		instagram: 'Instagram',
		telegram: 'Telegram'
	};

	let currentChannelName = $derived(
		() => ChannelMap[data.slug as keyof typeof ChannelMap] || data.slug
	);

	onMount(() => {
		// console.log('data', data);
	});
</script>

<div class="h-[calc(100dvh-70px)] w-full overflow-y-auto bg-background lg:h-full">
	<div class="relative flex h-full w-full">
		<div class="flex h-full w-full">
			<div class="flex h-full w-full flex-col border-r border-border px-10 pt-10 lg:w-[70%]">
				<div class="mb-10 flex h-max w-full flex-col lg:flex-row">
					<div class="flex w-full flex-col">
						<div class="mb-3 flex justify-between whitespace-nowrap">
							<div class="flex items-center gap-3">
								<Button
									variant="ghost"
									size="icon"
									onclick={() => goto(localizeHref('/app/integrations'))}
								>
									<ArrowLeft />
								</Button>
								<div class="text-xl font-medium capitalize">{currentChannelName()}</div>
							</div>
							<MobileGuide channel={currentChannelName()} type={data.slug} />
						</div>
						<div class="mb-5 text-sm text-muted-foreground">
							<span>
								Start communicating with your customers through {currentChannelName()}
							</span>
						</div>
						<div data-test="limit" class="mb-2 flex items-center justify-start">
							<Crown class="mr-2 mb-1 size-5 text-amber-600" />
							0
						</div>
					</div>
					<div class="flex w-full flex-row justify-end gap-4">
						<div class="flex w-max gap-4">
							<Button class="capitalize">
								<Plus class="text-3xl text-white" />
								<span class="ml-2"> Add {currentChannelName()} </span>
							</Button>
						</div>
					</div>
				</div>
				<Separator class="my-5" />
				<div class="scrollbar-primary h-full w-full space-y-5 overflow-y-auto pb-10">
					<div>
						<div class="flex flex-col items-center justify-center">
							<span class="mt-10 mb-2 text-xl font-semibold">Belum Ada {currentChannelName()}</span
							><span class="text-center text-sm text-muted-foreground">
								Silahkan integrasikan akun anda dengan menekan tombol 'Tambahkan {currentChannelName()}'
							</span>
						</div>
					</div>
				</div>
			</div>
			<DesktopGuide channel={currentChannelName()} type={data.slug} />
		</div>
	</div>
</div>
