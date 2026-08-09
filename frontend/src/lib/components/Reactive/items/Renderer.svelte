<script lang="ts">
	import type { CloudFile, Folder } from '$lib/types';
	import FileItem from './FileItem.svelte';
	import FolderItem from './FolderItem.svelte';

	interface Props {
		item: CloudFile | Folder;
		variant: 'list' | 'grid';
	}

	let { item, variant }: Props = $props();
	const isFile = (i: CloudFile | Folder): i is CloudFile => 'size' in i;
</script>

{#if isFile(item)}
	<FileItem file={item} {variant} />
{:else}
	<FolderItem folder={item} {variant} />
{/if}
