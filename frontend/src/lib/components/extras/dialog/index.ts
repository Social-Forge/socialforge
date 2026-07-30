import Root from '$lib/components/extras/dialog/dialog.svelte';
import Portal from '$lib/components/extras/dialog/dialog-portal.svelte';
import Title from '$lib/components/extras/dialog/dialog-title.svelte';
import Footer from '$lib/components/extras/dialog/dialog-footer.svelte';
import Header from '$lib/components/extras/dialog/dialog-header.svelte';
import Overlay from '$lib/components/extras/dialog/dialog-overlay.svelte';
import Content from '$lib/components/extras/dialog/dialog-content.svelte';
import Description from '$lib/components/extras/dialog/dialog-description.svelte';
import Trigger from '$lib/components/extras/dialog/dialog-trigger.svelte';
import Close from '$lib/components/extras/dialog/dialog-close.svelte';

export {
	Root,
	Title,
	Portal,
	Footer,
	Header,
	Trigger,
	Overlay,
	Content,
	Description,
	Close,
	//
	Root as Dialog,
	Title as DialogTitle,
	Portal as DialogPortal,
	Footer as DialogFooter,
	Header as DialogHeader,
	Trigger as DialogTrigger,
	Overlay as DialogOverlay,
	Content as DialogContent,
	Description as DialogDescription,
	Close as DialogClose
};
