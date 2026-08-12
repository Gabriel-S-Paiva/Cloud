import { describe, it, expect, vi, beforeEach } from 'vitest';
import { mockFile, mockFolder } from '$lib/mocks/testData';
import type { UserSummary } from '$lib/types';

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
import { endpoints } from '$lib/api';

describe('setContents', () => {
	beforeEach(() => {
		vi.clearAllMocks();
		driveContents.setContents([], []);
	});

	it('converts null folders into an empty array', () => {
		driveContents.setContents(null, []);
		expect(driveContents.folders).toEqual([]);
	});

	it('converts null files into an empty array', () => {
		driveContents.setContents([], null);
		expect(driveContents.files).toEqual([]);
	});

	it('allItems combines folders and files into one array', () => {
		const folder = mockFolder();
		const file = mockFile();

		driveContents.setContents([folder], [file]);

		expect(driveContents.allItems).toEqual([folder, file]);
	});
});

describe('loadSharableUsers', () => {
	beforeEach(() => {
		vi.clearAllMocks();
		driveContents.setContents([], []);
		driveContents.sharableUsers = [];
	});

	it('populates sharableUsers from a success response', async () => {
		const userList: UserSummary[] = [
			{ id: 1, username: 'user1' },
			{ id: 2, username: 'user2' }
		];
		vi.mocked(endpoints.getSharableUsers).mockResolvedValue(userList);

		await driveContents.loadSharableUsers();

		expect(driveContents.sharableUsers).toEqual(userList);
	});

	it('keeps the existing list and shows an error toast when the reload fails', async () => {
		const userList: UserSummary[] = [{ id: 1, username: 'user1' }];
		driveContents.sharableUsers = userList;
		vi.mocked(endpoints.getSharableUsers).mockRejectedValue(new Error('Error fetching users'));

		await driveContents.loadSharableUsers();

		expect(driveContents.sharableUsers).toEqual(userList);
		expect(toast.error).toHaveBeenCalled();
	});
});

describe('deleteFile', () => {
	beforeEach(() => {
		vi.clearAllMocks();
		driveContents.setContents([], []);
	});

	it('removes the file from driveContents.files on success', async () => {
		const files = [mockFile({ id: 1 }), mockFile({ id: 2 })];
		driveContents.setContents([], files);
		vi.mocked(endpoints.deleteFile).mockResolvedValue(undefined);

		await driveContents.deleteFile(1);

		expect(driveContents.files).toEqual([files[1]]);
		expect(toast.success).toHaveBeenCalled();
	});

	it('keeps the file in driveContents.files on failure', async () => {
		const files = [mockFile({ id: 1 }), mockFile({ id: 2 })];
		driveContents.setContents([], files);
		vi.mocked(endpoints.deleteFile).mockRejectedValue(new Error('Error deleting file'));

		await driveContents.deleteFile(1);

		expect(driveContents.files).toEqual(files);
		expect(toast.error).toHaveBeenCalled();
	});
});

describe('deleteFolder', () => {
	beforeEach(() => {
		vi.clearAllMocks();
		driveContents.setContents([], []);
	});

	it('removes the folder from driveContents.folders on success', async () => {
		const folders = [mockFolder({ id: 1 }), mockFolder({ id: 2 })];
		driveContents.setContents(folders, []);
		vi.mocked(endpoints.deleteFolder).mockResolvedValue(undefined);

		await driveContents.deleteFolder(1);

		expect(driveContents.folders).toEqual([folders[1]]);
		expect(toast.success).toHaveBeenCalled();
	});

	it('keeps the folder in driveContents.folders on failure', async () => {
		const folders = [mockFolder({ id: 1 }), mockFolder({ id: 2 })];
		driveContents.setContents(folders, []);
		vi.mocked(endpoints.deleteFolder).mockRejectedValue(new Error('Error deleting folder'));

		await driveContents.deleteFolder(1);

		expect(driveContents.folders).toEqual(folders);
		expect(toast.error).toHaveBeenCalled();
	});
});

describe('renameFile', () => {
	beforeEach(() => {
		vi.clearAllMocks();
		driveContents.setContents([], []);
	});

	it('returns early and makes no network call when the file id is not found', async () => {
		const files = [mockFile({ id: 1 }), mockFile({ id: 2 })];
		driveContents.setContents([], files);

		await driveContents.renameFile(999, 'New Name');

		expect(endpoints.updateFile).not.toHaveBeenCalled();
		expect(driveContents.files).toEqual(files);
	});

	it('updates the display name and shows a success toast', async () => {
		const file = mockFile({ id: 1 });
		driveContents.setContents([], [file]);
		vi.mocked(endpoints.updateFile).mockResolvedValue(undefined);

		await driveContents.renameFile(1, 'New Name');

		expect(endpoints.updateFile).toHaveBeenCalledTimes(1);
		expect(driveContents.files[0].displayName).toBe('New Name');
		expect(toast.success).toHaveBeenCalled();
	});

	it('rolls back the display name and shows an error toast on failure', async () => {
		const file = mockFile({ id: 1, displayName: 'Original Name' });
		driveContents.setContents([], [file]);
		vi.mocked(endpoints.updateFile).mockRejectedValue(new Error('Error renaming file'));

		await driveContents.renameFile(1, 'New Name');

		expect(driveContents.files[0].displayName).toBe('Original Name');
		expect(toast.error).toHaveBeenCalledTimes(1);
	});
});

describe('renameFolder', () => {
	beforeEach(() => {
		vi.clearAllMocks();
		driveContents.setContents([], []);
	});

	it('returns early and makes no network call when the folder id is not found', async () => {
		const folders = [mockFolder({ id: 1 }), mockFolder({ id: 2 })];
		driveContents.setContents(folders, []);

		await driveContents.renameFolder(999, 'New Name');

		expect(endpoints.updateFolder).not.toHaveBeenCalled();
		expect(driveContents.folders).toEqual(folders);
	});

	it('updates the display name and shows a success toast', async () => {
		const folder = mockFolder({ id: 1 });
		driveContents.setContents([folder], []);
		vi.mocked(endpoints.updateFolder).mockResolvedValue(undefined);

		await driveContents.renameFolder(1, 'New Name');

		expect(endpoints.updateFolder).toHaveBeenCalledTimes(1);
		expect(driveContents.folders[0].displayName).toBe('New Name');
		expect(toast.success).toHaveBeenCalled();
	});

	it('rolls back the display name and shows an error toast on failure', async () => {
		const folder = mockFolder({ id: 1, displayName: 'Original Name' });
		driveContents.setContents([folder], []);
		vi.mocked(endpoints.updateFolder).mockRejectedValue(new Error('Error renaming folder'));

		await driveContents.renameFolder(1, 'New Name');

		expect(driveContents.folders[0].displayName).toBe('Original Name');
		expect(toast.error).toHaveBeenCalledTimes(1);
	});
});

describe('shareFile', () => {
	beforeEach(() => {
		vi.clearAllMocks();
		driveContents.setContents([], []);
	});

	it('calls createShare with the file id and shows a success toast', async () => {
		vi.mocked(endpoints.createShare).mockResolvedValue({ id: 1 });

		await driveContents.shareFile(1, 1, 'Edit');

		expect(endpoints.createShare).toHaveBeenCalledWith(1, null, 1, 'Edit');
		expect(toast.success).toHaveBeenCalled();
	});

	it('shows an error toast on failure', async () => {
		vi.mocked(endpoints.createShare).mockRejectedValue(new Error('Error creating share'));

		await driveContents.shareFile(1, 1, 'Edit');

		expect(toast.error).toHaveBeenCalled();
	});
});

describe('shareFolder', () => {
	beforeEach(() => {
		vi.clearAllMocks();
		driveContents.setContents([], []);
	});

	it('calls createShare with the folder id and shows a success toast', async () => {
		vi.mocked(endpoints.createShare).mockResolvedValue({ id: 1 });

		await driveContents.shareFolder(1, 1, 'Edit');

		expect(endpoints.createShare).toHaveBeenCalledWith(null, 1, 1, 'Edit');
		expect(toast.success).toHaveBeenCalled();
	});

	it('shows an error toast on failure', async () => {
		vi.mocked(endpoints.createShare).mockRejectedValue(new Error('Error creating share'));

		await driveContents.shareFolder(1, 1, 'Edit');

		expect(toast.error).toHaveBeenCalled();
	});
});

describe('moveFile', () => {
	beforeEach(() => {
		vi.clearAllMocks();
		driveContents.setContents([], []);
	});

	it('removes the file from the current list and shows a success toast', async () => {
		const file = mockFile({ id: 1 });
		driveContents.setContents([], [file]);
		vi.mocked(endpoints.updateFile).mockResolvedValue(undefined);

		await driveContents.moveFile(1, 1);

		expect(driveContents.files).toEqual([]);
		expect(toast.success).toHaveBeenCalled();
	});

	it('keeps the file in the list and shows an error toast on failure', async () => {
		const file = mockFile({ id: 1 });
		driveContents.setContents([], [file]);
		vi.mocked(endpoints.updateFile).mockRejectedValue(new Error('Error moving file'));

		await driveContents.moveFile(1, 1);

		expect(driveContents.files).toEqual([file]);
		expect(toast.error).toHaveBeenCalled();
	});
});

describe('moveFolder', () => {
	beforeEach(() => {
		vi.clearAllMocks();
		driveContents.setContents([], []);
	});

	it('removes the folder from the current list and shows a success toast', async () => {
		const folder = mockFolder({ id: 1 });
		driveContents.setContents([folder], []);
		vi.mocked(endpoints.updateFolder).mockResolvedValue(undefined);

		await driveContents.moveFolder(1, 1);

		expect(driveContents.folders).toEqual([]);
		expect(toast.success).toHaveBeenCalled();
	});

	it('keeps the folder in the list and shows an error toast on failure', async () => {
		const folder = mockFolder({ id: 1 });
		driveContents.setContents([folder], []);
		vi.mocked(endpoints.updateFolder).mockRejectedValue(new Error('Error moving folder'));

		await driveContents.moveFolder(1, 1);

		expect(driveContents.folders).toEqual([folder]);
		expect(toast.error).toHaveBeenCalled();
	});
});
