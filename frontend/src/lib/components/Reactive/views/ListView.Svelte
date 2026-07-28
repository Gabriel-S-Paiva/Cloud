<script lang="ts">
	import type { CloudFile, Folder } from "$lib/types"
	import ItemRenderer from '../items/Renderer.svelte';

	interface Props {
		items: (CloudFile | Folder)[];
	}

	let { items }: Props = $props();
</script>

<ul>
    {#each items as item (item.id)}
        <ItemRenderer {item} variant="list"/>
    {/each}
</ul>