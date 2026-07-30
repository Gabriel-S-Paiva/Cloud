<script lang="ts">
	import { toast } from '$lib/stores/toast.svelte';
	import ToastItem from './ToastItem.svelte';
</script>

<div class="fixed bottom-4 right-4 z-50 flex flex-col gap-2 w-80">
	{#each toast.queue as item (item.id)}
		<ToastItem {item} />
	{/each}
</div>