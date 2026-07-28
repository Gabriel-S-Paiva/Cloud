<script lang="ts">
	import type { Snippet } from 'svelte';

	interface Props {
		open: boolean;
		onclose: () => void;
		children?: Snippet;
	}

	let { open, onclose, children }: Props = $props();

	// Close on Escape key
	const handleKeyDown = (e: KeyboardEvent) => {
		if (e.key === 'Escape' && open) onclose();
	};
</script>

<svelte:window onkeydown={handleKeyDown} />

{#if open}
	<div 
		class="fixed inset-0 bg-black/50 z-50 flex items-center justify-center p-4"
		onclick={onclose}
		role="presentation"
	>
		<div 
			class="bg-white rounded-lg p-6 max-w-md w-full shadow-xl"
			onclick={(e) => e.stopPropagation()}
			role="dialog"
		>
			{@render children?.()}
		</div>
	</div>
{/if}