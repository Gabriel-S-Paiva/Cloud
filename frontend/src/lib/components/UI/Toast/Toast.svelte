<script lang="ts">
	import { toast } from '$lib/stores/toast.svelte';
	import ToastItem from './ToastItem.svelte';
</script>

<div class="fixed right-4 bottom-4 z-50 flex w-80 flex-col gap-2">
	{#each toast.queue as item (item.id)}
		<ToastItem {item} />
	{/each}
</div>
