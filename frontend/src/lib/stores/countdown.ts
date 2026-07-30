import { writable, derived } from 'svelte/store';
import {
	// getTimeRemaining,
	// getTimeUntilExpiry,
	// parsePostgresDate,
	// willExpireInDays,
	// isExpired,
	// getDaysDifference,
	formatDate
	// getRelativeTime
} from '$lib/utils/time';

// export function createCountdownStore(expiryDate: string) {
// 	const { subscribe, set, update } = writable<CountdownStore>({
// 		timeRemaining: getTimeRemaining(expiryDate),
// 		formattedTime: getTimeUntilExpiry(expiryDate),
// 		isRunning: true
// 	});

// 	let intervalId: NodeJS.Timeout;

// 	const start = () => {
// 		intervalId = setInterval(() => {
// 			update((store) => {
// 				const timeRemaining = getTimeRemaining(expiryDate);

// 				return {
// 					timeRemaining,
// 					formattedTime: getTimeUntilExpiry(expiryDate),
// 					isRunning: !timeRemaining.isExpired
// 				};
// 			});
// 		}, 1000);
// 	};

// 	const stop = () => {
// 		if (intervalId) {
// 			clearInterval(intervalId);
// 		}
// 		update((store) => ({ ...store, isRunning: false }));
// 	};

// 	// Auto-start
// 	start();

// 	return {
// 		subscribe,
// 		start,
// 		stop,
// 		destroy: stop
// 	};
// }
// export const createReactiveCountdown = (expiryDate: string) => {
// 	return derived(
// 		writable(expiryDate),
// 		($expiryDate, set) => {
// 			const update = () => {
// 				const timeRemaining = getTimeRemaining($expiryDate);
// 				set(timeRemaining);
// 			};

// 			update();
// 			const intervalId = setInterval(update, 1000);

// 			return () => clearInterval(intervalId);
// 		},
// 		getTimeRemaining(expiryDate)
// 	);
// };
// export const DateUtils = {
// 	/**
// 	 * Check if date string is valid PostgreSQL date format
// 	 */
// 	isValidPostgresDate(dateString: string): boolean {
// 		try {
// 			parsePostgresDate(dateString);
// 			return true;
// 		} catch {
// 			return false;
// 		}
// 	},

// 	/**
// 	 * Check if trial will expire soon (5 days or less)
// 	 */
// 	isTrialEndingSoon(expiryDate: string, warningDays: number = 5): boolean {
// 		return willExpireInDays(expiryDate, warningDays) && !isExpired(expiryDate);
// 	},

// 	/**
// 	 * Get Status Expiration with label
// 	 */
// 	getExpirationStatus(expiryDate: string): {
// 		status: 'active' | 'warning' | 'expired';
// 		message: string;
// 		daysLeft: number;
// 	} {
// 		if (isExpired(expiryDate)) {
// 			return {
// 				status: 'expired',
// 				message: 'Expired',
// 				daysLeft: 0
// 			};
// 		}

// 		const daysLeft = getDaysDifference(new Date().toISOString(), expiryDate);

// 		if (daysLeft <= 5) {
// 			return {
// 				status: 'warning',
// 				message: `Ends in ${daysLeft} days`,
// 				daysLeft
// 			};
// 		}

// 		return {
// 			status: 'active',
// 			message: `Ends in ${daysLeft} days`,
// 			daysLeft
// 		};
// 	},

// 	/**
// 	 * Get Status Expiration with label (3 days or less)
// 	 */
// 	getExpirationStatusThreeDayFromNow(expiryDate: string): {
// 		status: 'active' | 'warning' | 'expired';
// 		message: string;
// 		daysLeft: number;
// 	} {
// 		if (isExpired(expiryDate)) {
// 			return {
// 				status: 'expired',
// 				message: 'Expired',
// 				daysLeft: 0
// 			};
// 		}

// 		const daysLeft = getDaysDifference(new Date().toISOString(), expiryDate);

// 		if (daysLeft <= 3) {
// 			return {
// 				status: 'warning',
// 				message: `Ends in ${daysLeft} days`,
// 				daysLeft
// 			};
// 		}

// 		return {
// 			status: 'active',
// 			message: `Ends in ${daysLeft} days`,
// 			daysLeft
// 		};
// 	},

// 	/**
// 	 * Format date string for display UI
// 	 */
// 	formatForDisplay(dateString: string): {
// 		full: string;
// 		relative: string;
// 		short: string;
// 	} {
// 		return {
// 			full: formatDate(dateString, 'PPPP'),
// 			relative: getRelativeTime(dateString),
// 			short: formatDate(dateString, 'dd/MM/yyyy HH:mm')
// 		};
// 	}
// };
