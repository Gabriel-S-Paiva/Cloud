import type {
	User,
	UserSummary,
	Folder,
	SharedFolder,
	CloudFile,
	SharedFile,
	FolderContents,
	SharedContents,
	RegisterRequest,
	PathSegment,
	ToastItem
} from '$lib/types';

// ==========================================
// Users & Summaries
// ==========================================

export const mockUser = (overrides?: Partial<User>): User => ({
	id: 1,
	username: 'johndoe',
	role: 'User',
	quota: 1000000,
	quotaUsed: 500000,
	rootFolderId: 10,
	...overrides
});

export const mockAdminUser = (overrides?: Partial<User>): User =>
	mockUser({ id: 2, username: 'admin', role: 'Admin', ...overrides });

export const mockUserSummary = (overrides?: Partial<UserSummary>): UserSummary => ({
	id: 1,
	username: 'johndoe',
	...overrides
});

// ==========================================
// Folders & Files
// ==========================================

export const mockFolder = (overrides?: Partial<Folder>): Folder => ({
	id: 100,
	displayName: 'Documents',
	ownedBy: 1,
	parentFolder: null,
	...overrides
});

export const mockFile = (overrides?: Partial<CloudFile>): CloudFile => ({
	id: 200,
	displayName: 'Report.pdf',
	ownedBy: 1,
	size: 1024,
	bytesReceived: 1024,
	status: 'Complete',
	contentType: 'application/pdf',
	uploadedAt: 1700000000,
	lastModified: 1700000000,
	parentFolder: null,
	...overrides
});

export const mockFolderContents = (overrides?: Partial<FolderContents>): FolderContents => ({
	folders: [mockFolder({ id: 100 }), mockFolder({ id: 101, displayName: 'Photos' })],
	files: [mockFile({ id: 200 }), mockFile({ id: 201, displayName: 'Image.png' })],
	...overrides
});

// ==========================================
// Shared Folders & Files
// ==========================================

export const mockSharedFolder = (overrides?: Partial<SharedFolder>): SharedFolder => ({
	id: 300,
	displayName: 'Shared Project',
	ownedBy: 2,
	parentFolder: null,
	shareId: 1001,
	sharedWith: 'johndoe',
	ownedByUsername: 'alice',
	permissions: 'View',
	...overrides
});

export const mockSharedFile = (overrides?: Partial<SharedFile>): SharedFile => ({
	id: 400,
	displayName: 'Design.fig',
	ownedBy: 2,
	size: 2048,
	bytesReceived: 2048,
	status: 'Complete',
	contentType: 'application/octet-stream',
	uploadedAt: 1700000000,
	lastModified: 1700000000,
	parentFolder: null,
	shareId: 2002,
	sharedWith: 'johndoe',
	ownedByUsername: 'alice',
	permissions: 'Edit',
	...overrides
});

export const mockSharedContents = (overrides?: Partial<SharedContents>): SharedContents => ({
	folders: [
		mockSharedFolder({ shareId: 1001, displayName: 'Incoming Folder 1' }),
		mockSharedFolder({ shareId: 1002, displayName: 'Incoming Folder 2' })
	],
	files: [
		mockSharedFile({ shareId: 2001, displayName: 'Incoming File 1.pdf' }),
		mockSharedFile({ shareId: 2002, displayName: 'Incoming File 2.png' })
	],
	...overrides
});

// ==========================================
// Utility / Misc
// ==========================================

export const mockRegisterRequest = (overrides?: Partial<RegisterRequest>): RegisterRequest => ({
	id: 1,
	username: 'pending_user',
	status: 'Pending',
	...overrides
});

export const mockPathSegment = (overrides?: Partial<PathSegment>): PathSegment => ({
	id: 100,
	displayName: 'Documents',
	...overrides
});

export const mockNestedPath: PathSegment[] = [
	{ id: 1, displayName: 'Root' },
	{ id: 10, displayName: 'Work' },
	{ id: 100, displayName: 'Projects' }
];

export const mockToastItem = (overrides?: Partial<ToastItem>): ToastItem => ({
	id: '1',
	message: 'Operation successful',
	durationMs: 5000,
	variant: 'success',
	...overrides
});
