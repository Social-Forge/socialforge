<script lang="ts">
	import * as Card from '$lib/components/ui/card';
	import { Button } from '$lib/components/ui/button';
	import { Badge } from '$lib/components/ui/badge';
	import { Check, Loader2 } from '@lucide/svelte';
	import type { PlanView, Subscription, Invoice } from '$lib/server/billing';

	let { data } = $props();
	const sub = $derived(data.subscription as Subscription | null);
	const plans = $derived((data.plans ?? []) as PlanView[]);
	const invoices = $derived((data.invoices ?? []) as Invoice[]);

	const currentCode = $derived(sub?.tenant?.subscription_plan ?? 'free');
	const aiCredits = $derived(sub?.tenant?.ai_credits ?? 0);

	let checkingOut = $state<string | null>(null);
	let errorMsg = $state('');

	function fmtIDR(n: number) {
		return 'Rp' + (n ?? 0).toLocaleString('id-ID');
	}
	function fmtDate(s?: string | null) {
		return s
			? new Date(s).toLocaleDateString('id-ID', { day: 'numeric', month: 'short', year: 'numeric' })
			: '-';
	}
	const statusColor: Record<string, string> = {
		paid: 'bg-emerald-500/15 text-emerald-600',
		pending: 'bg-amber-500/15 text-amber-600',
		expired: 'bg-muted text-muted-foreground',
		failed: 'bg-red-500/15 text-red-600'
	};

	async function subscribe(planCode: string) {
		errorMsg = '';
		checkingOut = planCode;
		try {
			const res = await fetch('/api/billing/checkout', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({
					kind: 'subscription',
					provider: 'xendit',
					plan_code: planCode,
					months: 1
				})
			});
			const json = await res.json();
			if (json.success && json.checkout_url) {
				window.location.href = json.checkout_url;
				return;
			}
			errorMsg = json.message || 'Gagal memulai pembayaran';
		} catch {
			errorMsg = 'Terjadi kesalahan jaringan';
		} finally {
			checkingOut = null;
		}
	}
</script>

<div class="mx-auto flex w-full max-w-5xl flex-col gap-6 p-6">
	<div>
		<h1 class="text-xl font-semibold">Langganan & Pembayaran</h1>
		<p class="text-sm text-muted-foreground">Kelola paket, kredit AI, dan riwayat tagihan.</p>
	</div>

	<!-- Current subscription -->
	<Card.Root>
		<Card.Header>
			<Card.Title>Paket Saat Ini</Card.Title>
		</Card.Header>
		<Card.Content class="flex flex-wrap items-center gap-x-10 gap-y-4">
			<div>
				<div class="text-xs text-muted-foreground">Paket</div>
				<div class="flex items-center gap-2 text-lg font-semibold capitalize">
					{currentCode}
					<Badge variant="secondary" class="capitalize">
						{sub?.tenant?.subscription_status ?? 'active'}
					</Badge>
				</div>
			</div>
			<div>
				<div class="text-xs text-muted-foreground">Kredit AI</div>
				<div class="text-lg font-semibold">{aiCredits.toLocaleString('id-ID')}</div>
			</div>
			<div>
				<div class="text-xs text-muted-foreground">Berlaku Hingga</div>
				<div class="text-lg font-semibold">{fmtDate(sub?.subscription?.current_period_end)}</div>
			</div>
		</Card.Content>
	</Card.Root>

	{#if errorMsg}
		<div class="rounded-md bg-red-500/10 px-4 py-2 text-sm text-red-600">{errorMsg}</div>
	{/if}

	<!-- Plans -->
	<div>
		<h2 class="mb-3 text-lg font-semibold">Ubah Paket</h2>
		<div class="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
			{#each plans as plan (plan.id)}
				{@const isCurrent = plan.code === currentCode}
				<Card.Root class="flex flex-col {plan.popular ? 'border-primary' : ''}">
					<Card.Header>
						<Card.Title class="flex items-center justify-between text-base">
							{plan.name}
							{#if plan.popular}<Badge>Populer</Badge>{/if}
						</Card.Title>
						<div class="mt-1 text-2xl font-bold">{plan.priceLabel}</div>
					</Card.Header>
					<Card.Content class="flex flex-1 flex-col gap-2 text-sm">
						{#each plan.benefits as b}
							<div class="flex items-center gap-2">
								<Check class="h-4 w-4 shrink-0 text-primary" /><span>{b}</span>
							</div>
						{/each}
					</Card.Content>
					<Card.Footer>
						{#if isCurrent}
							<Button variant="outline" class="w-full" disabled>Paket Aktif</Button>
						{:else if plan.price === 0}
							<Button variant="outline" class="w-full" disabled>Gratis</Button>
						{:else}
							<Button
								class="w-full"
								onclick={() => subscribe(plan.code)}
								disabled={checkingOut !== null}
							>
								{#if checkingOut === plan.code}
									<Loader2 class="mr-2 h-4 w-4 animate-spin" /> Memproses…
								{:else}
									Pilih Paket
								{/if}
							</Button>
						{/if}
					</Card.Footer>
				</Card.Root>
			{/each}
		</div>
	</div>

	<!-- Invoices -->
	<div>
		<h2 class="mb-3 text-lg font-semibold">Riwayat Tagihan</h2>
		<Card.Root>
			<Card.Content class="p-0">
				{#if invoices.length === 0}
					<div class="p-6 text-center text-sm text-muted-foreground">Belum ada tagihan.</div>
				{:else}
					<div class="overflow-x-auto">
						<table class="w-full text-sm">
							<thead class="border-b text-left text-muted-foreground">
								<tr>
									<th class="p-3 font-medium">No.</th>
									<th class="p-3 font-medium">Deskripsi</th>
									<th class="p-3 font-medium">Jumlah</th>
									<th class="p-3 font-medium">Status</th>
									<th class="p-3 font-medium">Tanggal</th>
									<th class="p-3"></th>
								</tr>
							</thead>
							<tbody>
								{#each invoices as inv (inv.id)}
									<tr class="border-b last:border-0">
										<td class="p-3">#{inv.number}</td>
										<td class="p-3">{inv.description}</td>
										<td class="p-3">{fmtIDR(inv.amount)}</td>
										<td class="p-3">
											<span
												class="rounded-full px-2 py-0.5 text-xs {statusColor[inv.status] ??
													'bg-muted'}"
											>
												{inv.status}
											</span>
										</td>
										<td class="p-3">{fmtDate(inv.created_at)}</td>
										<td class="p-3 text-right">
											{#if inv.status === 'pending' && inv.checkout_url}
												<a class="text-primary hover:underline" href={inv.checkout_url}>Bayar</a>
											{/if}
										</td>
									</tr>
								{/each}
							</tbody>
						</table>
					</div>
				{/if}
			</Card.Content>
		</Card.Root>
	</div>
</div>
