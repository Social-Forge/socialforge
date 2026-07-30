import type { RequestEvent } from '@sveltejs/kit';
import { BaseHandler } from './base';
import type {
	UpdateProfileInput,
	UpdatePasswordInput,
	ActivatedTwoFactorInput
} from '$lib/utils/validators';

export class UserHandler extends BaseHandler {
	constructor(protected readonly event: RequestEvent) {
		super(event);
	}
	async currentUser() {
		try {
			const response = await this.api.authRequest<UserResponse>('GET', '/user/protected/me');
			if (!response.success) {
				return null;
			}
			return response.data;
		} catch (error) {
			return null;
		}
	}
	async logout() {
		return await this.api.authRequest('POST', '/user/protected/logout');
	}
	async uploadAvatar(file: File) {
		const formData = new FormData();
		formData.append('avatar', file);

		return await this.api.multipartAuthRequest<{ avatar_url: string }>(
			'POST',
			'/user/protected/avatar',
			formData
		);
	}
	async updateProfile(data: UpdateProfileInput) {
		return await this.api.authRequest<UserResponse>('PUT', '/user/protected/profile', data);
	}
	async changePassword(data: UpdatePasswordInput) {
		return await this.api.authRequest('PUT', '/user/protected/password', data);
	}
	async enableTwoFactor(status: string) {
		const body = { status };
		return await this.api.authRequest<{ qr_code?: string; secret?: string }>(
			'POST',
			'/user/protected/two-factor/enable',
			body
		);
	}
	async verifyTwoFactor(data: ActivatedTwoFactorInput) {
		return await this.api.authRequest('POST', '/user/protected/two-factor/verify', data);
	}
}
