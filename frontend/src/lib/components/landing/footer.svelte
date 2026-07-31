<script lang="ts">
	import { SvelteDate } from 'svelte/reactivity';
	import { Separator } from '$lib/components/ui/separator';
	import { ChevronsDown } from '@lucide/svelte';

	let currentYear = $state(new SvelteDate().getFullYear());
</script>

<footer id="footer" class="container pt-24">
	<div class="rounded-2xl border bg-muted/50 p-10 dark:bg-card">
		<div class="grid grid-cols-2 gap-x-12 gap-y-8 md:grid-cols-4 xl:grid-cols-6">
			<div class="col-span-full xl:col-span-2">
				<a href="/" class="flex items-center font-bold">
					<ChevronsDown
						class="mr-2 h-9 w-9 rounded-lg border bg-linear-to-tr from-primary via-primary/70 to-primary text-white"
					/>
					<h3 class="text-2xl">Shadcn-Svelte</h3>
				</a>
			</div>

			<div class="flex flex-col gap-2">
				<h3 class="text-lg font-bold">Contact</h3>
				<div><a href="https://github.com" class="opacity-60 hover:opacity-100">Github</a></div>
				<div><a href="https://twitter.com" class="opacity-60 hover:opacity-100">Twitter</a></div>
				<div>
					<a href="https://instagram.com" class="opacity-60 hover:opacity-100">Instagram</a>
				</div>
			</div>

			<div class="flex flex-col gap-2">
				<h3 class="text-lg font-bold">Platforms</h3>
				<div><span class="cursor-default opacity-60 hover:opacity-100">iOS</span></div>
				<div><span class="cursor-default opacity-60 hover:opacity-100">Android</span></div>
				<div><span class="cursor-default opacity-60 hover:opacity-100">Web</span></div>
			</div>

			<div class="flex flex-col gap-2">
				<h3 class="text-lg font-bold">Help</h3>
				<div><a href="#contact" class="opacity-60 hover:opacity-100">Contact Us</a></div>
				<div><a href="#faq" class="opacity-60 hover:opacity-100">FAQ</a></div>
				<div><a href="#contact" class="opacity-60 hover:opacity-100">Feedback</a></div>
			</div>

			<div class="flex flex-col gap-2">
				<h3 class="text-lg font-bold">Socials</h3>
				<div><a href="https://twitch.tv" class="opacity-60 hover:opacity-100">Twitch</a></div>
				<div><a href="https://discord.com" class="opacity-60 hover:opacity-100">Discord</a></div>
				<div><a href="https://dribbble.com" class="opacity-60 hover:opacity-100">Dribbble</a></div>
			</div>
		</div>

		<Separator class="my-6" />

		<section>
			<h3>
				&copy; {currentYear} Designed and developed by
				<a
					target="_blank"
					href="https://github.com/leoMirandaa"
					class="border-primary text-primary transition-all hover:border-b-2"
					rel="noreferrer"
				>
					Social Forge
				</a>
			</h3>
		</section>
	</div>
</footer>
