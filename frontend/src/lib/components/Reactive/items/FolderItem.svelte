<script lang="ts">
	import { driveContents } from '$lib/stores/driveContents.svelte';
	import { navigation } from '$lib/stores/navigation.svelte';
	import { goto } from '$app/navigation';
	import Button from '$lib/components/UI/Button/Button.svelte';
	import Modal from '$lib/components/UI/Modal/Modal.svelte';
	import { toast } from '$lib/stores/toast.svelte';
	import type { Folder } from '$lib/types';
	import { Folder as FolderIcon } from '@lucide/svelte';

	let { folder, variant }: { folder: Folder; variant: 'list' | 'grid' } = $props();

	let isEditing = $state(false);
	let editName = $state(folder.displayName);
	$effect(() => {
		editName = folder.displayName;
	});

	let isModalOpen = $state(false);
	let activeModalTab = $state<'actions' | 'info' | 'delete' | 'share'>('actions');

	let shareTarget = $state<number | null>(null);
	let sharePermission = $state<'Edit' | 'View'>('View');

	const focus = (node: HTMLElement) => node.focus();

	const handleRename = async () => {
		if (!editName.trim() || editName === folder.displayName) {
			isEditing = false;
			return;
		}
		await driveContents.renameFolder(folder.id, editName);
		isEditing = false;
	};

	const handleShare = async () => {
		if (!shareTarget) return;
		await driveContents.shareFolder(folder.id, shareTarget, sharePermission);
		isModalOpen = false;
		shareTarget = null;
		sharePermission = 'View';
	};

	const handleKeyDown = (e: KeyboardEvent) => {
		if (e.key === 'Enter') handleRename();
		if (e.key === 'Escape') {
			editName = folder.displayName;
			isEditing = false;
		}
	};

	const handleContextMenu = (e: MouseEvent) => {
		e.preventDefault();
		activeModalTab = 'actions';
		isModalOpen = true;
	};

	const enterFolder = async (id: number, displayName: string): Promise<void> => {
		try {
			navigation.enter({ id, displayName });
			goto(`/home/${navigation.urlPath}`);
		} catch (err) {
			err instanceof Error ? toast.error(err.message) : toast.error('Folder fetch failed');
		}
	};

	let isDragOver = $state(false);

	const handleDrop = async (e: DragEvent) => {
		e.preventDefault();
		isDragOver = false;
		const raw = e.dataTransfer?.getData('application/json');
		if (!raw) return;
		const { type, id } = JSON.parse(raw) as { type: 'file' | 'folder'; id: number };
		if (type === 'folder' && id === folder.id) return;
		if (type === 'file') await driveContents.moveFile(id, folder.id);
		else await driveContents.moveFolder(id, folder.id);
	};
</script>

<div
	draggable="true"
	ondragstart={(e) => {
		e.dataTransfer?.setData('application/json', JSON.stringify({ type: 'folder', id: folder.id }));
		if (e.dataTransfer) e.dataTransfer.effectAllowed = 'move';
	}}
	ondragover={(e) => {
		e.preventDefault();
		if (e.dataTransfer) e.dataTransfer.dropEffect = 'move';
	}}
	ondragenter={(e) => {
		e.preventDefault();
		isDragOver = true;
	}}
	ondragleave={() => (isDragOver = false)}
	ondrop={handleDrop}
	class="group flex cursor-pointer items-center justify-between gap-3 rounded-lg border bg-surface-raised px-3 py-2.5 transition-colors select-none {isDragOver
		? 'border-accent bg-accent/10'
		: 'border-border hover:border-accent/40'}"
	ondblclick={() => enterFolder(folder.id, folder.displayName)}
	oncontextmenu={handleContextMenu}
	role="button"
	tabindex="0"
>
	<div class="flex items-center gap-3 truncate">
		<FolderIcon size={18} class="shrink-0 text-accent" />
		{#if isEditing}
			<input
				type="text"
				bind:value={editName}
				onkeydown={handleKeyDown}
				onblur={() => handleRename()}
				onclick={(e) => e.stopPropagation()}
				use:focus
				class="border-b border-accent bg-transparent text-sm text-text outline-none"
			/>
		{:else}
			<p
				class="truncate text-sm font-medium text-text"
				ondblclick={(e) => {
					e.stopPropagation();
					isEditing = true;
				}}
			>
				{folder.displayName}
			</p>
		{/if}
	</div>

	<button
		type="button"
		class="shrink-0 px-1 text-text-muted hover:text-text"
		onclick={(e) => {
			e.stopPropagation();
			activeModalTab = 'actions';
			isModalOpen = true;
		}}
	>
		⋮
	</button>
</div>

<Modal open={isModalOpen} onclose={() => (isModalOpen = false)}>
	{#if activeModalTab === 'actions'}
		<div class="flex flex-col gap-2">
			<Button variant="secondary" onclick={() => enterFolder(folder.id, folder.displayName)}
				>Open</Button
			>
			<Button variant="secondary" onclick={() => (activeModalTab = 'share')}>Share</Button>
			<Button variant="secondary" onclick={() => (activeModalTab = 'info')}>Info</Button>
			<Button variant="secondary" onclick={() => ((isEditing = true), (isModalOpen = false))}
				>Rename</Button
			>
			<Button variant="danger" onclick={() => (activeModalTab = 'delete')}>Delete</Button>
		</div>
	{:else if activeModalTab === 'share'}
		<div class="flex flex-col gap-3">
			<h3 class="font-display text-lg text-text">Share "{folder.displayName}"</h3>
			<select
				bind:value={shareTarget}
				class="h-10 rounded-md border border-border bg-surface-raised px-3 text-sm text-text"
			>
				<option value={null} disabled>Select a person</option>
				{#each driveContents.sharableUsers as user (user.id)}
					<option value={user.id}>{user.username}</option>
				{/each}
			</select>
			<div class="flex gap-2">
				<button
					type="button"
					class="h-9 flex-1 rounded-md border text-sm {sharePermission === 'View'
						? 'border-accent text-accent'
						: 'border-border text-text-muted'}"
					onclick={() => (sharePermission = 'View')}
				>
					View
				</button>
				<button
					type="button"
					class="h-9 flex-1 rounded-md border text-sm {sharePermission === 'Edit'
						? 'border-accent text-accent'
						: 'border-border text-text-muted'}"
					onclick={() => (sharePermission = 'Edit')}
				>
					Edit
				</button>
			</div>
			<div class="mt-1 flex gap-2">
				<Button variant="secondary" onclick={() => (activeModalTab = 'actions')}>Back</Button>
				<Button onclick={handleShare}>Confirm Share</Button>
			</div>
		</div>
	{:else if activeModalTab === 'info'}
		<div class="flex flex-col gap-2">
			<h3 class="font-display text-lg text-text">Folder Details</h3>
			<p class="text-sm text-text">
				<span class="text-text-muted">Name:</span>
				{folder.displayName}
			</p>
			<Button variant="secondary" onclick={() => (activeModalTab = 'actions')}>Back</Button>
		</div>
	{:else if activeModalTab === 'delete'}
		<div class="flex flex-col gap-3">
			<h3 class="font-display text-lg text-text">Confirm Delete</h3>
			<p class="text-sm text-text">
				Are you sure you want to delete <strong>{folder.displayName}</strong> and its contents?
			</p>
			<div class="flex gap-2">
				<Button variant="secondary" onclick={() => (activeModalTab = 'actions')}>Cancel</Button>
				<Button
					variant="danger"
					onclick={() => {
						driveContents.deleteFolder(folder.id);
						isModalOpen = false;
					}}
				>
					Confirm Delete
				</Button>
			</div>
		</div>
	{/if}
</Modal>
