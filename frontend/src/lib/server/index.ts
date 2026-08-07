import type { RequestEvent } from '@sveltejs/kit';
import { BaseHandler } from './base';
import { SessionHelper } from './session';
import { AuthHandler } from './auth';
import { UserHandler } from './user';
import { QueryHelper } from './query';
import { ConversationHandler } from './conversation';
import { ChannelHandler } from './channel';
import { BillingHandler } from './billing';
import { AIAgentHandler } from './ai-agent';
import { ContactHandler } from './contact';

export class ServiceHelper extends BaseHandler {
	public readonly session: InstanceType<typeof SessionHelper>;
	public readonly auth: InstanceType<typeof AuthHandler>;
	public readonly user: InstanceType<typeof UserHandler>;
	public readonly conversation: InstanceType<typeof ConversationHandler>;
	public readonly channel: InstanceType<typeof ChannelHandler>;
	public readonly billing: InstanceType<typeof BillingHandler>;
	public readonly aiAgent: InstanceType<typeof AIAgentHandler>;
	public readonly contact: InstanceType<typeof ContactHandler>;
	public readonly query: InstanceType<typeof QueryHelper>;
	constructor(event: RequestEvent) {
		super(event);
		this.session = new SessionHelper(event);
		this.auth = new AuthHandler(event);
		this.user = new UserHandler(event);
		this.conversation = new ConversationHandler(event);
		this.channel = new ChannelHandler(event);
		this.billing = new BillingHandler(event);
		this.aiAgent = new AIAgentHandler(event);
		this.contact = new ContactHandler(event);
		this.query = new QueryHelper(event);
	}
}
