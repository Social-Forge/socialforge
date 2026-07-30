import type { RequestEvent } from '@sveltejs/kit';
import { BaseHandler } from './base';
import type {
	RegisterInput,
	LoginInput,
	ForgotPasswordInput,
	ResetPasswordInput,
	VerifyTwoFactorInput
} from '$lib/utils/validators';

export class AuthHandler extends BaseHandler {
	constructor(protected readonly event: RequestEvent) {
		super(event);
	}
	async register(value: RegisterInput) {
		return await this.api.publicRequest('POST', '/auth/register', value);
	}
	async login(value: LoginInput) {
		return await this.api.publicRequest<LoginResponse>('POST', '/auth/login', value);
	}
	async forgot(value: ForgotPasswordInput) {
		return await this.api.publicRequest('POST', '/auth/forgot', value);
	}
	async verifyEmail(token: string) {
		return await this.api.publicRequest('POST', '/auth/verify-email', { token });
	}
	async resetPassword(value: ResetPasswordInput) {
		return await this.api.publicRequest('POST', '/auth/reset-password', value);
	}
	async oAuthRedirect(provider: string) {
		return await this.api.publicRequest('GET', `/auth/oauth/${provider}`);
	}
	async oAuthCallback(provider: string) {
		return await this.api.publicRequest<LoginResponse>('GET', `/auth/oauth/${provider}/callback`);
	}
	async refreshToken(refreshToken: string) {
		return await this.api.publicRequest<LoginResponse>('POST', '/auth/refresh-token', {
			refresh_token: refreshToken
		});
	}
	async verifyTwoFactor(value: VerifyTwoFactorInput) {
		return await this.api.publicRequest<LoginResponse>('POST', '/auth/verify-two-factor', value);
	}
}
