declare global {
	interface ApiResponse<T = any, M extends Record<string, any> = ApiMeta> {
		status: number;
		success: boolean;
		message: string;
		data?: T | null;
		error?: FetchError;
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
	interface FetchError {
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
	interface FetchOptions {
		method?: HttpMethod;
		body?: any;
		headers?: HeadersInit;
		swrKey?: string;
		params?: Record<string, string | number | boolean | undefined>;
		timeout?: number;
		retries?: number;
		cache?: boolean;
		limit?: number;
		page?: number;
	}

	interface QueryParams {
		page: number;
		limit: number;
		search?: string;
		sort_by?: string;
		order_by?: 'asc' | 'desc';
		status?: string;
		include_deleted?: boolean;
		with_relations?: boolean;
		with_delete_column?: boolean;
		is_active?: boolean;
		is_verified?: boolean;
		user_id?: string;
		includes?: string[]; // e.g., ['user', 'category']
		fields?: string[]; // e.g., ['id', 'name', 'email']
		date_from?: Date | string;
		date_to?: Date | string;
		extra?: Record<string, any>;
	}
	interface PaginatedResult<T> {
		data: T[];
		pagination: {
			current_page: number;
			total_pages: number;
			total_items: number;
			has_next: boolean;
			has_prev: boolean;
			limit: number;
		};
	}
	interface QueryState {
		page: number;
		limit: number;
		search: string;
		type: string;
		sort_by: string;
		order_by: 'asc' | 'desc';
		params: QueryParams;
		date_from: Date | string;
		date_to: Date | string;
	}
	interface QueryParamsConfig<T = any> {
		defaults: T;
		validators?: Partial<Record<keyof T, (value: any) => any>>;
	}
	interface SearchFieldConfig {
		field: string;
		type: 'string' | 'number' | 'date' | 'boolean';
		searchable?: boolean; // Default: true
	}
}

export {};

declare module 'sveltekit-autoimport';
