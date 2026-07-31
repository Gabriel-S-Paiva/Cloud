<script lang="ts">
	import { goto } from '$app/navigation';
	import { endpoints } from '$lib/api';
	import { auth } from '$lib/stores/auth.svelte';
    import { navigation } from '$lib/stores/navigation.svelte';
    import { toast } from '$lib/stores/toast.svelte';
	import type { UserSummary,  FolderContents } from '$lib/types';
	import { onMount } from 'svelte';
    import {page} from '$app/stores'
	import FileManager from '$lib/components/Reactive/FileManager.svelte';
	import Toast from '$lib/components/UI/Toast/Toast.svelte';

    // RENDERING
    let folderContents = $state<FolderContents | null>(null)
    
    // UPLOAD
    let selectedFiles = $state<FileList | null>(null)
    let uploadProgress = $state<{ uploaded: number; total: number } | null>(null);

    // SHARE
    let users = $state<UserSummary[] | null>(null)
    let shareSelections = $state<Record<number, { target: number | null; permission: 'Edit' | 'View' }>>({});
    

    onMount(async () => {
        const urlSegments = $page.params.path?.split('/').filter(Boolean) ?? [];

        if (urlSegments.length > 0 && navigation.path.length === 0) {
            goto('/home');
            return;
        }

        try{
            let folderId: number = navigation.currentFolderId ?? auth.user!.rootFolderId
            folderContents = await endpoints.getFolderContent(folderId)
            users = await endpoints.getSharableUsers()
            for (const item of [...folderContents.files, ...folderContents.folders]) {
                shareSelections[item.id] = { target: null, permission: 'View' };
            }
        } catch(err) {
            console.error(err)
            err instanceof Error ? toast.error(err.message) : toast.error('Folder Fecth Failed');
        }
    })

    const createFolder = async (): Promise<void> => {
        try {
            const id = navigation.currentFolderId ?? auth.user!.rootFolderId;
            await endpoints.createFolder('New Folder', id);
            folderContents = await endpoints.getFolderContent(id);
        } catch (err) {
            err instanceof Error ? toast.error(err.message) : toast.error('Folder Fecth Failed');
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
            err instanceof Error ? toast.error(err.message) : toast.error('Folder Fecth Failed');
        } finally {
            uploadProgress = null;
        }
    };

    function getSelection(id: number) {
        if (!shareSelections[id]) {
            shareSelections[id] = { target: null, permission: 'View' };
        }
        return shareSelections[id];
    }

    const shareFile = async (id: number) => {
        const sel = shareSelections[id];
        if (!sel?.target) return;
        try {
            await endpoints.createShare(id, null, sel.target, sel.permission);
        } catch (err) {
            err instanceof Error ? toast.error(err.message) : toast.error('Folder Fecth Failed');
        }
    };

    const shareFolder = async (id: number) => {
        const sel = shareSelections[id];
        if (!sel?.target) return;
        try {
            await endpoints.createShare(null, id, sel.target, sel.permission);
        } catch (err) {
            err instanceof Error ? toast.error(err.message) : toast.error('Folder Fecth Failed');
        }
    };

    $effect(() => {
        if (folderContents) {
            for (const item of [...folderContents.files, ...folderContents.folders]) {
                if (!shareSelections[item.id]) shareSelections[item.id] = { target: null, permission: 'View' };
            }
        }
        });
</script>

<div>
    <nav>
    <button onclick={() => goto(`/home`)}>Home</button>
        {#each navigation.path as segment, i}
            / <button onclick={() => { navigation.goToDepth(i); goto(`/home/${navigation.urlPath}`); }}>{segment.displayName}</button>
        {/each}
    </nav>
    {#if folderContents}
        <FileManager {...folderContents}/>
    {:else}
        No data
    {/if}
    <button onclick={() => createFolder()}>Create Folder</button>
    <input type="file" bind:files={selectedFiles} onchange={createFile} />


    {#if uploadProgress}
    <progress value={uploadProgress.uploaded} max={uploadProgress.total}></progress>
    <span>{Math.round((uploadProgress.uploaded / uploadProgress.total) * 100)}%</span>
    {/if}

    <Toast/>
</div>


