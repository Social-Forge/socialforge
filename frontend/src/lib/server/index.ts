import type { RequestEvent } from '@sveltejs/kit';
import { BaseHandler } from './base';
import { SessionHelper } from './session';
import { AuthHandler } from './auth';
import { UserHandler } from './user';

export class ServiceHelper extends BaseHandler {
	public readonly session: InstanceType<typeof SessionHelper>;
	public readonly auth: InstanceType<typeof AuthHandler>;
	public readonly user: InstanceType<typeof UserHandler>;
	constructor(event: RequestEvent) {
		super(event);
		this.session = new SessionHelper(event);
		this.auth = new AuthHandler(event);
		this.user = new UserHandler(event);
	}
}
