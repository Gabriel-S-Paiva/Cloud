<script lang="ts">
	import { endpoints } from '$lib/api';
	import { toast } from '$lib/stores/toast.svelte';
	import type { RegisterRequest } from '$lib/types';
	import { onMount } from 'svelte';
	import Button from '$lib/components/UI/Button/Button.svelte';

	let requests = $state<RegisterRequest[] | null>(null);

	onMount(async () => {
		try {
			requests = await endpoints.getRegisterRequests();
		} catch (err) {
			toast.error(err instanceof Error ? err.message : 'Failed to load requests');
		}
	});

	const approve = async (id: number) => {
		try {
			await endpoints.aproveRequest(id);
			requests = await endpoints.getRegisterRequests();
			toast.success('Request approved');
		} catch (err) {
			toast.error(err instanceof Error ? err.message : 'Failed to approve request');
		}
	};

	const reject = async (id: number) => {
		try {
			await endpoints.rejectRequest(id);
			requests = await endpoints.getRegisterRequests();
			toast.success('Request rejected');
		} catch (err) {
			toast.error(err instanceof Error ? err.message : 'Failed to reject request');
		}
	};
</script>

<div class="p-6">
	<h1 class="font-display text-2xl text-text mb-6">Registration Requests</h1>

	{#if requests && requests.length === 0}
		<p class="text-text-muted text-sm">No pending requests.</p>
	{:else if requests}
		<div class="flex flex-col gap-2">
			{#each requests as req (req.id)}
				<div class="flex items-center gap-3 px-3 py-2.5 border border-border rounded-lg bg-surface-raised">
					<span class="flex-1 text-sm text-text">{req.username}</span>
					<span
						class="text-xs px-2 py-0.5 rounded-full border"
						style={req.status === 'Rejected'
							? 'color: var(--danger); border-color: var(--danger);'
							: 'color: var(--text-muted); border-color: var(--border);'}
					>
						{req.status}
					</span>
					<Button variant="secondary" onclick={() => approve(req.id)}>Approve</Button>
					<Button variant="danger" onclick={() => reject(req.id)}>Reject</Button>
				</div>
			{/each}
		</div>
	{/if}
</div>