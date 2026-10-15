# Materi: FastAPI untuk Backend Asynchronous Berkinerja Tinggi

**Format:** teks · **Level:** Basic–Advanced · **Estimasi:** 60–90 menit

## Tujuan belajar

Setelah materi ini, pelajar dapat membangun API asynchronous modern berkinerja tinggi dengan FastAPI: memanfaatkan model validasi Pydantic v2, menguasai Dependency Injection bawaan (`Depends`), menangani operasi async I/O database, mengamankan endpoint dengan OAuth2 & JWT tokens, serta memanfaatkan dokumentasi interaktif Swagger/OpenAPI otomatis.

## Prasyarat

- Memahami Python 3.10+ (Type Hinting, async/await coroutines).
- Mengerti konsep protokol HTTP dan format pertukaran data JSON.
- Python 3.10+ terpasang di sistem.

## 1. Dasar: arsitektur ASGI dan ekosistem FastAPI

FastAPI adalah framework web modern berbasis standar tipe data Python yang berjalan di atas server **ASGI (Asynchronous Server Gateway Interface)** seperti Uvicorn. FastAPI menggabungkan dua kekuatan utama:
1. **Starlette:** Engine web asynchronous berkecepatan tinggi.
2. **Pydantic:** Pustaka parsing dan validasi data berbasis skema tipe data Python yang sangat cepat (ditulis ulang dalam Rust pada v2).

### Keunggulan Utama:
- **Kinerja Sangat Tinggi:** Setara dengan NodeJS dan Go dalam berbagai benchmark I/O.
- **Otomatisasi Dokumentasi:** Menghasilkan dokumentasi interaktif OpenAPI (Swagger UI di `/docs` dan ReDoc di `/redoc`) secara otomatis tanpa plugin tambahan.
- **Validasi Data Statis & Runtime:** Error otomatis dikembalikan dengan kode `422 Unprocessable Entity` jika tipe data tidak cocok.

## 2. Praktik dasar: Pydantic Schema, Endpoints, dan Dependency Injection

```python
from fastapi import FastAPI, Depends, HTTPException, status
from pydantic import BaseModel, Field
from typing import List, Optional
import uuid

app = FastAPI(
    title="LearnHub API",
    description="Backend API untuk platform e-learning",
    version="1.0.0"
)

# 1. Pydantic Models untuk Validasi Request & Response
class CourseBase(BaseModel):
    title: str = Field(..., min_length=5, max_length=150, example="Mastering FastAPI & Asyncio")
    description: str = Field(..., example="Panduan lengkap membangun API performa tinggi")
    price: float = Field(..., ge=0, example=250000.0)
    category: str = Field(..., example="Backend")

class CourseCreate(CourseBase):
    pass

class CourseResponse(CourseBase):
    id: str
    is_published: bool

    class Config:
        from_attributes = True

# Simulasi Database In-Memory
fake_db_courses = {}

# 2. Dependency Injection untuk Autentikasi / Sesi
async def verify_instructor_token(token: str = "default_token"):
    if token != "secret-instructor-token":
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail="Token instruktur tidak valid atau kadaluarsa"
        )
    return {"user_id": "inst_001", "role": "instructor"}

# 3. Path Operations
@app.post(
    "/api/v1/courses",
    response_model=CourseResponse,
    status_code=status.HTTP_201_CREATED,
    tags=["Courses"]
)
async def create_course(
    course_in: CourseCreate,
    current_user: dict = Depends(verify_instructor_token)
):
    course_id = f"crs_{uuid.uuid4().hex[:8]}"
    course_data = course_in.model_dump()
    course_record = {**course_data, "id": course_id, "is_published": False}
    fake_db_courses[course_id] = course_record
    return course_record

@app.get(
    "/api/v1/courses",
    response_model=List[CourseResponse],
    tags=["Courses"]
)
async def list_courses(skip: int = 0, limit: int = 20):
    all_courses = list(fake_db_courses.values())
    return all_courses[skip : skip + limit]
```

## 3. Fitur utama dan pola arsitektur

### Async vs Def (Sync) di FastAPI
FastAPI menangani fungsi rute `async def` dan `def` secara berbeda:
- `async def`: Dijalankan langsung di main thread Event Loop. Jangan panggil operasi blocking (seperti `time.sleep()` atau driver DB sync seperti `psycopg2`) di sini karena akan membekukan seluruh server.
- `def` (biasa): FastAPI secara otomatis mengirim fungsi ini ke threadpool eksternal (`anyio/starlette worker threads`) sehingga aman untuk kode blocking lama.

```python
# Benar untuk async I/O library (misal asyncpg, httpx)
@app.get("/async-data")
async def get_data():
    async with httpx.AsyncClient() as client:
        return (await client.get("https://httpbin.org/delay/1")).json()

# Benar untuk synchronous blocking operations
@app.get("/sync-data")
def get_sync_data():
    time.sleep(1) # Aman, dijalankan di threadpool terpisah
    return {"status": "ok"}
```

### Background Tasks
Menjalankan operasi pasca-respons (seperti mengirim email verifikasi) tanpa memperlambat waktu respon pengguna:

```python
from fastapi import BackgroundTasks

def send_welcome_email(email: str):
    # Logika kirim email lambat
    pass

@app.post("/register")
async def register_student(email: str, background_tasks: BackgroundTasks):
    background_tasks.add_task(send_welcome_email, email)
    return {"message": "Pendaftaran berhasil, email konfirmasi sedang dikirim"}
```

## 4. Praktik kerja: integrasi database asynchronous dengan SQLAlchemy 2.0 & Asyncpg

### Skenario
Aplikasi mengalami penurunan throughput drastis saat pengguna ramai mengakses katalog karena menggunakan driver database sinkron yang memblokir coroutine event loop.

### Tugas
Konfigurasikan SQLAlchemy 2.0 dengan mode `create_async_engine` dan `AsyncSession` menggunakan driver PostgreSQL `asyncpg`.

```python
from sqlalchemy.ext.asyncio import create_async_engine, async_sessionmaker, AsyncSession
from sqlalchemy.orm import DeclarativeBase, Mapped, mapped_column
from sqlalchemy import select

DATABASE_URL = "postgresql+asyncpg://postgres:postgres@localhost:5432/learnhub"

engine = create_async_engine(DATABASE_URL, echo=False, pool_size=20, max_overflow=10)
AsyncSessionLocal = async_sessionmaker(engine, expire_on_commit=False, class_=AsyncSession)

class Base(DeclarativeBase):
    pass

class CourseModel(Base):
    __tablename__ = "courses"
    id: Mapped[int] = mapped_column(primary_key=True)
    title: Mapped[str]

# Dependency yield AsyncSession per request
async def get_db() -> AsyncSession:
    async with AsyncSessionLocal() as session:
        try:
            yield session
        finally:
            await session.close()

@app.get("/courses-db")
async def get_courses_from_db(db: AsyncSession = Depends(get_db)):
    result = await db.execute(select(CourseModel).limit(50))
    courses = result.scalars().all()
    return courses
```

### Lima hint bertahap
1. Pastikan connection string diawali dengan `postgresql+asyncpg://` bukan `postgresql://`.
2. Gunakan `create_async_engine` dari modul `sqlalchemy.ext.asyncio`.
3. Selalu sertakan keyword `await` pada pemanggilan `db.execute()`, `db.commit()`, dan `session.close()`.
4. Definisikan dependency generator dengan `yield session` agar sesi otomatis ditutup setelah endpoint merespons.
5. Gunakan `AsyncSessionLocal(expire_on_commit=False)` agar atribut objek tetap dapat diakses setelah commit tanpa memicu lazy loading implisit yang gagal.

## 5. Lanjutan: OAuth2 Password Bearer, CORS, dan Uvicorn Tuning

- **Uvicorn Workers:** Jalankan server produksi di balik Nginx menggunakan beberapa worker CPU:
  ```bash
  uvicorn main:app --workers 4 --host 127.0.0.1 --port 8000
  ```
- **CORSMiddleware:** Konfigurasikan domain asal yang diizinkan untuk melindungi cookie dan header otentikasi.
- **Lifespan Events:** Gunakan `@asynccontextmanager` modern pada inisialisasi aplikasi untuk mengelola koneksi startup dan shutdown pool basis data.

## Pemeriksaan pemahaman

1. Apa bahaya memanggil library I/O blocking (seperti `requests.get()`) di dalam endpoint `async def` FastAPI?
2. Bagaimana Pydantic v2 meningkatkan performa parsing JSON di FastAPI?
3. Apa fungsi mekanisme `Depends()` di parameter fungsi rute FastAPI?
4. Kapan sebaiknya menggunakan `BackgroundTasks` bawaan FastAPI dibanding task queue eksternal seperti Celery?

**Kunci ringkas:** Library blocking di coroutine `async def` membekukan event loop single-thread utama sehingga request pengguna lain terhenti total; Pydantic v2 menulis ulang core engine validasinya dalam bahasa Rust sehingga parsing data 5-10 kali lebih cepat; `Depends()` menyediakan sistem dependency injection yang modular untuk mengekstrak logika otentikasi, koneksi DB, dan konfigurasi; `BackgroundTasks` cocok untuk tugas ringan dalam proses yang sama, sedangkan Celery dibutuhkan untuk tugas berdurasi panjang, komputasi berat, atau yang membutuhkan retry persisten jika server restart.

## Referensi

- [FastAPI Official Documentation](https://fastapi.tiangolo.com/)
- [Pydantic v2 Documentation](https://docs.pydantic.dev/latest/)
- [SQLAlchemy 2.0 Async Documentation](https://docs.sqlalchemy.org/en/20/orm/extensions/asyncio.html)
