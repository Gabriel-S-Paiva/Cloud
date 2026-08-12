import { describe, it, expect, vi, beforeEach } from 'vitest';
import type { SharedContents } from '$lib/types';
import { mockSharedFolder, mockSharedFile } from '$lib/mocks/testData';

vi.mock('$lib/api', () => ({
	endpoints: {
		getShareIncoming: vi.fn(),
		getShareOutgoing: vi.fn(),
		deleteShare: vi.fn(),
		updateShare: vi.fn()
	}
}));

vi.mock('$lib/stores/toast.svelte', () => ({
	toast: {
		success: vi.fn(),
		error: vi.fn()
	}
}));

import { sharedContents } from './sharedContents.svelte';
import { endpoints } from '$lib/api';
import { toast } from '$lib/stores/toast.svelte';

function resetStore() {
	sharedContents.incomingFolders = [];
	sharedContents.incomingFiles = [];
	sharedContents.outgoingFolders = [];
	sharedContents.outgoingFiles = [];
}

describe('derived contents', () => {
	beforeEach(resetStore);

	it('incomingItems merges incoming folders and files', () => {
		const folder = mockSharedFolder();
		const file = mockSharedFile();
		sharedContents.incomingFolders = [folder];
		sharedContents.incomingFiles = [file];

		expect(sharedContents.incomingItems).toEqual([folder, file]);
	});

	it('outgoingItems merges outgoing folders and files', () => {
		const folder = mockSharedFolder({ shareId: 5 });
		const file = mockSharedFile({ shareId: 6 });
		sharedContents.outgoingFolders = [folder];
		sharedContents.outgoingFiles = [file];

		expect(sharedContents.outgoingItems).toEqual([folder, file]);
	});
});

describe('loadIncoming', () => {
	beforeEach(() => {
		vi.clearAllMocks();
		resetStore();
	});

	it('populates incoming folders and files on success', async () => {
		const folder = mockSharedFolder();
		const file = mockSharedFile();
		vi.mocked(endpoints.getShareIncoming).mockResolvedValue({
			folders: [folder],
			files: [file]
		});

		await sharedContents.loadIncoming();

		expect(sharedContents.incomingFolders).toEqual([folder]);
		expect(sharedContents.incomingFiles).toEqual([file]);
	});

	it('converts a null folders field in the response to an empty array', async () => {
		const file = mockSharedFile();
		vi.mocked(endpoints.getShareIncoming).mockResolvedValue({
			folders: null,
			files: [file]
		} as unknown as SharedContents);

		await sharedContents.loadIncoming();

		expect(sharedContents.incomingFolders).toEqual([]);
		expect(sharedContents.incomingFiles).toEqual([file]);
	});

	it('leaves existing state untouched and shows an error toast on failure', async () => {
		const existingFolder = mockSharedFolder();
		sharedContents.incomingFolders = [existingFolder];
		vi.mocked(endpoints.getShareIncoming).mockRejectedValue(new Error('Network error'));

		await sharedContents.loadIncoming();

		expect(sharedContents.incomingFolders).toEqual([existingFolder]);
		expect(toast.error).toHaveBeenCalled();
	});
});

describe('loadOutgoing', () => {
	beforeEach(() => {
		vi.clearAllMocks();
		resetStore();
	});

	it('populates outgoing folders and files on success', async () => {
		const folder = mockSharedFolder();
		const file = mockSharedFile();
		vi.mocked(endpoints.getShareOutgoing).mockResolvedValue({
			folders: [folder],
			files: [file]
		});

		await sharedContents.loadOutgoing();

		expect(sharedContents.outgoingFolders).toEqual([folder]);
		expect(sharedContents.outgoingFiles).toEqual([file]);
	});

	it('converts a null files field in the response to an empty array', async () => {
		const folder = mockSharedFolder();
		vi.mocked(endpoints.getShareOutgoing).mockResolvedValue({
			folders: [folder],
			files: null
		} as unknown as SharedContents);

		await sharedContents.loadOutgoing();

		expect(sharedContents.outgoingFolders).toEqual([folder]);
		expect(sharedContents.outgoingFiles).toEqual([]);
	});

	it('leaves existing state untouched and shows an error toast on failure', async () => {
		const existingFile = mockSharedFile();
		sharedContents.outgoingFiles = [existingFile];
		vi.mocked(endpoints.getShareOutgoing).mockRejectedValue(new Error('Network error'));

		await sharedContents.loadOutgoing();

		expect(sharedContents.outgoingFiles).toEqual([existingFile]);
		expect(toast.error).toHaveBeenCalled();
	});
});

describe('unshare', () => {
	beforeEach(() => {
		vi.clearAllMocks();
		resetStore();
	});

	it('removes only the matching incoming item when direction is incoming', async () => {
		const target = mockSharedFolder({ shareId: 1 });
		const other = mockSharedFolder({ shareId: 2 });
		sharedContents.incomingFolders = [target, other];
		sharedContents.outgoingFolders = [mockSharedFolder({ shareId: 1 })]; // same id, different list
		vi.mocked(endpoints.deleteShare).mockResolvedValue(undefined);

		await sharedContents.unshare(1, 'incoming');

		expect(sharedContents.incomingFolders).toEqual([other]);
		expect(sharedContents.outgoingFolders).toHaveLength(1); // untouched
		expect(toast.success).toHaveBeenCalled();
	});

	it('removes only the matching outgoing item when direction is outgoing', async () => {
		const target = mockSharedFile({ shareId: 3 });
		const other = mockSharedFile({ shareId: 4 });
		sharedContents.outgoingFiles = [target, other];
		sharedContents.incomingFiles = [mockSharedFile({ shareId: 3 })]; // same id, different list
		vi.mocked(endpoints.deleteShare).mockResolvedValue(undefined);

		await sharedContents.unshare(3, 'outgoing');

		expect(sharedContents.outgoingFiles).toEqual([other]);
		expect(sharedContents.incomingFiles).toHaveLength(1); // untouched
		expect(toast.success).toHaveBeenCalled();
	});

	it('does not throw when the shareId is not found in either list', async () => {
		sharedContents.incomingFolders = [mockSharedFolder({ shareId: 1 })];
		vi.mocked(endpoints.deleteShare).mockResolvedValue(undefined);

		await expect(sharedContents.unshare(999, 'incoming')).resolves.not.toThrow();
		expect(sharedContents.incomingFolders).toHaveLength(1);
	});

	it('leaves lists untouched and shows an error toast on failure', async () => {
		const folder = mockSharedFolder({ shareId: 1 });
		sharedContents.incomingFolders = [folder];
		vi.mocked(endpoints.deleteShare).mockRejectedValue(new Error('Delete failed'));

		await sharedContents.unshare(1, 'incoming');

		expect(sharedContents.incomingFolders).toEqual([folder]);
		expect(toast.error).toHaveBeenCalled();
	});
});

describe('updatePermission', () => {
	beforeEach(() => {
		vi.clearAllMocks();
		resetStore();
	});

	it('does nothing and makes no network call when the shareId is not found', async () => {
		sharedContents.outgoingFolders = [mockSharedFolder({ shareId: 1 })];

		await sharedContents.updatePermission(999, 'Edit');

		expect(endpoints.updateShare).not.toHaveBeenCalled();
	});

	it('updates the permission on the matching outgoing folder on success', async () => {
		const folder = mockSharedFolder({ shareId: 1, permissions: 'View' });
		sharedContents.outgoingFolders = [folder];
		vi.mocked(endpoints.updateShare).mockResolvedValue(undefined);

		await sharedContents.updatePermission(1, 'Edit');

		expect(sharedContents.outgoingFolders[0].permissions).toBe('Edit');
		expect(toast.success).toHaveBeenCalled();
	});

	it('updates the permission on the matching outgoing file on success', async () => {
		const file = mockSharedFile({ shareId: 2, permissions: 'View' });
		sharedContents.outgoingFiles = [file];
		vi.mocked(endpoints.updateShare).mockResolvedValue(undefined);

		await sharedContents.updatePermission(2, 'Edit');

		expect(sharedContents.outgoingFiles[0].permissions).toBe('Edit');
		expect(toast.success).toHaveBeenCalled();
	});

	it('rolls back to the previous permission and shows an error toast on failure', async () => {
		const folder = mockSharedFolder({ shareId: 1, permissions: 'View' });
		sharedContents.outgoingFolders = [folder];
		vi.mocked(endpoints.updateShare).mockRejectedValue(new Error('Update failed'));

		await sharedContents.updatePermission(1, 'Edit');

		expect(sharedContents.outgoingFolders[0].permissions).toBe('View');
		expect(toast.error).toHaveBeenCalled();
	});
});
