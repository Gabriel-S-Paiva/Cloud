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
</script>

<div
	class="group border border-border rounded-lg bg-surface-raised hover:border-accent/40 transition-colors select-none cursor-pointer flex items-center justify-between gap-3 px-3 py-2.5"
	ondblclick={() => enterFolder(folder.id, folder.displayName)}
	oncontextmenu={handleContextMenu}
	role="button"
	tabindex="0"
>
	<div class="flex items-center gap-3 truncate">
		<FolderIcon size={18} class="text-accent shrink-0" />
		{#if isEditing}
			<input
				type="text"
				bind:value={editName}
				onkeydown={handleKeyDown}
				onblur={() => handleRename()}
				onclick={(e) => e.stopPropagation()}
				use:focus
				class="bg-transparent border-b border-accent outline-none text-sm text-text"
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
		class="text-text-muted hover:text-text px-1 shrink-0"
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
			<Button variant="secondary" onclick={() => enterFolder(folder.id, folder.displayName)}>Open</Button>
			<Button variant="secondary" onclick={() => (activeModalTab = 'share')}>Share</Button>
			<Button variant="secondary" onclick={() => (activeModalTab = 'info')}>Info</Button>
			<Button variant="secondary" onclick={() => (isEditing = true, isModalOpen = false)}>Rename</Button>
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
			<h3 class="font-display text-lg text-text">Folder Details</h3>
			<p class="text-sm text-text"><span class="text-text-muted">Name:</span> {folder.displayName}</p>
			<Button variant="secondary" onclick={() => (activeModalTab = 'actions')}>Back</Button>
		</div>
	{:else if activeModalTab === 'delete'}
		<div class="flex flex-col gap-3">
			<h3 class="font-display text-lg text-text">Confirm Delete</h3>
			<p class="text-sm text-text">Are you sure you want to delete <strong>{folder.displayName}</strong> and its contents?</p>
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