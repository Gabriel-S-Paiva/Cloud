import { describe, it, expect, vi, beforeEach } from 'vitest';

vi.mock('$lib/api', () => ({
	endpoints: {
		deleteFile: vi.fn(),
		deleteFolder: vi.fn(),
		updateFile: vi.fn(),
		updateFolder: vi.fn(),
		createShare: vi.fn(),
		getSharableUsers: vi.fn()
	}
}));

vi.mock('$lib/stores/toast.svelte', () => ({
	toast: {
		success: vi.fn(),
		error: vi.fn()
	}
}));

import { driveContents } from './driveContents.svelte';
import { toast } from '$lib/stores/toast.svelte';
import type { CloudFile, Folder, UserSummary } from '$lib/types';
import { endpoints } from '$lib/api';

describe('setContents', () => {
	beforeEach(() => {
		vi.clearAllMocks();
		driveContents.setContents([], []);
	});

	it('Converts null folders into empty array', async () => {
		driveContents.setContents(null, []);
		expect(driveContents.folders).toEqual([]);
	});

	it('Converts null files into empty array', async () => {
		driveContents.setContents([], null);
		expect(driveContents.files).toEqual([]);
	});

	it('allItems combines folder and files into one array', async () => {
		const mockFile: CloudFile = {
			id: 1,
			displayName: 'Test Files',
			ownedBy: 1,
			size: 15,
			bytesReceived: 15,
			status: 'Complete',
			contentType: 'type',
			uploadedAt: 1,
			lastModified: 1,
			parentFolder: null
		};
		const mockFolder: Folder = {
			id: 1,
			displayName: 'Test Files',
			ownedBy: 1,
			parentFolder: null
		};

		driveContents.setContents([mockFolder], [mockFile]);

		expect(driveContents.allItems).toEqual([mockFolder, mockFile]);
	});
});

describe('loadSharableUsers', () => {
	beforeEach(() => {
		vi.clearAllMocks();
		driveContents.setContents([], []);
		driveContents.sharableUsers = [];
	});

	it('Populates sharableUsers from a success response', async () => {
		const userList: UserSummary[] = [
			{
				id: 1,
				username: 'user1'
			},
			{
				id: 2,
				username: 'user2'
			}
		];
		vi.mocked(endpoints.getSharableUsers).mockResolvedValue(userList);

		await driveContents.loadSharableUsers();

		expect(driveContents.sharableUsers).toBe(userList);
	});

	it('sharableUsers keeps their list from error response', async () => {
		const userList: UserSummary[] = [
			{
				id: 1,
				username: 'user1'
			},
			{
				id: 2,
				username: 'user2'
			}
		];
		driveContents.sharableUsers = userList;

		vi.mocked(endpoints.getSharableUsers).mockRejectedValue(new Error('Error Fetching Users'));

		await driveContents.loadSharableUsers();

		expect(driveContents.sharableUsers).toBe(userList);
		expect(toast.error).toHaveBeenCalled();
	});
});

describe('deleteFile', () => {
	beforeEach(() => {
		vi.clearAllMocks();
		driveContents.setContents([], []);
	});

	it('Removes files from driveContents.files on sucess', async () => {
		const fileList: CloudFile[] = [
			{
				id: 1,
				displayName: 'Test Files',
				ownedBy: 1,
				size: 15,
				bytesReceived: 15,
				status: 'Complete',
				contentType: 'type',
				uploadedAt: 1,
				lastModified: 1,
				parentFolder: null
			},
			{
				id: 2,
				displayName: 'Test Files',
				ownedBy: 1,
				size: 15,
				bytesReceived: 15,
				status: 'Complete',
				contentType: 'type',
				uploadedAt: 1,
				lastModified: 1,
				parentFolder: null
			}
		];
		driveContents.setContents([], fileList);

		vi.mocked(endpoints.deleteFile).mockResolvedValue(undefined);
		await driveContents.deleteFile(1);

		expect(driveContents.files).toEqual([fileList[1]]);
		expect(toast.success).toHaveBeenCalled();
	});

	it('Keeps files from driveContents.files on error', async () => {
		const fileList: CloudFile[] = [
			{
				id: 1,
				displayName: 'Test Files',
				ownedBy: 1,
				size: 15,
				bytesReceived: 15,
				status: 'Complete',
				contentType: 'type',
				uploadedAt: 1,
				lastModified: 1,
				parentFolder: null
			},
			{
				id: 2,
				displayName: 'Test Files',
				ownedBy: 1,
				size: 15,
				bytesReceived: 15,
				status: 'Complete',
				contentType: 'type',
				uploadedAt: 1,
				lastModified: 1,
				parentFolder: null
			}
		];
		driveContents.setContents([], fileList);

		vi.mocked(endpoints.deleteFile).mockRejectedValue(new Error('Error Deleting File'));
		await driveContents.deleteFile(1);

		expect(driveContents.files).toEqual(fileList);
		expect(toast.error).toHaveBeenCalled();
	});
});

describe('deleteFile', () => {
	beforeEach(() => {
		vi.clearAllMocks();
		driveContents.setContents([], []);
	});

	it('Removes folder from driveContents.folder on sucess', async () => {
		const folderList: Folder[] = [
			{
				id: 1,
				displayName: 'Test Folder',
				ownedBy: 1,
				parentFolder: null
			},
			{
				id: 2,
				displayName: 'Test Folder',
				ownedBy: 1,
				parentFolder: null
			}
		];
		driveContents.setContents(folderList, null);

		vi.mocked(endpoints.deleteFolder).mockResolvedValue(undefined);
		await driveContents.deleteFolder(1);

		expect(driveContents.folders).toEqual([folderList[1]]);
		expect(toast.success).toHaveBeenCalled();
	});

	it('Keeps folders from driveContents.folders on error', async () => {
		const folderList: Folder[] = [
			{
				id: 1,
				displayName: 'Test Folder',
				ownedBy: 1,
				parentFolder: null
			},
			{
				id: 2,
				displayName: 'Test Folder',
				ownedBy: 1,
				parentFolder: null
			}
		];
		driveContents.setContents(folderList, null);

		vi.mocked(endpoints.deleteFolder).mockRejectedValue(new Error('Error Deleting Folder'));
		await driveContents.deleteFile(1);

		expect(driveContents.folders).toEqual(folderList);
		expect(toast.error).toHaveBeenCalled();
	});
});

describe('rename file', () => {
	beforeEach(() => {
		vi.clearAllMocks();
		driveContents.setContents([], []);
	});

	it('returns if file id isnt found', async () => {
		const fileList: CloudFile[] = [
			{
				id: 1,
				displayName: 'Test Files',
				ownedBy: 1,
				size: 15,
				bytesReceived: 15,
				status: 'Complete',
				contentType: 'type',
				uploadedAt: 1,
				lastModified: 1,
				parentFolder: null
			},
			{
				id: 2,
				displayName: 'Test Files',
				ownedBy: 1,
				size: 15,
				bytesReceived: 15,
				status: 'Complete',
				contentType: 'type',
				uploadedAt: 1,
				lastModified: 1,
				parentFolder: null
			}
		];
		driveContents.setContents([], fileList);

		driveContents.renameFile(3, 'New Name');

		expect(endpoints.updateFile).toHaveBeenCalledTimes(0);
		expect(driveContents.files).toEqual(fileList);
	});

	it('Sucess Rename and Toast', async () => {
		const file: CloudFile = {
			id: 1,
			displayName: 'Test Files',
			ownedBy: 1,
			size: 15,
			bytesReceived: 15,
			status: 'Complete',
			contentType: 'type',
			uploadedAt: 1,
			lastModified: 1,
			parentFolder: null
		};

		driveContents.setContents([], [file]);
		vi.mocked(endpoints.updateFile).mockResolvedValue(undefined);

		await driveContents.renameFile(1, 'New Name');

		expect(endpoints.updateFile).toHaveBeenCalledTimes(1);
		expect(driveContents.files[0]).toEqual(
			expect.objectContaining({
				id: 1,
				displayName: 'New Name'
			})
		);
		expect(toast.success).toHaveBeenCalled();
	});

	it('Failure Roolsback name and Toast', async () => {
		const file: CloudFile = {
			id: 1,
			displayName: 'Test Files',
			ownedBy: 1,
			size: 15,
			bytesReceived: 15,
			status: 'Complete',
			contentType: 'type',
			uploadedAt: 1,
			lastModified: 1,
			parentFolder: null
		};

		driveContents.setContents([], [file]);
		vi.mocked(endpoints.updateFile).mockRejectedValue(new Error('Error renaming File'));

		await driveContents.renameFile(1, 'New Name');

		expect(endpoints.updateFile).toHaveBeenCalledTimes(1);
		expect(driveContents.files[0]).toEqual(expect.objectContaining(file));
		expect(toast.error).toHaveBeenCalledTimes(1);
	});
});

describe('rename folder', () => {
	beforeEach(() => {
		vi.clearAllMocks();
		driveContents.setContents([], []);
	});

	it('returns if folder id isnt found', async () => {
		const folderList: Folder[] = [
			{
				id: 1,
				displayName: 'Test Folder',
				ownedBy: 1,
				parentFolder: null
			},
			{
				id: 2,
				displayName: 'Test Folder',
				ownedBy: 1,
				parentFolder: null
			}
		];
		driveContents.setContents(folderList, []);

		driveContents.renameFolder(3, 'New Name');

		expect(endpoints.updateFolder).toHaveBeenCalledTimes(0);
		expect(driveContents.folders).toEqual(folderList);
	});

	it('Sucess Rename and Toast', async () => {
		const folder: Folder = {
			id: 1,
			displayName: 'Test Files',
			ownedBy: 1,
			parentFolder: null
		};

		driveContents.setContents([folder], []);
		vi.mocked(endpoints.updateFolder).mockResolvedValue(undefined);

		await driveContents.renameFolder(1, 'New Name');

		expect(endpoints.updateFolder).toHaveBeenCalledTimes(1);
		expect(driveContents.folders[0]).toEqual(
			expect.objectContaining({
				id: 1,
				displayName: 'New Name'
			})
		);
		expect(toast.success).toHaveBeenCalled();
	});

	it('Failure Roolback name and Toast', async () => {
		const folder: Folder = {
			id: 1,
			displayName: 'Test Folder',
			ownedBy: 1,
			parentFolder: null
		};

		driveContents.setContents([folder], []);
		vi.mocked(endpoints.updateFolder).mockRejectedValue(new Error('Error renaming Folder'));

		await driveContents.renameFolder(1, 'New Name');

		expect(endpoints.updateFolder).toHaveBeenCalledTimes(1);
		expect(driveContents.folders[0]).toEqual(expect.objectContaining(folder));
		expect(toast.error).toHaveBeenCalledTimes(1);
	});
});
