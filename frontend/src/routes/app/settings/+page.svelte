<script lang="ts">
	import { SettingLayout } from '$lib/components/app/setting/index.js';

	let { data } = $props();

	let selectedTab = $state<string>('');
</script>

<SettingLayout>
	<div class="flex w-full flex-col gap-8">
		<div class="rounded-3xl bg-muted p-0 shadow-md lg:p-7">
			<!-- {#if selectedMenu?.child && selectedMenu.child?.length > 0}
								<div class="relative w-full max-w-114.5">
									<UnderlineTabs.Root value={selectedTab}>
										<UnderlineTabs.List>
											{#each selectedMenu.child as item, index (index)}
												<UnderlineTabs.Trigger value={item.title}>
													{item.title}
												</UnderlineTabs.Trigger>
											{/each}
										</UnderlineTabs.List>
										<UnderlineTabs.Content value={selectedTab}></UnderlineTabs.Content>
									</UnderlineTabs.Root>
								</div>
							{/if} -->
		</div>
	</div>
</SettingLayout>
