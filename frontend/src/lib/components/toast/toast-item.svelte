<script lang="ts" module>
	interface Props {
		type: 'success' | 'error' | 'warning' | 'info';
		message: string;
		duration: number;
		onDismiss: () => void;
	}
</script>

<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import { cubicOut } from 'svelte/easing';
	import { fly } from 'svelte/transition';
	import CheckIcon from '~icons/fluent-color/checkmark-circle-48';
	import DismissIcon from '~icons/fluent-color/dismiss-circle-48';
	import InfoIcon from '~icons/flat-color-icons/info';
	import WarningIcon from '~icons/fluent-color/warning-48';
	import XIcon from '~icons/lucide/x';

	let { type = 'info', message, duration, onDismiss }: Props = $props();

	let progress = $state(100);
	let timeoutId: ReturnType<typeof setTimeout>;
	let intervalId: ReturnType<typeof setInterval>;

	const typeColors = {
		success: 'bg-green-500',
		error: 'bg-red-500',
		warning: 'bg-yellow-500',
		info: 'bg-blue-500'
	};

	const typeIcons = {
		success: CheckIcon,
		error: DismissIcon,
		warning: WarningIcon,
		info: InfoIcon
	};

	const CurrentIcon = $derived(typeIcons[type]);

	onMount(() => {
		// Start progress bar
		const startTime = Date.now();
		const endTime = startTime + duration;

		intervalId = setInterval(() => {
			const remaining = endTime - Date.now();
			progress = (remaining / duration) * 100;

			if (remaining <= 0) {
				clearInterval(intervalId);
				dismiss();
			}
		}, 50);

		// Auto-dismiss after duration
		timeoutId = setTimeout(dismiss, duration);
	});

	function dismiss() {
		clearTimeout(timeoutId);
		clearInterval(intervalId);
		onDismiss();
	}

	onDestroy(() => {
		clearTimeout(timeoutId);
		clearInterval(intervalId);
	});
</script>

<div
	class="fixed top-0 right-4 z-999 w-full max-w-xs overflow-hidden rounded-lg shadow-lg"
	in:fly={{ y: 50, duration: 300, easing: cubicOut }}
	out:fly={{ y: 50, duration: 200, easing: cubicOut }}
>
	<div class="rounded-lg bg-white p-4 shadow-md dark:bg-neutral-800">
		<div class="flex w-full items-center justify-between">
			<div class="flex w-full items-center gap-2">
				<div class="rounded-full p-1">
					<CurrentIcon class="h-5 w-5" />
				</div>
				<div class="text-sm font-medium text-neutral-900 dark:text-neutral-200">
					{message}
				</div>
			</div>
			<div class="ml-4 flex shrink-0">
				<button
					type="button"
					class="inline-flex rounded-md text-neutral-600 hover:text-neutral-400 focus:outline-none dark:text-neutral-400 dark:hover:text-neutral-600"
					onclick={dismiss}
				>
					<span class="sr-only">Close</span>
					<XIcon class="h-5 w-5" />
				</button>
			</div>
		</div>
		<div class="mt-2">
			<div class="h-1 w-full rounded-full bg-neutral-200">
				<div
					class={`h-1 rounded-full ${typeColors[type]}`}
					style={`width: ${progress}%; transition: width 50ms linear`}
				></div>
			</div>
		</div>
	</div>
</div>
