<script lang="ts">
	import { onMount } from 'svelte';
	import { semesterService } from '$lib/services/semester';
	import type { Semester } from '$lib/types';
	import { toast } from '$lib/stores/toast.svelte';
	import { authState } from '$lib/stores/auth.svelte';
	import Card from '$lib/components/ui/Card.svelte';
	import Skeleton from '$lib/components/ui/Skeleton.svelte';
	import EmptyState from '$lib/components/ui/EmptyState.svelte';

	let semesters = $state<Semester[]>([]);
	let loading = $state(true);
	
	async function fetchData() {
		loading = true;
		try {
			const semRes = await semesterService.getAll();
			if (semRes.success && semRes.data) semesters = semRes.data;
		} catch (err: any) {
			toast.error(err.message || 'Gagal memuat data');
		} finally {
			loading = false;
		}
	}

	onMount(() => {
		
		if (!authState.token) {
			window.location.href = '/login';
			return;
		}
		fetchData();
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

<div class="page-container">
	<div class="page-header">
		<div class="header-content">
			<h1>Kurikulum (Fixed)</h1>
			<p class="text-muted">Daftar semester beserta mata kuliah yang sudah ditetapkan.</p>
		</div>
	</div>

	{#if loading}
		<div class="grid-container">
			<Skeleton />
			<Skeleton />
			<Skeleton />
		</div>
	{:else if semesters.length === 0}
		<EmptyState message="Belum ada data semester." />
	{:else}
		<div class="grid-container">
			{#each semesters as sem}
				<Card class={`semester-card ${sem.is_active ? 'active-semester' : ''}`}>
					<div class="semester-header">
						<div>
							<h2>Semester {sem.nomor}</h2>
							<div class="sks-badge">
								SKS: {sem.min_sks} - {sem.max_sks}
							</div>
						</div>
						{#if sem.is_active}
							<span class="status-badge active">Sedang Aktif</span>
						{/if}
					</div>
					
					<div class="semester-body">
						<div class="body-header">
							<h3>Mata Kuliah ({sem.mata_kuliah?.length || 0})</h3>
						</div>
						
						{#if !sem.mata_kuliah || sem.mata_kuliah.length === 0}
							<p class="text-muted text-sm">Belum ada mata kuliah.</p>
						{:else}
							{#each groupMataKuliahByProdi(sem) as group}
								<div class="prodi-group">
									<div class="prodi-title-container">
										<h4 class="prodi-title">{group.prodi}</h4>
										<div class="sks-info-badge {group.total_sks < group.min_sks || group.total_sks > group.max_sks ? 'sks-warning' : 'sks-ok'}">
											SKS Saat Ini: {group.total_sks} <span class="sks-separator">/</span> Batas: {group.min_sks} - {group.max_sks} SKS
										</div>
									</div>
									<div class="table-responsive">
										<table class="mk-table">
											<thead>
												<tr>
													<th width="5%">No</th>
													<th width="70%">Nama Mata Kuliah</th>
													<th width="10%">SKS</th>
													<th width="15%">Kategori</th>
												</tr>
											</thead>
											<tbody>
												{#each group.mks as smk, index}
													<tr>
														<td class="text-center text-muted">{index + 1}</td>
														<td class="font-medium">{smk.mata_kuliah?.name}</td>
														<td class="text-center">{smk.mata_kuliah?.sks}</td>
														<td class="text-center">
															{#if smk.kategori === 'pilihan'}
																<span class="badge-pilihan">Pilihan</span>
															{:else}
																<span class="badge-wajib">Wajib</span>
															{/if}
														</td>
													</tr>
												{/each}
											</tbody>
										</table>
									</div>
								</div>
							{/each}
						{/if}
					</div>
				</Card>
			{/each}
		</div>
	{/if}
</div>

<style>
	.page-container {
		display: flex;
		flex-direction: column;
		gap: 24px;
	}

	.page-header {
		display: flex;
		justify-content: space-between;
		align-items: flex-end;
		padding-bottom: 16px;
		border-bottom: 1px solid var(--border-color);
	}

	.header-content h1 {
		font-size: 1.8rem;
		font-weight: 700;
		color: var(--text-color);
		margin-bottom: 4px;
		letter-spacing: -0.02em;
	}

	.grid-container {
		display: grid;
		grid-template-columns: repeat(auto-fill, minmax(400px, 1fr));
		gap: 20px;
	}

	:global(.semester-card) {
		display: flex;
		flex-direction: column;
		padding: 0 !important;
		overflow: hidden;
	}

	:global(.active-semester) {
		border: 1px solid #10b981;
		box-shadow: 0 4px 12px rgba(16, 185, 129, 0.1);
	}

	.semester-header {
		padding: 20px;
		display: flex;
		justify-content: space-between;
		align-items: flex-start;
		border-bottom: 1px solid var(--surface-border);
		background: rgba(0, 0, 0, 0.01);
	}

	.semester-header h2 {
		font-size: 1.25rem;
		font-weight: 600;
		margin: 0 0 8px 0;
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

	.status-badge {
		padding: 4px 10px;
		border-radius: 20px;
		font-size: 0.75rem;
		font-weight: 600;
		background: var(--success-bg);
		color: var(--success-color);
	}

	.semester-body {
		padding: 20px;
		background: var(--surface-color);
	}

	.body-header {
		margin-bottom: 12px;
	}

	.body-header h3 {
		font-size: 0.95rem;
		font-weight: 600;
		margin: 0;
		color: var(--text-muted);
	}

	.prodi-group {
		margin-bottom: 24px;
	}

	.prodi-group:last-child {
		margin-bottom: 0;
	}

	.prodi-title-container {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin: 0 0 12px 0;
		padding-bottom: 8px;
		border-bottom: 1px solid var(--surface-border);
		flex-wrap: wrap;
		gap: 10px;
	}

	.prodi-title {
		font-size: 0.9rem;
		font-weight: 600;
		color: var(--text-main);
		margin: 0;
	}

	.sks-info-badge {
		font-size: 0.75rem;
		font-weight: 600;
		padding: 4px 10px;
		border-radius: 4px;
		display: flex;
		align-items: center;
		gap: 6px;
	}

	.sks-separator {
		color: rgba(0,0,0,0.2);
	}

	.sks-ok {
		background: rgba(16, 185, 129, 0.1);
		color: #10b981;
	}

	.sks-warning {
		background: rgba(245, 158, 11, 0.1);
		color: #f59e0b;
	}

	.table-responsive {
		width: 100%;
		overflow-x: auto;
		background: var(--secondary-color);
		border-radius: 8px;
		padding: 1px;
	}

	.mk-table {
		width: 100%;
		border-collapse: collapse;
		font-size: 0.85rem;
		background: var(--surface-color);
		border-radius: 7px;
		overflow: hidden;
	}

	.mk-table th,
	.mk-table td {
		padding: 12px 16px;
		text-align: left;
		border-bottom: 1px solid var(--surface-border);
	}

	.mk-table th {
		background: rgba(0, 0, 0, 0.02);
		font-weight: 600;
		color: var(--text-muted);
		text-transform: uppercase;
		font-size: 0.75rem;
		letter-spacing: 0.05em;
	}

	.mk-table tbody tr:last-child td {
		border-bottom: none;
	}

	.mk-table tbody tr:hover {
		background: rgba(0, 0, 0, 0.01);
	}

	.font-medium {
		font-weight: 500;
		color: var(--text-main);
	}
	
	.text-center {
		text-align: center !important;
	}

	.badge-wajib, .badge-pilihan {
		display: inline-block;
		font-size: 0.7rem;
		padding: 2px 6px;
		border-radius: 12px;
		font-weight: 600;
	}

	.badge-wajib {
		background: rgba(59, 130, 246, 0.1);
		color: #3b82f6;
	}

	.badge-pilihan {
		background: rgba(245, 158, 11, 0.1);
		color: #f59e0b;
	}
</style>
