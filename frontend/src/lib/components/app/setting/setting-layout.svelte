<script lang="ts" module>
	interface SettingMenu {
		title: string;
		href: string;
		desc: string;
		child?: SettingMenu[];
	}
</script>

<script lang="ts">
	import { IsMobile } from '$lib/hooks/is-mobile.svelte';
	import { NavDesktop, NavMobile } from '$lib/components/app/setting/index.js';
	import type { Snippet } from 'svelte';

	let { children }: { children?: Snippet<[]> } = $props();

	const settingsMenu: SettingMenu[] = [
		{
			title: 'Account',
			href: '/app/settings/account?tab=profile',
			desc: 'Manage your account settings',
			child: [
				{
					title: 'Account',
					href: '/app/settings/account?tab=profile',
					desc: 'Manage your account settings'
				},
				{
					title: 'Password',
					href: '/app/settings/account?tab=password',
					desc: 'Change your password securely'
				}
			]
		},
		{
			title: 'Message',
			href: '/app/settings/message',
			desc: 'Manage your message settings'
		},
		{
			title: 'Plan & Billing',
			href: '/app/settings/billing',
			desc: 'Manage your plan and billing settings'
		},
		{
			title: 'Automation',
			href: '/app/settings/automation?tab=auto-response',
			desc: 'Manage your automation settings',
			child: [
				{
					title: 'Auto Response',
					href: '/app/settings/automation?tab=auto-response',
					desc: 'Manage your auto response settings'
				},
				{
					title: 'Pipelead',
					href: '/app/settings/automation?tab=pipelead',
					desc: 'Manage your pipelead settings'
				},
				{
					title: 'CSAT',
					href: '/app/settings/automation?tab=csat',
					desc: 'Manage your CSAT settings'
				},
				{
					title: 'Signature Agent',
					href: '/app/settings/automation?tab=signature-agent',
					desc: 'Manage your signature agent settings'
				},
				{
					title: 'Message Completed',
					href: '/app/settings/automation?tab=message-completed',
					desc: 'Manage your message completed settings'
				},
				{
					title: 'Auto Assign Contact',
					href: '/app/settings/automation?tab=auto-assign-contact',
					desc: 'Manage your auto assign contact settings'
				}
			]
		},
		{
			title: 'Agent',
			href: '/app/settings/agent?tab=agent-rotator',
			desc: 'Manage your agent settings',
			child: [
				{
					title: 'Agent Rotator',
					href: '/app/settings/agent?tab=agent-rotator',
					desc: 'Manage your agent rotator settings'
				},
				{
					title: 'Additional Agent',
					href: '/app/settings/agent?tab=additional-agent',
					desc: 'With the additional agent feature, an agent can join a conversation even if it has already been picked up by another agent.'
				},
				{
					title: 'Switch Agent',
					href: '/app/settings/agent?tab=switch-agent',
					desc: 'If this feature is enabled, agents can transfer their assignments to other agents.'
				},
				{
					title: 'Sticky Agent',
					href: '/app/settings/agent?tab=sticky-agent',
					desc: "The system will automatically assign new conversations to the agent holding the lead's data."
				}
			]
		},
		{
			title: 'Contact',
			href: '/app/settings/contact?tab=general',
			desc: 'Manage your contact settings',
			child: [
				{
					title: 'Contact',
					href: '/app/settings/contact?tab=general',
					desc: 'Manage your contact settings'
				},
				{
					title: 'Column',
					href: '/app/settings/contact?tab=column',
					desc: 'You can add custom fields to contact data. These fields can also be used in broadcasts and automated messages.'
				},
				{
					title: 'Blocked',
					href: '/app/settings/contact?tab=blocked',
					desc: 'Manage your blocked contact settings.'
				}
			]
		},
		{
			title: 'Security',
			href: '/app/settings/security',
			desc: 'Manage your security settings'
		},
		{
			title: 'SLA',
			href: '/app/settings/sla-rule',
			desc: 'Manage your SLA-rule settings'
		},
		{
			title: 'Topic',
			href: '/app/settings/topic',
			desc: 'Manage your topic settings'
		},
		{
			title: 'Notification',
			href: '/app/settings/notification',
			desc: 'Manage your notification settings'
		}
	];

	let selectedMenu = $state<SettingMenu | null>(null);
	let isMobile = $derived(new IsMobile().current);
</script>

<div class="h-[calc(100dvh-70px)] w-full overflow-y-auto bg-background lg:h-full">
	<div class="relative flex h-full w-full">
		<div class="scrollbar-primary flex h-full w-full flex-col overflow-y-auto">
			<div class="relative h-full w-full lg:flex">
				<NavDesktop menus={settingsMenu} onMenuChange={(menu) => (selectedMenu = menu)} />
				{#if isMobile}
					<div class="px-4">
						<NavMobile menus={settingsMenu} onMenuChange={(menu) => (selectedMenu = menu)} />
					</div>
				{/if}
				<div
					class="block h-full w-full overflow-y-auto border-r border-border px-4 pt-4 pb-30 lg:max-h-full lg:bg-background lg:px-8 lg:py-8 lg:pt-10"
				>
					{@render children?.()}
				</div>
			</div>
		</div>
	</div>
</div>
