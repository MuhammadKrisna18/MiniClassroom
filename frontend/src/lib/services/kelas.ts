import { fetchApi } from './api';
import type { ApiResponse, Kelas, CreateKelasPayload, PengajuanKelas } from '$lib/types';

export const kelasService = {
	async getList(): Promise<ApiResponse<Kelas[]>> {
		return await fetchApi<ApiResponse<Kelas[]>>('/kelas', {
			method: 'GET'
		});
	},

	async getById(id: string): Promise<ApiResponse<Kelas>> {
		return await fetchApi<ApiResponse<Kelas>>(`/kelas/${id}`, {
			method: 'GET'
		});
	},

	async create(payload: CreateKelasPayload): Promise<ApiResponse<Kelas>> {
		return await fetchApi<ApiResponse<Kelas>>('/kelas', {
			method: 'POST',
			body: JSON.stringify(payload)
		});
	},

	async delete(id: string): Promise<ApiResponse<null>> {
		return await fetchApi<ApiResponse<null>>(`/kelas/${id}`, {
			method: 'DELETE'
		});
	},

	async requestKelas(kelasId: string, mataKuliahId: string): Promise<ApiResponse<PengajuanKelas>> {
		return await fetchApi<ApiResponse<PengajuanKelas>>('/kelas/request', {
			method: 'POST',
			body: JSON.stringify({ kelas_id: kelasId, mata_kuliah_id: mataKuliahId })
		});
	},

	async getMyPengajuan(): Promise<ApiResponse<PengajuanKelas[]>> {
		return await fetchApi<ApiResponse<PengajuanKelas[]>>('/kelas/pengajuan/me', {
			method: 'GET'
		});
	},

	async getAllPengajuan(): Promise<ApiResponse<PengajuanKelas[]>> {
		return await fetchApi<ApiResponse<PengajuanKelas[]>>('/kelas/pengajuan', {
			method: 'GET'
		});
	},

	async approvePengajuan(id: string): Promise<ApiResponse<null>> {
		return await fetchApi<ApiResponse<null>>(`/kelas/pengajuan/${id}/approve`, {
			method: 'POST'
		});
	},

	async rejectPengajuan(id: string): Promise<ApiResponse<null>> {
		return await fetchApi<ApiResponse<null>>(`/kelas/pengajuan/${id}/reject`, {
			method: 'POST'
		});
	},

	async getMahasiswaInKelas(id: string): Promise<ApiResponse<any[]>> {
		return await fetchApi<ApiResponse<any[]>>(`/kelas/pengajuan/${id}/mahasiswa`, {
			method: 'GET'
		});
	},

	async getMyJadwal(): Promise<ApiResponse<PengajuanKelas[]>> {
		return await fetchApi<ApiResponse<PengajuanKelas[]>>('/kelas/mahasiswa/my-jadwal', {
			method: 'GET'
		});
	},

	async getAvailableKelas(): Promise<ApiResponse<PengajuanKelas[]>> {
		return await fetchApi<ApiResponse<PengajuanKelas[]>>('/kelas/mahasiswa/krs/available', {
			method: 'GET'
		});
	},

	async ambilKelas(pengajuanId: string): Promise<ApiResponse<null>> {
		return await fetchApi<ApiResponse<null>>('/kelas/mahasiswa/krs', {
			method: 'POST',
			body: JSON.stringify({ pengajuan_id: pengajuanId })
		});
	},

	async getPertemuan(pengajuanId: string): Promise<ApiResponse<any[]>> {
		return await fetchApi<ApiResponse<any[]>>(`/pertemuan/pengajuan/${pengajuanId}`, {
			method: 'GET'
		});
	},

	async mulaiPertemuan(pengajuanId: string, judul: string): Promise<ApiResponse<any>> {
		return await fetchApi<ApiResponse<any>>('/pertemuan', {
			method: 'POST',
			body: JSON.stringify({ pengajuan_id: pengajuanId, judul })
		});
	},

	async akhiriPertemuan(pertemuanId: string): Promise<ApiResponse<null>> {
		return await fetchApi<ApiResponse<null>>(`/pertemuan/${pertemuanId}/end`, {
			method: 'PUT'
		});
	},

	async getAbsensi(pertemuanId: string): Promise<ApiResponse<any[]>> {
		return await fetchApi<ApiResponse<any[]>>(`/pertemuan/${pertemuanId}/absensi`, {
			method: 'GET'
		});
	},

	async submitAbsensi(pertemuanId: string, data: any[]): Promise<ApiResponse<null>> {
		return await fetchApi<ApiResponse<null>>(`/pertemuan/${pertemuanId}/absensi`, {
			method: 'PUT',
			body: JSON.stringify({ data })
		});
	},

	async getRekapKehadiran(pengajuanId: string): Promise<ApiResponse<any>> {
		let endpoint = `/pertemuan/rekap/${pengajuanId}`;
		
		
		
		
		const role = localStorage.getItem('role');
		if (role === 'admin') {
			endpoint = `/pertemuan/admin/rekap/${pengajuanId}`;
		}
		
		return await fetchApi<ApiResponse<any>>(endpoint, {
			method: 'GET'
		});
	}
};
