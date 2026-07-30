import Root from '$lib/components/extras/stepper/stepper.svelte';
import Nav from '$lib/components/extras/stepper/stepper-nav.svelte';
import Item from '$lib/components/extras/stepper/stepper-item.svelte';
import Trigger from '$lib/components/extras/stepper/stepper-trigger.svelte';
import Indicator from '$lib/components/extras/stepper/stepper-indicator.svelte';
import Separator from '$lib/components/extras/stepper/stepper-separator.svelte';
import Title from '$lib/components/extras/stepper/stepper-title.svelte';
import Description from '$lib/components/extras/stepper/stepper-description.svelte';
import Next from '$lib/components/extras/stepper/stepper-next.svelte';
import Previous from '$lib/components/extras/stepper/stepper-previous.svelte';

export {
	Root,
	Nav,
	Item,
	Trigger,
	Indicator,
	Separator,
	Title,
	Description,
	Next,
	Previous,
	//
	Root as Stepper,
	Nav as StepperNav,
	Item as StepperItem,
	Trigger as StepperTrigger,
	Indicator as StepperIndicator,
	Separator as StepperSeparator,
	Title as StepperTitle,
	Description as StepperDescription,
	Next as StepperNext,
	Previous as StepperPrevious
};
