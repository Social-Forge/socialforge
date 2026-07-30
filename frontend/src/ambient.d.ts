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
