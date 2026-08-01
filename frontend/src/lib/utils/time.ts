const translations: Record<string, Record<string, string>> = {
	en: {
		year: 'year',
		month: 'month',
		week: 'week',
		day: 'day',
		hour: 'hour',
		minute: 'minute',
		second: 'second',
		ago: 'ago',
		just_now: 'just now'
	},
	id: {
		year: 'tahun',
		month: 'bulan',
		week: 'minggu',
		day: 'hari',
		hour: 'jam',
		minute: 'menit',
		second: 'detik',
		ago: 'yang lalu',
		just_now: 'baru saja'
	}
	// Tambahkan bahasa lain sesuai kebutuhan
};
export function formatTime(input: string | Date) {
	const date = typeof input === 'string' ? new Date(input) : input;
	return date.toLocaleTimeString('en-US', { hour: 'numeric', minute: '2-digit' });
}

export const formatTimeAgo = (dateString?: string, lang: string = 'en') => {
	if (!dateString) return '';
	try {
		const now = new Date();
		const publishedDate = new Date(dateString);
		// Validasi tanggal
		if (isNaN(publishedDate.getTime())) {
			console.warn('Invalid date string:', dateString);
			return '';
		}

		const diffInSeconds = Math.floor((now.getTime() - publishedDate.getTime()) / 1000);
		// Tetapkan satuan waktu yang sesuai
		const units = {
			year: 31536000,
			month: 2592000,
			week: 604800,
			day: 86400,
			hour: 3600,
			minute: 60,
			second: 1
		};
		// Cari satuan waktu terbesar yang cocok
		for (const [unit, seconds] of Object.entries(units)) {
			const interval = Math.floor(diffInSeconds / seconds);
			if (interval >= 1) {
				// Tambahkan support multi-bahasa
				// Pilih bahasa default jika tidak tersedia
				const t = translations[lang] || translations.en;
				// Handle pluralization
				const unitStr = interval === 1 ? t?.[unit] : `${t?.[unit]}${lang === 'en' ? 's' : ''}`;
				return interval <= 5 && unit === 'second'
					? t?.just_now
					: `${interval} ${unitStr} ${t?.ago}`;
			}
		}

		return translations[lang]?.just_now || 'just now';
	} catch (error) {
		console.error('Error formatting time ago:', error);
		return '';
	}
};
export const formatDateString = (dateString: string) => {
	if (!dateString) return 'Not set';
	return new Date(dateString).toLocaleDateString('en-US', {
		year: 'numeric',
		month: 'long',
		day: 'numeric',
		hour: '2-digit',
		minute: '2-digit'
	});
};
export const formatDate = (date: Date) => {
	return new Intl.DateTimeFormat('id-ID', {
		weekday: 'short', // "Rab"
		day: '2-digit', // "24"
		month: 'short', // "Sep"
		year: 'numeric', // "2025"
		hour: '2-digit', // "08"
		minute: '2-digit', // "54"
		second: '2-digit', // "27"
		hour12: false // 24 jam format
	})
		.format(date)
		.replace(/,/g, '') // Hapus koma
		.replace(/:/g, '.') // Ganti : dengan .
		.replace(/\s+/g, ' '); // Normalize spaces
};
export function formatToYYYYMMDD(date: Date): string {
	// Mengambil tahun (UTC)
	const year = date.getUTCFullYear();

	// Mengambil bulan (UTC), ditambah 1 karena getUTCMonth() mengembalikan 0-11
	const month = date.getUTCMonth() + 1;

	// Mengambil tanggal (UTC)
	const day = date.getUTCDate();

	// Memastikan format menjadi dua digit (misal: 1 -> 01)
	const formattedMonth = String(month).padStart(2, '0');
	const formattedDay = String(day).padStart(2, '0');

	return `${year}-${formattedMonth}-${formattedDay}`;
}
export function formatToPostgresTimestamp(date: Date): string {
	// Gunakan UTC time, bukan local time
	const year = date.getUTCFullYear();
	const month = String(date.getUTCMonth() + 1).padStart(2, '0');
	const day = String(date.getUTCDate()).padStart(2, '0');
	const hours = String(date.getUTCHours()).padStart(2, '0');
	const minutes = String(date.getUTCMinutes()).padStart(2, '0');
	const seconds = String(date.getUTCSeconds()).padStart(2, '0');

	// Milliseconds dengan 6 digit seperti PostgreSQL
	const milliseconds = String(date.getUTCMilliseconds()).padStart(3, '0') + '000';

	// PostgreSQL expects +00 for UTC, not local timezone offset
	return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}.${milliseconds}+00`;
}
export function formatToPostgresTimestampV2(date: Date): string {
	// Gunakan UTC time
	const year = date.getUTCFullYear();
	const month = String(date.getUTCMonth() + 1).padStart(2, '0');
	const day = String(date.getUTCDate()).padStart(2, '0');
	const hours = String(date.getUTCHours()).padStart(2, '0');
	const minutes = String(date.getUTCMinutes()).padStart(2, '0');
	const seconds = String(date.getUTCSeconds()).padStart(2, '0');

	// Milliseconds dengan 6 digit seperti PostgreSQL
	const milliseconds = String(date.getUTCMilliseconds()).padStart(3, '0') + '000';

	// PostgreSQL expects +00 for UTC
	return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}.${milliseconds}+00`;
}
export function parsePostgresTimestamp(timestamp: string): Date | null {
	try {
		// Handle both formats: '2025-10-12 05:37:30.323317+00' dan ISO string
		if (timestamp.includes('T')) {
			// ISO string format
			return new Date(timestamp);
		} else {
			// PostgreSQL timestamp format - convert to ISO
			// Remove the timezone part and add 'Z' for UTC
			const withoutTimezone = timestamp.split('+')[0].split('-')[0];
			// const withoutTimezone = timestamp.split('+')[0];
			const isoString = withoutTimezone.replace(' ', 'T') + 'Z';
			return new Date(isoString);
		}
	} catch {
		return null;
	}
}
export function parsePostgresTimestampV2(timestamp: string): Date | null {
	try {
		if (timestamp.includes('T')) {
			// ISO string format
			return new Date(timestamp);
		} else {
			// PostgreSQL timestamp format - convert to ISO dengan UTC
			// Handle format: '2025-10-12 05:37:30.323317+00'
			const withoutTimezone = timestamp.split('+')[0].trim();
			const isoString = withoutTimezone.replace(' ', 'T') + 'Z';
			const date = new Date(isoString);

			// Validasi date
			if (isNaN(date.getTime())) {
				console.error('❌ Invalid date parsed:', timestamp);
				return null;
			}

			return date;
		}
	} catch (error) {
		console.error('❌ Error parsing PostgreSQL timestamp:', error);
		return null;
	}
}
export function formatDateRangeForPostgres(
	startDate: Date,
	endDate: Date
): { start: string; end: string } {
	return {
		start: formatToPostgresTimestamp(startDate),
		end: formatToPostgresTimestamp(endDate)
	};
}
export function startOfDay(date: Date): Date {
	const newDate = new Date(date);
	newDate.setUTCHours(0, 0, 0, 0);
	return newDate;
}
export function endOfDay(date: Date): Date {
	const newDate = new Date(date);
	newDate.setUTCHours(23, 59, 59, 999);
	return newDate;
}
export function localToUTCDate(date: Date): Date {
	return new Date(
		Date.UTC(
			date.getFullYear(),
			date.getMonth(),
			date.getDate(),
			date.getHours(),
			date.getMinutes(),
			date.getSeconds(),
			date.getMilliseconds()
		)
	);
}
export function convertDayToDateFormatYYYYMMDD(dayOfWeek: string): string | null {
	const days = ['Sunday', 'Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday'];
	const dayIndex = days.indexOf(dayOfWeek);

	if (dayIndex === -1) {
		console.error(`Invalid day of week: ${dayOfWeek}`);
		return null; // Or throw an error
	}

	const today = new Date();
	const currentDayIndex = today.getDay(); // 0 for Sunday, 1 for Monday, etc.

	// Calculate the difference in days to reach the target day of the week
	let diff = dayIndex - currentDayIndex;

	// If the target day is earlier in the week, adjust to the next occurrence
	if (diff < 0) {
		diff += 7;
	}

	today.setDate(today.getDate() + diff);

	// Format the date to YYYY-MM-DD
	const year = today.getFullYear();
	const month = (today.getMonth() + 1).toString().padStart(2, '0'); // Months are 0-indexed
	const day = today.getDate().toString().padStart(2, '0');

	return `${year}-${month}-${day}`;
}
export function convertDayStringToDateYYYYMMDD(
	dayString: string,
	targetDate: string
): string | null {
	const daysOfWeek: { [key: string]: number } = {
		sunday: 0,
		monday: 1,
		tuesday: 2,
		wednesday: 3,
		thursday: 4,
		friday: 5,
		saturday: 6
	};

	const dayIndex = daysOfWeek[dayString.toLowerCase()];

	if (dayIndex === undefined) {
		console.error('Invalid day string provided.');
		return null;
	}

	const [year, month, day] = targetDate.split('-').map(Number);
	// Month is 0-indexed in JavaScript Date objects
	const date = new Date(year, month - 1, day);

	// Set the day of the week
	date.setDate(date.getDate() - date.getDay() + dayIndex);

	// Format the date back to YYYY-MM-DD
	const formattedYear = date.getFullYear();
	const formattedMonth = (date.getMonth() + 1).toString().padStart(2, '0'); // Add 1 as month is 0-indexed
	const formattedDay = date.getDate().toString().padStart(2, '0');

	return `${formattedYear}-${formattedMonth}-${formattedDay}`;
}
export function getTodayFormatYYYYMMDD() {
	const now: Date = new Date();

	const year: number = now.getFullYear();
	const month: string = (now.getMonth() + 1).toString().padStart(2, '0'); // Months are 0-indexed, so add 1
	const day: string = now.getDate().toString().padStart(2, '0');

	return `${year}-${month}-${day}`;
}
export function getDateForDayYYYYMMDD(dayName: string, weekOffset: number = 0): string {
	const days = ['Minggu', 'Senin', 'Selasa', 'Rabu', 'Kamis', 'Jumat', 'Sabtu'];
	const dayIndex = days.indexOf(dayName);

	if (dayIndex === -1) {
		console.error(`Invalid day name: ${dayName}`);
		// Fallback ke hari ini
		const today = new Date();
		return today.toISOString().split('T')[0];
	}

	const today = new Date();
	const currentDayIndex = today.getDay(); // 0 for Sunday, 1 for Monday, etc.

	// Calculate the target date
	const targetDate = new Date(today);

	// Set ke awal minggu (Minggu)
	targetDate.setDate(today.getDate() - currentDayIndex);

	// Tambah offset minggu jika diperlukan
	if (weekOffset !== 0) {
		targetDate.setDate(targetDate.getDate() + weekOffset * 7);
	}

	// Tambah hari sesuai index
	targetDate.setDate(targetDate.getDate() + dayIndex);

	// Format to YYYY-MM-DD
	const year = targetDate.getFullYear();
	const month = String(targetDate.getMonth() + 1).padStart(2, '0');
	const day = String(targetDate.getDate()).padStart(2, '0');

	return `${year}-${month}-${day}`;
}
export function safeCreateDate(dateString: string): Date {
	try {
		const date = new Date(dateString);
		if (isNaN(date.getTime())) {
			console.warn(`Invalid date string: ${dateString}, using current date`);
			return new Date();
		}
		return date;
	} catch (error) {
		console.error(`Error creating date from: ${dateString}`, error);
		return new Date();
	}
}
export function formatSecondsToHMS(totalSeconds: number): string {
	const hours = Math.floor(totalSeconds / 3600);
	const minutes = Math.floor((totalSeconds % 3600) / 60);
	const seconds = totalSeconds % 60;

	const pad = (num: number): string => String(num).padStart(2, '0');

	return `${pad(hours)}:${pad(minutes)}:${pad(seconds)}`;
}
export function parseSecondToTimeComponents(totalSeconds: number): TimeComponents {
	const hours = Math.floor(totalSeconds / 3600);
	const minutes = Math.floor((totalSeconds % 3600) / 60);
	const seconds = totalSeconds % 60;

	// new Date(s * 1000).toISOString().substring(11, 19);
	return { hours, minutes, seconds };
}
