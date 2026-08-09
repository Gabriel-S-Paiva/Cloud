<script lang="ts">
	import { sharedContents } from '$lib/stores/sharedContents.svelte';
	import { endpoints } from '$lib/api';
	import { navigation } from '$lib/stores/navigation.svelte';
	import { goto } from '$app/navigation';
	import { toast } from '$lib/stores/toast.svelte';
	import Button from '$lib/components/UI/Button/Button.svelte';
	import Modal from '$lib/components/UI/Modal/Modal.svelte';
	import { Folder as FolderIcon, FileText } from '@lucide/svelte';
	import type { SharedFolder, SharedFile } from '$lib/types';

	let { item, direction }: { item: SharedFolder | SharedFile; direction: 'incoming' | 'outgoing' } =
		$props();

	const isFile = (i: SharedFolder | SharedFile): i is SharedFile => 'size' in i;
	const otherUser = direction === 'incoming' ? item.ownedByUsername : item.sharedWith;

	let confirmOpen = $state(false);

	const openFile = async (): Promise<void> => {
		try {
			const blob = await endpoints.getFileContent(item.id);
			const url = URL.createObjectURL(blob);
			window.open(url, '_blank');
			setTimeout(() => URL.revokeObjectURL(url), 10000);
		} catch (err) {
			toast.error(err instanceof Error ? err.message : 'Failed to open file');
		}
	};

	const enterFolder = async (): Promise<void> => {
		try {
			navigation.reset();
			navigation.enter({ id: item.id, displayName: item.displayName });
			goto(`/home/${navigation.urlPath}`);
		} catch (err) {
			toast.error(err instanceof Error ? err.message : 'Failed to open folder');
		}
	};

	const handleOpen = () => (isFile(item) ? openFile() : enterFolder());
</script>

<div
	class="flex cursor-pointer items-center gap-3 rounded-lg border border-border bg-surface-raised px-3 py-2.5 transition-colors hover:border-accent/40"
	ondblclick={handleOpen}
	role="button"
	tabindex="0"
>
	{#if isFile(item)}
		<FileText size={18} class="shrink-0 text-accent-secondary" />
	{:else}
		<FolderIcon size={18} class="shrink-0 text-accent" />
	{/if}

	<span class="flex-1 truncate text-sm text-text">{item.displayName}</span>

	<span class="hidden shrink-0 text-xs text-text-muted sm:inline">
		{direction === 'incoming' ? 'from' : 'with'}
		{otherUser}
	</span>

	{#if direction === 'outgoing'}
		<select
			value={item.permissions}
			onchange={(e) =>
				sharedContents.updatePermission(item.shareId, e.currentTarget.value as 'Edit' | 'View')}
			class="h-8 shrink-0 rounded-md border border-border bg-surface px-2 text-xs text-text"
		>
			<option value="View">View</option>
			<option value="Edit">Edit</option>
		</select>
	{:else}
		<span
			class="shrink-0 rounded-full border px-2 py-1 text-xs"
			style={item.permissions === 'Edit'
				? 'color: var(--accent-secondary); border-color: var(--accent-secondary);'
				: 'color: var(--text-muted); border-color: var(--border);'}
		>
			{item.permissions}
		</span>
	{/if}

	<Button variant="secondary" onclick={() => (confirmOpen = true)}>Unshare</Button>
</div>

<Modal open={confirmOpen} onclose={() => (confirmOpen = false)}>
	<div class="flex flex-col gap-3">
		<h3 class="font-display text-lg text-text">Remove share?</h3>
		<p class="text-sm text-text">
			{#if direction === 'incoming'}
				You'll lose access to <strong>{item.displayName}</strong>.
			{:else}
				<strong>{item.sharedWith}</strong> will lose access to <strong>{item.displayName}</strong>.
			{/if}
		</p>
		<div class="flex gap-2">
			<Button variant="secondary" onclick={() => (confirmOpen = false)}>Cancel</Button>
			<Button
				variant="danger"
				onclick={() => {
					sharedContents.unshare(item.shareId, direction);
					confirmOpen = false;
				}}
			>
				Confirm
			</Button>
		</div>
	</div>
</Modal>
