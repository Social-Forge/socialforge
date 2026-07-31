<script lang="ts">
	let { title }: { title: string } = $props();
</script>

<div
	class="relative inline-flex h-8 overflow-hidden rounded-full p-[1.5px] select-none focus:outline-none"
>
	<span
		class="absolute inset-[-1000%] animate-[spin_3s_linear_infinite] bg-[conic-gradient(from_90deg_at_50%_50%,#6d28d9_0%,#d8b4fe_50%,#6d28d9_100%)]"
	></span>
	<span
		class="inline-flex h-full w-full cursor-pointer items-center justify-center rounded-full bg-slate-950 px-4 py-1 text-sm font-medium text-white backdrop-blur-3xl"
	>
		{title}
	</span>
</div>
