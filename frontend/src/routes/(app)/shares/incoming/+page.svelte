<script lang="ts">
	import { onMount } from 'svelte';
	import { sharedContents } from '$lib/stores/sharedContents.svelte';
	import SharedItemRow from '$lib/components/Reactive/SharedItemRow.svelte';

	onMount(() => {
		sharedContents.loadIncoming();
	});
</script>

<div class="p-6">
	<h1 class="font-display text-2xl text-text mb-6">Shared With Me</h1>

	{#if sharedContents.incomingItems.length === 0}
		<p class="text-text-muted text-sm">Nothing has been shared with you yet.</p>
	{:else}
		<div class="flex flex-col gap-2">
			{#each sharedContents.incomingItems as item (item.shareId)}
				<SharedItemRow {item} direction="incoming" />
			{/each}
		</div>
	{/if}
</div>