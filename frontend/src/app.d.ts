/// <reference types="unplugin-icons/types/svelte" />
import type { RequestEvent, ResolveOptions, MaybePromise } from '@sveltejs/kit';
import { ApiHandler } from './lib/server/api';
import { ServiceHelper } from './lib/server/index';

declare global {
	namespace App {
		interface Error {
			code?: string;
		}
		interface Locals {
			api: ApiHandler;
			helper: ServiceHelper;
			safeGetUser: () => Promise<UserResponse | null>;
			user?: UserResponse | null;
			origin?: string;
			lang: string;
		}
		interface PageData {
			user?: UserResponse | null;
			success?: boolean;
			errors?: {
				code: string;
				message: string;
				details?: any;
			};
			messages?: string;
		}
		// interface PageState {}
		// interface Platform {}
	}
	interface RequestHandlerParams {
		event: RequestEvent;
		resolve: (event: RequestEvent) => MaybePromise<Response>;
		isAuthenticated: boolean;
		hasTenant: boolean;
		method: string;
		pathname: string;
	}
	interface RouteConfig {
		public?: boolean;
		roles?: ('superadmin' | 'tenant_owner' | 'supervisor' | 'agent')[];
		roleLevel?: (0 | 1 | 2 | 3)[];
		methods?: ('GET' | 'POST' | 'PUT' | 'PATCH' | 'DELETE')[];
		tenantRequired?: boolean;
	}
	interface RouteRules {
		[key: string]: RouteConfig;
	}

	interface ApiMiddlewareParams {
		event: RequestEvent;
		method: string;
		pathname: string;
		isAuthenticated: boolean;
		userRoleLevel: number;
		hasTenant: boolean;
	}
}

export {};
