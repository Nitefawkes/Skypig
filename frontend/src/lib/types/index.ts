export interface User {
	id: string;
	callsign: string;
	email: string;
	name: string;
	tier: 'free' | 'operator' | 'contester' | 'partner';
	createdAt: string;
}

export interface QSO {
	id: string;
	userId: string;
	callsign: string;
	frequency: number;
	band: string;
	mode: string;
	rst_sent: string;
	rst_received: string;
	qso_date: string;
	time_on: string;
	time_off?: string;
	grid_square?: string;
	country?: string;
	state?: string;
	county?: string;
	comment?: string;
	lotw_sent: boolean;
	lotw_confirmed: boolean;
	createdAt: string;
}

export interface PropagationData {
	timestamp: string;
	solar_flux: number;
	sunspot_number: number;
	a_index: number;
	k_index: number;
	xray_flux: string;
}

export interface BandCondition {
	band: string;
	condition: 'good' | 'fair' | 'poor';
	score: number;
	day_night: 'day' | 'night';
}
