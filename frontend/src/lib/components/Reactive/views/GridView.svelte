<script lang="ts">
	import type { CloudFile, Folder } from '$lib/types';
	import ItemRenderer from '../items/Renderer.svelte';

	interface Props {
		items: (CloudFile | Folder)[];
	}

	let { items }: Props = $props();
</script>

<div class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-6 gap-4 p-4">
	{#each items as item (item.id)}
		<ItemRenderer 
			{item} 
			variant="grid"  
		/>
	{/each}
</div>