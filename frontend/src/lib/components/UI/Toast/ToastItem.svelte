<script lang="ts">
    import type { ToastItem } from '$lib/types';
	import { toast } from '$lib/stores/toast.svelte';

	let { item }: { item: ToastItem } = $props();
	let progress = $state(100);

	$effect(() => {
		const startTime = Date.now();
		const interval = setInterval(() => {
			const elapsed = Date.now() - startTime;
			const remainingRatio = Math.max(0, 1 - elapsed / item.durationMs);
			progress = remainingRatio * 100;

			if (remainingRatio <= 0) {
				clearInterval(interval);
				toast.deleteById(item.id);
			}
		}, 16);
        
		return () => clearInterval(interval);
	});
</script>

<div>
	<p>{item.message}</p>
	<button onclick={() => toast.deleteById(item.id)}>✕</button>
	<progress value={progress} max="100"></progress>
</div>