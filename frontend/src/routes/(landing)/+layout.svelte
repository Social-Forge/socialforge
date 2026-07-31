<script lang="ts">
	import { HomeNavbar } from '$lib/components/landing';

	let { data, children } = $props();
</script>

<div
	class="min-full absolute inset-0 bg-[linear-gradient(to_right,#e5e5e5_1px,transparent_1px),linear-gradient(to_bottom,#e5e5e5_1px,transparent_1px)] mask-[radial-gradient(ellipse_60%_50%_at_50%_0%,#000_70%,transparent_110%)] bg-size-[3rem_3rem] dark:bg-[linear-gradient(to_right,#171717_1px,transparent_1px),linear-gradient(to_bottom,#171717_1px,transparent_1px)]"
></div>
<HomeNavbar lang={data.lang} />
<main class="relative z-0 mx-auto mt-20 w-full">
	{@render children()}
</main>
