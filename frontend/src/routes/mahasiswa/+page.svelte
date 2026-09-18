<script lang="ts">
	import { onMount } from 'svelte';
	import { fade } from 'svelte/transition';
	import { kelasService } from '$lib/services/kelas';
	import { toast } from '$lib/stores/toast.svelte';

	let bursaKelasList = $state<any[]>([]);
	let jadwalList = $state<any[]>([]);
	let loading = $state(true);
	let error = $state('');

	onMount(async () => {
		await loadData();
	});

	async function loadData() {
		loading = true;
		error = '';
		try {
			const [bursaRes, jadwalRes] = await Promise.all([
				kelasService.getAvailableKelas(),
				kelasService.getMyJadwal()
			]);
			
			if (bursaRes.success && bursaRes.data) {
				bursaKelasList = bursaRes.data;
			} else {
				error = bursaRes.message || 'Gagal memuat bursa kelas';
			}

			if (jadwalRes.success && jadwalRes.data) {
				jadwalList = jadwalRes.data;
			}
		} catch (e: any) {
			error = e.message || 'Terjadi kesalahan saat memuat data';
		} finally {
			loading = false;
		}
	}

	async function ambilKelas(pengajuanId: string) {
		try {
			const res = await kelasService.ambilKelas(pengajuanId);
			if (res.success) {
				toast.success('Berhasil mengambil kelas');
				await loadData(); // refresh data
			} else {
				toast.error(res.message || 'Gagal mengambil kelas');
			}
		} catch (e: any) {
			toast.error(e.message || 'Terjadi kesalahan sistem');
		}
	}
</script>

<svelte:head>
	<title>KRS Mahasiswa - SIAKAD Pro</title>
</svelte:head>

{#if loading}
	<div class="loading-state">
		<div class="spinner"></div>
		<p>Memuat data...</p>
	</div>
{:else if error}
	<div class="error-state">
		<svg xmlns="http://www.w3.org/2000/svg" width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
			<circle cx="12" cy="12" r="10"/><line x1="12" y1="8" x2="12" y2="12"/><line x1="12" y1="16" x2="12.01" y2="16"/>
		</svg>
		<p>{error}</p>
		<button class="btn btn-primary" onclick={() => loadData()}>Coba Lagi</button>
	</div>
{:else}
	<div class="page-header">
		<h1>Jadwal Kelas Anda</h1>
		<p>Kelas yang telah resmi Anda ambil pada semester ini.</p>
	</div>

	{#if jadwalList.length === 0}
		<div class="empty-state glass-panel">
			<div class="empty-icon">
				<svg xmlns="http://www.w3.org/2000/svg" width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
					<rect x="3" y="4" width="18" height="18" rx="2" ry="2"></rect><line x1="16" y1="2" x2="16" y2="6"></line><line x1="8" y1="2" x2="8" y2="6"></line><line x1="3" y1="10" x2="21" y2="10"></line>
				</svg>
			</div>
			<h2>Belum Ada Kelas</h2>
			<p>Silakan ambil kelas di Bursa KRS di bawah ini.</p>
		</div>
	{:else}
		<div class="grid-container" style="margin-bottom: 3rem;">
			{#each jadwalList as jdwl (jdwl.id)}
				<div class="mk-card glass-panel" in:fade>
					<div class="mk-header">
						<h3 class="mk-title">{jdwl.kelas?.name}</h3>
						<span class="sks-badge" style="background: var(--accent-color);">{jdwl.mata_kuliah?.sks} SKS</span>
					</div>
					<div class="mk-body">
						<div class="dosen-list" style="margin-bottom: 12px;">
							<div class="dosen-item">
								<div class="avatar" style="background: var(--primary-color);">
									{jdwl.dosen?.name?.charAt(0) || 'D'}
								</div>
								<span class="dosen-name">{jdwl.dosen?.name || 'Dosen'}</span>
							</div>
						</div>
						<p style="margin: 4px 0; font-size: 0.9rem; color: #475569;">
							<strong>{jdwl.mata_kuliah?.name}</strong>
						</p>
						<div style="display: flex; gap: 12px; font-size: 0.85rem; color: #64748b; margin-top: 12px; margin-bottom: 16px;">
							<span>📅 {jdwl.kelas?.hari}</span>
							<span>🕐 {jdwl.kelas?.jam_mulai} - {jdwl.kelas?.jam_selesai}</span>
						</div>
						<a href={`/mahasiswa/kelas/${jdwl.id}`} class="btn btn-primary" style="width: 100%; display: block; text-align: center; text-decoration: none;">Masuk Kelas</a>
					</div>
				</div>
			{/each}
		</div>
	{/if}

	<hr class="section-divider" style="margin: 3rem 0; border: none; border-top: 1px dashed #cbd5e1;" />

	<div class="page-header">
		<h1>Bursa KRS (Tersedia)</h1>
		<p>Daftar semua kelas yang tersedia untuk Anda ambil.</p>
	</div>

	{#if bursaKelasList.length === 0}
		<div class="empty-state glass-panel">
			<div class="empty-icon">
				<svg xmlns="http://www.w3.org/2000/svg" width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
					<path d="M4 19.5v-15A2.5 2.5 0 0 1 6.5 2H20v20H6.5a2.5 2.5 0 0 1 0-5H20"/>
				</svg>
			</div>
			<h2>Bursa Kelas Kosong</h2>
			<p>Tidak ada kelas yang disetujui di Program Studi Anda saat ini.</p>
		</div>
	{:else}
		<div class="grid-container">
			{#each bursaKelasList as bursa (bursa.id)}
				<div class="mk-card glass-panel" in:fade>
					<div class="mk-header">
						<h3 class="mk-title">{bursa.kelas?.name}</h3>
						<span class="sks-badge">{bursa.mata_kuliah?.sks} SKS</span>
					</div>
					<div class="mk-body">
						<div class="dosen-list" style="margin-bottom: 12px;">
							<div class="dosen-item">
								<div class="avatar" style="background: var(--primary-color);">
									{bursa.dosen?.name?.charAt(0) || 'D'}
								</div>
								<span class="dosen-name">{bursa.dosen?.name || 'Dosen'}</span>
							</div>
						</div>
						<p style="margin: 4px 0; font-size: 0.9rem; color: #475569;">
							<strong>{bursa.mata_kuliah?.name}</strong>
						</p>
						<div style="display: flex; gap: 12px; font-size: 0.85rem; color: #64748b; margin-top: 12px; margin-bottom: 16px;">
							<span>📅 {bursa.kelas?.hari}</span>
							<span>🕐 {bursa.kelas?.jam_mulai} - {bursa.kelas?.jam_selesai}</span>
						</div>
						
						<!-- Cek apakah sudah diambil -->
						{#if jadwalList.find(j => j.id === bursa.id)}
							<button class="btn" style="width: 100%; background: #e2e8f0; color: #64748b; cursor: not-allowed;" disabled>Sudah Diambil</button>
						{:else}
							<button class="btn btn-primary" style="width: 100%;" onclick={() => ambilKelas(bursa.id)}>Ambil Kelas</button>
						{/if}
					</div>
				</div>
			{/each}
		</div>
	{/if}
{/if}

<style>
	.page-header {
		margin-bottom: 2rem;
	}

	.page-header h1 {
		font-size: 1.75rem;
		font-weight: 700;
		color: var(--text-main);
		margin: 0 0 0.5rem 0;
	}

	.page-header p {
		color: var(--text-muted);
		margin: 0;
		font-size: 1.05rem;
	}

	.empty-state {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		padding: 5rem 2rem;
		text-align: center;
		border-radius: 16px;
		border: 1px dashed var(--border);
	}

	.empty-icon {
		display: flex;
		align-items: center;
		justify-content: center;
		width: 80px;
		height: 80px;
		background: white;
		border-radius: 50%;
		color: var(--primary);
		margin-bottom: 1.5rem;
		box-shadow: 0 4px 12px rgba(79, 70, 229, 0.15);
	}

	.grid-container {
		display: grid;
		grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
		gap: 1.5rem;
	}

	.mk-card {
		display: flex;
		flex-direction: column;
		padding: 1.5rem;
		border-radius: 12px;
		transition: transform 0.2s, box-shadow 0.2s;
	}

	.mk-card:hover {
		transform: translateY(-4px);
		box-shadow: var(--shadow-md);
	}

	.mk-header {
		display: flex;
		justify-content: space-between;
		align-items: flex-start;
		margin-bottom: 1rem;
		padding-bottom: 1rem;
		border-bottom: 1px solid var(--border);
	}

	.mk-title {
		font-size: 1.25rem;
		font-weight: 600;
		color: var(--text-main);
		margin: 0;
	}

	.sks-badge {
		background: rgba(79, 70, 229, 0.1);
		color: var(--primary);
		padding: 0.25rem 0.75rem;
		border-radius: 20px;
		font-size: 0.875rem;
		font-weight: 600;
	}

	.dosen-list {
		display: flex;
		flex-direction: column;
		gap: 0.75rem;
	}

	.dosen-item {
		display: flex;
		align-items: center;
		gap: 0.75rem;
	}

	.avatar {
		width: 32px;
		height: 32px;
		border-radius: 50%;
		background: linear-gradient(135deg, var(--primary), var(--primary-hover));
		color: white;
		display: flex;
		align-items: center;
		justify-content: center;
		font-weight: 600;
		font-size: 0.875rem;
	}

	.dosen-name {
		font-size: 0.95rem;
		color: var(--text-main);
		font-weight: 500;
	}

	.loading-state, .error-state {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		padding: 4rem 0;
		gap: 1rem;
		color: var(--text-muted);
	}

	.error-state {
		color: var(--danger);
	}

	.spinner {
		width: 40px;
		height: 40px;
		border: 3px solid rgba(79, 70, 229, 0.1);
		border-radius: 50%;
		border-top-color: var(--primary);
		animation: spin 1s linear infinite;
	}

	@keyframes spin {
		to { transform: rotate(360deg); }
	}

	.empty-state h2 {
		margin: 0 0 0.75rem 0;
		color: var(--text-main);
		font-size: 1.5rem;
	}

	.empty-state p {
		margin: 0;
		color: var(--text-muted);
		max-width: 500px;
		line-height: 1.6;
	}
</style>
