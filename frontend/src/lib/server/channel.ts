import { BaseHandler } from './base';

/** Channel as returned by GET /channels/protected. */
export interface Channel {
	id: string;
	tenant_id: string;
	division_id?: string | null;
	type: string; // telegram | whatsapp_waha | whatsapp_meta | messenger | instagram | webchat | linkchat
	name: string;
	status?: string;
	ai_agent_id?: string | null;
	created_at?: string;
}

export class ChannelHandler extends BaseHandler {
	/** List the tenant's configured channels. */
	async list(type?: string): Promise<Channel[]> {
		const suffix = type ? `?type=${encodeURIComponent(type)}` : '';
		const res = await this.api.authRequest<Channel[]>('GET', `/channels/protected/${suffix}`);
		if (!res.success || !res.data) return [];
		return res.data;
	}

	/** Count configured channels per type (for the integrations catalog). */
	async countsByType(): Promise<Record<string, number>> {
		const channels = await this.list();
		const counts: Record<string, number> = {};
		for (const c of channels) {
			counts[c.type] = (counts[c.type] ?? 0) + 1;
		}
		return counts;
	}

	async get(id: string): Promise<Channel | null> {
		const res = await this.api.authRequest<Channel>('GET', `/channels/protected/${id}`);
		return res.success && res.data ? res.data : null;
	}
}
