<script lang="ts">
	import { endpoints } from "$lib/api";
	import type { SharedContents } from "$lib/types";
	import { onMount } from "svelte";

    let fetchErr = $state<string | null>(null);
    let outgoingShares = $state<SharedContents|null>(null)
    let files = $derived(outgoingShares?.files ?? []); 
    let folders = $derived(outgoingShares?.folders ?? [])

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
</script>
<ul>
    {#each files as file}
        <li>{file.displayName} shared with {file.sharedWith}<button onclick={()=>unshare(file.shareId)}>Unshare</button></li>
    {/each}
</ul>

<ul>
    {#each folders as folder}
        <li>{folder.displayName} shared with {folder.sharedWith}<button onclick={()=>unshare(folder.shareId)}>Unshare</button></li>
    {/each}
</ul>

{#if fetchErr}
    {fetchErr}
{/if}

