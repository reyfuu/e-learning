# BRD — LearnHub

**Versi:** 1.3 · **Tanggal:** 15 Oktober 2026 · **Status:** Draft target bisnis; belum seluruhnya diterapkan

## 1. Ringkasan

LearnHub adalah platform e-learning multibahasa untuk pelajar, instruktur, dan admin. Platform menyelenggarakan kursus terstruktur, pembelajaran mandiri, penjualan kursus, pemantauan kemajuan, serta katalog jalur teknologi dari tingkat Basic sampai Advanced. Sasaran awal adalah pasar Indonesia dengan dukungan Bahasa Indonesia dan Inggris.

## 2. Tujuan bisnis

- Membuat materi teknologi terstruktur dan dapat dipelajari bertahap.
- Membantu instruktur menerbitkan, mengelola, dan memonetisasi kursus.
- Menyediakan pengalaman belajar yang dapat diukur melalui progress, penilaian, dan sertifikat.
- Menyediakan fondasi deploy yang mendukung frontend di Vercel dan layanan backend yang sesuai dengan runtime-nya.

Metrik acuan dari [PRD induk](../PRD_eLearning.md) adalah pengguna terdaftar (10.000/50.000), kursus aktif (200/1.000), MAU (5.000/25.000), gross revenue (Rp 50 juta/Rp 500 juta), completion rate (≥30%/≥45%), dan NPS (≥40/≥60%) untuk horizon 6/12 bulan. Angka tersebut masih target perencanaan dari PRD induk, bukan baseline terukur atau komitmen produksi; validasi kapasitas dan persetujuan bisnis diperlukan sebelum dipakai sebagai KPI resmi.

## 3. Pemangku kepentingan dan kebutuhan

| Pemangku kepentingan | Kebutuhan bisnis |
|---|---|
| Pelajar | Menemukan jalur belajar, mengikuti Basic lalu Advanced, menyimpan progres, mengerjakan tugas, memperoleh sertifikat |
| Pelajar/praktisi | Mencoba latihan coding di browser dengan umpan balik dan hint bertahap pada kasus yang menyerupai pekerjaan |
| Instruktur | Menyusun kursus dan lesson, mengunggah materi, melihat peserta dan performa kursus |
| Admin akademik | Meninjau kualitas, mengelola kategori dan katalog teknologi, menerbitkan atau menolak kursus |
| Operasional/keuangan | Memantau transaksi, refund, komisi, dan pencairan instruktur |
| Tim teknologi | Menjaga keamanan, ketersediaan, biaya, observability, dan pipeline rilis |

## 4. Ruang lingkup bisnis

### Termasuk

- Autentikasi, profil, role learner/instructor/admin, dan verifikasi akun.
- Katalog kursus, pencarian, kurikulum, enrollment, lesson, tugas/kuis, ulasan, progress, serta sertifikat.
- Pembuatan dan moderasi kursus instruktur.
- Pembayaran dan riwayat transaksi sebagai fase monetisasi.
- Katalog silabus teknologi Basic–Advanced yang tercantum di bawah.
- Latihan coding interaktif untuk pilot teknologi yang dipilih, dengan soal berbasis skenario kerja dan lima hint bertahap.
- Situs publik dan aplikasi learner/instructor/admin yang responsif.

### Di luar ruang lingkup rilis awal

- Aplikasi native mobile.
- Kelas live/video conference buatan sendiri.
- Menjamin sertifikat memiliki akreditasi eksternal.
- Menyediakan lingkungan IDE/cloud lab interaktif untuk semua bahasa sejak rilis pertama.
- Menyediakan isi lesson berupa buku lengkap; dokumen ini menetapkan kurikulum dan hasil belajar yang harus diproduksi.

## 5. Model operasional dan pendapatan

Platform menyediakan kursus gratis dan berbayar, promosi/kupon, biaya layanan platform, serta pencairan instruktur. Persentase komisi, gateway pembayaran, aturan pajak, refund, serta SLA pencairan merupakan kebijakan bisnis yang perlu disetujui sebelum fitur pembayaran diaktifkan. Kursus berbayar hanya memberi akses setelah status pembayaran tervalidasi melalui webhook terautentikasi.

## 6. Risiko dan asumsi bisnis

| Risiko/asumsi | Dampak | Penanganan |
|---|---|---|
| Katalog sangat luas, kualitas lesson tidak merata | Kepercayaan dan completion turun | Rubrik akademik, review ahli, versi silabus, pilot per rumpun |
| Video dan file menimbulkan biaya bandwidth/storage | Biaya meningkat | Object storage + CDN, batas ukuran, lifecycle policy, streaming adaptif |
| Deploy di Vercel dianggap mencakup Go API | Pilihan runtime/deploy tidak sesuai | Tegaskan Vercel untuk Next.js; Go API/worker di container platform terpisah |
| Pembayaran dan payout memiliki kewajiban operasional | Risiko finansial/kepatuhan | Tahapkan monetisasi, rekonsiliasi, audit log dan kebijakan tertulis |
| Target skala PRD belum memiliki uji kapasitas | Risiko SLA tidak realistis | Uji beban sebelum klaim kapasitas dan tetapkan SLO terukur |

## 7. Tahapan bisnis

1. **Fondasi:** identitas, katalog, course authoring, konten, enrollment gratis, progress.
2. **Pembelajaran:** kuis/tugas, ulasan, sertifikat, dashboard instruktur dan moderasi.
3. **Monetisasi:** checkout, webhook, promo, transaksi, komisi, refund dan payout.
4. **Skala:** analitik, rekomendasi, optimasi biaya, keamanan, serta ekspansi konten.

## 8. Kriteria keberhasilan

- Pelajar dapat menemukan kursus, mendaftar, belajar, dan melihat progres secara konsisten.
- Instruktur dapat mengajukan kursus dan admin dapat melakukan moderasi.
- Silabus setiap teknologi memiliki level Basic dan Advanced yang jelas, prasyarat, hasil belajar, dan proyek praktik.
- Pelajar dapat mencoba latihan coding pilot di browser dan memperoleh hasil evaluasi serta tepat lima hint bertahap.
- Rilis frontend dapat dilakukan melalui Vercel; API, database, worker, dan media memiliki deployment serta pemantauan yang terdokumentasi.
- Metrik completion, enrollment, error, availability, dan biaya dapat diukur sebelum target kuantitatif disahkan.

## 9. Hubungan dokumen dan status

`PRD_eLearning.md` adalah acuan kebutuhan produk LearnHub secara umum. PRD katalog silabus di folder ini memperinci inisiatif katalog teknologi; BRD ini menjelaskan tujuan bisnisnya, FRD menjabarkan perilaku produk, dan TRD menjabarkan target teknis. Dokumen-dokumen ini menetapkan sasaran, bukan bukti bahwa fitur sudah tersedia. Kondisi implementasi yang diperiksa pada 15 Oktober 2026 dicatat di bagian baseline TRD.
