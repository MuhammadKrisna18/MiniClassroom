# SIAKAD Pro System

<div align="center">
  <img src="https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white" alt="Go" />
  <img src="https://img.shields.io/badge/SvelteKit-FF3E00?style=for-the-badge&logo=svelte&logoColor=white" alt="SvelteKit" />
  <img src="https://img.shields.io/badge/PostgreSQL-316192?style=for-the-badge&logo=postgresql&logoColor=white" alt="PostgreSQL" />
  <img src="https://img.shields.io/badge/Redis-DC382D?style=for-the-badge&logo=redis&logoColor=white" alt="Redis" />
</div>

<br/>

Sebuah sistem informasi akademik terintegrasi (SIAKAD Pro) yang dikembangkan dengan fokus pada performa tinggi, struktur kode yang bersih (*Clean Architecture*), serta desain antarmuka modern (*Glassmorphism*).

> **Visi Proyek**: Mewujudkan infrastruktur sistem akademik yang *scalable*, aman, dan memberikan *user experience* terbaik melalui adaptasi teknologi Svelte 5 Runes dan kapabilitas tinggi Go Fiber.

---

## 📑 Daftar Isi
1. [Fitur Utama](#-fitur-utama)
2. [Teknologi & Kualitas](#-teknologi--kualitas)
3. [Arsitektur Sistem](#-arsitektur-sistem)
4. [Struktur Proyek](#-struktur-proyek)
5. [Instalasi & Penggunaan](#-instalasi--penggunaan)
6. [Alur Bisnis](#-alur-bisnis)

---

## 🚀 Fitur Utama

### 1. Manajemen Institusi Terpusat
- **Manajemen Program Studi & Kurikulum**: Pengelolaan kurikulum yang mengikat pada setiap program studi.
- **Manajemen Kelas Dinamis**: Mendukung jadwal *multi-schedule* dengan proteksi *conflict* (bentrok jadwal) secara *real-time*.
- **Portal Kelas & Sistem Absensi**: Dilengkapi **Manajemen Pertemuan Terpusat** (auto-increment pertemuan hingga maksimal 16), pembuatan **Kode Absensi 6-digit** dinamis berbatas waktu, dan terintegrasi dengan **Rekap Kehadiran Mahasiswa**.
- **Sistem Pengajuan Berjenjang**: Dosen mengajukan kelas dan mata kuliah, kemudian Admin melakukan *review* (Approve/Reject).

### 2. Validasi & Business Rules (Aturan Bisnis)
- **Maksimal Pertemuan**: Setiap kelas dibatasi maksimal 16 kali pertemuan. Sistem akan menolak pembuatan pertemuan ke-17.
- **Proteksi Jadwal Bentrok**: Dosen tidak dapat mengajukan kelas jika rentang waktu mengajar (hari dan jam) tumpang tindih dengan kelas lain yang sudah disetujui.
- **Otoritas Data Mengajar**: Mata kuliah yang diampu hanya dapat dikelola oleh dosen bersangkutan. Mahasiswa tidak memiliki hak untuk melihat kode absensi, mahasiswa hanya dapat mengirim kode yang didapatkan dari dosen.
- **Sistem Semester**: Data akademik dikaitkan dengan status semester aktif. Hanya satu semester yang bisa berstatus aktif pada satu waktu.
- **Perubahan Akun (NID & Email)**: Nomor Induk Dosen (NID) 5-digit akan digenerate permanen oleh sistem. Dosen juga hanya bisa melakukan pengubahan email melalui *Request Change* yang harus disetujui Admin.

### 3. Arsitektur & Keamanan
- **Single Source of Truth**: Logika bisnis dipusatkan secara eksklusif agar riwayat database selalu konsisten tanpa redundansi data.
- **Modern Authentication (RBAC)**: Autentikasi JWT yang memisahkan otorisasi antara **Admin**, **Dosen**, dan **Mahasiswa**.
- **Data Integrity**: Memanfaatkan *transactional operation* pada PostgreSQL (ACID) untuk mencegah anomali data.

### 4. User Experience (UX)
- **Role-based Dashboards**: Tampilan adaptif dengan navigasi *Glassmorphism*, transisi instan *client-side routing*.
- **Real-time Evaluasi**: Admin memiliki dasbor analitik untuk melihat utilitas akademik (misal: Mata Kuliah tanpa dosen).
- **Notifikasi Pintar**: *Toast Notification* global berbasis *state* Svelte 5 Runes yang reaktif.

---

## 🛠️ Teknologi & Kualitas

Sistem ini dibangun dengan memprioritaskan kebutuhan non-fungsional (*Non-Functional Requirements*).

| Komponen | Teknologi | Fokus Kualitas (*Quality Attributes*) |
| :--- | :--- | :--- |
| **Backend** | Go, Fiber (v2) | **Performance**: Eksekusi HTTP ultra-cepat, penggunaan (*footprint*) memori sangat rendah. |
| **Frontend** | Svelte 5 (Runes), Vite | **Usability**: UI reaktif dan efisien tanpa *overhead Virtual DOM*. |
| **Database** | PostgreSQL, GORM | **Data Integrity**: Skema selalu konsisten dengan fitur *Auto-Migrate*. |
| **Caching** | Redis | **Scalability**: Penyimpanan sesi dan token sementara berkecepatan tinggi. |
| **Keamanan** | JWT, Bcrypt | **Security**: Interaksi API *stateless* yang ketat dengan aturan *Role-Based Access*. |

---

## 🏗️ Arsitektur Sistem

Sistem ini menerapkan pola **Clean Architecture** pada Backend (Domain, Repository, Service, Delivery). Lapisan logika bisnis (Service) terisolasi penuh dari infrastruktur (Repository).

<details>
<summary><b>Lihat Diagram Arsitektur (Data Flow)</b></summary>

```mermaid
sequenceDiagram
    participant C as Client (Frontend)
    participant H as HTTP Delivery (Handler)
    participant S as Service Layer
    participant R as Repository Layer
    participant D as Database

    C->>H: 1. Request API
    H->>S: 2. Panggil fungsi Service (bawa struct domain)
    S->>S: 3. Eksekusi Logika Bisnis & Validasi
    S->>R: 4. Akses Data
    R->>D: 5. Query / Command (SQL)
    D-->>R: 6. Kembalikan Data Mentah
    R-->>S: 7. Kembalikan Model Domain
    S-->>H: 8. Kembalikan Response DTO
    H-->>C: 9. HTTP Response (JSON)
```
</details>

---

## 📂 Struktur Proyek

Seluruh komponen disusun agar mudah dikembangkan dan dimodifikasi (*Maintainability* tinggi).

<details>
<summary><b>Jelajahi Hierarki Direktori Utama</b></summary>

```text
📦 Project Root
 ┣ 📂 cmd/api              # Main entry point aplikasi Go
 ┣ 📂 config               # Konfigurasi environment (DB, Redis, JWT)
 ┣ 📂 docs                 # Dokumentasi auto-generated (Swagger)
 ┣ 📂 frontend             # Direktori frontend SvelteKit
 ┃  ┣ 📂 src/lib
 ┃  ┃  ┣ 📂 components     # Dashboard, Splash, UI Reusable (Card, Modal, Badge)
 ┃  ┃  ┣ 📂 services       # API Service layer untuk fetch data
 ┃  ┃  ┗ 📂 stores         # Global State Svelte 5 Runes (auth, toast)
 ┃  ┣ 📂 src/routes        # Halaman Aplikasi berbasis File-system Routing
 ┃  ┗ 📜 app.css           # File core CSS dengan variable desain
 ┣ 📂 internal             # Core logic dari backend SIAKAD Pro
 ┃  ┣ 📂 app               # Registrasi aplikasi dan middleware (Fiber)
 ┃  ┣ 📂 modules           # Modul domain (auth, kelas, matakuliah, programstudi)
 ┃  ┗ 📂 shared            # Logic shared (DB, Cache, Error Handler, Response)
```
</details>

---

## 💻 Instalasi & Penggunaan

### 1. Persyaratan Sistem
Pastikan sistem Anda telah memiliki instalasi:
- **Go** (v1.20+)
- **Node.js** (v18+)
- **PostgreSQL** (berjalan pada port `5432`)
- **Redis** (berjalan pada port `6379`)

### 2. Konfigurasi & Menjalankan Backend
Buat file konfigurasi `.env` di direktori utama (Root):
```env
PORT=8080
ENV=development

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=postgres
DB_SSLMODE=disable

REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

JWT_SECRET=supersecret-jwt-key
```

Lalu nyalakan server *Backend*:
```bash
go run cmd/api/main.go
# Server berjalan di: http://localhost:8080
# Dokumentasi Swagger: http://localhost:8080/api/v1/swagger/
```

### 3. Menjalankan Frontend
Buka terminal/konsol baru, navigasi ke direktori frontend, instal dependensi, lalu jalankan:
```bash
cd frontend
npm install
npm run dev
# Frontend berjalan di: http://localhost:5173
```

---

## ⚖️ Alur Bisnis

Berikut adalah alur operasional utama SIAKAD Pro mulai dari registrasi dosen hingga pengelolaan sesi kelas.

<details>
<summary><b>Lihat Diagram Flowchart Aturan Bisnis</b></summary>

```mermaid
flowchart TD
    %% Roles
    Dosen([Dosen])
    Admin([Admin])

    %% Business Entities & Actions
    subgraph Pengajuan Mata Kuliah & Kelas
        ReqMK(Dosen mengajukan Mata Kuliah)
        CheckMK{Status MK?}
        StatusPendingMK(Status MK: Pending)
        ApproveMK(Admin Approve MK)
        MKTaken(MK Menjadi Milik Dosen)
        
        ReqKelas(Dosen mengajukan Kelas)
        CheckKelas{Jadwal & Prodi Valid?}
        StatusPendingKelas(Status Kelas: Pending)
        ApproveKelas(Admin Approve Kelas)
        KelasTaken(Portal Kelas Terbuka)
    end

    subgraph Manajemen Akun
        RegDosen(Registrasi Akun Baru)
        GenNID(NID 5-Digit Permanen Dibuat)
        ReqEmail(Dosen Request Ganti Email)
        ApproveEmail(Admin Approve Request)
    end

    %% Flows
    Dosen --> RegDosen
    RegDosen --> GenNID
    Dosen --> ReqEmail
    ReqEmail --> ApproveEmail
    Admin --> ApproveEmail

    Dosen --> ReqMK
    ReqMK --> CheckMK
    CheckMK -- Belum Diambil --> StatusPendingMK
    CheckMK -- Sudah Diambil --> RejectSystem(Ditolak Sistem)
    StatusPendingMK --> Admin
    Admin -- Setuju --> ApproveMK
    ApproveMK --> MKTaken

    MKTaken --> ReqKelas
    ReqKelas --> CheckKelas
    CheckKelas -- Valid --> StatusPendingKelas
    CheckKelas -- Bentrok/Prodi Beda --> RejectSystem
    StatusPendingKelas --> Admin
    Admin -- Setuju --> ApproveKelas
    ApproveKelas --> KelasTaken

    style GenNID fill:#dcfce7,stroke:#166534
    style MKTaken fill:#dcfce7,stroke:#166534
    style KelasTaken fill:#dcfce7,stroke:#166534
    style RejectSystem fill:#fee2e2,stroke:#991b1b
```
</details>

---
*Dikembangkan untuk arsitektur modern web dengan standar kebersihan kode dan fungsionalitas akademik tinggi.*
