<script lang="ts">
	import { endpoints } from '$lib/api';
    import { navigation } from '$lib/stores/navigation.svelte';
    import { goto } from '$app/navigation';
	import Button from '$lib/components/UI/Button/Button.svelte';
	import Modal from '$lib/components/UI/Modal/Modal.svelte';
	import { toast } from '$lib/stores/toast.svelte';
	import type { Folder } from '$lib/types';
	import {Folder as FolderIcon} from '@lucide/svelte';

	interface Props {
		folder: Folder;
	}

	let { folder }: {folder: Folder} = $props();

	let isEditing = $state(false);
	let editName = $state(folder.displayName);
	$effect(() => {
		editName = folder.displayName;
	});

	let isModalOpen = $state(false);
	let activeModalTab = $state<'actions' | 'info' | 'delete'>('actions');

	const focus = (node: HTMLElement) => node.focus();

	const handleRename = async () => {
		if (!editName.trim() || editName === folder.displayName) {
			isEditing = false;
			return;
		}

		try {
			await endpoints.updateFolder(folder.id, editName);
			folder.displayName = editName;
			isEditing = false;
			toast.add('Folder renamed');
		} catch (err) {
			toast.add(err instanceof Error ? err.message : 'Failed to rename folder');
			editName = folder.displayName;
		}
	};

	const handleDelete = async () => {
		try {
			await endpoints.deleteFolder(folder.id);
			toast.add('Folder deleted');
		} catch (err) {
			toast.add(err instanceof Error ? err.message : 'Failed to delete folder');
		}
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
            err instanceof Error ? toast.add(err.message) : toast.add('Folder fetch failed');
        }
    };
</script>

<div 
	class="folder-card border p-4 rounded hover:bg-gray-50 select-none cursor-pointer flex items-center justify-between"
	ondblclick={() => enterFolder(folder.id, folder.displayName)}
	oncontextmenu={handleContextMenu}
	role="button"
	tabindex="0"
>
	<div class="flex items-center gap-3 truncate">
		<div class="file-icon font-bold text-lg">
			<FolderIcon/>
		</div>
		{#if isEditing}
			<input 
				type="text" 
				bind:value={editName} 
				onkeydown={handleKeyDown} 
				onblur={() => handleRename()}
				onclick={(e) => e.stopPropagation()} 
				use:focus
			/>
		{:else}
			<p 
				class="truncate font-medium" 
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
		class="px-2"
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
			<Button onclick={() => enterFolder(folder.id, folder.displayName)}>Open</Button>
			<Button onclick={() => (activeModalTab = 'info')}>Info</Button>
			<Button onclick={() => (isEditing = true, isModalOpen = false)}>Rename</Button>
			<Button variant="danger" onclick={() => (activeModalTab = 'delete')}>Delete</Button>
		</div>
	{:else if activeModalTab === 'info'}
		<div>
			<h3 class="font-bold mb-2">Folder Details</h3>
			<p><strong>Name:</strong> {folder.displayName}</p>
			<Button onclick={() => (activeModalTab = 'actions')}>Back</Button>
		</div>
	{:else if activeModalTab === 'delete'}
		<div>
			<h3 class="font-bold mb-2">Confirm Delete</h3>
			<p>Are you sure you want to delete <strong>{folder.displayName}</strong> and its contents?</p>
			<div class="flex gap-2 mt-4">
				<Button variant="secondary" onclick={() => (activeModalTab = 'actions')}>Cancel</Button>
				<Button 
					variant="danger" 
					onclick={() => {
						handleDelete();
						isModalOpen = false;
					}}
				>
					Confirm Delete
				</Button>
			</div>
		</div>
	{/if}
</Modal>