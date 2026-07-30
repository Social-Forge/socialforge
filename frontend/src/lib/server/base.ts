import type { RequestEvent } from '@sveltejs/kit';
import slugify from 'slugify';

export class BaseHandler {
	protected event: RequestEvent;
	constructor(event: RequestEvent) {
		this.event = event;
	}
	protected get api() {
		return this.event.locals.api;
	}
	protected get helper() {
		return this.event.locals.helper;
	}
	protected get fetch() {
		return this.event.fetch;
	}
	protected handleError(error: unknown): Error {
		if (error instanceof Error) {
			throw error;
		}
		throw new Error('Internal server error');
	}
	protected safeSlug(value: unknown, fallback: string): string {
		const base = typeof value === 'string' && value.trim().length > 0 ? value : fallback;
		// only character, number and '-' is expected
		return slugify(base, {
			lower: true,
			remove: /[^\w\s-]/g,
			strict: true
		});
	}
	protected safeText(value: unknown, fallback: string): string {
		return typeof value === 'string' && value.trim().length > 0 ? value : fallback;
	}
	protected toISO(value?: string | Date | null): string | null {
		if (value === undefined || value === null) return null;
		if (typeof value === 'string') {
			const s = value.trim();
			if (!s || s === '?' || s.toLowerCase() === 'unknown' || s.toLowerCase() === 'n/a')
				return null;
			const d = new Date(s);
			if (isNaN(d.getTime())) return null;
			return d.toISOString();
		}
		const d = new Date(value);
		if (isNaN(d.getTime())) return null;
		return d.toISOString();
	}
	protected generateDateRange(start: string | Date, end: string | Date): Date[] {
		const startDate = new Date(start);
		const endDate = new Date(end);
		const dateRange: Date[] = [];

		const currentData = new Date(startDate);
		currentData.setHours(0, 0, 0, 0);
		endDate.setHours(0, 0, 0, 0);

		while (currentData <= endDate) {
			dateRange.push(new Date(currentData));
			currentData.setDate(currentData.getDate() + 1);
		}

		return dateRange;
	}
	protected generateDateRangeString(start: string | Date, end: string | Date): DateRangeChartLabel {
		const startDate = new Date(start);
		const endDate = new Date(end);
		const dates: Date[] = [];
		const labels: string[] = [];

		const currentDate = new Date(startDate);
		currentDate.setHours(0, 0, 0, 0);
		endDate.setHours(0, 0, 0, 0);

		while (currentDate <= endDate) {
			dates.push(new Date(currentDate));
			labels.push(this.formatChartLabel(currentDate));
			currentDate.setDate(currentDate.getDate() + 1);
		}

		return { dates, labels };
	}
	protected formatChartLabel(date: Date): string {
		const option: Intl.DateTimeFormatOptions = {
			weekday: 'short',
			day: 'numeric',
			month: 'short'
		};
		return date.toLocaleDateString('en-US', option);
	}
	protected formatDaysOnly(date: Date): string {
		const days = ['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat'];
		return days[date.getDay()];
	}
	protected formatFullDate(date: Date): string {
		const month = date.toLocaleString('en-US', { month: 'short' });
		const day = date.getDate();
		const year = date.getFullYear();
		return `${month} ${day}, ${year}`;
	}
}
