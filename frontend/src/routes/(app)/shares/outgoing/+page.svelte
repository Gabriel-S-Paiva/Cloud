<script lang="ts">
	import { endpoints } from "$lib/api";
	import type { FolderContents } from "$lib/types";
	import { onMount } from "svelte";

    let fetchErr = $state<string | null>(null);
    let outgoingShares = $state<FolderContents|null>(null)
    let files = $derived(outgoingShares?.files ?? []); 

    onMount( async () => {
        try{
            outgoingShares = await endpoints.getShareOutgoing()
        } catch(err) {
            fetchErr = err instanceof Error ? err.message : 'Folder Fecth Failed';
        }
    })
</script>

{#each files as file}
    {file.displayName}
{/each}
