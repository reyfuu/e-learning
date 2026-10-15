# Materi: MySQL 8.x dan Engine Transaksional InnoDB

**Format:** teks · **Level:** Basic–Advanced · **Estimasi:** 60–90 menit

## Tujuan belajar

Setelah materi ini, pelajar dapat merancang, mengoptimalkan, dan mengelola database MySQL 8.x: memahami arsitektur storage engine InnoDB (Clustered Index, Buffer Pool, Redo Log, Undo Log), merancang skema relasional ACID yang efisien, mengoptimalkan query lambat menggunakan `EXPLAIN ANALYZE`, mengelola penguncian (*Locking & Deadlocks*), serta memahami replikasi Master-Replica dan strategi backup.

## Prasyarat

- Memahami konsep dasar relational database dan sintaks dasar SQL.
- Memiliki MySQL Server 8.0+ terpasang di sistem atau menggunakan Docker.
- Memahami konsep dasar transaksi dan integritas data (ACID).

## 1. Dasar: arsitektur MySQL dan Storage Engine InnoDB

MySQL menggunakan arsitektur modular yang memisahkan lapisan SQL parser/optimizer dari **Storage Engine** yang menangani penyimpanan fisik ke disk. Di MySQL 8, **InnoDB** adalah storage engine default yang menyediakan kepatuhan penuh terhadap prinsip ACID.

### Konsep Kunci InnoDB:
- **Clustered Index (Primary Key):** Baris data fisik tabel disimpan secara berurutan langsung di dalam struktur B+Tree Primary Key. Membaca data via Primary Key adalah operasi paling cepat.
- **Secondary Index:** Indeks non-primary menyimpan nilai kolom yang diindeks beserta referensi ke nilai Primary Key-nya (memerlukan *bookmark lookup*).
- **InnoDB Buffer Pool:** Area memori RAM terbesar tempat InnoDB melakukan caching data tabel dan indeks.
- **Redo Log (WAL) & Undo Log:** Redo log menjamin *Durability* saat crash mendadak, sedangkan Undo log menyediakan *Rollback* dan MVCC snapshot isolation.

## 2. Praktik dasar: Skema Relasional, Foreign Key, dan Transaction

```sql
-- 1. Skema Tabel dengan InnoDB
CREATE TABLE instructors (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(150) NOT NULL UNIQUE,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE courses (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    instructor_id INT UNSIGNED NOT NULL,
    title VARCHAR(200) NOT NULL,
    price DECIMAL(10, 2) NOT NULL DEFAULT 0.00,
    status ENUM('draft', 'published', 'archived') NOT NULL DEFAULT 'draft',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_courses_instructor 
        FOREIGN KEY (instructor_id) REFERENCES instructors(id) 
        ON DELETE RESTRICT ON UPDATE CASCADE,
    INDEX idx_instructor_status (instructor_id, status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 2. Transaksi Aman ACID
START TRANSACTION;

INSERT INTO instructors (name, email) VALUES ('Ahmad Dahlan', 'ahmad@learnhub.id');
SET @instructor_id = LAST_INSERT_ID();

INSERT INTO courses (instructor_id, title, price, status)
VALUES (@instructor_id, 'Arsitektur MySQL 8 untuk Skala Besar', 300000.00, 'published');

COMMIT;
```

## 3. Fitur utama dan pola arsitektur

### Optimasi Query dengan `EXPLAIN ANALYZE` (MySQL 8.0+)
MySQL 8 menyediakan `EXPLAIN ANALYZE` yang mengeksekusi query dan menampilkan waktu eksekusi riil, iterator tree, dan jumlah baris data yang dibaca:

```sql
EXPLAIN ANALYZE
SELECT c.title, c.price, i.name AS instructor_name
FROM courses c
INNER JOIN instructors i ON c.instructor_id = i.id
WHERE c.status = 'published' AND c.price > 100000;
```

Indikator penting:
- **`type: ALL` (Table Scan):** MySQL membaca seluruh tabel; membutuhkan penambahan indeks.
- **`type: ref` atau `eq_ref`:** MySQL membaca data dengan indeks yang sangat efisien.
- **`Using filesort`:** MySQL harus melakukan sorting manual di memori/disk tambahan; optimasi dengan indeks komposit yang menyertakan urutan kolom `ORDER BY`.

### Covering Index
Jika semua kolom yang diminta dalam `SELECT` sudah ada di dalam Secondary Index, MySQL tidak perlu melakukan *bookmark lookup* ke Clustered Index tabel utama:

```sql
-- Indeks komposit covering:
CREATE INDEX idx_courses_covering ON courses (status, price, title);

-- Query ini 100% dipenuhi oleh indeks tanpa membaca tabel fisik (Using index)
SELECT title, price FROM courses WHERE status = 'published';
```

## 4. Praktik kerja: menganalisis dan menyelesaikan Deadlock

### Skenario
Dua proses transaksi pembayaran berjalan konkuren. Transaksi A memperbarui tabel pendaftaran lalu mengunci saldo dompet pengguna, sedangkan Transaksi B memperbarui saldo dompet pengguna lalu mencoba mengunci tabel pendaftaran. Terjadi error: `ERROR 1213 (40001): Deadlock found when trying to get lock; try restarting transaction`.

### Tugas
Analisis penyebab deadlock dari engine status, dan ubah urutan penguncian tabel di aplikasi agar konsisten.

```sql
-- Periksa rincian transaksi penyebab deadlock terakhir:
SHOW ENGINE INNODB STATUS\G
```

Solusi arsitektural:
1. **Urutan Penguncian Konsisten:** Selalu kunci tabel dalam urutan yang sama di seluruh kode aplikasi (misalnya selalu kunci `wallets` terlebih dahulu, baru kemudian `enrollments`).
2. **Kueri Singkat:** Jangan menaruh jeda logika bisnis panjang di antara `START TRANSACTION` dan `COMMIT`.
3. **Pemberian Indeks Tepat:** Pastikan klausa `WHERE` pada query update menggunakan kolom berindeks agar InnoDB hanya mengunci baris (*Row-level lock*) bukan mengunci rentang (*Gap lock* / *Table lock*).

```sql
-- Transaksi A dan Transaksi B sama-sama mengunci wallet terlebih dahulu:
START TRANSACTION;
SELECT * FROM wallets WHERE user_id = 42 FOR UPDATE;
-- ... update wallets ...
SELECT * FROM enrollments WHERE id = 101 FOR UPDATE;
-- ... update enrollments ...
COMMIT;
```

### Lima hint bertahap
1. Deadlock terjadi saat dua transaksi saling menunggu kunci yang dipegang satu sama lain.
2. InnoDB secara otomatis mendeteksi deadlock dan membatalkan (*rollback*) transaksi yang memiliki jumlah modifikasi baris paling sedikit.
3. Aplikasi backend harus selalu memiliki logika *retry* otomatis saat menangani error kode 1213.
4. Jangan pernah menggunakan transaksi besar yang menyentuh ribuan baris data dalam satu waktu.
5. Perhatikan bahwa update pada kolom tanpa indeks akan mengunci seluruh tabel (*full table lock*).

## 5. Lanjutan: Replikasi GTID, Binlog, dan Tuning Buffer Pool

- **InnoDB Buffer Pool Size:** Pada server database dedicated, setel `innodb_buffer_pool_size` sekitar 60–75% dari total kapasitas RAM server.
- **GTID Replication (Global Transaction Identifiers):** Mempermudah pelacakan transaksi antar Master dan Replica tanpa harus mencatat file offset binlog manual.
- **Percona XtraBackup:** Alat backup fisik online open-source tanpa downtime dan non-blocking untuk basis data MySQL terabytes.

## Pemeriksaan pemahaman

1. Apa perbedaan arsitektural antara Clustered Index dan Secondary Index pada InnoDB?
2. Mengapa urutan penguncian data yang seragam mencegah terjadinya Deadlock?
3. Apa fungsi komponen Redo Log pada database MySQL?
4. Apa yang dimaksud dengan "Covering Index" dan mengapa performanya sangat tinggi?

**Kunci ringkas:** Clustered index menyimpan seluruh baris data fisik di daun B+Tree primary key, sedangkan secondary index hanya menyimpan kolom indeks dan nilai primary key; Urutan penguncian yang konsisten menjamin dependensi sirkular antar transaksi tidak pernah terbentuk sehingga deadlock dihindari; Redo Log mencatat transaksi yang belum sempat di-flush dari RAM Buffer Pool ke disk agar data tidak hilang saat server mati mendadak; Covering index memuat semua kolom yang dibutuhkan oleh query SELECT sehingga MySQL tidak perlu melakukan secondary lookup ke tabel utama.

## Referensi

- [MySQL 8.0 Reference Manual](https://dev.mysql.com/doc/refman/8.0/en/)
- [High Performance MySQL by Baron Schwartz et al.](https://www.oreilly.com/library/view/high-performance-mysql/9781492080510/)
- [MySQL InnoDB Locking Explanation](https://dev.mysql.com/doc/refman/8.0/en/innodb-locking.html)
