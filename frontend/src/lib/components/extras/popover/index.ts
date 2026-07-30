import Root from '$lib/components/extras/popover/popover.svelte';
import Close from '$lib/components/extras/popover/popover-close.svelte';
import Content from '$lib/components/extras/popover/popover-content.svelte';
import Description from '$lib/components/extras/popover/popover-description.svelte';
import Header from '$lib/components/extras/popover/popover-header.svelte';
import Title from '$lib/components/extras/popover/popover-title.svelte';
import Trigger from '$lib/components/extras/popover/popover-trigger.svelte';
import Portal from '$lib/components/extras/popover/popover-portal.svelte';

export {
	Root,
	Content,
	Description,
	Header,
	Title,
	Trigger,
	Close,
	Portal,
	//
	Root as Popover,
	Content as PopoverContent,
	Description as PopoverDescription,
	Header as PopoverHeader,
	Title as PopoverTitle,
	Trigger as PopoverTrigger,
	Close as PopoverClose,
	Portal as PopoverPortal
};
