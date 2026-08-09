<script lang="ts">
	import { auth } from '$lib/stores/auth.svelte';
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	import Sidebar from '$lib/components/UI/Nav/Sidebar.svelte';
	import Toast from '$lib/components/UI/Toast/Toast.svelte';

	let { children } = $props();
	let checked = $state(false);

	onMount(async () => {
		await auth.checkSession();
		checked = true;
	});

	$effect(() => {
		if (checked && !auth.isLoggedIn) {
			goto('/login');
		}
	});
</script>

{#if checked}
	<div class="flex min-h-screen bg-surface">
		<Sidebar />
		<main class="flex-1 overflow-y-auto">
			{@render children()}
		</main>
	</div>
	<Toast />
{/if}
