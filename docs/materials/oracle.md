# Materi: Oracle Database dan Rekayasa Prosedural PL/SQL Enterprise

**Format:** teks · **Level:** Basic–Advanced · **Estimasi:** 60–90 menit

## Tujuan belajar

Setelah materi ini, pelajar dapat memahami dan mengoperasikan Oracle Database (19c / 23c): memahami arsitektur memori instance (SGA/PGA) dan struktur penyimpanan fisik (Tablespaces, Datafiles), menulis prosedur dan paket PL/SQL yang modular, memanfaatkan Sequences dan Triggers, mengoptimalkan query dengan Oracle Execution Plans (`AUTOTRACE` / `EXPLAIN PLAN`), serta memahami konsep dasar Real Application Clusters (RAC) dan Data Guard.

## Prasyarat

- Memahami konsep dasar relational database dan sintaks standar SQL ANSI.
- Memiliki akses ke Oracle Database (Express Edition/Free 23c) atau Docker image Oracle.
- SQL*Plus atau Oracle SQL Developer terpasang di sistem.

## 1. Dasar: arsitektur Oracle Database dan Instance

Oracle Database memisahkan antara **Database** (file fisik di disk: Datafiles, Redo Log Files, Control Files) dan **Instance** (proses latar belakang OS dan struktur memori RAM).

### Komponen Memori & Proses Oracle:
- **SGA (System Global Area):** Area memori bersama yang diakses oleh semua proses server. Komponen utamanya mencakup *Database Buffer Cache*, *Shared Pool* (caching SQL execution plan dan data dictionary), dan *Redo Log Buffer*.
- **PGA (Program Global Area):** Area memori privat yang dialokasikan khusus untuk setiap koneksi proses server (untuk operasi sorting dan hash join).
- **Background Processes:** Termasuk DBWn (Database Writer), LGWR (Log Writer), CKPT (Checkpoint), dan SMON (System Monitor).

## 2. Praktik dasar: Skema Relasional, Sequence, dan Paket PL/SQL

```sql
-- 1. Pembuatan Sequence dan Tabel
CREATE SEQUENCE seq_courses_id START WITH 1000 INCREMENT BY 1 NOCACHE;

CREATE TABLE courses (
    course_id NUMBER PRIMARY KEY,
    title VARCHAR2(200) NOT NULL,
    price NUMBER(10, 2) DEFAULT 0.00 NOT NULL,
    status VARCHAR2(20) DEFAULT 'DRAFT' NOT NULL,
    created_at TIMESTAMP DEFAULT SYSTIMESTAMP NOT NULL,
    CONSTRAINT chk_course_status CHECK (status IN ('DRAFT', 'PUBLISHED', 'ARCHIVED'))
);

-- 2. Paket PL/SQL (Package Specification)
CREATE OR REPLACE PACKAGE course_pkg AS
    -- Deklarasi prosedur publik
    PROCEDURE create_course(
        p_title IN VARCHAR2,
        p_price IN NUMBER,
        p_course_id OUT NUMBER
    );

    FUNCTION get_course_price(p_course_id IN NUMBER) RETURN NUMBER;
END course_pkg;
/

-- 3. Paket PL/SQL (Package Body)
CREATE OR REPLACE PACKAGE BODY course_pkg AS
    PROCEDURE create_course(
        p_title IN VARCHAR2,
        p_price IN NUMBER,
        p_course_id OUT NUMBER
    ) IS
    BEGIN
        IF p_price < 0 THEN
            RAISE_APPLICATION_ERROR(-20001, 'Harga kursus tidak boleh negatif.');
        END IF;

        p_course_id := seq_courses_id.NEXTVAL;

        INSERT INTO courses (course_id, title, price, status, created_at)
        VALUES (p_course_id, p_title, p_price, 'DRAFT', SYSTIMESTAMP);

        COMMIT;
    EXCEPTION
        WHEN OTHERS THEN
            ROLLBACK;
            RAISE;
    END create_course;

    FUNCTION get_course_price(p_course_id IN NUMBER) RETURN NUMBER IS
        v_price NUMBER(10, 2);
    BEGIN
        SELECT price INTO v_price FROM courses WHERE course_id = p_course_id;
        RETURN v_price;
    EXCEPTION
        WHEN NO_DATA_FOUND THEN
            RETURN 0.00;
    END get_course_price;
END course_pkg;
/
```

Eksekusi PL/SQL via anonymous block:
```sql
DECLARE
    v_new_id NUMBER;
BEGIN
    course_pkg.create_course('Arsitektur Oracle Database 23c', 500000, v_new_id);
    DBMS_OUTPUT.PUT_LINE('Kursus baru berhasil dibuat dengan ID: ' || v_new_id);
END;
/
```

## 3. Fitur utama dan pola arsitektur

### Optimasi Query dan Execution Plan dengan `EXPLAIN PLAN`
Menganalisis keputusan Cost-Based Optimizer (CBO) Oracle:

```sql
EXPLAIN PLAN FOR
SELECT title, price FROM courses WHERE status = 'PUBLISHED' ORDER BY created_at DESC;

-- Lihat tabel rencana eksekusi:
SELECT * FROM TABLE(DBMS_XPLAN.DISPLAY);
```

Indikator penting:
- **TABLE ACCESS FULL:** Memindai seluruh segmen tabel; pastikan ada indeks jika tabel besar.
- **INDEX UNIQUE SCAN vs INDEX RANGE SCAN:** Akses indeks Primary Key/Unique vs pencarian non-unique.
- **Cost & Cardinality:** Estimasi biaya komputasi I/O dan jumlah baris data yang dievaluasi optimizer.

### Transaksi dan Locking di Oracle
Oracle tidak pernah mengunci data saat operasi pembacaan (`SELECT` tidak pernah mengunci data, dan pembacaan tidak pernah memblokir penulisan). Oracle mengandalkan segmen **Undo Tablespace** untuk menyajikan gambaran data masa lalu (*Consistent Read / ORA-01555 Snapshot too old*).

## 4. Praktik kerja: menyelesaikan error `ORA-01555: snapshot too old`

### Skenario
Sebuah batch job pelaporan data keuangan berjalan selama 45 menit membaca tabel transaksi bernilai jutaan baris. Di tengah proses, program gagal dengan pesan error `ORA-01555: snapshot too old: rollback segment number with name "_SYSSMU..." too small`.

### Tugas
Pahami penyebab ORA-01555 (data undo yang dibutuhkan query pembacaan tertimpa oleh transaksi update lain) dan konfigurasikan parameter retensi Undo Tablespace yang memadai.

```sql
-- 1. Periksa ukuran dan retensi Undo saat ini
SHOW PARAMETER UNDO;

-- 2. Tingkatkan periode retensi undo (dalam detik, misal 3600 detik = 1 jam)
ALTER SYSTEM SET UNDO_RETENTION = 3600 SCOPE=BOTH;

-- 3. Pastikan Undo Tablespace dikonfigurasi dengan RETENTION GUARANTEE
ALTER TABLESPACE undotbs1 RETENTION GUARANTEE;
```

### Lima hint bertahap
1. Oracle menggunakan Undo Tablespace untuk merekonstruksi data versi lama saat query panjang dimulai.
2. Jika transaksi konkuren sering melakukan `COMMIT`, blok data di Undo tablespace dapat ditimpa jika retensi habis atau ruang penuh.
3. Menyetel `RETENTION GUARANTEE` memaksa Oracle mempertahankan undo data meskipun tablespace penuh (transaksi update baru akan dibatalkan daripada menimpa undo yang masih dibutuhkan query aktif).
4. Hindari pemanggilan `COMMIT` di dalam loop cursor (*commit within a loop*); lakukan commit per batch besar.
5. Jalankan query analitik berat di waktu senggang (*off-peak hours*) atau di instance read-only Active Data Guard.

## 5. Lanjutan: RAC (Real Application Clusters) dan Data Guard

- **Oracle RAC:** Memungkinkan beberapa server independen (instance) mengakses satu set database fisik yang sama secara simultan melalui shared disk dan interconnect Cache Fusion (*High Availability & Horizontal Scale*).
- **Oracle Data Guard:** Menyediakan replikasi fisik standby otomatis untuk disaster recovery (DR) antar-data center.
- **Automatic Workload Repository (AWR):** Laporan analitik komprehensif performa database yang mencatat Top SQL, wait events, dan bottleneck I/O.

## Pemeriksaan pemahaman

1. Apa perbedaan antara Oracle Database dan Oracle Instance?
2. Mengapa query SELECT di Oracle tidak pernah memblokir operasi UPDATE pada baris data yang sama?
3. Apa penyebab utama munculnya pesan error `ORA-01555: snapshot too old`?
4. Apa keunggulan mengelompokkan prosedur ke dalam PL/SQL Package dibanding membuat prosedur mandiri terpisah?

**Kunci ringkas:** Instance adalah alokasi RAM (SGA) dan background processes di memori sedangkan Database adalah kumpulan file fisik di storage disk; Oracle menggunakan Undo segments untuk menyajikan read consistency snapshot tanpa memerlukan shared read lock; ORA-01555 terjadi ketika query jangka panjang membutuhkan gambar data lama dari Undo segment namun blok tersebut sudah terlanjur ditimpa oleh transaksi lain; Package membungkus fungsionalitas terkait dalam satu namespace, memisahkan spesifikasi dari implementasi body, dan seluruh paket dimuat ke shared pool memori secara utuh dalam satu panggilan.

## Referensi

- [Oracle Database Documentation (23c / 19c)](https://docs.oracle.com/en/database/oracle/oracle-database/)
- [Oracle Base - DBA and Developer Guides](https://oracle-base.com/)
- [Tom Kyte's Expert Oracle Database Architecture](https://www.apress.com/gp/book/9781430262985)
