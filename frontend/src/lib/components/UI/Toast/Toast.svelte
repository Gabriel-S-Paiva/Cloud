<script lang="ts">
	import { toast } from '$lib/stores/toast.svelte';
	import ToastItem from './ToastItem.svelte';
</script>

<div>
	{#each toast.queue as item (item.id)}
		<ToastItem {item} />
	{/each}
</div>