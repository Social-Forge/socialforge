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

	let tab = $state(page.url.searchParams.get('tab') ?? 'general');

	async function onTabChange(tab: string) {
		tab = tab;
		await goto(localizeHref(`/app/settings/contact?tab=${tab}`));
	}
</script>

<SettingLayout>
	<div class="flex w-full flex-col gap-8">
		<div class="rounded-md bg-card p-2 shadow-md lg:p-7">
			<UnderlineTabs.Root bind:value={tab} onValueChange={(tab) => onTabChange(tab)}>
				<UnderlineTabs.List>
					<UnderlineTabs.Trigger value="general">General</UnderlineTabs.Trigger>
					<UnderlineTabs.Trigger value="column">Column</UnderlineTabs.Trigger>
					<UnderlineTabs.Trigger value="blocked">Blocked</UnderlineTabs.Trigger>
				</UnderlineTabs.List>
				<UnderlineTabs.Content value="general">General</UnderlineTabs.Content>
				<UnderlineTabs.Content value="column">Column</UnderlineTabs.Content>
				<UnderlineTabs.Content value="blocked">Blocked</UnderlineTabs.Content>
			</UnderlineTabs.Root>
		</div>
	</div>
</SettingLayout>
