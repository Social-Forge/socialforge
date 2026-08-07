<script lang="ts">
	import { page } from '$app/state';
	import { goto, invalidateAll } from '$app/navigation';
	import { SvelteURL } from 'svelte/reactivity';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Badge } from '$lib/components/ui/badge/index.js';
	import { ContactFormFilter } from '$lib/components/app/contact/index.js';
	import { Separator } from '$lib/components/ui/separator/index.js';
	import { createQueryContactManager, updateUrlParams } from '$lib/stores/query.js';
	import {
		Search,
		CirclePlus,
		BrushCleaning,
		ChevronLeft,
		ChevronRight,
		Ban,
		Trash2
	} from '@lucide/svelte';
	import type { Contact, PageMeta } from '$lib/server/contact';

	let { data } = $props();
	const queryManager = createQueryContactManager();
	let query = $derived(queryManager.parse(page.url));

	const contacts = $derived((data.contacts ?? []) as Contact[]);
	const meta = $derived(data.meta as PageMeta);

	// svelte-ignore state_referenced_locally
	let searchTerm = $state((data.search ?? '') as string);

	async function updateQuery(updates: Record<string, unknown>, resetPage = false) {
		await updateUrlParams(goto, page.url, updates, {
			resetPage,
			replaceState: true,
			invalidateAll: true
		});
	}
	async function applySearch() {
		await updateQuery({ search: searchTerm }, true);
	}
	async function goPage(p: number) {
		if (p < 1 || (meta && p > meta.total_pages)) return;
		await updateQuery({ page: p });
	}
	async function onReset() {
		const url = new SvelteURL(page.url);
		url.search = '';
		searchTerm = '';
		await goto(url.toString(), { replaceState: true, invalidateAll: true });
	}

	function nstr(v: unknown): string {
		if (v && typeof v === 'object' && 'String' in v)
			return (v as any).Valid ? (v as any).String : '';
		return (v as string) ?? '';
	}
	function initials(name: string) {
		return (name || '?').trim().slice(0, 2).toUpperCase();
	}
	async function toggleBlock(c: Contact) {
		await fetch(`/api/contacts/${c.id}`, {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify({ blocked: !c.is_blocked })
		});
		await invalidateAll();
	}
	async function removeContact(c: Contact) {
		if (!confirm(`Hapus kontak "${c.display_name}"?`)) return;
		await fetch(`/api/contacts/${c.id}`, { method: 'DELETE' });
		await invalidateAll();
	}
</script>

<div class="h-[calc(100dvh-70px)] w-full overflow-y-auto bg-background lg:h-full">
	<div class="m-6 lg:m-10">
		<div class="flex flex-col gap-4">
			<div class="flex items-center justify-between">
				<div class="text-base font-medium lg:text-xl">Contacts</div>
				<div class="text-sm text-muted-foreground">{meta?.total ?? 0} kontak</div>
			</div>
			<div class="flex flex-col items-center justify-center gap-4 md:flex-row md:justify-between">
				<div class="flex items-center gap-2">
					<div class="relative">
						<Search class="absolute top-1/2 left-2 size-4 -translate-y-1/2 text-muted-foreground" />
						<Input
							type="text"
							placeholder="Search Contact"
							class="ps-9"
							bind:value={searchTerm}
							onkeydown={(e) => e.key === 'Enter' && applySearch()}
						/>
					</div>
					<Button variant="outline" onclick={applySearch}>Cari</Button>
					<ContactFormFilter {updateQuery} />
				</div>
				<div class="flex items-center gap-2">
					<Button><CirclePlus /><span class="ml-2 hidden md:block">Add Contact</span></Button>
					<Button variant="destructive" onclick={onReset}>
						<BrushCleaning /><span class="ml-2 hidden md:block">Reset</span>
					</Button>
				</div>
			</div>
		</div>
		<Separator class="my-6" />

		{#if contacts.length === 0}
			<div class="flex flex-col items-center justify-center py-24 text-center select-none">
				<div class="text-lg font-bold">Belum ada kontak.</div>
				<div class="mt-2 text-sm text-muted-foreground">
					Kontak akan muncul otomatis saat pelanggan mengirim pesan, atau tambah manual.
				</div>
			</div>
		{:else}
			<div class="overflow-x-auto rounded-lg border">
				<table class="w-full text-sm">
					<thead class="border-b bg-muted/40 text-left text-muted-foreground">
						<tr>
							<th class="p-3 font-medium">Nama</th>
							<th class="p-3 font-medium">External ID</th>
							<th class="p-3 font-medium">Status</th>
							<th class="p-3 font-medium">Dibuat</th>
							<th class="p-3"></th>
						</tr>
					</thead>
					<tbody>
						{#each contacts as c (c.id)}
							<tr class="border-b last:border-0 hover:bg-muted/30">
								<td class="p-3">
									<div class="flex items-center gap-3">
										<div
											class="flex size-9 items-center justify-center rounded-full bg-primary/10 text-xs font-semibold text-primary"
										>
											{initials(c.display_name)}
										</div>
										<span class="font-medium">{c.display_name}</span>
									</div>
								</td>
								<td class="p-3 text-muted-foreground">{c.external_id}</td>
								<td class="p-3">
									{#if c.is_blocked}
										<Badge variant="destructive">Diblokir</Badge>
									{:else}
										<Badge variant="secondary">Aktif</Badge>
									{/if}
								</td>
								<td class="p-3 text-muted-foreground">
									{new Date(c.created_at).toLocaleDateString('id-ID', {
										day: 'numeric',
										month: 'short',
										year: 'numeric'
									})}
								</td>
								<td class="p-3">
									<div class="flex justify-end gap-1">
										<Button
											variant="ghost"
											size="icon"
											title={c.is_blocked ? 'Buka blokir' : 'Blokir'}
											onclick={() => toggleBlock(c)}
										>
											<Ban class="size-4" />
										</Button>
										<Button
											variant="ghost"
											size="icon"
											class="text-red-600"
											title="Hapus"
											onclick={() => removeContact(c)}
										>
											<Trash2 class="size-4" />
										</Button>
									</div>
								</td>
							</tr>
						{/each}
					</tbody>
				</table>
			</div>

			<!-- Pagination -->
			<div class="mt-4 flex items-center justify-between">
				<div class="text-sm text-muted-foreground">
					Halaman {meta.page} dari {Math.max(meta.total_pages, 1)}
				</div>
				<div class="flex items-center gap-2">
					<Button
						variant="outline"
						size="sm"
						disabled={meta.page <= 1}
						onclick={() => goPage(meta.page - 1)}
					>
						<ChevronLeft class="size-4" /> Sebelumnya
					</Button>
					<Button
						variant="outline"
						size="sm"
						disabled={!meta.has_more}
						onclick={() => goPage(meta.page + 1)}
					>
						Berikutnya <ChevronRight class="size-4" />
					</Button>
				</div>
			</div>
		{/if}
	</div>
</div>
