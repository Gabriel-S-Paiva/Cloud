<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import { endpoints } from '$lib/api';
	import { auth } from '$lib/stores/auth.svelte';
	import { navigation } from '$lib/stores/navigation.svelte';
	import { toast } from '$lib/stores/toast.svelte';
	import { driveContents } from '$lib/stores/driveContents.svelte';
	import { onMount } from 'svelte';
	import FileManager from '$lib/components/Reactive/FileManager.svelte';

	let loaded = $state(false);
	let selectedFiles = $state<FileList | null>(null);
	let uploadProgress = $state<{ uploaded: number; total: number } | null>(null);
	let fetchToken = 0;

	onMount(() => {
		const urlSegments = page.params.path?.split('/').filter(Boolean) ?? [];
		if (urlSegments.length > 0 && navigation.path.length === 0) {
			goto('/home');
			return;
		}
		driveContents.loadSharableUsers();
	});

	$effect(() => {
		const folderId = navigation.currentFolderId ?? auth.user?.rootFolderId;
		if (folderId == null) return;

		const token = ++fetchToken;
		loaded = false;

		endpoints
			.getFolderContent(folderId)
			.then((contents) => {
				if (token !== fetchToken) return; // a newer navigation already superseded this
				driveContents.setContents(contents.folders, contents.files);
				loaded = true;
			})
			.catch((err) => {
				if (token !== fetchToken) return;
				toast.error(err instanceof Error ? err.message : 'Folder fetch failed');
			});
	});

	const refetch = async () => {
		const folderId = navigation.currentFolderId ?? auth.user!.rootFolderId;
		const contents = await endpoints.getFolderContent(folderId);
		driveContents.setContents(contents.folders, contents.files);
	};

	const createFolder = async (): Promise<void> => {
		try {
			const folderId = navigation.currentFolderId ?? auth.user!.rootFolderId;
			await endpoints.createFolder('New Folder', folderId);
			await refetch();
		} catch (err) {
			toast.error(err instanceof Error ? err.message : 'Failed to create folder');
		}
	};

	const createFile = async () => {
		if (!selectedFiles) return;
		const file = selectedFiles[0];
		try {
			const parentId = navigation.currentFolderId ?? auth.user!.rootFolderId;
			const { id: newFileId } = await endpoints.createFile(
				file.name,
				parentId,
				file.size,
				file.type || 'application/octet-stream'
			);

			const CHUNK_SIZE = 20 * 1024 * 1024;
			let offset = 0;
			uploadProgress = { uploaded: 0, total: file.size };

			while (offset < file.size) {
				const chunk = file.slice(offset, offset + CHUNK_SIZE);
				await endpoints.uploadChunk(newFileId, chunk);
				offset += CHUNK_SIZE;
				uploadProgress = { uploaded: Math.min(offset, file.size), total: file.size };
			}

			await refetch();
		} catch (err) {
			toast.error(err instanceof Error ? err.message : 'Failed to upload file');
		} finally {
			uploadProgress = null;
		}
	};

	let dragOverCrumb = $state<number | null>(null);

	const handleCrumbDrop = async (e: DragEvent, targetFolderId: number) => {
		e.preventDefault();
		dragOverCrumb = null;
		const raw = e.dataTransfer?.getData('application/json');
		if (!raw) return;
		const { type, id } = JSON.parse(raw) as { type: 'file' | 'folder'; id: number };
		if (type === 'file') await driveContents.moveFile(id, targetFolderId);
		else await driveContents.moveFolder(id, targetFolderId);
	};
</script>

<div class="p-6">
	<nav class="flex items-center gap-1 text-sm text-text-muted mb-4">
		<button
			class="hover:text-text px-1.5 py-0.5 rounded {dragOverCrumb === auth.user?.rootFolderId ? 'bg-accent/10 text-accent' : ''}"
			onclick={() => {
					navigation.goToDepth(-1)
					goto('/home')
				}}
			ondragover={(e) => e.preventDefault()}
			ondragenter={(e) => {
				e.preventDefault();
				dragOverCrumb = auth.user!.rootFolderId;
			}}
			ondragleave={() => (dragOverCrumb = null)}
			ondrop={(e) => handleCrumbDrop(e, auth.user!.rootFolderId)}
		>
			Home
		</button>
		{#each navigation.path as segment, i}
			<span>/</span>
			<button
				class="hover:text-text px-1.5 py-0.5 rounded {dragOverCrumb === segment.id ? 'bg-accent/10 text-accent' : ''}"
				onclick={() => {
					navigation.goToDepth(i);
					goto(`/home/${navigation.urlPath}`);
				}}
				ondragover={(e) => e.preventDefault()}
				ondragenter={(e) => {
					e.preventDefault();
					dragOverCrumb = segment.id;
				}}
				ondragleave={() => (dragOverCrumb = null)}
				ondrop={(e) => handleCrumbDrop(e, segment.id)}
			>
				{segment.displayName}
			</button>
		{/each}
	</nav>

	<div class="flex items-center gap-2 mb-4">
		<button class="h-9 px-3 rounded-md bg-accent text-surface text-sm font-medium" onclick={createFolder}>
			+ New Folder
		</button>
		<label class="h-9 px-3 rounded-md border border-border text-sm text-text flex items-center cursor-pointer">
			Upload
			<input type="file" class="hidden" bind:files={selectedFiles} onchange={createFile} />
		</label>
	</div>

	{#if uploadProgress}
		<div class="mb-4">
			<div class="h-1 bg-border rounded-full overflow-hidden">
				<div
					class="h-full bg-accent transition-[width] duration-150"
					style="width: {(uploadProgress.uploaded / uploadProgress.total) * 100}%"
				></div>
			</div>
			<span class="text-xs text-text-muted">
				{Math.round((uploadProgress.uploaded / uploadProgress.total) * 100)}%
			</span>
		</div>
	{/if}

	{#if loaded}
		<FileManager />
	{:else}
		<p class="text-text-muted text-sm">Loading…</p>
	{/if}
</div>