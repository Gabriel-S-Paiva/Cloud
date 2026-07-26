<script lang="ts">
	import { goto } from '$app/navigation';
	import { endpoints } from '$lib/api';
	import { auth } from '$lib/stores/auth.svelte';
    import { navigation } from '$lib/stores/navigation.svelte';
	import type { FolderContents } from '$lib/types';
	import { onMount } from 'svelte';

    let fetchErr = $state<string | null>(null);
    let folderContents = $state<FolderContents | null>(null)

    let folders = $derived(folderContents?.folders ?? []);
    let files = $derived(folderContents?.files ?? []);

    onMount(async () => {
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

    const openFile = (id: number): void => {
        endpoints.getFileContent(id, true)
    }

    const createFolder = (): void => {
        let id = navigation.currentFolderId ? navigation.currentFolderId : auth.user!.rootFolderId
        endpoints.createFolder('New Folder', id)
    }
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
    {/each}
    {#each folders as folder}
        <button type="button" onclick={() => enterFolder(folder.id, folder.displayName)}>{folder.displayName}</button>
    {/each}

    <button onclick={() => createFolder()}>Create Folder</button>
    <button>Create File</button>
</div>


