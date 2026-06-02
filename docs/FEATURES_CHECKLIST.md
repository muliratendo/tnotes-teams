# TNotes Teams Web — Feature Checklist & Implementation Status

This checklist tracks the implementation status of features based on the PRD.

---

## Phase 1: Foundation (Backend Setup & Database)

| ID | Requirement | Priority | Status | Notes |
| :--- | :--- | :--- | :---: | :--- |
| **FR-1** | **Board Management** | | | |
| FR-1.1 | Create boards with names, descriptions, and themes | P0 | 🔄 Scaffolded | Ent schema & repo defined |
| FR-1.2 | Delete boards (Admin-only) with confirmation | P0 | 🔄 Scaffolded | Ent schema & repo defined |
| FR-1.3 | Board setting changes (color, archived, title) | P1 | 🔄 Scaffolded | Ent schema & repo defined |
| FR-1.4 | Board switching / sidebar navigation | P2 | ❌ Pending | |
| **FR-2** | **List & Card Operations** | | | |
| FR-2.1 | Create unlimited lists per board | P0 | ✅ Implemented | `POST /api/boards/{boardID}/lists` |
| FR-2.2 | Create cards with title, description, due date, labels | P0 | 🔄 Scaffolded | Ent schema & repo defined |
| FR-2.3 | Drag-and-drop cards between lists | P0 | ❌ Pending | Frontend implementation needed |
| FR-2.4 | Drag-and-drop cards within lists (reordering) | P0 | ❌ Pending | Backend normalization + Frontend DnD |
| FR-2.5 | Delete cards (Member+) | P0 | 🔄 Scaffolded | |
| FR-2.6 | Duplicate cards with subtasks | P1 | ❌ Pending | |
| **FR-7** | **Permission System** | | | |
| FR-7.1 | Role assignment per workspace / board | P0 | 🔄 Scaffolded | Role enum defined in join schema |
| FR-7.2 | Invite members via email with role selection | P0 | ❌ Pending | API handler needed |
| FR-7.3 | Role change audit logs | P1 | ❌ Pending | |
| FR-7.4 | Permission enforcement on backend middleware | P0 | 🔄 Scaffolded | `RequireRole` middleware initialized |

---

## Phase 2: Card Operations, Quick Notes & Comments

| ID | Requirement | Priority | Status | Notes |
| :--- | :--- | :--- | :---: | :--- |
| **FR-3** | **Card Detail View** | | | |
| FR-3.1 | Drawer/Modal panel with full details | P0 | ❌ Pending | |
| FR-3.2 | Unified Quick Notes (embedded in card check views) | P0 | 🔄 Scaffolded | Ent Schema defined |
| FR-3.3 | Progress Bar auto-computed from subtask completion | P0 | ❌ Pending | Service logic needed |
| FR-3.4 | Subtask checklist with completion toggling | P0 | 🔄 Scaffolded | Ent Schema defined |
| FR-3.5 | Comments section with @mentions | P0 | 🔄 Scaffolded | Ent Schema defined |
| FR-3.6 | Text link/attachments support | P1 | ❌ Pending | |
| FR-3.7 | Activity logs / audit trail per board | P2 | 🔄 Scaffolded | Ent Schema defined |

---

## Phase 3: Real-Time Collaboration & Team Engine

| ID | Requirement | Priority | Status | Notes |
| :--- | :--- | :--- | :---: | :--- |
| **FR-4** | **Live Co-Presence System** | | | |
| FR-4.1 | Live peer avatars showing online users | P0 | ❌ Pending | WS + Presence Hub needed |
| FR-4.2 | Dynamic cursor labels identifying editing cards | P0 | ❌ Pending | WS broadcasts needed |
| FR-4.3 | Real-time card movement visualization | P0 | ❌ Pending | |
| FR-4.4 | Real-time checklist/note edit synchronization | P0 | ❌ Pending | |
| FR-4.5 | Real-time comment posting and viewing | P0 | ❌ Pending | |
| FR-4.6 | Typing indicators in comments | P1 | ❌ Pending | |
| **FR-5** | **Team Update Engine** | | | |
| FR-5.1 | Go scheduler background worker for mock actions | P0 | ❌ Pending | Worker skeleton needed |
| FR-5.2 | Reconnection mechanism with backoff | P0 | ❌ Pending | Frontend WebSocket client |
| FR-5.3 | LWW deterministic merge rules | P1 | ❌ Pending | Sync Service |

---

## Phase 4: Offline-First Caching & Sync

| ID | Requirement | Priority | Status | Notes |
| :--- | :--- | :--- | :---: | :--- |
| **FR-6** | **Offline Caching & Sync** | | | |
| FR-6.1 | Toggle-able live sync switch in UI | P0 | ❌ Pending | |
| FR-6.2 | Offline read and write operations | P0 | ❌ Pending | |
| FR-6.3 | Mutation queue stored locally in IndexedDB | P0 | ❌ Pending | Dexie/custom DB queue |
| FR-6.4 | Instant bulk synchronization when back online | P0 | ✅ Implemented | `POST /api/boards/{boardID}/sync` + IndexedDB queue |
| FR-6.5 | Conflict resolution UI / notifications | P1 | ❌ Pending | |
| FR-6.6 | Caching static assets with Service Worker | P0 | ❌ Pending | |
| FR-6.7 | IndexedDB state persistence | P0 | ❌ Pending | |

---

## Phase 5: Data Portability & Visual Experience

| ID | Requirement | Priority | Status | Notes |
| :--- | :--- | :--- | :---: | :--- |
| **FR-8** | **Backup & Export/Import** | | | |
| FR-8.1 | Board export to JSON | P0 | ✅ Implemented | `GET /api/boards/{boardID}/export` |
| FR-8.2 | Copy-to-clipboard backup integration | P0 | ❌ Pending | |
| FR-8.3 | Import board from JSON file/payload | P0 | ❌ Pending | Schema validation check on Go backend |
| FR-8.4 | CSV export for cards | P1 | ✅ Implemented | `GET /api/boards/{boardID}/export/csv` |
| **FR-9** | **Search & Filtering** | | | |
| FR-9.1 | Full-text search across titles, description, notes, comments | P1 | ❌ Pending | |
| FR-9.2 | Filters by labels, assignee, due date, completion | P1 | ❌ Pending | |
| **FR-10**| **UI & Design (Indigo Cyber)** | | | |
| FR-10.1| Material 3 design elements | P0 | ❌ Pending | |
| FR-10.2| Custom warm theme "Indigo Cyber" | P0 | ❌ Pending | Primary `#6366F1` with dark accents |
| FR-10.3| Edge-to-edge full viewport design | P0 | ❌ Pending | |
| FR-10.4| Dynamic dark/light mode toggle | P0 | ❌ Pending | |
| FR-10.5| Spring drag-and-drop animation mechanics | P0 | ❌ Pending | Framer Motion physics |

---

## DevOps & Non-Functional Verification

| ID | Requirement | Status | Notes |
| :--- | :--- | :---: | :--- |
| **NFR-1** | Go Backend Dockerfile (multi-stage) | ✅ Implemented | `deploy/Dockerfile.backend` |
| **NFR-2** | React Frontend Dockerfile (serve static/nginx) | ✅ Implemented | `deploy/Dockerfile.frontend` |
| **NFR-3** | `docker-compose.yml` (PostgreSQL + Backend + Frontend) | ✅ Implemented | `deploy/docker-compose.yml` |
| **NFR-4** | Local Dev instructions and script verification | ✅ Implemented | `README.md` + `Makefile` |
| **NFR-5** | Go Service Unit tests (Permissions, Progress, Sync) | ❌ Pending | |
| **NFR-6** | E2E integration test verification (Playwright/Cypress) | ❌ Pending | |
