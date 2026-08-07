<script lang="ts" module>
	type QueryOption = {
		page: number;
		limit: number;
		search: string;
		channel?: string;
		label?: string;
		assign_agent?: string;
		status: string;
		sort_by: string;
		order_by: 'asc' | 'desc';
		date_from: string;
		date_to: string;
	};
</script>

<script lang="ts">
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import { Button, buttonVariants } from '$lib/components/ui/button/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import * as Select from '$lib/components/ui/select/index.js';
	import { SlidersHorizontal } from '@lucide/svelte';

	let {
		meta,
		updateQuery
	}: {
		meta?: ApiMeta;
		updateQuery: (updates: Partial<QueryOption>, resetPage: boolean) => Promise<void>;
	} = $props();

	let channelValue = $state('all');
	let agentValue = $state('all');
	let assignAgentId = $state('all');
	const channelItems = [
		{
			label: 'Whatsapp Waha',
			value: 'whatsapp_waha'
		},
		{
			label: 'Whatsapp API',
			value: 'whatsapp_meta'
		},
		{
			label: 'Messenger',
			value: 'messenger'
		},
		{
			label: 'Instagram',
			value: 'instagram'
		},
		{
			label: 'Telegram',
			value: 'telegram'
		}
	];
</script>

<Dialog.Root>
	<Dialog.Trigger type="button" class={buttonVariants({ variant: 'outline', size: 'icon' })}>
		<SlidersHorizontal />
	</Dialog.Trigger>
	<Dialog.Content class="sm:max-w-[425px]">
		<Dialog.Header>
			<Dialog.Title>Contact Filter</Dialog.Title>
			<Dialog.Description>
				Please fill in some field values to search for contact data.
			</Dialog.Description>
		</Dialog.Header>
		<div class="grid gap-4">
			<div class="grid gap-3">
				<Label for="channel">Channel Name</Label>
				<Select.Root type="single" bind:value={channelValue} name="channel">
					<Select.Trigger class="w-full capitalize">
						{channelValue || 'Select channel'}
					</Select.Trigger>
					<Select.Content class="w-full">
						{#each channelItems as ch, index (index)}
							<Select.Item value={ch.value}>{ch.label}</Select.Item>
						{/each}
					</Select.Content>
				</Select.Root>
			</div>
			<div class="grid gap-3">
				<Label for="label">Label Name</Label>
				<Select.Root type="single" bind:value={agentValue} name="label">
					<Select.Trigger class="w-full capitalize">
						{agentValue || 'Select label'}
					</Select.Trigger>
					<Select.Content class="w-full">
						<Select.Item value="new_customer">Pelanggan Baru</Select.Item>
					</Select.Content>
				</Select.Root>
			</div>
			<div class="grid gap-3">
				<Label for="assign_agent">Agent Name</Label>
				<Select.Root type="single" bind:value={assignAgentId} name="assign_agent">
					<Select.Trigger class="w-full capitalize">
						{assignAgentId || 'Select Agent'}
					</Select.Trigger>
					<Select.Content class="w-full">
						<Select.Item value="agent_id">Sarah</Select.Item>
					</Select.Content>
				</Select.Root>
			</div>
		</div>
		<Dialog.Footer>
			<Dialog.Close type="button" class={buttonVariants({ variant: 'outline' })}>
				Cancel
			</Dialog.Close>
			<Button type="button">Save changes</Button>
		</Dialog.Footer>
	</Dialog.Content>
</Dialog.Root>
