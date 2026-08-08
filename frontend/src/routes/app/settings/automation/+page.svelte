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

	let tab = $state(page.url.searchParams.get('tab') ?? 'auto-response');

	async function onTabChange(tab: string) {
		tab = tab;
		await goto(localizeHref(`/app/settings/automation?tab=${tab}`));
	}
</script>

<SettingLayout>
	<div class="flex w-full flex-col gap-8">
		<div class="rounded-md bg-card p-2 shadow-md lg:p-7">
			<UnderlineTabs.Root bind:value={tab} onValueChange={(tab) => onTabChange(tab)}>
				<UnderlineTabs.List>
					<UnderlineTabs.Trigger value="auto-response">Auto Response</UnderlineTabs.Trigger>
					<UnderlineTabs.Trigger value="pipelead">Pipelead</UnderlineTabs.Trigger>
					<UnderlineTabs.Trigger value="csat">CSAT</UnderlineTabs.Trigger>
					<UnderlineTabs.Trigger value="signature-agent">Signature Agent</UnderlineTabs.Trigger>
					<UnderlineTabs.Trigger value="message-completed">Message Completed</UnderlineTabs.Trigger>
					<UnderlineTabs.Trigger value="auto-assign-contact"
						>Auto Assign Contact</UnderlineTabs.Trigger
					>
				</UnderlineTabs.List>
				<UnderlineTabs.Content value="auto-response">Auto Response</UnderlineTabs.Content>
				<UnderlineTabs.Content value="pipelead">Pipelead</UnderlineTabs.Content>
				<UnderlineTabs.Content value="csat">CSAT</UnderlineTabs.Content>
				<UnderlineTabs.Content value="signature-agent">Signature Agent</UnderlineTabs.Content>
				<UnderlineTabs.Content value="message-completed">Message Completed</UnderlineTabs.Content>
				<UnderlineTabs.Content value="auto-assign-contact"
					>Auto Assign Contact</UnderlineTabs.Content
				>
			</UnderlineTabs.Root>
		</div>
	</div>
</SettingLayout>
