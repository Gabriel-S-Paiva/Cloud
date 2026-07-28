<script lang="ts">
	import type { FolderContents } from '$lib/types';
	import ListView from './views/ListView.svelte';
	import GridView from './views/GridView.svelte';
    import Button from '$lib/components/UI/Button/Button.svelte'

	

	let { folders, files }: FolderContents = $props();
	let viewMode = $state<'list' | 'grid'>('list');
	let allItems = $derived([...folders, ...files]);
</script>

<div>
	<div>
		<Button 
			variant={viewMode === 'list' ? 'primary' : 'secondary'} 
			onclick={() => viewMode = 'list'}
		>
			List
		</Button>
		<Button 
			variant={viewMode === 'grid' ? 'primary' : 'secondary'} 
			onclick={() => viewMode = 'grid'}
		>
			Grid
		</Button>
	</div>

    {#if viewMode === 'list'}
		<ListView items={allItems}/>
	{:else}
		<GridView items={allItems} />
	{/if}
</div>