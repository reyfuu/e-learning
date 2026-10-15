# Materi: Python 3 Modern dan Rekayasa Kode Bersih

**Format:** teks · **Level:** Basic–Advanced · **Estimasi:** 60–90 menit

## Tujuan belajar

Setelah materi ini, pelajar dapat menulis kode Python modern (Python 3.11+) yang idiomatis: memahami sistem tipe data dan model objek data model (*dunder methods*), menguasai generators, iterators, dan context managers, memanfaatkan Type Hinting tingkat lanjut dengan modul `typing`, mengimplementasikan pemrograman asynchronous dengan `asyncio`, serta memahami manajemen memori dan GIL (*Global Interpreter Lock*).

## Prasyarat

- Memahami konsep dasar logika pemrograman (variabel, loop, fungsi).
- Memiliki Python 3.11+ terpasang di sistem (`python3 --version`).
- Memahami konsep virtual environment (`venv`).

## 1. Dasar: filosofi "Zen of Python" dan Data Model

Python dirancang dengan prinsip eksplisit, sederhana, dan mudah dibaca (*The Zen of Python*). Di Python, segala sesuatunya adalah objek. Perilaku objek diatur oleh protokol khusus yang disebut **Dunder Methods** (Magic Methods, seperti `__init__`, `__str__`, `__iter__`, `__enter__`).

### Fitur Kunci Python Modern:
- **Type Hints:** Anotasi tipe statis (`x: int = 5`, `Union`, `Optional`, `TypeVar`, `Self`).
- **Pattern Matching (PEP 634):** Pernyataan `match / case` untuk mendestrukturisasi struktur data kompleks.
- **Generators & `yield`:** Streaming data bernilai besar secara hemat memori (*lazy evaluation*).
- **Asynchronous I/O (`asyncio`):** Model eksekusi berbasis single-thread coroutine event loop untuk operasi I/O konkuren.

## 2. Praktik dasar: Type Hints, Dataclass, dan Generators

```python
from dataclasses import dataclass, field
from datetime import datetime
from typing import Generator, List, Optional
import uuid

@dataclass(frozen=True)
class Student:
    name: str
    email: str
    id: str = field(default_factory=lambda: str(uuid.uuid4()))
    enrolled_at: datetime = field(default_factory=datetime.utcnow)

@dataclass
class Course:
    title: str
    price: float
    tags: List[str] = field(default_factory=list)
    is_published: bool = False

    def apply_discount(self, percentage: float) -> None:
        if not (0 <= percentage <= 100):
            raise ValueError("Persentase diskon harus antara 0 dan 100.")
        self.price -= self.price * (percentage / 100)

# Generator hemat memori untuk memproses ribuan data
def stream_large_dataset(total_records: int) -> Generator[dict, None, None]:
    for i in range(total_records):
        # Data dihasilkan satu per satu tanpa disimpan sekaligus di memori RAM
        yield {
            "record_id": i,
            "status": "processed" if i % 2 == 0 else "pending"
        }

if __name__ == "__main__":
    course = Course(title="Python Advanced Engineering", price=500_000.0, tags=["python", "backend"])
    course.apply_discount(20)
    print(f"Harga setelah diskon: Rp {course.price:,.2f}")

    # Konsumsi generator secara efisien
    generator = stream_large_dataset(1_000_000)
    first_item = next(generator)
    print("Item pertama:", first_item)
```

## 3. Fitur utama dan pola arsitektur

### Context Managers (`with` statement)
Menjamin alokasi dan pelepasan resource (file socket, database lock) secara deterministik:

```python
from contextlib import contextmanager
import time

@contextmanager
def execution_timer(label: str):
    start = time.perf_counter()
    try:
        yield
    finally:
        elapsed = time.perf_counter() - start
        print(f"[{label}] Selesai dalam {elapsed:.4f} detik")

# Penggunaan:
with execution_timer("Parsing File Log"):
    time.sleep(0.1)
```

### Pemrograman Asynchronous dengan `asyncio`
Menjalankan ribuan tugas I/O tanpa memblokir thread utama:

```python
import asyncio
import httpx

async def fetch_course_metadata(course_id: str) -> dict:
    async with httpx.AsyncClient() as client:
        response = await client.get(f"https://api.learnhub.id/courses/{course_id}")
        return response.json()

async def main_async():
    course_ids = ["py_101", "go_201", "ts_301"]
    tasks = [fetch_course_metadata(cid) for cid in course_ids]
    results = await asyncio.gather(*tasks, return_exceptions=True)
    print("Hasil fetch serentak:", results)

# asyncio.run(main_async())
```

## 4. Praktik kerja: optimasi memory footprint dengan `__slots__`

### Skenario
Aplikasi analisis data LearnHub memuat 2.000.000 objek pelacakan event klik pengguna ke dalam memori. Server kehabisan RAM (OOM Crash) karena setiap objek Python biasa memiliki dictionary internal `__dict__` yang berukuran besar.

### Tugas
Terapkan `__slots__` pada kelas event untuk menghilangkan `__dict__` dinamis, menghemat penggunaan memori hingga 60%.

```python
import sys

# Kelas biasa (menggunakan dictionary internal)
class StandardEvent:
    def __init__(self, user_id: str, action: str, timestamp: float):
        self.user_id = user_id
        self.action = action
        self.timestamp = timestamp

# Kelas teroptimasi dengan __slots__
class OptimizedEvent:
    __slots__ = ('user_id', 'action', 'timestamp')

    def __init__(self, user_id: str, action: str, timestamp: float):
        self.user_id = user_id
        self.action = action
        self.timestamp = timestamp

if __name__ == "__main__":
    e1 = StandardEvent("usr_1", "click", 1700000.0)
    e2 = OptimizedEvent("usr_1", "click", 1700000.0)

    size_e1 = sys.getsizeof(e1) + sys.getsizeof(e1.__dict__)
    size_e2 = sys.getsizeof(e2)

    print(f"Ukuran memori StandardEvent: {size_e1} bytes")
    print(f"Ukuran memori OptimizedEvent: {size_e2} bytes")
```

### Lima hint bertahap
1. Objek kelas Python standar menyimpan atribut instance-nya dalam dictionary `__dict__`.
2. Dictionary memiliki alokasi memori tambahan (*hash table overhead*) untuk mendukung penambahan atribut secara dinamis.
3. Mendefinisikan tuple `__slots__` memberitahu interpreter untuk mengalokasikan array tetap untuk atribut yang terdaftar saja.
4. Kelas dengan `__slots__` tidak mengizinkan penambahan atribut baru di luar daftar yang dideklarasikan.
5. Jika membuat subclass, kelas anak juga harus mendefinisikan `__slots__` agar tidak otomatis membuat `__dict__`.

## 5. Lanjutan: GIL, Multiprocessing, dan Profiling CPython

- **Global Interpreter Lock (GIL):** Mutex di CPython yang mencegah multi-threading mengeksekusi bytecode CPython secara paralel murni pada CPU-bound tasks.
- **Multiprocessing vs Asyncio:** Gunakan `asyncio` untuk operasi I/O-bound (HTTP, database) dan `multiprocessing` atau modul `concurrent.futures.ProcessPoolExecutor` untuk tugas intensif CPU.
- **Memory Profiling:** Gunakan modul bawaan `tracemalloc` dan alat pihak ketiga `cProfile` untuk mendeteksi memory leak dan bottleneck performa fungsi.

## Pemeriksaan pemahaman

1. Apa perbedaan cara kerja `return` dan `yield` pada fungsi Python?
2. Bagaimana mekanisme Global Interpreter Lock (GIL) memengaruhi thread Python?
3. Apa keuntungan menggunakan `@dataclass(frozen=True)`?
4. Kapan kita sebaiknya menggunakan `asyncio` dibanding multithreading?

**Kunci ringkas:** `return` menghentikan fungsi dan mengembalikan satu nilai, sedangkan `yield` menjeda fungsi dan mempertahankan state lokalnya untuk menghasilkan urutan nilai berikutnya secara lazy; GIL membatasi hanya satu thread yang dapat mengeksekusi bytecode Python pada satu waktu sehingga multi-threading tidak meningkatkan performa CPU-bound murni; `frozen=True` membuat instance dataclass menjadi immutable dan dapat di-hash (*hashable*) sebagai key dictionary; `asyncio` lebih ideal untuk ribuan operasi jaringan I/O konkuren karena tidak memerlukan context switching tingkat kernel thread OS.

## Referensi

- [Python Official Documentation](https://docs.python.org/3/)
- [Real Python Tutorials](https://realpython.com/)
- [Fluent Python by Luciano Ramalho](https://www.oreilly.com/library/view/fluent-python-2nd/9781492056348/)
