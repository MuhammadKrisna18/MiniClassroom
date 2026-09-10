<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import { onMount } from 'svelte';
	import { authState } from '$lib/stores/auth.svelte';

	let { children } = $props();

	onMount(() => {
		if (!authState.token) {
			goto('/login');
		}
	});

	function handleLogout(e: Event) {
		e.preventDefault();
		authState.logout();
	}
</script>

<div class="layout-wrapper">

	<aside class="sidebar glass-panel">
		<div class="sidebar-brand">
			<div class="brand-logo">
				<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="m12 3-1.912 5.813a2 2 0 0 1-1.275 1.275L3 12l5.813 1.912a2 2 0 0 1 1.275 1.275L12 21l1.912-5.813a2 2 0 0 1 1.275-1.275L21 12l-5.813-1.912a2 2 0 0 1-1.275-1.275L12 3Z"/></svg>
			</div>
			<div class="brand-text">
				<h1>SIAKAD</h1>
				<p>Pro</p>
			</div>
		</div>

		<nav class="sidebar-nav">
			<a href="/dashboard" class="nav-item {$page.url.pathname === '/dashboard' ? 'active' : ''}">
				<svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect width="7" height="9" x="3" y="3" rx="1"/><rect width="7" height="5" x="14" y="3" rx="1"/><rect width="7" height="9" x="14" y="12" rx="1"/><rect width="7" height="5" x="3" y="16" rx="1"/></svg>
				<span>Dashboard</span>
			</a>
			
			{#if authState.role === 'admin'}
				<div class="nav-section-title">ADMINISTRASI</div>
				<a href="/dashboard/dosen" class="nav-item {$page.url.pathname.includes('/dashboard/dosen') ? 'active' : ''}">
					<svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2"/><circle cx="9" cy="7" r="4"/><path d="M22 21v-2a4 4 0 0 0-3-3.87"/><path d="M16 3.13a4 4 0 0 1 0 7.75"/></svg>
					<span>Dosen</span>
				</a>
				<a href="/dashboard/mahasiswa" class="nav-item {$page.url.pathname.includes('/dashboard/mahasiswa') ? 'active' : ''}">
					<svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2"/><circle cx="9" cy="7" r="4"/><line x1="19" y1="8" x2="19" y2="14"/><line x1="16" y1="11" x2="22" y2="11"/></svg>
					<span>Mahasiswa</span>
				</a>
				<a href="/dashboard/evaluasi" class="nav-item {$page.url.pathname.includes('/dashboard/evaluasi') ? 'active' : ''}">
					<svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/><polyline points="14 2 14 8 20 8"/><line x1="16" y1="13" x2="8" y2="13"/><line x1="16" y1="17" x2="8" y2="17"/><polyline points="10 9 9 9 8 9"/></svg>
					<span>Evaluasi</span>
				</a>
			{/if}
			
			{#if authState.role === 'admin'}
				<div class="nav-section-title">MANAJEMEN AKADEMIK</div>
				<a href="/dashboard/semester" class="nav-item {$page.url.pathname === '/dashboard/semester' ? 'active' : ''}">
					<svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M4 19.5v-15A2.5 2.5 0 0 1 6.5 2H20v20H6.5a2.5 2.5 0 0 1 0-5H20"/><path d="M8 7h6"/><path d="M8 11h8"/></svg>
					<span>Semester</span>
				</a>
				<a href="/dashboard/kurikulum" class="nav-item {$page.url.pathname.includes('/dashboard/kurikulum') ? 'active' : ''}">
					<svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M4 22h14a2 2 0 0 0 2-2V7.5L14.5 2H6a2 2 0 0 0-2 2v4"/><polyline points="14 2 14 8 20 8"/><path d="m3 15 2 2 4-4"/></svg>
					<span>Kurikulum (Fix)</span>
				</a>
			{/if}

			<div class="nav-section-title">AKADEMIK</div>
			<a href="/dashboard/matakuliah" class="nav-item {$page.url.pathname.includes('/dashboard/matakuliah') ? 'active' : ''}">
				<svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M4 19.5v-15A2.5 2.5 0 0 1 6.5 2H20v20H6.5a2.5 2.5 0 0 1 0-5H20"/></svg>
				<span>Mata Kuliah</span>
			</a>
			<a href="/dashboard/kelas" class="nav-item {$page.url.pathname.includes('/dashboard/kelas') ? 'active' : ''}">
				<svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="3" width="18" height="18" rx="2" ry="2"/><line x1="3" y1="9" x2="21" y2="9"/><line x1="9" y1="21" x2="9" y2="9"/></svg>
				<span>Kelas</span>
			</a>
		</nav>

		<div class="sidebar-footer">
			<div class="sidebar-user-info">
				<div class="sidebar-user-avatar">
					{#if authState.profile?.photo_url}
						<img src={`http://localhost:8080${authState.profile.photo_url}`} alt="Avatar" class="sidebar-avatar-img" />
					{:else}
						{authState.profile?.name?.charAt(0).toUpperCase() || 'U'}
					{/if}
				</div>
				<div class="sidebar-user-details">
					<span class="sidebar-user-name">{authState.profile?.name || 'User'}</span>
					<span class="sidebar-user-role">{authState.role?.toUpperCase() || 'USER'}</span>
				</div>
			</div>
			<button class="btn-logout" onclick={handleLogout}>
				<svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"/><polyline points="16 17 21 12 16 7"/><line x1="21" x2="9" y1="12" y2="12"/></svg>
				<span>Log out</span>
			</button>
		</div>
	</aside>

	<div class="main-wrapper">

		<header class="top-header glass-panel">
			<div class="header-breadcrumb">
				<span class="text-muted">Home</span>
				<span class="separator">/</span>
				<span class="current-page">
					{#if $page.url.pathname === '/dashboard'}
						Dashboard
					{:else if $page.url.pathname.includes('/dashboard/dosen')}
						Manajemen Dosen
					{:else if $page.url.pathname.includes('/dashboard/evaluasi')}
						Evaluasi Akademik
					{:else if $page.url.pathname.includes('/dashboard/matakuliah')}
						Mata Kuliah
					{:else if $page.url.pathname.includes('/dashboard/kelas')}
						Kelas
					{:else if $page.url.pathname.includes('/dashboard/profile')}
						Profil
					{/if}
				</span>
			</div>
			
			<a href="/dashboard/profile" class="header-profile">
				<div class="user-role-badge">
					{authState.role?.toUpperCase() || 'USER'}
				</div>
				<div class="user-avatar">
					{#if authState.profile?.photo_url}
						<img src={`http://localhost:8080${authState.profile.photo_url}`} alt="Avatar" class="avatar-img-nav" />
					{:else}
						{authState.profile?.name?.charAt(0).toUpperCase() || 'U'}
					{/if}
				</div>
			</a>
		</header>

		<main class="content-area">
			{@render children()}
		</main>
	</div>
</div>

<style>
	.layout-wrapper {
		display: flex;
		min-height: 100vh;
		background: var(--bg-color);
	}

	.sidebar {
		width: 260px;
		display: flex;
		flex-direction: column;
		background: var(--surface-color);
		border-right: 1px solid var(--surface-border);
		position: sticky;
		top: 0;
		height: 100vh;
		z-index: 40;
	}

	.sidebar-brand {
		padding: 32px 24px;
		display: flex;
		align-items: center;
		gap: 16px;
	}

	.brand-logo {
		width: 34px;
		height: 34px;
		background: var(--primary-color);
		border-radius: var(--radius-sm);
		display: flex;
		align-items: center;
		justify-content: center;
		color: white;
		box-shadow: 0 4px 10px rgba(79, 70, 229, 0.3);
	}

	.brand-text h1 {
		font-family: 'Outfit', sans-serif;
		font-size: 1.35rem;
		font-weight: 700;
		color: var(--text-main);
		line-height: 1.1;
		letter-spacing: -0.02em;
	}

	.brand-text p {
		font-family: 'Outfit', sans-serif;
		font-size: 0.8rem;
		color: var(--primary-color);
		font-weight: 700;
		text-transform: uppercase;
		letter-spacing: 0.15em;
		margin-top: 2px;
	}

	.sidebar-nav {
		flex: 1;
		padding: 16px;
		display: flex;
		flex-direction: column;
		gap: 6px;
		overflow-y: auto;
	}

	.nav-section-title {
		font-size: 0.75rem;
		font-weight: 700;
		color: var(--text-muted);
		text-transform: uppercase;
		letter-spacing: 0.08em;
		padding: 24px 12px 10px 12px;
	}

	.nav-item {
		display: flex;
		align-items: center;
		gap: 12px;
		padding: 10px 14px;
		border-radius: var(--radius-md);
		color: var(--text-muted);
		text-decoration: none;
		font-weight: 500;
		font-size: 0.95rem;
		transition: all 0.2s ease;
	}

	.nav-item:hover {
		background: var(--secondary-hover);
		color: var(--text-main);
	}

	.nav-item.active {
		background: var(--primary-light);
		color: var(--primary-color);
		font-weight: 600;
	}

	.sidebar-footer {
		padding: 16px;
		border-top: 1px solid var(--surface-border);
		display: flex;
		flex-direction: column;
		gap: 12px;
	}

	.sidebar-user-info {
		display: flex;
		align-items: center;
		gap: 12px;
		padding: 8px;
	}

	.sidebar-user-avatar {
		width: 36px;
		height: 36px;
		border-radius: 50%;
		background: #111827;
		color: white;
		display: flex;
		align-items: center;
		justify-content: center;
		font-weight: 600;
		font-size: 0.85rem;
		overflow: hidden;
		flex-shrink: 0;
	}

	.sidebar-avatar-img {
		width: 100%;
		height: 100%;
		object-fit: cover;
	}

	.sidebar-user-details {
		display: flex;
		flex-direction: column;
		min-width: 0;
	}

	.sidebar-user-name {
		font-size: 0.85rem;
		font-weight: 600;
		color: var(--text-main);
		white-space: nowrap;
		overflow: hidden;
		text-overflow: ellipsis;
	}

	.sidebar-user-role {
		font-size: 0.7rem;
		font-weight: 500;
		color: var(--text-muted);
		letter-spacing: 0.05em;
	}

	.btn-logout {
		display: flex;
		align-items: center;
		justify-content: center;
		gap: 8px;
		width: 100%;
		background: transparent;
		color: var(--text-muted);
		border: 1px solid var(--surface-border);
		padding: 12px;
		border-radius: var(--radius-md);
		font-family: inherit;
		font-size: 0.9rem;
		font-weight: 600;
		cursor: pointer;
		transition: all 0.2s;
	}

	.btn-logout:hover {
		background: rgba(239, 68, 68, 0.05);
		color: var(--error-color);
		border-color: rgba(239, 68, 68, 0.2);
	}

	.main-wrapper {
		flex: 1;
		display: flex;
		flex-direction: column;
		min-width: 0;
		background: var(--surface-color);
	}

	
	.top-header {
		height: 64px;
		padding: 0 40px;
		display: flex;
		align-items: center;
		justify-content: space-between;
		border-bottom: 1px solid var(--surface-border);
		position: sticky;
		top: 0;
		z-index: 30;
		background: rgba(255, 255, 255, 0.8);
		backdrop-filter: blur(12px);
		-webkit-backdrop-filter: blur(12px);
	}

	.header-breadcrumb {
		display: flex;
		align-items: center;
		gap: 8px;
		font-size: 0.95rem;
		font-weight: 500;
	}

	.separator {
		color: var(--surface-border);
	}

	.current-page {
		color: var(--text-main);
		font-weight: 600;
	}

	.header-profile {
		display: flex;
		align-items: center;
		gap: 16px;
		text-decoration: none;
		cursor: pointer;
		padding: 4px 8px;
		border-radius: var(--radius-full);
		transition: background-color 0.2s;
	}

	.header-profile:hover {
		background-color: rgba(0, 0, 0, 0.03);
	}

	.user-role-badge {
		font-size: 0.7rem;
		font-weight: 600;
		color: var(--text-main);
		background: #f3f4f6;
		padding: 4px 10px;
		border-radius: var(--radius-full);
		letter-spacing: 0.05em;
	}

	.user-avatar {
		width: 32px;
		height: 32px;
		border-radius: 50%;
		background: #111827;
		color: white;
		display: flex;
		align-items: center;
		justify-content: center;
		font-weight: 500;
		font-size: 0.9rem;
		overflow: hidden;
	}

	.avatar-img-nav {
		width: 100%;
		height: 100%;
		object-fit: cover;
		aspect-ratio: 1/1;
	}

	.content-area {
		flex: 1;
		padding: 32px;
		width: 100%;
	}
</style>
