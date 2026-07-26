<script lang="ts">
	import { endpoints } from '$lib/api';
    import type { RegisterRequest } from '$lib/types'
	import { onMount } from 'svelte';

    let requests = $state<RegisterRequest[] |null>(null)
    let fetchErr = $state<string | null>(null);

    onMount(async () => {
        try{
            requests = await endpoints.getRegisterRequests()
        } catch (err) {
            fetchErr = err instanceof Error ? err.message : 'Folder Fecth Failed';
        }
    })

    const aprove = async (id: number) => {
        try {
            await endpoints.aproveRequest(id)
            requests = await endpoints.getRegisterRequests();
        } catch(err) {
            fetchErr = err instanceof Error ? err.message : 'Folder Fecth Failed';
        }
    }

    const reject = async (id: number) => {
        try {
            await endpoints.rejectRequest(id)
            requests = await endpoints.getRegisterRequests();
        } catch(err) {
            fetchErr = err instanceof Error ? err.message : 'Folder Fecth Failed';
        }
    }
</script>

{#if requests}
    {#each requests as req}
        <p>{req.username}</p>
        <button onclick={() => aprove(req.id)}>Accept</button>
        <button onclick={() => reject(req.id)}>Reject</button>
    {/each}
{/if}
