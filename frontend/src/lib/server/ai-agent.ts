import { BaseHandler } from './base';

export interface AiPersona {
	agent_name?: string;
	soul?: string;
	tone?: string;
	gender?: string;
	greeting?: string;
	language?: string;
	[k: string]: unknown;
}
export interface AiSafety {
	handoff_to_human?: boolean;
	sensitive_topics?: string[];
	[k: string]: unknown;
}
export interface AiGuardrails {
	rules?: string[];
	[k: string]: unknown;
}

export interface AIAgent {
	id: string;
	name: string;
	provider: 'claude' | 'openai' | 'google';
	model: string;
	system_prompt: string;
	persona?: AiPersona | null;
	safety?: AiSafety | null;
	guardrails?: AiGuardrails | null;
	temperature: number;
	max_tokens: number;
	auto_reply_enabled: boolean;
	is_active: boolean;
	created_at?: string;
	updated_at?: string;
}

export interface AIAgentPayload {
	name: string;
	provider: string;
	model?: string;
	system_prompt: string;
	persona?: AiPersona;
	safety?: AiSafety;
	guardrails?: AiGuardrails;
	temperature?: number;
	max_tokens?: number;
	auto_reply_enabled?: boolean;
	is_active?: boolean;
}

export interface AIKnowledge {
	id: string;
	ai_agent_id: string;
	title: string;
	content: string;
	token_count: number;
}
export interface AIPlaybook {
	id: string;
	name: string;
	keywords: string[];
	instruction: string;
	asset_ids: string[];
	priority: number;
	is_active: boolean;
}
export interface AIAsset {
	id: string;
	name: string;
	type: 'image' | 'video' | 'document';
	storage_key: string;
	mime_type?: { String: string; Valid: boolean } | string | null;
	size?: unknown;
	description?: { String: string; Valid: boolean } | string | null;
}

type ApiRes<T> = { success: boolean; message?: string; data?: T };

export class AIAgentHandler extends BaseHandler {
	private base = '/ai-agents/protected';

	// ---- Agents ----
	async list(): Promise<AIAgent[]> {
		const res = await this.api.authRequest<AIAgent[]>('GET', `${this.base}/`);
		return res.success && res.data ? res.data : [];
	}
	async get(id: string): Promise<AIAgent | null> {
		const res = await this.api.authRequest<AIAgent>('GET', `${this.base}/${id}`);
		return res.success ? (res.data ?? null) : null;
	}
	async create(payload: AIAgentPayload): Promise<ApiRes<AIAgent>> {
		return this.api.authRequest<AIAgent>('POST', `${this.base}/`, payload);
	}
	async update(id: string, payload: AIAgentPayload): Promise<ApiRes<AIAgent>> {
		return this.api.authRequest<AIAgent>('PUT', `${this.base}/${id}`, payload);
	}
	async remove(id: string): Promise<ApiRes<null>> {
		return this.api.authRequest<null>('DELETE', `${this.base}/${id}`);
	}

	// ---- Nested resources (knowledge | playbooks | assets) ----
	async listResource<T>(agentId: string, resource: string): Promise<T[]> {
		const res = await this.api.authRequest<T[]>('GET', `${this.base}/${agentId}/${resource}/`);
		return res.success && res.data ? res.data : [];
	}
	async createResource<T>(agentId: string, resource: string, payload: unknown): Promise<ApiRes<T>> {
		return this.api.authRequest<T>('POST', `${this.base}/${agentId}/${resource}/`, payload);
	}
	async updateResource<T>(
		agentId: string,
		resource: string,
		childId: string,
		payload: unknown
	): Promise<ApiRes<T>> {
		return this.api.authRequest<T>('PUT', `${this.base}/${agentId}/${resource}/${childId}`, payload);
	}
	async deleteResource(agentId: string, resource: string, childId: string): Promise<ApiRes<null>> {
		return this.api.authRequest<null>('DELETE', `${this.base}/${agentId}/${resource}/${childId}`);
	}
}
