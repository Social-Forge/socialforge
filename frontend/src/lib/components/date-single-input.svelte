<script lang="ts">
	import type { ClassValue } from 'svelte/elements';
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
	import { cn } from '$lib/utils.js';
	import { buttonVariants } from '$lib/components/ui/button/index.js';
	import { Calendar } from '$lib/components/ui/calendar/index.js';
	import * as Popover from '$lib/components/ui/popover/index.js';
	import { formatToPostgresTimestampV2, parsePostgresTimestampV2 } from '$lib/utils/time.js';
	import { CalendarIcon } from '@lucide/svelte';

	let {
		modelValue = $bindable(),
		onchange,
		name,
		disabled,
		placeholder = 'Pick a date',
		class: className
	}: {
		modelValue?: string | Date | null;
		onchange?: (value: string | null) => void;
		name?: string;
		disabled?: boolean;
		placeholder?: string;
		class?: ClassValue;
	} = $props();

	let contentRef = $state<HTMLElement | null>(null);
	let isOpen = $state(false);
	const df = new DateFormatter('en-US', {
		dateStyle: 'long'
	});
	const timezone = getLocalTimeZone();

	const createCalendarDateFromDate = (date: Date): CalendarDate => {
		return new CalendarDate(date.getFullYear(), date.getMonth() + 1, date.getDate());
	};

	const dateToCalendarDate = (date: Date): ZonedDateTime => {
		return fromDate(date, timezone);
	};

	const dateValueToDate = (dateValue: DateValue): Date => {
		if (dateValue instanceof CalendarDate) {
			// Return Date dengan waktu set ke tengah malam UTC
			return new Date(Date.UTC(dateValue.year, dateValue.month - 1, dateValue.day));
		}

		try {
			const calendarDateTime = toCalendarDateTime(dateValue);
			const zonedDateTime = toZoned(calendarDateTime, timezone);
			return zonedDateTime.toDate();
		} catch {
			return new Date();
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

	// PERBAIKAN: Format untuk PostgreSQL timestamptz dengan waktu start of day
	const formatForPostgres = (date: Date): string => {
		// Set waktu ke start of day (00:00:00) untuk single date picker
		const startOfDay = new SvelteDate(date);
		startOfDay.setHours(0, 0, 0, 0);
		return formatToPostgresTimestampV2(startOfDay);
	};

	let value = $state<DateValue | undefined>();

	$effect(() => {
		if (!modelValue) {
			value = undefined;
		} else {
			let dateValue: ZonedDateTime | null = null;

			if (typeof modelValue === 'string') {
				dateValue = parsePostgresToCalendarDate(modelValue);
			} else if (modelValue instanceof Date) {
				dateValue = dateToCalendarDate(modelValue);
			}

			value = dateValue || undefined;
		}
	});

	const handleDateChange = (date: DateValue | undefined) => {
		if (!date) {
			modelValue = null;
			onchange?.(null);
			isOpen = false; // Tutup popover setelah clear
			return;
		}

		try {
			const jsDate = dateValueToDate(date);
			const formattedDate = formatForPostgres(jsDate); // Gunakan format PostgreSQL yang konsisten

			modelValue = formattedDate;
			onchange?.(formattedDate);

			// Close popover setelah select
			setTimeout(() => {
				isOpen = false;
			}, 300);
		} catch (error) {
			console.error('Error converting date:', error);
			modelValue = null;
			onchange?.(null);
		}
	};

	const displayText = $derived.by(() => {
		if (!value) return placeholder;

		try {
			const jsDate = dateValueToDate(value);
			return df.format(jsDate);
		} catch {
			return placeholder;
		}
	});

	const hasValue = $derived(!!value);
</script>

{#if name && modelValue}
	<input type="hidden" {name} value={modelValue} />
{/if}
<Popover.Root bind:open={isOpen}>
	<Popover.Trigger
		class={cn(
			buttonVariants({
				variant: 'outline'
			}),
			'w-full justify-start bg-white text-left font-normal text-black dark:bg-neutral-700 dark:text-white',
			!hasValue && 'text-muted-foreground',
			className
		)}
		{disabled}
	>
		<CalendarIcon class="mr-2 h-4 w-4" />
		{displayText}
	</Popover.Trigger>
	<Popover.Content bind:ref={contentRef} class="w-auto p-0" align="start">
		<Calendar
			type="single"
			bind:value
			weekdayFormat="short"
			numberOfMonths={1}
			{disabled}
			onValueChange={handleDateChange}
		/>
	</Popover.Content>
</Popover.Root>
