<script lang="ts">
	import type { Snippet } from 'svelte';
	import type { HTMLButtonAttributes } from 'svelte/elements';

	interface Props extends HTMLButtonAttributes {
		variant?: 'primary' | 'secondary' | 'danger';
		children?: Snippet;
	}

	let { 
		variant = 'primary', 
		children, 
		type = 'button', 
		class: className = '', 
		...restProps 
	}: Props = $props();
</script>

<button
	{type}
	data-variant={variant}
	class={className}
	{...restProps}
>
	{@render children?.()}
</button>