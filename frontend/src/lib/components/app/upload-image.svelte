<script lang="ts" module>
	type UploadedFile = {
		name: string;
		type: string;
		size: number;
		uploadedAt: number;
		url: Promise<string>;
	};
</script>

<script lang="ts">
	import { invalidateAll } from '$app/navigation';
	import { onDestroy } from 'svelte';
	import { SvelteDate } from 'svelte/reactivity';
	import { Button, buttonVariants } from '$lib/components/ui/button/index.js';
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import { CameraIcon, XIcon } from '@lucide/svelte';
	import { displaySize, MEGABYTE } from '$lib/components/extras/file-drop-zone/index.js';
	import * as FileDropZone from '$lib/components/extras/file-drop-zone/index.js';
	import { Progress } from '$lib/components/ui/progress/index.js';
	import * as Empty from '$lib/components/ui/empty/index.js';
	import { Spinner } from '$lib/components/ui/spinner/index.js';
	import { toast } from '$lib/stores/toast.js';

	let {
		fileName,
		fieldName,
		apiUrl,
		method = 'POST',
		label,
		placeholder,
		onchange
	}: {
		fileName?: string;
		fieldName: string;
		apiUrl: string;
		method?: 'POST' | 'PUT';
		label?: string;
		placeholder?: string;
		onchange?: (image: string) => void;
	} = $props();

	let open = $state(false);
	let files = $state<UploadedFile[]>([]);
	let date = new SvelteDate();
	let avatarFile = $state<File | null>(null);
	let isUploading = $state(false);

	const onUpload: FileDropZone.FileDropZoneRootProps['onUpload'] = async (files) => {
		await Promise.allSettled(files.map((file) => uploadFile(file)));
	};
	const onFileRejected: FileDropZone.FileDropZoneRootProps['onFileRejected'] = async ({
		reason,
		file
	}) => {
		toast.error(`${file.name} failed to upload: ${reason}!`);
	};

	const uploadFile = async (file: File) => {
		if (files.find((f) => f.name === file.name)) return;
		const urlPromise = new Promise<string>((resolve) => {
			sleep(1000).then(() => resolve(URL.createObjectURL(file)));
		});

		files.push({
			name: `${new Date().getTime()}_${fileName?.toLowerCase() || 'user'}`,
			type: file.type,
			size: file.size,
			uploadedAt: Date.now(),
			url: urlPromise
		});
		avatarFile = file;
		await urlPromise;
	};
	function sleep(durationMs: number): Promise<void> {
		return new Promise((res) => setTimeout(res, durationMs));
	}

	async function uploadToServer() {
		if (!avatarFile) return;

		try {
			isUploading = true;
			const formData = new FormData();
			formData.append(fieldName, avatarFile);

			const res = await fetch(apiUrl, {
				method: method,
				body: formData
			});
			if (!res.ok) {
				throw new Error(await res.text());
			}
			const result = (await res.json()) as any;
			onchange?.(result.data);
			toast.success(`Image uploaded successfully: ${result.data}`);
			isUploading = false;
		} catch (error: any) {
			toast.error(`Image upload failed: ${error.message || 'Unknown error'}`);
		} finally {
			isUploading = false;
			await invalidateAll();
			open = false;
			files = [];
		}
	}
	onDestroy(async () => {
		for (const file of files) {
			URL.revokeObjectURL(await file.url);
		}
	});

	$effect(() => {
		const interval = setInterval(() => {
			date.setTime(Date.now());
		}, 10);
		return () => {
			clearInterval(interval);
		};
	});
</script>

<Dialog.Root bind:open onOpenChange={(val) => val && invalidateAll()}>
	<Dialog.Trigger
		class={buttonVariants({
			variant: 'outline',
			size: 'icon',
			className: 'absolute right-0 bottom-0 h-8 w-8 cursor-pointer rounded-full'
		})}
	>
		<CameraIcon />
	</Dialog.Trigger>
	<Dialog.Content>
		<Dialog.Header>
			<Dialog.Title>{label}</Dialog.Title>
			<Dialog.Description>
				{placeholder}
			</Dialog.Description>
		</Dialog.Header>
		<div class="flex w-full flex-col gap-2 p-6">
			{#if isUploading}
				<Empty.Root class="w-full">
					<Empty.Header>
						<Empty.Media variant="icon">
							<Spinner />
						</Empty.Media>
						<Empty.Title>Uploading...</Empty.Title>
						<Empty.Description>
							Please wait while we process your request. Do not refresh the page.
						</Empty.Description>
					</Empty.Header>
				</Empty.Root>
			{:else}
				<FileDropZone.Root
					{onUpload}
					{onFileRejected}
					maxFileSize={5 * MEGABYTE}
					fileCount={files.length}
					accept="image/*"
					maxFiles={1}
					disabled={files.length > 0 || isUploading}
				>
					<FileDropZone.Trigger />
				</FileDropZone.Root>
				<div class="flex flex-col gap-2">
					{#each files as file, i (file.name)}
						<div class="flex place-items-center justify-between gap-2">
							<div class="flex place-items-center gap-2">
								{#await file.url then src}
									<div class="relative size-9 overflow-clip">
										<img
											{src}
											alt={file.name}
											class="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 overflow-clip"
										/>
									</div>
								{/await}
								<div class="flex flex-col">
									<span>{file.name}</span>
									<span class="text-xs text-muted-foreground">{displaySize(file.size)}</span>
								</div>
							</div>
							{#await file.url}
								<Progress
									class="h-2 w-full grow"
									value={((date.getTime() - file.uploadedAt) / 1000) * 100}
									max={100}
								/>
							{:then url}
								<Button
									variant="outline"
									size="icon"
									onclick={() => {
										URL.revokeObjectURL(url);
										files = [...files.slice(0, i), ...files.slice(i + 1)];
									}}
								>
									<XIcon />
								</Button>
							{/await}
						</div>
					{/each}
				</div>
			{/if}
		</div>
		<Dialog.Footer>
			<Button type="button" onclick={uploadToServer}>Upload</Button>
			<Dialog.Close>
				<Button variant="destructive" size="default">Close</Button>
			</Dialog.Close>
		</Dialog.Footer>
	</Dialog.Content>
</Dialog.Root>
