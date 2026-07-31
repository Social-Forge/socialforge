<script lang="ts" module>
	interface Props {
		title?: string;
		description?: string;
		variant?:
			| 'primary'
			| 'neutral'
			| 'info'
			| 'success'
			| 'danger'
			| 'warning'
			| 'purple'
			| 'pink'
			| 'darken'
			| 'none';
		shadow?: 'none' | 'small' | 'medium' | 'large' | 'xl';
		children?: Snippet<[]>;
		spotlight?: boolean;
		spotlightIntensity?: 'subtle' | 'medium' | 'strong';
		class?: ClassValue;
		useBorder?: boolean;
	}
</script>

<script lang="ts">
	import { onMount } from 'svelte';
	import type { Snippet } from 'svelte';
	import type { ClassValue } from 'svelte/elements';
	import { cn } from '$lib/utils.js';

	let {
		title,
		description,
		variant = 'neutral',
		shadow = 'medium',
		children,
		spotlight = false,
		spotlightIntensity = 'medium',
		class: className,
		useBorder = true
	}: Props = $props();

	let containerRef = $state<HTMLDivElement | null>(null);
	let spotlightRef = $state<HTMLDivElement | null>(null);
	let isHovered = $state(false);

	// Container background & text colors
	const containerClass = $derived(() => {
		const baseClasses = 'relative overflow-hidden';
		switch (variant) {
			case 'primary':
				return `${baseClasses} bg-gradient-to-br from-white via-emerald-100 to-emerald-200 dark:from-transparent dark:via-emerald-900/50 dark:to-emerald-800 text-emerald-900 dark:text-emerald-50  ${useBorder ? ' border border-emerald-300 dark:border-emerald-700' : ''}`;
			case 'neutral':
				return `${baseClasses} bg-gradient-to-br from-white via-neutral-100 to-neutral-200 dark:from-transparent dark:via-neutral-900/50  dark:to-neutral-800 text-neutral-900 dark:text-neutral-50  ${useBorder ? 'border border-neutral-300 dark:border-neutral-700' : ''}`;
			case 'info':
				return `${baseClasses} bg-gradient-to-br from-white via-blue-100 to-blue-200 dark:from-transparent dark:via-blue-900/50  dark:to-blue-800 text-blue-900 dark:text-blue-50  ${useBorder ? 'border border-blue-300 dark:border-blue-700' : ''}`;
			case 'success':
				return `${baseClasses} bg-gradient-to-br from-white via-green-100 to-green-200 dark:from-transparent dark:via-green-900/50  dark:to-green-800 text-green-900 dark:text-green-50  ${useBorder ? 'border border-green-300 dark:border-green-700' : ''}`;
			case 'danger':
				return `${baseClasses} bg-gradient-to-br from-white via-red-100 to-red-200 dark:from-transparent dark:via-red-900 dark:to-red-800/50  text-red-900 dark:text-red-50  ${useBorder ? 'border border-red-300 dark:border-red-700' : ''}`;
			case 'warning':
				return `${baseClasses} bg-gradient-to-br from-white via-amber-100 to-amber-200 dark:from-transparent dark:via-amber-900/50  dark:to-amber-800 text-amber-900 dark:text-amber-50  ${useBorder ? 'border border-amber-300 dark:border-amber-700' : ''}`;
			case 'purple':
				return `${baseClasses} bg-gradient-to-br from-white via-purple-100 to-purple-200 dark:from-transparent dark:via-purple-900/50  dark:to-purple-800 text-purple-900 dark:text-purple-50  ${useBorder ? 'border border-purple-300 dark:border-purple-700' : ''}`;
			case 'pink':
				return `${baseClasses} bg-gradient-to-br from-white via-pink-100 to-pink-200 dark:from-transparent dark:via-pink-900/50  dark:to-pink-800 text-pink-900 dark:text-pink-50  ${useBorder ? 'border border-pink-300 dark:border-pink-700' : ''}`;
			case 'darken':
				return `${baseClasses} bg-gradient-to-br from-neutral-100 via-neutral-200 to-neutral-300 dark:from-black/10 dark:via-neutral-900/30 dark:to-neutral-800/50 text-white  ${useBorder ? 'border border-neutral-300 dark:border-neutral-700' : ''}`;
			case 'none':
				return `${baseClasses} bg-white dark:bg-black text-black dark:text-white  ${useBorder ? 'border border-neutral-300 dark:border-neutral-700' : ''}`;
			default:
				return `${baseClasses} bg-gradient-to-br from-white via-neutral-100 to-neutral-200 dark:from-transparent dark:via-neutral-900/50  dark:to-neutral-800 text-neutral-900 dark:text-neutral-50  ${useBorder ? 'border border-neutral-300 dark:border-neutral-700' : ''}`;
		}
	});

	// Spotlight gradient colors
	const spotlightColor = $derived(() => {
		switch (variant) {
			case 'primary':
				return '#388e5f'; // emerald-500
			case 'neutral':
				return '#737373'; // neutral-500
			case 'info':
				return '#3b82f6'; // blue-500
			case 'success':
				return '#22c55e'; // green-500
			case 'danger':
				return '#ef4444'; // red-500
			case 'warning':
				return '#f59e0b'; // amber-500
			case 'purple':
				return '#a855f7'; // purple-500
			case 'pink':
				return '#ec4899'; // pink-500
			case 'darken':
				return '#000'; // black
			case 'none':
				return '#737373';
			default:
				return '#737373';
		}
	});

	// Spotlight size based on intensity
	const spotlightSize = $derived(() => {
		switch (spotlightIntensity) {
			case 'subtle':
				return { size: '350px', blur: '150px' };
			case 'medium':
				return { size: '500px', blur: '200px' };
			case 'strong':
				return { size: '700px', blur: '280px' };
			default:
				return { size: '500px', blur: '200px' };
		}
	});

	// Border color
	const borderClass = $derived(() => {
		switch (variant) {
			case 'primary':
				return 'border-emerald-300 dark:border-emerald-600';
			case 'neutral':
				return 'border-neutral-300 dark:border-neutral-600';
			case 'info':
				return 'border-blue-300 dark:border-blue-600';
			case 'success':
				return 'border-green-300 dark:border-green-600';
			case 'danger':
				return 'border-red-300 dark:border-red-600';
			case 'warning':
				return 'border-amber-300 dark:border-amber-600';
			case 'purple':
				return 'border-purple-300 dark:border-purple-600';
			case 'pink':
				return 'border-pink-300 dark:border-pink-600';
			case 'darken':
				return 'border-black dark:border-white';
			case 'none':
				return 'border-neutral-300 dark:border-neutral-600';
			default:
				return 'border-neutral-300 dark:border-neutral-600';
		}
	});

	// Shadow
	const shadowClass = $derived(() => {
		switch (shadow) {
			case 'none':
				return '';
			case 'small':
				return 'shadow-sm';
			case 'medium':
				return 'shadow-md';
			case 'large':
				return 'shadow-lg';
			case 'xl':
				return 'shadow-xl';
			default:
				return 'shadow-md';
		}
	});

	function handleMouseMove(e: MouseEvent) {
		if (!spotlight || !spotlightRef || !containerRef) return;

		const rect = containerRef.getBoundingClientRect();
		const x = e.clientX - rect.left;
		const y = e.clientY - rect.top;

		spotlightRef.style.setProperty('--mouse-x', `${x}px`);
		spotlightRef.style.setProperty('--mouse-y', `${y}px`);
	}

	function handleMouseEnter() {
		isHovered = true;
	}

	function handleMouseLeave() {
		isHovered = false;
	}

	onMount(() => {
		if (spotlight && containerRef) {
			containerRef.addEventListener('mousemove', handleMouseMove);
			containerRef.addEventListener('mouseenter', handleMouseEnter);
			containerRef.addEventListener('mouseleave', handleMouseLeave);

			return () => {
				containerRef?.removeEventListener('mousemove', handleMouseMove);
				containerRef?.removeEventListener('mouseenter', handleMouseEnter);
				containerRef?.removeEventListener('mouseleave', handleMouseLeave);
			};
		}
	});
</script>

<div
	bind:this={containerRef}
	class={cn(
		'relative h-full w-full space-y-4 rounded-xl p-4 transition-all duration-300 md:p-6',
		containerClass(),
		shadowClass(),
		className
	)}
>
	<!-- Spotlight Effect -->
	{#if spotlight}
		<div
			bind:this={spotlightRef}
			class="spotlight-overlay pointer-events-none absolute inset-0 z-0 opacity-0 transition-opacity duration-300"
			class:opacity-100={isHovered}
			style="
				--spotlight-color: {spotlightColor()};
				--spotlight-size: {spotlightSize().size};
				--spotlight-blur: {spotlightSize().blur};
			"
		></div>
	{/if}

	<!-- Content -->
	<div class="relative z-10 min-h-min space-y-4">
		{#if title || description}
			<div class={cn('space-y-2 border-b pb-3 md:pb-4', useBorder ? borderClass() : '')}>
				{#if title}
					<h2 class="text-xl font-bold md:text-2xl">
						{title}
					</h2>
				{/if}
				{#if description}
					<p class="text-sm opacity-70 md:text-base">{description}</p>
				{/if}
			</div>
		{/if}

		{@render children?.()}
	</div>
</div>

<style>
	.spotlight-overlay {
		--mouse-x: 50%;
		--mouse-y: 50%;
		--spotlight-color: #737373;
		--spotlight-size: 500px;
		--spotlight-blur: 200px;

		background: radial-gradient(
			circle var(--spotlight-size) at var(--mouse-x) var(--mouse-y),
			var(--spotlight-color),
			transparent 50%
		);
		filter: blur(var(--spotlight-blur));
	}
</style>
