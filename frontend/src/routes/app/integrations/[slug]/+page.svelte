<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Separator } from '$lib/components/ui/separator/index.js';
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
							<div
								class="flex w-full cursor-pointer justify-end font-medium text-primary lg:hidden"
							>
								Tutorial
							</div>
						</div>
						<div class="mb-5 text-sm text-muted-foreground">
							<span>
								Mulailah berkomunikasi dengan pelanggan Anda melalui {currentChannelName()}
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
								<span class="ml-2"> Tambahkan {currentChannelName()} </span>
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
						<div class="mb-3 text-xl font-medium">Integrasi Channel Link Chat</div>
						<div class="text-muted-foreground"></div>
					</div>
					<div
						class="h-auto w-full overflow-y-auto rounded-xl border-t border-border px-5 pt-5 pb-10"
					>
						<div class="mb-6 flex justify-between gap-2 font-bold">
							Bagaimana cara menambahkan Link Chat? <span class="cursor-pointer"
								><svg
									xmlns="http://www.w3.org/2000/svg"
									xmlns:xlink="http://www.w3.org/1999/xlink"
									aria-hidden="true"
									role="img"
									class="iconify iconify--ph text-2xl"
									width="1em"
									height="1em"
									viewBox="0 0 256 256"
									><path
										fill="currentColor"
										d="M210.83 162.83a4 4 0 0 1-5.66 0L128 85.66l-77.17 77.17a4 4 0 0 1-5.66-5.66l80-80a4 4 0 0 1 5.66 0l80 80a4 4 0 0 1 0 5.66"
									></path></svg
								></span
							>
						</div>
						<div class="space-y-4">
							<div class="flex w-full items-center gap-4">
								<div
									class="mt-1 flex h-5 w-5 shrink-0 items-center justify-center rounded-full bg-primary select-none"
								>
									<div class="text-white">1</div>
								</div>
								<div class="text-base-content">Klik Tambahkan Link Chat</div>
							</div>
							<div class="flex w-full items-center gap-4">
								<div
									class="mt-1 flex h-5 w-5 shrink-0 items-center justify-center rounded-full bg-primary select-none"
								>
									<div class="text-white">2</div>
								</div>
								<div class="text-base-content">Masukan Nama Link Chat</div>
							</div>
							<div class="flex w-full items-center gap-4">
								<div
									class="mt-1 flex h-5 w-5 shrink-0 items-center justify-center rounded-full bg-primary select-none"
								>
									<div class="text-white">3</div>
								</div>
								<div class="text-base-content">Klik Simpan</div>
							</div>
							<div class="flex w-full items-center gap-4">
								<div
									class="mt-1 flex h-5 w-5 shrink-0 items-center justify-center rounded-full bg-primary select-none"
								>
									<div class="text-white">4</div>
								</div>
								<div class="text-base-content">Selesai.</div>
							</div>
						</div>
					</div>
				</div>
			</div>
		</div>
	</div>
</div>
