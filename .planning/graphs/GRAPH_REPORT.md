# Graph Report - learn-2  (2026-10-15)

## Corpus Check
- 77 files · ~125,281 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 618 nodes · 835 edges · 41 communities (34 shown, 7 thin omitted)
- Extraction: 100% EXTRACTED · 0% INFERRED · 0% AMBIGUOUS
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `e02e51f1`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- context.Context
- useAuthStore
- BRD — LearnHub
- 📚 Product Requirements Document (PRD)
- Learn Platform Backend
- dependencies
- UserHandler
- compilerOptions
- gorm.io/gorm.DB
- devDependencies
- cn
- AuthService
- PRD — LearnHub: Katalog Silabus Teknologi
- PracticeService
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
- Communities (31 total, 7 thin omitted)
- TRD — LearnHub
- 🔒 Security Features
- 🔑 API Endpoints
- 📊 Key Implementation Details
- 📦 Dependencies
- 🐛 Common Issues
- 🗄️ Database
- 🚀 Deployment
- 🚀 Quick Start

## God Nodes (most connected - your core abstractions)
1. `CourseRepository` - 24 edges
2. `Communities (31 total, 7 thin omitted)` - 23 edges
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
- `main()` --calls--> `NewPracticeHandler()`  [EXTRACTED]
  apps/api/main.go → apps/api/internal/handler/practice_handler.go

## Import Cycles
- None detected.

## Communities (41 total, 7 thin omitted)

### Community 0 - "context.Context"
Cohesion: 0.06
Nodes (26): Course, CourseResponse, CreateCourseRequest, CreateModuleRequest, ModuleResponse, UpdateCourseStatusRequest, User, Course (+18 more)

### Community 1 - "useAuthStore"
Cohesion: 0.06
Nodes (32): AdminDashboard(), Course, Course, CourseDetailPage(), Module, AddCoursePage(), CoursesPage(), DashboardPage() (+24 more)

### Community 2 - "BRD — LearnHub"
Cohesion: 0.07
Nodes (24): 1. Ringkasan, 2. Tujuan bisnis, 3. Pemangku kepentingan dan kebutuhan, 4. Ruang lingkup bisnis, 5. Model operasional dan pendapatan, 6. Risiko dan asumsi bisnis, 7. Tahapan bisnis, 8. Kriteria keberhasilan (+16 more)

### Community 3 - "📚 Product Requirements Document (PRD)"
Cohesion: 0.05
Nodes (44): 10. Rencana Pengembangan (Milestone), 11. Asumsi & Risiko, 12. Open Questions untuk Review, 1. Ringkasan Produk, 2. Tujuan Bisnis & Metrik Keberhasilan, 3. User Persona, 4.1 Modul Autentikasi & Manajemen Pengguna, 4.2 Modul Manajemen Kursus (Instructor) (+36 more)

### Community 4 - "Learn Platform Backend"
Cohesion: 0.20
Nodes (10): 📚 Additional Resources, 🏗️ Architecture, 🛠️ Development Commands, Learn Platform Backend, 📈 Performance Optimization, 📁 Project Structure, 🔗 Related Documentation, 🔄 Request Flow (+2 more)

### Community 5 - "dependencies"
Cohesion: 0.06
Nodes (35): dependencies, axios, class-variance-authority, clsx, framer-motion, @hookform/resolvers, js-cookie, lucide-react (+27 more)

### Community 6 - "UserHandler"
Cohesion: 0.39
Nodes (3): echo.Context, NewUserHandler(), UserHandler

### Community 7 - "compilerOptions"
Cohesion: 0.07
Nodes (28): compilerOptions, allowJs, esModuleInterop, incremental, isolatedModules, jsx, lib, module (+20 more)

### Community 8 - "gorm.io/gorm.DB"
Cohesion: 0.06
Nodes (23): getEnv(), Load(), Course, User, Course, User, UpdateUserRequest, User (+15 more)

### Community 9 - "devDependencies"
Cohesion: 0.08
Nodes (25): devDependencies, eslint, eslint-config-next, tailwindcss, @tailwindcss/postcss, @types/node, @types/react, @types/react-dom (+17 more)

### Community 10 - "cn"
Cohesion: 0.12
Nodes (15): Sidebar(), SidebarProps, Button, ButtonProps, Badge(), BadgeProps, Card(), CardProps (+7 more)

### Community 11 - "AuthService"
Cohesion: 0.13
Nodes (11): echo.Context, NewAuthHandler(), AuthResponse, Claims, LoginRequest, RegisterRequest, AuthService, NewAuthService() (+3 more)

### Community 12 - "PRD — LearnHub: Katalog Silabus Teknologi"
Cohesion: 0.08
Nodes (24): 10. Kriteria penerimaan rilis awal, 11. Ketergantungan dan keputusan yang masih perlu ditetapkan, 1. Visi dan ringkasan produk, 2. Sasaran dan bukan sasaran, 3. Persona dan kebutuhan, 4. Taksonomi dan pola silabus, 5. Katalog silabus, 6. Persyaratan produk dan UX (+16 more)

### Community 13 - "PracticeService"
Cohesion: 0.15
Nodes (11): echo.Context, NewPracticeHandler(), firstNonEmpty(), PracticeService, NewPracticeService(), statusMessage(), net/http.Client, PracticeHandler (+3 more)

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

### Community 31 - "Communities (31 total, 7 thin omitted)"
Cohesion: 0.06
Nodes (33): Communities (31 total, 7 thin omitted), Community 0 - "context.Context", Community 10 - "cn", Community 11 - "AuthService", Community 12 - "PRD — LearnHub: Katalog Silabus Teknologi", Community 13 - "PracticeService", Community 14 - "CourseHandler", Community 15 - "scripts" (+25 more)

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

## Knowledge Gaps
- **280 isolated node(s):** `learn-platform/backend`, `RefreshTokenRequest`, `CreateUserRequest`, `runnerResult`, `eslintConfig` (+275 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **7 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `📚 Product Requirements Document (PRD)` connect `📚 Product Requirements Document (PRD)` to `BRD — LearnHub`?**
  _High betweenness centrality (0.020) - this node is a cross-community bridge._
- **Why does `main()` connect `gorm.io/gorm.DB` to `context.Context`, `UserHandler`, `AuthService`, `PracticeService`, `CourseHandler`?**
  _High betweenness centrality (0.014) - this node is a cross-community bridge._
- **Why does `PRD — LearnHub: Katalog Silabus Teknologi` connect `PRD — LearnHub: Katalog Silabus Teknologi` to `BRD — LearnHub`?**
  _High betweenness centrality (0.012) - this node is a cross-community bridge._
- **What connects `learn-platform/backend`, `RefreshTokenRequest`, `CreateUserRequest` to the rest of the system?**
  _280 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `context.Context` be split into smaller, more focused modules?**
  _Cohesion score 0.06086956521739131 - nodes in this community are weakly interconnected._
- **Should `useAuthStore` be split into smaller, more focused modules?**
  _Cohesion score 0.06462585034013606 - nodes in this community are weakly interconnected._
- **Should `BRD — LearnHub` be split into smaller, more focused modules?**
  _Cohesion score 0.07389162561576355 - nodes in this community are weakly interconnected._