import Root from '$lib/components/extras/tabs/tabs.svelte';
import Content from '$lib/components/extras/tabs/tabs-content.svelte';
import List, { tabsListVariants, type TabsListVariant } from '$lib/components/extras/tabs/tabs-list.svelte';
import Trigger from '$lib/components/extras/tabs/tabs-trigger.svelte';

export {
	Root,
	Content,
	List,
	Trigger,
	tabsListVariants,
	type TabsListVariant,
	//
	Root as Tabs,
	Content as TabsContent,
	List as TabsList,
	Trigger as TabsTrigger
};
