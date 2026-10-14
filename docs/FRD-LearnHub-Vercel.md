# FRD — LearnHub

**Versi:** 1.3 · **Tanggal:** 15 Oktober 2026 · **Status:** Draft kebutuhan target · **Acuan:** PRD_eLearning.md, DESIGN.md, dan [PRD katalog silabus](PRD-LearnHub-Vercel.md)

> FR-01–FR-22 adalah kebutuhan target dan prioritas rilis, bukan klaim bahwa fungsi tersebut sudah diimplementasikan. Baseline yang diverifikasi ada di [TRD](TRD-LearnHub-Vercel.md#baseline-implementasi-repositori).

## 1. Aktor dan hak akses

| Aktor | Kapabilitas utama |
|---|---|
| Pengunjung | Melihat halaman publik, kategori, pencarian, detail kursus dan preview gratis |
| Learner | Enrollment, mengakses materi sesuai hak, menyimpan progress, mengerjakan kuis, memberi ulasan dan mengunduh sertifikat |
| Instructor | Membuat/mengubah kursus miliknya, menyusun seksi/lesson, mengajukan moderasi, melihat analitik kursus |
| Admin | Mengelola user, kategori, katalog, kursus pending, transaksi dan audit |
| Payment provider | Mengirim notifikasi status transaksi melalui webhook yang diverifikasi |

## 2. Kebutuhan fungsional

| ID | Kebutuhan | Prioritas |
|---|---|---|
| FR-01 | Pengguna dapat registrasi, login, logout, verifikasi email, reset password, dan login Google | P0 |
| FR-02 | Sistem menerapkan RBAC server-side dan memeriksa kepemilikan resource | P0 |
| FR-03 | Pengunjung dapat menelusuri, mencari, memfilter, dan membuka detail kursus | P0 |
| FR-04 | Instruktur dapat membuat draft kursus, mengatur metadata, seksi, lesson dan urutan | P0 |
| FR-05 | Instruktur dapat menyimpan lesson video, teks, lampiran, kuis, dan proyek praktik | P0 |
| FR-06 | Instruktur mengajukan kursus untuk moderasi; admin menerima atau menolak dengan alasan | P0 |
| FR-07 | Learner dapat enrollment kursus gratis atau memperoleh akses setelah pembayaran sukses | P0 |
| FR-08 | Learner dapat melanjutkan video dari posisi terakhir dan menandai lesson selesai | P0 |
| FR-09 | Sistem menyimpan progress per lesson dan menghitung progress kursus | P0 |
| FR-10 | Kursus dapat memiliki level Basic/Intermediate/Advanced, prasyarat, bahasa, tag dan jalur belajar | P0 |
| FR-11 | Kuis memberikan hasil dan umpan balik sesuai aturan kursus; attempt dicatat | P1 |
| FR-12 | Learner yang memenuhi aturan penyelesaian dapat menerima sertifikat dengan kode verifikasi | P1 |
| FR-13 | Learner dapat mengulas kursus yang diikuti; satu ulasan aktif per learner/kursus | P1 |
| FR-14 | Sistem menyediakan checkout, histori pembayaran, status transaksi dan penanganan webhook idempoten | P1 |
| FR-15 | Admin dapat mengelola kategori, katalog teknologi, user, status kursus, promo dan laporan | P1 |
| FR-16 | Sistem mengirim notifikasi email dan in-app untuk event penting | P2 |
| FR-17 | Antarmuka mendukung Bahasa Indonesia dan Inggris, perangkat mobile, keyboard dan screen reader | P0 |
| FR-18 | Admin dapat mengubah silabus dengan versi; enrollment yang ada tetap menunjuk versi yang dipelajari | P1 |
| FR-19 | Learner dapat membuka editor kode di browser, menjalankan solusi pada runner terisolasi, dan menerima hasil dari test yang ditetapkan | P1 |
| FR-20 | Latihan coding memakai skenario kerja yang realistis dan dianonimkan, dengan tujuan, input/output, batasan, serta kriteria penerimaan yang jelas | P1 |
| FR-21 | Setiap latihan menyediakan lima hint bertingkat; learner membuka hint satu per satu dan progres hint dicatat per attempt | P1 |
| FR-22 | Pilot runner mendukung Go dan JavaScript/TypeScript; bahasa lain ditambahkan setelah evaluasi keamanan dan biaya | P1 |

## 3. Alur inti

### Belajar kursus

1. Pengunjung mencari jalur/kursus dan melihat prasyarat, level, bahasa, durasi perkiraan dan kurikulum.
2. Learner login/registrasi lalu enrollment gratis atau checkout.
3. Server memvalidasi enrollment/akses berbayar sebelum memberikan lesson privat atau URL media bertanda tangan.
4. Progress disimpan berkala dan saat lesson selesai; klien dapat melanjutkan dari posisi terakhir.
5. Setelah kriteria kursus terpenuhi, learner mengerjakan asesmen akhir dan menerima sertifikat jika lulus.

### Penerbitan kursus

1. Instructor membuat draft dan menyusun seksi/lesson.
2. Sistem memvalidasi metadata wajib, urutan, aset, transkrip/aksesibilitas yang disyaratkan, dan asesmen.
3. Instructor mengajukan review; admin menyetujui atau menolak dengan catatan.
4. Kursus approved dapat dipublikasikan; perubahan substantif setelah terbit dapat memerlukan review ulang.

### Pembayaran

1. Learner meminta checkout; server menghitung harga dan diskon dari data tepercaya.
2. Server membuat transaksi pending dan sesi pembayaran melalui provider.
3. Webhook diverifikasi signature, timestamp dan idempotency; browser redirect tidak dianggap bukti pembayaran.
4. Status sukses memberi entitlement; kegagalan/refund mencabut atau memperbarui entitlement sesuai kebijakan.

### Latihan coding interaktif

1. Learner membuka soal berbasis skenario kerja, membaca spesifikasi dan contoh, lalu mengedit starter code di browser.
2. Learner menjalankan atau mengirim solusi; runner mengeksekusi kode di lingkungan sekali pakai yang terisolasi dan mengembalikan hasil test serta batasan yang dilanggar.
3. Learner dapat membuka tepat lima hint secara berurutan; hint memberi petunjuk progresif tanpa langsung memberikan solusi.
4. Attempt, hasil, bahasa, versi soal, dan hint yang dibuka disimpan untuk progres belajar. Kode sumber tidak dijalankan di proses API.

## 4. Aturan bisnis

- Hanya course published yang dapat ditemukan publik.
- Draft/pending/rejected hanya dapat dilihat pemilik dan admin yang berwenang.
- Learner hanya dapat mengakses preview atau lesson kursus yang dimiliki/enrolled dan belum dicabut.
- Instructor hanya dapat mengubah kursus miliknya; admin memiliki kewenangan moderasi yang diaudit.
- Progress bersifat per learner dan lesson; update bersifat idempoten dan tidak boleh menurunkan completion tanpa aturan eksplisit.
- Level silabus mencerminkan kompetensi: Basic membangun fondasi dan aplikasi sederhana; Advanced menangani arsitektur, kualitas, keamanan, performa, skala atau topik spesialis.
- Sertifikat diterbitkan berdasarkan aturan penyelesaian yang eksplisit dan menyimpan versi kursus.
- Nilai, mata uang, pajak, refund window, komisi dan payout mengikuti kebijakan yang disetujui sebelum monetisasi.
- Setiap latihan memiliki starter code, spesifikasi input/output, sample case, hidden test, batas waktu/memori, rubrik, dan tepat lima hint berurutan.
- Skenario kerja berupa kasus umum yang dianonimkan; materi tidak menyalin kode privat, data pelanggan, atau soal interview proprietary dari perusahaan.
- Runner menolak akses jaringan, secret, host database, filesystem host, dan proses di luar sandbox; limit waktu, CPU, memori, output, serta ukuran submission diterapkan.
- Hasil hidden test tidak membocorkan test case atau expected output secara penuh; umpan balik menunjukkan kategori kegagalan yang membantu tanpa membocorkan solusi.

## 5. Data dan validasi penting

Entitas target: users, profiles, courses, course_versions, sections, lessons, lesson_assets, quizzes, questions, attempts, code_exercises, code_submissions, hint_unlocks, enrollments, lesson_progress, transactions, entitlements, reviews, certificates, categories, technology_tracks, notifications, audit_events. Tidak semua entitas tersebut tersedia di implementasi saat ini; lihat baseline TRD.

Validasi meliputi format dan keunikan email/slug, batas panjang input, MIME dan ukuran file, harga nonnegatif, rating 1–5, posisi lesson unik dalam seksi, kepemilikan resource, prasyarat enrollment, serta signature webhook. Data sensitif tidak dikembalikan di response publik.

## 6. Status dan error

- Course: `draft → pending → published`; admin dapat mengembalikan `pending → rejected`, dan published dapat diarsipkan.
- Payment target: `pending → success | failed | expired`; `success → refunded` hanya melalui alur refund yang diaudit. Nilai `success` mengikuti nama status pada model transaksi dan skema PRD induk.
- Error API memakai format konsisten: kode stabil, pesan yang aman untuk pengguna, request ID; detail internal hanya di log.
- Respons umum: 400 validasi, 401 autentikasi, 403 otorisasi, 404 resource, 409 konflik/idempotency, 422 aturan bisnis, 429 rate limit, 5xx kegagalan internal.

## 7. Acceptance criteria ringkas

- Akses role dan kepemilikan diperiksa di backend untuk setiap operasi terlindungi.
- Course published dapat dibaca publik; lesson terkunci tidak bocor melalui API, HTML, atau URL permanen.
- Progress bertahan setelah reload dan lintas perangkat setelah login.
- Webhook berulang tidak menggandakan transaksi atau entitlement.
- Submission hanya berjalan di runner terisolasi, mengembalikan hasil dalam limit yang ditentukan, dan tidak dapat mengakses resource platform.
- Tiap latihan memiliki tepat lima hint yang terbuka bertahap; learner tidak dapat membuka hint berikutnya sebelum hint sebelumnya.
- Silabus teknologi bisa ditelusuri berdasarkan kategori/level dan menunjukan prasyarat serta hasil belajar.
- Form dapat digunakan dengan keyboard, label yang dapat dibaca screen reader, pesan validasi, dan layout responsif.
