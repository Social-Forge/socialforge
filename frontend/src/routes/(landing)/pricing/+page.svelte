<script lang="ts">
	import * as Card from '$lib/components/ui/card';
	import { Button } from '$lib/components/ui/button';
	import { Badge } from '$lib/components/ui/badge';
	import { Check } from '@lucide/svelte';
	import { localizeHref } from '$lib/paraglide/runtime.js';
	import type { PlanView } from '$lib/server/billing';

	let { data } = $props();
	const plans = $derived((data.plans ?? []) as PlanView[]);
</script>

<section class="mx-auto w-full max-w-6xl px-4 py-16 md:py-24">
	<div class="mb-12 flex flex-col items-center gap-3 text-center">
		<Badge variant="secondary">Harga</Badge>
		<h1 class="text-3xl font-bold tracking-tight md:text-4xl">Pilih paket yang tepat</h1>
		<p class="max-w-xl text-muted-foreground">
			Mulai gratis, tingkatkan kapan saja. Semua harga dalam Rupiah per bulan.
		</p>
	</div>

	{#if plans.length === 0}
		<p class="text-center text-muted-foreground">Paket belum tersedia.</p>
	{:else}
		<div class="grid gap-6 md:grid-cols-2 lg:grid-cols-4">
			{#each plans as plan (plan.id)}
				<Card.Root class="relative flex flex-col {plan.popular ? 'border-primary shadow-lg' : ''}">
					{#if plan.popular}
						<div class="absolute -top-3 left-1/2 -translate-x-1/2">
							<Badge>Populer</Badge>
						</div>
					{/if}
					<Card.Header>
						<Card.Title class="text-lg">{plan.name}</Card.Title>
						<div class="mt-2 flex items-end gap-1">
							<span class="text-3xl font-bold">{plan.priceLabel}</span>
							{#if plan.price > 0}
								<span class="pb-1 text-sm text-muted-foreground">/bulan</span>
							{/if}
						</div>
					</Card.Header>
					<Card.Content class="flex flex-1 flex-col gap-3">
						<ul class="flex flex-col gap-2 text-sm">
							{#each plan.benefits as benefit}
								<li class="flex items-center gap-2">
									<Check class="h-4 w-4 shrink-0 text-primary" />
									<span>{benefit}</span>
								</li>
							{/each}
						</ul>
					</Card.Content>
					<Card.Footer>
						<Button
							href={localizeHref(plan.price === 0 ? '/signup' : '/app/settings/billing')}
							variant={plan.popular ? 'default' : 'outline'}
							class="w-full"
						>
							{plan.price === 0 ? 'Mulai Gratis' : 'Pilih Paket'}
						</Button>
					</Card.Footer>
				</Card.Root>
			{/each}
		</div>
	{/if}
</section>
