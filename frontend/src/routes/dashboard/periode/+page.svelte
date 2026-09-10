<script lang="ts">
	import { onMount } from 'svelte';
	import { periodeService, type PeriodeAkademik } from '$lib/services/periode';
	import { toast } from '$lib/stores/toast.svelte';
	import { authState } from '$lib/stores/auth.svelte';
	import Card from '$lib/components/ui/Card.svelte';
	import Skeleton from '$lib/components/ui/Skeleton.svelte';
	import EmptyState from '$lib/components/ui/EmptyState.svelte';

	let periodes = $state<PeriodeAkademik[]>([]);
	let loading = $state(true);
	let submitting = $state(false);
	let showForm = $state(false);
	let newPeriode = $state({ tahun: '2024/2025', jenis: 'ganjil' });

	async function fetchData() {
		loading = true;
		try {
			const res = await periodeService.getAll();
			if (res.success && res.data) {
				periodes = res.data;
			}
		} catch (err: any) {
			toast.error(err.message || 'Gagal memuat data periode');
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
		if (!newPeriode.tahun || !newPeriode.jenis) {
			toast.error('Harap isi semua kolom');
			return;
		}

		submitting = true;
		try {
			const res = await periodeService.create(newPeriode);
			if (res.success) {
				toast.success('Periode berhasil dibuat');
				showForm = false;
				fetchData();
			}
		} catch (err: any) {
			toast.error(err.message || 'Gagal membuat periode');
		} finally {
			submitting = false;
		}
	}

	async function handleActivate(id: string) {
		try {
			const res = await periodeService.activate(id);
			if (res.success) {
				toast.success('Periode berhasil diaktifkan');
				fetchData();
			}
		} catch (err: any) {
			toast.error(err.message || 'Gagal mengaktifkan periode');
		}
	}

	async function handleDelete(id: string) {
		if (!confirm('Hapus periode ini?')) return;
		try {
			const res = await periodeService.delete(id);
			if (res.success) {
				toast.success('Periode berhasil dihapus');
				fetchData();
			}
		} catch (err: any) {
			toast.error(err.message || 'Gagal menghapus periode');
		}
	}
</script>

<div class="page-container animate-fade-in">
	<div class="page-header">
		<div>
			<h1 class="page-title">Periode Akademik</h1>
			<p class="page-subtitle">Kelola tahun ajaran dan aktifkan periode Ganjil/Genap.</p>
		</div>
		<button class="btn btn-primary" onclick={() => showForm = !showForm}>
			{showForm ? 'Batal' : '+ Tambah Periode'}
		</button>
	</div>

	{#if showForm}
		<Card class="form-card animate-slide-up" style="margin-bottom: 24px; padding: 24px;">
			<h2 style="margin-bottom: 16px;">Periode Baru</h2>
			<form onsubmit={(e) => { e.preventDefault(); handleCreate(); }} class="periode-form">
				<div class="form-group">
					<label for="tahun">Tahun Ajaran (mis: 2024/2025)</label>
					<input type="text" id="tahun" bind:value={newPeriode.tahun} required class="input" placeholder="2024/2025" />
				</div>
				<div class="form-group">
					<label for="jenis">Jenis Semester</label>
					<select id="jenis" bind:value={newPeriode.jenis} class="input">
						<option value="ganjil">Ganjil</option>
						<option value="genap">Genap</option>
						<option value="pendek">Semester Pendek / Antara</option>
					</select>
				</div>
				<div class="form-actions">
					<button type="submit" class="btn btn-primary" disabled={submitting}>
						{submitting ? 'Menyimpan...' : 'Simpan Periode'}
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
	{:else if periodes.length === 0}
		<EmptyState message="Belum ada periode akademik. Silakan tambah." />
	{:else}
		<div class="grid-container">
			{#each periodes as p}
				<Card class={`periode-card ${p.is_active ? 'active-periode' : ''}`}>
					<div class="periode-header">
						<div>
							<h2 style="margin: 0;">{p.tahun} - {p.jenis.toUpperCase()}</h2>
							<p class="text-muted text-sm" style="margin-top: 4px;">Dibuat: {new Date(p.created_at).toLocaleDateString('id-ID')}</p>
						</div>
						{#if p.is_active}
							<span class="status-badge active">Aktif</span>
						{:else}
							<button class="btn-activate" onclick={() => handleActivate(p.id)}>Aktifkan</button>
						{/if}
					</div>
					
					<div class="periode-footer">
						<button class="btn-delete" onclick={() => handleDelete(p.id)} disabled={p.is_active}>
							Hapus
						</button>
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

	.periode-form {
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

	.periode-card {
		padding: 20px;
		display: flex;
		flex-direction: column;
		gap: 16px;
		border-top: 4px solid transparent;
		transition: all 0.3s ease;
	}

	.periode-card.active-periode {
		border-top-color: var(--primary-color);
		background: linear-gradient(to bottom right, rgba(255, 255, 255, 1), rgba(var(--primary-rgb), 0.03));
		box-shadow: 0 8px 30px rgba(var(--primary-rgb), 0.1);
	}

	.periode-header {
		display: flex;
		justify-content: space-between;
		align-items: flex-start;
		border-bottom: 1px solid var(--border-color);
		padding-bottom: 12px;
	}

	.periode-header h2 {
		font-size: 1.25rem;
		font-weight: 700;
	}

	.status-badge {
		padding: 4px 10px;
		border-radius: 20px;
		font-size: 0.75rem;
		font-weight: 600;
		background: var(--success-bg);
		color: var(--success-color);
	}

	.btn-activate {
		background: var(--secondary-color);
		color: var(--text-main);
		border: 1px solid var(--surface-border);
		padding: 4px 12px;
		border-radius: 6px;
		font-size: 0.8rem;
		font-weight: 500;
		cursor: pointer;
		transition: all 0.2s ease;
	}

	.btn-activate:hover {
		background: var(--primary-color);
		color: white;
		border-color: var(--primary-color);
	}

	.periode-footer {
		display: flex;
		justify-content: flex-end;
	}

	.btn-delete {
		background: transparent;
		color: var(--error-color);
		border: 1px solid rgba(239, 68, 68, 0.3);
		padding: 4px 12px;
		border-radius: 6px;
		font-size: 0.8rem;
		font-weight: 500;
		cursor: pointer;
		transition: all 0.2s ease;
	}

	.btn-delete:hover:not(:disabled) {
		background: var(--error-color);
		color: white;
	}

	.btn-delete:disabled {
		opacity: 0.5;
		cursor: not-allowed;
	}
</style>
