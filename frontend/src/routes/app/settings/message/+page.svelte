<script lang="ts">
	import { SettingLayout } from '$lib/components/app/setting/index.js';
	import * as Card from '$lib/components/ui/card';

	let { data } = $props();
	const user = $derived(data.user as UserResponse | undefined);
</script>

<SettingLayout>
	<div class="flex w-full flex-col gap-8">
		<div class="rounded-md bg-card p-2 shadow-md lg:p-7">
			<h2 class="mb-3 text-lg font-semibold">Message Settings</h2>
		</div>
	</div>
</SettingLayout>
