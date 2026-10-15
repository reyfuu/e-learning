# Materi: Bahasa Pemrograman Ruby dan Idiom Berorientasi Objek

**Format:** teks · **Level:** Basic–Advanced · **Estimasi:** 60–90 menit

## Tujuan belajar

Setelah materi ini, pelajar dapat menulis kode Ruby idiomatis: memahami prinsip "Everything is an Object", menguasai Blocks, Procs, dan Lambdas, memanfaatkan modul Mixin untuk komposisi perilaku, mengelola dependensi dengan Bundler/Gems, memahami dasar metaprogramming yang aman, serta menerapkan penanganan error dan konvensi penamaan kode Ruby.

## Prasyarat

- Memahami konsep dasar logika pemrograman (variabel, percabangan, fungsi).
- Memiliki Ruby 3.2+ terpasang di sistem (`ruby -v`).
- Memahami konsep dasar Object-Oriented Programming (OOP).

## 1. Dasar: filosofi Ruby dan prinsip kebahagiaan programmer

Ruby diciptakan oleh Yukihiro Matsumoto ("Matz") dengan tujuan utama membuat pemrogram bahagia dan produktif. Sintaksnya menyerupai bahasa manusia alami dan fleksibel. Di Ruby, **semua hal adalah objek murni**—bahkan angka primitif (`5.times`), boolean, dan nil adalah instance dari kelas.

### Fitur Kunci:
- **Pure OOP:** Tidak ada tipe primitif murni; semuanya mewarisi kelas `Object` (atau `BasicObject`).
- **Blocks & Closures:** Kemampuan meneruskan blok kode ke method menggunakan keyword `yield` atau objek `Proc` / `lambda`.
- **Duck Typing:** "If it walks like a duck and quacks like a duck, it's a duck." Fokus pada method yang dimiliki objek, bukan hirarki pewarisan formal.
- **Mixins via Modules:** Pengganti pewarisan ganda (*multiple inheritance*) melalui `include` dan `extend`.

## 2. Praktik dasar: Class, Modules, Blocks, dan Enumerable

```ruby
# frozen_string_literal: true

# Modul Mixin untuk menambahkan kemampuan audit
module Auditable
  def audit_log(action)
    puts "[AUDIT - #{Time.now.utc}] #{self.class} (ID: #{id}) -> #{action}"
  end
end

class Course
  include Auditable # Menambahkan method instance dari modul

  attr_reader :id, :title
  attr_accessor :price, :published

  def __init__; end # Bukan python! Di Ruby gunakan initialize

  def initialize(id:, title:, price:)
    @id = id
    @title = title
    @price = price.to_f
    @published = false
  end

  def publish!
    @published = true
    audit_log("Kursus berhasil dipublikasikan")
  end

  # Method yang menerima block dengan yield
  def with_discount(discount_percent)
    raise ArgumentError, "Diskon tidak valid" unless (0..100).cover?(discount_percent)

    discounted_price = @price - (@price * (discount_percent / 100.0))
    yield(discounted_price) if block_given?
  end
end

# Eksekusi praktik
course = Course.new(id: "rb_101", title: "Pemrograman Ruby Modern", price: 300_000)
course.publish!

# Memanfaatkan block
course.with_discount(15) do |new_price|
  puts "Harga promo khusus: Rp #{new_price.to_i}"
end

# Pemanfaatan Enumerable
courses = [
  Course.new(id: "1", title: "Ruby", price: 200_000),
  Course.new(id: "2", title: "Rails", price: 400_000),
  Course.new(id: "3", title: "Sinatra", price: 150_000)
]

total_value = courses.map(&:price).reduce(0, :+)
puts "Total nilai katalog: Rp #{total_value.to_i}"
```

## 3. Fitur utama dan pola arsitektur

### Perbedaan Block, Proc, dan Lambda
Ketiganya adalah closure di Ruby, namun memiliki perbedaan penting:
- **Block:** Bukan objek kelas satu; hanya dapat diteruskan satu per method (`def foo(&block)` atau `yield`).
- **Proc:** Objek (`Proc.new { ... }`), tidak memeriksa jumlah argumen dengan ketat, dan `return` di dalamnya akan keluar dari method yang memanggilnya.
- **Lambda:** Dibuat dengan `->(x) { ... }`, memeriksa jumlah argumen dengan ketat, dan `return` hanya mengembalikan nilai dari lambda tersebut.

```ruby
# Lambda memeriksa argumen secara ketat
strict_formatter = ->(name, score) { "#{name.upcase}: #{score} pts" }
puts strict_formatter.call("Budi", 95) # Berhasil
# strict_formatter.call("Budi") # ArgumentError!
```

### Metaprogramming Dasar dengan `define_method` dan `method_missing`
Ruby memungkinkan modifikasi kelas dan pembuatan method saat runtime:

```ruby
class UserRoleManager
  ROLES = %w[admin instructor student].freeze

  ROLES.each do |role_name|
    # Membuat method admin?, instructor?, student? secara dinamis
    define_method("#{role_name}?") do
      @role == role_name
    end
  end

  def initialize(role)
    @role = role
  end
end
```

## 4. Praktik kerja: penanganan memory leak pada gem dan monkey patching

### Skenario
Aplikasi mengalami kegagalan fungsi login setelah menginstal library pihak ketiga yang melakukan *monkey patching* (menimpa) method `String#strip` bawaan Ruby secara ceroboh, mengubah perilaku pembersihan input form.

### Tugas
Terapkan fitur **Refinements** modern Ruby untuk membatasi cakupan perubahan method hanya pada modul lokal tertentu tanpa mencemari namespace global.

```ruby
# frozen_string_literal: true

module StringSanitizerExtension
  refine String do
    def sanitize_input
      strip.gsub(/<script.*?>.*?<\/script>/i, "")
    end
  end
end

class RegistrationService
  # Hanya aktifkan monkey patch di dalam kelas ini saja
  using StringSanitizerExtension

  def process(raw_username)
    # Method sanitize_input hanya ada di dalam scope kelas ini
    clean_username = raw_username.sanitize_input
    puts "Username bersih: #{clean_username}"
  end
end
```

### Lima hint bertahap
1. Monkey patching langsung pada kelas inti (`class String; def ...; end`) berbahaya karena berefek global ke semua gem lain.
2. Fitur `refine Module do ... end` memungkinkan modifikasi kelas yang terisolasi.
3. Keyword `using ExtensionModule` mengaktifkan modifikasi hanya di dalam file atau kelas yang memanggilnya.
4. Gunakan `frozen_string_literal: true` di baris pertama file untuk menghindari alokasi memori string berulang.
5. Gunakan bundler-audit untuk memeriksa kerentanan keamanan pada daftar gem di `Gemfile.lock`.

## 5. Lanjutan: Garbage Collection, YJIT, dan Concurrency (Ractor)

- **YJIT (Yet Another Ruby JIT):** JIT compiler bawaan Ruby 3+ yang mempercepat aplikasi web Rails/Sinatra hingga 20-40%. Aktifkan dengan flag `--yjit`.
- **Ractors (Ruby 3 Concurrency):** Model konkurensi berbasis actor tanpa memory-sharing yang dapat mengeksekusi thread paralel murni tanpa terhalang GVL (Global VM Lock).
- **Bundler:** Manajer paket standar untuk mengunci versi dependency dan mengisolasi environment aplikasi (`bundle exec`).

## Pemeriksaan pemahaman

1. Mengapa di Ruby dikatakan "Everything is an Object"?
2. Apa perbedaan mendasar antara `Proc` dan `Lambda`?
3. Apa keunggulan `Refinements` dibanding monkey patching konvensional?
4. Apa manfaat menyertakan komentar `# frozen_string_literal: true` di awal file Ruby?

**Kunci ringkas:** Semua tipe data termasuk integer dan method adalah turunan kelas objek dan memiliki method; Lambda memvalidasi jumlah argumen secara ketat dan pernyataan `return` hanya keluar dari scope lambda itu sendiri, sedangkan Proc tidak ketat dan return-nya melompat dari method induk; Refinements membatasi perubahan method pada file/kelas tertentu sehingga tidak merusak pustaka lain di scope global; `frozen_string_literal: true` membuat semua literal string menjadi immutable sehingga menghemat alokasi memori objek string baru pada setiap eksekusi baris kode.

## Referensi

- [Official Ruby Programming Documentation](https://www.ruby-lang.org/en/documentation/)
- [Ruby Style Guide by Bozhidar Batsov](https://rubystyle.guide/)
- [Ruby Monstas Primer](https://rubymonstas.org/)
