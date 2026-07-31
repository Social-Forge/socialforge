<script lang="ts">
	import { browser } from '$app/environment';
	import { setMode, mode } from 'mode-watcher';
	import { Button } from '$lib/components/ui/button';
	import Icon from '@iconify/svelte';

	const switchTheme = () => {
		if (mode.current === 'light') {
			setMode('dark');
		} else {
			setMode('light');
		}
	};
	const startViewTransition = (event: MouseEvent) => {
		if (!browser) return;
		if (!document.startViewTransition) {
			switchTheme();
			return;
		}
		const x = event.clientX;
		const y = event.clientY;
		const endRadius = Math.hypot(
			Math.max(x, window.innerWidth - x),
			Math.max(y, window.innerHeight - y)
		);

		const transition = document.startViewTransition(() => {
			switchTheme();
		});
		transition.ready.then(() => {
			const duration = 600;
			document.documentElement.animate(
				{
					clipPath: [`circle(0px at ${x}px ${y}px)`, `circle(${endRadius}px at ${x}px ${y}px)`]
				},
				{
					duration: duration,
					easing: 'cubic-bezier(.76,.32,.29,.99)',
					pseudoElement: '::view-transition-new(root)'
				}
			);
		});
	};
</script>

<Button onclick={startViewTransition} variant="ghost" size="icon" class="shadow-2xl">
	<Icon
		icon="line-md:moon-rising-filled-alt-loop"
		class="absolute h-8 w-8 scale-100 rotate-0 text-sky-600 transition-all dark:scale-0 dark:-rotate-90"
	/>
	<Icon
		icon="line-md:moon-filled-alt-to-sunny-filled-loop-transition"
		class="absolute h-8 w-8 scale-0 rotate-90 text-yellow-500 transition-all dark:scale-100 dark:rotate-0"
	/>
	<span class="sr-only">Toggle theme</span>
</Button>

