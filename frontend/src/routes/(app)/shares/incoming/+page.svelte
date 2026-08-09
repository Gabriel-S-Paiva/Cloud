<script lang="ts">
	import { onMount } from 'svelte';
	import { sharedContents } from '$lib/stores/sharedContents.svelte';
	import SharedItemRow from '$lib/components/Reactive/SharedItemRow.svelte';

	onMount(() => {
		sharedContents.loadIncoming();
	});
</script>

<div class="p-6">
	<h1 class="mb-6 font-display text-2xl text-text">Shared With Me</h1>

	{#if sharedContents.incomingItems.length === 0}
		<p class="text-sm text-text-muted">Nothing has been shared with you yet.</p>
	{:else}
		<div class="flex flex-col gap-2">
			{#each sharedContents.incomingItems as item (item.shareId)}
				<SharedItemRow {item} direction="incoming" />
			{/each}
		</div>
	{/if}
</div>
