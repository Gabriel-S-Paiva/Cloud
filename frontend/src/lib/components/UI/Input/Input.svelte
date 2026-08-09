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

<div
	class="flex h-10 items-center gap-2 rounded-md border border-border bg-surface-raised px-3
            transition-colors focus-within:border-accent"
>
	{#if left}
		<span class="shrink-0 text-text-muted">{@render left()}</span>
	{/if}
	<input
		{type}
		{placeholder}
		bind:value
		class="flex-1 border-none bg-transparent text-sm text-text outline-none placeholder:text-text-muted focus:ring-0"
	/>
	{#if right}
		<span class="shrink-0 text-text-muted">{@render right()}</span>
	{/if}
</div>
