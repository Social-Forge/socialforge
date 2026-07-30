export function parseQueryParams<T extends Record<string, any>>(
	url: URL,
	config: QueryParamsConfig<T>
): T {
	const params: any = { ...config.defaults };

	Object.keys(config.defaults).forEach((key) => {
		const value = url.searchParams.get(key);

		if (value !== null) {
			// Apply validator if exists
			if (config.validators && config.validators[key as keyof T]) {
				params[key] = config.validators[key as keyof T]!(value);
			} else {
				// Auto-detect type based on default value
				const defaultValue = config.defaults[key];

				if (typeof defaultValue === 'number') {
					params[key] = Number(value) || defaultValue;
				} else if (typeof defaultValue === 'boolean') {
					params[key] = value === 'true';
				} else {
					params[key] = value;
				}
			}
		}
	});

	return params as T;
}
export function buildUrlWithParams(
	baseUrl: string | URL,
	params: Record<string, any>,
	options: {
		skipEmpty?: boolean;
		skipDefaults?: boolean;
		defaults?: Record<string, any>;
	} = {}
): URL {
	const url = typeof baseUrl === 'string' ? new URL(baseUrl) : new URL(baseUrl.toString());

	Object.entries(params).forEach(([key, value]) => {
		// Skip empty values if configured
		if (options.skipEmpty && (value === '' || value === null || value === undefined)) {
			url.searchParams.delete(key);
			return;
		}

		// Skip default values if configured
		if (options.skipDefaults && options.defaults && value === options.defaults[key]) {
			url.searchParams.delete(key);
			return;
		}

		// Set param
		if (value !== undefined && value !== null) {
			url.searchParams.set(key, String(value));
		}
	});

	return url;
}
export async function updateUrlParams(
	goto: (url: string, opts?: any) => Promise<void>,
	currentUrl: URL,
	updates: Record<string, any>,
	options: {
		resetPage?: boolean;
		replaceState?: boolean;
		invalidateAll?: boolean;
	} = {}
): Promise<void> {
	const url = new URL(currentUrl);

	// Reset page to 1 if configured
	if (options.resetPage) {
		updates.page = 1;
	}

	// Update params
	Object.entries(updates).forEach(([key, value]) => {
		if (value !== undefined && value !== null && value !== '') {
			url.searchParams.set(key, String(value));
		} else {
			url.searchParams.delete(key);
		}
	});

	// Navigate
	await goto(url.toString(), {
		replaceState: options.replaceState ?? true,
		invalidateAll: options.invalidateAll ?? true
	});
}
export function clearParams(url: URL, keys: string[]): URL {
	const newUrl = new URL(url);
	keys.forEach((key) => newUrl.searchParams.delete(key));
	return newUrl;
}
export function getAllParams(url: URL): Record<string, string> {
	return Object.fromEntries(url.searchParams.entries());
}
export function hasParam(url: URL, key: string): boolean {
	return url.searchParams.has(key);
}
export function getParam(url: URL, key: string, fallback: string = ''): string {
	return url.searchParams.get(key) || fallback;
}
export function getParamAsNumber(url: URL, key: string, fallback: number = 0): number {
	const value = url.searchParams.get(key);
	const parsed = Number(value);
	return isNaN(parsed) ? fallback : parsed;
}
export function getParamAsBoolean(url: URL, key: string, fallback: boolean = false): boolean {
	const value = url.searchParams.get(key);
	if (value === null) return fallback;
	return value === 'true' || value === '1';
}
export function getParamAsArray(url: URL, key: string, separator: string = ','): string[] {
	const value = url.searchParams.get(key);
	if (!value) return [];
	return value.split(separator).filter(Boolean);
}
export function mergeParams(
	url: URL,
	newParams: Record<string, any>,
	strategy: 'replace' | 'merge' = 'merge'
): URL {
	const result = new URL(url);

	if (strategy === 'replace') {
		// Clear existing params
		result.search = '';
	}

	// Add/update params
	Object.entries(newParams).forEach(([key, value]) => {
		if (value !== undefined && value !== null && value !== '') {
			result.searchParams.set(key, String(value));
		}
	});

	return result;
}
export function objectToQueryString(
	obj: Record<string, any>,
	options: {
		skipEmpty?: boolean;
		encode?: boolean;
	} = {}
): string {
	const params = new URLSearchParams();

	Object.entries(obj).forEach(([key, value]) => {
		if (options.skipEmpty && (value === '' || value === null || value === undefined)) {
			return;
		}

		if (value !== undefined && value !== null) {
			params.set(key, String(value));
		}
	});

	const queryString = params.toString();
	return options.encode ? queryString : decodeURIComponent(queryString);
}
export function queryStringToObject(queryString: string): Record<string, string> {
	const params = new URLSearchParams(queryString);
	return Object.fromEntries(params.entries());
}
export function validatePaginationParams(params: { page?: number; limit?: number }): {
	page: number;
	limit: number;
} {
	return {
		page: Math.max(1, params.page || 1),
		limit: Math.min(Math.max(1, params.limit || 10), 100) // Max 100
	};
}
export function buildFilterSummary(
	params: Record<string, any>,
	labels: Record<string, string>
): Array<{ key: string; label: string; value: any }> {
	return Object.entries(params)
		.filter(([key, value]) => value !== undefined && value !== null && value !== '')
		.map(([key, value]) => ({
			key,
			label: labels[key] || key,
			value
		}));
}
export function createUrlStateManager<T extends Record<string, any>>(config: {
	defaults: T;
	validators?: Partial<Record<keyof T, (value: any) => any>>;
}) {
	return {
		/**
		 * Parse dari URL
		 */
		parse: (url: URL): T => {
			return parseQueryParams(url, config);
		},

		/**
		 * Build URL dari state
		 */
		build: (baseUrl: URL, state: T): URL => {
			return buildUrlWithParams(baseUrl, state, {
				skipEmpty: true,
				skipDefaults: true,
				defaults: config.defaults
			});
		},

		/**
		 * Validate state
		 */
		validate: (state: Partial<T>): T => {
			const result: any = { ...config.defaults };

			Object.keys(state).forEach((key) => {
				const value = state[key as keyof T];
				if (config.validators && config.validators[key as keyof T]) {
					result[key] = config.validators[key as keyof T]!(value);
				} else {
					result[key] = value;
				}
			});

			return result as T;
		}
	};
}
export function createQueryManager() {
	const end = new Date();
	const start = new Date();
	start.setDate(end.getDate() - 30);
	return createUrlStateManager({
		defaults: {
			page: 1,
			limit: 10,
			search: '',
			type: 'ALL',
			status: 'ALL',
			sort_by: 'created_at',
			order_by: 'desc' as 'asc' | 'desc',
			date_from: start.toISOString(),
			date_to: end.toISOString()
		},
		validators: {
			page: (v) => Math.max(1, Number(v) || 1),
			limit: (v) => Math.min(Math.max(1, Number(v) || 10), 100),
			order_by: (v) => (v === 'asc' ? 'asc' : 'desc')
		}
	});
}
export function createVideosQueryManager() {
	const end = new Date();
	const start = new Date();
	start.setDate(end.getDate() - 30);
	return createUrlStateManager({
		defaults: {
			page: 1,
			limit: 10,
			search: '',
			status: 'all',
			tags: '',
			categori: '',
			actor: '',
			embed: '',
			sources: 'all',
			release_date: '',
			sort_by: 'created_at',
			order_by: 'desc' as 'asc' | 'desc',
			date_from: start.toISOString(),
			date_to: end.toISOString()
		},
		validators: {
			page: (v) => Math.max(1, Number(v) || 1),
			limit: (v) => Math.min(Math.max(1, Number(v) || 10), 100),
			order_by: (v) => (v === 'asc' ? 'asc' : 'desc')
		}
	});
}
