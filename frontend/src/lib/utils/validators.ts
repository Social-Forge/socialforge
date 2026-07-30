import { z } from 'zod';

const RESERVED_TENANT_SUBDOMAINS = new Set(['www', 'member', 'admin', 'api', 'auth']);
const TENANT_SUBDOMAIN_PATTERN = /^[a-z0-9](?:[a-z0-9-]{1,61}[a-z0-9])?$/;
const MAX_FILE_SIZE = 5 * 1024 * 1024;
const ACCEPTED_IMAGE_TYPES = ['image/jpeg', 'image/png', 'image/webp'];
const imageUploadSchema = z
	.instanceof(File)
	.refine((file) => file.size <= MAX_FILE_SIZE, 'Max size 5MB.')
	.refine(
		(file) => ACCEPTED_IMAGE_TYPES.includes(file.type),
		'Only .jpg, .png, and .webp supported.'
	);
const singleImageSchema = z.instanceof(File).superRefine(async (file, ctx) => {
	if (!ACCEPTED_IMAGE_TYPES.includes(file.type)) {
		ctx.addIssue({ code: 'custom', message: 'Only .jpg, .png, and .webp supported.' });
		return;
	}
	if (file.size > MAX_FILE_SIZE) {
		ctx.addIssue({ code: 'custom', message: 'Max file size 5MB.' });
		return;
	}
	if (typeof window !== 'undefined') {
		const { width, height } = await validateDimensions(file);
		if (width > 1920 || height > 1080) {
			ctx.addIssue({ code: 'custom', message: 'Max image resolution 1920x1080.' });
		}
	}
});
const validateDimensions = (file: File): Promise<{ width: number; height: number }> => {
	return new Promise((resolve) => {
		if (typeof window === 'undefined') return resolve({ width: 0, height: 0 }); // Skip jika di Server
		const img = new Image();
		img.src = URL.createObjectURL(file);
		img.onload = () => {
			resolve({ width: img.width, height: img.height });
			URL.revokeObjectURL(img.src);
		};
		img.onerror = () => resolve({ width: 0, height: 0 });
	});
};
const fileSchema = z.instanceof(File).superRefine(async (file, ctx) => {
	if (!ACCEPTED_IMAGE_TYPES.includes(file.type)) {
		ctx.addIssue({ code: 'custom', message: 'Only .jpg, .png, and .webp supported.' });
		return;
	}

	if (file.size > MAX_FILE_SIZE) {
		ctx.addIssue({ code: 'custom', message: 'Max size 5MB.' });
		return;
	}

	const { width, height } = await validateDimensions(file);
	if (width > 1920 || height > 1080) {
		ctx.addIssue({ code: 'custom', message: 'Max 1920x1080' });
	}
});

export const loginSchema = z.object({
	email: z.string().email('Invalid email').min(1),
	password: z.string().min(1)
});
export const registerSchema = z
	.object({
		name: z.string().min(1),
		email: z.string().email('Invalid email').min(1),
		subdomain: z
			.string()
			.trim()
			.toLowerCase()
			.min(3, 'Subdomain must be at least 3 characters')
			.max(63, 'Subdomain must be at most 63 characters')
			.regex(TENANT_SUBDOMAIN_PATTERN, 'Subdomain can only contain letters, numbers, and hyphens')
			.refine((value) => !RESERVED_TENANT_SUBDOMAINS.has(value), {
				message: 'Subdomain is reserved'
			}),
		password: z.string().min(1),
		confirmPassword: z.string().min(1)
	})
	.refine((data) => data.password === data.confirmPassword, {
		message: 'Passwords do not match',
		path: ['confirmPassword']
	});
export const forgotPasswordSchema = z.object({
	email: z.string().email('Invalid email').min(1)
});
export const resetPasswordSchema = z
	.object({
		token: z.string().min(1),
		password: z.string().min(1),
		confirmPassword: z.string().min(1)
	})
	.refine((data) => data.password === data.confirmPassword, {
		message: 'Passwords do not match',
		path: ['confirmPassword']
	});

export type LoginInput = z.infer<typeof loginSchema>;
export type RegisterInput = z.infer<typeof registerSchema>;
export type ForgotPasswordInput = z.infer<typeof forgotPasswordSchema>;
export type ResetPasswordInput = z.infer<typeof resetPasswordSchema>;

export const platformSettingsSchema = z.object({
	publisher_share_percent: z.coerce.number().min(0).max(100),
	referral_percent: z.coerce.number().min(0).max(50),
	min_payout_usd: z.coerce.number().min(1).max(1000),
	default_countdown_sec: z.coerce.number().int().min(0).max(120),
	guest_links_per_hour: z.coerce.number().int().min(1).max(100),
	enable_register: z.boolean().default(true),
	default_placements: z
		.array(
			z.enum([
				'popunder',
				'banner_header',
				'banner_footer',
				'banner_sidebar',
				'native_in_content',
				'smartlink_button',
				'social_bar'
			])
		)
		.max(7)
});
export const updatePasswordSchema = z
	.object({
		currentPassword: z
			.string({ error: 'Current password is required' })
			.min(1, { message: 'Current password is required' })
			.min(6, { message: 'Current password must be at least 6 characters long' })
			.transform((value) => value.replaceAll(/\s+/g, '')),
		newPassword: z
			.string({ error: 'New password is required' })
			.min(6, { message: 'New password must be at least 6 characters long' })
			.transform((value) => value.replaceAll(/\s+/g, '')),
		confirmPassword: z
			.string({ error: 'Confirm password is required' })
			.nonempty({ message: 'Confirm password is required' })
			.transform((value) => value.replaceAll(/\s+/g, ''))
	})
	.refine((data) => data.newPassword === data.confirmPassword, {
		message: 'Passwords do not match',
		path: ['confirmPassword']
	});
export const updateProfileSchema = z.object({
	name: z.string().min(1),
	email: z.string().email('Invalid email').min(1)
});

export type PlatformSettingsInput = z.infer<typeof platformSettingsSchema>;
export type Placement = PlatformSettingsInput['default_placements'][number];
export type UpdatePasswordInput = z.infer<typeof updatePasswordSchema>;
export type UpdateProfileInput = z.infer<typeof updateProfileSchema>;
