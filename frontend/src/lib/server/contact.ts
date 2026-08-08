import { BaseHandler } from './base';

export interface Contact {
	id: string;
	channel_id: string;
	external_id: string;
	display_name: string;
	avatar_url?: { String: string; Valid: boolean } | string | null;
	is_blocked: boolean;
	created_at: string;
	updated_at: string;
}

export interface ContactListParams {
	page?: number;
	perPage?: number;
	search?: string;
	channelId?: string;
}

export class ContactHandler extends BaseHandler {
	async list(params: ContactListParams = {}): Promise<{ contacts: Contact[]; meta: PageMeta }> {
		const qs = new URLSearchParams();
		qs.set('page', String(params.page ?? 1));
		qs.set('per_page', String(params.perPage ?? 20));
		if (params.search) qs.set('search', params.search);
		if (params.channelId) qs.set('channel_id', params.channelId);
		const res = await this.api.authRequest<Contact[]>('GET', `/contacts/protected/?${qs}`);
		const meta = (res.meta as PageMeta) ?? {
			page: 1,
			per_page: 20,
			total: 0,
			total_pages: 0,
			has_more: false
		};
		return { contacts: res.success && res.data ? res.data : [], meta };
	}

	async setBlocked(id: string, blocked: boolean) {
		return this.api.authRequest(
			'POST',
			`/contacts/protected/${id}/${blocked ? 'block' : 'unblock'}`
		);
	}

	async remove(id: string) {
		return this.api.authRequest('DELETE', `/contacts/protected/${id}`);
	}
}

