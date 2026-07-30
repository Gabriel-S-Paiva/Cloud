<script lang="ts">
  import type { Snippet } from 'svelte';

  let {
    type = 'text',
    placeholder = 'Enter text...',
    value = $bindable(''),
    left,
    right
  }: {
    type?: string;
    placeholder?: string;
    value?: string;
    left?: Snippet;
    right?: Snippet;
  } = $props();
</script>

<div class="flex items-center gap-2 h-10 px-3 rounded-md border border-border bg-surface-raised
            focus-within:border-accent transition-colors">
  {#if left}
    <span class="text-text-muted shrink-0">{@render left()}</span>
  {/if}
  <input
    {type}
    {placeholder}
    bind:value
    class="flex-1 bg-transparent outline-none border-none focus:ring-0 text-sm text-text placeholder:text-text-muted"
  />
  {#if right}
    <span class="text-text-muted shrink-0">{@render right()}</span>
  {/if}
</div>