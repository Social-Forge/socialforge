import Root from '$lib/components/extras/tooltip/tooltip.svelte';
import Trigger from '$lib/components/extras/tooltip/tooltip-trigger.svelte';
import Content from '$lib/components/extras/tooltip/tooltip-content.svelte';
import Provider from '$lib/components/extras/tooltip/tooltip-provider.svelte';
import Portal from '$lib/components/extras/tooltip/tooltip-portal.svelte';

export {
	Root,
	Trigger,
	Content,
	Provider,
	Portal,
	//
	Root as Tooltip,
	Content as TooltipContent,
	Trigger as TooltipTrigger,
	Provider as TooltipProvider,
	Portal as TooltipPortal
};
