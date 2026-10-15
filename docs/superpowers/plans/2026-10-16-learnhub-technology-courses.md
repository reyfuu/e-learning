# LearnHub technology courses Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Publish 21 Indonesian-language technology courses with structured lesson text that authenticated learners can read in LearnHub.

**Architecture:** Keep course material in a version-controlled JSON catalog embedded by the Go seed package. An explicit seed command inserts missing courses, sections, and lessons into PostgreSQL and requires `--apply` before writes. The authenticated course API returns lesson text, and the existing Next.js detail page renders it as escaped plain text. Production publication requires a configured backend API URL and an existing instructor owner.

**Tech Stack:** Go 1.22, Echo, GORM, PostgreSQL, Next.js App Router, TypeScript.

**Spec:** `docs/superpowers/specs/2026-10-16-learnhub-tech-courses-design.md`

## Global Constraints

- Courses use Indonesian (`language: id`), are free (`price: 0`), and are published.
- Seed only missing stable slugs; never overwrite or delete existing course data.
- Use an existing instructor ID; do not create a sample user or credentials.
- Do not guess the production API URL or expose database credentials to the frontend.
- Preserve existing module fields and lesson completion behavior.
- Lesson bodies are plain text and rendered as escaped text, never injected as HTML.
- Leave existing untracked files in `docs/materials/` untouched.

## Review Focus

- Duplicate seed runs skip existing slugs without changing their content; the seed summary reports created and skipped counts.
- Missing instructor ID or database URL stops the seed before writes.
- Empty course, section, lesson title, or lesson body is rejected before database writes.
- Draft courses stay out of the catalog, and non-owner learners cannot fetch draft lesson text.
- Lesson text is rendered as escaped text, including strings that resemble HTML.

---

### Task 1: Author the 21-course catalog

**Files:**
- Create: `apps/api/internal/seed/technology-courses.json`

**Interfaces:**
- Catalog entries contain `slug`, `title`, `description`, `outcomes`, `prerequisites`, ordered `sections`, and ordered `lessons` with `title`, `duration_seconds`, and plain-text `content`.
- Each course has fundamentals, applied practice, and capstone sections, with three lessons per section.

- [ ] Add all requested titles: Google Cloud, Azure, Docker, Kubernetes, Terraform, Redis, Jenkins, GitHub Actions, Selenium, Appium, TestRail, PyTorch, LangChain, LangGraph, CrewAI, QLoRA, MCP, AI Agents, Postman, Linux, and Git.
- [ ] Write course-specific outcomes, prerequisites, and Indonesian lesson text for all 189 lessons.

### Task 2: Add safe explicit course seeding

**Files:**
- Create: `apps/api/internal/seed/courses.go`
- Create: `apps/api/cmd/seed-courses/main.go`

**Interfaces:**
- `CourseSpec` and nested section/lesson structs decode the embedded JSON catalog.
- `SeedCourses(ctx context.Context, db *gorm.DB, instructorID string, courses []CourseSpec, apply bool) (SeedSummary, error)` validates all input, verifies the instructor, counts existing slugs, and inserts missing records in one transaction only when `apply` is true.
- The command requires `DATABASE_URL` and `COURSE_SEED_INSTRUCTOR_ID`; without `--apply`, it performs a dry run and writes nothing.

- [ ] Embed and decode the JSON catalog with Go `embed` and validate it before database writes.
- [ ] Reject catalogs unless they contain exactly 21 unique slugs and all required fields and lesson bodies are non-empty.
- [ ] Require the configured owner to exist and have the instructor role.
- [ ] On apply, insert missing courses, sections, and lessons with published/free/Indonesian values; leave existing slugs untouched.
- [ ] Print created/skipped counts and exit non-zero on invalid config or failed database operations.
- [ ] Run the dry-run command against a configured database and confirm the summary reports 21 new courses on an empty target.

### Task 3: Return published catalog and readable lessons

**Files:**
- Modify: `apps/api/internal/repository/course_repository.go`
- Modify: `apps/api/internal/model/course.go`
- Modify: `apps/api/internal/service/course_service.go`

**Interfaces:**
- `CourseResponse` adds `description` while retaining existing fields.
- `ModuleResponse` adds `section_title`, `content_type`, and nullable `content`, while retaining progress fields.
- Course detail resolves an instructor-owned course first; otherwise it returns only published courses.

- [ ] Filter catalog queries to `status = published`.
- [ ] Map descriptions and ordered section/lesson text into additive response fields.
- [ ] Preserve lesson completion state and keep draft details available only to their instructor owner; fall back from an instructor-owned lookup only on not-found, not on database errors.

### Task 4: Render course descriptions and lesson text

**Files:**
- Modify: `apps/web/src/app/dashboard/courses/[id]/page.tsx`
- Modify: `apps/web/src/app/dashboard/catalog/page.tsx`
- Modify: `apps/web/src/lib/api.ts`

**Interfaces:**
- Course and lesson response types include `description`, `section_title`, `content_type`, and `content`.
- Learners can select lessons by section and read their text; completion controls remain available.

- [ ] Add the response fields to the TypeScript interfaces.
- [ ] Render section headings and selectable lesson titles; keep completion controls separate from lesson selection.
- [ ] Show selected lesson text with preserved line breaks using React text rendering, not `dangerouslySetInnerHTML`.
- [ ] Keep lesson selection keyboard-accessible and completion controls separately operable.
- [ ] Show the course description in the catalog and label the lesson count accurately.

### Task 5: Connect and verify production publication

**Files:**
- Modify: Vercel production environment configuration only after the API origin is confirmed.

- [ ] Confirm the deployed API origin and existing instructor ID without printing credentials.
- [ ] Set `NEXT_PUBLIC_API_URL` to the confirmed API origin and deploy the frontend/API as needed.
- [ ] Review the seed dry-run target; only then run with `--apply` against the intended database.
- [ ] Verify the live catalog has all 21 titles and open one lesson to confirm text renders.
- [ ] If no deployed API or instructor account is available, stop before production writes and report the missing prerequisite.
