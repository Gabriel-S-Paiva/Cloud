<script lang="ts">
	import { endpoints } from "$lib/api";
	import type { FolderContents } from "$lib/types";
	import { onMount } from "svelte";

    let fetchErr = $state<string | null>(null);
    let incomingShares = $state<FolderContents|null>(null)
    let files = $derived(incomingShares?.files ?? []); 

    onMount( async () => {
        try{
            incomingShares = await endpoints.getShareIncoming()
        } catch(err) {
            fetchErr = err instanceof Error ? err.message : 'Folder Fecth Failed';
        }
    })
</script>

{#each files as file}
    {file.displayName}
{/each}
