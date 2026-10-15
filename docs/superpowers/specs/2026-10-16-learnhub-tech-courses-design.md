# LearnHub technology courses design

Date: 2026-10-16
Status: awaiting user review

## Goal

Publish 21 Indonesian-language technology courses in LearnHub, with structured sections and lesson text, and make the lessons readable from the deployed course experience.

Topics: Google Cloud, Azure, Docker, Kubernetes, Terraform, Redis, Jenkins, GitHub Actions, Selenium, Appium, TestRail, PyTorch, LangChain, LangGraph, CrewAI, QLoRA, MCP, AI Agents, Postman, Linux, and Git.

## Current constraints

- The live frontend at `https://e-learning-rosy-eight.vercel.app` returns the catalog page, but its production bundle has an empty `NEXT_PUBLIC_API_URL`. The frontend cannot fetch catalog or lesson data until a production API URL is configured.
- The Go API has `Course`, `Section`, and `Lesson` models. Lessons support text content, but the create-module service ignores lesson body fields, course responses do not return lesson content, and the course detail page only renders lesson titles and completion controls.
- `GetAllPublishedCourses` currently does not filter by published status.
- Course ownership is required by the database, so seeded records must belong to an existing instructor account.
- Existing work under `docs/materials/` is outside this change and must remain untouched.

## Recommended design

1. Add a structured, version-controlled course catalog containing the 21 course descriptions, outcomes, prerequisites, ordered sections, and lesson bodies. Each course will use a consistent Indonesian learning path with fundamentals, applied practice, and a capstone section. Lesson count may vary where the subject needs it; every listed lesson must have readable text.
2. Add an explicit API seed command. It requires a configured existing instructor ID and database URL, creates only missing course slugs, and never overwrites an existing course, section, or lesson. Seeded courses are free (`price: 0`), Indonesian (`language: id`), and published. Seeding is an operator action, not an API startup side effect.
3. Restrict the catalog query to published courses. Include descriptions and lesson text in the authenticated course detail response, and render lesson content in the existing course detail page while keeping completion tracking.
4. Configure the frontend production environment with the actual API base URL and deploy the API with its production database. Do not guess a backend URL or place database credentials in frontend configuration.

## Data flow

The seed command reads the checked-in catalog and inserts missing records into PostgreSQL. The catalog endpoint returns only published course summaries. The authenticated course detail endpoint returns ordered sections and lessons, including text content. The course page displays a selected lesson and its completion state. Existing enrollment and progress records remain intact.

## Safety and compatibility

- Stable slugs make seeding repeatable. Existing records are never updated or deleted by the seed command.
- No sample user, password, or credential is created.
- Draft courses stay out of the catalog. Seeded courses are free, so the published lesson text is accessible to authenticated learners.
- Existing `modules` response fields and progress behavior remain compatible; lesson content and course description are additive fields.
- If no production API URL or instructor account is available, course content can be prepared in source but cannot be truthfully reported as live.

## Acceptance criteria

- The 21 requested course titles are present exactly once in the published catalog.
- Every course has a coherent ordered curriculum and non-empty Indonesian lesson text.
- Re-running the seed command creates no duplicates and does not alter existing course data.
- Authenticated learners can open a course and read each lesson; completion tracking still works.
- The deployed frontend requests the configured production API, and the live catalog displays all 21 seeded courses.

## Verification plan

Validate the catalog data structure and stable slugs, run the seed command against a non-production database first, inspect API catalog and detail responses, then verify the live catalog and a lesson page after deployment. Do not run the production seed until the target API and instructor account are confirmed.
