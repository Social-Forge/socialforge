<script lang="ts">
	import { Input } from '$lib/components/ui/input/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu/index.js';
	import { Separator } from '$lib/components/ui/separator/index.js';

	import { Search, CirclePlus, Settings, Plus, Minus } from '@lucide/svelte';

	let { data } = $props();
</script>

<div class="h-[calc(100dvh-70px)] w-full overflow-y-auto bg-background lg:h-full">
	<div class="m-10">
		<div class="flex flex-col gap-4">
			<div class="text-base font-medium lg:text-xl">List Agent</div>
			<div class="flex items-center justify-between gap-4">
				<div class="relative">
					<Search class="absolute top-1/2 left-2 -translate-y-1/2 text-muted-foreground" />
					<Input
						type="text"
						placeholder="Search User"
						class="ps-10"
						data-test="agent-input-search"
						data-cy="agent-input-search"
					/>
				</div>
				<div class="flex items-center gap-3">
					<Button>
						<CirclePlus />
						<span class="hidden md:block">Add Agent</span>
					</Button>

					<DropdownMenu.Root>
						<DropdownMenu.Trigger>
							{#snippet child({ props })}
								<Button {...props} variant="outline" size="icon">
									<Settings />
								</Button>
							{/snippet}
						</DropdownMenu.Trigger>
						<DropdownMenu.Content class="w-56" align="start">
							<DropdownMenu.Group>
								<DropdownMenu.Item>
									<Plus />
									Increase Limit
								</DropdownMenu.Item>
								<DropdownMenu.Item>
									<Minus />
									Reduce Limit
								</DropdownMenu.Item>
							</DropdownMenu.Group>
						</DropdownMenu.Content>
					</DropdownMenu.Root>
				</div>
			</div>
		</div>
		<Separator class="my-6" />
		<div
			class="mt-[15%] flex h-full w-full flex-col items-center justify-center select-none lg:mt-[10%]"
		>
			<div class="mt-10 text-lg font-bold">You haven't added a user yet.</div>
			<div class="mt-2 text-center text-sm text-muted-foreground">
				Please press the 'Add User' button to add a user.
			</div>
			<Button class="mt-3">
				<CirclePlus />
				<span>Add User</span>
			</Button>
		</div>
	</div>
</div>
