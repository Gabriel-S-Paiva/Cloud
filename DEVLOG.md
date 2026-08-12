# Devlog

Day-by-day notes from building Owned Cloud. Cleaned up from paper, kept honest bugs and dead ends included.

---

## Day 1 - The Spark & Tech Stack

**The Spark:**
Looked at a pile of old hardware sitting around, too slow to be primary machines, too good to throw away. Decided to repurpose them into a self-hosted personal cloud drive for a small group of family/friends (~5 users).

**Initial Design & Pain Points:**
Identified four theoretical bottlenecks right out of the gate:
- **Server Performance:** Old hardware will struggle with heavy runtimes.
- **Remote Access:** Need a way to reach the server outside the local LAN.
- **File Sizes:** Handling large uploads reliably on constrained specs.
- **Security:** Exposing home network hardware to the public internet.

**Tech Stack Decisions:**
Selected a stack focused on lightweight performance and personal learning goals:
- **Frontend:** Svelte (Learning opportunity / lightweight UI)
- **Backend:** Go (Learning opportunity / fast, low-memory compiled binary)
- **Database:** SQLite (Simple, file-based, minimal overhead for 5 users)
- **Containerization:** Docker
- **Networking/DNS:** DuckDNS + Cloudflare (Theoretical plan for DDNS + SSL/Tunneling)

> **Hindsight Note (Looking back from later in the build):**
> *DuckDNS and Cloudflare don't pair cleanly together. DuckDNS was ultimately dropped to avoid exposing the home router directly. Cloudflare Tunnels are on hold until I claim an domain, so right now the app is strictly deployed to local LAN.*

---

## Day 2 - Hello World & First Schema Draft

**Time Spent:** ~3 hours

**Focus:** Bootstrapping the Go backend and drafting the core database schema.

### 1. Backend: The API is Alive
Before touching deployment or hosting, I needed an actual server running. Started with the fundamental Go HTTP handler pattern:

```go
package main

import (
    "fmt"
    "net/http"
)

func hello(w http.ResponseWriter, req *http.Request) {
    fmt.Fprintf(w, "hello\n")
}

func main() {
    http.HandleFunc("/hello", hello)
    http.ListenAndServe(":8080", nil)
}
```

**Learned / Go vs. Express.js:**
Coming from Express, passing the `http.Request` as a pointer (`*http.Request`) felt odd at first. My hypothesis: Go passes request pointers so middlewares can attach context or modify state across the request lifecycle without copying the whole object. *(Target for tomorrow: investigate Go context and middleware patterns to confirm this).*

### 2. Database Schema (First Draft & Refinements)
Sketched out the initial entities (`Users`, `Files`, `Shares`) and caught a few design flaws almost immediately.

* **Initial Sketch:**
  * **User:** `username` (PK), `hashed_password`, `role`, `quota`
  * **Files:** `id`, `display_name`, `owned_by` (FK), `size`, `uploaded_at`, `last_modified`, `path_file`, `folder_path`
  * **Share:** `file_id` (FK), `shared_with` (FK), `owner` (FK)

* **Refinements & Realizations:**
  * **Keys:** Switched `username` (PK) to a unique `id`. Using `username` as a primary key creates cascading update nightmares if a user changes their handle or email.
  * **Data Types:** Changed `size` from `decimal` to `int` (storing raw bytes).
  * **Storage Cleanup:** Dropped `path_file`—it's useless in a flat object/file storage pattern.



### Open Debates & Next Steps
- [ ] **Folders:** Storing `folder_path` as a string is going to lead to messy path-parsing logic down the line. Need to figure out a cleaner folder tree structure (e.g., parent IDs or adjacency lists).
- [ ] **Quotas:** Should `quota_used` be a calculated aggregate query (`SUM(size)`) or a stored column updated on every upload/delete?
- [ ] **Auth / Sharing:** Determine share permission granularities and user sign-up approval workflows.

---

## Day 3

---

## Day 4 - Baby Steps & CLI Quick Wins

**Time Spent:** ~30 minutes

**Focus:** Quick API wins, error handling helpers, and direct SQLite database seeding.

### Wins & Progress
* **Conflict Handling:** Added a `409 Conflict` HTTP status code when attempting to create a user with a duplicate username.
* **User Endpoint:** Built a functional (though static) `/user/me` endpoint that hardcodes reading user ID `1` from the database.
* **JSON Error Helper:** Refactored error handling across the backend. Replaced raw text/string error responses with a standardized JSON error helper function for consistent API output.
* **Database CLI Tinkering:** Popped into the SQLite CLI directly on `data.db` to manually seed a user entry without needing an API endpoint ready for user creation—always satisfying to manipulate the raw DB state.


### Targets for Next Time
- [ ] Implement full auth flow: login endpoint, session management, and auth middlewares.
- [ ] Build user approval/rejection workflows for new signups.

---

## Day 5 - Session Auth & Architectural Pivot

**Time Spent:** 4 hours

**Focus:** Settling on an authentication strategy, executing schema updates, and building server-side session management.

### Architectural Decision: Session Tokens vs. JWTs
Started the session evaluating auth options. JWT was the default choice, but after stepping back and looking at the project goals, implementing it felt like premature optimization.

* **The JWT Myth for Small Apps:** The primary selling point of JWTs is stateful-free authorization without hitting the database on every request. However:
  1. Given this project's scale (~5 users on local hardware), database lookups are fast and negligible in overhead.
  2. If implementing refresh tokens for better UX and security, you end up needing database checks anyway.
* **The Solution:** Switched to traditional **server-side session tokens**. Clean, easy to invalidate, and simple to maintain.


### Progress & Implementation
* **Schema Updates:** Added a `sessions` table to track valid logins:
  * Fields: `id` (token/UUID), `user_id` (FK to user), `expires_at` (timestamp)
* **Auth Flow:** Built the full login route—verifying credentials, generating session tokens, and persisting them to SQLite.
* **Milestone:** Full session-based authentication backend logic is officially complete!


### Targets for Next Time
- [ ] Build the auth middleware to validate incoming session tokens.
- [ ] Fix the hardcoded user ID `1` in `GET /user/me` using the authenticated session context.
- [ ] Set up admin-only route protection.
- [ ] Build an endpoint to list pending registration requests for admin approval.

---

## Day 6 - Middleware & Go Handler Patterns

**Time Spent:** 3 hours

**Focus:** Implementing authentication and authorization middlewares, plus learning Go's native HTTP handler wrapping patterns.

### Progress & Wins
* **Auth Middlewares:** Successfully built both **Authentication** (validating the session token) and **Authorization** (checking user roles/permissions) middlewares.
* **Refactoring:** Protected key endpoints so unauthenticated or unauthorized requests are rejected early before hitting core route logic.



### Key Concept: Go vs. Express Middlewares
Coming from Express.js, middleware architecture works quite differently under the hood:

* **Express.js:** Uses continuation passing via a `next()` argument (`app.use((req, res, next) => ...)`).
* **Go Standard Library:** Uses **higher-order function wrapping (decorators)**. A middleware takes an `http.Handler` as an argument and returns a new `http.Handler`. 

To pass control down the stack, the middleware explicitly calls `next.ServeHTTP(w, r)` on the inner handler. It feels very clean once it clicks—no global middleware stack needed, just explicit function composition.



### Targets for Next Time
- [ ] Implement admin-only routes (e.g., listing pending register requests).
- [ ] Update `GET /user/me` to read the user identity from the request context provided by the auth middleware.

## Day 7 - Folder Hierarchy & Go SQL Types

**Focus:** Designing the folder system and handling nullable foreign keys in Go.

### Architectural Decision: Folders Before Files
Even though files are the core deliverable, folder structures had to come first. Since files maintain a reference to their parent folder, the database tables and endpoints for folder management needed to exist before building out upload logic.

---

### Progress & Technical Challenges

* **Go Concept: Handling `NULL` with `sql.NullInt64`**
  * Root folders don't have a parent ID (they are `NULL` in SQLite). In Go, attempting to scan a `NULL` database value into a standard `int64` variable throws a runtime error.
  * **Solution:** Used Go's native `sql.NullInt64` type from `database/sql`. It provides a struct containing both `.Int64` (the value) and `.Valid` (a boolean indicating if the field is non-null).

* **Folder Content Retrieval:**
  * Fetching the contents of a directory required building a unified response struct capable of holding arrays of both `File` structs and `Folder` structs in a single JSON payload for the frontend.

---

### Targets for Next Time
- [ ] Implement nested folder creation endpoints.
- [ ] Connect file metadata records to parent folder IDs.

## Day 8 — Go Slices, Late-Night API Grind & Sharing Logic

**Time Spent:** Late night (~2 AM finish)

**Focus:** Clearing out boilerplate endpoints, diving deep into Go slices under the hood, and drafting file sharing queries.

### Go Concept: Arrays vs. Slices & Mutation Side-Effects
Go arrays are fixed-size, so dynamic collections use **slices**. Under the hood, a slice isn't a direct data container,it's a header struct containing:
1. A **pointer** to an underlying array
2. The current **length**
3. The maximum **capacity**

**The Gotcha:** Because a slice holds a *pointer* to the underlying array, slicing an existing slice doesn't copy the data. Both slices point to the exact same array in memory—meaning mutations in one will mutate the other!

---

### Progress & Brainstorming

* **Boilerplate Cleanup:** Mindlessly ground through standard CRUD/utility endpoints to clear out backend technical debt.
* **Late-Night Sharing Logic:** Sketched out the SQL joining logic for file access:
  * *Shared with me:* `WHERE shared_with = session.user_id`
  * *Files I've shared:* `WHERE owner = session.user_id` (JOIN on `shares` table)

---

### Remaining Backend Checklist
- [ ] **Quota Bookkeeping:** Update user storage usage on file insert and delete.
- [ ] **Sharing & Ownership:** Integrate share permissions into ownership validation middleware.
- [ ] **File I/O:** Real disk read/write stream handling for file uploads & downloads.
- [ ] **Admin Utilities:** Finish minor administrative endpoints.

*Next up: File I/O stream handling, then onto the Svelte frontend!*

## Day 9 - Quick Handler Polish

**Focus:** Miscellaneous handler cleanup and tightening up route signatures.

* Ground through minor backend handler tweaks to ensure consistency before finalizing API endpoints.

---

## Day 10 - API "Complete", Schema Revisions & Svelte Onboarding

**Focus:** Finalizing the backend, adapting the schema for upload states, and setting up CORS for frontend development.

### Progress & Wins
* **Upload State Schema Tweaks:** Updated file metadata tables to support asynchronous/chunked file uploads. Added fields to track upload status (e.g., `pending`, `completed`), progress signals, and incoming upload sizes.
* **CORS Middleware:** Implemented a Cross-Origin Resource Sharing (CORS) middleware on the Go server to prevent cross-origin fetch/XHR blocking when developing the frontend locally on a different port.
* **Milestone Achieved:** The initial backend API feature set is officially "built"! 🎉

---

### Targets for Next Time (Frontend Phase)
- [ ] Initialize the Svelte/SvelteKit frontend project.
- [ ] Set up basic API client fetch wrappers for authentication and session management.
- [ ] Build the initial login / registration view.

## Day 11 - Back to Pen and Paper: Wireframing the UI

**Focus:** Designing the multi-pane dashboard layout and user experience flow before committing to code.

### UI/UX Design Session
Stepped away from the code today to focus on visual hierarchy and workflow. For a project like this (a cloud drive management interface), jumping straight into Svelte/CSS is a recipe for messy re-layouts later.

**Wireframe Notes:**
* **Dashboard Goal:** Needs a classic, intuitive file manager layout: left sidebar, main file browser grid/list, and a third pane for file details/previews (similar to Windows Explorer or macOS Finder).
* **Navigation:** Sidebar should handle navigation (My Drive, Shares, Recents) and administrative sections (Users, Approvals).
* **Core Interaction:** Main panel must support classic drag-and-drop file uploading and intuitive right-click context menus for actions (Rename, Share, Delete).

---

### Initial Sketches

I created two quick sketches today:

1.  **Sketch 1: Main Dashboard Layout** - Roughing out the tripartite layout and responsive behavior.
    > *[Insert/Link photo of Sketch 1 here, e.g., images/day11_sketch_dashboard.png]*

2.  **Sketch 2: Login And Register** - some quick form work.
    > *[Insert/Link photo of Sketch 2 here, e.g., images/day11_sketch_upload_mgmt.png]*

---

### Targets for Next Time (Frontend Build Phase)
- [ ] Initialize the Svelte/SvelteKit frontend project.
- [ ] Implement the base CSS layout/grid system derived from today's sketches.
- [ ] Build the initial login / registration views.

## Day 12 - Rescuing Legacy Hardware & A Code-First Learning Plan

**Focus:** Rescuing lost photos off legacy hardware, followed by structuring a non-traditional frontend learning workflow.

### The Side Quest: Legacy Hardware Recovery Mission
Wasted most of my day and brainpower fighting an ancient tablet to rescue old childhood photos before the device died for good.

* **The Constraints:** No working web browser, no Bluetooth, no USB support, and all legacy share methods were dead.
* **The Hack:** 
  1. Spun up a mini HTTP file share on my phone so the tablet could download a file manager APK over Wi-Fi.
  2. Used the APK to launch a remote manager on the tablet.
  3. Connected the tablet to a PC over an FTP server via a phone hotspot (bypassing router security) to finally pull all the photos off safely.

By the time the photos were backed up, I was completely exhausted. I shifted to light planning so the day wasn't a total wash.

---

### Frontend Strategy: Code-First Learning Pipeline
Normally, I'd draft a happy path, build Figma designs, and then write code. But since **learning Svelte** is a primary objective, I'm flipping the sequence to avoid juggling framework learning, architecture, and visual design all at once:

1. **Phase 1: Pure Code (Raw Functionality)**
   * Get features working first with raw HTML. No styling distractions.
2. **Phase 2: Architecture & Refactoring**
   * Organize the working code into clean components so I can leave the project for 6 months and come back without getting lost.
3. **Phase 3: Visual Design & Styling**
   * Coat the raw, working application in a clean layer of Tailwind CSS.

---

### Breaking Down Phase 1 (Pure Code Objectives)
- [ ] **Folder Routing:** Implement route parameters for navigating folder hierarchies.
- [ ] **Data Fetching:** Wire up backend API calls to Svelte state/stores.
- [ ] **Rendering:** Display raw directory contents and metadata lists.
- [ ] **Navigation:** Build forward/back/up directory traversal logic.

## Day 13 - Navigation, Breadcrumbs & API Refactoring

**Focus:** Building frontend directory navigation, adding basic folder CRUD, and fixing API response payloads on the fly.

### Progress & Wins
* **Navigation & Breadcrumbs:** Built out directory traversal logic with dynamic breadcrumb paths so users can navigate up and down folder trees.
* **Admin Features:** Admin user approval/rejection workflows are officially functional in the UI.
* **Milestone Progress:** The Home/Drive view is operational! Only `/profile` and `/shares` remain before moving out of the raw HTML phase and into code architecture/refactoring.

---

### Pain Point & Immediate Refactor: Optimistic UI Updates
Caught a major inefficiency in the backend API design while wiring up folder creation in Svelte:

* **The Problem:** The `POST` endpoint to create a folder originally just returned a generic `200 OK`. This forced the frontend to make a second, wasteful API request (`GET /folders`) just to fetch the newly created folder's details and render it in the UI.
* **The Fix:** Refactored creation endpoints immediately to return the newly generated resource `id` in the JSON response payload.
* **The Result:** The frontend can now immediately construct/derive the new folder object locally and append it to the active UI state without waiting on a second network request.

---

### Backlog Note: Drag-and-Drop
* **Feature Deferred:** Decided to hold off on implementing drag-and-drop file/folder moving until **Phase 2 (Architecture)** or **Phase 3 (UI/Tailwind)**. Trying to handle native drag-and-drop HTML events on unstyled elements creates messy boilerplate that will just need to be rewritten once component structures solidify.

---

### Targets for Next Time
- [ ] Build out the `/shares` page (shared files/folders view).
- [ ] Build out the `/profile` page (user settings & quota indicators).
- [ ] Prepare code base for Phase 2: Architecture & Component Cleanup.

## Day 14 — The Backend Returns: Ownership Checks & Type Refactoring

**Focus:** Fixing backend file ownership verification and syncing API types across the stack.

### Progress & Technical Debt Cleanup
* **Ownership Validator Fix:** Returned to the Go backend to patch issues with ownership checks on shared files. The original middleware was too strict and blocked valid shared access patterns.
* **API Payload Refactor:** Realized mid-implementation that several sharing endpoints were returning incorrect response structures. Refactored the Go handler signatures and updated response objects.
* **Frontend Sync (`api.ts`):** Synchronized the Svelte TypeScript interfaces/types in `api.ts` to match the updated Go backend responses.

---

## Day 15 - Sharing Workflows Complete & The Component Pivot

**Focus:** Wrapping up share permissions, identifying a key user-discovery gap, and preparing for Svelte component architecture.

### Progress & Wins
* **Share CRUD Operational:** Finished implementing both the **Create Share** and **Update Share** (permissions modification) features across the API and frontend!

---

### Pain Point / Missing Link: User Discovery
While testing the sharing flow, hit an obvious missing UX requirement:
* **The Problem:** You can create a share for "someone," but there is currently no backend endpoint or UI component to search/list users on the system to share *with*. 
* **Backlog Item:** Need to add a lightweight user search endpoint.

---

### API Cleanup & Next Steps
* **Pruning Dead Endpoints:** Considering dropping `getFolder` and `getFile` single-item endpoints if full directory tree payloads satisfy all current UI needs.
* **Missing Minor Features:** Need to wire up a proper Logout action in the sidebar.

**Milestone Achieved:** With the backend stable and all core CRUD features functional, Phase 1 (Pure Code) is essentially wrapped! Time to move into **Phase 2 (Component Architecture)** and **Phase 3 (Tailwind Styling)**. 

## Day 16 - Tailwind Styling, Toast Notifications & Roadmap to v1

**Time Spent:** 1 hour

**Focus:** Kicking off Phase 3 (Styling & Componentization), building the login page UI, and mapping out releases v1 through v3.

### Progress & Wins
* **Login UI Complete:** Coated the raw login view with clean Tailwind CSS styling.
* **Toast Notification System:** Integrated toast alerts on the frontend to provide real-time feedback for authentication success/failure states.

---

### Immediate Todo List
1. Add color variants to toast alerts (success green, error red, warning yellow).
2. Build styled Svelte components for file and folder items.
3. Implement the styled sidebar navigation layout.
4. Fix the `/home` root folder navigation routing bug.

---

### Project Roadmap & Strategy
* **v1 (Release Target):** Implement Drag-and-Drop file uploads/moving, complete a full end-to-end bug sweep, and deploy to local hardware.
* **v2 (Infrastructure):** Add automated backend/frontend test suites and set up a CI/CD deployment pipeline.
* **v3 (Iterative Polish):** Introduce deferred secondary features to test and validate the automated CI/CD pipeline in a real-world workflow.

## Day 17 - The Final Stretch Stretch

**Focus:** Manual testing, edge-case bug hunting, and final UI/UX polish ahead of v1 deployment.

### The "90/90 Rule" of Development
Deployment feels incredibly close, but as every developer knows: *the first 90% of the code accounts for the first 90% of the development time. The remaining 10% of the code accounts for the other 90% of the development time.*

Every time a late bug gets squashed, two small UX papercuts show up in manual testing. 

---

### Current Polish & Bug Sweep
* **Manual Testing:** Clicking through every authentication, folder navigation, and file action edge-case to make sure no unhandled errors crash the UI.
* **Late Bug Fixes:** Patching minor layout overflows, state sync glitches when switching directories, and toast notification timing.
* **Final Touchups:** Double-checking mobile/small screen responsiveness and cleaning up console log noise.

---

### Targets Before Deployment 🚀
- [ ] Fix the `/home` root folder navigation routing glitch.
- [ ] Implement native drag-and-drop file upload/moving in the styled file view.
- [ ] Build release Docker container and run initial local deployment test on target hardware.

## Day 18 - Architectural Epiphany: Modular Feature Extensions

**Focus:** Designing a low-power primary server with optional high-spec microservice extensions.

### The Challenge: Always-On Server vs. Heavy Features
I want 100% uptime, which means running the cloud drive on my always-on, low-spec old PC. However, old hardware can't handle modern heavy tasks (like local AI photo tagging, video transcoding, or smart search models).

---

### The Idea: Dynamic Capability-Based Modules
Instead of compromising on features or forcing the old PC to suffer, the architecture can support **optional local extensions**:

```text
+-------------------+
|     Frontend      |
+-------------------+
     |         |
     |         +---------------------------+
     v (Core Drive API)                    v (Optional AI API)
+-------------------+             +-------------------+
|  Old PC (Core)    |             |  New PC (AI/Heavy)|
|  * Always On      |             |  * Online Sometimes|
|  * Storage & Auth |             |  * ML / Tagging   |
+-------------------+             +-------------------+
```

### How It Works in Practice
* **The Core (Old PC):** Runs 24/7. Handles file management, directory navigation, upload/download, and user authentication.
* **The Add-on (New PC):** Runs an AI or heavy-processing microservice whenever the desktop is turned on.
* **Smart UI Integration:** The frontend checks if the New PC's service is reachable. If online, advanced features (like "Auto-tag Images with AI") appear. If offline, the core cloud drive functions 100% normally without any missing dependencies.

*Zero code written today, but a super clean concept for adding heavy desktop-powered features to a low-spec home lab without breaking uptime.*
---

## Day 19 - The Bug War (Nearing Victory)

**Focus:** Final bug hunting and edge-case squashing ahead of environment setup.

### The Grind
* **Bug Sweeping:** Spent the day fighting late-stage bugs in state sync and edge cases. Every time one gets resolved, another minor edge case pops up—but the list is shrinking rapidly.
* **Pre-Deployment Prep:** Getting the codebase stabilized so the server environment setup can start smooth.

---

## Day 20 - Development Frozen & Machine Setup Begins

**Focus:** Closing active feature development and transitioning to server deployment preparation.

### Major Milestone
* **Development Freeze:** Officially froze the codebase for v1! No new features—all focus shifts to getting the deployment target machine configured and prepped.
* Exhausted after a full-day push, but the project has officially crossed the finish line from "dev code" to "deployable software."

---

## Day 21 - Dockerization & LAN Deployment

**Focus:** Containerizing the stack and getting the cloud drive live on the home network.

### Accomplishments
* **Docker Setup:** Wrote and configured the `Dockerfile` and `docker-compose` setup to run the Go backend and built Svelte frontend together.
* **LAN Live:** Successfully deployed the containers onto the home server! The cloud drive is now accessible across local devices via LAN IP.

---

## Day 22 - Live Environment Bug Patching

**Focus:** Fixing real-world bugs discovered during actual network usage.

### Post-Deploy Fixes
* Deploying to real hardware on a real local network exposed a few environment-specific glitches (CORS headers, file pathing, and network upload stream quirks).
* Applied hotfixes directly to the live container environment to stabilize v1.

## Day 23 - Linting, Automated Testing & Self-Hosted CI/CD Pipeline 🚀

**Focus:** Code cleanup, static analysis, and building an automated self-hosted deployment pipeline.

### Progress & Polish
* **Code Quality & Linting:** Cleaned up residual backend debt and ran full frontend linter/formatter passes using **ESLint** and **Prettier** to enforce consistent style across the codebase.
* **Backend Unit Tests:** Added initial Go test suites to validate core logic during automated builds.

---

### Major Milestone: Self-Hosted CI/CD Pipeline
Built and verified a complete automated integration and deployment workflow from scratch:

1. **Self-Hosted Runner:** Set up a local runner instance inside the home network (LAN) so deployment triggers can target the home server directly without exposing internal ports to the public web.
2. **CI (Continuous Integration):**
   * Automatically runs `go test` on incoming backend commits.
   * Runs ESLint + Prettier checks on the frontend code.
3. **CD (Continuous Deployment):**
   * Automatically builds and deploys the latest container release straight to the local hardware once checks pass.

> **The Developer Experience:** Watching those first few GitHub Actions pipeline runs go through live is nerve-wracking—hoping every step turns green. But seeing the full pipeline succeed end-to-end makes it all worth it!

## Day 24 - Vitest Integration & Branch Protection Rules

**Time Spent:** 3 hours

**Focus:** Introducing unit testing to the Svelte frontend and locking down repo workflow habits.

### Progress & Wins
* **Frontend Unit Testing with Vitest:** Integrated [Vitest](https://vitest.dev/) into the frontend setup to test core JavaScript/TypeScript modules and enforce code reliability before deployment.
* **API Service Coverage:** Started writing test cases for `api.ts` to mock backend responses and verify network call handling, error states, and data mapping.

---

### Workflow & Repo Hardening
* **Branch Protection Rules:** Configured rulesets on the main repository branch:
  * **Protected `main`:** Direct pushes to `main` are now blocked.
  * **Pull Request Workflow:** All changes must pass through feature branches and open PRs.
  * **CI Gatekeeper:** Merges require green status checks from the automated CI pipeline (Go tests + ESLint + Vitest).

> **Habit Building:** Enforcing branch protection forces good git discipline—even on solo projects—and ensures broken code never hits the self-hosted deployment runner automatically.

## Day 25 — Git Hooks with Husky & Frontend Test Expansion

**Focus:** Automating pre-commit code formatting to keep Git history clean, and expanding Svelte store/state test coverage.

### Workflow Upgrade: Fixing CI Pain with Husky
Repeatedly running into a frustrating workflow bottleneck on GitHub Actions:

* **The Problem:** I'd commit a feature, push it, and watch the CI pipeline fail simply because I forgot to manually run formatting/linting scripts beforehand. Fixing it required pushing messy extra commits just to clean up formatting (`"fix linting"`).
* **The Solution:** Installed and configured **Husky** to automate pre-commit Git hooks.
* **How It Works:** Before any commit is created locally, Husky automatically executes `pnpm run format` (and lint checks). If the formatting or linting checks fail or break, Husky cancels/rolls back the commit before it ever hits remote tracking. 

> **Result:** No more failed CI builds due to whitespace or code style mistakes, and a completely clean commit history!

---

### Frontend Testing Progress (Vitest)
Significantly expanded test coverage across the application state layer:

- [x] **`api.ts`** — Network call payloads, headers, and error handling fully covered.
- [x] **`auth.svelte.ts`** — Session state, login/logout reactive state transitions verified.
- [x] **`driveContents.svelte.ts`** — File/folder list mutations and local optimistic updates fully tested.
- [ ] **`navigation.svelte.ts`** — Breadcrumb path parsing and directory traversal logic *(In Progress)*.