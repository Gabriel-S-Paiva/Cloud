<script lang="ts">
	import type { CloudFile, Folder } from "$lib/types"
	import ItemRenderer from '../items/Renderer.svelte';

	interface Props {
		items: (CloudFile | Folder)[];
	}

	let { items }: Props = $props();

	const isFile = (item: CloudFile | Folder): item is CloudFile => 'size' in item;
</script>

<ul>
	{#each items as item (isFile(item) ? `file-${item.id}` : `folder-${item.id}`)}
		<ItemRenderer {item} variant="list" />
	{/each}
</ul>