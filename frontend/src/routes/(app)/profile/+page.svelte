<script lang="ts">
	import { auth } from '$lib/stores/auth.svelte';
	import { goto } from '$app/navigation';
	import Button from '$lib/components/UI/Button/Button.svelte';
	import { LogOut } from '@lucide/svelte';

	const formatBytes = (bytes: number) => {
		if (bytes === 0) return '0 B';
		const units = ['B', 'KB', 'MB', 'GB', 'TB'];
		const i = Math.floor(Math.log(bytes) / Math.log(1024));
		return `${(bytes / 1024 ** i).toFixed(1)} ${units[i]}`;
	};

	const quotaPercent = $derived(
		auth.user ? Math.min(100, (auth.user.quotaUsed / auth.user.quota) * 100) : 0
	);

	const logout = async () => {
		await auth.logout();
		goto('/login');
	};
</script>

<div class="max-w-lg p-6">
	{#if auth.user}
		<div class="mb-8 flex items-center gap-4">
			<div
				class="flex h-14 w-14 shrink-0 items-center justify-center rounded-full bg-accent font-display text-2xl text-surface"
			>
				{auth.user.username[0].toUpperCase()}
			</div>
			<div>
				<h1 class="font-display text-xl text-text">{auth.user.username}</h1>
				<span
					class="mt-1 inline-block rounded-full border px-2 py-0.5 text-xs"
					style={auth.user.role === 'Admin'
						? 'color: var(--accent-secondary); border-color: var(--accent-secondary);'
						: 'color: var(--text-muted); border-color: var(--border);'}
				>
					{auth.user.role}
				</span>
			</div>
		</div>

		<div class="mb-6 rounded-lg border border-border bg-surface-raised p-4">
			<div class="mb-2 flex items-baseline justify-between">
				<span class="text-sm text-text">Storage</span>
				<span class="text-xs text-text-muted">
					{formatBytes(auth.user.quotaUsed)} of {formatBytes(auth.user.quota)}
				</span>
			</div>
			<div class="h-1.5 overflow-hidden rounded-full bg-border">
				<div
					class="h-full transition-[width] duration-300"
					style="width: {quotaPercent}%; background: {quotaPercent > 90
						? 'var(--danger)'
						: 'var(--accent)'};"
				></div>
			</div>
		</div>

		<Button variant="secondary" onclick={logout}>
			<span class="flex items-center gap-2"><LogOut size={16} /> Logout</span>
		</Button>
	{/if}
</div>
