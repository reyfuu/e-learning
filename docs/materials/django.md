# Materi: Django Framework dan Django REST Framework (DRF)

**Format:** teks · **Level:** Basic–Advanced · **Estimasi:** 60–90 menit

## Tujuan belajar

Setelah materi ini, pelajar dapat merancang backend berbasis Django dan Django REST Framework (DRF): memahami pola arsitektur MVT (*Model-View-Template*), memanfaatkan Django ORM secara optimal (QuerySets, Eager Loading via `select_related` / `prefetch_related`), membangun REST API type-safe dengan Serializers, serta mengonfigurasi autentikasi, permissions, dan Django Admin.

## Prasyarat

- Memahami bahasa Python modern (OOP, decorators, list comprehension, exceptions).
- Mengerti konsep dasar database relasional (tabel, primary key, foreign key).
- Python 3.10+ terpasang di sistem lokal.

## 1. Dasar: arsitektur MVT dan "Batteries-Included"

Django mengusung filosofi *The web framework for perfectionists with deadlines*. Menggunakan pola **MVT (Model-View-Template)**, Django menyediakan segalanya secara terintegrasi (*batteries-included*): sistem autentikasi, ORM tangguh, panel administrasi otomatis (Django Admin), proteksi CSRF/SQL Injection bawaan, dan sistem migrasi database otomatis.

### Peran MVT:
- **Model:** Mendefinisikan struktur database dan relasi antar-tabel dalam kelas Python murni.
- **View:** Memproses logika bisnis, menerima HTTP request, dan mengembalikan HTTP response (atau data JSON pada DRF).
- **Template:** Lapisan presentasi HTML (atau digantikan oleh API Serializer jika dipakai sebagai headless backend).

## 2. Praktik dasar: Model, Serializer, dan DRF ViewSet

```python
# courses/models.py
from django.db import models
from django.contrib.auth.models import User

class Category(models.Model):
    name = models.CharField(max_length=100, unique=True)
    slug = models.SlugField(max_length=100, unique=True)

    def __str__(self):
        return self.name

class Course(models.Model):
    title = models.CharField(max_length=200)
    description = models.TextField()
    price = models.DecimalField(max_digits=10, decimal_places=2)
    instructor = models.ForeignKey(User, on_delete=models.CASCADE, related_name='courses')
    category = models.ForeignKey(Category, on_delete=models.PROTECT, related_name='courses')
    is_published = models.BooleanField(default=False)
    created_at = models.DateTimeField(auto_now_add=True)

    class Meta:
        ordering = ['-created_at']

    def __str__(self):
        return self.title
```

```python
# courses/serializers.py
from rest_framework import serializers
from .models import Course, Category

class CategorySerializer(serializers.ModelSerializer):
    class Meta:
        model = Category
        fields = ['id', 'name', 'slug']

class CourseSerializer(serializers.ModelSerializer):
    instructor_name = serializers.ReadOnlyField(source='instructor.username')
    category = CategorySerializer(read_only=True)
    category_id = serializers.PrimaryKeyRelatedField(
        queryset=Category.objects.all(), source='category', write_only=True
    )

    class Meta:
        model = Course
        fields = [
            'id', 'title', 'description', 'price', 
            'instructor_name', 'category', 'category_id', 'is_published'
        ]
```

```python
# courses/views.py
from rest_framework import viewsets, permissions
from .models import Course
from .serializers import CourseSerializer

class CourseViewSet(viewsets.ModelViewSet):
    serializer_class = CourseSerializer
    permission_classes = [permissions.IsAuthenticatedOrReadOnly]

    def get_queryset(self):
        # Optimasi query dengan select_related untuk mencegah N+1
        return Course.objects.select_related('instructor', 'category').filter(is_published=True)

    def perform_create(self, serializer):
        # Otomatis kaitkan instruktur dengan user yang sedang login
        serializer.save(instructor=self.request.user)
```

## 3. Fitur utama dan pola arsitektur

### Optimasi Query Django ORM: `select_related` vs `prefetch_related`
Salah satu penyebab utama kelambatan Django di produksi adalah query N+1:
- `select_related()`: Menggunakan SQL `INNER/LEFT JOIN` untuk relasi satu-ke-satu (*OneToOne*) atau banyak-ke-satu (*ForeignKey*).
- `prefetch_related()`: Menjalankan query SQL terpisah kedua dan menggabungkannya di Python untuk relasi banyak-ke-banyak (*ManyToMany*) atau satu-ke-banyak terbalik.

```python
# Buruk (N+1 queries):
courses = Course.objects.all()
for c in courses:
    print(c.category.name) # Memicu query SQL baru per baris

# Baik (1 SQL JOIN query):
courses = Course.objects.select_related('category').all()
for c in courses:
    print(c.category.name) # Diambil langsung dari memori
```

### Custom Permissions di DRF
Membuat batasan akses agar hanya pembuat kursus yang dapat mengedit:

```python
from rest_framework import permissions

class IsInstructorOrReadOnly(permissions.BasePermission):
    def has_object_permission(self, request, view, obj):
        if request.method in permissions.SAFE_METHODS:
            return True
        return obj.instructor == request.user
```

## 4. Praktik kerja: optimasi transaksi keuangan dengan Database Atomic Transaction

### Skenario
Saat checkout kursus, sistem memotong saldo dompet siswa, menambahkan komisi ke instruktur, dan membuat record transaksi pendaftaran. Jika langkah ketiga gagal karena koneksi putus, saldo siswa sudah terlanjur terpotong tanpa ada akses ke materi (*inconsistent state*).

### Tugas
Bungkus seluruh alur operasi ke dalam `transaction.atomic()` dan gunakan `select_for_update()` untuk mencegah double spending.

```python
from django.db import transaction
from django.core.exceptions import ValidationError
from .models import Course, Enrollment, Wallet

def process_course_checkout(user, course_id):
    with transaction.atomic():
        # Kunci baris saldo dompet dengan SELECT ... FOR UPDATE
        wallet = Wallet.objects.select_for_update().get(user=user)
        course = Course.objects.get(id=course_id)

        if wallet.balance < course.price:
            raise ValidationError("Saldo dompet Anda tidak mencukupi.")

        # 1. Kurangi saldo siswa
        wallet.balance -= course.price
        wallet.save()

        # 2. Tambah saldo instruktur
        instructor_wallet = Wallet.objects.select_for_update().get(user=course.instructor)
        instructor_wallet.balance += course.price * 0.8 # 80% royalti
        instructor_wallet.save()

        # 3. Buat entri pendaftaran
        enrollment = Enrollment.objects.create(
            student=user,
            course=course,
            status='ACTIVE'
        )

        return enrollment
```

### Lima hint bertahap
1. Gunakan context manager `with transaction.atomic():` dari `django.db`.
2. Jika ada exception yang dilempar di dalam blok atomic, Django secara otomatis melakukan database `ROLLBACK`.
3. `select_for_update()` mengunci record pada level database hingga blok transaksi selesai.
4. Jangan pernah menaruh operasi I/O jaringan lambat (seperti kirim email atau panggil API pihak ketiga) di dalam blok `atomic()` karena akan memperpanjang waktu penguncian tabel.
5. Kirim event email di luar blok transaksi atau gunakan `transaction.on_commit(lambda: send_email())`.

## 5. Lanjutan: Celery untuk Background Tasks, Caching, dan Security Checklist

- **Celery & Redis:** Integrasikan Celery untuk memproses tugas lambat di background (pembuatan sertifikat PDF, encoding video).
- **Django Caching Framework:** Gunakan Redis cache untuk menyimpan hasil query kompleks via `@method_decorator(cache_page(60 * 15))`.
- **Production Security:** Jalankan `python manage.py check --deploy` untuk memeriksa konfigurasi keamanan (`SECURE_SSL_REDIRECT`, `SESSION_COOKIE_SECURE`, `CSRF_COOKIE_SECURE`).

## Pemeriksaan pemahaman

1. Apa perbedaan mendasar antara `select_related` dan `prefetch_related` di Django ORM?
2. Kapan `transaction.on_commit()` harus digunakan dalam alur bisnis?
3. Apa keunggulan ModelViewSet di Django REST Framework dibanding APIView biasa?
4. Mengapa penting mengunci baris data dengan `select_for_update()` saat memproses saldo keuangan?

**Kunci ringkas:** `select_related` menggunakan SQL JOIN untuk relasi single-valued sedangkan `prefetch_related` mengeksekusi query tambahan terpisah untuk relasi multi-valued; `transaction.on_commit()` memastikan task sampingan (seperti kirim email atau webhook) hanya dieksekusi setelah transaksi database berhasil di-COMMIT secara permanen; `ModelViewSet` mengotomatisasi operasi CRUD standar (list, create, retrieve, update, destroy) dalam satu kelas tanpa duplikasi boilerplate kode; `select_for_update()` mencegah race condition di mana request bersamaan membaca saldo lama sebelum saldo baru disimpan.

## Referensi

- [Django Official Documentation](https://docs.djangoproject.com/)
- [Django REST Framework Documentation](https://www.django-rest-framework.org/)
- [Django ORM CookBook](https://books.agiliq.com/projects/django-orm-cookbook/en/latest/)
