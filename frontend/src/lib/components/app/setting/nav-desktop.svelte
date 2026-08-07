<script lang="ts" module>
	interface SettingMenu {
		title: string;
		href: string;
		desc: string;
		child?: SettingMenu[];
	}
</script>

<script lang="ts">
	import { page } from '$app/state';
	import { localizeHref } from '$lib/paraglide/runtime.js';
	import { ChevronRight } from '@lucide/svelte';
	import { cn } from '$lib/utils';

	let {
		menus,
		onMenuChange
	}: {
		menus: SettingMenu[];
		onMenuChange?: (menu: SettingMenu) => void;
	} = $props();
</script>

<div class="hidden h-full w-[20%] border-r border-border bg-muted px-3 py-7.5 lg:block lg:px-10">
	<div class="mb-4 text-lg font-medium text-primary">Settings</div>
	<div class="space-y-2">
		{#each menus as item, index (index)}
			<a
				href={localizeHref(item.href)}
				class={cn(
					'-ml-2.5 flex w-[calc(100%+20px)] cursor-pointer items-center justify-start gap-3 rounded-2xl px-2.5 py-3 text-sm hover:bg-background active:bg-background',
					{
						'bg-background': page.url.pathname.startsWith(item.href)
					}
				)}
				onclick={() => {
					onMenuChange?.(item);
				}}
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
