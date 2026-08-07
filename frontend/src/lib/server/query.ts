import { type RequestEvent } from '@sveltejs/kit';
import { BaseHandler } from './base';
import { formatToPostgresTimestampV2 } from '$lib/utils/time.js';

export class QueryHelper<T = any> extends BaseHandler {
	private whereConditions: any[] = [];

	constructor(event: RequestEvent) {
		super(event);
	}

	/**
	 * Add search conditions untuk berbagai tipe field
	 *
	 * @param search - Search term
	 * @param fields - Array of field configurations
	 *
	 * @example
	 * builder.addSearch('action', [
	 *   { field: 'name', type: 'string' },
	 *   { field: 'slug', type: 'string' },
	 *   { field: 'tmdb_id', type: 'number' }
	 * ]);
	 */
	addSearch(search: string | null, fields: SearchFieldConfig[]): this {
		if (!search || search.trim() === '') return this;

		const searchConditions: any[] = [];

		for (const config of fields) {
			// Skip if not searchable
			if (config.searchable === false) continue;

			const { field, type } = config;

			switch (type) {
				case 'string':
					// String fields: use contains with case-insensitive
					searchConditions.push({
						[field]: {
							contains: search,
							mode: 'insensitive'
						}
					});
					break;

				case 'number':
					// Number fields: convert search to number and use equals
					const searchAsNumber = parseInt(search);
					if (!isNaN(searchAsNumber)) {
						searchConditions.push({
							[field]: searchAsNumber
						});
					}
					break;

				case 'date':
					// Date fields: try to parse as date
					const searchAsDate = new Date(search);
					if (!isNaN(searchAsDate.getTime())) {
						// Search for same day
						const startOfDay = new Date(searchAsDate.setHours(0, 0, 0, 0));
						const endOfDay = new Date(searchAsDate.setHours(23, 59, 59, 999));

						searchConditions.push({
							[field]: {
								gte: startOfDay,
								lte: endOfDay
							}
						});
					}
					break;

				case 'boolean':
					// Boolean fields: convert search to boolean
					const searchLower = search.toLowerCase();
					if (searchLower === 'true' || searchLower === '1' || searchLower === 'yes') {
						searchConditions.push({ [field]: true });
					} else if (searchLower === 'false' || searchLower === '0' || searchLower === 'no') {
						searchConditions.push({ [field]: false });
					}
					break;
			}
		}

		// Add OR conditions if any
		if (searchConditions.length > 0) {
			this.whereConditions.push({
				OR: searchConditions
			});
		}

		return this;
	}

	/**
	 * Add filter condition
	 */
	addFilter(field: string, value: any, operator: 'equals' | 'in' | 'contains' = 'equals'): this {
		if (value === undefined || value === null) return this;

		switch (operator) {
			case 'equals':
				this.whereConditions.push({ [field]: value });
				break;
			case 'in':
				this.whereConditions.push({ [field]: { in: Array.isArray(value) ? value : [value] } });
				break;
			case 'contains':
				this.whereConditions.push({
					[field]: {
						contains: value,
						mode: 'insensitive'
					}
				});
				break;
		}

		return this;
	}

	/**
	 * Add date range filter
	 */
	addDateRange(field: string, from?: string | Date, to?: string | Date): this {
		if (!from && !to) return this;

		const condition: any = {};

		if (from) {
			condition.gte = new Date(from);
		}

		if (to) {
			condition.lte = new Date(to);
		}

		this.whereConditions.push({
			[field]: condition
		});

		return this;
	}

	/**
	 * Add number range filter
	 */
	addNumberRange(field: string, min?: number, max?: number): this {
		if (min === undefined && max === undefined) return this;

		const condition: any = {};

		if (min !== undefined) {
			condition.gte = min;
		}

		if (max !== undefined) {
			condition.lte = max;
		}

		this.whereConditions.push({
			[field]: condition
		});

		return this;
	}

	/**
	 * Add custom condition
	 */
	addCondition(condition: any): this {
		if (condition) {
			this.whereConditions.push(condition);
		}
		return this;
	}

	/**
	 * Add soft delete filter
	 */
	excludeDeleted(field: string = 'deleted_at'): this {
		this.whereConditions.push({ [field]: null });
		return this;
	}

	/**
	 * Build final WHERE clause
	 */
	build(): any {
		if (this.whereConditions.length === 0) {
			return {};
		}

		if (this.whereConditions.length === 1) {
			return this.whereConditions[0];
		}

		return { AND: this.whereConditions };
	}

	/**
	 * Build ORDER BY clause
	 */
	buildOrderBy(sortBy?: string, orderBy: 'asc' | 'desc' = 'desc'): any {
		if (!sortBy) {
			return { created_at: 'desc' };
		}

		return { [sortBy]: orderBy };
	}

	/**
	 * Build pagination params
	 */
	buildPagination(page: number = 1, limit: number = 10): { skip: number; take: number } {
		const skip = (page - 1) * limit;
		return { skip, take: limit };
	}

	/**
	 * Build paginated result
	 */
	buildPaginatedResult<T>(
		data: T[],
		totalCount: number,
		page: number,
		limit: number
	): {
		data: T[];
		pagination: {
			current_page: number;
			total_pages: number;
			total_items: number;
			has_next: boolean;
			has_prev: boolean;
			limit: number;
		};
	} {
		const totalPages = Math.ceil(totalCount / limit);

		return {
			data,
			pagination: {
				current_page: page,
				total_pages: totalPages,
				total_items: totalCount,
				has_next: page < totalPages,
				has_prev: page > 1,
				limit
			}
		};
	}

	parseQueryParams(url: URL): QueryParams {
		const params = url.searchParams;

		const defaultDateRange = this.getDefaultDateRange();

		const calculateParams = {
			page: parseInt(params.get('page') || '1'),
			limit: parseInt(params.get('limit') || '10'),
			search: params.get('search') || undefined,
			sort_by: params.get('sort_by') || 'created_at',
			order_by: (params.get('order_by') as 'asc' | 'desc') || 'desc',
			date_from: params.get('date_from') || defaultDateRange.start,
			date_to: params.get('date_to') || defaultDateRange.end,
			status: params.get('status') || undefined,
			type: params.get('type') || undefined,
			include_deleted: params.get('include_deleted') === 'true',
			with_relations: params.get('with_relations') === 'true',
			with_delete_column: params.get('with_delete_column') === 'true',
			is_active: params.get('is_active') === 'true',
			is_verified: params.get('is_verified') === 'true',
			user_id: params.get('user_id') || undefined,
			viewer_id: params.get('viewer_id') || undefined,
			includes: params.get('includes')?.split(',') || undefined,
			fields: params.get('fields')?.split(',') || undefined,
			genre: params.get('genre') || undefined,
			extra: {},
			// for comment
			parent_id: params.get('parent_id') || undefined,
			include_replies: params.get('include_replies') === 'true',
			depth: params.get('depth') ? parseInt(params.get('depth') || '0') : undefined,
			// For Blog Post
			tag: params.get('tag') || undefined,
			series: params.get('series') || undefined,
			year: params.get('year') ? parseInt(params.get('year') || '0') : undefined,
			month: params.get('month') ? parseInt(params.get('month') || '0') : undefined
		};

		// if (calculateParams.page < 1) calculateParams.page = 1;
		// if (calculateParams.limit < 1) calculateParams.limit = 10;
		// if (calculateParams.limit > 100) calculateParams.limit = 100;
		// if (calculateParams?.year < 0) calculateParams.year = undefined;
		// if (calculateParams?.month < 0 || calculateParams?.month > 12) calculateParams.month = undefined;

		return calculateParams;
	}

	getDefaultDateRange() {
		const end = new Date();
		const start = new Date();
		start.setDate(end.getDate() - 30);

		start.setHours(0, 0, 0, 0);
		end.setHours(23, 59, 59, 999);

		return {
			start: start.toISOString(),
			end: end.toISOString()
		};
	}

	getPrevoiusPeriod(dateFrom: string, dateTo: string): { start: string; end: string } {
		const start = new Date(dateFrom);
		const end = new Date(dateTo);

		const timeDiff = end.getTime() - start.getTime();

		const previoudDateStart = new Date(start.getTime() - timeDiff - 1);
		const previoudDateEnd = new Date(start.getTime() - 1);

		return {
			start: formatToPostgresTimestampV2(previoudDateStart),
			end: formatToPostgresTimestampV2(previoudDateEnd)
		};
	}

	clone(): QueryHelper {
		const clone = new QueryHelper(this.event);
		clone.whereConditions = [...this.whereConditions];
		return clone;
	}
}
