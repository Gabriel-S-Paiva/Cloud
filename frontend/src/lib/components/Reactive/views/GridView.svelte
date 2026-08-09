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
	<p class="p-4 text-sm text-text-muted">No content</p>
{:else}
	<div class="grid grid-cols-2 gap-4 p-4 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-6">
		{#each items as item (isFile(item) ? `file-${item.id}` : `folder-${item.id}`)}
			<ItemRenderer {item} variant="grid" />
		{/each}
	</div>
{/if}
