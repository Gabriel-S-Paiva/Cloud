import { endpoints } from '$lib/api';
import { toast } from '$lib/stores/toast.svelte';
import type { SharedFolder, SharedFile } from '$lib/types';

class SharedContentsStore {
	incomingFolders = $state<SharedFolder[]>([]);
	incomingFiles = $state<SharedFile[]>([]);
	outgoingFolders = $state<SharedFolder[]>([]);
	outgoingFiles = $state<SharedFile[]>([]);

	get incomingItems() {
		return [...this.incomingFolders, ...this.incomingFiles];
	}

	get outgoingItems() {
		return [...this.outgoingFolders, ...this.outgoingFiles];
	}

	async loadIncoming() {
		try {
			const res = await endpoints.getShareIncoming();
			this.incomingFolders = res.folders;
			this.incomingFiles = res.files;
		} catch (err) {
			toast.error(err instanceof Error ? err.message : 'Failed to load shared items');
		}
	}

	async loadOutgoing() {
		try {
			const res = await endpoints.getShareOutgoing();
			this.outgoingFolders = res.folders;
			this.outgoingFiles = res.files;
		} catch (err) {
			toast.error(err instanceof Error ? err.message : 'Failed to load shared items');
		}
	}

	async unshare(shareId: number, direction: 'incoming' | 'outgoing') {
		try {
			await endpoints.deleteShare(shareId);
			if (direction === 'incoming') {
				this.incomingFolders = this.incomingFolders.filter((f) => f.shareId !== shareId);
				this.incomingFiles = this.incomingFiles.filter((f) => f.shareId !== shareId);
			} else {
				this.outgoingFolders = this.outgoingFolders.filter((f) => f.shareId !== shareId);
				this.outgoingFiles = this.outgoingFiles.filter((f) => f.shareId !== shareId);
			}
			toast.success('Share removed');
		} catch (err) {
			toast.error(err instanceof Error ? err.message : 'Failed to remove share');
		}
	}

	async updatePermission(shareId: number, permission: 'Edit' | 'View') {
		const item =
			this.outgoingFolders.find((f) => f.shareId === shareId) ??
			this.outgoingFiles.find((f) => f.shareId === shareId);
		if (!item) return;
		const previous = item.permissions;
		try {
			await endpoints.updateShare(shareId, permission);
			item.permissions = permission;
			toast.success('Permission updated');
		} catch (err) {
			toast.error(err instanceof Error ? err.message : 'Failed to update permission');
			item.permissions = previous;
		}
	}
}

export const sharedContents = new SharedContentsStore();