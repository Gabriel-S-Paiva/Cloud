<script lang="ts">
	import type { CloudFile, Folder } from '$lib/types';
	import ItemRenderer from '../items/Renderer.svelte';

	interface Props {
		items: (CloudFile | Folder)[];
	}

	let { items }: Props = $props();

	const isFile = (item: CloudFile | Folder): item is CloudFile => 'size' in item;
</script>

{#if items.length === 0}
	<p class="text-text-muted text-sm p-4">No content</p>
{:else}
	<ul>
		{#each items as item (isFile(item) ? `file-${item.id}` : `folder-${item.id}`)}
			<ItemRenderer {item} variant="list" />
		{/each}
	</ul>
{/if}