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

<button {type} data-variant={variant} class="btn {className}" {...restProps}>
	{@render children?.()}
</button>

<style>
	.btn {
		height: 2.5rem;
		padding: 0 1rem;
		border-radius: 0.375rem;
		font-size: 0.875rem;
		font-weight: 500;
		font-family: var(--font-body);
		transition:
			opacity 0.15s ease,
			background-color 0.15s ease;
	}
	.btn:active {
		transform: scale(0.98);
	}
	.btn:disabled {
		opacity: 0.5;
		pointer-events: none;
	}

	.btn[data-variant='primary'] {
		background: var(--accent);
		color: var(--surface);
	}
	.btn[data-variant='primary']:hover {
		opacity: 0.9;
	}

	.btn[data-variant='secondary'] {
		background: transparent;
		color: var(--text);
		border: 0.5px solid var(--border);
	}
	.btn[data-variant='secondary']:hover {
		background: var(--surface-raised);
		border-color: var(--text-muted);
	}

	.btn[data-variant='danger'] {
		background: transparent;
		color: var(--danger);
		border: 0.5px solid var(--danger);
	}
	.btn[data-variant='danger']:hover {
		background: var(--danger);
		color: var(--surface);
	}
</style>
