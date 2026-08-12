import { describe, it, expect, vi, beforeEach } from 'vitest';
import { mockUser, mockAdminUser } from '$lib/mocks/testData';

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

	it('populates user on success', async () => {
		const user = mockUser();
		vi.mocked(endpoints.login).mockResolvedValue(undefined);
		vi.mocked(endpoints.getMe).mockResolvedValue(user);

		await auth.login('user', 'epicPassword');

		expect(auth.user).toEqual(user);
		expect(endpoints.login).toHaveBeenCalledWith('user', 'epicPassword');
		expect(endpoints.getMe).toHaveBeenCalledTimes(1);
	});

	it('leaves user as null when checkSession fails', async () => {
		vi.mocked(endpoints.login).mockResolvedValue(undefined);
		vi.mocked(endpoints.getMe).mockRejectedValue(new Error('Unauthorized'));

		await auth.login('user', 'epicPassword');

		expect(auth.user).toBeNull();
		expect(endpoints.login).toHaveBeenCalledWith('user', 'epicPassword');
		expect(endpoints.getMe).toHaveBeenCalledTimes(1);
	});
});

describe('session getters', () => {
	beforeEach(() => {
		vi.clearAllMocks();
		auth.user = null;
	});

	it('isLoggedIn is true when a user is set', () => {
		auth.user = mockUser();
		expect(auth.isLoggedIn).toBe(true);
	});

	it('isLoggedIn is false when user is null', () => {
		expect(auth.isLoggedIn).toBe(false);
	});

	it('isAdmin is false when user is null', () => {
		expect(auth.isAdmin).toBe(false);
	});

	it('isAdmin is false for a User role', () => {
		auth.user = mockUser();
		expect(auth.isAdmin).toBe(false);
	});

	it('isAdmin is true for an Admin role', () => {
		auth.user = mockAdminUser();
		expect(auth.isAdmin).toBe(true);
	});
});

describe('logout', () => {
	beforeEach(() => {
		vi.clearAllMocks();
		auth.user = null;
	});

	it('sets user to null on successful logout', async () => {
		auth.user = mockUser();
		vi.mocked(endpoints.logout).mockResolvedValueOnce(undefined);

		await auth.logout();

		expect(auth.user).toBeNull();
	});

	it('sets user to null even when the API call fails', async () => {
		auth.user = mockUser();
		vi.mocked(endpoints.logout).mockRejectedValueOnce(new Error('Network error'));

		await auth.logout();

		expect(auth.user).toBeNull();
	});
});
