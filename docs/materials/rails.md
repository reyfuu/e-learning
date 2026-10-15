# Materi: Ruby on Rails untuk Pengembangan Web Penuh Konvensi

**Format:** teks · **Level:** Basic–Advanced · **Estimasi:** 60–90 menit

## Tujuan belajar

Setelah materi ini, pelajar dapat membangun aplikasi web dan API skala penuh menggunakan Ruby on Rails (v7/v8): memahami filosofi *Convention over Configuration* dan *Don't Repeat Yourself (DRY)*, menguasai Active Record ORM (relasi, validasi, migrasi, scopes), membangun controller dan route RESTful, mengelola background jobs dengan Sidekiq/SolidQueue, serta menerapkan caching dan keamanan bawaan Rails.

## Prasyarat

- Memahami bahasa pemrograman Ruby (OOP, blocks, modules, exceptions).
- Memahami konsep basis data relasional (PostgreSQL/SQLite) dan protokol HTTP.
- Ruby 3.2+ dan Rails 7+ terpasang di sistem.

## 1. Dasar: filosofi "The Rails Way"

Ruby on Rails adalah framework web pelopor konsep **Convention over Configuration (CoC)**. Rails mengasumsikan keputusan desain terbaik (nama tabel jamak, nama model tunggal, struktur direktori tetap, routing standar) sehingga developer tidak perlu menulis konfigurasi XML atau boilerplate berulang.

### Tiga Pilar Rails:
- **Active Record:** Pola arsitektur di mana satu kelas Model mewakili satu tabel database, dan setiap instance objek mewakili satu baris data (*row*).
- **Action Controller & Routing:** Memetakan verb HTTP ke tujuh aksi REST standar (`index`, `show`, `new`, `create`, `edit`, `update`, `destroy`).
- **Action View / Hotwire:** Lapisan presentasi modern menggunakan Turbo Streams dan Stimulus JS untuk interaktivitas real-time tanpa SPA yang kompleks.

## 2. Praktik dasar: Model, Migrasi, dan REST API Controller

```ruby
# db/migrate/20261016000000_create_courses.rb
class CreateCourses < ActiveRecord::Migration[7.1]
  def change
    create_table :courses do |t|
      t.string :title, null: false
      t.text :description
      t.decimal :price, precision: 10, scale: 2, null: false, default: 0.0
      t.references :instructor, null: false, foreign_key: { to_table: :users }
      t.boolean :published, default: false

      t.timestamps
    end

    add_index :courses, :title
  end
end
```

```ruby
# app/models/course.rb
class Course < ApplicationRecord
  belongs_to :instructor, class_name: 'User'
  has_many :enrollments, dependent: :destroy
  has_many :students, through: :enrollments, source: :user

  # Validasi
  validates :title, presence: true, length: { minimum: 5, maximum: 150 }
  validates :price, numericality: { greater_than_or_equal_to: 0 }

  # Scopes untuk query reusable yang bersih
  scope :published, -> { where(published: true) }
  scope :recent, -> { order(created_at: :desc) }
end
```

```ruby
# app/controllers/api/v1/courses_controller.rb
module Api
  module V1
    class CoursesController < ActionController::API
      before_action :set_course, only: %i[show update destroy]

      # GET /api/v1/courses
      def index
        # Eager loading untuk mencegah masalah N+1
        @courses = Course.includes(:instructor).published.recent.limit(20)
        render json: @courses, include: { instructor: { only: %i[id name email] } }
      end

      # POST /api/v1/courses
      def create
        @course = Course.new(course_params)
        if @course.save
          render json: @course, status: :created
        else
          render json: { errors: @course.errors.full_messages }, status: :unprocessable_entity
        end
      end

      # GET /api/v1/courses/:id
      def show
        render json: @course
      end

      private

      def set_course
        @course = Course.find(params[:id])
      rescue ActiveRecord::RecordNotFound
        render json: { error: 'Kursus tidak ditemukan' }, status: :not_found
      end

      # Strong Parameters untuk mencegah mass-assignment vulnerability
      def course_params
        params.require(:course).permit(:title, :description, :price, :instructor_id, :published)
      end
    end
  end
end
```

## 3. Fitur utama dan pola arsitektur

### Strong Parameters untuk Keamanan Mass-Assignment
Rails secara ketat memblokir atribut yang tidak didaftarkan secara eksplisit untuk mencegah pengguna memanipulasi field sensitif (seperti `is_admin`):

```ruby
def user_params
  params.require(:user).permit(:name, :email, :password) # is_admin tidak diizinkan di sini!
end
```

### Background Processing dengan Active Job
Rails menyediakan abstraksi Active Job terpadu yang dapat dihubungkan ke backend queue seperti SolidQueue atau Sidekiq:

```ruby
# app/jobs/enrollment_notification_job.rb
class EnrollmentNotificationJob < ApplicationJob
  queue_as :default

  def perform(enrollment_id)
    enrollment = Enrollment.find(enrollment_id)
    CourseMailer.welcome_student(enrollment).deliver_now
  end
end

# Panggil di controller:
# EnrollmentNotificationJob.perform_later(@enrollment.id)
```

## 4. Praktik kerja: mendeteksi dan menyelesaikan query N+1 dengan Bullet Gem

### Skenario
Halaman beranda katalog LearnHub memuat 30 kursus beserta nama instrukturnya. Response time lambat karena Active Record mengeksekusi 1 query untuk mengambil daftar kursus dan 30 query tambahan untuk setiap instruktur.

### Tugas
Gunakan method `.includes()` pada Active Record untuk menerapkan *eager loading*, dan konfigurasikan gem Bullet di lingkungan development untuk memberi peringatan otomatis jika ada query N+1 yang tidak disengaja.

```ruby
# Di Gemfile (group :development)
# gem 'bullet'

# config/environments/development.rb
config.after_initialize do
  Bullet.enable = true
  Bullet.alert = true
  Bullet.bullet_logger = true
  Bullet.console = true
end

# Solusi pada controller:
# Sebelum (Buruk): Course.all
# Sesudah (Optimal):
@courses = Course.includes(:instructor, :category).all
```

### Lima hint bertahap
1. Periksa log server Rails (`log/development.log`) untuk melihat baris SQL yang berulang-ulang dieksekusi.
2. Gunakan `includes(:relation)` untuk membiarkan Rails memilih antara query SQL JOIN atau subquery `WHERE IN`.
3. Gunakan `preload(:relation)` jika ingin memaksakan dua query terpisah tanpa JOIN.
4. Gunakan `eager_load(:relation)` jika ingin memaksakan SQL `LEFT OUTER JOIN`.
5. Bullet gem akan menampilkan popup alert di browser atau terminal jika mendeteksi pemanggilan relasi tanpa eager loading.

## 5. Lanjutan: Caching, Turbo/Hotwire, dan Production Deployment

- **Russian Doll Caching:** Memanfaatkan caching fragment bersarang dengan key digest otomatis (`<% cache @course do %>`).
- **Turbo Drive & Frames:** Mempercepat navigasi web konvensional menjadi secepat SPA tanpa harus menulis kode JavaScript manual.
- **Puma Server:** Web server multi-threaded dan multi-process bawaan Rails yang sangat optimal untuk deployment container Docker.

## Pemeriksaan pemahaman

1. Apa arti prinsip "Convention over Configuration" pada Ruby on Rails?
2. Mengapa Strong Parameters sangat krusial dalam controller Rails?
3. Apa perbedaan fungsi antara `includes`, `preload`, dan `eager_load` di Active Record?
4. Bagaimana Active Job memfasilitasi penggantian antrean background worker tanpa mengubah kode aplikasi?

**Kunci ringkas:** CoC menetapkan standar default penamaan dan struktur file sehingga pengembang fokus pada logika bisnis tanpa menulis konfigurasi berulang; Strong Parameters memfilter atribut input berbahaya dari serangan mass-assignment; `preload` selalu menggunakan 2 query terpisah, `eager_load` selalu memakai satu query LEFT JOIN, dan `includes` secara cerdas memilih di antara keduanya; Active Job adalah adapter layer yang memisahkan definisi pekerjaan dari implementasi backend spesifik (Sidekiq, SolidQueue, SQS).

## Referensi

- [Ruby on Rails Official Guides](https://guides.rubyonrails.org/)
- [Active Record Querying Guide](https://guides.rubyonrails.org/active_record_querying.html)
- [Hotwire Rails Documentation](https://hotwired.dev/)
