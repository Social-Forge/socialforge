<script lang="ts">
	import { page } from '$app/state';
	import { goto } from '$app/navigation';
	import { invalidateAll } from '$app/navigation';
	import * as UnderlineTabs from '$lib/components/extras/underline-tabs';
	import { SettingLayout } from '$lib/components/app/setting/index.js';
	import * as Card from '$lib/components/ui/card';
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import { Loader2, Check } from '@lucide/svelte';
	import { localizeHref } from '$lib/paraglide/runtime.js';

	let { data } = $props();
	const user = $derived(data.user as UserResponse | undefined);
	let tab = $state(page.url.searchParams.get('tab') ?? 'profile');

	// Profile form
	let fullName = $state('');
	let email = $state('');
	let phone = $state('');
	$effect(() => {
		fullName = user?.full_name ?? '';
		email = user?.email ?? '';
		phone = user?.phone ?? '';
	});

	let savingProfile = $state(false);
	let profileMsg = $state<{ ok: boolean; text: string } | null>(null);

	async function saveProfile() {
		savingProfile = true;
		profileMsg = null;
		try {
			const res = await fetch('/api/user/profile', {
				method: 'PUT',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ full_name: fullName, email, phone })
			});
			const j = await res.json();
			profileMsg = {
				ok: j.success,
				text: j.success ? 'Profil tersimpan' : j.message || 'Gagal menyimpan'
			};
			if (j.success) await invalidateAll();
		} finally {
			savingProfile = false;
		}
	}

	// Password form
	let currentPassword = $state('');
	let newPassword = $state('');
	let confirmPassword = $state('');
	let savingPassword = $state(false);
	let passwordMsg = $state<{ ok: boolean; text: string } | null>(null);

	async function savePassword() {
		passwordMsg = null;
		if (newPassword !== confirmPassword) {
			passwordMsg = { ok: false, text: 'Konfirmasi password tidak cocok' };
			return;
		}
		savingPassword = true;
		try {
			const res = await fetch('/api/user/password', {
				method: 'PUT',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({
					current_password: currentPassword,
					new_password: newPassword,
					confirm_password: confirmPassword
				})
			});
			const j = await res.json();
			passwordMsg = {
				ok: j.success,
				text: j.success ? 'Password diperbarui' : j.message || 'Gagal memperbarui'
			};
			if (j.success) {
				currentPassword = '';
				newPassword = '';
				confirmPassword = '';
			}
		} finally {
			savingPassword = false;
		}
	}

	function initials(name: string) {
		return (name || '?').trim().slice(0, 2).toUpperCase();
	}

	async function onTabChange(tab: string) {
		tab = tab;
		await goto(localizeHref(`/app/settings/account?tab=${tab}`));
	}
</script>

<SettingLayout>
	<div class="flex w-full flex-col gap-8">
		<div class="rounded-md bg-card p-2 shadow-md lg:p-7">
			<UnderlineTabs.Root bind:value={tab} onValueChange={(tab) => onTabChange(tab)}>
				<UnderlineTabs.List>
					<UnderlineTabs.Trigger value="profile">Profil Akun</UnderlineTabs.Trigger>
					<UnderlineTabs.Trigger value="password">Password Akun</UnderlineTabs.Trigger>
				</UnderlineTabs.List>
				<UnderlineTabs.Content value="profile">
					<Card.Root>
						<Card.Header>
							<Card.Title>Profil Akun</Card.Title>
							<Card.Description>Perbarui informasi akun Anda.</Card.Description>
						</Card.Header>
						<Card.Content class="flex flex-col gap-4">
							<div class="flex items-center gap-4">
								<div
									class="flex size-16 items-center justify-center rounded-full bg-primary/10 text-lg font-semibold text-primary"
								>
									{initials(fullName)}
								</div>
								<div>
									<div class="font-medium">{fullName || '—'}</div>
									<div class="text-sm text-muted-foreground">{email}</div>
								</div>
							</div>
							<div class="flex flex-col gap-1.5">
								<Label>Nama Lengkap</Label>
								<Input bind:value={fullName} placeholder="Nama lengkap" />
							</div>
							<div class="flex flex-col gap-1.5">
								<Label>Email</Label>
								<Input type="email" bind:value={email} placeholder="email@contoh.com" />
							</div>
							<div class="flex flex-col gap-1.5">
								<Label>Nomor Telepon</Label>
								<Input bind:value={phone} placeholder="0812xxxxxxx" />
							</div>
							{#if profileMsg}
								<div
									class="flex items-center gap-2 rounded-md px-3 py-2 text-sm {profileMsg.ok
										? 'bg-emerald-500/10 text-emerald-600'
										: 'bg-red-500/10 text-red-600'}"
								>
									{#if profileMsg.ok}<Check class="h-4 w-4" />{/if}{profileMsg.text}
								</div>
							{/if}
						</Card.Content>
						<Card.Footer>
							<Button onclick={saveProfile} disabled={savingProfile || !fullName.trim()}>
								{#if savingProfile}<Loader2 class="mr-2 h-4 w-4 animate-spin" />{/if}
								Simpan Perubahan
							</Button>
						</Card.Footer>
					</Card.Root>
				</UnderlineTabs.Content>
				<UnderlineTabs.Content value="password">
					<Card.Root>
						<Card.Header>
							<Card.Title>Ubah Password</Card.Title>
							<Card.Description>Gunakan password yang kuat dan unik.</Card.Description>
						</Card.Header>
						<Card.Content class="flex flex-col gap-4">
							<div class="flex flex-col gap-1.5">
								<Label>Password Saat Ini</Label>
								<Input
									type="password"
									bind:value={currentPassword}
									autocomplete="current-password"
								/>
							</div>
							<div class="flex flex-col gap-1.5">
								<Label>Password Baru</Label>
								<Input type="password" bind:value={newPassword} autocomplete="new-password" />
							</div>
							<div class="flex flex-col gap-1.5">
								<Label>Konfirmasi Password Baru</Label>
								<Input type="password" bind:value={confirmPassword} autocomplete="new-password" />
							</div>
							{#if passwordMsg}
								<div
									class="rounded-md px-3 py-2 text-sm {passwordMsg.ok
										? 'bg-emerald-500/10 text-emerald-600'
										: 'bg-red-500/10 text-red-600'}"
								>
									{passwordMsg.text}
								</div>
							{/if}
						</Card.Content>
						<Card.Footer>
							<Button
								onclick={savePassword}
								disabled={savingPassword || !currentPassword || !newPassword}
							>
								{#if savingPassword}<Loader2 class="mr-2 h-4 w-4 animate-spin" />{/if}
								Perbarui Password
							</Button>
						</Card.Footer>
					</Card.Root>
				</UnderlineTabs.Content>
			</UnderlineTabs.Root>
		</div>
	</div>
</SettingLayout>
