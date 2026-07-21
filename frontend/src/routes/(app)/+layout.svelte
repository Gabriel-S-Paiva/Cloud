<script lang="ts">
  import { auth } from '$lib/stores/auth.svelte';
  import { goto } from '$app/navigation';
  import { onMount } from 'svelte';

  let { children } = $props();
  let checked = $state(false);

  onMount(async () => {
    await auth.checkSession();
    if (!auth.isLoggedIn) {
      goto('/login');
    }
    checked = true;
  });
</script>

{#if checked}
  {@render children()}
{/if}