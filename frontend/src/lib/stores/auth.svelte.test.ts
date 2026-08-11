import type { User } from '$lib/types';
import { describe, it, expect, vi, beforeEach } from 'vitest';

vi.mock('$lib/api', () => ({
	endpoints: {
		login: vi.fn(),
		logout: vi.fn(),
		getMe: vi.fn(),
		register: vi.fn()
	}
}));

const { endpoints } = await import('$lib/api');
const { auth } = await import('./auth.svelte');

describe('login', () => {
	beforeEach(() => {
		vi.clearAllMocks();
		auth.user = null;
	});

	it('Populates user', async () => {
		const mockUser: User = {
			id: 1,
			username: 'user',
			role: 'User',
			quota: 2,
			quotaUsed: 0,
			rootFolderId: 15
		};

		vi.mocked(endpoints.login).mockResolvedValue(undefined);
		vi.mocked(endpoints.getMe).mockResolvedValue(mockUser);

		await auth.login('user', 'epicPassword');

		expect(auth.user).toEqual(mockUser);

		expect(endpoints.login).toHaveBeenCalledWith('user', 'epicPassword');
		expect(endpoints.getMe).toHaveBeenCalledTimes(1);
	});

	it('Leaves user as null when failed', async () => {
		vi.mocked(endpoints.login).mockResolvedValue(undefined);
		vi.mocked(endpoints.getMe).mockRejectedValue(new Error('Unauthorized'));

		await auth.login('user', 'epicPassword');

		expect(auth.user).toBeNull();

		expect(endpoints.login).toHaveBeenCalledWith('user', 'epicPassword');
		expect(endpoints.getMe).toHaveBeenCalledTimes(1);
	});
});

describe('session', () => {
	beforeEach(() => {
		vi.clearAllMocks();
		auth.user = null;
	});

	it('IsLoggedIn returns true when user is logged', async () => {
		const mockUser: User = {
			id: 1,
			username: 'user',
			role: 'User',
			quota: 2,
			quotaUsed: 0,
			rootFolderId: 15
		};

		auth.user = mockUser;

		expect(auth.isLoggedIn).toEqual(true);
	});

	it('IsLoggedIn returns false when user is empty', async () => {
		expect(auth.isLoggedIn).toEqual(false);
	});

	it('IsAdmin return false when user is null', async () => {
		expect(auth.isAdmin).toEqual(false);
	});

	it('IsAdmin return false when user is User', async () => {
		const mockUser: User = {
			id: 1,
			username: 'user',
			role: 'User',
			quota: 2,
			quotaUsed: 0,
			rootFolderId: 15
		};

		auth.user = mockUser;

		expect(auth.isAdmin).toEqual(false);
	});

	it('IsAdmin return true when user is Admin', async () => {
		const mockUser: User = {
			id: 1,
			username: 'user',
			role: 'Admin',
			quota: 2,
			quotaUsed: 0,
			rootFolderId: 15
		};

		auth.user = mockUser;

		expect(auth.isAdmin).toEqual(true);
	});
});

describe('logout', () => {
	beforeEach(() => {
		vi.clearAllMocks();
		auth.user = null;
	});

	it('sets user to null regardless of api result', async () => {
		auth.user = { id: 1, username: 'user', role: 'User', quota: 2, quotaUsed: 0, rootFolderId: 15 };

		vi.mocked(endpoints.logout).mockResolvedValueOnce(undefined);
		await auth.logout();
		expect(auth.user).toBeNull();
	});

	it('set user to null on error', async () => {
		auth.user = { id: 1, username: 'user', role: 'User', quota: 2, quotaUsed: 0, rootFolderId: 15 };

		vi.mocked(endpoints.logout).mockRejectedValueOnce(new Error('Network error'));
		await auth.logout();
		expect(auth.user).toBeNull();
	});
});
