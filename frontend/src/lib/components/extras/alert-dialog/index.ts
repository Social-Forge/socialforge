import Root from '$lib/components/extras/alert-dialog/alert-dialog.svelte';
import Portal from '$lib/components/extras/alert-dialog/alert-dialog-portal.svelte';
import Trigger from '$lib/components/extras/alert-dialog/alert-dialog-trigger.svelte';
import Title from '$lib/components/extras/alert-dialog/alert-dialog-title.svelte';
import Action from '$lib/components/extras/alert-dialog/alert-dialog-action.svelte';
import Cancel from '$lib/components/extras/alert-dialog/alert-dialog-cancel.svelte';
import Footer from '$lib/components/extras/alert-dialog/alert-dialog-footer.svelte';
import Header from '$lib/components/extras/alert-dialog/alert-dialog-header.svelte';
import Overlay from '$lib/components/extras/alert-dialog/alert-dialog-overlay.svelte';
import Content from '$lib/components/extras/alert-dialog/alert-dialog-content.svelte';
import Description from '$lib/components/extras/alert-dialog/alert-dialog-description.svelte';
import Media from '$lib/components/extras/alert-dialog/alert-dialog-media.svelte';

export {
	Root,
	Title,
	Action,
	Cancel,
	Portal,
	Footer,
	Header,
	Trigger,
	Overlay,
	Content,
	Description,
	Media,
	//
	Root as AlertDialog,
	Title as AlertDialogTitle,
	Action as AlertDialogAction,
	Cancel as AlertDialogCancel,
	Portal as AlertDialogPortal,
	Footer as AlertDialogFooter,
	Header as AlertDialogHeader,
	Trigger as AlertDialogTrigger,
	Overlay as AlertDialogOverlay,
	Content as AlertDialogContent,
	Description as AlertDialogDescription,
	Media as AlertDialogMedia
};
