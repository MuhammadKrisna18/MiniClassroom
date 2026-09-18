export interface ApiResponse<T = any> {
	success: boolean;
	message: string;
	data?: T;
}

export interface User {
	id: string;
	name: string;
	email: string;
	role: 'admin' | 'dosen';
	photo_url?: string;
	created_at: string;
}

export interface UserProfile {
	id: string;
	name: string;
	nickname?: string;
	nid?: string;
	email: string;
	role: 'admin' | 'dosen';
	program_studi_id?: string;
	program_studi?: ProgramStudi;
	pending_email?: string;
	photo_url?: string;
	created_at: string;
}

export interface EmailChangeRequest {
	id: string;
	user_id: string;
	user?: UserProfile;
	new_email: string;
	status: 'pending' | 'approved' | 'rejected';
	created_at: string;
}

export interface Dosen {
	id: string;
	name: string;
	email: string;
	photo_url?: string;
	created_at: string;
}

export interface LoginResponse {
	token: string;
	role: 'admin' | 'dosen';
}

export interface ProgramStudi {
	id: string;
	name: string;
	code: string;
}

export interface MataKuliah {
	id: string;
	name: string;
	sks: number;
	program_studi_id: string;
	program_studi?: ProgramStudi;
	pengajuan?: PengajuanMataKuliah[];
	created_at: string;
	updated_at: string;
}

export interface PengajuanMataKuliah {
	id: string;
	dosen_id: string;
	dosen?: UserProfile;
	mata_kuliah_id: string;
	mata_kuliah?: MataKuliah;
	status: 'pending' | 'approved' | 'rejected' | 'offered';
	code: string;
	created_at: string;
}

export interface Kelas {
	id: string;
	name: string;
	capacity: number;
	hari: string;
	jam_mulai: string;
	jam_selesai: string;
	program_studi_id: string;
	program_studi?: ProgramStudi;
	pengajuan?: PengajuanKelas[];
	created_at: string;
}

export interface PengajuanKelas {
	id: string;
	dosen_id: string;
	dosen?: UserProfile;
	kelas_id: string;
	kelas?: Kelas;
	mata_kuliah_id: string;
	mata_kuliah?: MataKuliah;
	status: 'pending' | 'approved' | 'rejected';
	code: string;
	created_at: string;
}

export interface CreateKelasPayload {
	name: string;
	capacity: number;
	hari: string;
	jam_mulai: string;
	jam_selesai: string;
	program_studi_id: string;
}

export interface Semester {
	id: string;
	nomor: number;
	min_sks: number;
	max_sks: number;
	is_active: boolean;
	mata_kuliah?: SemesterMataKuliah[];
	sks_prodi?: SemesterSKSProdi[];
	created_at: string;
}

export interface SemesterSKSProdi {
	id: string;
	semester_id: string;
	program_studi_id: string;
	program_studi?: ProgramStudi;
	min_sks: number;
	max_sks: number;
}

export interface SemesterMataKuliah {
	id: string;
	semester_id: string;
	mata_kuliah_id: string;
	mata_kuliah?: MataKuliah;
	kategori: string;
	created_at: string;
}

export interface Pertemuan {
	id: string;
	semester_id: string;
	kelas_id: string;
	dosen_id: string;
	nomor_pertemuan: number;
	topik: string;
	status: 'belum' | 'selesai';
	tanggal_selesai?: string;
	created_at: string;
}
