<script lang="ts">
	import type { ClassValue } from 'svelte/elements';
	import type { Snippet } from 'svelte';
	import { Motion } from 'svelte-motion';
	import { onMount } from 'svelte';

	let {
		class: className = 'size-full',
		children,
		variant = 'fade-up',
		delay = 0,
		duration = 0.6,
		distance = 20,
		once = true,
		threshold = 0.3
	}: {
		class?: ClassValue;
		children?: Snippet;
		variant?: 'fade-in' | 'fade-up' | 'fade-down' | 'fade-left' | 'fade-right';
		delay?: number;
		duration?: number;
		distance?: number;
		once?: boolean;
		threshold?: number;
	} = $props();

	let ref: HTMLDivElement;
	let isInView = $state(false);

	let currentVariant = $derived.by(() => {
		const variants = {
			'fade-in': {
				initial: { opacity: 0 },
				animate: { opacity: 1 }
			},
			'fade-up': {
				initial: { opacity: 0, y: distance },
				animate: { opacity: 1, y: 0 }
			},
			'fade-down': {
				initial: { opacity: 0, y: -distance },
				animate: { opacity: 1, y: 0 }
			},
			'fade-left': {
				initial: { opacity: 0, x: distance },
				animate: { opacity: 1, x: 0 }
			},
			'fade-right': {
				initial: { opacity: 0, x: -distance },
				animate: { opacity: 1, x: 0 }
			}
		};
		return variants[variant];
	});

	onMount(() => {
		const observer = new IntersectionObserver(
			(entries) => {
				entries.forEach((entry) => {
					if (entry.isIntersecting) {
						isInView = true;
						if (once) {
							observer.unobserve(entry.target);
						}
					} else if (!once) {
						isInView = false;
					}
				});
			},
			{
				threshold: threshold,
				rootMargin: '0px'
			}
		);

		if (ref) {
			observer.observe(ref);
		}

		return () => {
			if (ref) {
				observer.unobserve(ref);
			}
		};
	});
</script>

<div bind:this={ref}>
	<Motion
		initial={currentVariant.initial}
		animate={isInView ? currentVariant.animate : currentVariant.initial}
		transition={{
			duration: duration,
			delay: delay,
			ease: 'easeOut',
			type: 'spring',
			stiffness: 260,
			damping: 20
		}}
		let:motion
	>
		<div class={className} use:motion>
			{@render children?.()}
		</div>
	</Motion>
</div>
