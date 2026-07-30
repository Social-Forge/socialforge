import Root from '$lib/components/extras/drawer/drawer.svelte';
import Content from '$lib/components/extras/drawer/drawer-content.svelte';
import Description from '$lib/components/extras/drawer/drawer-description.svelte';
import Overlay from '$lib/components/extras/drawer/drawer-overlay.svelte';
import Footer from '$lib/components/extras/drawer/drawer-footer.svelte';
import Header from '$lib/components/extras/drawer/drawer-header.svelte';
import Title from '$lib/components/extras/drawer/drawer-title.svelte';
import NestedRoot from '$lib/components/extras/drawer/drawer-nested.svelte';
import Close from '$lib/components/extras/drawer/drawer-close.svelte';
import Trigger from '$lib/components/extras/drawer/drawer-trigger.svelte';
import Portal from '$lib/components/extras/drawer/drawer-portal.svelte';

export {
	Root,
	NestedRoot,
	Content,
	Description,
	Overlay,
	Footer,
	Header,
	Title,
	Trigger,
	Portal,
	Close,

	//
	Root as Drawer,
	NestedRoot as DrawerNestedRoot,
	Content as DrawerContent,
	Description as DrawerDescription,
	Overlay as DrawerOverlay,
	Footer as DrawerFooter,
	Header as DrawerHeader,
	Title as DrawerTitle,
	Trigger as DrawerTrigger,
	Portal as DrawerPortal,
	Close as DrawerClose
};
