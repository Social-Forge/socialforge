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
	password: z
		.string({ error: 'Password is required' })
		.min(1, { message: 'Password is required' })
		.min(6, { message: 'Password must be at least 6 characters long' })
		.transform((value) => value.replaceAll(/\s+/g, '')),
	remember_me: z.boolean().optional().default(false)
});
export const registerSchema = z
	.object({
		full_name: z
			.string({ error: 'First name is required' })
			.min(3, 'First name must be at least 3 characters long')
			.nonempty('First name is required'),
		email: z
			.string({ error: 'Email is required' })
			.email('Email is not valid')
			.nonempty('Email is required'),
		phone: z
			.string()
			.regex(
				/^\+\d{1,4}[\d\s-]{6,15}$/,
				'Phone must start with country code (e.g., +1) and contain only numbers, spaces, or dashes'
			)
			.transform((val) => (val ? val.replace(/[\s-]/g, '') : ''))
			.refine((val) => !val || /^\+\d{1,4}\d{4,13}$/.test(val), {
				message: 'Phone must start with country code and contain only numbers after cleaning'
			})
			.refine((val) => !val || val.length >= 8, {
				message: 'Phone must be at least 8 characters long'
			})
			.refine((val) => !val || val.length <= 16, {
				message: 'Phone must be at most 16 characters long'
			})
			.optional()
			.or(z.literal('')),

		password: z
			.string({ error: 'Password is required' })
			.min(1, { message: 'Password is required' })
			.min(6, { message: 'Password must be at least 6 characters long' })
			.regex(/[A-Z]/, {
				message: 'Password must contain at least one uppercase letter'
			})
			.regex(/[0-9]/, { message: 'Password must contain at least one number' })
			.transform((value) => value.replaceAll(/\s+/g, '')),

		confirm_password: z
			.string({ error: 'Confirm password is required' })
			.nonempty({ message: 'Confirm password is required' })
			.transform((value) => value.replaceAll(/\s+/g, ''))
	})
	.superRefine((data, ctx) => {
		if (data.password !== data.confirm_password) {
			ctx.addIssue({
				path: ['confirm_password'],
				code: z.ZodIssueCode.custom,
				message: 'Password and confirm password must be the same'
			});
		}
	});
export const forgotPasswordSchema = z.object({
	email: z.string().email('Invalid email').min(1)
});
export const resetPasswordSchema = z
	.object({
		new_password: z
			.string()
			.min(6, 'Password must be at least 6 characters')
			.transform((value) => value.replaceAll(/\s+/g, '')),
		confirm_password: z
			.string()
			.nonempty('Confirm password is required')
			.transform((value) => value.replaceAll(/\s+/g, '')),
		token: z.string().nonempty('Token is required')
	})
	.superRefine((data, ctx) => {
		if (data.new_password != data.confirm_password) {
			ctx.addIssue({
				path: ['confirm_password'],
				code: z.ZodIssueCode.custom,
				message: 'Password and confirm password must be the same'
			});
		}
	});
export const verifyTwoFactorSchema = z.object({
	token: z.string().nonempty('Two factor authentication token is required'),
	otp: z.string().nonempty('One time password is required')
});
export const updateProfileSchema = z.object({
	full_name: z
		.string({ error: 'Full name is required' })
		.min(3, 'Full name must be at least 3 characters long')
		.nonempty('Full name is required'),
	email: z
		.string({ error: 'Email is required' })
		.email('Email is not valid')
		.nonempty('Email is required'),
	username: z
		.string({ error: 'Username is required' })
		.min(3, 'Username must be at least 3 characters long')
		.regex(/^[a-z0-9]+$/, 'Username must be lowercase letters and numbers only, without spaces')
		.nonempty('Username is required'),
	phone: z
		.string()
		.regex(
			/^\+\d{1,4}[\d\s-]{6,15}$/,
			'Phone must start with country code (e.g., +1) and contain only numbers, spaces, or dashes'
		)
		.transform((val) => (val ? val.replace(/[\s-]/g, '') : ''))
		.refine((val) => !val || /^\+\d{1,4}\d{4,13}$/.test(val), {
			message: 'Phone must start with country code and contain only numbers after cleaning'
		})
		.refine((val) => !val || val.length >= 8, {
			message: 'Phone must be at least 8 characters long'
		})
		.refine((val) => !val || val.length <= 16, {
			message: 'Phone must be at most 16 characters long'
		})
		.optional()
		.or(z.literal(''))
});
export const updateTenantSchema = z.object({
	id: z.string().nonempty('Tenant ID is required'),
	name: z
		.string({ error: 'Name is required' })
		.min(3, 'Name must be at least 3 characters long')
		.nonempty('Name is required'),
	slug: z
		.string({ error: 'Slug is required' })
		.min(3, 'Slug must be at least 3 characters long')
		.nonempty('Slug is required'),
	subdomain: z.string().nullable().optional(),
	description: z.string().nullable().optional()
});
export const updatePasswordSchema = z.object({
	current_password: z
		.string({ error: 'Current password is required' })
		.min(1, { message: 'Current password is required' })
		.min(6, { message: 'Current password must be at least 6 characters long' })
		.transform((value) => value.replaceAll(/\s+/g, '')),
	new_password: z
		.string({ error: 'New password is required' })
		.min(6, { message: 'New password must be at least 6 characters long' })
		.transform((value) => value.replaceAll(/\s+/g, '')),
	confirm_password: z
		.string({ error: 'Confirm password is required' })
		.nonempty({ message: 'Confirm password is required' })
		.transform((value) => value.replaceAll(/\s+/g, ''))
});
export const activatedTwoFactorSchema = z.object({
	code: z
		.string()
		.min(6, 'Two factor authentication code must be at least 6 characters long')
		.max(6, 'Two factor authentication code must be at most 6 characters long')
		.nonempty('Two factor authentication code is required')
});

export type LoginInput = z.infer<typeof loginSchema>;
export type RegisterInput = z.infer<typeof registerSchema>;
export type ForgotPasswordInput = z.infer<typeof forgotPasswordSchema>;
export type ResetPasswordInput = z.infer<typeof resetPasswordSchema>;
export type VerifyTwoFactorInput = z.infer<typeof verifyTwoFactorSchema>;
export type ActivatedTwoFactorInput = z.infer<typeof activatedTwoFactorSchema>;

export type UpdatePasswordInput = z.infer<typeof updatePasswordSchema>;
export type UpdateProfileInput = z.infer<typeof updateProfileSchema>;
