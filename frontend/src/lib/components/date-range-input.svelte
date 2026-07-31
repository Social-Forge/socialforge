<script lang="ts">
	import { SvelteDate } from 'svelte/reactivity';
	import {
		DateFormatter,
		CalendarDate,
		type DateValue,
		type ZonedDateTime,
		getLocalTimeZone,
		toCalendarDateTime,
		toZoned,
		fromDate
	} from '@internationalized/date';
	import type { DateRange } from 'bits-ui';
	import { cn } from '$lib/utils.js';
	import { buttonVariants } from '$lib/components/ui/button/index.js';
	import { RangeCalendar } from './ui/range-calendar';
	import * as Popover from '$lib/components/ui/popover/index.js';
	import { formatToPostgresTimestampV2, parsePostgresTimestampV2 } from '$lib/utils/time.js';
	import { CalendarIcon } from '@lucide/svelte';

	let {
		modelValue = $bindable(),
		onchange,
		disabled
	}: {
		modelValue?: { start: string; end: string };
		onchange?: (value: { start: string; end: string } | null) => void;
		disabled?: boolean;
	} = $props();

	const df = new DateFormatter('en-US', {
		dateStyle: 'medium'
	});

	let contentRef = $state<HTMLElement | null>(null);
	let isOpen = $state(false);
	const timezone = getLocalTimeZone();

	const createCalendarDateFromDate = (date: Date): CalendarDate => {
		return new CalendarDate(date.getFullYear(), date.getMonth() + 1, date.getDate());
	};

	const dateToCalendarDate = (date: Date): ZonedDateTime => {
		return fromDate(date, timezone);
	};

	const dateValueToDate = (dateValue: DateValue): Date => {
		if (dateValue instanceof CalendarDate) {
			return new Date(dateValue.year, dateValue.month - 1, dateValue.day);
		}

		try {
			const calendarDateTime = toCalendarDateTime(dateValue);
			const zonedDateTime = toZoned(calendarDateTime, timezone);
			return zonedDateTime.toDate();
		} catch {
			return new Date();
		}
	};
	const stringToCalendarDate = (dateString: string): CalendarDate | null => {
		try {
			// Handle both PostgreSQL timestamp and ISO string
			const date = new Date(dateString);
			if (isNaN(date.getTime())) return null;
			return createCalendarDateFromDate(date);
		} catch {
			return null;
		}
	};
	const parsePostgresToCalendarDate = (timestamp: string): ZonedDateTime | null => {
		try {
			const date = parsePostgresTimestampV2(timestamp);
			if (!date || isNaN(date.getTime())) return null;

			// Gunakan function yang benar untuk convert Date ke CalendarDate
			return dateToCalendarDate(date);
		} catch {
			return null;
		}
	};
	const formatForPostgres = (date: Date, isEndDate: boolean = false): string => {
		if (isEndDate) {
			const endOfDay = new SvelteDate(date);
			endOfDay.setHours(23, 59, 59, 999);
			return formatToPostgresTimestampV2(endOfDay);
		} else {
			const startOfDay = new SvelteDate(date);
			startOfDay.setHours(0, 0, 0, 0);
			return formatToPostgresTimestampV2(startOfDay);
		}
	};

	const defaultStart = new CalendarDate(2024, 1, 1);
	const defaultEnd = defaultStart.add({ days: 7 });

	let value = $state<DateRange>({
		start: defaultStart,
		end: defaultEnd
	});

	$effect(() => {
		if (modelValue?.start && modelValue.end) {
			const startCal = parsePostgresToCalendarDate(modelValue.start);
			const endCal = parsePostgresToCalendarDate(modelValue.end);

			if (startCal && endCal) {
				value = {
					start: startCal,
					end: endCal
				};
			}

			// const startDate = parsePostgresTimestampV2(modelValue.start);
			// const endDate = parsePostgresTimestampV2(modelValue.end);

			// if (startDate && endDate) {
			// 	const startCal = createCalendarDateFromDate(startDate);
			// 	const endCal = createCalendarDateFromDate(endDate);

			// 	value = {
			// 		start: startCal,
			// 		end: endCal
			// 	};
			// }
		}
	});

	const handleDateChange = (range: DateRange) => {
		if (!range.start || !range.end) {
			return;
		}

		try {
			const startDate = dateValueToDate(range.start);
			const endDate = dateValueToDate(range.end);

			const result = {
				start: formatForPostgres(startDate, false), // Start of day
				end: formatForPostgres(endDate, true) // End of day
			};

			// Update internal value
			value = range;

			// Update modelValue for 2-way binding
			modelValue = result;

			// Call onchange callback if provided
			if (onchange) {
				onchange(result);
			}

			// Close popover setelah select
			setTimeout(() => {
				isOpen = false;
			}, 300);
		} catch (error) {
			console.error('Error converting date range:', error);
			if (onchange) {
				onchange(null);
			}
		}
	};

	const displayText = $derived.by(() => {
		if (!value?.start) return 'Pick a date range';

		try {
			const startDate = dateValueToDate(value.start);

			if (value.end) {
				const endDate = dateValueToDate(value.end);
				const display = `${df.format(startDate)} - ${df.format(endDate)}`;
				return display;
			} else {
				const display = df.format(startDate);
				return display;
			}
		} catch {
			const display = 'Pick a date range';
			return display;
		}
	});

	// Computed to check if we have a valid selection
	const hasValue = $derived(value?.start && value?.end);
</script>

<Popover.Root bind:open={isOpen}>
	<Popover.Trigger
		class={cn(
			buttonVariants({
				variant: 'outline'
			}),
			'w-full justify-start text-left font-normal md:w-75',
			!hasValue && 'text-muted-foreground'
		)}
		{disabled}
	>
		<CalendarIcon class="mr-2 h-4 w-4" />
		{displayText}
	</Popover.Trigger>
	<Popover.Content bind:ref={contentRef} class="w-auto p-0" align="start">
		<RangeCalendar
			bind:value
			class="rounded-md border"
			weekdayFormat="short"
			numberOfMonths={2}
			{disabled}
			onValueChange={handleDateChange}
		/>
	</Popover.Content>
</Popover.Root>
