import type { ChatMessage, ChatSummary, Contact, DocumentFileSummary } from './types';

export const chats: ChatSummary[] = [
	{
		id: 'design-team',
		name: 'Design Team',
		avatarUrl: 'https://i.pravatar.cc/150?img=47',
		isGroup: true,
		memberCount: 12,
		onlineCount: 5,
		lastMessagePreview: 'Asiap maszeehh!',
		lastMessageAt: '2026-08-01T04:30:00',
		pinned: true,
		lastMessageStatus: 'read',
		channel: 'whatsapp_waha',
		agentName: 'Desi Irawati'
	},
	{
		id: 'random-james-1',
		name: 'Random James',
		avatarUrl: 'https://i.pravatar.cc/150?img=12',
		lastMessagePreview: 'Asiap maszeehh!',
		lastMessageAt: '2026-08-01T04:30:00',
		lastMessageStatus: 'read',
		channel: 'whatsapp_meta',
		labels: [
			{
				id: '1',
				label: 'Pesanan Baru',
				color: '#00c335'
			},
			{
				id: '2',
				label: 'Pembayaran Tertunda',
				color: '#ae662a'
			}
		]
	},
	{
		id: 'john-doe-1',
		name: 'John Doe',
		avatarUrl: 'https://i.pravatar.cc/150?img=33',
		lastMessagePreview: 'Asiap maszeehh!',
		lastMessageAt: '2026-08-01T04:30:00',
		unreadCount: 2,
		channel: 'messenger',
		agentName: 'Sarah Putri'
	},
	{
		id: 'bob-spinkler',
		name: 'Bob Spinkler',
		lastMessagePreview: 'Asiap maszeehh!',
		lastMessageAt: '2026-08-01T04:30:00',
		channel: 'instagram'
	},
	{
		id: 'column-martin',
		name: 'Column Martin',
		avatarUrl: 'https://i.pravatar.cc/150?img=25',
		lastMessagePreview: 'Asiap maszeehh!',
		lastMessageAt: '2026-08-01T04:30:00',
		lastMessageStatus: 'read',
		channel: 'telegram',
		labels: [
			{
				id: '1',
				label: 'Menunggu Pembayaran',
				color: '#d78d00'
			}
		]
	},
	{
		id: 'random-james-2',
		name: 'Random James',
		avatarUrl: 'https://i.pravatar.cc/150?img=14',
		lastMessagePreview: 'Asiap maszeehh!',
		lastMessageAt: '2026-08-01T04:30:00',
		lastMessageStatus: 'read',
		channel: 'whatsapp_meta',
		agentName: 'Ayu Anisa'
	},
	{
		id: 'john-doe-2',
		name: 'John Doe',
		avatarUrl: 'https://i.pravatar.cc/150?img=51',
		lastMessagePreview: 'Asiap maszeehh!',
		lastMessageAt: '2026-08-01T04:30:00',
		channel: 'messenger'
	}
];

export const activeContact: Contact = {
	id: 'catherine-richardson',
	name: 'Catherine Richardson',
	avatarUrl: 'https://i.pravatar.cc/150?img=32',
	phone: '+01-222-364522',
	email: 'catherine.richardson@gmail.com',
	address: '1134 Ridder Park Road, San Fransisco, CA 94851',
	location: 'San Francisco, CA',
	online: true
};

export const documents: DocumentFileSummary[] = [
	{ id: 'doc-1', name: 'Effects-of-global-warming.docx', sizeBytes: 73933, kind: 'DOCX' },
	{ id: 'doc-2', name: 'Effects-of-global-warming.docx', sizeBytes: 73933, kind: 'DOCX' },
	{ id: 'doc-3', name: 'Effects-of-global-warming.docx', sizeBytes: 73933, kind: 'DOCX' }
];

export const lastMedia: string[] = [
	'https://images.unsplash.com/photo-1441974231531-c6227db76b6e?w=200&q=60',
	'https://images.unsplash.com/photo-1470071459604-3b5ec3a7fe05?w=200&q=60',
	'https://images.unsplash.com/photo-1500534623283-312aade485b7?w=200&q=60'
];

export const messages: ChatMessage[] = [
	{
		id: 'm1',
		chatId: 'design-team',
		direction: 'out',
		timestamp: '2026-08-01T04:28:00',
		status: 'read',
		text: 'This is Content! This is Content! This is Content! This is Content!'
	},
	{
		id: 'm2',
		chatId: 'design-team',
		direction: 'out',
		timestamp: '2026-08-01T04:29:00',
		status: 'read',
		text: 'This is Content! This is Content! This is Content! This is Content!'
	},
	{
		id: 'm3',
		chatId: 'design-team',
		direction: 'out',
		timestamp: '2026-08-01T04:29:30',
		status: 'delivered',
		media: [
			{
				kind: 'image',
				url: 'https://images.unsplash.com/photo-1465146344425-f00d5f5c8f07?w=800&q=70',
				name: 'meadow.jpg'
			}
		]
	},
	{
		id: 'm4',
		chatId: 'design-team',
		direction: 'in',
		senderName: 'Putri Tanjak',
		senderAvatarUrl: 'https://i.pravatar.cc/150?img=44',
		timestamp: '2026-08-01T04:30:00',
		text: 'This is Content! This is Content! This is Content! This is Content!'
	},
	{
		id: 'm5',
		chatId: 'design-team',
		direction: 'in',
		senderName: 'Putri Tanjak',
		senderAvatarUrl: 'https://i.pravatar.cc/150?img=44',
		timestamp: '2026-08-01T04:31:00',
		text: 'Sending over the brand deck for review, let me know your thoughts by EOD.',
		linkPreview: {
			url: 'https://www.figma.com/file/design-team-brand-deck',
			title: 'Brand Deck — Q3 Refresh',
			description: 'Figma file with the updated logo lockups, color tokens and type scale.',
			siteName: 'figma.com',
			imageUrl: 'https://images.unsplash.com/photo-1561070791-2526d30994b5?w=400&q=60'
		}
	},
	{
		id: 'm6',
		chatId: 'design-team',
		direction: 'out',
		timestamp: '2026-08-01T04:32:00',
		status: 'sent',
		media: [
			{
				kind: 'document',
				url: '#',
				name: 'Effects-of-global-warming.docx',
				sizeBytes: 73933,
				mimeType: 'DOCX'
			}
		],
		text: 'Here is the doc we discussed earlier.'
	},
	{
		id: 'm7',
		chatId: 'design-team',
		direction: 'in',
		senderName: 'Random James',
		senderAvatarUrl: 'https://i.pravatar.cc/150?img=12',
		timestamp: '2026-08-01T04:33:00',
		media: [
			{
				kind: 'video',
				url: 'https://cdn.coverr.co/videos/coverr-a-video-of-a-city-4k',
				thumbnailUrl: 'https://images.unsplash.com/photo-1483959651481-dc75b89291f4?w=800&q=70',
				durationSeconds: 42,
				name: 'city-drone.mp4'
			}
		]
	},
	{
		id: 'm8',
		chatId: 'design-team',
		direction: 'in',
		senderName: 'Design Team Bot',
		senderAvatarUrl: 'https://i.pravatar.cc/150?img=68',
		timestamp: '2026-08-01T04:34:00',
		template: {
			headerText: 'Order Update',
			bodyText: 'Your asset request #4821 has been approved and is ready for download.',
			footerText: 'Design Team Automations',
			buttons: [
				{ label: 'Download', kind: 'url' },
				{ label: 'View details', kind: 'reply' }
			]
		}
	}
];
