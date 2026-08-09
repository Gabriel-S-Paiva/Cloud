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
		class="fixed inset-0 z-50 flex items-center justify-center bg-black/40 p-4"
		onclick={onclose}
		role="presentation"
	>
		<!-- svelte-ignore a11y_click_events_have_key_events -->
		<div
			class="w-full max-w-md rounded-xl border border-border bg-surface-raised p-6"
			style="box-shadow: 0 8px 24px rgb(0 0 0 / 0.12);"
			role="dialog"
			aria-modal="true"
			tabindex="-1"
			onclick={(e) => e.stopPropagation()}
		>
			{@render children?.()}
		</div>
	</div>
{/if}
