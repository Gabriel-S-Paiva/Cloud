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
			requests = [];
		}
	});

	const approve = async (id: number) => {
		try {
			await endpoints.aproveRequest(id);
			requests = requests?.filter((req) => req.id !== id) ?? [];
			toast.success('Request approved');
		} catch (err) {
			toast.error(err instanceof Error ? err.message : 'Failed to approve request');
		}
	};

	const reject = async (id: number) => {
		try {
			await endpoints.rejectRequest(id);
			requests = requests?.filter((req) => req.id !== id) ?? [];
			toast.success('Request rejected');
		} catch (err) {
			toast.error(err instanceof Error ? err.message : 'Failed to reject request');
		}
	};
</script>

<div class="p-6">
	<h1 class="mb-6 font-display text-2xl text-text">Registration Requests</h1>

	{#if requests === null}
		<p class="text-sm text-text-muted">Loading requests...</p>
	{:else if requests.length === 0}
		<p class="text-sm text-text-muted">No pending requests.</p>
	{:else}
		<div class="flex flex-col gap-2">
			{#each requests as req (req.id)}
				<div
					class="flex items-center gap-3 rounded-lg border border-border bg-surface-raised px-3 py-2.5"
				>
					<span class="flex-1 text-sm text-text">{req.username}</span>
					<span
						class="rounded-full border px-2 py-0.5 text-xs"
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
