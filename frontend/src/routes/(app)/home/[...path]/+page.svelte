<script lang="ts">
	import { goto } from '$app/navigation';
	import { endpoints } from '$lib/api';
	import { auth } from '$lib/stores/auth.svelte';
    import { navigation } from '$lib/stores/navigation.svelte';
	import type { FolderContents } from '$lib/types';
	import { onMount } from 'svelte';
    import {page} from '$app/stores'

    // ERROR
    let fetchErr = $state<string | null>(null);

    // RENDERING
    let folderContents = $state<FolderContents | null>(null)
    let folders = $derived(folderContents?.folders ?? []);
    let files = $derived(folderContents?.files ?? []); 
    
    // UPLOAD
    let selectedFiles = $state<FileList | null>(null)
    let uploadProgress = $state<{ uploaded: number; total: number } | null>(null);

    onMount(async () => {
        const urlSegments = $page.params.path?.split('/').filter(Boolean) ?? [];

        if (urlSegments.length > 0 && navigation.path.length === 0) {
            goto('/home');
            return;
        }

        try{
            let folderId: number = navigation.currentFolderId ?? auth.user!.rootFolderId
            folderContents = await endpoints.getFolderContent(folderId)
        } catch(err) {
            fetchErr = err instanceof Error ? err.message : 'Folder Fecth Failed';
        }
    })

    const enterFolder = async (id: number, displayName: string): Promise<void> => {
        try {
            navigation.enter({ id, displayName });
            folderContents = await endpoints.getFolderContent(id);
            goto(`/home/${navigation.urlPath}`);
        } catch (err) {
            fetchErr = err instanceof Error ? err.message : 'Folder fetch failed';
        }
    };

    const openFile = async (id: number, download = false): Promise<void> => {
        try {
            const blob = await endpoints.getFileContent(id, download);
            const url = URL.createObjectURL(blob);
            window.open(url, '_blank');
        } catch (err) {
            fetchErr = err instanceof Error ? err.message : 'File fetch failed';
        }
    };

    const createFolder = async (): Promise<void> => {
        try {
            const id = navigation.currentFolderId ?? auth.user!.rootFolderId;
            await endpoints.createFolder('New Folder', id);
            folderContents = await endpoints.getFolderContent(id);
        } catch (err) {
            fetchErr = err instanceof Error ? err.message : 'Could not create folder';
        }
    };

    const createFile = async () => {
        if (!selectedFiles) return;
        const file = selectedFiles[0];
        try {
            const parentId = navigation.currentFolderId ?? auth.user!.rootFolderId;
            const { id: newFileId } = await endpoints.createFile(file.name, parentId, file.size, file.type || 'application/octet-stream');

            const CHUNK_SIZE = 20 * 1024 * 1024;
            let offset = 0;
            uploadProgress = { uploaded: 0, total: file.size };

            while (offset < file.size) {
                const chunk = file.slice(offset, offset + CHUNK_SIZE);
                await endpoints.uploadChunk(newFileId, chunk);
                offset += CHUNK_SIZE;
                uploadProgress = { uploaded: Math.min(offset, file.size), total: file.size };
            }

            folderContents = await endpoints.getFolderContent(parentId);
        } catch (err) {
            fetchErr = err instanceof Error ? err.message : 'Upload failed';
        } finally {
            uploadProgress = null;
        }
    };

    const deleteFile = async (id: number) => {
        try {
            await endpoints.deleteFile(id);
            const pid = navigation.currentFolderId ?? auth.user!.rootFolderId;
            folderContents = await endpoints.getFolderContent(pid);
        } catch (err) {
            fetchErr = err instanceof Error ? err.message : 'Delete failed';
        }
        };

    const renameFile = async (id: number, newName: string) => {
        try {
            await endpoints.updateFile(id, newName);
            const pid = navigation.currentFolderId ?? auth.user!.rootFolderId;
            folderContents = await endpoints.getFolderContent(pid);
        } catch (err) {
            fetchErr = err instanceof Error ? err.message : 'Rename failed';
        }
    };

    const deleteFolder = async (id: number) => {
        try {
            await endpoints.deleteFolder(id);
            const pid = navigation.currentFolderId ?? auth.user!.rootFolderId;
            folderContents = await endpoints.getFolderContent(pid);
        } catch (err) {
            fetchErr = err instanceof Error ? err.message : 'Delete failed';
        }
        };

    const renameFolder = async (id: number, newName: string) => {
        try {
            await endpoints.updateFolder(id, newName);
            const pid = navigation.currentFolderId ?? auth.user!.rootFolderId;
            folderContents = await endpoints.getFolderContent(pid);
        } catch (err) {
            fetchErr = err instanceof Error ? err.message : 'Rename failed';
        }
    };
</script>

<div>
    <nav>
    <button onclick={() => goto(`/home`)}>Home</button>
        {#each navigation.path as segment, i}
            / <button onclick={() => { navigation.goToDepth(i); goto(`/home/${navigation.urlPath}`); }}>{segment.displayName}</button>
        {/each}
    </nav>
    {#each files as file}
        <button type="button" onclick={() => openFile(file.id)}>{file.displayName}</button>
        <button type="button" onclick={() => deleteFile(file.id)}>Delete</button>
        <input bind:value={file.displayName} />
        <button type="button" onclick={() => renameFile(file.id, file.displayName)}>Rename</button>
    {/each}
    {#each folders as folder}
        <button type="button" onclick={() => enterFolder(folder.id, folder.displayName)}>{folder.displayName}</button>
        <button type="button" onclick={() => deleteFolder(folder.id)}>Delete Folder</button>
        <input bind:value={folder.displayName}><button type="button" onclick={() => renameFolder(folder.id,folder.displayName)}>Rename</button>
    {/each}

    <button onclick={() => createFolder()}>Create Folder</button>
    <input type="file" bind:files={selectedFiles} onchange={createFile} />


    {#if uploadProgress}
    <progress value={uploadProgress.uploaded} max={uploadProgress.total}></progress>
    <span>{Math.round((uploadProgress.uploaded / uploadProgress.total) * 100)}%</span>
    {/if}
</div>


