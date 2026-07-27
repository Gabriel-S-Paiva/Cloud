<script lang="ts">
    import { auth } from "$lib/stores/auth.svelte";
</script>

{#if auth.user}
    {#each Object.entries(auth.user) as [key, value]}
        {key}:{value}
    {/each}
{/if}