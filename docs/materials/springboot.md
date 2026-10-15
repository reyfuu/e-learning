# Materi: Spring Boot 3 untuk Layanan REST Enterprise

**Format:** teks · **Level:** Basic–Advanced · **Estimasi:** 60–90 menit

## Tujuan belajar

Setelah materi ini, pelajar dapat merancang layanan backend enterprise berbasis Spring Boot 3: memahami konsep Inversion of Control (IoC) dan Dependency Injection, membangun RESTful API dengan Spring MVC, mengelola persistensi data menggunakan Spring Data JPA & Hibernate, mengamankan endpoint dengan Spring Security (JWT / OAuth2), serta memantau kesehatan aplikasi dengan Spring Boot Actuator.

## Prasyarat

- Memahami Java modern (Java 17+, OOP, Streams, Annotations).
- Memahami konsep dasar relational database dan SQL.
- Maven atau Gradle terpasang di sistem.

## 1. Dasar: ekosistem Spring Framework dan Spring Boot

Spring Boot menyederhanakan pengembangan aplikasi Spring dengan mengadopsi prinsip **Convention over Configuration**. Melalui fitur *Starters*, *Auto-configuration*, dan *Embedded Server* (seperti Tomcat bawaan), pengembang dapat langsung menjalankan aplikasi tanpa konfigurasi XML yang rumit.

### Anotasi Esensial:
- `@SpringBootApplication`: Mengaktifkan konfigurasi otomatis, component scanning, dan konfigurasi Spring Boot.
- `@RestController`: Menggabungkan `@Controller` dan `@ResponseBody` untuk menghasilkan respons JSON langsung.
- `@Service` & `@Repository`: Stereotype anotasi untuk mendefinisikan komponen logika bisnis dan akses data.
- `@Autowired` / Konstruktor Injection: Menyerahkan resolusi dependensi ke Spring ApplicationContext.

## 2. Praktik dasar: REST Controller, DTO Validation, dan Spring Data JPA

```java
// Entity
package com.learnhub.entity;

import jakarta.persistence.*;
import jakarta.validation.constraints.*;
import java.math.BigDecimal;

@Entity
@Table(name = "courses")
public class Course {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;

    @NotBlank(message = "Judul kursus wajib diisi")
    @Column(nullable = false)
    private String title;

    @NotNull
    @DecimalMin(value = "0.0", inclusive = true)
    @Column(nullable = false)
    private BigDecimal price;

    // Constructors, Getters, and Setters
    public Course() {}

    public Course(String title, BigDecimal price) {
        this.title = title;
        this.price = price;
    }

    public Long getId() { return id; }
    public String getTitle() { return title; }
    public void setTitle(String title) { this.title = title; }
    public BigDecimal getPrice() { return price; }
    public void setPrice(BigDecimal price) { this.price = price; }
}
```

```java
// Repository
package com.learnhub.repository;

import com.learnhub.entity.Course;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;
import java.util.List;

@Repository
public interface CourseRepository extends JpaRepository<Course, Long> {
    List<Course> findByTitleContainingIgnoreCase(String keyword);
}
```

```java
// Controller
package com.learnhub.controller;

import com.learnhub.entity.Course;
import com.learnhub.repository.CourseRepository;
import jakarta.validation.Valid;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;
import org.springframework.web.server.ResponseStatusException;

import java.util.List;

@RestController
@RequestMapping("/api/v1/courses")
public class CourseController {

    private final CourseRepository courseRepository;

    // Rekomendasi Spring: Constructor Injection
    public CourseController(CourseRepository courseRepository) {
        this.courseRepository = courseRepository;
    }

    @GetMapping
    public List<Course> getAllCourses(@RequestParam(required = false) String search) {
        if (search != null && !search.isBlank()) {
            return courseRepository.findByTitleContainingIgnoreCase(search);
        }
        return courseRepository.findAll();
    }

    @PostMapping
    public ResponseEntity<Course> createCourse(@Valid @RequestBody Course course) {
        Course saved = courseRepository.save(course);
        return ResponseEntity.status(HttpStatus.CREATED).body(saved);
    }

    @GetMapping("/{id}")
    public Course getCourseById(@PathVariable Long id) {
        return courseRepository.findById(id)
            .orElseThrow(() -> new ResponseStatusException(HttpStatus.NOT_FOUND, "Kursus tidak ditemukan"));
    }
}
```

## 3. Fitur utama dan pola arsitektur

### Pengelolaan Transaksi dengan `@Transactional`
Menjaga konsistensi data atomik dalam operasi bisnis multi-langkah:

```java
@Service
public class EnrollmentService {

    private final CourseRepository courseRepository;
    private final EnrollmentRepository enrollmentRepository;

    public EnrollmentService(CourseRepository courseRepository, EnrollmentRepository enrollmentRepository) {
        this.courseRepository = courseRepository;
        this.enrollmentRepository = enrollmentRepository;
    }

    @Transactional(rollbackFor = Exception.class)
    public void enrollStudent(Long studentId, Long courseId) {
        Course course = courseRepository.findById(courseId)
            .orElseThrow(() -> new IllegalArgumentException("Kursus tidak valid"));

        // Lakukan pemotongan saldo / mutasi
        enrollmentRepository.save(new Enrollment(studentId, course.getId()));
    }
}
```

### Global Exception Handler dengan `@RestControllerAdvice`
Menangkap exception terpusat dan mengembalikan format error standar RFC 7807:

```java
@RestControllerAdvice
public class GlobalExceptionHandler {

    @ExceptionHandler(ResponseStatusException.class)
    public ResponseEntity<Map<String, Object>> handleResponseStatus(ResponseStatusException ex) {
        return ResponseEntity.status(ex.getStatusCode()).body(Map.of(
            "timestamp", Instant.now(),
            "status", ex.getStatusCode().value(),
            "error", ex.getReason()
        ));
    }
}
```

## 4. Praktik kerja: optimasi hibernate query dan deteksi N+1 problem

### Skenario
Endpoint `/api/v1/instructors` memuat 50 instruktur beserta daftar kursus yang mereka ajar. Setiap akses endpoint memicu 51 query SQL terpisah ke PostgreSQL, menyebabkan latency respons melonjak hingga 1,8 detik.

### Tugas
Gunakan anotasi `@EntityGraph` atau query kustom `JOIN FETCH` di Spring Data JPA untuk memuat data relasi dalam 1 query SQL efisien.

```java
package com.learnhub.repository;

import com.learnhub.entity.Instructor;
import org.springframework.data.jpa.repository.EntityGraph;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Query;
import java.util.List;

public interface InstructorRepository extends JpaRepository<Instructor, Long> {

    // Cara 1: Menggunakan EntityGraph
    @EntityGraph(attributePaths = {"courses"})
    List<Instructor> findAll();

    // Cara 2: Menggunakan explicit JOIN FETCH
    @Query("SELECT DISTINCT i FROM Instructor i LEFT JOIN FETCH i.courses")
    List<Instructor> findAllWithCoursesFetch();
}
```

### Lima hint bertahap
1. Periksa log SQL Hibernate dengan mengaktifkan properti `spring.jpa.show-sql=true` di `application.properties`.
2. Relasi bertipe `@OneToMany` memiliki default fetch type `LAZY`.
3. Mengakses koleksi relasi lazy di dalam loop memicu satu query baru per elemen (*N+1 problem*).
4. `JOIN FETCH` di JPQL memberitahu Hibernate untuk melakukan SQL `LEFT OUTER JOIN` dan menginstansiasi objek anak dalam 1 query.
5. Gunakan `DISTINCT` pada JPQL untuk mencegah duplikasi entitas induk di memori hasil query join.

## 5. Lanjutan: Spring Security 6, JWT, dan Actuator

- **Spring Security 6:** Menggunakan DSL lambda modern pada `SecurityFilterChain` tanpa memperpanjang kelas adapter deprecated.
- **Actuator:** Endpoint pemantauan sistem (`/actuator/health`, `/actuator/metrics`, `/actuator/prometheus`) untuk integrasi dengan Grafana.
- **Spring Boot Native (GraalVM):** Kompilasi kode Java langsung ke biner native Linux untuk waktu startup sub-detik (<50ms) dan penggunaan memori sangat kecil.

## Pemeriksaan pemahaman

1. Mengapa Constructor Injection lebih disukai daripada Field Injection (`@Autowired` di variabel field)?
2. Bagaimana cara kerja `@Transactional` di balik layar menggunakan Spring AOP Proxy?
3. Apa perbedaan antara `FetchType.EAGER` dan `FetchType.LAZY` di Hibernate?
4. Mengapa kita harus berhati-hati saat menggunakan `findAll()` pada entitas dengan banyak relasi?

**Kunci ringkas:** Constructor Injection memfasilitasi *immutability* (field dapat dideklarasikan `final`), memudahkan unit testing tanpa refleksi, dan mencegah *circular dependency*; Spring membungkus kelas `@Transactional` dalam dynamic proxy yang membuka transaksi sebelum metode berjalan dan melakukan commit/rollback setelahnya; `EAGER` selalu memuat relasi saat itu juga, sedangkan `LAZY` hanya memuat data ketika getter relasi pertama kali dipanggil; `findAll()` tanpa filter pagination atau tanpa fetch join dapat memicu penarikan jutaan baris data ke memori RAM hingga JVM crash.

## Referensi

- [Spring Boot Reference Documentation](https://docs.spring.io/spring-boot/docs/current/reference/html/)
- [Spring Guides: Building a RESTful Web Service](https://spring.io/guides/gs/rest-service/)
- [Vlad Mihalcea on Hibernate Performance](https://vladmihalcea.com/)
