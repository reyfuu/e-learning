# Materi: Java Modern (Java 17/21 LTS) dan Pemrograman Berorientasi Objek

**Format:** teks · **Level:** Basic–Advanced · **Estimasi:** 60–90 menit

## Tujuan belajar

Setelah materi ini, pelajar dapat memahami fitur-fitur Java modern (Java 17 & 21 LTS): menguasai prinsip Object-Oriented Programming (OOP) yang bersih, memanfaatkan Records, Pattern Matching, Stream API, Generics, Concurrency (termasuk Virtual Threads / Project Loom), serta manajemen memori JVM dan Exception Handling.

## Prasyarat

- Memahami konsep dasar algoritma dan struktur data.
- Memiliki JDK 17 atau 21 terpasang di sistem (`java -version`).
- Memahami konsep dasar build tools seperti Maven atau Gradle.

## 1. Dasar: evolusi Java modern dan JVM

Java adalah bahasa berorientasi objek yang berjalan di atas Java Virtual Machine (JVM), mengusung filosofi *Write Once, Run Anywhere* (WORA). Java modern telah bertransformasi dengan sintaks yang jauh lebih ringkas, tipe data *immutable* bawaan, dan model konkurensi revolusioner.

### Fitur Kunci Java 17 & 21:
- **Records:** Cara ringkas untuk mendefinisikan kelas pembawa data *immutable* tanpa *boilerplate* getter/setter/toString.
- **Pattern Matching for Switch & `instanceof`:** Mempermudah pengecekan dan ekstraksi tipe data.
- **Sealed Classes:** Mengontrol secara ketat kelas mana saja yang diizinkan mewarisi suatu *interface* atau *abstract class*.
- **Virtual Threads (Project Loom di Java 21):** Thread ringan dengan jejak memori minim untuk throughput I/O tingkat tinggi.

## 2. Praktik dasar: Record, Sealed Interface, dan Stream API

```java
package com.learnhub.model;

import java.time.Instant;
import java.util.List;
import java.util.stream.Collectors;

// 1. Sealed Interface untuk memodelkan tipe hasil yang pasti
public sealed interface EnrollmentResult permits EnrollmentSuccess, EnrollmentFailure {}

record EnrollmentSuccess(String studentId, String courseId, Instant enrolledAt) implements EnrollmentResult {}
record EnrollmentFailure(String reason, int errorCode) implements EnrollmentResult {}

// 2. Record Course
public record Course(
    String id,
    String title,
    double price,
    boolean isPublished,
    List<String> tags
) {}

// 3. Pemrosesan Data dengan Stream API
public class CourseAnalytics {

    public static List<String> getPremiumCourseTitles(List<Course> courses, double minPrice) {
        return courses.stream()
            .filter(Course::isPublished)
            .filter(c -> c.price() >= minPrice)
            .map(Course::title)
            .sorted()
            .collect(Collectors.toList());
    }

    public static double calculateAveragePrice(List<Course> courses) {
        return courses.stream()
            .mapToDouble(Course::price)
            .average()
            .orElse(0.0);
    }
}
```

## 3. Fitur utama dan pola arsitektur

### Pattern Matching untuk Switch
Menyederhanakan kontrol alur berbasis tipe data objek secara deklaratif:

```java
public class EnrollmentNotifier {

    public static String formatResult(EnrollmentResult result) {
        return switch (result) {
            case EnrollmentSuccess s -> 
                "Selamat! Pendaftaran berhasil untuk siswa " + s.studentId() + " pada kursus " + s.courseId();
            case EnrollmentFailure f -> 
                "Pendaftaran gagal (" + f.errorCode() + "): " + f.reason();
        };
    }
}
```

### Generics dan Collections Framework
Java menyediakan hirarki `Collection` yang kuat (`List`, `Set`, `Map`) dengan keamanan tipe generik:

```java
import java.util.*;

public class Repository<T> {
    private final Map<String, T> storage = new ConcurrentHashMap<>();

    public void save(String id, T item) {
        storage.put(id, item);
    }

    public Optional<T> findById(String id) {
        return Optional.ofNullable(storage.get(id));
    }
}
```

## 4. Praktik kerja: konkurensi modern dengan Virtual Threads (Java 21)

### Skenario
Aplikasi pemrosesan sertifikat harus memverifikasi status kelulusan 10.000 pelajar secara simultan ke API eksternal. Menggunakan Thread standar OS (`Platform Threads`) menyebabkan konsumsi memori membengkak hingga kehabisan batas memori (*OutOfMemoryError: unable to create native thread*).

### Tugas
Manfaatkan `Executors.newVirtualThreadPerTaskExecutor()` untuk mengeksekusi ribuan tugas I/O blocking secara konkuren dengan konsumsi RAM yang sangat rendah.

```java
package com.learnhub.service;

import java.net.URI;
import java.net.http.HttpClient;
import java.net.http.HttpRequest;
import java.net.http.HttpResponse;
import java.util.List;
import java.util.concurrent.*;

public class CertificateBatchService {

    private final HttpClient httpClient = HttpClient.newHttpClient();

    public void processAllCertificates(List<String> studentIds) {
        // Menggunakan Virtual Threads Executor (Java 21)
        try (var executor = Executors.newVirtualThreadPerTaskExecutor()) {
            List<Future<String>> futures = studentIds.stream()
                .map(id -> executor.submit(() -> verifyAndIssue(id)))
                .toList();

            for (Future<String> future : futures) {
                try {
                    String result = future.get(5, TimeUnit.SECONDS);
                    System.out.println("Hasil: " + result);
                } catch (TimeoutException | ExecutionException | InterruptedException e) {
                    System.err.println("Gagal memproses item: " + e.getMessage());
                }
            }
        }
    }

    private String verifyAndIssue(String studentId) throws Exception {
        var request = HttpRequest.newBuilder()
            .uri(URI.create("https://api.learnhub.id/verify/" + studentId))
            .GET()
            .build();

        var response = httpClient.send(request, HttpResponse.BodyHandlers.ofString());
        return "Siswa " + studentId + " status: " + response.statusCode();
    }
}
```

### Lima hint bertahap
1. Platform thread memetakan 1:1 ke kernel thread OS dengan alokasi stack sekitar 1 MB per thread.
2. Virtual thread dikelola oleh JVM runtime di user-space dengan footprint memori hanya beberapa kilobyte.
3. Virtual thread sangat ideal untuk tugas yang bersifat I/O-bound (panggilan jaringan atau query database).
4. Gunakan `try-with-resources` pada `ExecutorService` untuk menjamin semua task selesai sebelum scope keluar.
5. Hindari `synchronized` block yang panjang pada virtual thread jika memicu *thread pinning*; gunakan `ReentrantLock` sebagai alternatif.

## 5. Lanjutan: JVM Memory Management, GC Tuning, dan Profiling

- **JVM Memory Areas:** Heap Memory (Eden, Survivor, Old Gen) dan Non-Heap (Metaspace, Thread Stacks).
- **Garbage Collectors:** G1GC (default), ZGC (ultra-low latency untuk heap besar), dan Shenandoah.
- **Java Flight Recorder (JFR):** Profiler bawaan JVM ber-overhead sangat rendah untuk mendeteksi alokasi memori berlebih dan CPU lock contention.

## Pemeriksaan pemahaman

1. Apa perbedaan utama Java Record dengan kelas Java konvensional?
2. Mengapa Sealed Classes penting dalam perancangan arsitektur domain?
3. Kapan Virtual Threads memberikan keuntungan performa terbesar dibanding Platform Threads?
4. Mengapa `Optional` diperkenalkan di Java dan apa aturan penggunaannya?

**Kunci ringkas:** Record secara otomatis menghasilkan konstruktor, getter, `equals()`, `hashCode()`, dan `toString()`, serta bersifat immutable; Sealed Classes membatasi hierarki pewarisan sehingga kompiler dapat memverifikasi kelengkapan kasus switch secara menyeluruh; Virtual Threads unggul pada tugas I/O-bound throughput tinggi karena tidak memblokir kernel thread OS; `Optional` mencegah bug `NullPointerException` dengan mewajibkan penanganan nilai kosong secara eksplisit pada return value method.

## Referensi

- [Oracle Java SE Documentation](https://docs.oracle.com/en/java/)
- [Java 21 Feature Guides (OpenJDK)](https://openjdk.org/projects/jdk/21/)
- [Baeldung Java Tutorials](https://www.baeldung.com/category/java/)
