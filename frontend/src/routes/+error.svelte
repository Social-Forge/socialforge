<script lang="ts">
	import { page } from '$app/state';
	import { localizeHref } from '$lib/paraglide/runtime';
</script>

<svelte:head>
	<title>{page.status} - {page.error?.message}</title>
	<meta name="description" content="Error {page.status}: {page.error?.message}" />
	<meta name="viewport" content="width=device-width, initial-scale=1" />
	<link rel="icon" href="/favicon.ico" type="image/x-icon" />
</svelte:head>
<main class="relative flex flex-col items-center justify-center px-4">
	<div class="mx-auto flex h-screen flex-col items-center justify-center">
		<div class="flex h-full flex-col items-center justify-center">
			<span
				class="not-found rounded-md bg-linear-to-br from-violet-400 to-purple-600 px-3.5 py-1 text-sm font-medium text-neutral-50"
			>
				{page.status}
			</span>
			<h1 class="mt-5 text-3xl font-bold text-neutral-50 md:text-5xl">
				{page.status === 404 ? 'Not Found' : page.status === 403 ? 'Forbidden' : 'Error'}
			</h1>
			<a href={localizeHref('/')} class="mt-8">
				<button
					class="cursor-pointer rounded-md bg-linear-to-br from-violet-400 to-purple-600 p-4 text-sm font-medium text-neutral-50 active:scale-95"
				>
					Back to homepage
				</button>
			</a>
		</div>
	</div>
</main>
