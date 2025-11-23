const API_BASE_URL = import.meta.env.VITE_API_URL || '/api/v1';

interface ApiOptions extends RequestInit {
	token?: string;
}

async function request<T>(endpoint: string, options: ApiOptions = {}): Promise<T> {
	const { token, ...fetchOptions } = options;

	const headers: HeadersInit = {
		'Content-Type': 'application/json',
		...fetchOptions.headers
	};

	if (token) {
		headers['Authorization'] = `Bearer ${token}`;
	}

	const response = await fetch(`${API_BASE_URL}${endpoint}`, {
		...fetchOptions,
		headers
	});

	if (!response.ok) {
		const error = await response.json().catch(() => ({ message: response.statusText }));
		throw new Error(error.message || 'API request failed');
	}

	return response.json();
}

export const api = {
	// Auth
	login: (provider: string) => request(`/auth/login/${provider}`),
	logout: () => request('/auth/logout', { method: 'POST' }),

	// QSOs
	getQSOs: (filters?: object, token?: string) =>
		request(`/qso?${new URLSearchParams(filters as any)}`, { token }),
	createQSO: (qso: object, token?: string) =>
		request('/qso', { method: 'POST', body: JSON.stringify(qso), token }),
	updateQSO: (id: string, qso: object, token?: string) =>
		request(`/qso/${id}`, { method: 'PUT', body: JSON.stringify(qso), token }),
	deleteQSO: (id: string, token?: string) => request(`/qso/${id}`, { method: 'DELETE', token }),

	// Propagation
	getCurrentPropagation: () => request('/propagation/current'),
	getBandConditions: () => request('/propagation/bands'),

	// User
	getUserProfile: (token: string) => request('/user/profile', { token }),
	updateUserSettings: (settings: object, token: string) =>
		request('/user/settings', { method: 'PUT', body: JSON.stringify(settings), token }),

	// Health
	health: () => request('/health')
};

export default api;
