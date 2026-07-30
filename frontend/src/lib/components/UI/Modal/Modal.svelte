<script lang="ts">
	import type { Snippet } from 'svelte';

	interface Props {
		open: boolean;
		onclose: () => void;
		children?: Snippet;
	}

	let { open, onclose, children }: Props = $props();

	const handleKeyDown = (e: KeyboardEvent) => {
		if (e.key === 'Escape' && open) onclose();
	};
</script>

<svelte:window onkeydown={handleKeyDown} />

{#if open}
	<div 
		class="fixed inset-0 bg-black/40 z-50 flex items-center justify-center p-4"
		onclick={onclose}
		role="presentation"
	>
		<div 
			class="bg-surface-raised border border-border rounded-xl p-6 max-w-md w-full"
			style="box-shadow: 0 8px 24px rgb(0 0 0 / 0.12);"
			onclick={(e) => e.stopPropagation()}
			role="dialog"
		>
			{@render children?.()}
		</div>
	</div>
{/if}