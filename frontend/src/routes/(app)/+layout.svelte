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


<nav class="w-20">
  <button onclick={() => goto('/home')}>Home</button>
  <button onclick={() => goto('/shares/incoming')}>Shared With Me</button>
  <button onclick={() => goto('/shares/outgoing')}>Manage Shares</button>
  <button onclick={() => goto('/profile')}>Profile</button>
</nav>
{#if checked}
  {@render children()}
{/if}