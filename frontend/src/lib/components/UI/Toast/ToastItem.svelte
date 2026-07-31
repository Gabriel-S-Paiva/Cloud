<script lang="ts">
  import type { ToastItem } from '$lib/types';
  import { toast } from '$lib/stores/toast.svelte';
  import { CircleCheck, Info, TriangleAlert, CircleX } from '@lucide/svelte';

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

  const icons = {
    info: Info,
    success: CircleCheck,
    warning: TriangleAlert,
    error: CircleX
  };

  const Icon = $derived(icons[item.variant]);
</script>

<div
  class="flex gap-3 bg-surface-raised border border-border rounded-lg overflow-hidden relative"
  style="box-shadow: 0 4px 12px rgb(0 0 0 / 0.1);"
>
  <div class="w-1 shrink-0 toast-bar" data-variant={item.variant}></div>

  <div class="flex items-start gap-2 py-3 pr-3 flex-1 min-w-0">
    <Icon size={16} class="toast-icon shrink-0 mt-0.5" data-variant={item.variant} />
    <p class="text-sm text-text flex-1">{item.message}</p>
    <button
      onclick={() => toast.deleteById(item.id)}
      class="text-text-muted hover:text-text text-sm leading-none shrink-0"
      aria-label="Dismiss"
    >
      ✕
    </button>
  </div>

  <div
    class="absolute bottom-0 left-1 right-0 h-0.5 toast-bar transition-[width] duration-75 ease-linear"
    data-variant={item.variant}
    style="width: {progress}%"
  ></div>
</div>

<style>
  .toast-bar[data-variant='success'] { background: var(--accent); }
  .toast-bar[data-variant='warning'] { background: var(--accent-secondary); }
  .toast-bar[data-variant='error']   { background: var(--danger); }
  .toast-bar[data-variant='info']    { background: var(--info); }

  :global(.toast-icon[data-variant='success']) { color: var(--accent); }
  :global(.toast-icon[data-variant='warning']) { color: var(--accent-secondary); }
  :global(.toast-icon[data-variant='error'])   { color: var(--danger); }
  :global(.toast-icon[data-variant='info'])    { color: var(--info); }
</style>