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
	let selectedMkIds = $state<string[]>([]);
	let selectedKategori = $state('wajib');

	let editModalOpen = $state(false);
	let editSemesterId = $state('');
	let editSKSProdi = $state<{program_studi_id: string, min_sks: number, max_sks: number}[]>([]);
	
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
		selectedMkIds = [];
		selectedKategori = 'wajib';
		assignModalOpen = true;
	}

	async function handleAssign() {
		if (selectedMkIds.length === 0) return;
		submitting = true;
		try {
			for (const mkId of selectedMkIds) {
				await semesterService.assignMataKuliah(assignSemesterId, mkId, selectedKategori);
			}
			toast.success('Mata kuliah berhasil ditambahkan ke semester');
			assignModalOpen = false;
			fetchData();
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
		editSKSProdi = prodiList.map(prodi => {
			const existing = sem.sks_prodi?.find(p => p.program_studi_id === prodi.id);
			return {
				program_studi_id: prodi.id,
				min_sks: existing ? existing.min_sks : sem.min_sks,
				max_sks: existing ? existing.max_sks : sem.max_sks
			};
		});
		editModalOpen = true;
	}

	async function handleEdit() {
		for (const cfg of editSKSProdi) {
			if (cfg.min_sks >= cfg.max_sks) {
				toast.error('Minimal SKS harus lebih kecil dari Maksimal SKS untuk semua prodi');
				return;
			}
		}
		
		submitting = true;
		try {
			const res = await semesterService.setSKSProdi(editSemesterId, editSKSProdi);
			if (res.success) {
				toast.success('Batas SKS Prodi berhasil diperbarui');
				editModalOpen = false;
				fetchData();
			}
		} catch (err: any) {
			toast.error(err.message || 'Gagal memperbarui batas SKS');
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

	function groupMataKuliahByProdi(semester: Semester) {
		const mata_kuliah = semester.mata_kuliah || [];
		
		const sharedMks: any[] = [];
		const coreMks: any[] = [];
		
		for (const smk of mata_kuliah) {
			const code = smk.mata_kuliah?.program_studi?.code || '';
			if (code === 'DEPT' || code === 'MKUB') {
				sharedMks.push(smk);
			} else {
				coreMks.push(smk);
			}
		}
		
		const coreProdisMap = new Map();
		
		if (semester.sks_prodi) {
			for (const sp of semester.sks_prodi) {
				const code = sp.program_studi?.code;
				if (code && code !== 'DEPT' && code !== 'MKUB') {
					coreProdisMap.set(sp.program_studi_id, {
						id: sp.program_studi_id,
						name: sp.program_studi?.name || 'Unknown',
						code: code
					});
				}
			}
		}
		
		for (const smk of coreMks) {
			const p = smk.mata_kuliah?.program_studi;
			if (p && !coreProdisMap.has(p.id) && p.code !== 'DEPT' && p.code !== 'MKUB') {
				coreProdisMap.set(p.id, {
					id: p.id,
					name: p.name,
					code: p.code
				});
			}
		}
		
		if (coreProdisMap.size === 0 && sharedMks.length > 0) {
			coreProdisMap.set('TI', { id: 'TI', name: 'Teknik Informatika', code: 'TI' });
			coreProdisMap.set('RPL', { id: 'RPL', name: 'Rekayasa Perangkat Lunak', code: 'RPL' });
			coreProdisMap.set('RKA', { id: 'RKA', name: 'Rekayasa Kecerdasan Artificial', code: 'RKA' });
		}
		
		const groups = [];
		for (const prodi of coreProdisMap.values()) {
			const specificMks = coreMks.filter(smk => smk.mata_kuliah?.program_studi_id === prodi.id || smk.mata_kuliah?.program_studi?.code === prodi.code);
			const allMksForProdi = [...specificMks, ...sharedMks];
			
			const total_sks = allMksForProdi.reduce((sum, smk) => sum + (smk.mata_kuliah?.sks || 0), 0);
			
			let min_sks = semester.min_sks;
			let max_sks = semester.max_sks;
			if (semester.sks_prodi) {
				const limits = semester.sks_prodi.find(sp => sp.program_studi_id === prodi.id || sp.program_studi?.code === prodi.code);
				if (limits) {
					min_sks = limits.min_sks;
					max_sks = limits.max_sks;
				}
			}
			
			if (allMksForProdi.length > 0 || (semester.sks_prodi && semester.sks_prodi.find(sp => sp.program_studi_id === prodi.id || sp.program_studi?.code === prodi.code))) {
			    groups.push({
				    prodi: prodi.name,
				    prodi_id: prodi.id,
				    total_sks,
				    mks: allMksForProdi,
				    min_sks,
				    max_sks
			    });
			}
		}
		
		return groups;
	}
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
						<div style="display: flex; justify-content: space-between; align-items: center; width: 100%;">
							<h2 style="margin: 0; display: flex; align-items: center; gap: 8px;">
								<span style="font-size: 1.4rem; font-weight: 800; color: var(--primary-color);">Semester {sem.nomor}</span>
							</h2>
							<button class="btn-icon-tiny" onclick={() => openEditModal(sem)} title="Edit SKS Semester">
								<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"></path><path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"></path></svg>
							</button>
						</div>
						<div class="sks-badge-container">
							{#if sem.sks_prodi && sem.sks_prodi.length > 0}
								{#each sem.sks_prodi as p}
									<div class="sks-badge" title="Min: {p.min_sks}, Max: {p.max_sks}">
										<span class="prodi-code">{p.program_studi?.code || p.program_studi?.name}</span>
										<span class="sks-range">{p.min_sks}-{p.max_sks}</span>
									</div>
								{/each}
							{:else}
								<div class="sks-badge">
									<span class="prodi-code">Default</span>
									<span class="sks-range">{sem.min_sks}-{sem.max_sks}</span>
								</div>
							{/if}
						</div>
					</div>
					
					<div class="semester-body">
						<div class="body-header">
							<h3 style="display: flex; align-items: center; gap: 8px;">
								<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M4 19.5A2.5 2.5 0 0 1 6.5 17H20"></path><path d="M6.5 2H20v20H6.5A2.5 2.5 0 0 1 4 19.5v-15A2.5 2.5 0 0 1 6.5 2z"></path></svg>
								Mata Kuliah <span class="badge-count">{sem.mata_kuliah?.length || 0}</span>
							</h3>
							<button class="btn-icon-small btn-tambah" onclick={() => openAssignModal(sem.id)} title="Tambah MK">
								<svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="12" y1="5" x2="12" y2="19"></line><line x1="5" y1="12" x2="19" y2="12"></line></svg>
								Tambah
							</button>
						</div>
						
						{#if !sem.mata_kuliah || sem.mata_kuliah.length === 0}
							<div class="empty-mk-state">
								<p>Belum ada mata kuliah.</p>
							</div>
						{:else}
							<div class="mk-list-container">
								{#each groupMataKuliahByProdi(sem) as group}
									<div class="prodi-section">
										<h4 class="prodi-section-title">
											{group.prodi}
											<span class="badge-count-small">{group.mks.length} MK</span>
										</h4>
										<div class="table-responsive">
											<table class="mk-table">
												<thead>
													<tr>
														<th>Mata Kuliah</th>
														<th style="width: 60px; text-align: center;">SKS</th>
														<th style="width: 80px; text-align: center;">Sifat</th>
														<th style="width: 40px; text-align: right;">Aksi</th>
													</tr>
												</thead>
												<tbody>
													{#each group.mks as smk}
														<tr>
															<td class="mk-name-cell" title={smk.mata_kuliah?.name}>{smk.mata_kuliah?.name}</td>
															<td style="text-align: center; font-weight: 600;">{smk.mata_kuliah?.sks}</td>
															<td style="text-align: center;">
																{#if smk.kategori === 'pilihan'}
																	<span class="badge-pilihan">Pilihan</span>
																{:else}
																	<span class="badge-wajib">Wajib</span>
																{/if}
															</td>
															<td style="text-align: right;">
																<button class="btn-icon-tiny btn-reject" onclick={() => handleUnassign(sem.id, smk.mata_kuliah_id)} title="Hapus MK">
																	<svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="3 6 5 6 21 6"></polyline><path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path></svg>
																</button>
															</td>
														</tr>
													{/each}
												</tbody>
											</table>
										</div>
									</div>
								{/each}
							</div>
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
				<div class="checkbox-list" style="max-height: 200px; overflow-y: auto; border: 1px solid var(--surface-border); border-radius: var(--radius-md); padding: 12px; background: var(--secondary-color);">
					{#if availableMataKuliah().length === 0}
						<p class="text-muted" style="margin: 0; font-size: 0.9rem;">Tidak ada mata kuliah tersedia</p>
					{:else}
						{#each availableMataKuliah() as mk}
							<label class="checkbox-item" style="display: flex; align-items: flex-start; gap: 12px; margin-bottom: 12px; cursor: pointer;">
								<input type="checkbox" value={mk.id} bind:group={selectedMkIds} style="margin-top: 4px; cursor: pointer;" />
								<div style="display: flex; flex-direction: column;">
									<span style="font-size: 0.95rem; font-weight: 500; color: var(--text-main);">{mk.name}</span>
									<span style="font-size: 0.8rem; color: var(--text-muted);">{mk.sks} SKS • {mk.program_studi?.name}</span>
								</div>
							</label>
						{/each}
					{/if}
				</div>
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
				<button class="btn btn-primary" onclick={handleAssign} disabled={selectedMkIds.length === 0 || submitting}>
					{submitting ? 'Menambahkan...' : 'Tambahkan'}
				</button>
			</div>
		</div>
	</div>
{/if}

{#if editModalOpen}
	<div class="modal-backdrop" transition:fade={{ duration: 200 }}>
		<div class="modal-content animate-slide-up">
			<h2>Edit SKS per Prodi</h2>
			<p class="text-muted" style="margin-bottom: 20px;">Sesuaikan batas minimal dan maksimal SKS untuk masing-masing program studi di semester ini.</p>
			
			<form onsubmit={(e) => { e.preventDefault(); handleEdit(); }}>
				<div style="max-height: 400px; overflow-y: auto; padding-right: 8px;">
					{#each editSKSProdi as cfg}
						<div style="margin-bottom: 16px; padding: 12px; background: var(--secondary-color); border-radius: 6px; border: 1px solid var(--surface-border);">
							<h3 style="margin: 0 0 12px 0; font-size: 1rem;">
								{prodiList.find(p => p.id === cfg.program_studi_id)?.name}
							</h3>
							<div style="display: grid; grid-template-columns: 1fr 1fr; gap: 12px;">
								<div class="form-group" style="margin: 0;">
									<label for="edit_min_{cfg.program_studi_id}" style="font-size: 0.85rem;">Minimal SKS</label>
									<input type="number" id="edit_min_{cfg.program_studi_id}" bind:value={cfg.min_sks} min="1" required class="input" />
								</div>
								<div class="form-group" style="margin: 0;">
									<label for="edit_max_{cfg.program_studi_id}" style="font-size: 0.85rem;">Maksimal SKS</label>
									<input type="number" id="edit_max_{cfg.program_studi_id}" bind:value={cfg.max_sks} min="1" required class="input" />
								</div>
							</div>
						</div>
					{/each}
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
		grid-template-columns: repeat(auto-fill, minmax(360px, 1fr));
		gap: 24px;
	}

	.semester-card {
		padding: 24px;
		display: flex;
		flex-direction: column;
		gap: 20px;
		border-top: 4px solid transparent;
		transition: transform 0.2s ease, box-shadow 0.2s ease;
		background: var(--surface-color);
		border: 1px solid var(--surface-border);
		box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.05), 0 2px 4px -1px rgba(0, 0, 0, 0.03);
	}

	.semester-card:hover {
		transform: translateY(-2px);
		box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.1), 0 4px 6px -2px rgba(0, 0, 0, 0.05);
	}

	.semester-header {
		display: flex;
		flex-direction: column;
		gap: 8px;
		border-bottom: 1px solid var(--border-color);
		padding-bottom: 16px;
	}

	.sks-badge-container {
		display: flex;
		flex-wrap: wrap;
		gap: 8px;
	}

	.sks-badge {
		display: inline-flex;
		align-items: center;
		background: var(--secondary-color);
		border-radius: 6px;
		border: 1px solid var(--surface-border);
		overflow: hidden;
		font-size: 0.75rem;
		font-weight: 600;
	}

	.prodi-code {
		background: var(--primary-color);
		color: white;
		padding: 2px 6px;
	}

	.sks-range {
		padding: 2px 6px;
		color: var(--text-color);
	}

	.body-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: 16px;
	}

	.body-header h3 {
		font-size: 1rem;
		font-weight: 600;
		margin: 0;
		color: var(--text-color);
	}

	.badge-count {
		background: var(--secondary-color);
		color: var(--text-color);
		padding: 2px 8px;
		border-radius: 9999px;
		font-size: 0.75rem;
		font-weight: 700;
	}

	.btn-icon-small {
		background: transparent;
		border: 1px solid var(--surface-border);
		cursor: pointer;
		font-size: 0.8rem;
		font-weight: 600;
		padding: 4px 10px;
		border-radius: 6px;
		transition: all 0.2s;
		display: inline-flex;
		align-items: center;
		gap: 6px;
		color: var(--text-color);
	}

	.btn-tambah {
		background: var(--primary-color);
		color: white;
		border: none;
	}

	.btn-tambah:hover {
		background: var(--primary-hover);
		transform: translateY(-1px);
	}

	.btn-icon-tiny {
		background: transparent;
		border: 1px solid var(--surface-border);
		cursor: pointer;
		font-size: 0.8rem;
		padding: 6px;
		border-radius: 6px;
		color: var(--text-muted);
		transition: all 0.2s;
		display: inline-flex;
		align-items: center;
		justify-content: center;
	}

	.btn-icon-tiny:hover {
		color: var(--primary-color);
		border-color: var(--primary-color);
		background: rgba(79, 70, 229, 0.05);
	}

	.btn-reject:hover {
		color: #ef4444;
		border-color: #ef4444;
		background: rgba(239, 68, 68, 0.05);
	}

	.empty-mk-state {
		text-align: center;
		padding: 24px 0;
		background: var(--secondary-color);
		border-radius: 8px;
		border: 1px dashed var(--surface-border);
	}

	.empty-mk-state p {
		margin: 0;
		color: var(--text-muted);
		font-size: 0.9rem;
	}

	.mk-list-container {
		max-height: 280px;
		overflow-y: auto;
		padding-right: 4px;
	}

	.mk-list-container::-webkit-scrollbar {
		width: 4px;
	}
	.mk-list-container::-webkit-scrollbar-track {
		background: var(--secondary-color);
		border-radius: 4px;
	}
	.mk-list-container::-webkit-scrollbar-thumb {
		background: var(--border-color);
		border-radius: 4px;
	}

	.mk-list-container {
		max-height: 350px;
		overflow-y: auto;
		padding-right: 4px;
		display: flex;
		flex-direction: column;
		gap: 20px;
	}

	.mk-list-container::-webkit-scrollbar {
		width: 4px;
	}
	.mk-list-container::-webkit-scrollbar-track {
		background: var(--secondary-color);
		border-radius: 4px;
	}
	.mk-list-container::-webkit-scrollbar-thumb {
		background: var(--border-color);
		border-radius: 4px;
	}

	.prodi-section-title {
		font-size: 0.9rem;
		font-weight: 700;
		color: var(--primary-color);
		margin: 0 0 10px 0;
		display: flex;
		align-items: center;
		gap: 8px;
		padding-bottom: 6px;
		border-bottom: 2px solid var(--secondary-color);
	}

	.badge-count-small {
		background: var(--secondary-color);
		color: var(--text-muted);
		padding: 2px 6px;
		border-radius: 4px;
		font-size: 0.7rem;
		font-weight: 600;
	}

	.table-responsive {
		overflow-x: auto;
	}

	.mk-table {
		width: 100%;
		border-collapse: separate;
		border-spacing: 0;
		font-size: 0.85rem;
	}

	.mk-table th, .mk-table td {
		padding: 10px 12px;
		border-bottom: 1px solid var(--surface-border);
	}

	.mk-table th {
		text-align: left;
		font-weight: 600;
		color: var(--text-muted);
		background: var(--secondary-color);
		text-transform: uppercase;
		font-size: 0.7rem;
		letter-spacing: 0.05em;
	}

	.mk-table th:first-child {
		border-top-left-radius: 6px;
		border-bottom-left-radius: 6px;
	}

	.mk-table th:last-child {
		border-top-right-radius: 6px;
		border-bottom-right-radius: 6px;
	}

	.mk-table tr:last-child td {
		border-bottom: none;
	}

	.mk-table tbody tr {
		transition: background-color 0.15s;
	}

	.mk-table tbody tr:hover {
		background: rgba(0, 0, 0, 0.02);
	}

	.mk-name-cell {
		font-weight: 500;
		color: var(--text-color);
		max-width: 180px;
		white-space: nowrap;
		overflow: hidden;
		text-overflow: ellipsis;
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
