import { endpoints } from '$lib/api';
import type { User } from '$lib/types';

class AuthClass {
	user = $state<User | null>(null);

	get isLoggedIn() {
		return this.user !== null;
	}

	get isAdmin() {
		return this.user?.role === 'Admin';
	}

	async register(username: string, password: string) {
		await endpoints.register(username, password);
	}

	async login(username: string, password: string) {
		await endpoints.login(username, password);
		await this.checkSession();
	}

	async logout() {
		try {
			await endpoints.logout();
		} catch {
			// No need to treat error if something went wrong revoke session in the frontend
		} finally {
			this.user = null;
		}
	}

	async checkSession() {
		try {
			this.user = await endpoints.getMe();
		} catch {
			this.user = null;
		}
	}
}

export const auth = new AuthClass();
