import { endpoints } from '$lib/api';
import { toast } from '$lib/stores/toast.svelte';
import type { CloudFile, Folder, UserSummary } from '$lib/types';

class DriveContents {
	folders = $state<Folder[]>([]);
	files = $state<CloudFile[]>([]);
	sharableUsers = $state<UserSummary[]>([]);

	get allItems() {
		return [...this.folders, ...this.files];
	}

	setContents(folders: Folder[] | null, files: CloudFile[] | null) {
		this.folders = folders ?? [];
		this.files = files ?? [];
	}

	async loadSharableUsers() {
		try {
			this.sharableUsers = await endpoints.getSharableUsers();
		} catch (err) {
			toast.error(err instanceof Error ? err.message : 'Failed to load users');
		}
	}

	async deleteFile(id: number) {
		try {
			await endpoints.deleteFile(id);
			this.files = this.files.filter((f) => f.id !== id);
			toast.success('File deleted');
		} catch (err) {
			toast.error(err instanceof Error ? err.message : 'Failed to delete file');
		}
	}

	async deleteFolder(id: number) {
		try {
			await endpoints.deleteFolder(id);
			this.folders = this.folders.filter((f) => f.id !== id);
			toast.success('Folder deleted');
		} catch (err) {
			toast.error(err instanceof Error ? err.message : 'Failed to delete folder');
		}
	}

	async renameFile(id: number, displayName: string) {
		const file = this.files.find((f) => f.id === id);
		if (!file) return;
		const previous = file.displayName;
		try {
			await endpoints.updateFile(id, displayName);
			file.displayName = displayName;
			toast.success('File renamed');
		} catch (err) {
			toast.error(err instanceof Error ? err.message : 'Failed to rename file');
			file.displayName = previous;
		}
	}

	async renameFolder(id: number, displayName: string) {
		const folder = this.folders.find((f) => f.id === id);
		if (!folder) return;
		const previous = folder.displayName;
		try {
			await endpoints.updateFolder(id, displayName);
			folder.displayName = displayName;
			toast.success('Folder renamed');
		} catch (err) {
			toast.error(err instanceof Error ? err.message : 'Failed to rename folder');
			folder.displayName = previous;
		}
	}

	async shareFile(id: number, targetUserId: number, permission: 'Edit' | 'View') {
		try {
			await endpoints.createShare(id, null, targetUserId, permission);
			toast.success('File shared');
		} catch (err) {
			toast.error(err instanceof Error ? err.message : 'Failed to share file');
		}
	}

	async shareFolder(id: number, targetUserId: number, permission: 'Edit' | 'View') {
		try {
			await endpoints.createShare(null, id, targetUserId, permission);
			toast.success('Folder shared');
		} catch (err) {
			toast.error(err instanceof Error ? err.message : 'Failed to share folder');
		}
	}

	async moveFile(id: number, newParentId: number) {
		try {
			await endpoints.updateFile(id, null, newParentId);
			this.files = this.files.filter((f) => f.id !== id);
			toast.success('File moved');
		} catch (err) {
			toast.error(err instanceof Error ? err.message : 'Failed to move file');
		}
	}

	async moveFolder(id: number, newParentId: number) {
		try {
			await endpoints.updateFolder(id, null, newParentId);
			this.folders = this.folders.filter((f) => f.id !== id);
			toast.success('Folder moved');
		} catch (err) {
			toast.error(err instanceof Error ? err.message : 'Failed to move folder');
		}
	}
}

export const driveContents = new DriveContents();