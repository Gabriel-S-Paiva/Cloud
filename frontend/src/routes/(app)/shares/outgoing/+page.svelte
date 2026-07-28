<script lang="ts">
	import { endpoints } from "$lib/api";
	import type { SharedContents } from "$lib/types";
	import { onMount } from "svelte";

    let fetchErr = $state<string | null>(null);
    let outgoingShares = $state<SharedContents|null>(null)
    let files = $derived(outgoingShares?.files ?? []); 
    let folders = $derived(outgoingShares?.folders ?? [])
    let permissionSelections = $state<Record<number, 'Edit' | 'View'>>({});

    onMount( async () => {
        try{
            outgoingShares = await endpoints.getShareOutgoing()
        } catch(err) {
            fetchErr = err instanceof Error ? err.message : 'Shared With Me Fecth Failed';
        }
    })

    const unshare = async (shareId: number) => {
        try {
            await endpoints.deleteShare(shareId)
            outgoingShares = await endpoints.getShareOutgoing()
        } catch(err) {
            fetchErr = err instanceof Error ? err.message : ''
        }
    }

    $effect(() => {
    if (outgoingShares) {
        for (const item of [...outgoingShares.files, ...outgoingShares.folders]) {
        if (!permissionSelections[item.shareId]) {
            permissionSelections[item.shareId] = item.permissions as 'Edit' | 'View';
        }
        }
    }
    });

    const changePermission = async (shareId: number) => {
    try {
        await endpoints.updateShare(shareId, permissionSelections[shareId]);
        outgoingShares = await endpoints.getShareOutgoing();
    } catch (err) {
        fetchErr = err instanceof Error ? err.message : 'Failed to update permission';
    }
    };
</script>
<ul>
    {#each files as file}
    <li>
        {file.displayName} shared with {file.sharedWith}
        <button onclick={() => unshare(file.shareId)}>Unshare</button>
    </li>
    <label><input type="radio" name={`perm-${file.shareId}`} value="View" bind:group={permissionSelections[file.shareId]}/> View</label>
    <label><input type="radio" name={`perm-${file.shareId}`} value="Edit" bind:group={permissionSelections[file.shareId]}/> Edit</label>
    <button onclick={() => changePermission(file.shareId)}>Update Permission</button>
{/each}
</ul>

<ul>
    {#each folders as folder}
        <li>{folder.displayName} shared with {folder.sharedWith}<button onclick={()=>unshare(folder.shareId)}>Unshare</button></li>
        <label><input type="radio" name={`perm-${folder.shareId}`} value="View" bind:group={permissionSelections[folder.id]} /> View</label>
        <label><input type="radio" name={`perm-${folder.shareId}`} value="Edit" bind:group={permissionSelections[folder.id]} /> Edit</label>
        <button onclick={() => changePermission(folder.shareId)}>Update Permission</button>
    {/each}
</ul>

{#if fetchErr}
    {fetchErr}
{/if}

