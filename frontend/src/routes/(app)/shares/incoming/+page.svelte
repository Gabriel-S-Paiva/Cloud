<script lang="ts">
	import { endpoints } from "$lib/api";
	import type { SharedContents } from "$lib/types";
    import { goto } from "$app/navigation";
    import { navigation } from "$lib/stores/navigation.svelte";
	import { onMount } from "svelte";

    let fetchErr = $state<string | null>(null);
    let incomingShares = $state<SharedContents|null>(null)
    let files = $derived(incomingShares?.files ?? []); 
    let folders = $derived(incomingShares?.folders ?? [])

    onMount( async () => {
        try{
            incomingShares = await endpoints.getShareIncoming()
        } catch(err) {
            fetchErr = err instanceof Error ? err.message : 'Shared With Me Fecth Failed';
        }
    })

    const openFile = async (id: number) => {
        try {
            await endpoints.getFileContent(id, true)
            incomingShares = await endpoints.getShareIncoming()
        } catch(err) {
            fetchErr = err instanceof Error ? err.message : ''
        }
    }

    const renameFile = async (id: number, displayName: string) => {
        try {
            await endpoints.updateFile(id, displayName)
            incomingShares = await endpoints.getShareIncoming()
        } catch(err) {
            fetchErr = err instanceof Error ? err.message : ''
        }
    }

    const deleteFile = async (id: number) => {
        try {
            await endpoints.deleteFile(id)
            incomingShares = await endpoints.getShareIncoming()
        } catch(err) {
            fetchErr = err instanceof Error ? err.message : ''
        }
    }

    const enterFolder = async (id: number, displayName: string) => {
        navigation.reset();
        navigation.enter({ id, displayName });
        goto(`/home/${navigation.urlPath}`);
    }

    const renameFolder = async (id: number, displayName: string) => {
        try {
            await endpoints.updateFolder(id, displayName)
            incomingShares = await endpoints.getShareIncoming()
        } catch(err) {
            fetchErr = err instanceof Error ? err.message : ''
        }
    }

    const deleteFolder = async (id: number) => {
        try {
            await endpoints.deleteFolder(id)
            incomingShares = await endpoints.getShareIncoming()
        } catch(err) {
            fetchErr = err instanceof Error ? err.message : ''
        }
    }

    const unshare = async (shareId: number) => {
        try {
            await endpoints.deleteShare(shareId)
            incomingShares = await endpoints.getShareIncoming()
        } catch(err) {
            fetchErr = err instanceof Error ? err.message : ''
        }
    }
</script>

{#each files as file}
    <button onclick={() => openFile(file.id)}>
        {file.displayName}, owned by: {file.ownedByUsername}
    </button>
    {#if file.permissions == "Edit"}<button onclick={() => deleteFile(file.id)}>Delete {file.displayName}</button>{/if}
    <button onclick={() => unshare(file.shareId)}>Leave Share</button>
    {#if file.permissions == "Edit"}<input bind:value={file.displayName}><button type="button" onclick={() => renameFile(file.id,file.displayName)}>Rename</button>{/if}
{/each}

{#each folders as folder}
    <button onclick={() => enterFolder(folder.id, folder.displayName)}>
        {folder.displayName}, owned by: {folder.ownedByUsername}
    </button>
    {#if folder.permissions == "Edit"}<button onclick={() => deleteFolder(folder.id)}>Delete {folder.displayName}</button>{/if}
    <button onclick={() => unshare(folder.shareId)}>Leave Share</button>
    {#if folder.permissions == "Edit"}<input bind:value={folder.displayName}><button type="button" onclick={() => renameFolder(folder.id,folder.displayName)}>Rename</button>{/if}
{/each}

{#if fetchErr}
    {fetchErr}
{/if}

