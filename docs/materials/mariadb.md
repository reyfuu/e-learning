# Materi: MariaDB Server dan Galera Cluster untuk High Availability

**Format:** teks · **Level:** Basic–Advanced · **Estimasi:** 60–90 menit

## Tujuan belajar

Setelah materi ini, pelajar dapat merancang dan mengelola basis data relasional dengan MariaDB (v10.11 / v11): memahami asal-usul dan perbedaan fitur dengan MySQL, menguasai beragam storage engine alternatif (Aria, ColumnStore, InnoDB), memanfaatkan fitur temporal tables dan sequence bawaan, serta memahami konsep replikasi multi-master menggunakan Galera Cluster.

## Prasyarat

- Memahami konsep dasar relational database dan sintaks SQL standar.
- Memiliki akses ke MariaDB Server lokal atau kontainer Docker.
- Memahami konsep dasar High Availability (HA) dan arsitektur database.

## 1. Dasar: MariaDB sebagai evolusi mandiri open-source

MariaDB dibuat oleh pembuat asli MySQL (Michael "Monty" Widenius) sebagai alternatif yang sepenuhnya open-source setelah akuisisi MySQL oleh Oracle. MariaDB mempertahankan kompatibilitas protokol dan port (3306), namun mengembangkan optimizer kueri mandiri, engine penyimpanan baru, dan fitur-fitur inovatif.

### Fitur Khas MariaDB:
- **Storage Engine Ekstensif:** Selain InnoDB, MariaDB menyertakan **Aria** (pengganti MyISAM yang tahan crash), **ColumnStore** (untuk analitik OLAP dan big data), dan **Spider** (sharding tabel transparan).
- **Galera Cluster:** Replikasi synchronous multi-master multi-node bawaan tanpa dependensi pihak ketiga.
- **System-Versioned Tables:** Pelacakan riwayat perubahan data otomatis berdasarkan timestamp masa lalu (*temporal tables*).
- **Native Sequences:** Objek sequence SQL standar (`CREATE SEQUENCE ...`) tanpa harus bergantung pada `AUTO_INCREMENT` kolom.

## 2. Praktik dasar: System-Versioned Tables dan Window Functions

```sql
-- 1. Membuat tabel dengan pelacakan riwayat otomatis (System-Versioned Table)
CREATE TABLE course_pricing (
    course_id INT NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    updated_by VARCHAR(50) NOT NULL,
    PRIMARY KEY (course_id)
) WITH SYSTEM VERSIONING;

-- Sisipkan data awal
INSERT INTO course_pricing (course_id, price, updated_by)
VALUES (101, 200000.00, 'admin_budi');

-- Ubah harga kursus (MariaDB otomatis mengarsipkan data lama ke history partition di disk)
UPDATE course_pricing
SET price = 250000.00, updated_by = 'admin_siti'
WHERE course_id = 101;

-- 2. Query melihat data pada waktu masa lalu (Time Travel Query)
SELECT * FROM course_pricing; -- Menampilkan harga aktif saat ini (250000)

SELECT * FROM course_pricing 
FOR SYSTEM_TIME AS OF TIMESTAMP'2026-10-15 10:00:00'
WHERE course_id = 101; -- Mengembalikan harga lama 200000 secara otomatis!

-- Query melihat riwayat lengkap audit perubahan
SELECT course_id, price, updated_by, row_start, row_end
FROM course_pricing FOR SYSTEM_TIME ALL
WHERE course_id = 101;
```

## 3. Fitur utama dan pola arsitektur

### Storage Engine Aria
Aria adalah engine default untuk tabel sistem internal dan tabel sementara (*temporary tables*) di MariaDB:
- Tahan terhadap crash mendadak (*crash-safe*) karena menggunakan transaksi berbasis log.
- Memiliki throughput baca yang sangat tinggi untuk data read-heavy dibandingkan MyISAM lama.

### Galera Cluster (Multi-Master Synchronous Replication)
Berbeda dengan replikasi Master-Replica standar MySQL yang bersifat *asynchronous*:
- Setiap node di dalam klaster adalah Master yang dapat menerima operasi penulisan (*write*).
- Menggunakan protokol *Certification-based Replication* berbasis kuorum.
- Tidak ada data loss (*zero slave lag*) saat satu node server mati mendadak.

## 4. Praktik kerja: migrasi data dan penanganan komparabilitas MySQL vs MariaDB

### Skenario
Tim pengembang memindahkan dump database dari MySQL 8 ke MariaDB 10.11. Terjadi error saat mengimpor skema karena fungsi JSON dan beberapa variabel sistem tidak kompatibel.

### Tugas
Konfigurasikan mode kompatibilitas SQL MariaDB dan pahami pemetaan tipe data `JSON` (yang di MariaDB merupakan alias untuk `LONGTEXT` dengan constraint `JSON_VALID()`).

```sql
-- Di MariaDB, JSON adalah teks dengan validasi check:
CREATE TABLE course_reviews (
    id INT AUTO_INCREMENT PRIMARY KEY,
    course_id INT NOT NULL,
    rating INT CHECK (rating BETWEEN 1 AND 5),
    details JSON, -- Alias untuk LONGTEXT CHECK (JSON_VALID(details))
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Penggunaan fungsi JSON teroptimasi:
INSERT INTO course_reviews (course_id, rating, details)
VALUES (101, 5, '{"pros": "Materi sangat jelas", "recommended": true}');

-- Ekstraksi field JSON:
SELECT id, rating, JSON_VALUE(details, '$.pros') AS review_text
FROM course_reviews
WHERE JSON_VALUE(details, '$.recommended') = 'true';
```

### Lima hint bertahap
1. Periksa variabel `sql_mode` di `my.cnf` atau session untuk menyelaraskan aturan sintaks.
2. Di MySQL 8, JSON disimpan dalam format biner internal; di MariaDB, JSON disimpan sebagai string teks dengan indeks fungsional atau virtual columns.
3. Buat *Virtual Column* yang diindeks jika sering melakukan filter berdasarkan properti di dalam JSON.
4. Gunakan perintah `mariadb-dump` alih-alih `mysqldump` untuk memastikan fitur-fitur spesifik MariaDB (seperti Sequence dan System Versioning) tercadangkan dengan benar.
5. Pada Galera Cluster, setiap tabel **wajib** memiliki Primary Key.

## 5. Lanjutan: Galera Arbitrator, Optimizer Cost, dan Monitoring

- **Galera Arbitrator (garbd):** Node ringan tanpa data lokal yang bertindak sebagai pemutus suara (*tie-breaker quorum*) pada klaster 2-node untuk mencegah kondisi *split-brain*.
- **MariaDB Thread Pool:** Fitur bawaan gratis di versi Community MariaDB (berbeda dengan MySQL yang mewajibkan lisensi Enterprise) untuk menangani hingga puluhan ribu koneksi simultan.
- **Performance Schema & User Statistics:** Aktifkan `userstat=1` untuk memantau konsumsi CPU dan aktivitas I/O per user.

## Pemeriksaan pemahaman

1. Apa itu System-Versioned Tables di MariaDB dan apa manfaatnya bagi sistem audit?
2. Apa perbedaan mendasar antara replikasi standar asynchronous dengan Galera Cluster?
3. Mengapa Storage Engine Aria lebih unggul dibanding MyISAM?
4. Mengapa pada Galera Cluster setiap tabel wajib memiliki Primary Key?

**Kunci ringkas:** System-Versioned Tables secara otomatis menyimpan dan mengarsipkan seluruh riwayat perubahan record tanpa perlu membuat tabel audit manual atau trigger database; Galera Cluster mereplikasi data secara synchronous multi-master sehingga seluruh node memiliki data identik seketika tanpa jeda replikasi (*replication lag*); Aria memiliki kemampuan crash recovery otomatis menggunakan write-ahead logging sedangkan MyISAM mudah korup jika server mati mendadak; Galera membutuhkan Primary Key untuk melacak baris data yang diubah saat memvalidasi sertifikasi konflik transaksi (*row certification*) antar-node klaster.

## Referensi

- [MariaDB Official Knowledge Base](https://mariadb.com/kb/en/)
- [Galera Cluster Documentation](https://galeracluster.com/library/)
- [MariaDB System-Versioned Tables Guide](https://mariadb.com/kb/en/system-versioned-tables/)
