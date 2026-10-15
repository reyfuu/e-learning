# Materi: MongoDB dan Arsitektur Basis Data Berorientasi Dokumen

**Format:** teks · **Level:** Basic–Advanced · **Estimasi:** 60–90 menit

## Tujuan belajar

Setelah materi ini, pelajar dapat merancang dan mengelola basis data NoSQL menggunakan MongoDB (v6/v7): memahami pemodelan dokumen BSON (Embedding vs Referencing), melakukan operasi CRUD tingkat lanjut, menguasai Aggregation Pipeline untuk analitik dan transformasi data, mengoptimalkan indeks (Compound, Multikey, Text), serta memahami arsitektur Replica Sets dan Sharding.

## Prasyarat

- Memahami format data JSON dan manipulasi struktur objek/array.
- Memiliki MongoDB Community Server atau akun MongoDB Atlas.
- Alat antarmuka query (`mongosh` CLI atau MongoDB Compass) terpasang.

## 1. Dasar: model data dokumen NoSQL dan BSON

MongoDB adalah basis data non-relasional (NoSQL) berorientasi dokumen. Alih-alih baris dan kolom yang kaku, data disimpan dalam bentuk koleksi (**Collections**) dari dokumen-dokumen fleksibel berformat **BSON (Binary JSON)**.

### Perbandingan Konsep SQL vs MongoDB:
| SQL (Relasional) | MongoDB (NoSQL) |
|---|---|
| Database | Database |
| Table | Collection |
| Row / Record | Document (BSON) |
| Column | Field |
| Primary Key (`id`) | `_id` (ObjectId 12-byte) |
| Foreign Key / JOIN | Embedding / `$lookup` |

### Embedding vs Referencing:
- **Embedding (Denormalisasi):** Menyimpan data anak langsung di dalam dokumen induk. Ideal untuk relasi "1-ke-sedikit" dan data yang selalu dibaca bersamaan (contoh: daftar modul di dalam kursus).
- **Referencing (Normalisasi):** Menyimpan ID referensi ke dokumen di koleksi lain. Ideal untuk relasi "1-ke-banyak besar" atau relasi "banyak-ke-banyak" yang sering diperbarui secara independen (contoh: log audit siswa).

## 2. Praktik dasar: Operasi CRUD dan Document Schema

```javascript
// mongosh (MongoDB Shell)

// 1. Menyisipkan Dokumen (Create)
db.courses.insertOne({
  title: "Arsitektur NoSQL dengan MongoDB",
  slug: "arsitektur-nosql-mongodb",
  price: 350000,
  instructor: {
    name: "Budi Santoso",
    email: "budi@learnhub.id"
  },
  tags: ["nosql", "mongodb", "backend", "database"],
  modules: [
    { title: "Pengenalan BSON", duration_mins: 15, is_free: true },
    { title: "Aggregation Pipeline", duration_mins: 45, is_free: false }
  ],
  is_published: true,
  enrollment_count: 120,
  created_at: new Date()
});

// 2. Query Pencarian (Read)
db.courses.find(
  {
    is_published: true,
    price: { $lte: 400000 },
    tags: "backend" // Pencarian di dalam array otomatis didukung
  },
  { title: 1, price: 1, "instructor.name": 1 } // Proyeksi kolom
).sort({ price: 1 }).limit(10);

// 3. Pembaruan Dokumen Atomik (Update)
db.courses.updateOne(
  { slug: "arsitektur-nosql-mongodb" },
  {
    $inc: { enrollment_count: 1 }, // Operasi increment numerik secara atomik
    $push: { tags: "high-demand" }  // Menambahkan elemen ke dalam array
  }
);
```

## 3. Fitur utama dan pola arsitektur

### Aggregation Pipeline
Aggregation Pipeline memproses dokumen melalui tahapan-tahapan berurutan (*stages*) untuk menghasilkan data rangkuman yang telah ditransformasi:

```javascript
db.courses.aggregate([
  // Tahap 1: Filter kursus yang aktif
  { $match: { is_published: true } },

  // Tahap 2: Pecah array tags menjadi dokumen terpisah per elemen
  { $unwind: "$tags" },

  // Tahap 3: Kelompokkan per tag dan hitung rata-rata harga serta total kursus
  {
    $group: {
      _id: "$tags",
      total_courses: { $sum: 1 },
      avg_price: { $avg: "$price" },
      total_students: { $sum: "$enrollment_count" }
    }
  },

  // Tahap 4: Filter hanya kategori yang memiliki lebih dari 2 kursus
  { $match: { total_courses: { $gt: 2 } } },

  // Tahap 5: Urutkan berdasarkan total pelajar terbanyak
  { $sort: { total_students: -1 } }
]);
```

### Strategi Indexing di MongoDB
Tanpa indeks, MongoDB melakukan *COLLSCAN* (Collection Scan) yang memindai setiap dokumen di disk:
- **Single Field Index:** `db.courses.createIndex({ slug: 1 }, { unique: true })`
- **Compound Index:** `db.courses.createIndex({ is_published: 1, created_at: -1 })`
- **Multikey Index:** Dibuat otomatis saat kolom yang diindeks berupa array (misal `tags`).
- **Text Index:** Mendukung pencarian teks kata kunci (`$text: { $search: "database" }`).

## 4. Praktik kerja: optimasi kueri dan diagnosa COLLSCAN dengan `explain()`

### Skenario
Endpoint katalog produk pada aplikasi LearnHub memakan waktu 3,5 detik saat koleksi dokumen mencapai 500.000 data. Server CPU MongoDB melonjak hingga 100%.

### Tugas
Gunakan `.explain("executionStats")` untuk mendeteksi `COLLSCAN`, lalu buat Compound Index dengan aturan **ESR (Equality, Sort, Range)**.

```javascript
// Analisis rencana eksekusi
db.courses.find({
  is_published: true,
  price: { $gte: 100000, $lte: 500000 }
}).sort({ created_at: -1 }).explain("executionStats");
```

Terapkan Aturan ESR Indexing:
1. **Equality:** Letakkan field filter nilai pasti di awal indeks (`is_published`).
2. **Sort:** Letakkan field yang di-sort di bagian tengah (`created_at`).
3. **Range:** Letakkan filter rentang di posisi paling akhir (`price`).

```javascript
// Buat compound index mengikuti aturan ESR:
db.courses.createIndex({
  is_published: 1,
  created_at: -1,
  price: 1
});
```

### Lima hint bertahap
1. Periksa atribut `stage` di `executionStages`; jika bernilai `COLLSCAN`, berarti query memindai seluruh koleksi.
2. Indeks yang baik menghasilkan tahap `IXSCAN` yang diikuti oleh `FETCH`.
3. Bandingkan nilai `totalDocsExamined` dengan `nReturned`; angka idealnya adalah mendekati 1:1.
4. Aturan ESR (Equality, Sort, Range) menjamin MongoDB tidak perlu melakukan in-memory sorting (*blocking sort*).
5. Hindari membuat terlalu banyak indeks pada satu koleksi karena akan memperlambat operasi `INSERT` dan memakan RAM.

## 5. Lanjutan: Replica Sets, Sharding, dan Transaksi Multi-Dokumen

- **Replica Sets:** Klaster node (Primary-Secondary-Arbiter) yang menyediakan failover otomatis dan replikasi data asynchronous.
- **Multi-Document ACID Transactions:** Sejak MongoDB 4.0, MongoDB mendukung transaksi ACID multi-dokumen multi-koleksi dengan session (`session.startTransaction()`).
- **Sharding:** Partisi horizontal dataset besar ke beberapa server/shards menggunakan Shard Key untuk penskalaan throughput baca/tulis tak terbatas.

## Pemeriksaan pemahaman

1. Apa keuntungan dan kelemahan menyimpan data relasi anak secara Embedding di MongoDB?
2. Bagaimana aturan "ESR" membantu dalam perancangan Compound Index?
3. Apa perbedaan tahapan `$match` dan `$filter` dalam Aggregation Pipeline?
4. Kapan kita membutuhkan Multi-Document ACID Transactions di MongoDB?

**Kunci ringkas:** Embedding meningkatkan kecepatan pembacaan karena data diambil dalam satu operasi baca tanpa join, namun memiliki batasan ukuran dokumen (maksimal 16 MB) dan risiko redundansi data saat update; Aturan ESR (Equality, Sort, Range) memastikan indeks digunakan secara optimal baik untuk filter kesetaraan, pengurutan tanpa filesort, maupun penyaringan rentang; `$match` memfilter dokumen di tingkat pipeline sedangkan `$filter` menyaring elemen di dalam array dokumen; Transaksi multi-dokumen dibutuhkan saat perubahan state pada beberapa dokumen berbeda harus bersifat all-or-nothing (misal transfer saldo antar-akun).

## Referensi

- [MongoDB Official Documentation](https://www.mongodb.com/docs/)
- [MongoDB University Free Courses](https://learn.mongodb.com/)
- [MongoDB Indexing Strategies & ESR Rule](https://www.mongodb.com/docs/manual/tutorial/equality-sort-range-rule/)
