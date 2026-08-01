export type MessageStatus = 'sending' | 'sent' | 'delivered' | 'read' | 'failed';

export type MessageDirection = 'in' | 'out';

export interface Contact {
	id: string;
	name: string;
	avatarUrl?: string;
	phone?: string;
	email?: string;
	address?: string;
	location?: string;
	online?: boolean;
}

export interface LinkPreview {
	url: string;
	title: string;
	description?: string;
	imageUrl?: string;
	siteName?: string;
}

export interface MediaAttachment {
	kind: 'image' | 'video' | 'document' | 'audio';
	url: string;
	name?: string;
	sizeBytes?: number;
	durationSeconds?: number;
	mimeType?: string;
	thumbnailUrl?: string;
}

export interface TemplateButton {
	label: string;
	kind: 'reply' | 'url' | 'call';
}

export interface TemplateContent {
	headerText?: string;
	bodyText: string;
	footerText?: string;
	buttons?: TemplateButton[];
}

/**
 * Mirrors the shape of WhatsApp Cloud API message types closely enough to
 * drive the UI: text, media (image/video/document/audio), text+media caption,
 * link preview, text+link preview, template, location and reaction.
 */
export type ChatMessage = {
	id: string;
	chatId: string;
	direction: MessageDirection;
	senderName?: string;
	senderAvatarUrl?: string;
	timestamp: string;
	status?: MessageStatus;
	text?: string;
	media?: MediaAttachment[];
	linkPreview?: LinkPreview;
	template?: TemplateContent;
	location?: { label: string; lat: number; lng: number };
	replyTo?: { id: string; preview: string; senderName?: string };
	reactions?: { emoji: string; count: number }[];
};

export interface ChatSummary {
	id: string;
	name: string;
	avatarUrl?: string;
	isGroup?: boolean;
	memberCount?: number;
	onlineCount?: number;
	lastMessagePreview: string;
	lastMessageAt: string;
	unreadCount?: number;
	pinned?: boolean;
	lastMessageStatus?: MessageStatus;
	presence?: 'online' | 'offline';
	channel: 'whatsapp_waha' | 'whatsapp_meta' | 'messenger' | 'instagram' | 'telegram';
	agentName?: string;
	labels?: ChipLabel[];
}

export interface DocumentFileSummary {
	id: string;
	name: string;
	sizeBytes: number;
	kind: 'DOCX' | 'PDF' | 'XLSX' | 'PPTX';
}

export interface ChipLabel {
	id: string;
	label: string;
	color: string; // Hex color atau tailwind color class
}
