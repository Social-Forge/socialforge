<script lang="ts">
	import * as Avatar from '$lib/components/ui/avatar/index.js';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu/index.js';
	import { localizeHref } from '$lib/paraglide/runtime';
	import {
		LogOutIcon,
		SparklesIcon,
		BadgeCheckIcon,
		BellIcon,
		CreditCardIcon
	} from '@lucide/svelte';

	let { user }: { user?: UserResponse | null } = $props();

	async function logout() {
		try {
			const response = await fetch('/api/user/logout', {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json'
				}
			});
			if (!response.ok) {
				throw new Error('Logout failed');
			}
			const redirectTarget = `${window.location.pathname}${window.location.search}`;
			const signInHref = localizeHref(`/signin?redirect=${encodeURIComponent(redirectTarget)}`);

			// Use a full navigation so auth pages mount fresh after tearing down the app shell.
			window.location.assign(signInHref);
		} catch (error) {
			console.error('Logout failed:', error);
		}
	}
</script>

<DropdownMenu.Root>
	<DropdownMenu.Trigger>
		{#snippet children()}
			<Avatar.Root class="size-8 rounded-lg">
				<Avatar.Image src={user?.avatar_url} alt={user?.full_name || 'User'} />
				<Avatar.Fallback class="rounded-lg"
					>{user?.full_name.slice(0, 2).toUpperCase() || 'CN'}</Avatar.Fallback
				>
			</Avatar.Root>
		{/snippet}
	</DropdownMenu.Trigger>
	<DropdownMenu.Content
		class="w-(--bits-dropdown-menu-anchor-width) min-w-56 rounded-lg"
		side="right"
		align="end"
		sideOffset={4}
	>
		<DropdownMenu.Label class="p-0 font-normal">
			<div class="flex items-center gap-2 px-1 py-1.5 text-start text-sm">
				<Avatar.Root class="size-8 rounded-lg">
					<Avatar.Image src={user?.avatar_url} alt={user?.full_name || 'User'} />
					<Avatar.Fallback class="rounded-lg"
						>{user?.full_name.slice(0, 2).toUpperCase() || 'CN'}</Avatar.Fallback
					>
				</Avatar.Root>
				<div class="grid flex-1 text-start text-sm leading-tight">
					<span class="truncate font-medium">{user?.full_name}</span>
					<span class="truncate text-xs">{user?.email}</span>
				</div>
			</div>
		</DropdownMenu.Label>
		<DropdownMenu.Separator />
		<DropdownMenu.Group>
			<DropdownMenu.Item>
				<SparklesIcon />
				Upgrade to Pro
			</DropdownMenu.Item>
		</DropdownMenu.Group>
		<DropdownMenu.Separator />
		<DropdownMenu.Group>
			<DropdownMenu.Item>
				<BadgeCheckIcon />
				Account
			</DropdownMenu.Item>
			<DropdownMenu.Item>
				<CreditCardIcon />
				Billing
			</DropdownMenu.Item>
			<DropdownMenu.Item>
				<BellIcon />
				Notifications
			</DropdownMenu.Item>
		</DropdownMenu.Group>
		<DropdownMenu.Separator />
		<DropdownMenu.Item onSelect={logout}>
			<LogOutIcon />
			Log out
		</DropdownMenu.Item>
	</DropdownMenu.Content>
</DropdownMenu.Root>
