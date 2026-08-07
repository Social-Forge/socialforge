import { BaseHandler } from './base';

export interface PlanFeatures {
	divisions?: number;
	agents?: number;
	ai_agents?: number;
	ai_credits?: number;
	quick_replies?: number;
	waha_whatsapp?: number;
	meta_whatsapp?: number;
	meta_messenger?: number;
	instagram?: number;
	telegram?: number;
	webchat?: number;
	linkchat?: number;
	[k: string]: number | undefined;
}

export interface Plan {
	id: string;
	code: string;
	name: string;
	price: number;
	currency: string;
	interval: string;
	features: PlanFeatures;
	is_active: boolean;
	sort: number;
}

/** Plan enriched with UI display fields. */
export interface PlanView extends Plan {
	priceLabel: string;
	popular: boolean;
	benefits: string[];
}

export interface Subscription {
	subscription: {
		status: string;
		current_period_end?: string | null;
		plan_id: string;
	} | null;
	plan: Plan | null;
	tenant: { subscription_plan: string; subscription_status: string; ai_credits: number } | null;
}

export interface Invoice {
	id: string;
	number: number;
	status: string;
	amount: number;
	currency: string;
	description: string;
	provider: string;
	checkout_url?: string | null;
	paid_at?: string | null;
	created_at: string;
}

export interface CheckoutRequest {
	kind: 'subscription' | 'addon';
	provider: 'xendit' | 'midtrans' | 'paypal';
	plan_code?: string;
	months?: number;
	addon_type?: 'channel_slot' | 'agent_slot' | 'ai_credits';
	quantity?: number;
}

function formatIDR(n: number): string {
	if (n === 0) return 'Gratis';
	return 'Rp' + n.toLocaleString('id-ID');
}

function buildBenefits(f: PlanFeatures): string[] {
	const out: string[] = [];
	const line = (n: number | undefined, label: string) => {
		if (n == null) return;
		out.push(`${n.toLocaleString('id-ID')} ${label}`);
	};
	line(f.agents, 'agen');
	line(f.divisions, 'divisi');
	line(f.ai_agents, 'AI agent');
	if (f.ai_credits != null) out.push(`${f.ai_credits.toLocaleString('id-ID')} kredit AI / bulan`);
	const channels =
		(f.waha_whatsapp ?? 0) +
		(f.meta_whatsapp ?? 0) +
		(f.meta_messenger ?? 0) +
		(f.instagram ?? 0) +
		(f.telegram ?? 0) +
		(f.webchat ?? 0) +
		(f.linkchat ?? 0);
	out.push(`${channels} slot channel`);
	line(f.quick_replies, 'quick reply');
	return out;
}

export function toPlanView(p: Plan): PlanView {
	return {
		...p,
		priceLabel: formatIDR(p.price),
		popular: p.code === 'pro',
		benefits: buildBenefits(p.features || {})
	};
}

export class BillingHandler extends BaseHandler {
	/** Public plan catalog (active plans), enriched for display. */
	async plans(): Promise<PlanView[]> {
		const res = await this.api.publicRequest<Plan[]>('GET', '/plans/');
		if (!res.success || !res.data) return [];
		return res.data.map(toPlanView);
	}

	/** Current tenant subscription + plan + limits. */
	async subscription(): Promise<Subscription | null> {
		const res = await this.api.authRequest<Subscription>('GET', '/billing/protected/subscription');
		return res.success ? (res.data ?? null) : null;
	}

	/** Invoice history. */
	async invoices(): Promise<Invoice[]> {
		const res = await this.api.authRequest<Invoice[]>('GET', '/billing/protected/invoices');
		return res.success && res.data ? res.data : [];
	}

	/** Start a checkout; returns the hosted-checkout URL to redirect to. */
	async checkout(req: CheckoutRequest): Promise<{ checkoutUrl: string } | { error: string }> {
		const res = await this.api.authRequest<{ checkout_url: string }>(
			'POST',
			'/billing/protected/checkout',
			req
		);
		if (!res.success || !res.data?.checkout_url) {
			return { error: res.message || 'Checkout failed' };
		}
		return { checkoutUrl: res.data.checkout_url };
	}
}
