# Katalog Materi Pembelajaran Teknis — LearnHub

Dokumentasi materi pembelajaran (*learning materials*) terstruktur untuk platform e-learning **LearnHub**. Setiap modul materi mencakup **Tujuan Belajar**, **Prasyarat**, **Konsep Dasar**, **Praktik Implementasi Kode**, **Fitur Utama & Pola Arsitektur**, **Praktik Kerja (Studi Kasus Nyata)** dengan lima hint bertahap, **Lanjutan / Produksi**, **Pemeriksaan Pemahaman**, dan **Referensi Resmi**.

---

## Daftar Materi Berdasarkan Domain

### A. Bahasa & Framework Frontend
1. [JavaScript Modern (ES6+ & Runtime)](javascript.md) — Runtime, Event Loop, closures, asynchronous execution, DOM, dan modul.
2. [TypeScript](typescript.md) — Sistem tipe statis, Generics, Discriminated Unions, Type Narrowing, dan type safety.
3. [Next.js](nextjs.md) — App Router, Server vs Client Components, Server Actions, Caching/ISR, dan Streaming Suspense.
4. [Vue.js 3](vue.md) — Composition API (`ref`, `reactive`), reusable Composables, Pinia state management, dan Vue Router.
5. [React](react.md) — Rekonsiliasi Virtual DOM, hooks dasar & custom hooks, state architecture, dan memoization.

### B. Bahasa & Framework Backend
6. [Modern PHP 8.x](php.md) — Strict types, OOP modern, match expression, PDO prepared statements, dan standar PSR.
7. [Laravel](laravel.md) — Arsitektur MVC, Eloquent ORM & Eager Loading (N+1 fix), Form Requests, Queues, dan Transactions.
8. [Express.js](express.md) — Middleware pipeline, request validation Zod, centralized error handling, JWT auth, dan security headers.
9. [NestJS](nestjs.md) — Arsitektur modular enterprise, Dependency Injection, DTO validation pipes, guards, dan interceptors.
10. [AdonisJS](adonisjs.md) — TypeScript-first full-stack MVC, Lucid ORM, VineJS validation, dan Bouncer authorization.
11. [Java Modern (17/21 LTS)](java.md) — Records, Sealed Classes, Pattern Matching, Stream API, Generics, dan Virtual Threads (Project Loom).
12. [Spring Boot 3](springboot.md) — Inversion of Control / DI, Spring MVC REST API, Spring Data JPA, `@Transactional`, dan Actuator.
13. [Go (Golang)](go.md) — Structs, implicit interfaces, explicit error values, Goroutines, Channels, dan context cancellation.
14. [Echo (Go)](echo.md) — Radix tree router, context binding, custom validator, centralized error handler, dan graceful shutdown.
15. [Gin (Go)](gin.md) — Fast HTTP router, model binding & tag validation, middleware pipeline, dan safe context concurrency.
16. [Fiber (Go)](fiber.md) — Engine Fasthttp zero memory allocation, Express-like syntax, rate limiting, dan buffer lifecycle.
17. [Python 3 Modern](python.md) — Python data model, type hints, dataclasses, generators (`yield`), context managers, dan `asyncio`.
18. [Flask](flask.md) — Application Factory pattern, Blueprints, request context vs app context, Flask-SQLAlchemy, dan Gunicorn WSGI.
19. [Django & DRF](django.md) — Pola MVT, Django ORM (`select_related` vs `prefetch_related`), DRF ViewSets/Serializers, dan `atomic` transactions.
20. [FastAPI](fastapi.md) — Server ASGI Uvicorn, validasi Pydantic v2, dependency injection `Depends()`, async database, dan OpenAPI docs.
21. [Ruby](ruby.md) — Pure OOP, blocks/procs/lambdas, enumerable, modules mixins, dan refinements.
22. [Ruby on Rails](rails.md) — Convention over Configuration, Active Record, migrations, strong parameters, dan active job.
23. [C# 12](csharp.md) — Records, primary constructors, collection expressions, LINQ, pattern matching, dan `ReadOnlySpan<T>` zero-allocation.
24. [.NET 8 (ASP.NET Core)](dotnet.md) — Minimal APIs, EF Core 8 (`AsNoTracking`), service lifetimes, custom middleware, dan `IExceptionHandler`.
25. [C++ Modern (17/20)](cpp.md) — RAII, smart pointers (`unique_ptr`, `shared_ptr`), move semantics, STL containers, dan sanitizers.

### C. Mobile Development
26. [Flutter & Dart](flutter.md) — Widget Tree (Stateless/Stateful), Impeller rendering, state management BLoC, dan optimasi 60/120 FPS.

### D. Basis Data & Backend Platform
27. [PostgreSQL](postgresql.md) — Relational schema & constraints, MVCC, indexing (B-Tree, GIN, Partial), `EXPLAIN (ANALYZE, BUFFERS)`, dan VACUUM bloat mitigation.
28. [MySQL 8.x](mysql.md) — Storage engine InnoDB, Clustered vs Secondary Index, `EXPLAIN ANALYZE`, covering index, dan resolusi deadlock.
29. [MariaDB](mariadb.md) — Storage engine Aria, System-Versioned Tables (time travel queries), Galera Cluster synchronous multi-master, dan kompatibilitas MySQL.
30. [MongoDB](mongodb.md) — NoSQL document model BSON, embedding vs referencing, Aggregation Pipeline, dan aturan ESR compound indexing.
31. [Supabase](supabase.md) — Backend-as-a-Service berbasis PostgreSQL, Row Level Security (RLS) policies, Supabase Auth, Storage, dan Realtime WebSocket.
32. [Oracle Database & PL/SQL](oracle.md) — Arsitektur SGA/PGA, paket PL/SQL, Sequences, Execution Plans, dan mitigasi error ORA-01555.
33. [Firebase & Cloud Firestore](firebase.md) — Realtime NoSQL document listeners (`onSnapshot`), Security Rules deklaratif, Firebase Storage, dan Cloud Functions.

### E. Cloud, Platform Engineering & DevOps
34. [AWS Cloud Architecture](aws.md) — Shared Responsibility Model, IAM Least Privilege, VPC Multi-tier Subnets, S3 presigned URLs, dan Incident Response.
35. [Nginx](nginx.md) — Reverse proxy, server blocks, static asset handling, SSL/TLS, load balancing, dan diagnosis `502 Bad Gateway`.
