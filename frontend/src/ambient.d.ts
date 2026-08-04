declare global {
	interface ApiResponse<T = any, M extends Record<string, any> = ApiMeta> {
		status: number;
		success: boolean;
		message: string;
		data?: T | null;
		error?: ApiError;
		meta?: M;
		headers?: Headers;
	}
	interface ApiMeta {
		page: number;
		limit: number;
		total_rows: number;
		total_pages: number;
		has_prev: boolean;
		has_next: boolean;
	}
	interface ApiError {
		code: string;
		message?: string;
		redirect_url?: string;
		details?: any;
		retryable?: boolean;
		timestamp?: string;
	}
	interface ErrorResponse extends ApiResponse<undefined, undefined> {
		error: {
			code: string;
			details?: Record<string, unknown>;
			redirect_url?: Record<string, unknown>;
		};
	}
	type HttpMethod = 'GET' | 'POST' | 'PATCH' | 'PUT' | 'DELETE';

	interface QueryParams {
		page: number;
		limit: number;
		search?: string;
		sort_by?: string;
		order_by?: 'asc' | 'desc';
		status?: string;
		include_deleted?: boolean;
		is_active?: boolean;
		is_verified?: boolean;
		tenant_id?: string;
		user_id?: string;
		division_id?: string;
		extra?: Record<string, any>;
	}
	interface QueryState {
		params: QueryParams;
		meta: ApiMeta | null;
		isLoading: boolean;
	}
	const DEFAULT_PAGINATION: QueryParams = {
		page: 1,
		limit: 10,
		order_by: 'desc'
	};

	interface CsrfTokenResponse extends ApiResponse<string> {
		data?: any;
	}
	// ================================
	// Database Types
	// ================================
	type RoleName = 'superadmin' | 'tenant_owner' | 'supervisor' | 'agent';
	type RoleLevel = 1 | 2 | 3 | 4;
	type UserStatus = 'active' | 'inactive' | 'suspended';
	type SubscriptionPlan = 'free' | 'starter' | 'pro' | 'enterprise';
	type SubscriptionStatus = 'active' | 'suspended' | 'canceled' | 'expired';
	type DivisionRoutingType = 'equal' | 'percentage' | 'priority';
	type User = {
		id: string;
		email: string;
		full_name: string;
		phone?: string;
		avatar_url?: string;
		two_fa_secret?: string;
		status: UserStatus;
		email_verified_at?: string;
		last_login_at?: string;
		created_at: string;
		updated_at: string;
		deleted_at?: string;
		user_tenants?: UserTenant[];
	};
	type Tenant = {
		id: string;
		name: string;
		slug: string;
		max_divisions: number;
		max_agents: number;
		max_quick_replies: number;
		max_waha_whatsapp: number;
		max_meta_whatsapp: number;
		max_meta_messenger: number;
		max_instagram: number;
		max_telegram: number;
		max_webchat: number;
		max_linkchat: number;
		ai_credits: number;
		subscription_plan: SubscriptionPlan;
		subscription_status: SubscriptionStatus;
		trial_ends_at?: string;
		is_active: boolean;
		settings?: Record<string, any>;
		created_at: string;
		updated_at: string;
		deleted_at?: string;
	};
	type Role = {
		id: string;
		name: RoleName;
		slug: string;
		level: RoleLevel;
		description?: string;
		created_at: string;
		updated_at: string;
	};
	type Division = {
		id: string;
		tenant_id: string;
		name: string;
		slug: string;
		description?: string;
		created_at: string;
		updated_at: string;
		routing_type: DivisionRoutingType;
		routing_config?: Record<string, any>;
		is_active: boolean;
		link_url?: string;
		created_at: string;
		updated_at: string;
		deleted_at?: string;
		members?: DivisionMember[];
	};
	type DivisionMember = {
		id: string;
		user_tenant_id: string;
		division_id: string;
		is_active: boolean;
		joined_at: string;
		created_at: string;
		updated_at: string;
		deleted_at?: string;
		user_tenant?: UserTenant;
		division?: Division;
	};
	type UserTenant = {
		id: string;
		user_id: string;
		tenant_id: string;
		role_id: string;
		created_at: string;
		updated_at: string;
		user?: User;
		tenant?: Tenant;
		role?: Role;
		division_members?: DivisionMember[];
	};
	type OAuthProvider = {
		id: string;
		user_id: string;
		provider_name: string;
		provider_id: string;
		created_at: string;
		updated_at: string;
	};
	type Channel = {
		id: string;
		tenant_id: string;
		division_id: string;
		ai_agent_id?: string;
		type: 'whatsapp_waha' | 'whatsapp_meta' | 'messenger' | 'instagram' | 'telegram';
		name: string;
		status: 'disconnected' | 'connected' | 'connecting' | 'failed';
		external_id?: string;
		waha_engine?: string;
		waha_session_name?: string;
		webhook_secret?: string;
		credentials?: Record<string, any>;
		settings?: Record<string, any>;
		created_at: string;
		updated_at: string;
	};
	type Contact = {
		id: string;
		tenant_id: string;
		channel_id: string;
		external_id: string;
		display_name: string;
		avatar_url?: string;
		is_blocked: boolean;
		attributes?: Record<string, any>;
		created_at: string;
		updated_at: string;
	};
	type Label = {
		id: string;
		tenant_id: string;
		name: string;
		color: string;
		created_at: string;
		updated_at: string;
	};
	type QuickReply = {
		id: string;
		tenant_id: string;
		shortcut: string;
		content_type: 'text' | 'image' | 'video' | 'document';
		body?: string;
		media: Record<string, any>;
		created_at: string;
		updated_at: string;
	};
	type AutoResponse = {
		id: string;
		tenant_id: string;
		channel_id: string;
		content_type: 'text' | 'image' | 'video' | 'document';
		body?: string;
		media: Record<string, any>;
		created_at: string;
		updated_at: string;
	};
	type AIAgent = {
		id: string;
		tenant_id: string;
		name: string;
		provider: 'claude' | 'openai' | 'google' | 'openrouter';
		model: string;
		system_prompt: string;
		persona?: Record<string, any>;
		safety?: Record<string, any>;
		guardrails?: Record<string, any>;
		temperature: number;
		max_tokens: number;
		auto_reply_enabled: boolean;
		working_hours?: Record<string, any>;
		is_active: boolean;
		created_at: string;
		updated_at: string;
	};
	type AgentWorkingHours = {
		id: string;
		tenant_id: string;
		user_id: string;
		day_of_week: number;
		start_time: string;
		end_time: string;
		is_active: boolean;
		created_at: string;
		updated_at: string;
	};
	type AIAsset = {
		id: string;
		tenant_id: string;
		ai_agent_id: string;
		name: string;
		type: string;
		storage_key: string;
		mime_type?: string;
		size?: number;
		description?: string;
		created_at: string;
		updated_at: string;
	};
	type AICreditLedger = {
		id: string;
		tenant_id: string;
		conversation_id?: string;
		message_id?: string;
		delta: number;
		balance_after: number;
		reason: string;
		model?: string;
		input_tokens: number;
		output_tokens: number;
		cost_usd: number;
		credit: number;
		created_at: string;
		updated_at: string;
	};
	type AIKnowledge = {
		id: string;
		tenant_id: string;
		ai_agent_id: string;
		title: string;
		content: string;
		token_count: number;
		created_at: string;
		updated_at: string;
	};
	type AIPlaybook = {
		id: string;
		tenant_id: string;
		ai_agent_id: string;
		name: string;
		keywords: string[];
		instruction: string;
		asset_ids: string[];
		priority: number;
		is_active: boolean;
		created_at: string;
		updated_at: string;
	};

	type Convertation = {
		id: string;
		tenant_id: string;
		channel_id: string;
		contact_id: string;
		assigned_agent_id?: string;
		status: 'open' | 'unassigned' | 'completed' | 'archived';
		is_pinned: boolean;
		is_archived: boolean;
		unread_count: number;
		last_message_at?: string;
		service_window_expires_at?: string;
		metadata?: Record<string, any>;
		created_at: string;
		updated_at: string;
	};
	type ConversationLabel = {
		id: string;
		tenant_id: string;
		conversation_id: string;
		label_id: string;
		created_at: string;
		updated_at: string;
	};
	type Message = {
		id: string;
		tenant_id: string;
		conversation_id: string;
		sender_id?: string;
		direction: 'in' | 'out';
		sender_type: 'contact' | 'agent' | 'ai' | 'system';
		content_type:
			| 'text'
			| 'image'
			| 'video'
			| 'audio'
			| 'document'
			| 'location'
			| 'contact'
			| 'sticker'
			| 'button'
			| 'list'
			| 'template';
		body?: string;
		media?: Record<string, any>;
		provider_message_id?: string;
		status: 'pending' | 'sent' | 'delivered' | 'read' | 'failed';
		reply_to_id?: string;
		is_pinned: boolean;
		error?: string;
		edited_at?: string;
		created_at: string;
		updated_at: string;
		deleted_at?: string;
	};
	type MessageOutbox = {
		id: string;
		tenant_id: string;
		message_id: string;
		status: string;
		attempts: number;
		max_attempts: number;
		next_retry_at?: string;
		last_error?: string;
		created_at: string;
		updated_at: string;
	};
	// ================================
	// Response Types
	// ================================
	type UserResponse = {
		id: string;
		tenant_id: string;
		role_id: string;
		email: string;
		full_name: string;
		phone?: string;
		avatar_url?: string;
		status: UserStatus;
		is_active: boolean;
		is_verified: boolean;
		email_verified_at?: string;
		last_login_at?: string;
		created_at: string;
		updated_at: string;
		tenant?: Tenant;
		user_tenant?: UserTenant;
		role?: Role;
	};
	type LoginResponse = {
		access_token: string;
		refresh_token: string;
		two_fa_token?: string;
		expires_in: number;
		token_type: string;
		expires_refresh_in: number;
		user: UserResponse;
		status: string;
	};
}

export {};

declare module 'sveltekit-autoimport';
