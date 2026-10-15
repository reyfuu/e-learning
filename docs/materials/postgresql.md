# Materi: PostgreSQL untuk Basis Data Relasional Skala Produksi

**Format:** teks · **Level:** Basic–Advanced · **Estimasi:** 60–90 menit

## Tujuan belajar

Setelah materi ini, pelajar dapat merancang, mengoptimalkan, dan mengoperasikan basis data PostgreSQL (v15/v16): memahami integritas data relasional (ACID, constraints, foreign keys), menguasai indexing tingkat lanjut (B-Tree, GIN, Partial Index), menganalisis rencana eksekusi query dengan `EXPLAIN (ANALYZE, BUFFERS)`, mengelola concurrency dengan tingkat isolasi transaksi (*MVCC & Locks*), serta memahami partisi tabel dan manajemen backup.

## Prasyarat

- Memahami konsep dasar relational database dan sintaks dasar SQL (SELECT, INSERT, UPDATE, DELETE).
- Memiliki PostgreSQL terpasang di sistem atau menggunakan container Docker PostgreSQL.
- Alat antarmuka query terpasang (`psql` CLI, DBeaver, atau pgAdmin).

## 1. Dasar: arsitektur PostgreSQL dan MVCC

PostgreSQL adalah sistem basis data relasional objek (*object-relational database management system* / ORDBMS) yang paling tangguh dan kaya fitur di dunia open-source. PostgreSQL menerapkan standar SQL yang sangat ketat dan mengusung arsitektur **MVCC (Multi-Version Concurrency Control)**.

### Karakteristik Utama:
- **MVCC:** Pembacaan data tidak pernah memblokir penulisan data, dan penulisan data tidak pernah memblokir pembacaan. Setiap transaksi melihat snapshot data yang konsisten pada waktu tertentu.
- **Tipe Data Ekstensif:** Mendukung data terstruktur (INT, VARCHAR, TIMESTAMP), semi-terstruktur (`JSONB`), biner, UUID, hingga geospasial (`PostGIS`).
- **Ekstensibilitas Tinggi:** Mendukung procedural languages (PL/pgSQL), indeks kustom, dan plugin ekstensi (`pg_stat_statements`, `pgcrypto`).

## 2. Praktik dasar: Skema Relasional, Constraints, dan JSONB

```sql
-- 1. Pembuatan Ekstensi UUID
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- 2. Pembuatan Tabel Users & Courses dengan Constraint Ketat
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) NOT NULL UNIQUE,
    full_name VARCHAR(150) NOT NULL,
    role VARCHAR(50) NOT NULL CHECK (role IN ('student', 'instructor', 'admin')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE courses (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    instructor_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    title VARCHAR(200) NOT NULL,
    slug VARCHAR(200) NOT NULL UNIQUE,
    price NUMERIC(12, 2) NOT NULL CHECK (price >= 0),
    is_published BOOLEAN NOT NULL DEFAULT FALSE,
    metadata JSONB DEFAULT '{}'::jsonb, -- Dukungan JSONB fleksibel
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Indeks B-Tree untuk pencarian cepat dan relasi FK
CREATE INDEX idx_courses_instructor_id ON courses(instructor_id);
CREATE INDEX idx_courses_published ON courses(is_published) WHERE is_published = TRUE; -- Partial Index

-- Indeks GIN untuk pencarian di dalam dokumen JSONB
CREATE INDEX idx_courses_metadata ON courses USING gin (metadata);
```

Manipulasi data JSONB:
```sql
-- Masukkan kursus dengan atribut JSONB
INSERT INTO courses (instructor_id, title, slug, price, is_published, metadata)
VALUES (
    'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11',
    'Mastering PostgreSQL 16',
    'mastering-postgresql-16',
    350000.00,
    TRUE,
    '{"level": "advanced", "duration_hours": 12, "tags": ["database", "sql", "performance"]}'
);

-- Query data berdasarkan atribut di dalam JSONB
SELECT title, metadata->>'level' AS level
FROM courses
WHERE metadata->'tags' ? 'performance';
```

## 3. Fitur utama dan pola arsitektur

### Analisis Query Plan dengan `EXPLAIN (ANALYZE, BUFFERS)`
Jangan menebak penyebab query lambat; baca rencana eksekusinya:

```sql
EXPLAIN (ANALYZE, BUFFERS, COSTS)
SELECT * FROM courses
WHERE is_published = TRUE
ORDER BY created_at DESC
LIMIT 10;
```

Parameter penting output EXPLAIN:
- **Seq Scan (Sequential Scan):** PostgreSQL membaca seluruh tabel dari disk baris demi baris (buruk pada tabel besar).
- **Index Scan / Index Only Scan:** PostgreSQL menggunakan indeks B-Tree untuk langsung melompat ke blok data yang relevan (sangat cepat).
- **Buffers: shared hit vs read:** `shared hit` berarti data sudah ada di RAM (shared buffer), sedangkan `read` berarti Postgres harus membaca dari disk fisik.

### Concurrency dan Tingkat Isolasi Transaksi
PostgreSQL mendukung 3 tingkat isolasi transaksi standar:
1. `Read Committed` (default): Hanya membaca data yang sudah di-commit.
2. `Repeatable Read`: Menjamin snapshot data tidak berubah selama transaksi berjalan.
3. `Serializable`: Tingkat tertinggi yang menjamin transaksi berperilaku seolah-olah dieksekusi secara berurutan (*strictly sequential*).

## 4. Praktik kerja: optimasi query lambat dan mitigasi Table Bloat dengan VACUUM

### Skenario
Tabel riwayat aktivitas `user_activity_logs` berisi 10 juta baris. Query pencarian aktivitas pengguna memakan waktu 4,2 detik dan disk server cepat penuh meskipun banyak data lama yang telah dihapus.

### Tugas
Temukan penyebab ketiadaan indeks, buat indeks B-Tree yang tepat, serta pahami cara kerja `VACUUM` untuk mereklamasi ruang kosong (*table bloat*) yang ditinggalkan operasi UPDATE/DELETE.

```sql
-- 1. Buat Composite Index untuk mengeliminasi Sequential Scan
CREATE INDEX CONCURRENTLY idx_user_activity_user_date 
ON user_activity_logs (user_id, created_at DESC);

-- 2. Bersihkan dead tuples tanpa mengunci tabel
VACUUM (VERBOSE, ANALYZE) user_activity_logs;
```

### Lima hint bertahap
1. Operasi `UPDATE` di PostgreSQL sebenarnya menyisipkan tuple baru dan menandai tuple lama sebagai "dead" (*MVCC append-only*).
2. Jika dead tuples menumpuk tanpa dibersihkan oleh worker Autovacuum, tabel mengalami fenomena *bloat*.
3. Selalu buat indeks di tabel produksi yang sedang aktif dengan keyword `CREATE INDEX CONCURRENTLY` agar tidak mengunci penulisan data (*exclusive lock*).
4. Pastikan daemon `autovacuum` selalu aktif di konfigurasi `postgresql.conf`.
5. Komposit indeks `(A, B)` dapat digunakan untuk query `WHERE A = ?` atau `WHERE A = ? AND B = ?`, tetapi tidak efektif untuk query `WHERE B = ?`.

## 5. Lanjutan: Partisi Tabel (Table Partitioning), Replikasi, dan PgBouncer

- **Table Partitioning:** Pecah tabel berukuran puluhan gigabyte berdasarkan rentang tanggal (`PARTITION BY RANGE (created_at)`) agar query hanya membaca partisi yang relevan (*partition pruning*).
- **Connection Pooling dengan PgBouncer:** Setiap koneksi baru ke PostgreSQL membuat satu proses Linux (`fork`). Gunakan PgBouncer dalam mode *transaction pooling* untuk melayani ribuan koneksi konkuren dengan RAM minimal.
- **Physical Streaming Replication:** Replikasi data real-time menggunakan Write-Ahead Log (WAL) ke standby replica server untuk failover ketersediaan tinggi (*High Availability*).

## Pemeriksaan pemahaman

1. Bagaimana cara kerja MVCC di PostgreSQL saat terjadi operasi UPDATE?
2. Mengapa pembuatan indeks pada tabel produksi besar wajib menggunakan klausa `CONCURRENTLY`?
3. Apa perbedaan antara tipe data `JSON` dan `JSONB` di PostgreSQL?
4. Kapan connection pooler seperti PgBouncer mutlak diperlukan?

**Kunci ringkas:** MVCC tidak menimpa baris lama melainkan membuat baris versi baru dan menandai baris lama sebagai dead tuple; `CONCURRENTLY` membangun indeks tanpa mengunci operasi INSERT, UPDATE, dan DELETE pada tabel sehingga layanan tidak mengalami downtime; `JSON` menyimpan teks mentah sedangkan `JSONB` mengompilasi JSON ke format biner terindeks yang jauh lebih cepat untuk di-query; PgBouncer mutlak dibutuhkan saat aplikasi web memiliki ratusan hingga ribuan thread koneksi, karena setiap koneksi langsung ke PostgreSQL memakan proses OS terpisah dan sekitar 10MB memori per koneksi.

## Referensi

- [PostgreSQL Official Documentation](https://www.postgresql.org/docs/current/)
- [Use The Index, Luke! - Guide to Database Performance](https://use-the-index-luke.com/)
- [PostgreSQL Wiki on VACUUM and Bloat](https://wiki.postgresql.org/wiki/Vacuum_execution_metadata)
