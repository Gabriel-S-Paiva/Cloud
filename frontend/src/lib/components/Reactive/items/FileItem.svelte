<script lang="ts">
	import { driveContents } from '$lib/stores/driveContents.svelte';
	import { endpoints } from '$lib/api';
	import Button from '$lib/components/UI/Button/Button.svelte';
	import Modal from '$lib/components/UI/Modal/Modal.svelte';
	import { toast } from '$lib/stores/toast.svelte';
	import type { CloudFile } from '$lib/types';
	import { FileText } from '@lucide/svelte';

	let { file, variant }: { file: CloudFile; variant: 'list' | 'grid' } = $props();

	let isEditing = $state(false);
	let editName = $state(file.displayName);
	$effect(() => {
		editName = file.displayName;
	});
	let isModalOpen = $state(false);
	let activeModalTab = $state<'actions' | 'info' | 'delete' | 'share'>('actions');

	let shareTarget = $state<number | null>(null);
	let sharePermission = $state<'Edit' | 'View'>('View');

	const openFile = async (download = false): Promise<void> => {
		try {
			const blob = await endpoints.getFileContent(file.id, download);
			const url = URL.createObjectURL(blob);
			window.open(url, '_blank');
			setTimeout(() => URL.revokeObjectURL(url), 10000);
		} catch (err) {
			toast.error(err instanceof Error ? err.message : 'Failed to fetch file');
		}
	};

	const handleRename = async () => {
		if (!editName.trim() || editName === file.displayName) {
			isEditing = false;
			return;
		}
		await driveContents.renameFile(file.id, editName);
		isEditing = false;
	};

	const handleShare = async () => {
		if (!shareTarget) return;
		await driveContents.shareFile(file.id, shareTarget, sharePermission);
		isModalOpen = false;
		shareTarget = null;
		sharePermission = 'View';
	};

	const openModal = (tab: typeof activeModalTab) => {
		activeModalTab = tab;
		isModalOpen = true;
	};

	const handleKeyDown = (e: KeyboardEvent) => {
		if (e.key === 'Enter') handleRename();
		if (e.key === 'Escape') {
			editName = file.displayName;
			isEditing = false;
		}
	};

	const handleContextMenu = (e: MouseEvent) => {
		e.preventDefault();
		openModal('actions');
	};
</script>

<div
	class="group border border-border rounded-lg bg-surface-raised hover:border-accent/40 transition-colors select-none cursor-pointer {variant === 'list' ? 'flex items-center gap-3 px-3 py-2.5' : 'p-4 flex flex-col items-center gap-2 text-center'}"
	ondblclick={() => openFile(false)}
	oncontextmenu={handleContextMenu}
	role="button"
	tabindex="0"
>
	<FileText size={variant === 'grid' ? 32 : 18} class="text-accent-secondary shrink-0" />

	<div class="{variant === 'grid' ? 'w-full' : 'flex-1 min-w-0 flex items-center justify-between gap-2'}">
		{#if isEditing}
			<input
				type="text"
				bind:value={editName}
				onkeydown={handleKeyDown}
				onblur={handleRename}
				onclick={(e) => e.stopPropagation()}
				autofocus
				class="w-full bg-transparent border-b border-accent outline-none text-sm text-text"
			/>
		{:else}
			<p
				class="truncate text-sm text-text {variant === 'grid' ? 'w-full' : ''}"
				ondblclick={(e) => {
					e.stopPropagation();
					isEditing = true;
				}}
			>
				{file.displayName}
			</p>
		{/if}

		{#if variant === 'list'}
			<button
				type="button"
				class="text-text-muted hover:text-text px-1 shrink-0"
				onclick={(e) => {
					e.stopPropagation();
					openModal('actions');
				}}
			>
				⋮
			</button>
		{/if}
	</div>
</div>

<Modal open={isModalOpen} onclose={() => (isModalOpen = false)}>
	{#if activeModalTab === 'actions'}
		<div class="flex flex-col gap-2">
			<Button variant="secondary" onclick={() => openFile(true)}>Download</Button>
			<Button variant="secondary" onclick={() => (activeModalTab = 'share')}>Share</Button>
			<Button variant="secondary" onclick={() => (activeModalTab = 'info')}>Info</Button>
			<Button variant="secondary" onclick={() => (isEditing = true, isModalOpen = false)}>Rename</Button>
			<Button variant="danger" onclick={() => (activeModalTab = 'delete')}>Delete</Button>
		</div>
	{:else if activeModalTab === 'share'}
		<div class="flex flex-col gap-3">
			<h3 class="font-display text-lg text-text">Share "{file.displayName}"</h3>
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
					class="flex-1 h-9 rounded-md text-sm border {sharePermission === 'View' ? 'border-accent text-accent' : 'border-border text-text-muted'}"
					onclick={() => (sharePermission = 'View')}
				>
					View
				</button>
				<button
					type="button"
					class="flex-1 h-9 rounded-md text-sm border {sharePermission === 'Edit' ? 'border-accent text-accent' : 'border-border text-text-muted'}"
					onclick={() => (sharePermission = 'Edit')}
				>
					Edit
				</button>
			</div>
			<div class="flex gap-2 mt-1">
				<Button variant="secondary" onclick={() => (activeModalTab = 'actions')}>Back</Button>
				<Button onclick={handleShare}>Confirm Share</Button>
			</div>
		</div>
	{:else if activeModalTab === 'info'}
		<div class="flex flex-col gap-2">
			<h3 class="font-display text-lg text-text">File Details</h3>
			<p class="text-sm text-text"><span class="text-text-muted">Name:</span> {file.displayName}</p>
			<p class="text-sm text-text"><span class="text-text-muted">Size:</span> {(file.size / 1024 / 1024).toFixed(2)} MB</p>
			<p class="text-sm text-text"><span class="text-text-muted">Type:</span> {file.contentType}</p>
			<Button variant="secondary" onclick={() => (activeModalTab = 'actions')}>Back</Button>
		</div>
	{:else if activeModalTab === 'delete'}
		<div class="flex flex-col gap-3">
			<h3 class="font-display text-lg text-text">Confirm Delete</h3>
			<p class="text-sm text-text">Are you sure you want to delete <strong>{file.displayName}</strong>?</p>
			<div class="flex gap-2">
				<Button variant="secondary" onclick={() => (activeModalTab = 'actions')}>Cancel</Button>
				<Button
					variant="danger"
					onclick={() => {
						driveContents.deleteFile(file.id);
						isModalOpen = false;
					}}
				>
					Confirm Delete
				</Button>
			</div>
		</div>
	{/if}
</Modal>