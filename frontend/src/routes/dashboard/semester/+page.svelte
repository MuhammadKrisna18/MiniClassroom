<script lang="ts">
	import { onMount } from 'svelte';
	import { fade } from 'svelte/transition';
	import { semesterService } from '$lib/services/semester';
	import { matakuliahService } from '$lib/services/matakuliah';
	import { programStudiService } from '$lib/services/programstudi';
	import type { Semester, MataKuliah, ProgramStudi } from '$lib/types';
	import { toast } from '$lib/stores/toast.svelte';
	import { authState } from '$lib/stores/auth.svelte';
	import Card from '$lib/components/ui/Card.svelte';
	import Skeleton from '$lib/components/ui/Skeleton.svelte';
	import EmptyState from '$lib/components/ui/EmptyState.svelte';

	let semesters = $state<Semester[]>([]);
	let allMataKuliah = $state<MataKuliah[]>([]);
	let prodiList = $state<ProgramStudi[]>([]);
	let loading = $state(true);
	let submitting = $state(false);

	let showForm = $state(false);
	let newSemester = $state({ nomor: 1, min_sks: 18, max_sks: 24 });
	
	let assignModalOpen = $state(false);
	let assignSemesterId = $state('');
	let selectedProdiId = $state('');
	let selectedMkId = $state('');
	let selectedKategori = $state('wajib');

	let editModalOpen = $state(false);
	let editSemesterId = $state('');
	let editSemester = $state({ min_sks: 18, max_sks: 24 });
	
	async function fetchData() {
		loading = true;
		try {
			const [semRes, mkRes, prodiRes] = await Promise.all([
				semesterService.getAll(),
				matakuliahService.getList(),
				programStudiService.getList()
			]);
			if (semRes.success && semRes.data) semesters = semRes.data;
			if (mkRes.success && mkRes.data) allMataKuliah = mkRes.data;
			if (prodiRes.success && prodiRes.data) prodiList = prodiRes.data;
		} catch (err: any) {
			toast.error(err.message || 'Gagal memuat data');
		} finally {
			loading = false;
		}
	}

	onMount(() => {
		if (authState.role !== 'admin') {
			window.location.href = '/dashboard';
			return;
		}
		fetchData();
	});

	async function handleCreate() {
		if (newSemester.nomor < 1 || newSemester.nomor > 8) {
			toast.error('Nomor semester harus antara 1-8');
			return;
		}
		if (newSemester.min_sks >= newSemester.max_sks) {
			toast.error('Minimal SKS harus lebih kecil dari Maksimal SKS');
			return;
		}

		submitting = true;
		try {
			const res = await semesterService.create(newSemester);
			if (res.success) {
				toast.success('Semester berhasil dibuat');
				showForm = false;
				newSemester = { nomor: semesters.length + 1 > 8 ? 8 : semesters.length + 1, min_sks: 18, max_sks: 24 };
				fetchData();
			}
		} catch (err: any) {
			toast.error(err.message || 'Gagal membuat semester');
		} finally {
			submitting = false;
		}
	}



	function openAssignModal(semesterId: string) {
		assignSemesterId = semesterId;
		selectedProdiId = '';
		selectedMkId = '';
		selectedKategori = 'wajib';
		assignModalOpen = true;
	}

	async function handleAssign() {
		if (!selectedMkId) return;
		submitting = true;
		try {
			const res = await semesterService.assignMataKuliah(assignSemesterId, selectedMkId, selectedKategori);
			if (res.success) {
				toast.success('Mata kuliah berhasil ditambahkan ke semester');
				assignModalOpen = false;
				fetchData();
			}
		} catch (err: any) {
			toast.error(err.message || 'Gagal menambahkan mata kuliah');
		} finally {
			submitting = false;
		}
	}

	async function handleUnassign(semesterId: string, mkId: string) {
		if (!confirm('Hapus mata kuliah dari semester ini?')) return;
		try {
			const res = await semesterService.unassignMataKuliah(semesterId, mkId);
			if (res.success) {
				toast.success('Mata kuliah berhasil dihapus');
				fetchData();
			}
		} catch (err: any) {
			toast.error(err.message || 'Gagal menghapus mata kuliah');
		}
	}

	function openEditModal(sem: Semester) {
		editSemesterId = sem.id;
		editSemester = { min_sks: sem.min_sks, max_sks: sem.max_sks };
		editModalOpen = true;
	}

	async function handleEdit() {
		if (editSemester.min_sks >= editSemester.max_sks) {
			toast.error('Minimal SKS harus lebih kecil dari Maksimal SKS');
			return;
		}
		submitting = true;
		try {
			const res = await semesterService.update(editSemesterId, editSemester);
			if (res.success) {
				toast.success('Semester berhasil diperbarui');
				editModalOpen = false;
				fetchData();
			}
		} catch (err: any) {
			toast.error(err.message || 'Gagal memperbarui semester');
		} finally {
			submitting = false;
		}
	}
	
	
	let availableMataKuliah = $derived(() => {
		const sem = semesters.find(s => s.id === assignSemesterId);
		const assignedIds = new Set(sem?.mata_kuliah?.map(sm => sm.mata_kuliah_id) || []);
		
		let available = allMataKuliah.filter(mk => !assignedIds.has(mk.id));
		
		if (selectedProdiId) {
			available = available.filter(mk => mk.program_studi_id === selectedProdiId);
		}
		
		return available;
	});
</script>

<div class="page-container animate-fade-in">
	<div class="page-header">
		<div>
			<h1 class="page-title">Manajemen Semester</h1>
			<p class="page-subtitle">Atur semester, batas SKS, dan mata kuliah tiap semester.</p>
		</div>
		<button class="btn btn-primary" onclick={() => showForm = !showForm}>
			{showForm ? 'Batal' : '+ Tambah Semester'}
		</button>
	</div>

	{#if showForm}
		<Card class="form-card animate-slide-up" style="margin-bottom: 24px; padding: 24px;">
			<h2 style="margin-bottom: 16px;">Semester Baru</h2>
			<form onsubmit={(e) => { e.preventDefault(); handleCreate(); }} class="semester-form">
				<div class="form-group">
					<label for="nomor">Nomor Semester (1-8)</label>
					<input type="number" id="nomor" bind:value={newSemester.nomor} min="1" max="8" required class="input" />
				</div>
				<div class="form-group">
					<label for="min_sks">Minimal SKS</label>
					<input type="number" id="min_sks" bind:value={newSemester.min_sks} min="1" required class="input" />
				</div>
				<div class="form-group">
					<label for="max_sks">Maksimal SKS</label>
					<input type="number" id="max_sks" bind:value={newSemester.max_sks} min="1" required class="input" />
				</div>
				<div class="form-actions">
					<button type="submit" class="btn btn-primary" disabled={submitting}>
						{submitting ? 'Menyimpan...' : 'Simpan Semester'}
					</button>
				</div>
			</form>
		</Card>
	{/if}

	{#if loading}
		<div class="grid-container">
			<Skeleton />
			<Skeleton />
			<Skeleton />
		</div>
	{:else if semesters.length === 0}
		<EmptyState message="Belum ada semester yang ditambahkan. Silakan tambah semester baru." />
	{:else}
		<div class="grid-container">
			{#each semesters as sem}
				<Card class="semester-card">
					<div class="semester-header">
						<div>
							<div style="display: flex; align-items: center; gap: 8px; margin-bottom: 8px;">
								<h2 style="margin: 0;">Semester {sem.nomor}</h2>
								<button class="btn-icon-tiny" onclick={() => openEditModal(sem)} title="Edit SKS Semester" style="opacity: 1; font-weight: 500; font-size: 0.8rem; padding: 4px 8px; border: 1px solid var(--surface-border); border-radius: 4px;">
									Edit
								</button>
							</div>
							<div class="sks-badge">
								SKS: {sem.min_sks} - {sem.max_sks}
							</div>
					</div>
					
					<div class="semester-body">
						<div class="body-header">
							<h3>Mata Kuliah ({sem.mata_kuliah?.length || 0})</h3>
							<button class="btn-icon-small" onclick={() => openAssignModal(sem.id)} title="Tambah MK" style="opacity: 1; font-weight: 500; font-size: 0.8rem; padding: 4px 8px; border: 1px solid var(--surface-border); border-radius: 4px;">
								Tambah
							</button>
						</div>
						
						{#if !sem.mata_kuliah || sem.mata_kuliah.length === 0}
							<p class="text-muted text-sm">Belum ada mata kuliah.</p>
						{:else}
							<ul class="mk-list">
								{#each sem.mata_kuliah as smk}
									<li>
										<div class="mk-info">
											<span class="mk-name">{smk.mata_kuliah?.name}</span>
											<span class="mk-meta">{smk.mata_kuliah?.sks} SKS • {smk.mata_kuliah?.program_studi?.name}</span>
											{#if smk.kategori === 'pilihan'}
												<span class="badge-pilihan">Pilihan</span>
											{:else}
												<span class="badge-wajib">Wajib</span>
											{/if}
										</div>
										<button class="btn-icon-tiny btn-reject" onclick={() => handleUnassign(sem.id, smk.mata_kuliah_id)} title="Hapus MK" style="opacity: 1; font-weight: 500; font-size: 0.8rem; padding: 4px 8px; border-radius: 4px;">
											Hapus
										</button>
									</li>
								{/each}
							</ul>
						{/if}
					</div>
				</Card>
			{/each}
		</div>
	{/if}
</div>

{#if assignModalOpen}
	<div class="modal-backdrop" transition:fade={{ duration: 200 }}>
		<div class="modal-content animate-slide-up">
			<h2>Tambahkan Mata Kuliah</h2>
			<p class="text-muted" style="margin-bottom: 20px;">Pilih mata kuliah yang tersedia untuk semester ini.</p>
			
			<div class="form-group">
				<label for="prodiSelect">Pilih Program Studi</label>
				<select id="prodiSelect" bind:value={selectedProdiId} class="input">
					<option value="">-- Semua Program Studi --</option>
					{#each prodiList as prodi}
						<option value={prodi.id}>{prodi.name}</option>
					{/each}
				</select>
			</div>

			<div class="form-group">
				<label for="mkSelect">Pilih Mata Kuliah</label>
				<select id="mkSelect" bind:value={selectedMkId} class="input">
					<option value="">-- Pilih Mata Kuliah --</option>
					{#each availableMataKuliah() as mk}
						<option value={mk.id}>{mk.name} ({mk.sks} SKS) - {mk.program_studi?.name}</option>
					{/each}
				</select>
			</div>

			<div class="form-group">
				<label for="mkKategori">Kategori Mata Kuliah</label>
				<select id="mkKategori" bind:value={selectedKategori} class="input">
					<option value="wajib">Wajib / Prodi</option>
					<option value="pilihan">Pilihan</option>
				</select>
			</div>
			
			<div class="modal-actions">
				<button class="btn btn-secondary" onclick={() => assignModalOpen = false} disabled={submitting}>Batal</button>
				<button class="btn btn-primary" onclick={handleAssign} disabled={!selectedMkId || submitting}>
					{submitting ? 'Menambahkan...' : 'Tambahkan'}
				</button>
			</div>
		</div>
	</div>
{/if}

{#if editModalOpen}
	<div class="modal-backdrop" transition:fade={{ duration: 200 }}>
		<div class="modal-content animate-slide-up">
			<h2>Edit SKS Semester</h2>
			<p class="text-muted" style="margin-bottom: 20px;">Sesuaikan batas minimal dan maksimal SKS untuk semester ini.</p>
			
			<form onsubmit={(e) => { e.preventDefault(); handleEdit(); }}>
				<div class="form-group">
					<label for="edit_min_sks">Minimal SKS</label>
					<input type="number" id="edit_min_sks" bind:value={editSemester.min_sks} min="1" required class="input" />
				</div>
				<div class="form-group">
					<label for="edit_max_sks">Maksimal SKS</label>
					<input type="number" id="edit_max_sks" bind:value={editSemester.max_sks} min="1" required class="input" />
				</div>
				
				<div class="form-actions" style="margin-top: 24px;">
					<button type="button" class="btn btn-secondary" onclick={() => editModalOpen = false} disabled={submitting}>
						Batal
					</button>
					<button type="submit" class="btn btn-primary" disabled={submitting}>
						{submitting ? 'Menyimpan...' : 'Simpan Perubahan'}
					</button>
				</div>
			</form>
		</div>
	</div>
{/if}

<style>
	.page-container {
		display: flex;
		flex-direction: column;
		gap: 24px;
	}

	.page-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
	}

	.page-title {
		font-size: 1.5rem;
		font-weight: 700;
		color: var(--text-color);
		margin: 0 0 4px 0;
	}

	.page-subtitle {
		color: var(--text-muted);
		margin: 0;
		font-size: 0.95rem;
	}

	.semester-form {
		display: grid;
		grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
		gap: 16px;
		align-items: end;
	}

	.grid-container {
		display: grid;
		grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
		gap: 20px;
	}

	.semester-card {
		padding: 20px;
		display: flex;
		flex-direction: column;
		gap: 16px;
		border-top: 4px solid transparent;
		transition: all 0.3s ease;
	}

	.semester-header {
		display: flex;
		justify-content: space-between;
		align-items: flex-start;
		border-bottom: 1px solid var(--border-color);
		padding-bottom: 12px;
	}

	.semester-header h2 {
		font-size: 1.25rem;
		font-weight: 700;
		margin: 0 0 4px 0;
	}

	.sks-badge {
		display: inline-block;
		background: var(--secondary-color);
		padding: 2px 8px;
		border-radius: 4px;
		font-size: 0.8rem;
		font-weight: 600;
		color: var(--text-muted);
	}

	.body-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: 12px;
	}

	.body-header h3 {
		font-size: 0.95rem;
		font-weight: 600;
		margin: 0;
		color: var(--text-muted);
	}

	.btn-icon-small {
		background: transparent;
		border: none;
		cursor: pointer;
		font-size: 1rem;
		padding: 4px;
		border-radius: 4px;
		transition: background 0.2s;
	}

	.btn-icon-small:hover {
		background: var(--secondary-hover);
	}

	.btn-icon-tiny {
		background: transparent;
		border: none;
		cursor: pointer;
		font-size: 0.8rem;
		padding: 2px;
		border-radius: 4px;
		opacity: 0.5;
		transition: all 0.2s;
	}

	.btn-icon-tiny:hover {
		opacity: 1;
		background: rgba(239, 68, 68, 0.1);
	}

	.mk-list {
		list-style: none;
		padding: 0;
		margin: 0;
		display: flex;
		flex-direction: column;
		gap: 8px;
		max-height: 250px;
		overflow-y: auto;
	}

	.mk-list li {
		display: flex;
		justify-content: space-between;
		align-items: flex-start;
		padding: 8px;
		background: var(--secondary-color);
		border-radius: 6px;
	}

	.mk-info {
		display: flex;
		flex-direction: column;
	}

	.mk-name {
		font-weight: 500;
		font-size: 0.9rem;
	}

	.mk-meta {
		font-size: 0.75rem;
		color: var(--text-muted);
		margin-bottom: 4px;
	}

	.modal-backdrop {
		position: fixed;
		top: 0;
		left: 0;
		right: 0;
		bottom: 0;
		background: rgba(15, 23, 42, 0.5);
		backdrop-filter: blur(4px);
		display: flex;
		align-items: center;
		justify-content: center;
		z-index: 1000;
	}

	.modal-content {
		background: var(--surface-color);
		border-radius: var(--radius-lg);
		padding: 32px;
		width: 90%;
		max-width: 500px;
		box-shadow: 0 20px 40px rgba(0, 0, 0, 0.2);
	}

	.modal-content h2 {
		margin: 0 0 8px 0;
		font-size: 1.5rem;
	}

	.modal-actions {
		display: flex;
		justify-content: flex-end;
		gap: 12px;
		margin-top: 24px;
	}
</style>
