# Materi: Go (Golang) untuk Rekayasa Sistem Berkinerja Tinggi

**Format:** teks · **Level:** Basic–Advanced · **Estimasi:** 60–90 menit

## Tujuan belajar

Setelah materi ini, pelajar dapat merancang program dan backend service berkinerja tinggi menggunakan Go 1.22+: memahami model tipe, struct, dan implicit interfaces; menguasai idiom penanganan error eksplisit; memanfaatkan model konkurensi CSP (*Goroutines* & *Channels*); mengelola *cancellation* dan batas waktu dengan `context.Context`; serta merancang paket HTTP standar yang tangguh.

## Prasyarat

- Memahami konsep dasar pemrograman (variabel, pointer, fungsi).
- Memiliki Go 1.22+ terpasang di sistem (`go version`).
- Memahami konsep dasar thread sistem operasi dan soket jaringan.

## 1. Dasar: filosofi kesederhanaan Go

Go dirancang oleh Google dengan penekanan pada **kesederhanaan, keterbacaan, dan efisiensi kompilasi**. Go tidak memiliki *inheritance* kelas, tidak memiliki *exception try/catch*, dan tidak memiliki overloading fungsi. Sebagai gantinya, Go mengandalkan *composition*, *implicit interfaces*, dan *explicit error values*.

### Fitur Kunci:
- **Goroutine:** Utas ringan (user-space thread) yang dikelola oleh Go runtime Scheduler (M:N scheduler) dengan alokasi stack awal hanya 2 KB.
- **Channel:** Saluran komunikasi ber-tipe untuk berbagi data antar-goroutine (*"Do not communicate by sharing memory; instead, share memory by communicating"*).
- **Interface Implisit:** Tipe apa pun yang mengimplementasikan semua method dalam suatu interface secara otomatis memenuhi interface tersebut tanpa deklarasi `implements`.
- **Error as Values:** Error adalah nilai biasa bertipe interface `error`, ditangani eksplisit dengan idiom `if err != nil`.

## 2. Praktik dasar: Struct, Interface, Error Handling, dan Concurrency

```go
package main

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

// Definisi interface implisit
type CourseNotifier interface {
	NotifyStudent(ctx context.Context, studentEmail, message string) error
}

type EmailService struct {
	Sender string
}

func (e *EmailService) NotifyStudent(ctx context.Context, studentEmail, message string) error {
	select {
	case <-ctx.Done():
		return ctx.Err() // Menangani timeout atau pembatalan
	case <-time.After(100 * time.Millisecond): // Simulasi latency jaringan
		fmt.Printf("[Email sent from %s to %s]: %s\n", e.Sender, studentEmail, message)
		return nil
	}
}

type CourseEnrollment struct {
	StudentID string
	CourseID  string
}

func ProcessEnrollmentsConcurrently(ctx context.Context, notifier CourseNotifier, enrollments []CourseEnrollment) error {
	var wg sync.WaitGroup
	errChan := make(chan error, len(enrollments))

	for _, enrollment := range enrollments {
		wg.Add(1)

		// Jalankan goroutine untuk tiap pemrosesan
		go func(enr CourseEnrollment) {
			defer wg.Done()

			msg := fmt.Sprintf("Anda telah terdaftar pada kursus %s", enr.CourseID)
			if err := notifier.NotifyStudent(ctx, enr.StudentID+"@learnhub.id", msg); err != nil {
				errChan <- fmt.Errorf("gagal notifikasi %s: %w", enr.StudentID, err)
			}
		}(enrollment)
	}

	wg.Wait()
	close(errChan)

	// Kumpulkan error jika ada
	var aggregatedErr error
	for err := range errChan {
		aggregatedErr = errors.Join(aggregatedErr, err)
	}

	return aggregatedErr
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	service := &EmailService{Sender: "no-reply@learnhub.id"}
	items := []CourseEnrollment{
		{StudentID: "std_01", CourseID: "go_101"},
		{StudentID: "std_02", CourseID: "nextjs_201"},
	}

	if err := ProcessEnrollmentsConcurrently(ctx, service, items); err != nil {
		fmt.Printf("Terjadi error: %v\n", err)
		return
	}

	fmt.Println("Semua pendaftaran berhasil diproses secara serentak.")
}
```

## 3. Fitur utama dan pola arsitektur

### Pola Fan-out / Fan-in dan Worker Pool
Mengontrol jumlah goroutine konkuren agar tidak membanjiri resource memori atau koneksi database:

```go
func WorkerPool(jobs <-chan int, results chan<- int, workerCount int) {
	var wg sync.WaitGroup
	for w := 1; w <= workerCount; w++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for job := range jobs {
				// Kerjakan tugas
				results <- job * 2
			}
		}(w)
	}
	wg.Wait()
	close(results)
}
```

### Context Cancellation dan Deadline
Paket `context` adalah standar Go untuk propagasi sinyal pembatalan, batas waktu (*deadline/timeout*), dan metadata lintas batas API:

```go
ctx, cancel := context.WithTimeout(parentCtx, 500*time.Millisecond)
defer cancel()

req, _ := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.external.com/data", nil)
resp, err := http.DefaultClient.Do(req)
```

## 4. Praktik kerja: mendeteksi dan mencegah Data Race dengan Race Detector

### Skenario
Layanan web menghitung jumlah pengunjung aktif secara global menggunakan variabel integer biasa. Saat diuji beban tinggi (*stress test*), nilai counter menjadi kacau dan program kadang crash tak terduga.

### Tugas
Temukan masalah *race condition* menggunakan flag bawaan `-race` dan perbaiki menggunakan `sync/atomic` atau `sync.Mutex`.

```go
package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

type SafeVisitorCounter struct {
	// Solusi 1: Atomic (Sangat efisien untuk nilai primitif)
	atomicCount int64

	// Solusi 2: sync.RWMutex (Untuk struktur data kompleks)
	mu    sync.RWMutex
	stats map[string]int
}

func NewCounter() *SafeVisitorCounter {
	return &SafeVisitorCounter{
		stats: make(map[string]int),
	}
}

func (c *SafeVisitorCounter) Increment() {
	atomic.AddInt64(&c.atomicCount, 1)
}

func (c *SafeVisitorCounter) Count() int64 {
	return atomic.LoadInt64(&c.atomicCount)
}

func (c *SafeVisitorCounter) TrackPage(page string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.stats[page]++
}

func (c *SafeVisitorCounter) GetPageViews(page string) int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.stats[page]
}
```

### Lima hint bertahap
1. Jalankan unit test atau aplikasi dengan flag `go run -race .` atau `go test -race ./...`.
2. Dua atau lebih goroutine yang mengakses memori yang sama di mana minimal satu operasi berupa operasi penulisan (*write*) tanpa sinkronisasi akan memicu race condition.
3. Gunakan paket `sync/atomic` untuk operasi matematika sederhana pada angka.
4. Gunakan `sync.RWMutex` jika pembacaan (*RLock*) jauh lebih sering terjadi dibandingkan penulisan (*Lock*).
5. Jangan pernah menyalin (*pass by value*) objek yang mengandung `sync.Mutex`.

## 5. Lanjutan: Garbage Collector, Profiling pprof, dan Escape Analysis

- **Escape Analysis:** Periksa apakah variabel dialokasikan di stack atau memori heap menggunakan `go build -gcflags="-m"`. Alokasi di stack bebas dari beban Garbage Collector.
- **pprof:** Profiler CPU dan Heap bawaan Go via paket `net/http/pprof` untuk menganalisis fungsi mana yang memakan waktu eksekusi paling besar di lingkungan produksi.
- **Graceful Shutdown:** Tangani sinyal OS `os.Interrupt` atau `syscall.SIGTERM` untuk mematikan `http.Server` tanpa memutus koneksi klien yang aktif.

## Pemeriksaan pemahaman

1. Mengapa Go tidak menyertakan mekanisme Exception (try/catch)?
2. Bagaimana cara kerja M:N Go runtime scheduler dalam mengelola ribuan goroutine?
3. Apa perbedaan mendasar antara Unbuffered Channel dan Buffered Channel?
4. Mengapa penting selalu memanggil `defer cancel()` saat menggunakan `context.WithTimeout`?

**Kunci ringkas:** Error ditangani sebagai nilai biasa agar alur kegagalan terlihat jelas dan tidak ada eksekusi yang melompati kode tanpa disadari; Go scheduler memetakan M Goroutines ke N OS threads menggunakan work-stealing algorithm; Unbuffered channel memblokir pengirim hingga penerima siap (sinkronisasi langsung), sedangkan Buffered channel membolehkan pengiriman hingga kapasitas antrean penuh; `defer cancel()` segera melepaskan timer internal dan resource yang diasosiasikan dengan context saat fungsi selesai.

## Referensi

- [The Go Programming Language Official Docs](https://go.dev/doc/)
- [Effective Go](https://go.dev/doc/effective_go)
- [Uber Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md)
