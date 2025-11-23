import { writable } from 'svelte/store';
import type { User } from '$types';

interface AuthState {
	user: User | null;
	token: string | null;
	isAuthenticated: boolean;
}

function createAuthStore() {
	const { subscribe, set, update } = writable<AuthState>({
		user: null,
		token: null,
		isAuthenticated: false
	});

	return {
		subscribe,
		login: (user: User, token: string) => {
			set({ user, token, isAuthenticated: true });
			if (typeof localStorage !== 'undefined') {
				localStorage.setItem('auth_token', token);
				localStorage.setItem('user', JSON.stringify(user));
			}
		},
		logout: () => {
			set({ user: null, token: null, isAuthenticated: false });
			if (typeof localStorage !== 'undefined') {
				localStorage.removeItem('auth_token');
				localStorage.removeItem('user');
			}
		},
		init: () => {
			if (typeof localStorage !== 'undefined') {
				const token = localStorage.getItem('auth_token');
				const userStr = localStorage.getItem('user');
				if (token && userStr) {
					try {
						const user = JSON.parse(userStr);
						set({ user, token, isAuthenticated: true });
					} catch (e) {
						console.error('Failed to parse stored user:', e);
					}
				}
			}
		}
	};
}

export const auth = createAuthStore();
