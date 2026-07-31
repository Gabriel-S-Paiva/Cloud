<script lang="ts">
	import { driveContents } from '$lib/stores/driveContents.svelte';
	import ListView from './views/ListView.svelte';
	import GridView from './views/GridView.svelte';
	import Button from '$lib/components/UI/Button/Button.svelte';

	let viewMode = $state<'list' | 'grid'>('list');
</script>

<div>
	<div>
		<Button variant={viewMode === 'list' ? 'primary' : 'secondary'} onclick={() => (viewMode = 'list')}>
			List
		</Button>
		<Button variant={viewMode === 'grid' ? 'primary' : 'secondary'} onclick={() => (viewMode = 'grid')}>
			Grid
		</Button>
	</div>

	{#if viewMode === 'list'}
		<ListView items={driveContents.allItems} />
	{:else}
		<GridView items={driveContents.allItems} />
	{/if}
</div>