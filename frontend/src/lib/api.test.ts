import { describe, it, expect, vi, beforeEach } from 'vitest';

vi.mock('$env/static/public', () => ({
	PUBLIC_API_BASE: 'http://test-api'
}));

const gotoMock = vi.fn();
vi.mock('$app/navigation', () => ({
	goto: gotoMock
}));

vi.mock('$app/paths', () => ({
	resolve: (path: string) => path
}));

const { endpoints } = await import('./api');

function mockFetchOnce(response: { ok?: boolean; status?: number; body?: string }) {
    // Create Fake Response
	const fetchMock = vi.fn().mockResolvedValue({
		ok: response.ok ?? true,
		status: response.status ?? 200,
		text: async () => response.body ?? '',
		blob: async () => new Blob([response.body ?? ''])
	});
    // Replaces Global fetch function
	vi.stubGlobal('fetch', fetchMock);
	return fetchMock;
}

beforeEach(() => {
	vi.restoreAllMocks();
	gotoMock.mockClear();
});

describe('request()', () => {
	it('sends the correct method, headers, and body for a POST', async () => {
		const fetchMock = mockFetchOnce({ body: '' });

		await endpoints.login('gabi', 'hunter2');

		expect(fetchMock).toHaveBeenCalledWith(
			'http://test-api/login',
			expect.objectContaining({
				method: 'POST',
				credentials: 'include',
				body: JSON.stringify({ username: 'gabi', password: 'hunter2' }),
				headers: expect.objectContaining({ 'Content-Type': 'application/json' })
			})
		);
	});

	it('parses a JSON response body on success', async () => {
		mockFetchOnce({
			body: JSON.stringify({ id: 1, username: 'gabi', role: 'User' })
		});

		const user = await endpoints.getMe();

		expect(user).toEqual({ id: 1, username: 'gabi', role: 'User' });
	});

	it('returns undefined for an empty successful response body', async () => {
		mockFetchOnce({ body: '' });

		const result = await endpoints.logout();

		expect(result).toBeUndefined();
	});

	it('throws the error message from a JSON error body', async () => {
		mockFetchOnce({
			ok: false,
			status: 400,
			body: JSON.stringify({ error: 'Username already taken' })
		});

		await expect(endpoints.register('gabi', 'hunter2')).rejects.toThrow(
			'Username already taken'
		);
	});

	it('falls back to a generic message when the error body is not JSON', async () => {
		mockFetchOnce({
			ok: false,
			status: 500,
			body: 'not json'
		});

		await expect(endpoints.getMe()).rejects.toThrow('Request failed with status 500');
	});

	it('redirects to login and throws on a 401', async () => {
		mockFetchOnce({ ok: false, status: 401, body: '' });

		await expect(endpoints.getMe()).rejects.toThrow('Session expired. Please log in again.');
		expect(gotoMock).toHaveBeenCalledWith('/login');
	});
});

describe('getFileContent()', () => {
	it('requests without the download flag by default', async () => {
		const fetchMock = mockFetchOnce({ body: 'file bytes' });

		await endpoints.getFileContent(42);

		expect(fetchMock).toHaveBeenCalledWith(
			'http://test-api/files/42/content',
			expect.objectContaining({ credentials: 'include' })
		);
	});

	it('appends ?download=true when download is requested', async () => {
		const fetchMock = mockFetchOnce({ body: 'file bytes' });

		await endpoints.getFileContent(42, true);

		expect(fetchMock).toHaveBeenCalledWith(
			'http://test-api/files/42/content?download=true',
			expect.objectContaining({ credentials: 'include' })
		);
	});

	it('throws on a failed response instead of returning an empty blob', async () => {
		mockFetchOnce({ ok: false, status: 404, body: '' });

		await expect(endpoints.getFileContent(42)).rejects.toThrow('Failed to fetch file content');
	});
});