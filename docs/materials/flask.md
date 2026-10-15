# Materi: Flask Microframework untuk Layanan Web Python Fleksibel

**Format:** teks · **Level:** Basic–Advanced · **Estimasi:** 60–90 menit

## Tujuan belajar

Setelah materi ini, pelajar dapat mengembangkan aplikasi web dan REST API yang modular menggunakan Flask: memahami siklus request context dan application context, menerapkan Application Factory Pattern, mengorganisasi kode menggunakan Blueprints, mengintegrasikan basis data dengan Flask-SQLAlchemy, serta mengamankan konfigurasi dan deployment di lingkungan produksi dengan Gunicorn.

## Prasyarat

- Memahami bahasa Python modern (fungsi, decorators, dictionary, exception handling).
- Memahami protokol HTTP dan arsitektur REST.
- Python 3.10+ terpasang di sistem lokal.

## 1. Dasar: arsitektur Flask dan filosofi microframework

Flask adalah microframework berbasis Werkzeug (routing WSGI) dan Jinja2 (template engine). Dikatakan "micro" bukan berarti fiturnya sedikit, melainkan Flask tidak memaksakan struktur direktori tertentu ataupun ORM bawaan, memberikan kebebasan penuh kepada developer untuk memilih pustaka yang sesuai.

### Konsep Konteks di Flask:
- **Application Context (`current_app`, `g`):** Terikat pada siklus hidup aplikasi. Objek `g` digunakan untuk menyimpan data global per-request (misal: koneksi DB atau identitas user).
- **Request Context (`request`, `session`):** Terikat pada siklus satu request HTTP tertentu dari klien.

## 2. Praktik dasar: Application Factory Pattern dan Blueprints

Pola *Application Factory* memungkinkan pembuatan instance aplikasi Flask secara dinamis (sangat memudahkan unit testing dan multi-environment configuration):

```python
# app/__init__.py
from flask import Flask
from flask_sqlalchemy import SQLAlchemy

db = SQLAlchemy()

def create_app(config_name="development"):
    app = Flask(__name__)
    
    # Konfigurasi dasar
    app.config["SQLALCHEMY_DATABASE_URI"] = "sqlite:///learnhub.db"
    app.config["SQLALCHEMY_TRACK_MODIFICATIONS"] = False
    app.config["SECRET_KEY"] = "super-secret-key-change-in-production"

    # Inisialisasi ekstensi
    db.init_app(app)

    # Registrasi Blueprints
    from app.courses.routes import courses_bp
    app.register_blueprint(courses_bp, url_prefix="/api/v1/courses")

    return app
```

```python
# app/courses/routes.py
from flask import Blueprint, jsonify, request, abort
from app import db

courses_bp = Blueprint("courses", __name__)

# Model sederhana
class Course(db.Model):
    __tablename__ = "courses"
    id = db.Column(db.Integer, primary_key=True)
    title = db.Column(db.String(150), nullable=False)
    price = db.Column(db.Float, nullable=False)

    def to_dict(self):
        return {"id": self.id, "title": self.title, "price": self.price}

@courses_bp.route("", methods=["GET"])
def list_courses():
    courses = Course.query.all()
    return jsonify([c.to_dict() for c in courses]), 200

@courses_bp.route("", methods=["POST"])
def create_course():
    data = request.get_json()
    if not data or "title" not in data or "price" not in data:
        abort(400, description="Payload JSON wajib menyertakan 'title' dan 'price'")

    new_course = Course(title=data["title"], price=float(data["price"]))
    db.session.add(new_course)
    db.session.commit()

    return jsonify({"success": True, "data": new_course.to_dict()}), 201
```

## 3. Fitur utama dan pola arsitektur

### Decorators Kustom untuk Autentikasi API
Membuat decorator Python untuk melindungi endpoint:

```python
from functools import wraps
from flask import request, jsonify

def token_required(f):
    @wraps(f)
    def decorated(*args, **kwargs):
        token = request.headers.get("X-API-TOKEN")
        if not token or token != "valid-learnhub-secret":
            return jsonify({"error": "Akses ditolak: token tidak valid"}), 401
        return f(*args, **kwargs)
    return decorated

@courses_bp.route("/protected-module", methods=["GET"])
@token_required
def get_protected_data():
    return jsonify({"data": "Materi eksklusif instruktur"})
```

### Global Error Handling
Mengembalikan format respons error standar JSON alih-alih halaman HTML default:

```python
@app.errorhandler(404)
def handle_not_found(e):
    return jsonify({"error": "Resource tidak ditemukan", "code": 404}), 404

@app.errorhandler(500)
def handle_server_error(e):
    return jsonify({"error": "Terjadi kesalahan internal server", "code": 500}), 500
```

## 4. Praktik kerja: mencegah thread-safety leak pada objek `g`

### Skenario
Developer menyimpan data sesi user ke variabel modul global tingkat atas (`user_session = {}`). Ketika dua pengguna mengakses server Flask yang berjalan di mode multi-threaded, data transaksi antar-pengguna saling tertukar (*state contamination*).

### Tugas
Pindahkan state per-request ke dalam objek thread-safe bawaan Flask yaitu `flask.g`, dan pastikan pembersihan dilakukan via teardown hook.

```python
from flask import g, request

@app.before_request
def authenticate_and_set_user():
    # Ambil header atau cookie
    api_key = request.headers.get("X-API-KEY")
    if api_key:
        # Objek g terisolasi untuk request thread saat ini saja
        g.current_user = {"id": "usr_42", "role": "student"}
    else:
        g.current_user = None

@app.teardown_appcontext
def cleanup_resources(exception=None):
    # Dijalankan otomatis saat request selesai
    db_session = g.pop("db_session", None)
    if db_session is not None:
        db_session.close()
```

### Lima hint bertahap
1. Jangan pernah menggunakan variabel global Python biasa untuk menyimpan status request atau identitas user.
2. `flask.g` menggunakan context local (di belakang layar menggunakan `contextvars` Python) yang terisolasi per request.
3. `@app.before_request` dijalankan sebelum handler rute utama menerima request.
4. `@app.teardown_appcontext` selalu dipanggil bahkan jika terjadi error/exception saat eksekusi handler.
5. Objek `g` otomatis di-reset saat request berakhir sehingga tidak akan bocor ke request berikutnya.

## 5. Lanjutan: Deployment WSGI dengan Gunicorn dan Nginx

- **WSGI vs Development Server:** Jangan pernah menjalankan `app.run()` di server produksi. Gunakan server WSGI seperti Gunicorn:
  ```bash
  gunicorn -w 4 -b 127.0.0.1:8000 "app:create_app('production')"
  ```
- **Database Migrations:** Gunakan Flask-Migrate (berbasis Alembic) untuk versioning skema basis data secara otomatis.
- **Testing:** Manfaatkan `app.test_client()` untuk menulis unit test dan integration test end-to-end tanpa perlu menyalakan server HTTP fisik.

## Pemeriksaan pemahaman

1. Mengapa Application Factory Pattern lebih direkomendasikan dibanding instance `app = Flask(__name__)` global?
2. Apa perbedaan fungsi antara `current_app` dan `g` di Flask?
3. Mengapa server development bawaan Flask (`flask run`) tidak boleh digunakan di production?
4. Bagaimana Blueprint membantu dalam modularitas aplikasi Flask berskala besar?

**Kunci ringkas:** Application Factory mempermudah pengujian dengan konfigurasi berbeda tanpa konflik state instance global; `current_app` mengacu pada instance aplikasi yang sedang aktif sedangkan `g` menyimpan variabel sementara yang terisolasi per-request; server bawaan tidak dioptimalkan untuk beban konkuren tinggi dan tidak memiliki proteksi keamanan/timeout sekuat WSGI container (seperti Gunicorn); Blueprint memecah routing, view, dan template ke modul-modul independen yang dapat didaftarkan dengan prefix URL terpisah.

## Referensi

- [Flask Official Documentation](https://flask.palletsprojects.com/)
- [Flask Mega-Tutorial by Miguel Grinberg](https://blog.miguelgrinberg.com/post/the-flask-mega-tutorial-part-i-hello-world)
- [Gunicorn WSGI Server Documentation](https://gunicorn.org/)
