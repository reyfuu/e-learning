# Materi: Modern C++ (C++17/20) dan Rekayasa Sistem Berkinerja Tinggi

**Format:** teks · **Level:** Basic–Advanced · **Estimasi:** 60–90 menit

## Tujuan belajar

Setelah materi ini, pelajar dapat menulis kode Modern C++ (C++17/20) yang aman dan berkecepatan tinggi: memahami model memori Stack vs Heap, menerapkan idiom RAII (*Resource Acquisition Is Initialization*), menguasai Smart Pointers (`std::unique_ptr`, `std::shared_ptr`), memanfaatkan Move Semantics dan rvalue references, serta menggunakan pustaka STL (*Standard Template Library*) tanpa memory leaks.

## Prasyarat

- Memahami konsep dasar logika pemrograman (tipe data, percabangan, fungsi).
- Memiliki kompiler modern (GCC 11+, Clang 13+, atau MSVC) terpasang di sistem (`g++ --version`).
- Memahami konsep kompilasi (preprocessor, compilation, linking).

## 1. Dasar: filosofi "Zero-Overhead Abstraction" dan RAII

C++ dirancang oleh Bjarne Stroustrup dengan prinsip *Zero-Overhead Principle*: "Apa yang tidak Anda gunakan, tidak perlu Anda bayar. Dan apa yang Anda gunakan, tidak bisa Anda tulis lebih baik dengan tangan." C++ tidak menggunakan Garbage Collector; sebaliknya, manajemen resource dikendalikan secara deterministik melalui **RAII**.

### Konsep Kunci:
- **RAII:** Kepemilikan resource (memori, file handle, socket, mutex) diikat pada masa hidup (*lifetime*) objek di stack. Ketika objek keluar dari scope, destruktornya otomatis dipanggil untuk membersihkan resource.
- **Smart Pointers (C++11/14/17):**
  - `std::unique_ptr`: Kepemilikan tunggal eksklusif (tanpa overhead memori tambahan).
  - `std::shared_ptr`: Kepemilikan bersama dengan penghitung referensi (*reference counting*).
  - `std::weak_ptr`: Referensi non-owning untuk memutus siklus sirkular pada shared pointer.
- **Move Semantics:** Memindahkan kepemilikan resource alih-alih menyalin data secara mendalam (*deep copy*), menghemat siklus komputasi CPU.

## 2. Praktik dasar: RAII, Smart Pointers, dan STL Containers

```cpp
#include <iostream>
#include <memory>
#include <vector>
#include <string>
#include <algorithm>

// 1. Kelas Course dengan RAII
class Course {
private:
    std::string id;
    std::string title;
    double price;

public:
    Course(std::string id, std::string title, double price)
        : id(std::move(id)), title(std::move(title)), price(price) {
        std::cout << "[Constructed] Kursus: " << this->title << "\n";
    }

    ~Course() {
        std::cout << "[Destructed] Kursus: " << this->title << "\n";
    }

    const std::string& getTitle() const { return title; }
    double getPrice() const { return price; }
};

// 2. Fungsi pemrosesan modern dengan Smart Pointers
void processCourses() {
    // Alokasi memori unik yang aman tanpa resiko memory leak
    auto c1 = std::make_unique<Course>("cpp_101", "Modern C++ Core", 350000.0);
    auto c2 = std::make_unique<Course>("sys_201", "Linux System Programming", 450000.0);

    std::vector<std::unique_ptr<Course>> catalog;
    catalog.push_back(std::move(c1)); // Pindahkan kepemilikan ke vektor
    catalog.push_back(std::move(c2));

    // Iterasi modern dengan range-based for loop
    std::cout << "\n--- Katalog Kursus ---\n";
    for (const auto& coursePtr : catalog) {
        std::cout << "- " << coursePtr->getTitle() 
                  << " (Rp " << coursePtr->getPrice() << ")\n";
    }

    // Begitu vector catalog keluar dari scope fungsi ini, 
    // semua Course otomatis di-destruct tanpa perlu memanggil `delete` manual!
}

int main() {
    processCourses();
    std::cout << "Fungsi selesai, memori telah dibersihkan secara aman.\n";
    return 0;
}
```

## 3. Fitur utama dan pola arsitektur

### Move Semantics dan `std::move`
Sebelum C++11, pengembalian objek besar menyalin seluruh buffer memori. Dengan move semantics (`&&`), data cukup ditransfer pointer-nya:

```cpp
class DynamicBuffer {
    size_t size;
    int* data;
public:
    // Move Constructor
    DynamicBuffer(DynamicBuffer&& other) noexcept
        : size(other.size), data(other.data) {
        other.size = 0;
        other.data = nullptr; // Lepaskan kepemilikan dari objek asal
    }
};
```

### Modern Features: `std::optional`, `std::variant`, dan Structured Binding (C++17)

```cpp
#include <optional>
#include <tuple>

std::optional<std::string> findCourseCode(int id) {
    if (id == 42) return "CPP-ADV";
    return std::nullopt; // Menggantikan pointer NULL yang berbahaya
}

auto getCourseMeta() {
    return std::make_tuple(101, "C++ Architecture", true);
}

void demoFeatures() {
    auto [id, name, published] = getCourseMeta(); // Structured binding
    std::cout << "ID: " << id << ", Name: " << name << "\n";
}
```

## 4. Praktik kerja: mendiagnosis dan mencegah Segmentation Fault serta Memory Leak

### Skenario
Aplikasi pemrosesan audio/video streaming crash di server Linux dengan error `Segmentation Fault (core dumped)` setelah berjalan selama 3 jam, disertai lonjakan konsumsi memori yang tidak pernah turun.

### Tugas
Gunakan AddressSanitizer (ASan) pada tahap kompilasi untuk mendeteksi *dangling pointer* atau *use-after-free*, dan ubah raw pointer (`new`/`delete`) menjadi smart pointer yang aman.

```bash
# Kompilasi dengan AddressSanitizer & LeakSanitizer
g++ -std=c++17 -fsanitize=address,undefined -g -O1 main.cpp -o app_sanitized
./app_sanitized
```

Refactoring pencegahan bug:
```cpp
// SALAH (Raw Pointer rawan memory leak jika terjadi exception sebelum delete):
Course* raw = new Course("id", "Leak", 100);
if (someError) return; // MEMORY LEAK! delete raw tidak pernah dipanggil
delete raw;

// BENAR (Smart Pointer menjamin deallokasi):
auto safe = std::make_unique<Course>("id", "Safe", 100);
if (someError) return; // AMAN! Destruktor safe otomatis dipanggil saat stack unwind
```

### Lima hint bertahap
1. Jangan pernah menggunakan `new` dan `delete` mentah di kode modern C++; gunakan `std::make_unique` atau `std::make_shared`.
2. Aktifkan `-fsanitize=address` pada compiler GCC/Clang untuk menangkap akses memori ilegal tepat pada baris terjadinya.
3. Hindari circular reference pada `std::shared_ptr` dengan menggunakan `std::weak_ptr` pada salah satu arah relasi.
4. Jangan pernah mengembalikan referensi atau pointer ke variabel lokal yang hidup di stack.
5. Gunakan `const` correctness pada method yang tidak mengubah state internal objek.

## 5. Lanjutan: C++20 Concepts, Coroutines, dan Concurrency

- **Concepts (C++20):** Validasi batasan tipe template pada compile time dengan pesan error yang manusiawi:
  ```cpp
  template<typename T>
  requires std::integral<T>
  T add(T a, T b) { return a + b; }
  ```
- **Concurrency:** Manfaatkan `std::jthread` (otomatis join saat destruksi di C++20), `std::atomic`, dan `std::mutex` untuk sinkronisasi thread yang aman.
- **Valgrind:** Tool analisis profil memori untuk mendeteksi leak di level binary eksekusi.

## Pemeriksaan pemahaman

1. Apa itu prinsip RAII dan bagaimana prinsip ini mencegah kebocoran resource?
2. Kapan kita harus memilih `std::unique_ptr` dibanding `std::shared_ptr`?
3. Mengapa Move Constructor sebaiknya diberi penanda `noexcept`?
4. Apa manfaat `std::optional` dibanding mengembalikan pointer `nullptr` saat data tidak ditemukan?

**Kunci ringkas:** RAII mengikat siklus hidup resource ke lifetime objek di stack sehingga destruktor selalu otomatis membersihkan resource saat keluar scope; Gunakan `std::unique_ptr` secara default karena tanpa overhead memori (*zero cost*), dan hanya gunakan `std::shared_ptr` jika objek benar-benar memiliki multi-ownership lintas modul; `noexcept` memungkinkan kontainer STL seperti `std::vector` melakukan realokasi menggunakan pemindahan (*move*) yang cepat alih-alih penyalinan lambat demi menjamin strong exception guarantee; `std::optional` secara eksplisit mendokumentasikan kemungkinan ketiadaan nilai tanpa resiko runtime crash dereferensi pointer null.

## Referensi

- [cppreference.com - The Definitive C++ Reference](https://en.cppreference.com/)
- [C++ Core Guidelines by Bjarne Stroustrup and Herb Sutter](https://isocpp.github.io/CppCoreGuidelines/CppCoreGuidelines)
- [Effective Modern C++ by Scott Meyers](https://www.oreilly.com/library/view/effective-modern-c/9781491903988/)
