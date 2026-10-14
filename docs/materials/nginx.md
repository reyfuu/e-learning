# Materi: Nginx untuk Reverse Proxy dan Deployment Web

**Format:** teks · **Level:** Basic–Advanced · **Estimasi:** 60–90 menit

## Tujuan belajar

Setelah materi ini, pelajar dapat menjelaskan alur reverse proxy, mengarahkan trafik HTTP ke aplikasi, menyajikan file statis, menguji konfigurasi sebelum reload, dan mendiagnosis masalah umum seperti `502 Bad Gateway`.

## Prasyarat

- Memahami HTTP, DNS, port, dan alamat IP.
- Memiliki aplikasi web yang berjalan pada port lokal, misalnya `127.0.0.1:8080`.
- Memiliki Nginx di Linux atau memakai container lokal.

## 1. Dasar: apa yang dilakukan Nginx?

Nginx menerima koneksi dari browser lalu memilih cara menangani request berdasarkan host, path, dan aturan konfigurasi. Nginx dapat menyajikan file statis secara langsung atau meneruskan request ke aplikasi di belakangnya.

Pada reverse proxy, browser hanya berkomunikasi dengan Nginx. Nginx meneruskan request ke aplikasi, lalu mengembalikan response aplikasi ke browser. Pola ini memudahkan pengelolaan TLS, routing beberapa aplikasi, batas request, dan log akses di satu lapisan.

Alur sederhana:

```text
Browser → DNS/domain → Nginx :80/:443 → aplikasi :8080
```

Istilah penting:

- **Server block:** aturan untuk satu host atau domain.
- **Location:** aturan yang memilih perilaku berdasarkan path.
- **Upstream:** aplikasi tujuan di balik reverse proxy.
- **Access log:** catatan request yang diterima.
- **Error log:** catatan masalah Nginx dan koneksi ke upstream.

## 2. Praktik dasar: teruskan request ke aplikasi

Contoh server block untuk aplikasi yang mendengarkan pada `127.0.0.1:8080`:

```nginx
server {
    listen 80;
    server_name learnhub.example.com;

    location / {
        proxy_pass http://127.0.0.1:8080;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

`proxy_set_header` meneruskan konteks request yang dibutuhkan aplikasi, misalnya host asli, alamat klien, dan skema HTTP/HTTPS. Jangan menganggap nilai header dari klien tepercaya; aplikasi harus mengandalkan proxy yang dikonfigurasi dan memvalidasi trust proxy sesuai lingkungannya.

Sebelum menerapkan perubahan:

```bash
sudo nginx -t
sudo systemctl reload nginx
```

`nginx -t` memeriksa sintaks dan file yang dirujuk. Jika pemeriksaan gagal, jangan reload. Reload menerapkan konfigurasi tanpa menghentikan worker yang sedang melayani koneksi.

## 3. File statis dan pemilihan path

Nginx juga dapat menyajikan aset statis tanpa melewati aplikasi:

```nginx
location /assets/ {
    alias /srv/learnhub/assets/;
    expires 7d;
    add_header Cache-Control "public";
}
```

Gunakan `root` ketika path request ditambahkan ke direktori dasar, dan `alias` ketika path location diganti dengan direktori tujuan. Perbedaan trailing slash dan path dapat menyebabkan file yang benar justru menghasilkan 404. Uji URL dan path hasilnya sebelum rilis.

## 4. Praktik kerja: diagnosis `502 Bad Gateway`

### Skenario

Setelah deployment, halaman login menampilkan `502 Bad Gateway`. Nginx hidup dan DNS sudah menuju server yang benar. Tim mengonfirmasi aplikasi seharusnya mendengarkan di `127.0.0.1:8080`.

### Tugas

Temukan penyebabnya dan susun langkah perbaikan tanpa membuka port aplikasi ke internet.

### Langkah diagnosis

1. Periksa error log Nginx untuk pesan seperti `connect() failed` atau `upstream timed out`.
2. Pastikan proses aplikasi aktif dan bind pada alamat/port yang dikonfigurasi.
3. Dari server yang sama, uji `curl -i http://127.0.0.1:8080/health`.
4. Bandingkan alamat, port, protokol, dan path pada `proxy_pass` dengan konfigurasi aplikasi.
5. Setelah memperbaiki sumber masalah, jalankan `sudo nginx -t`, lalu reload dan uji request publik.

Jangan mengatasi 502 dengan membuka port aplikasi ke seluruh internet. Pertahankan aplikasi pada loopback atau jaringan privat, lalu biarkan Nginx menjadi pintu masuk publik.

### Lima hint bertahap

1. `502` berarti Nginx menerima request, tetapi belum memperoleh response valid dari upstream.
2. Periksa apakah aplikasi benar-benar aktif; proses Nginx yang sehat tidak menjamin aplikasi ikut berjalan.
3. Jalankan health check dari mesin yang sama dengan Nginx agar DNS publik dan firewall tidak mengaburkan diagnosis.
4. Cocokkan alamat bind aplikasi dengan host di `proxy_pass`; `localhost` pada container berbeda menunjuk ke container itu sendiri.
5. Perbaiki alamat atau jaringan upstream yang salah, pastikan health check berhasil, lalu uji konfigurasi dan reload Nginx.

## 5. Lanjutan: HTTPS, load balancing, dan hardening

- **HTTPS/TLS:** terminasi TLS di Nginx; gunakan sertifikat valid dan redirect HTTP ke HTTPS setelah sertifikat siap.
- **Load balancing:** kelompokkan beberapa backend dalam `upstream`, lalu pantau health dan distribusi koneksi. Load balancing tidak menggantikan retry/idempotency pada aplikasi.
- **Timeout:** sesuaikan `proxy_connect_timeout`, `proxy_read_timeout`, dan batas body dengan kebutuhan endpoint; jangan menaikkannya tanpa memahami dampak antrean dan penggunaan resource.
- **Rate limiting:** batasi endpoint yang rawan penyalahgunaan, terutama login dan API publik, lalu pantau false positive.
- **Header keamanan:** terapkan header di satu lapisan secara konsisten. Hindari duplikasi yang menghasilkan nilai bertentangan dari aplikasi dan proxy.
- **Operasi:** simpan konfigurasi dalam version control, uji sebelum reload, pantau access/error log, dan siapkan rollback ke konfigurasi terakhir yang valid.

## Pemeriksaan pemahaman

1. Apa perbedaan reverse proxy dan menyajikan file statis?
2. Mengapa `nginx -t` dijalankan sebelum reload?
3. Sebutkan dua penyebab umum `502` selain salah sintaks Nginx.
4. Mengapa port aplikasi sebaiknya tidak dibuka langsung ke internet ketika Nginx menjadi reverse proxy?

**Kunci ringkas:** reverse proxy meneruskan request ke aplikasi, sedangkan file statis dilayani langsung oleh Nginx; `nginx -t` mencegah konfigurasi tidak valid diterapkan; 502 sering muncul karena aplikasi mati, bind/port salah, jaringan container keliru, atau upstream timeout; membatasi aplikasi ke loopback mengurangi permukaan akses publik.

## Referensi

- [Nginx Beginner’s Guide](https://nginx.org/en/docs/beginners_guide.html)
- [Nginx `ngx_http_proxy_module`](https://nginx.org/en/docs/http/ngx_http_proxy_module.html)
- [Nginx `ngx_http_upstream_module`](https://nginx.org/en/docs/http/ngx_http_upstream_module.html)
