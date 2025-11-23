/**
 * Shared TypeScript types for Ham-Radio Cloud
 * These types mirror the backend Go models and API responses
 */

// User & Authentication

export type UserTier = 'free' | 'operator' | 'contester' | 'partner';

export interface User {
	id: number;
	callsign: string;
	email: string;
	name?: string;
	grid_square?: string;
	qrz_verified: boolean;
	tier: UserTier;
	qso_limit: number;
	qso_count: number;
	created_at: string;
	updated_at: string;
	last_login_at?: string;
}

export interface AuthResponse {
	user: User;
	token: string;
	refresh_token: string;
}

// QSO (Contact Log)

export interface QSO {
	id: number;
	user_id: number;
	callsign: string;
	operator_call?: string;
	station_callsign?: string;
	qso_date: string;
	time_on: string;
	time_off?: string;
	band?: string;
	band_rx?: string;
	freq?: number;
	freq_rx?: number;
	mode?: string;
	submode?: string;
	rst_sent?: string;
	rst_rcvd?: string;
	name?: string;
	qth?: string;
	gridsquare?: string;
	country?: string;
	dxcc?: number;
	state?: string;
	county?: string;
	comment?: string;
	notes?: string;
	tx_pwr?: number;
	rx_pwr?: number;
	prop_mode?: string;
	sat_name?: string;
	sat_mode?: string;
	contest_id?: string;
	stx?: number;
	srx?: number;
	lotw_qsl_sent?: 'Y' | 'N' | 'R';
	lotw_qsl_rcvd?: 'Y' | 'N';
	lotw_qslrdate?: string;
	eqsl_qsl_sent?: 'Y' | 'N';
	eqsl_qsl_rcvd?: 'Y' | 'N';
	created_at: string;
	updated_at: string;
}

export interface QSOFilter {
	callsign?: string;
	band?: string;
	mode?: string;
	country?: string;
	date_from?: string;
	date_to?: string;
	limit?: number;
	offset?: number;
}

// Propagation

export interface PropagationData {
	timestamp: string;
	solar_flux: number;
	sunspot_number: number;
	a_index: number;
	k_index: number;
	x_ray_flux?: string;
	conditions: BandConditions;
}

export interface BandConditions {
	[band: string]: 'poor' | 'fair' | 'good' | 'excellent';
}

// SDR

export interface SDRStream {
	id: string;
	name: string;
	location: string;
	url: string;
	frequency_range: [number, number];
	type: 'kiwi' | 'websdr' | 'openwebrx';
	status: 'online' | 'offline';
	listeners?: number;
}

// API Response Wrappers

export interface ApiResponse<T> {
	data: T;
	meta?: {
		page?: number;
		per_page?: number;
		total?: number;
	};
}

export interface ApiError {
	error: {
		code: string;
		message: string;
		details?: Record<string, unknown>;
	};
}

// Ham Radio Specific Constants

export const BANDS = [
	'160m',
	'80m',
	'60m',
	'40m',
	'30m',
	'20m',
	'17m',
	'15m',
	'12m',
	'10m',
	'6m',
	'2m',
	'70cm'
] as const;

export type Band = (typeof BANDS)[number];

export const MODES = [
	'SSB',
	'CW',
	'FM',
	'AM',
	'FT8',
	'FT4',
	'PSK31',
	'RTTY',
	'SSTV',
	'JT65',
	'JT9'
] as const;

export type Mode = (typeof MODES)[number];

export const SUBMODES = ['USB', 'LSB', 'BPSK', 'QPSK'] as const;

export type Submode = (typeof SUBMODES)[number];
