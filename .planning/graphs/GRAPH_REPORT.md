# Graph Report - learn-2  (2026-10-16)

## Corpus Check
- 115 files · ~279,915 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 1293 nodes · 1510 edges · 87 communities (80 shown, 7 thin omitted)
- Extraction: 100% EXTRACTED · 0% INFERRED · 0% AMBIGUOUS
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `a06c4752`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- time.Time
- useAuthStore
- BRD — LearnHub
- 📚 Product Requirements Document (PRD)
- Learn Platform Backend
- dependencies
- context.Context
- compilerOptions
- User
- devDependencies
- cn
- AuthService
- PRD — LearnHub: Katalog Silabus Teknologi
- main
- CourseHandler
- scripts
- web/README.md
- AI Agent Guidelines — LearnHub
- 📐 System Design Document — LearnHub
- Claude System Instructions — LearnHub
- app/layout.tsx
- Setup LearnHub
- run.sh
- eslint.config.mjs
- next.config.ts
- postcss.config.mjs
- RTK.md
- setup-postgres.sh
- learn-platform/backend
- Communities (87 total, 7 thin omitted)
- TRD — LearnHub
- 🔒 Security Features
- 🔑 API Endpoints
- 📊 Key Implementation Details
- 📦 Dependencies
- 🐛 Common Issues
- 🗄️ Database
- 🚀 Deployment
- 🚀 Quick Start
- CourseService
- Materi: AWS Cloud Architecture dan Rekayasa Layanan Terkelola
- Materi: Nginx untuk Reverse Proxy dan Deployment Web
- Materi: JavaScript Modern (ES6+ dan Runtime Browser/Node.js)
- Materi: MongoDB dan Arsitektur Basis Data Berorientasi Dokumen
- Materi: AdonisJS untuk TypeScript-First Full-Stack MVC
- Materi: Modern C++ (C++17/20) dan Rekayasa Sistem Berkinerja Tinggi
- Materi: C# 12 dan Pemrograman Berorientasi Objek Modern
- Materi: Django Framework dan Django REST Framework (DRF)
- Materi: ASP.NET Core 8 Web API dan Ekosistem .NET Modern
- Materi: Echo Framework untuk REST API Cepat dan Ringan di Go
- Materi: Express.js untuk RESTful API dan Middleware Pipeline
- Materi: FastAPI untuk Backend Asynchronous Berkinerja Tinggi
- Materi: Firebase Platform dan Cloud Firestore NoSQL
- Materi: Flask Microframework untuk Layanan Web Python Fleksibel
- Materi: Flutter dan Dart untuk Aplikasi Multiplatform Modern
- Materi: Go (Golang) untuk Rekayasa Sistem Berkinerja Tinggi
- Materi: Java Modern (Java 17/21 LTS) dan Pemrograman Berorientasi Objek
- Materi: Laravel Framework untuk REST API dan Aplikasi Enterprise
- Materi: MariaDB Server dan Galera Cluster untuk High Availability
- Materi: MySQL 8.x dan Engine Transaksional InnoDB
- Materi: NestJS untuk Arsitektur Backend Skala Enterprise
- Materi: Oracle Database dan Rekayasa Prosedural PL/SQL Enterprise
- Materi: Modern PHP 8.x dan Rekayasa Backend Berbasis Standar
- Materi: PostgreSQL untuk Basis Data Relasional Skala Produksi
- Materi: Python 3 Modern dan Rekayasa Kode Bersih
- Materi: Ruby on Rails untuk Pengembangan Web Penuh Konvensi
- Materi: React Core, Hooks, dan Arsitektur Komponen
- Materi: Bahasa Pemrograman Ruby dan Idiom Berorientasi Objek
- Materi: Spring Boot 3 untuk Layanan REST Enterprise
- Materi: Supabase untuk Backend-as-a-Service Berbasis PostgreSQL
- Materi: Vue.js 3 dan Ekosistem Composition API
- Materi: Fiber Framework Berbasis Fasthttp untuk Throughput Maksimal
- Materi: Gin Web Framework untuk API Berkecepatan Tinggi di Go
- Materi: Next.js App Router dan Full-Stack Modern Web
- Graph Report - learn-2  (2026-10-16)
- Materi: TypeScript untuk Aplikasi Skala Produksi
- UserHandler
- materials/README.md
- LearnHub technology courses design
- FRD — LearnHub
- gorm.io/gorm.DB
- Review Focus
- 4. Praktik kerja: diagnosis `502 Bad Gateway`
- 4. Praktik kerja: implementasi Custom Auth Guard dan Roles Decorator

## God Nodes (most connected - your core abstractions)
1. `Communities (87 total, 7 thin omitted)` - 78 edges
2. `CourseRepository` - 24 edges
3. `Course` - 17 edges
4. `Learn Platform Backend` - 17 edges
5. `compilerOptions` - 16 edges
6. `CourseService` - 15 edges
7. `useAuthStore` - 15 edges
8. `User` - 14 edges
9. `UserRepository` - 14 edges
10. `main()` - 14 edges

## Surprising Connections (you probably didn't know these)
- `main()` --calls--> `NewAuthHandler()`  [EXTRACTED]
  apps/api/main.go → apps/api/internal/handler/auth_handler.go
- `CourseHandler` --references--> `CourseService`  [EXTRACTED]
  apps/api/internal/handler/course_handler.go → apps/api/internal/service/course_service.go
- `NewCourseHandler()` --references--> `CourseService`  [EXTRACTED]
  apps/api/internal/handler/course_handler.go → apps/api/internal/service/course_service.go
- `main()` --calls--> `NewCourseHandler()`  [EXTRACTED]
  apps/api/main.go → apps/api/internal/handler/course_handler.go
- `UserHandler` --references--> `UserService`  [EXTRACTED]
  apps/api/internal/handler/user_handler.go → apps/api/internal/service/user_service.go

## Import Cycles
- None detected.

## Communities (87 total, 7 thin omitted)

### Community 0 - "time.Time"
Cohesion: 0.10
Nodes (13): Course, User, Course, User, Course, User, gorm.io/gorm.DeletedAt, time.Time (+5 more)

### Community 1 - "useAuthStore"
Cohesion: 0.06
Nodes (33): AdminDashboard(), Course, Course, CourseDetailPage(), Module, AddCoursePage(), CoursesPage(), DashboardPage() (+25 more)

### Community 2 - "BRD — LearnHub"
Cohesion: 0.17
Nodes (12): 1. Ringkasan, 2. Tujuan bisnis, 3. Pemangku kepentingan dan kebutuhan, 4. Ruang lingkup bisnis, 5. Model operasional dan pendapatan, 6. Risiko dan asumsi bisnis, 7. Tahapan bisnis, 8. Kriteria keberhasilan (+4 more)

### Community 3 - "📚 Product Requirements Document (PRD)"
Cohesion: 0.05
Nodes (44): 10. Rencana Pengembangan (Milestone), 11. Asumsi & Risiko, 12. Open Questions untuk Review, 1. Ringkasan Produk, 2. Tujuan Bisnis & Metrik Keberhasilan, 3. User Persona, 4.1 Modul Autentikasi & Manajemen Pengguna, 4.2 Modul Manajemen Kursus (Instructor) (+36 more)

### Community 4 - "Learn Platform Backend"
Cohesion: 0.20
Nodes (10): 📚 Additional Resources, 🏗️ Architecture, 🛠️ Development Commands, Learn Platform Backend, 📈 Performance Optimization, 📁 Project Structure, 🔗 Related Documentation, 🔄 Request Flow (+2 more)

### Community 5 - "dependencies"
Cohesion: 0.06
Nodes (35): dependencies, axios, class-variance-authority, clsx, framer-motion, @hookform/resolvers, js-cookie, lucide-react (+27 more)

### Community 6 - "context.Context"
Cohesion: 0.13
Nodes (8): Course, User, Lesson, LessonProgress, User, Section, CourseRepository, context.Context

### Community 7 - "compilerOptions"
Cohesion: 0.07
Nodes (28): compilerOptions, allowJs, esModuleInterop, incremental, isolatedModules, jsx, lib, module (+20 more)

### Community 8 - "User"
Cohesion: 0.15
Nodes (7): UpdateUserRequest, User, UserRepository, NewUserRepository(), UserService, NewUserService(), CreateUserRequest

### Community 9 - "devDependencies"
Cohesion: 0.08
Nodes (25): devDependencies, eslint, eslint-config-next, tailwindcss, @tailwindcss/postcss, @types/node, @types/react, @types/react-dom (+17 more)

### Community 10 - "cn"
Cohesion: 0.12
Nodes (15): Sidebar(), SidebarProps, Button, ButtonProps, Badge(), BadgeProps, Card(), CardProps (+7 more)

### Community 11 - "AuthService"
Cohesion: 0.12
Nodes (12): echo.Context, NewAuthHandler(), AuthResponse, Claims, LoginRequest, RegisterRequest, AuthService, NewAuthService() (+4 more)

### Community 12 - "PRD — LearnHub: Katalog Silabus Teknologi"
Cohesion: 0.08
Nodes (24): 10. Kriteria penerimaan rilis awal, 11. Ketergantungan dan keputusan yang masih perlu ditetapkan, 1. Visi dan ringkasan produk, 2. Sasaran dan bukan sasaran, 3. Persona dan kebutuhan, 4. Taksonomi dan pola silabus, 5. Katalog silabus, 6. Persyaratan produk dan UX (+16 more)

### Community 13 - "main"
Cohesion: 0.10
Nodes (16): getEnv(), Load(), echo.Context, NewPracticeHandler(), firstNonEmpty(), PracticeService, NewPracticeService(), statusMessage() (+8 more)

### Community 14 - "CourseHandler"
Cohesion: 0.27
Nodes (3): echo.Context, NewCourseHandler(), CourseHandler

### Community 15 - "scripts"
Cohesion: 0.18
Nodes (10): name, packageManager, private, scripts, api:dev, build, dev, dev:all (+2 more)

### Community 16 - "web/README.md"
Cohesion: 0.20
Nodes (7): Deploy on Vercel, Getting Started, Learn More, Build dan deploy, LearnHub, Menjalankan lokal, Struktur

### Community 17 - "AI Agent Guidelines — LearnHub"
Cohesion: 0.22
Nodes (8): AI Agent Guidelines — LearnHub, Konteks Proyek, Pedoman Backend (Go), Pedoman Frontend (Next.js), Penyelesaian Sesi, Perilaku Umum, Tech Stack Utama, Tooling dan Workflow

### Community 18 - "📐 System Design Document — LearnHub"
Cohesion: 0.22
Nodes (8): 1. Arsitektur Sistem, 2. Tech Stack, 3. Skema Database (Utama), 4. Struktur API (REST), 5. Non-Functional Requirements (NFR), Backend, Frontend, 📐 System Design Document — LearnHub

### Community 19 - "Claude System Instructions — LearnHub"
Cohesion: 0.29
Nodes (6): 1. Project Context, 2. Technology Stack, 3. Frontend Coding Guidelines, 4. Backend Coding Guidelines, 5. Interaction Rules, Claude System Instructions — LearnHub

### Community 20 - "app/layout.tsx"
Cohesion: 0.40
Nodes (3): inter, metadata, Providers()

### Community 21 - "Setup LearnHub"
Cohesion: 0.33
Nodes (5): Instalasi dan development, Prasyarat, Setup LearnHub, Struktur, Vercel

### Community 31 - "Communities (87 total, 7 thin omitted)"
Cohesion: 0.03
Nodes (78): Communities (87 total, 7 thin omitted), Community 0 - "gorm.io/gorm.DB", Community 10 - "cn", Community 11 - "main", Community 12 - "PRD — LearnHub: Katalog Silabus Teknologi", Community 13 - "PracticeService", Community 14 - "CourseHandler", Community 15 - "scripts" (+70 more)

### Community 32 - "TRD — LearnHub"
Cohesion: 0.13
Nodes (15): 1. Prinsip arsitektur, 2. Topologi deployment, 3. Komponen aplikasi, 4. Kontrak API target, 5. Autentikasi dan keamanan, 6. Deploy, CI/CD, dan konfigurasi, 7. Observability dan reliability, 8. Batas performa dan kapasitas (+7 more)

### Community 33 - "🔒 Security Features"
Cohesion: 0.33
Nodes (6): Authentication, Authorization, Error Handling, Input Validation, Middleware, 🔒 Security Features

### Community 34 - "🔑 API Endpoints"
Cohesion: 0.40
Nodes (5): 🔑 API Endpoints, Authentication, Health Check, Latihan coding (protected), Users (Protected)

### Community 35 - "📊 Key Implementation Details"
Cohesion: 0.40
Nodes (5): Error Responses, JWT Tokens, 📊 Key Implementation Details, Password Security, User Roles

### Community 36 - "📦 Dependencies"
Cohesion: 0.50
Nodes (4): Authentication, Core, 📦 Dependencies, Security

### Community 37 - "🐛 Common Issues"
Cohesion: 0.50
Nodes (4): 🐛 Common Issues, Database Connection Error, JWT Secret Error, Port Already in Use

### Community 38 - "🗄️ Database"
Cohesion: 0.67
Nodes (3): Connection, 🗄️ Database, Schema

### Community 39 - "🚀 Deployment"
Cohesion: 0.67
Nodes (3): 🚀 Deployment, Docker, Environment Variables for Production

### Community 40 - "🚀 Quick Start"
Cohesion: 0.67
Nodes (3): Installation, Prerequisites, 🚀 Quick Start

### Community 41 - "CourseService"
Cohesion: 0.20
Nodes (8): CourseResponse, CreateCourseRequest, CreateModuleRequest, ModuleResponse, UpdateCourseStatusRequest, CourseService, NewCourseService(), slugify()

### Community 42 - "Materi: AWS Cloud Architecture dan Rekayasa Layanan Terkelola"
Cohesion: 0.11
Nodes (17): 1. Dasar: arsitektur global AWS dan Shared Responsibility Model, 2. Praktik dasar: IAM Least Privilege dan S3 Bucket Hardening, 3. Fitur utama dan pola arsitektur, 4. Praktik kerja: mitigasi insiden kebocoran kredensial root dan audit trail, 5. Lanjutan: Enkripsi KMS, Auto Scaling, dan Well-Architected Pillars, Desain Jaringan Virtual Private Cloud (VPC), Komputasi Serverless dengan AWS Lambda, Langkah Tanggap Darurat: (+9 more)

### Community 43 - "Materi: Nginx untuk Reverse Proxy dan Deployment Web"
Cohesion: 0.22
Nodes (9): 1. Dasar: apa yang dilakukan Nginx?, 2. Praktik dasar: teruskan request ke aplikasi, 3. File statis dan pemilihan path, 5. Lanjutan: HTTPS, load balancing, dan hardening, Materi: Nginx untuk Reverse Proxy dan Deployment Web, Pemeriksaan pemahaman, Prasyarat, Referensi (+1 more)

### Community 44 - "Materi: JavaScript Modern (ES6+ dan Runtime Browser/Node.js)"
Cohesion: 0.11
Nodes (17): 1. Dasar: ekosistem dan konsep inti JavaScript, 2. Praktik dasar: sintaks modern ES6+ dan Async/Await, 3. Fitur utama dan pola arsitektur, 4. Praktik kerja: optimasi memory leak dan debouncing input, 5. Lanjutan: performa, keamanan, dan best practices, Closures dan Prototypes, Event Loop: Microtasks vs Macrotasks, Konsep kunci: (+9 more)

### Community 45 - "Materi: MongoDB dan Arsitektur Basis Data Berorientasi Dokumen"
Cohesion: 0.11
Nodes (17): 1. Dasar: model data dokumen NoSQL dan BSON, 2. Praktik dasar: Operasi CRUD dan Document Schema, 3. Fitur utama dan pola arsitektur, 4. Praktik kerja: optimasi kueri dan diagnosa COLLSCAN dengan `explain()`, 5. Lanjutan: Replica Sets, Sharding, dan Transaksi Multi-Dokumen, Aggregation Pipeline, Embedding vs Referencing:, Lima hint bertahap (+9 more)

### Community 46 - "Materi: AdonisJS untuk TypeScript-First Full-Stack MVC"
Cohesion: 0.12
Nodes (16): 1. Dasar: filosofi AdonisJS v6, 2. Praktik dasar: Lucid Model, VineJS Validator, dan Controller, 3. Fitur utama dan pola arsitektur, 4. Praktik kerja: migrasi skema dan pagination pada dataset besar, 5. Lanjutan: event emitter, queue, dan health check, Karakteristik Kunci:, Lima hint bertahap, Materi: AdonisJS untuk TypeScript-First Full-Stack MVC (+8 more)

### Community 47 - "Materi: Modern C++ (C++17/20) dan Rekayasa Sistem Berkinerja Tinggi"
Cohesion: 0.12
Nodes (16): 1. Dasar: filosofi "Zero-Overhead Abstraction" dan RAII, 2. Praktik dasar: RAII, Smart Pointers, dan STL Containers, 3. Fitur utama dan pola arsitektur, 4. Praktik kerja: mendiagnosis dan mencegah Segmentation Fault serta Memory Leak, 5. Lanjutan: C++20 Concepts, Coroutines, dan Concurrency, Konsep Kunci:, Lima hint bertahap, Materi: Modern C++ (C++17/20) dan Rekayasa Sistem Berkinerja Tinggi (+8 more)

### Community 48 - "Materi: C# 12 dan Pemrograman Berorientasi Objek Modern"
Cohesion: 0.12
Nodes (16): 1. Dasar: arsitektur Common Language Runtime (CLR), 2. Praktik dasar: Record, Primary Constructor, dan LINQ, 3. Fitur utama dan pola arsitektur, 4. Praktik kerja: optimasi alokasi string dengan `ReadOnlySpan<char>`, 5. Lanjutan: Asynchronous Task vs ValueTask, dan GC Generations, Fitur Kunci C# 12:, Lima hint bertahap, Materi: C# 12 dan Pemrograman Berorientasi Objek Modern (+8 more)

### Community 49 - "Materi: Django Framework dan Django REST Framework (DRF)"
Cohesion: 0.12
Nodes (16): 1. Dasar: arsitektur MVT dan "Batteries-Included", 2. Praktik dasar: Model, Serializer, dan DRF ViewSet, 3. Fitur utama dan pola arsitektur, 4. Praktik kerja: optimasi transaksi keuangan dengan Database Atomic Transaction, 5. Lanjutan: Celery untuk Background Tasks, Caching, dan Security Checklist, Custom Permissions di DRF, Lima hint bertahap, Materi: Django Framework dan Django REST Framework (DRF) (+8 more)

### Community 50 - "Materi: ASP.NET Core 8 Web API dan Ekosistem .NET Modern"
Cohesion: 0.12
Nodes (16): 1. Dasar: arsitektur ASP.NET Core dan HTTP Pipeline, 2. Praktik dasar: Minimal APIs, EF Core, dan DTO Validation, 3. Fitur utama dan pola arsitektur, 4. Praktik kerja: implementasi Global Exception Handling Middleware, 5. Lanjutan: Health Checks, Rate Limiting, dan Kestrel Tuning, Custom Middleware Pipeline, Lima hint bertahap, Materi: ASP.NET Core 8 Web API dan Ekosistem .NET Modern (+8 more)

### Community 51 - "Materi: Echo Framework untuk REST API Cepat dan Ringan di Go"
Cohesion: 0.12
Nodes (16): 1. Dasar: arsitektur Echo Framework, 2. Praktik dasar: setup Server, Struct Binding, dan Validator, 3. Fitur utama dan pola arsitektur, 4. Praktik kerja: implementasi Graceful Shutdown pada Echo Server, 5. Lanjutan: profiling, Sub-domain routing, dan template rendering, Centralized HTTP Error Handler, Custom Middleware di Echo, Keunggulan Echo: (+8 more)

### Community 52 - "Materi: Express.js untuk RESTful API dan Middleware Pipeline"
Cohesion: 0.12
Nodes (16): 1. Dasar: siklus request-response dan middleware, 2. Praktik dasar: setup Express Server dengan Zod Validation, 3. Fitur utama dan pola arsitektur, 4. Praktik kerja: implementasi JWT Authentication & Role-Based Access Control (RBAC), 5. Lanjutan: rate limiting, compression, dan graceful shutdown, Async Error Handling Tanpa Try-Catch Berulang, Karakteristik Middleware:, Lima hint bertahap (+8 more)

### Community 53 - "Materi: FastAPI untuk Backend Asynchronous Berkinerja Tinggi"
Cohesion: 0.12
Nodes (16): 1. Dasar: arsitektur ASGI dan ekosistem FastAPI, 2. Praktik dasar: Pydantic Schema, Endpoints, dan Dependency Injection, 3. Fitur utama dan pola arsitektur, 4. Praktik kerja: integrasi database asynchronous dengan SQLAlchemy 2.0 & Asyncpg, 5. Lanjutan: OAuth2 Password Bearer, CORS, dan Uvicorn Tuning, Async vs Def (Sync) di FastAPI, Background Tasks, Keunggulan Utama: (+8 more)

### Community 54 - "Materi: Firebase Platform dan Cloud Firestore NoSQL"
Cohesion: 0.12
Nodes (16): 1. Dasar: ekosistem Firebase dan Cloud Firestore, 2. Praktik dasar: Firestore CRUD dan Realtime Listeners (SDK v9+ Modular), 3. Fitur utama dan pola arsitektur, 4. Praktik kerja: optimasi tagihan dan mencegah infinite read loop, 5. Lanjutan: Firestore Indexes, App Check, dan Backup Otomatis, Cloud Functions untuk Operasi Terproteksi, Firebase Security Rules: Menjaga Keamanan Langsung dari Klien, Komponen Utama Firebase: (+8 more)

### Community 55 - "Materi: Flask Microframework untuk Layanan Web Python Fleksibel"
Cohesion: 0.12
Nodes (16): 1. Dasar: arsitektur Flask dan filosofi microframework, 2. Praktik dasar: Application Factory Pattern dan Blueprints, 3. Fitur utama dan pola arsitektur, 4. Praktik kerja: mencegah thread-safety leak pada objek `g`, 5. Lanjutan: Deployment WSGI dengan Gunicorn dan Nginx, Decorators Kustom untuk Autentikasi API, Global Error Handling, Konsep Konteks di Flask: (+8 more)

### Community 56 - "Materi: Flutter dan Dart untuk Aplikasi Multiplatform Modern"
Cohesion: 0.12
Nodes (16): 1. Dasar: arsitektur engine Flutter dan bahasa Dart, 2. Praktik dasar: StatelessWidget, StatefulWidget, dan HTTP Fetching, 3. Fitur utama dan pola arsitektur, 4. Praktik kerja: optimasi jank dan rebuilding widget berlebih, 5. Lanjutan: Platform Channels, Isolates, dan Profiling, Konsep Inti Flutter & Dart:, Lima hint bertahap, Materi: Flutter dan Dart untuk Aplikasi Multiplatform Modern (+8 more)

### Community 57 - "Materi: Go (Golang) untuk Rekayasa Sistem Berkinerja Tinggi"
Cohesion: 0.12
Nodes (16): 1. Dasar: filosofi kesederhanaan Go, 2. Praktik dasar: Struct, Interface, Error Handling, dan Concurrency, 3. Fitur utama dan pola arsitektur, 4. Praktik kerja: mendeteksi dan mencegah Data Race dengan Race Detector, 5. Lanjutan: Garbage Collector, Profiling pprof, dan Escape Analysis, Context Cancellation dan Deadline, Fitur Kunci:, Lima hint bertahap (+8 more)

### Community 58 - "Materi: Java Modern (Java 17/21 LTS) dan Pemrograman Berorientasi Objek"
Cohesion: 0.12
Nodes (16): 1. Dasar: evolusi Java modern dan JVM, 2. Praktik dasar: Record, Sealed Interface, dan Stream API, 3. Fitur utama dan pola arsitektur, 4. Praktik kerja: konkurensi modern dengan Virtual Threads (Java 21), 5. Lanjutan: JVM Memory Management, GC Tuning, dan Profiling, Fitur Kunci Java 17 & 21:, Generics dan Collections Framework, Lima hint bertahap (+8 more)

### Community 59 - "Materi: Laravel Framework untuk REST API dan Aplikasi Enterprise"
Cohesion: 0.12
Nodes (16): 1. Dasar: arsitektur dan filosofi Laravel, 2. Praktik dasar: REST API Controller, Form Request, dan Eloquent, 3. Fitur utama dan pola arsitektur, 4. Praktik kerja: optimasi transaksi keuangan dengan Database Transaction & Pessimistic Locking, 5. Lanjutan: caching, API Resources, dan otentikasi Sanctum, Background Jobs dan Queues, Komponen Kunci:, Lima hint bertahap (+8 more)

### Community 60 - "Materi: MariaDB Server dan Galera Cluster untuk High Availability"
Cohesion: 0.12
Nodes (16): 1. Dasar: MariaDB sebagai evolusi mandiri open-source, 2. Praktik dasar: System-Versioned Tables dan Window Functions, 3. Fitur utama dan pola arsitektur, 4. Praktik kerja: migrasi data dan penanganan komparabilitas MySQL vs MariaDB, 5. Lanjutan: Galera Arbitrator, Optimizer Cost, dan Monitoring, Fitur Khas MariaDB:, Galera Cluster (Multi-Master Synchronous Replication), Lima hint bertahap (+8 more)

### Community 61 - "Materi: MySQL 8.x dan Engine Transaksional InnoDB"
Cohesion: 0.12
Nodes (16): 1. Dasar: arsitektur MySQL dan Storage Engine InnoDB, 2. Praktik dasar: Skema Relasional, Foreign Key, dan Transaction, 3. Fitur utama dan pola arsitektur, 4. Praktik kerja: menganalisis dan menyelesaikan Deadlock, 5. Lanjutan: Replikasi GTID, Binlog, dan Tuning Buffer Pool, Covering Index, Konsep Kunci InnoDB:, Lima hint bertahap (+8 more)

### Community 62 - "Materi: NestJS untuk Arsitektur Backend Skala Enterprise"
Cohesion: 0.17
Nodes (12): 1. Dasar: arsitektur modular NestJS, 2. Praktik dasar: DTO, Controller, dan Service dengan Dependency Injection, 3. Fitur utama dan pola arsitektur, 5. Lanjutan: Microservices, TypeORM/Prisma, dan Swagger OpenAPI, Alur Eksekusi Request Lifecycle di NestJS, Blok Pembangun Utama:, Dependency Injection Hierarkis, Materi: NestJS untuk Arsitektur Backend Skala Enterprise (+4 more)

### Community 63 - "Materi: Oracle Database dan Rekayasa Prosedural PL/SQL Enterprise"
Cohesion: 0.12
Nodes (16): 1. Dasar: arsitektur Oracle Database dan Instance, 2. Praktik dasar: Skema Relasional, Sequence, dan Paket PL/SQL, 3. Fitur utama dan pola arsitektur, 4. Praktik kerja: menyelesaikan error `ORA-01555: snapshot too old`, 5. Lanjutan: RAC (Real Application Clusters) dan Data Guard, Komponen Memori & Proses Oracle:, Lima hint bertahap, Materi: Oracle Database dan Rekayasa Prosedural PL/SQL Enterprise (+8 more)

### Community 64 - "Materi: Modern PHP 8.x dan Rekayasa Backend Berbasis Standar"
Cohesion: 0.12
Nodes (16): 1. Dasar: evolusi PHP modern, 2. Praktik dasar: OOP Modern, Validasi Tipe, dan Match, 3. Fitur utama dan pola arsitektur, 4. Praktik kerja: implementasi Dependency Injection dan custom Exception Handler, 5. Lanjutan: Composer autoloading, testing, dan performa, Fitur Kunci PHP 8+:, Interaksi Database Aman dengan PDO dan Prepared Statements, Lima hint bertahap (+8 more)

### Community 65 - "Materi: PostgreSQL untuk Basis Data Relasional Skala Produksi"
Cohesion: 0.12
Nodes (16): 1. Dasar: arsitektur PostgreSQL dan MVCC, 2. Praktik dasar: Skema Relasional, Constraints, dan JSONB, 3. Fitur utama dan pola arsitektur, 4. Praktik kerja: optimasi query lambat dan mitigasi Table Bloat dengan VACUUM, 5. Lanjutan: Partisi Tabel (Table Partitioning), Replikasi, dan PgBouncer, Analisis Query Plan dengan `EXPLAIN (ANALYZE, BUFFERS)`, Concurrency dan Tingkat Isolasi Transaksi, Karakteristik Utama: (+8 more)

### Community 66 - "Materi: Python 3 Modern dan Rekayasa Kode Bersih"
Cohesion: 0.12
Nodes (16): 1. Dasar: filosofi "Zen of Python" dan Data Model, 2. Praktik dasar: Type Hints, Dataclass, dan Generators, 3. Fitur utama dan pola arsitektur, 4. Praktik kerja: optimasi memory footprint dengan `__slots__`, 5. Lanjutan: GIL, Multiprocessing, dan Profiling CPython, Context Managers (`with` statement), Fitur Kunci Python Modern:, Lima hint bertahap (+8 more)

### Community 67 - "Materi: Ruby on Rails untuk Pengembangan Web Penuh Konvensi"
Cohesion: 0.12
Nodes (16): 1. Dasar: filosofi "The Rails Way", 2. Praktik dasar: Model, Migrasi, dan REST API Controller, 3. Fitur utama dan pola arsitektur, 4. Praktik kerja: mendeteksi dan menyelesaikan query N+1 dengan Bullet Gem, 5. Lanjutan: Caching, Turbo/Hotwire, dan Production Deployment, Background Processing dengan Active Job, Lima hint bertahap, Materi: Ruby on Rails untuk Pengembangan Web Penuh Konvensi (+8 more)

### Community 68 - "Materi: React Core, Hooks, dan Arsitektur Komponen"
Cohesion: 0.12
Nodes (16): 1. Dasar: model mental React dan Virtual DOM, 2. Praktik dasar: Hooks Esensial dan Form Handling, 3. Fitur utama dan pola arsitektur, 4. Praktik kerja: mencegah unnecessary re-renders dan cascading updates, 5. Lanjutan: Concurrent React, Suspense, dan Error Boundaries, Aturan Utama Hooks:, Context API vs External Store (Zustand), Lima hint bertahap (+8 more)

### Community 69 - "Materi: Bahasa Pemrograman Ruby dan Idiom Berorientasi Objek"
Cohesion: 0.12
Nodes (16): 1. Dasar: filosofi Ruby dan prinsip kebahagiaan programmer, 2. Praktik dasar: Class, Modules, Blocks, dan Enumerable, 3. Fitur utama dan pola arsitektur, 4. Praktik kerja: penanganan memory leak pada gem dan monkey patching, 5. Lanjutan: Garbage Collection, YJIT, dan Concurrency (Ractor), Fitur Kunci:, Lima hint bertahap, Materi: Bahasa Pemrograman Ruby dan Idiom Berorientasi Objek (+8 more)

### Community 70 - "Materi: Spring Boot 3 untuk Layanan REST Enterprise"
Cohesion: 0.12
Nodes (16): 1. Dasar: ekosistem Spring Framework dan Spring Boot, 2. Praktik dasar: REST Controller, DTO Validation, dan Spring Data JPA, 3. Fitur utama dan pola arsitektur, 4. Praktik kerja: optimasi hibernate query dan deteksi N+1 problem, 5. Lanjutan: Spring Security 6, JWT, dan Actuator, Anotasi Esensial:, Global Exception Handler dengan `@RestControllerAdvice`, Lima hint bertahap (+8 more)

### Community 71 - "Materi: Supabase untuk Backend-as-a-Service Berbasis PostgreSQL"
Cohesion: 0.12
Nodes (16): 1. Dasar: filosofi "Open Source Firebase Alternative", 2. Praktik dasar: Supabase Client SDK dan Row Level Security (RLS), 3. Fitur utama dan pola arsitektur, 4. Praktik kerja: mencegah kebocoran data akibat lupa mengaktifkan RLS, 5. Lanjutan: Supabase Edge Functions (Deno) dan Vector Search (`pgvector`), Database Triggers untuk Sinkronisasi `auth.users` ke `public.profiles`, Komponen Inti Supabase:, Lima hint bertahap (+8 more)

### Community 72 - "Materi: Vue.js 3 dan Ekosistem Composition API"
Cohesion: 0.12
Nodes (16): 1. Dasar: reaktivitas Vue 3 dan Composition API, 2. Praktik dasar: Single File Component (SFC) dengan `<script setup>`, 3. Fitur utama dan pola arsitektur, 4. Praktik kerja: optimasi form reaktif dan debounce pencarian, 5. Lanjutan: performa, teleports, dan testing, Konsep Kunci:, Lima hint bertahap, Materi: Vue.js 3 dan Ekosistem Composition API (+8 more)

### Community 73 - "Materi: Fiber Framework Berbasis Fasthttp untuk Throughput Maksimal"
Cohesion: 0.12
Nodes (15): 1. Dasar: arsitektur Fiber dan Fasthttp, 2. Praktik dasar: setup Fiber Server dan Routing, 3. Fitur utama dan pola arsitektur, 4. Praktik kerja: custom rate limiter middleware berbasis token bucket, 5. Lanjutan: WebSocket, Monitor dashboard, dan Fiber Prefork, Jebakan Memori Zero Allocation: `c.Immutable()` vs Salinan Buffer, Konsep Kunci:, Lima hint bertahap (+7 more)

### Community 74 - "Materi: Gin Web Framework untuk API Berkecepatan Tinggi di Go"
Cohesion: 0.12
Nodes (15): 1. Dasar: performa dan kesederhanaan Gin, 2. Praktik dasar: setup Router, Model Binding, dan Validation, 3. Fitur utama dan pola arsitektur, 4. Praktik kerja: penanganan context concurrency yang aman, 5. Lanjutan: Production Mode, Testing, dan Clean Architecture, Fitur Kunci:, Lima hint bertahap, Materi: Gin Web Framework untuk API Berkecepatan Tinggi di Go (+7 more)

### Community 75 - "Materi: Next.js App Router dan Full-Stack Modern Web"
Cohesion: 0.12
Nodes (15): 1. Dasar: arsitektur Next.js App Router, 2. Praktik dasar: Server Component Fetching & Server Actions, 3. Fitur utama dan pola arsitektur, 4. Praktik kerja: optimasi streaming UI dengan Suspense & Skeleton, 5. Lanjutan: Caching, Middleware, dan Metadata SEO, Batasan Penting:, Lima hint bertahap, Materi: Next.js App Router dan Full-Stack Modern Web (+7 more)

### Community 76 - "Graph Report - learn-2  (2026-10-16)"
Cohesion: 0.18
Nodes (10): Community Hubs (Navigation), Corpus Check, God Nodes (most connected - your core abstractions), Graph Freshness, Graph Report - learn-2  (2026-10-16), Import Cycles, Knowledge Gaps, Suggested Questions (+2 more)

### Community 77 - "Materi: TypeScript untuk Aplikasi Skala Produksi"
Cohesion: 0.12
Nodes (16): 1. Dasar: mengapa TypeScript?, 2. Praktik dasar: Type, Interface, dan Generics, 3. Fitur utama dan pola arsitektur, 4. Praktik kerja: validasi tipe API eksternal dan Type Guard, 5. Lanjutan: Conditional Types dan Mapped Types, Discriminated Unions (Tagged Unions), Konsep kunci:, Lima hint bertahap (+8 more)

### Community 78 - "UserHandler"
Cohesion: 0.39
Nodes (3): echo.Context, NewUserHandler(), UserHandler

### Community 79 - "materials/README.md"
Cohesion: 0.20
Nodes (7): A. Bahasa & Framework Frontend, B. Bahasa & Framework Backend, C. Mobile Development, D. Basis Data & Backend Platform, Daftar Materi Berdasarkan Domain, E. Cloud, Platform Engineering & DevOps, Katalog Materi Pembelajaran Teknis — LearnHub

### Community 80 - "LearnHub technology courses design"
Cohesion: 0.22
Nodes (8): Acceptance criteria, Current constraints, Data flow, Goal, LearnHub technology courses design, Recommended design, Safety and compatibility, Verification plan

### Community 81 - "FRD — LearnHub"
Cohesion: 0.17
Nodes (12): 1. Aktor dan hak akses, 2. Kebutuhan fungsional, 3. Alur inti, 4. Aturan bisnis, 5. Data dan validasi penting, 6. Status dan error, 7. Acceptance criteria ringkas, Belajar kursus (+4 more)

### Community 82 - "gorm.io/gorm.DB"
Cohesion: 0.15
Nodes (7): Course, Enrollment, User, AuthRepository, NewAuthRepository(), NewCourseRepository(), gorm.io/gorm.DB

### Community 83 - "Review Focus"
Cohesion: 0.22
Nodes (8): Global Constraints, LearnHub technology courses Implementation Plan, Review Focus, Task 1: Author the 21-course catalog, Task 2: Add safe explicit course seeding, Task 3: Return published catalog and readable lessons, Task 4: Render course descriptions and lesson text, Task 5: Connect and verify production publication

### Community 85 - "4. Praktik kerja: diagnosis `502 Bad Gateway`"
Cohesion: 0.40
Nodes (5): 4. Praktik kerja: diagnosis `502 Bad Gateway`, Langkah diagnosis, Lima hint bertahap, Skenario, Tugas

### Community 86 - "4. Praktik kerja: implementasi Custom Auth Guard dan Roles Decorator"
Cohesion: 0.50
Nodes (4): 4. Praktik kerja: implementasi Custom Auth Guard dan Roles Decorator, Lima hint bertahap, Skenario, Tugas

## Knowledge Gaps
- **773 isolated node(s):** `learn-platform/backend`, `RefreshTokenRequest`, `CreateUserRequest`, `runnerResult`, `eslintConfig` (+768 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **7 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `📚 Product Requirements Document (PRD)` connect `📚 Product Requirements Document (PRD)` to `PRD-LearnHub-Vercel.md`?**
  _High betweenness centrality (0.032) - this node is a cross-community bridge._
- **What connects `learn-platform/backend`, `RefreshTokenRequest`, `CreateUserRequest` to the rest of the system?**
  _773 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `time.Time` be split into smaller, more focused modules?**
  _Cohesion score 0.1038961038961039 - nodes in this community are weakly interconnected._
- **Should `useAuthStore` be split into smaller, more focused modules?**
  _Cohesion score 0.06448979591836734 - nodes in this community are weakly interconnected._
- **Should `📚 Product Requirements Document (PRD)` be split into smaller, more focused modules?**
  _Cohesion score 0.045454545454545456 - nodes in this community are weakly interconnected._
- **Should `dependencies` be split into smaller, more focused modules?**
  _Cohesion score 0.05714285714285714 - nodes in this community are weakly interconnected._
- **Should `context.Context` be split into smaller, more focused modules?**
  _Cohesion score 0.13333333333333333 - nodes in this community are weakly interconnected._