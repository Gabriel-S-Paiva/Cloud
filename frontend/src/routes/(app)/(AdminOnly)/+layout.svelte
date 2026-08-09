<script lang="ts">
	import { auth } from '$lib/stores/auth.svelte';
	import { goto } from '$app/navigation';
	let { children } = $props();

	if (!auth.isAdmin) {
		goto('/home');
	}
</script>

{#if auth.isAdmin}
	{@render children()}
{/if}
