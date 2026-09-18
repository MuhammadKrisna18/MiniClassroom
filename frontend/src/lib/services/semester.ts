import { fetchApi } from './api';
import type { ApiResponse, Semester, SemesterMataKuliah, Pertemuan } from '../types';

export const semesterService = {
	getAll: async () => {
		return fetchApi<ApiResponse<Semester[]>>('/semester');
	},
	getById: async (id: string) => {
		return fetchApi<ApiResponse<Semester>>(`/semester/${id}`);
	},
	create: async (data: { nomor: number; min_sks: number; max_sks: number }) => {
		return fetchApi<ApiResponse<Semester>>('/semester', {
			method: 'POST',
			body: JSON.stringify(data)
		});
	},
	update: async (id: string, data: { min_sks?: number; max_sks?: number }) => {
		return fetchApi<ApiResponse<Semester>>(`/semester/${id}`, {
			method: 'PUT',
			body: JSON.stringify(data)
		});
	},
	delete: async (id: string) => {
		return fetchApi<ApiResponse<null>>(`/semester/${id}`, {
			method: 'DELETE'
		});
	},
	activate: async (id: string) => {
		return fetchApi<ApiResponse<null>>(`/semester/${id}/activate`, {
			method: 'PUT'
		});
	},
	assignMataKuliah: async (id: string, mata_kuliah_id: string, kategori: string = 'wajib') => {
		return fetchApi<ApiResponse<SemesterMataKuliah>>(`/semester/${id}/matakuliah`, {
			method: 'POST',
			body: JSON.stringify({ mata_kuliah_id, kategori })
		});
	},
	unassignMataKuliah: async (id: string, mkId: string) => {
		return fetchApi<ApiResponse<null>>(`/semester/${id}/matakuliah/${mkId}`, {
			method: 'DELETE'
		});
	},
	getPertemuan: async (kelasId: string, semesterId: string) => {
		return fetchApi<ApiResponse<Pertemuan[]>>(`/semester/pertemuan?kelas_id=${kelasId}&semester_id=${semesterId}`);
	},
	catatPertemuan: async (kelasId: string, semesterId: string, topik: string) => {
		return fetchApi<ApiResponse<Pertemuan>>('/semester/pertemuan', {
			method: 'POST',
			body: JSON.stringify({ kelas_id: kelasId, semester_id: semesterId, topik })
		});
	},
	setSKSProdi: async (id: string, configs: { program_studi_id: string; min_sks: number; max_sks: number }[]) => {
		return fetchApi<ApiResponse<any>>(`/semester/${id}/sks-prodi`, {
			method: 'PUT',
			body: JSON.stringify({ configs })
		});
	}
};
