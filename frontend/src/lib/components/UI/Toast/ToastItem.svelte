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

<div class="bg-surface-raised border border-border rounded-lg p-3 overflow-hidden relative"
     style="box-shadow: 0 4px 12px rgb(0 0 0 / 0.1);">
	<div class="flex items-start justify-between gap-3">
		<p class="text-sm text-text flex-1">{item.message}</p>
		<button
			onclick={() => toast.deleteById(item.id)}
			class="text-text-muted hover:text-text text-sm leading-none"
			aria-label="Dismiss"
		>
			✕
		</button>
	</div>
	<div class="absolute bottom-0 left-0 h-0.5 bg-accent transition-[width] duration-75 ease-linear"
	     style="width: {progress}%"></div>
</div>