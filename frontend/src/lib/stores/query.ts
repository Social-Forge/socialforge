import { writable, derived } from 'svelte/store';
import { browser } from '$app/environment';

const createQueryStore = (init: QueryParams) => {
	const { subscribe, set, update } = writable<QueryState>({
		params: init,
		meta: null,
		isLoading: false
	});

	return {
		subscribe,
		setParams: (params: Partial<QueryParams>) => {
			update((state) => ({
				...state,
				params: { ...state.params, ...params, page: 1 }
			}));
		},
		setPage: (page: number) => {
			update((state) => ({
				...state,
				params: { ...state.params, page }
			}));
		},
		setMeta: (meta: ApiMeta) => {
			update((state) => ({ ...state, meta }));
		},
		setLoading: (isLoading: boolean) => {
			update((state) => ({ ...state, isLoading }));
		},
		reset: () => {
			set({
				params: init,
				meta: null,
				isLoading: false
			});
		},
		syncWithURL: (url: URL) => {
			if (!browser) return;

			const params: Partial<QueryParams> = {};
			// Basic pagination
			const page = url.searchParams.get('page');
			const limit = url.searchParams.get('limit');

			if (page) params.page = parseInt(page);
			if (limit) params.limit = parseInt(limit);

			// Filters
			const search = url.searchParams.get('search');
			const status = url.searchParams.get('status');
			const sort_by = url.searchParams.get('sort_by');
			const order_by = url.searchParams.get('order_by');

			if (search) params.search = search;
			if (status) params.status = status;
			if (sort_by) params.sort_by = sort_by;
			if (order_by) params.order_by = order_by as 'asc' | 'desc';

			// Boolean filters
			const include_deleted = url.searchParams.get('include_deleted');
			const is_active = url.searchParams.get('is_active');
			const is_verified = url.searchParams.get('is_verified');

			if (include_deleted) params.include_deleted = include_deleted === 'true';
			if (is_active) params.is_active = is_active === 'true';
			if (is_verified) params.is_verified = is_verified === 'true';

			// IDs
			const tenant_id = url.searchParams.get('tenant_id');
			const user_id = url.searchParams.get('user_id');
			const division_id = url.searchParams.get('division_id');

			if (tenant_id) params.tenant_id = tenant_id;
			if (user_id) params.user_id = user_id;
			if (division_id) params.division_id = division_id;

			update((state) => ({
				...state,
				params: { ...state.params, ...params }
			}));
		},
		updateURL: (url: URL) => {
			if (!browser) return;

			update((state) => {
				const newUrl = new URL(url);
				// Clear existing params
				newUrl.searchParams.delete('page');
				newUrl.searchParams.delete('limit');
				newUrl.searchParams.delete('search');
				newUrl.searchParams.delete('status');
				newUrl.searchParams.delete('sort_by');
				newUrl.searchParams.delete('order_by');
				newUrl.searchParams.delete('include_deleted');
				newUrl.searchParams.delete('is_active');
				newUrl.searchParams.delete('is_verified');
				newUrl.searchParams.delete('tenant_id');
				newUrl.searchParams.delete('user_id');
				newUrl.searchParams.delete('division_id');
				// Set new params
				const { params } = state;
				if (params.page > 1) newUrl.searchParams.set('page', params.page.toString());
				if (params.limit !== 10) newUrl.searchParams.set('limit', params.limit.toString());
				if (params.search) newUrl.searchParams.set('search', params.search);
				if (params.status) newUrl.searchParams.set('status', params.status);
				if (params.sort_by) newUrl.searchParams.set('sort_by', params.sort_by);
				if (params.order_by && params.order_by !== 'desc')
					newUrl.searchParams.set('order_by', params.order_by);

				if (params.include_deleted) newUrl.searchParams.set('include_deleted', 'true');
				if (params.is_active !== undefined)
					newUrl.searchParams.set('is_active', params.is_active.toString());
				if (params.is_verified !== undefined)
					newUrl.searchParams.set('is_verified', params.is_verified.toString());

				if (params.tenant_id) newUrl.searchParams.set('tenant_id', params.tenant_id);
				if (params.user_id) newUrl.searchParams.set('user_id', params.user_id);
				if (params.division_id) newUrl.searchParams.set('division_id', params.division_id);

				window.history.replaceState(null, '', newUrl.toString());

				return state;
			});
		}
	};
};
export const createQueryDerivedStores = (queryStore: ReturnType<typeof createQueryStore>) => {
	const currentPage = derived(queryStore, ($store) => $store.params.page);
	const totalPages = derived(queryStore, ($store) => $store.meta?.total_pages || 0);
	const hasNext = derived(queryStore, ($store) => $store.meta?.has_next || false);
	const hasPrev = derived(queryStore, ($store) => $store.meta?.has_prev || false);
	const isLoading = derived(queryStore, ($store) => $store.isLoading);

	return {
		currentPage,
		totalPages,
		hasNext,
		hasPrev,
		isLoading
	};
};

export default createQueryDerivedStores;
