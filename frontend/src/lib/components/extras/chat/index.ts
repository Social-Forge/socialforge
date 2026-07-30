import Bubble from '$lib/components/extras/chat/chat-bubble.svelte';
import BubbleMessage from '$lib/components/extras/chat/chat-bubble-message.svelte';
import BubbleAvatar from '$lib/components/extras/chat/chat-bubble-avatar.svelte';
import List from '$lib/components/extras/chat/chat-list.svelte';
import * as Avatar from '$lib/components/extras/avatar';

const BubbleAvatarImage = Avatar.Image;
const BubbleAvatarFallback = Avatar.Fallback;

export { List, Bubble, BubbleMessage, BubbleAvatar, BubbleAvatarImage, BubbleAvatarFallback };

export type {
	ChatListPropsWithoutHTML,
	ChatListProps,
	ChatBubblePropsWithoutHTML,
	ChatBubbleProps,
	ChatBubbleMessagePropsWithoutHTML,
	ChatBubbleMessageProps
} from '$lib/components/extras/chat/types';
