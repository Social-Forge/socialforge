<script lang="ts" module>
	interface SettingMenu {
		title: string;
		href: string;
		desc: string;
		child?: SettingMenu[];
	}
</script>

<script lang="ts">
	import { localizeHref } from '$lib/paraglide/runtime.js';
	import { ChevronRight, ChevronLeft } from '@lucide/svelte';
	import * as Sheet from '$lib/components/ui/sheet/index.js';
	import { buttonVariants } from '$lib/components/ui/button/index.js';

	let {
		menus,
		onMenuChange
	}: {
		menus: SettingMenu[];
		onMenuChange?: (menu: SettingMenu) => void;
	} = $props();

	let selectedMenu = $state<SettingMenu | null>(null);

	function handleMenuChange(menu: SettingMenu) {
		selectedMenu = menu;
		onMenuChange?.(menu);
	}
</script>

<Sheet.Root>
	<Sheet.Trigger class={buttonVariants({ variant: 'ghost', class: 'mt-6' })}>
		<ChevronLeft />
		{selectedMenu?.title}
	</Sheet.Trigger>
	<Sheet.Content side="left">
		<Sheet.Header>
			<Sheet.Title class="sr-only">{selectedMenu?.title}</Sheet.Title>
			<Sheet.Description class="sr-only">{selectedMenu?.desc}</Sheet.Description>
		</Sheet.Header>

		<div class="space-y-2 px-4">
			<div class="mb-4 text-lg font-medium text-primary">Settings</div>
			<div class="space-y-2">
				{#each menus as item, index (index)}
					<a
						href={localizeHref(item.href)}
						class="-ml-2.5 flex w-[calc(100%+20px)] cursor-pointer items-center justify-start gap-3 rounded-2xl px-2.5 py-3 text-sm hover:bg-background active:bg-background"
						onclick={() => handleMenuChange(item)}
					>
						<div class="w-full select-none">
							<div class="text-sm font-semibold">
								<div class="flex items-center justify-between gap-2">
									<div>{item.title}</div>
									{#if item.child && item.child?.length > 0}
										<ChevronRight class="size-4" color="currentColor" />
									{/if}
								</div>
							</div>
						</div>
					</a>
				{/each}
			</div>
		</div>
	</Sheet.Content>
</Sheet.Root>
