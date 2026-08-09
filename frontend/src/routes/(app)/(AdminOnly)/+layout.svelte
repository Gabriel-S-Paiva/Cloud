<script lang="ts">
	import { auth } from '$lib/stores/auth.svelte';
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	let { children } = $props();

	if (!auth.isAdmin) {
		goto(resolve('/home/'));
	}
</script>

{#if auth.isAdmin}
	{@render children()}
{/if}
