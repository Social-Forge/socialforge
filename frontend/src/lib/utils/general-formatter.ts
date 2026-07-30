/* eslint-disable no-useless-escape */
export function randomInt(min: number, max: number): number {
	return Math.floor(Math.random() * (max - min + 1)) + min;
}
export function randomFrom<T>(array: T[]): T {
	return array[Math.floor(Math.random() * array.length)]!;
}
export function stringToBooleanStrict(str: string): boolean {
	return str.toLowerCase() === 'true';
}
function replaceRandomChar(password: string, charSet: string): string {
	const randomIndex = Math.floor(Math.random() * password.length);
	const randomCharFromSet = charSet[Math.floor(Math.random() * charSet.length)];
	return (
		password.substring(0, randomIndex) + randomCharFromSet + password.substring(randomIndex + 1)
	);
}
export function generateRandomPassword(
	length: number = 12,
	includeUppercase: boolean = true,
	includeLowercase: boolean = true,
	includeNumbers: boolean = true,
	includeSymbols: boolean = true
): string {
	const uppercaseChars = 'ABCDEFGHIJKLMNOPQRSTUVWXYZ';
	const lowercaseChars = 'abcdefghijklmnopqrstuvwxyz';
	const numberChars = '0123456789';
	const symbolChars = '!@#$%^&*()-_=+[{]}\\|;:\'",<.>/?`~';

	let allChars = '';
	if (includeUppercase) allChars += uppercaseChars;
	if (includeLowercase) allChars += lowercaseChars;
	if (includeNumbers) allChars += numberChars;
	if (includeSymbols) allChars += symbolChars;

	if (allChars.length === 0) {
		throw new Error('At least one character type must be selected for password generation.');
	}

	let password = '';
	for (let i = 0; i < length; i++) {
		const randomIndex = Math.floor(Math.random() * allChars.length);
		password += allChars[randomIndex];
	}

	// Ensure at least one of each selected type is present (optional, for stronger passwords)
	if (includeUppercase && !/[A-Z]/.test(password)) {
		password = replaceRandomChar(password, uppercaseChars);
	}
	if (includeLowercase && !/[a-z]/.test(password)) {
		password = replaceRandomChar(password, lowercaseChars);
	}
	if (includeNumbers && !/[0-9]/.test(password)) {
		password = replaceRandomChar(password, numberChars);
	}
	if (includeSymbols && !/[!@#$%^&*()\-\_=+\[\]{}\\|;:'",<.>\/?`~]/.test(password)) {
		password = replaceRandomChar(password, symbolChars);
	}

	return password;
}

export function calculatePercentage(current: number, previous: number) {
	if (previous === 0) {
		return {
			value: current > 0 ? '+100%' : '0%',
			status: current > 0 ? 'New data' : 'No change',
			trend: current > 0 ? 'up' : 'neutral'
		};
	}

	const percentage = ((current - previous) / previous) * 100;
	const absPercentage = Math.abs(percentage);
	const value = `${percentage >= 0 ? '+' : ''}${absPercentage.toFixed(1)}%`;

	let status: string;
	let trend: 'up' | 'down' | 'neutral';

	if (percentage > 0) {
		status = `Increased ${absPercentage.toFixed(1)}% this period`;
		trend = 'up';
	} else if (percentage < 0) {
		status = `Decreased ${absPercentage.toFixed(1)}% this period`;
		trend = 'down';
	} else {
		status = 'No change this period';
		trend = 'neutral';
	}

	return { value, status, trend };
}
export function postgresArrayToStringArray(data: any): string[] {
	const cleanedString = data.slice(1, -1);
	const stringArray: string[] = cleanedString.split(', ');
	return stringArray;
}
export function compareVersions(a: string, b: string): number {
	const aParts = a.split('.').map(Number);
	const bParts = b.split('.').map(Number);

	for (let i = 0; i < Math.max(aParts.length, bParts.length); i++) {
		const aVal = aParts[i] || 0;
		const bVal = bParts[i] || 0;
		if (aVal > bVal) return 1;
		if (aVal < bVal) return -1;
	}
	return 0;
}
