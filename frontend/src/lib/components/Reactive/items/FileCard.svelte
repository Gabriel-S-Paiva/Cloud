<script lang="ts">
	import { endpoints } from '$lib/api';
	import Button from '$lib/components/UI/Button/Button.svelte';
	import Modal from '$lib/components/UI/Modal/Modal.svelte';
	import { toast } from '$lib/stores/toast.svelte';
	import type { CloudFile } from '$lib/types';
	import { FileText } from '@lucide/svelte';



	let { file }: {file: CloudFile} = $props();

	// Local component state
	let isEditing = $state(false);
	let editName = $state(file.displayName);
    $effect(() => {
        editName = file.displayName;
    });
	let isModalOpen = $state(false);
	let activeModalTab = $state<'actions' | 'info' | 'delete'>('actions');

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

		try {
			await endpoints.updateFile(file.id, editName);
			file.displayName = editName; // Optimistic update
			isEditing = false;
			toast.success('File renamed');
		} catch (err) {
			toast.error(err instanceof Error ? err.message : 'Failed to rename file');
			editName = file.displayName; // Revert on failure
		}
	};

    const deleteFile = async (id: number) => {
        try {
            await endpoints.deleteFile(id);
            toast.success('File Deleted')
        } catch (err) {
            toast.error(err instanceof Error ? err.message : 'Failed to rename file');
        }
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
		activeModalTab = 'actions';
		isModalOpen = true;
	};
</script>

<div 
	class="file-card border p-4 rounded hover:bg-gray-50 select-none cursor-pointer"
	ondblclick={() => openFile(false)}
	oncontextmenu={handleContextMenu}
	role="button"
	tabindex="0"
>
	<div class="file-icon font-bold text-lg">
		<FileText/>
	</div>

	<div class="file-details mt-2">
		{#if isEditing}
			<input 
				type="text" 
				bind:value={editName} 
				onkeydown={handleKeyDown} 
				onblur={handleRename}
				onclick={(e) => e.stopPropagation()} 
				autofocus 
			/>
		{:else}
			<p 
				class="truncate" 
				ondblclick={(e) => {
					e.stopPropagation();
					isEditing = true;
				}}
			>
				{file.displayName}
			</p>
		{/if}

		<button 
			type="button" 
			onclick={(e) => {
				e.stopPropagation();
				activeModalTab = 'actions';
				isModalOpen = true;
			}}
		>
			⋮
		</button>
	</div>
</div>

<Modal open={isModalOpen} onclose={() => (isModalOpen = false)}>
	{#if activeModalTab === 'actions'}
		<div class="flex flex-col gap-2">
			<Button onclick={() => openFile(true)}>Download</Button>
			<Button onclick={() => (activeModalTab = 'info')}>Info</Button>
			<Button onclick={() => (isEditing = true, isModalOpen = false)}>Rename</Button>
			<Button variant="danger" onclick={() => (activeModalTab = 'delete')}>Delete</Button>
		</div>
	{:else if activeModalTab === 'info'}
		<div>
			<h3>File Details</h3>
			<p><strong>Name:</strong> {file.displayName}</p>
			<p><strong>Size:</strong> {(file.size / 1024 / 1024).toFixed(2)} MB</p>
			<p><strong>Type:</strong> {file.contentType}</p>
			<Button onclick={() => (activeModalTab = 'actions')}>Back</Button>
		</div>
	{:else if activeModalTab === 'delete'}
		<div>
			<h3>Confirm Delete</h3>
			<p>Are you sure you want to delete <strong>{file.displayName}</strong>?</p>
			<Button variant="secondary" onclick={() => (activeModalTab = 'actions')}>Cancel</Button>
			<Button 
				variant="danger" 
				onclick={() => {
					deleteFile(file.id);
					isModalOpen = false;
				}}
			>
				Confirm Delete
			</Button>
		</div>
	{/if}
</Modal>