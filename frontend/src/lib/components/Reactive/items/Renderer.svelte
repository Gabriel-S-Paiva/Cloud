<script lang="ts">
    import type { CloudFile, Folder } from "$lib/types";
    import FileCard from "./FileCard.svelte";
    import FileRow from "./FileRow.svelte";
    import FolderCard from "./FolderCard.svelte";
    import FolderRow from "./FolderRow.svelte";

    interface Props {
        item: CloudFile | Folder
        variant: 'list' | 'grid'
    }

    let { item, variant }: Props = $props()
    const isFile = (i: CloudFile|Folder): i is CloudFile => 'size' in i
</script>

{#if isFile(item)}
    {#if variant == 'list'}
        <FileRow file={item}/>
    {:else}
        <FileCard file={item}/>
    {/if}
{:else}
    {#if variant == 'list'}
        <FolderRow folder={item}/>
    {:else}
        <FolderCard folder={item}/>
    {/if}
{/if}