# 📚 Product Requirements Document (PRD)
## Platform E-Learning — *LearnHub*

**Versi:** 1.0.0  
**Tanggal:** 18 Juni 2026  
**Status:** Draft untuk Review  
**Tech Stack:** Next.js 14 · Go 1.22 · PostgreSQL 16 · Tailwind CSS 3

---

## 1. Ringkasan Produk

**LearnHub** adalah platform e-learning modern berbasis web yang memungkinkan instruktur membuat dan menjual kursus online, serta pelajar mengakses materi belajar secara fleksibel — kapan saja dan di mana saja. Platform ini dirancang untuk segmen pasar Indonesia dengan antarmuka multibahasa (Bahasa Indonesia & Inggris).

### Visi
> *"Menjadikan pendidikan berkualitas dapat diakses oleh siapa saja, di mana saja."*

### Misi
- Menyediakan platform belajar yang intuitif, cepat, dan terjangkau.
- Memberdayakan instruktur lokal untuk memonetisasi keahlian mereka.
- Mendukung pembelajaran terstruktur dengan fitur tracking progres yang canggih.

---

## 2. Tujuan Bisnis & Metrik Keberhasilan

| Metrik | Target (6 Bulan) | Target (12 Bulan) |
|---|---|---|
| Pengguna Terdaftar | 10.000 | 50.000 |
| Kursus Aktif | 200 | 1.000 |
| Monthly Active Users | 5.000 | 25.000 |
| Gross Revenue | Rp 50 juta | Rp 500 juta |
| Course Completion Rate | ≥ 30% | ≥ 45% |
| NPS (Net Promoter Score) | ≥ 40 | ≥ 60 |

---

## 3. User Persona

### 👨‍🎓 Persona 1 — Pelajar (Learner)
- **Nama:** Andi, 22 tahun, mahasiswa tingkat akhir
- **Kebutuhan:** Belajar skill baru (programming, desain, marketing) untuk meningkatkan daya saing kerja
- **Pain Point:** Materi di YouTube tidak terstruktur; biaya bootcamp terlalu mahal
- **Goals:** Mendapat sertifikasi yang diakui, belajar sesuai jadwal sendiri

### 👩‍🏫 Persona 2 — Instruktur (Instructor)
- **Nama:** Sari, 35 tahun, profesional di bidang Data Science
- **Kebutuhan:** Platform yang mudah untuk upload materi & mendapatkan penghasilan pasif
- **Pain Point:** Platform lain mengambil komisi besar (>30%); dashboard analitik kurang informatif
- **Goals:** Menjangkau lebih banyak siswa, mendapat penghasilan tambahan

### 🏢 Persona 3 — Admin Platform
- **Kebutuhan:** Mengelola konten, pengguna, pembayaran, dan laporan bisnis
- **Goals:** Memastikan kualitas kursus dan pengalaman pengguna tetap optimal

---

## 4. Fitur Utama

### 4.1 Modul Autentikasi & Manajemen Pengguna

| Fitur | Prioritas | Keterangan |
|---|---|---|
| Registrasi & Login (Email/Password) | P0 | JWT-based authentication |
| Login dengan Google OAuth | P0 | Social login |
| Verifikasi Email | P0 | Konfirmasi akun via email |
| Lupa Password / Reset Password | P0 | Token-based reset |
| Manajemen Profil | P1 | Avatar, bio, link sosial |
| Role-Based Access Control | P0 | Learner, Instructor, Admin |
| Two-Factor Authentication (2FA) | P2 | TOTP via Google Authenticator |

### 4.2 Modul Manajemen Kursus (Instructor)

| Fitur | Prioritas | Keterangan |
|---|---|---|
| Dashboard Instruktur | P0 | Overview statistik kursus |
| Buat & Edit Kursus | P0 | Rich text editor, thumbnail, kategori |
| Manajemen Seksi & Pelajaran | P0 | Drag-and-drop reorder |
| Upload Video Materi | P0 | Direct upload ke cloud storage |
| Tambah Materi Teks/PDF | P1 | Rich text & file attachment |
| Kuis & Ujian | P1 | Multiple choice, true/false |
| Penetapan Harga | P0 | Gratis, berbayar, atau diskon |
| Preview Kursus Gratis | P1 | Pelajaran pertama dapat diakses tanpa bayar |
| Sertifikat Penyelesaian | P1 | Auto-generated PDF |
| Analitik Kursus | P1 | Pendapatan, enrollments, rating |

### 4.3 Modul Pelajar (Learner)

| Fitur | Prioritas | Keterangan |
|---|---|---|
| Beranda & Rekomendasi Kursus | P0 | Filter berdasarkan kategori & popularitas |
| Halaman Detail Kursus | P0 | Deskripsi, kurikulum, review, instruktur |
| Pencarian & Filter Kursus | P0 | Full-text search, filter harga, kategori, rating |
| Pembelian Kursus | P0 | Integrasi Midtrans / Xendit |
| Video Player Kustom | P0 | Auto-resume, kecepatan putar, subtitle |
| Tracking Progres Belajar | P0 | Progress bar per kursus & pelajaran |
| Pengerjaan Kuis | P1 | Nilai instan setelah submit |
| Unduh Sertifikat | P1 | PDF dengan QR code verifikasi |
| Daftar Wishlist | P2 | Simpan kursus untuk dibeli nanti |
| Diskusi / Komentar | P1 | Forum per kursus |
| Ulasan & Rating Kursus | P0 | Bintang 1–5 + komentar |
| Notifikasi | P1 | Email & in-app notifications |

### 4.4 Modul Pembayaran

| Fitur | Prioritas | Keterangan |
|---|---|---|
| Pembelian Langsung | P0 | Kartu kredit, transfer bank, e-wallet |
| Kode Promo / Kupon | P1 | Diskon persentase atau nominal |
| Riwayat Transaksi | P0 | Learner & Instructor |
| Pencairan Dana Instruktur | P1 | Withdraw dengan jadwal tetap |
| Split Payment Otomatis | P1 | Platform fee + instruktur |
| Invoice & Bukti Pembayaran | P1 | PDF otomatis |

### 4.5 Modul Admin

| Fitur | Prioritas | Keterangan |
|---|---|---|
| Dashboard Admin | P0 | KPI bisnis (revenue, user, kursus) |
| Manajemen Pengguna | P0 | Suspend, verifikasi, ubah role |
| Moderasi Kursus | P0 | Approve/reject kursus baru |
| Manajemen Kategori | P0 | CRUD kategori & subkategori |
| Manajemen Promo | P1 | Buat & kelola kode promo |
| Laporan Keuangan | P1 | Export CSV/PDF |
| Log Aktivitas | P2 | Audit trail sistem |

---

## 5. Arsitektur Sistem

```
┌─────────────────────────────────────────────────────────┐
│                    CLIENT LAYER                         │
│   Next.js 14 (App Router) + Tailwind CSS + shadcn/ui   │
│   React Query · Zustand · React Hook Form               │
└──────────────────────┬──────────────────────────────────┘
                       │ HTTPS / REST + WebSocket
┌──────────────────────▼──────────────────────────────────┐
│                    API GATEWAY                          │
│              Nginx (Reverse Proxy + TLS)                │
└──────┬───────────────┬────────────────┬─────────────────┘
       │               │                │
┌──────▼──────┐ ┌──────▼──────┐ ┌──────▼──────┐
│  Auth       │ │  Course     │ │  Payment    │
│  Service    │ │  Service    │ │  Service    │
│  (Go)       │ │  (Go)       │ │  (Go)       │
└──────┬──────┘ └──────┬──────┘ └──────┬──────┘
       │               │                │
       └───────────────▼────────────────┘
                       │
       ┌───────────────▼────────────────┐
       │         PostgreSQL 16          │
       │  (Primary + Read Replicas)     │
       └────────────────────────────────┘
                       │
       ┌───────────────▼────────────────┐
       │    Supporting Services         │
       │  Redis (Cache + Sessions)      │
       │  MinIO / S3 (Media Storage)    │
       │  SMTP / Mailgun (Email)        │
       └────────────────────────────────┘
```

### 5.1 Frontend — Next.js 14 + Tailwind CSS

| Komponen | Teknologi |
|---|---|
| Framework | Next.js 14 (App Router) |
| Styling | Tailwind CSS 3.x + shadcn/ui |
| State Management | Zustand |
| Data Fetching | TanStack Query (React Query) v5 |
| Form Handling | React Hook Form + Zod |
| Video Player | Video.js atau custom HTML5 |
| Rich Text Editor | TipTap |
| Animasi | Framer Motion |
| Internasionalisasi | next-intl |
| Testing | Jest + React Testing Library + Playwright |

### 5.2 Backend — Go

| Komponen | Teknologi |
|---|---|
| HTTP Framework | Gin atau Echo |
| ORM / Query Builder | sqlc + pgx |
| Authentication | JWT (golang-jwt) + bcrypt |
| File Upload | Multipart + streaming ke S3 |
| Task Queue | Asynq (Redis-based) |
| Logging | Zap |
| Config | Viper |
| Testing | testify + httptest |
| API Documentation | Swagger (swaggo) |

### 5.3 Database — PostgreSQL 16

| Komponen | Teknologi |
|---|---|
| Database | PostgreSQL 16 |
| Migrasi | golang-migrate |
| Caching | Redis 7 |
| Connection Pool | pgBouncer |

---

## 6. Skema Database (Utama)

```sql
-- USERS
CREATE TABLE users (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email       VARCHAR(255) UNIQUE NOT NULL,
    password    VARCHAR(255),           -- nullable untuk OAuth
    full_name   VARCHAR(255) NOT NULL,
    avatar_url  TEXT,
    role        VARCHAR(20) NOT NULL DEFAULT 'learner', -- learner | instructor | admin
    is_verified BOOLEAN DEFAULT FALSE,
    created_at  TIMESTAMPTZ DEFAULT NOW(),
    updated_at  TIMESTAMPTZ DEFAULT NOW()
);

-- COURSES
CREATE TABLE courses (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    instructor_id UUID REFERENCES users(id),
    title        VARCHAR(500) NOT NULL,
    slug         VARCHAR(500) UNIQUE NOT NULL,
    description  TEXT,
    thumbnail_url TEXT,
    price        NUMERIC(12, 2) DEFAULT 0,
    level        VARCHAR(20),            -- beginner | intermediate | advanced
    status       VARCHAR(20) DEFAULT 'draft', -- draft | pending | published | archived
    category_id  UUID REFERENCES categories(id),
    language     VARCHAR(10) DEFAULT 'id',
    total_duration INTEGER,             -- in seconds
    created_at   TIMESTAMPTZ DEFAULT NOW(),
    updated_at   TIMESTAMPTZ DEFAULT NOW()
);

-- SECTIONS
CREATE TABLE sections (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    course_id  UUID REFERENCES courses(id) ON DELETE CASCADE,
    title      VARCHAR(500) NOT NULL,
    position   INTEGER NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- LESSONS
CREATE TABLE lessons (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    section_id   UUID REFERENCES sections(id) ON DELETE CASCADE,
    title        VARCHAR(500) NOT NULL,
    content_type VARCHAR(20) NOT NULL,  -- video | text | quiz | file
    video_url    TEXT,
    content      TEXT,
    duration     INTEGER,               -- in seconds
    position     INTEGER NOT NULL,
    is_preview   BOOLEAN DEFAULT FALSE,
    created_at   TIMESTAMPTZ DEFAULT NOW()
);

-- ENROLLMENTS
CREATE TABLE enrollments (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id      UUID REFERENCES users(id),
    course_id    UUID REFERENCES courses(id),
    enrolled_at  TIMESTAMPTZ DEFAULT NOW(),
    completed_at TIMESTAMPTZ,
    UNIQUE(user_id, course_id)
);

-- LESSON PROGRESS
CREATE TABLE lesson_progress (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id       UUID REFERENCES users(id),
    lesson_id     UUID REFERENCES lessons(id),
    is_completed  BOOLEAN DEFAULT FALSE,
    last_position INTEGER DEFAULT 0,    -- video resume in seconds
    updated_at    TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(user_id, lesson_id)
);

-- TRANSACTIONS
CREATE TABLE transactions (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id        UUID REFERENCES users(id),
    course_id      UUID REFERENCES courses(id),
    amount         NUMERIC(12, 2) NOT NULL,
    platform_fee   NUMERIC(12, 2) NOT NULL,
    status         VARCHAR(20) DEFAULT 'pending', -- pending | success | failed | refunded
    payment_method VARCHAR(50),
    gateway_ref    VARCHAR(255),
    paid_at        TIMESTAMPTZ,
    created_at     TIMESTAMPTZ DEFAULT NOW()
);

-- REVIEWS
CREATE TABLE reviews (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    UUID REFERENCES users(id),
    course_id  UUID REFERENCES courses(id),
    rating     SMALLINT CHECK (rating BETWEEN 1 AND 5),
    comment    TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(user_id, course_id)
);

-- CATEGORIES
CREATE TABLE categories (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name        VARCHAR(255) NOT NULL,
    slug        VARCHAR(255) UNIQUE NOT NULL,
    parent_id   UUID REFERENCES categories(id),
    icon_url    TEXT,
    created_at  TIMESTAMPTZ DEFAULT NOW()
);

-- CERTIFICATES
CREATE TABLE certificates (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id      UUID REFERENCES users(id),
    course_id    UUID REFERENCES courses(id),
    issued_at    TIMESTAMPTZ DEFAULT NOW(),
    cert_url     TEXT,
    verify_code  VARCHAR(100) UNIQUE NOT NULL,
    UNIQUE(user_id, course_id)
);
```

---

## 7. Desain API (REST)

### Autentikasi
| Method | Endpoint | Deskripsi |
|---|---|---|
| POST | `/api/v1/auth/register` | Registrasi pengguna baru |
| POST | `/api/v1/auth/login` | Login & dapatkan JWT |
| POST | `/api/v1/auth/logout` | Invalidasi token |
| POST | `/api/v1/auth/refresh` | Refresh access token |
| POST | `/api/v1/auth/forgot-password` | Kirim email reset |
| POST | `/api/v1/auth/reset-password` | Reset password dengan token |
| GET  | `/api/v1/auth/google` | Redirect ke Google OAuth |
| GET  | `/api/v1/auth/google/callback` | Callback OAuth Google |

### Kursus
| Method | Endpoint | Deskripsi |
|---|---|---|
| GET | `/api/v1/courses` | List kursus (publik, dengan filter) |
| GET | `/api/v1/courses/:slug` | Detail kursus |
| POST | `/api/v1/courses` | Buat kursus (Instructor) |
| PUT | `/api/v1/courses/:id` | Update kursus (Instructor) |
| DELETE | `/api/v1/courses/:id` | Hapus kursus (Instructor/Admin) |
| POST | `/api/v1/courses/:id/publish` | Submit untuk review (Instructor) |
| GET | `/api/v1/courses/:id/curriculum` | Kurikulum kursus |
| POST | `/api/v1/courses/:id/sections` | Tambah seksi |
| POST | `/api/v1/sections/:id/lessons` | Tambah pelajaran |
| GET | `/api/v1/courses/:id/reviews` | Ulasan kursus |
| POST | `/api/v1/courses/:id/reviews` | Beri ulasan (Learner) |

### Enrollment & Progress
| Method | Endpoint | Deskripsi |
|---|---|---|
| POST | `/api/v1/enrollments` | Enroll ke kursus |
| GET | `/api/v1/my/enrollments` | Kursus yang diikuti |
| GET | `/api/v1/my/enrollments/:courseId/progress` | Progres kursus |
| POST | `/api/v1/lessons/:id/progress` | Update progres pelajaran |

### Pembayaran
| Method | Endpoint | Deskripsi |
|---|---|---|
| POST | `/api/v1/payments/checkout` | Buat transaksi |
| GET | `/api/v1/payments/history` | Riwayat transaksi |
| POST | `/api/v1/payments/webhook` | Webhook payment gateway |

### Admin
| Method | Endpoint | Deskripsi |
|---|---|---|
| GET | `/api/v1/admin/dashboard` | Statistik platform |
| GET | `/api/v1/admin/users` | List pengguna |
| PATCH | `/api/v1/admin/users/:id` | Update status pengguna |
| GET | `/api/v1/admin/courses/pending` | Kursus menunggu review |
| POST | `/api/v1/admin/courses/:id/approve` | Approve kursus |
| POST | `/api/v1/admin/courses/:id/reject` | Reject kursus |

---

## 8. Struktur Proyek

```
learnhub/
├── frontend/                        # Next.js 14 App
│   ├── src/
│   │   ├── app/                     # App Router pages
│   │   │   ├── (auth)/              # Login, Register
│   │   │   ├── (learner)/           # Dashboard pelajar
│   │   │   ├── (instructor)/        # Dashboard instruktur
│   │   │   ├── admin/               # Panel admin
│   │   │   └── courses/             # Halaman kursus publik
│   │   ├── components/              # Shared components
│   │   │   ├── ui/                  # shadcn/ui components
│   │   │   ├── course/              # Course-specific components
│   │   │   ├── player/              # Video player
│   │   │   └── layout/              # Header, Footer, Sidebar
│   │   ├── lib/                     # Utilities, API clients
│   │   ├── hooks/                   # Custom React hooks
│   │   ├── stores/                  # Zustand stores
│   │   └── types/                   # TypeScript types
│   ├── tailwind.config.ts
│   ├── next.config.ts
│   └── package.json
│
├── backend/                         # Go Monorepo
│   ├── cmd/
│   │   └── api/                     # Entry point
│   ├── internal/
│   │   ├── auth/                    # Auth service
│   │   ├── course/                  # Course service
│   │   ├── payment/                 # Payment service
│   │   ├── user/                    # User service
│   │   ├── admin/                   # Admin service
│   │   └── middleware/              # JWT, CORS, rate limiting
│   ├── pkg/
│   │   ├── database/                # DB connection & migrations
│   │   ├── storage/                 # S3/MinIO client
│   │   ├── email/                   # Email service
│   │   └── utils/                   # Shared utilities
│   ├── db/
│   │   ├── migrations/              # SQL migration files
│   │   └── queries/                 # sqlc query files
│   ├── go.mod
│   └── go.sum
│
├── infra/                           # Infrastructure
│   ├── docker-compose.yml
│   ├── docker-compose.prod.yml
│   ├── nginx/
│   │   └── nginx.conf
│   └── k8s/                         # Kubernetes manifests (opsional)
│
└── docs/                            # Documentation
    ├── api/                         # Swagger/OpenAPI spec
    └── architecture/
```

---

## 9. Non-Functional Requirements

### 9.1 Performa
- API response time ≤ 200ms untuk endpoint utama (p95)
- Halaman web Time-to-Interactive (TTI) ≤ 3 detik
- Video playback mulai dalam ≤ 2 detik

### 9.2 Skalabilitas
- Mendukung hingga 10.000 concurrent users pada fase awal
- Arsitektur stateless pada backend untuk horizontal scaling
- CDN untuk aset statis dan video

### 9.3 Keamanan
- HTTPS wajib di semua endpoint
- JWT dengan expiry pendek (15 menit) + refresh token (7 hari)
- Rate limiting pada endpoint login & registrasi
- SQL injection prevention via parameterized queries (sqlc)
- XSS protection via proper sanitization
- CORS policy yang ketat
- File upload validation (tipe, ukuran)

### 9.4 Ketersediaan
- Target uptime: 99.5% (≤ 3.65 jam downtime/tahun)
- Health check endpoint pada semua service
- Graceful shutdown pada backend

### 9.5 Aksesibilitas
- WCAG 2.1 Level AA compliance
- Keyboard navigation support
- Screen reader compatibility

---

## 10. Rencana Pengembangan (Milestone)

### 🚀 Phase 1 — MVP (Bulan 1–3)
- [ ] Setup infrastruktur (Docker, DB, CI/CD)
- [ ] Autentikasi & manajemen pengguna
- [ ] CRUD kursus, seksi, dan pelajaran
- [ ] Upload & streaming video
- [ ] Enrollment kursus (gratis)
- [ ] Tracking progres belajar
- [ ] Halaman publik kursus & pencarian

### 💳 Phase 2 — Monetisasi (Bulan 4–5)
- [ ] Integrasi Midtrans / Xendit
- [ ] Sistem kursus berbayar
- [ ] Dashboard pendapatan instruktur
- [ ] Pencairan dana instruktur
- [ ] Kode promo & diskon

### 🎓 Phase 3 — Engagement (Bulan 6–8)
- [ ] Sistem kuis & ujian
- [ ] Sertifikat penyelesaian otomatis
- [ ] Sistem ulasan & rating
- [ ] Forum diskusi per kursus
- [ ] Notifikasi email & in-app

### ⚙️ Phase 4 — Admin & Scale (Bulan 9–12)
- [ ] Dashboard admin lengkap
- [ ] Moderasi konten
- [ ] Laporan keuangan & analitik
- [ ] Optimasi performa (caching, CDN)
- [ ] Mobile responsiveness audit
- [ ] Audit keamanan

---

## 11. Asumsi & Risiko

### Asumsi
- Instruktur bertanggung jawab atas kualitas konten mereka sendiri
- Pembayaran awal fokus pada pasar Indonesia (IDR)
- Hosting menggunakan cloud provider (AWS/GCP/DigitalOcean)
- Video hosting menggunakan S3-compatible storage

### Risiko

| Risiko | Dampak | Probabilitas | Mitigasi |
|---|---|---|---|
| Kapasitas video streaming melebihi estimasi | Tinggi | Sedang | Gunakan CDN + transcoding asynkron |
| Fraud pembayaran | Tinggi | Rendah | Integrasi fraud detection Midtrans |
| Kualitas kursus rendah | Sedang | Sedang | Proses review admin sebelum publish |
| Data breach | Tinggi | Rendah | Enkripsi, audit, penetration test |
| Lambatnya adopsi instruktur | Tinggi | Sedang | Program onboarding & insentif awal |

---

## 12. Open Questions untuk Review

> [!IMPORTANT]
> Pertanyaan berikut perlu dijawab sebelum memulai development:

1. **Payment Gateway**: Midtrans atau Xendit? Atau keduanya?
2. **Komisi Platform**: Berapa % yang diambil platform dari setiap transaksi instruktur?
3. **Video Hosting**: Hosting sendiri (MinIO/S3) atau menggunakan layanan pihak ketiga seperti Bunny.net / Mux?
4. **Deployment**: Apakah akan di-deploy ke VPS tunggal (Docker Compose) atau Kubernetes?
5. **Live Session**: Apakah ada fitur live class (webinar) di roadmap jangka pendek?
6. **Mobile App**: Apakah ada rencana mengembangkan native mobile app?
7. **Bahasa**: Apakah platform multibahasa (ID + EN) dari awal, atau hanya Bahasa Indonesia dulu?
8. **Logo & Nama Brand**: Sudahkah ada nama brand final selain "LearnHub"?

---

*Dokumen ini adalah living document dan akan diperbarui seiring perkembangan proyek.*
