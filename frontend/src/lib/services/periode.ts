import { fetchApi } from './api';
import type { ApiResponse } from '../types';

export interface PeriodeAkademik {
	id: string;
	tahun: string;
	jenis: 'ganjil' | 'genap' | 'pendek';
	is_active: boolean;
	created_at: string;
	updated_at: string;
}

export const periodeService = {
	getAll: async () => {
		return fetchApi<ApiResponse<PeriodeAkademik[]>>('/periode');
	},
	getActive: async () => {
		return fetchApi<ApiResponse<PeriodeAkademik>>('/periode/active');
	},
	create: async (data: { tahun: string; jenis: string }) => {
		return fetchApi<ApiResponse<PeriodeAkademik>>('/periode', {
			method: 'POST',
			body: JSON.stringify(data)
		});
	},
	delete: async (id: string) => {
		return fetchApi<ApiResponse<null>>(`/periode/${id}`, {
			method: 'DELETE'
		});
	},
	activate: async (id: string) => {
		return fetchApi<ApiResponse<null>>(`/periode/${id}/active`, {
			method: 'POST'
		});
	}
};
