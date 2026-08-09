<script lang="ts">
	import { onMount } from 'svelte';
	import { sharedContents } from '$lib/stores/sharedContents.svelte';
	import SharedItemRow from '$lib/components/Reactive/SharedItemRow.svelte';

	onMount(() => {
		sharedContents.loadOutgoing();
	});
</script>

<div class="p-6">
	<h1 class="mb-6 font-display text-2xl text-text">Manage Shares</h1>

	{#if sharedContents.outgoingItems.length === 0}
		<p class="text-sm text-text-muted">You haven't shared anything yet.</p>
	{:else}
		<div class="flex flex-col gap-2">
			{#each sharedContents.outgoingItems as item (item.shareId)}
				<SharedItemRow {item} direction="outgoing" />
			{/each}
		</div>
	{/if}
</div>
