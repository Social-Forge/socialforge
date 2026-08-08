<script lang="ts">
	import { page } from '$app/state';
	import { goto } from '$app/navigation';
	import { invalidateAll } from '$app/navigation';
	import * as UnderlineTabs from '$lib/components/extras/underline-tabs';
	import { SettingLayout } from '$lib/components/app/setting/index.js';
	import * as Card from '$lib/components/ui/card';
	import { localizeHref } from '$lib/paraglide/runtime.js';

	let { data } = $props();
	const user = $derived(data.user as UserResponse | undefined);

	let tab = $state(page.url.searchParams.get('tab') ?? 'agent-rotator');

	async function onTabChange(tab: string) {
		tab = tab;
		await goto(localizeHref(`/app/settings/agent?tab=${tab}`));
	}
</script>

<SettingLayout>
	<div class="flex w-full flex-col gap-8">
		<div class="rounded-md bg-card p-2 shadow-md lg:p-7">
			<UnderlineTabs.Root bind:value={tab} onValueChange={(tab) => onTabChange(tab)}>
				<UnderlineTabs.List>
					<UnderlineTabs.Trigger value="agent-rotator">Agent Rotator</UnderlineTabs.Trigger>
					<UnderlineTabs.Trigger value="additional-agent">Additional Agent</UnderlineTabs.Trigger>
					<UnderlineTabs.Trigger value="switch-agent">Switch Agent</UnderlineTabs.Trigger>
					<UnderlineTabs.Trigger value="sticky-agent">Sticky Agent</UnderlineTabs.Trigger>
				</UnderlineTabs.List>
				<UnderlineTabs.Content value="agent-rotator">Agent Rotator</UnderlineTabs.Content>
				<UnderlineTabs.Content value="additional-agent">Additional Agent</UnderlineTabs.Content>
				<UnderlineTabs.Content value="switch-agent">Switch Agent</UnderlineTabs.Content>
				<UnderlineTabs.Content value="sticky-agent">Sticky Agent</UnderlineTabs.Content>
			</UnderlineTabs.Root>
		</div>
	</div>
</SettingLayout>
