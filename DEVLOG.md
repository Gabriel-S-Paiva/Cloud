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