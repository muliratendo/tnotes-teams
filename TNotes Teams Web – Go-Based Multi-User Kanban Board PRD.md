# "\# Product Requirements Document (PRD)

## TNotes Teams Web – Go-Based Multi-User Kanban Board

**Version:** 1.0
**Date:** June 1, 2026
**Status:** Draft for Review

***

## 1. Executive Summary

### 1.1 Product Vision

TNotes Teams Web is a **real-time, collaborative Multi-User Kanban Board** application built with **Go (Golang)**. It enables teams to work together with simulated co-presence, offline-first capabilities, and robust permission controls for agile project management.

## 2. Functional Requirements

### 2.1 Core Kanban Features

#### FR-1: Board Management

| ID | Requirement | Priority |
| :-- | :-- | :-- |
| FR-1.1 | Create unlimited boards with custom names and descriptions | P0 |
| FR-1.2 | Delete boards (Admin-only) with confirmation dialog | P0 |
| FR-1.3 | Board-level settings: color themes, visibility, archive | P1 |
| FR-1.4 | Drag-and-drop board switching in sidebar | P2 |

#### FR-2: List \& Card Operations

| ID | Requirement | Priority |
| :-- | :-- | :-- |
| FR-2.1 | Create unlimited lists per board (e.g., "To Do", "In Progress", "Done") | P0 |
| FR-2.2 | Create cards with title, description, due date, labels | P0 |
| FR-2.3 | Drag-and-drop cards between lists with spring animation | P0 |
| FR-2.4 | Drag-and-drop cards within same list (reordering) | P0 |
| FR-2.5 | Delete cards (Member+ permissions) | P0 |
| FR-2.6 | Duplicate cards with all subtasks/notes | P1 |

#### FR-3: Card Detail View

| ID | Requirement | Priority |
| :-- | :-- | :-- |
| FR-3.1 | Modal/panel view with full card details | P0 |
| FR-3.2 | **Unified Quick Notes**: Embedded note fields in checklist view | P0 |
| FR-3.3 | **Linear Percent Progress Bar**: Auto-computed from subtask completion | P0 |
| FR-3.4 | Subtasks with checkbox completion tracking | P0 |
| FR-3.5 | Comments section with @mentions | P0 |
| FR-3.6 | Attachments (file uploads, links) | P1 |
| FR-3.7 | Activity log showing all changes with timestamps | P2 |

### 2.2 Real-Time Collaboration Features

#### FR-4: Live Co-Presence System

| ID | Requirement | Priority | Technical Implementation |
| :-- | :-- | :-- | :-- |
| FR-4.1 | **Live peer avatars** (e.g., @sarah, @alex, @michael) showing online users | P0 | WebSocket presence channel |
| FR-4.2 | **Dynamic cursor labels** highlighting cards being edited | P0 | WebSocket cursor broadcast |
| FR-4.3 | Real-time card movement visualization (see others drag cards) | P0 | WebSocket position updates |
| FR-4.4 | Real-time checklist edit synchronization | P0 | WebSocket CRUD operations |
| FR-4.5 | Real-time comment posting and viewing | P0 | WebSocket message broadcast |
| FR-4.6 | Typing indicators in comments ("@sarah is typing...") | P1 | WebSocket typing events |

#### FR-5: Team Update Engine (Production Version)

| ID | Requirement | Priority | Notes |
| :-- | :-- | :-- | :-- |
| FR-5.1 | Replace Android's simulated updates with **real WebSocket/Pusher** | P0 | Uses Go WebSocket server |
| FR-5.2 | Automatic reconnection with exponential backoff | P0 | Handle network failures |
| FR-5.3 | Conflict resolution: Last-write-wins with server timestamp | P1 | Prevent data loss |
| FR-5.4 | Message queue for offline actions (sync on reconnect) | P0 | IndexedDB + Service Worker |

### 2.3 Offline-First Architecture

#### FR-6: Offline Caching \& Sync

| ID | Requirement | Priority |
| :-- | :-- | :-- |
| FR-6.1 | **Toggle-able live sync switch** (enable/disable real-time sync) | P0 |
| FR-6.2 | Full offline mode: read/write without internet | P0 |
| FR-6.3 | **Offline Cached Queue**: Queue all mutations locally | P0 |
| FR-6.4 | Instant synchronization when back online | P0 |
| FR-6.5 | Conflict resolution UI for simultaneous offline edits | P1 |
| FR-6.6 | Service Worker for caching static assets + API responses | P0 |
| FR-6.7 | IndexedDB for local board/card/storage | P0 |

### 2.4 Permission System

#### FR-7: Team Permissions Contexts

| Role | Permissions | UI Behavior |
| :-- | :-- | :-- |
| **Admin** | Full workspace command: create/edit/delete boards, lists, cards, manage members, delete workspace | All fields editable |
| **Member** | Write permissions: update cards, comment, complete subtasks, move cards | Write fields enabled |
| **Viewer** | Read-only: view boards, cards, comments | Write fields disabled + **informative popup snackbars** explaining restrictions |

| ID | Requirement | Priority |
| :-- | :-- | :-- |
| FR-7.1 | Role assignment per workspace/board | P0 |
| FR-7.2 | Invite members via email with role selection | P0 |
| FR-7.3 | Role change audit log | P1 |
| FR-7.4 | Permission enforcement on backend (Go middleware) | P0 |

### 2.5 Data Management

#### FR-8: Backup \& Export/Import

| ID | Requirement | Priority | Format |
| :-- | :-- | :-- | :-- |
| FR-8.1 | **Data Backup Exporter**: Generate JSON of active boards | P0 | JSON |
| FR-8.2 | Copy-paste ready export (clipboard integration) | P0 | – |
| FR-8.3 | **Importer**: Restore boards from JSON backup | P0 | JSON |
| FR-8.4 | Spreadsheet-compatible export (CSV for cards) | P1 | CSV |
| FR-8.5 | Full workspace backup (all boards + members) | P2 | JSON |
| FR-8.6 | Automatic daily cloud backup (admin configurable) | P2 | – |

#### FR-9: Search \& Filtering

| ID | Requirement | Priority |
| :-- | :-- | :-- |
| FR-9.1 | Full-text search across card titles, descriptions, notes, comments | P1 |
| FR-9.2 | Filter by labels, assignee, due date, completion status | P1 |
| FR-9.3 | Saved filter presets | P2 |

### 2.6 Visual \& UI Requirements

#### FR-10: Material 3 Visual Workspace

| ID | Requirement | Priority |
| :-- | :-- | :-- |
| FR-10.1 | **Material 3 design system** with custom theme | P0 |
| FR-10.2 | **Custom warm theme: "Indigo Cyber"** (matching Android) | P0 |
| FR-10.3 | **Edge-to-edge support** (full viewport usage) | P0 |
| FR-10.4 | **Custom dynamic colors** (light/dark mode + theme switching) | P0 |
| FR-10.5 | **Shadow elevations** for cards/lists (depth hierarchy) | P0 |
| FR-10.6 | **Spring animation parameters** for drag-and-drop | P0 |
| FR-10.7 | Responsive design: desktop (primary), tablet, mobile | P0 |
| FR-10.8 | Keyboard shortcuts (e.g., `N` for new card, `Delete` for remove) | P1 |


***

## 3. Non-Functional Requirements

### 3.1 Performance

| ID | Requirement | Target |
| :-- | :-- | :-- |
| NFR-1 | Page load time (first contentful paint) | < 2s on 3G |
| NFR-2 | WebSocket message latency (real-time updates) | < 100ms |
| NFR-3 | Drag-and-drop responsiveness | < 16ms frame rate |
| NFR-4 | Support boards with 500+ cards, 50+ lists | – |
| NFR-5 | Support 100+ concurrent users per board | – |

### 3.2 Security

| ID | Requirement | Implementation |
| :-- | :-- | :-- |
| NFR-6 | Authentication (JWT/OAuth2) | Go + Golang-jwt |
| NFR-7 | HTTPS everywhere | TLS 1.3 |
| NFR-8 | SQL injection prevention | Prepared statements (Gorm/sqlx) |
| NFR-9 | XSS prevention | Content Security Policy + input sanitization |
| NFR-10 | CSRF protection | SameSite cookies + CSRF tokens |
| NFR-11 | Role-based access control (RBAC) | Go middleware |
| NFR-12 | API rate limiting | Go + Redis |

### 3.3 Reliability \& Availability

| ID | Requirement | Target |
| :-- | :-- | :-- |
| NFR-13 | Uptime | 99.9% |
| NFR-14 | Data durability | 99.999999999% (11 nines) |
| NFR-15 | Automatic failover | < 30s RTO |
| NFR-16 | Backup retention | 30 days |

### 3.4 Scalability

| ID | Requirement | Architecture |
| :-- | :-- | :-- |
| NFR-17 | Horizontal scaling | Go microservices + Kubernetes |
| NFR-18 | WebSocket scaling | Redis pub/sub for fan-out |
| NFR-19 | Database scaling | PostgreSQL read replicas |
| NFR-20 | CDN for static assets | Cloudflare/AWS CloudFront |

### 3.5 Testing \& Quality

| ID | Requirement | Coverage Target |
| :-- | :-- | :-- |
| NFR-21 | Unit tests (Go) | 90%+ code coverage |
| NFR-22 | Integration tests (API + DB) | 85%+ coverage |
| NFR-23 | End-to-end tests (Playwright/Cypress) | Critical user flows |
| NFR-24 | WebSocket connection tests | All reconnection scenarios |
| NFR-25 | Offline mode tests | All sync/queue scenarios |
| NFR-26 | Screenshot/visual regression tests | All key screens |
| NFR-27 | Load testing (100 concurrent users) | Pass criteria |


***

## 4. Technical Architecture

### 4.1 Technology Stack

#### Backend (Go)

| Layer | Technology | Rationale |
| :-- | :-- | :-- |
| **Language** | Go 1.21+ | Performance, concurrency, type safety |
| **Web Framework** | Gin Echo or Fiber | High-performance HTTP routing |
| **WebSocket** | Gorilla WebSocket | Production-grade WebSocket library |
| **ORM** | Gorm or sqlx | Database abstraction (PostgreSQL) |
| **Database** | PostgreSQL 15+ | ACID compliance, JSONB support |
| **Cache** | Redis | Session storage + WebSocket pub/sub |
| **Authentication** | Golang-jwt + bcrypt | JWT tokens + password hashing |
| **Queue** | RabbitMQ or Redis Streams | Offline sync queue processing |
| **File Storage** | AWS S3 or MinIO | Card attachments |
| **Search** | PostgreSQL FTS or Meilisearch | Full-text search |

#### Frontend (React + TypeScript)

| Layer | Technology | Rationale |
| :-- | :-- | :-- |
| **Framework** | React 18+ | Component model, ecosystem |
| **Language** | TypeScript 5+ | Type safety (matching Go's strictness) |
| **State Management** | Zustand or Jotai | Lightweight, async support |
| **UI Library** | MUI (Material-UI) or Chakra UI | Material 3 support |
| **Drag-and-Drop** | @dnd-kit/core | Modern, accessible DnD |
| **Real-Time** | Socket.io-client | WebSocket client with reconnection |
| **Offline** | Dexie.js (IndexedDB wrapper) | Local caching |
| **Service Worker** | Workbox | Offline asset caching |
| **Testing** | Vitest + Playwright | Unit + E2E testing |

#### DevOps \& Infrastructure

| Component | Technology |
| :-- | :-- |
| **Containerization** | Docker + Docker Compose |
| **Orchestration** | Kubernetes (minikube for dev, EKS/GKE for prod) |
| **CI/CD** | GitHub Actions |
| **Monitoring** | Prometheus + Grafana |
| **Logging** | ELK Stack (Elasticsearch, Logstash, Kibana) |
| **Deployment** | Vercel (frontend) + Render/AWS (backend) |

### 4.2 System Architecture Diagram

```
┌─────────────────────────────────────────────────────────────────┐
│                        CLIENT (Browser)                         │
│  ┌───────────────────────────────────────────────────────────┐ │
│  │  React + TypeScript (Material 3 UI, Spring Animations)    │ │
│  │  ├─ Zustand (State)                                       │ │
│  │  ├─ Socket.io-client (WebSocket)                          │ │
│  │  ├─ Dexie.js (IndexedDB Offline Cache)                    │ │
│  │  └─ Service Worker (Workbox)                              │ │
│  └───────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────────┘
                              │
                              │ HTTPS + WebSocket
                              │
┌─────────────────────────────────────────────────────────────────┐
│                     LOAD BALANCER (Nginx)                       │
└─────────────────────────────────────────────────────────────────┘
                              │
                              │
        ┌─────────────────────┴─────────────────────┐
        │                                           │
┌───────▼────────┐                         ┌───────▼────────┐
│   Go API       │                         │  Go WebSocket  │
│   Server       │                         │   Server       │
│   (Gin/Echo)   │                         │ (Gorilla WS)   │
│                │                         │                │
│ ├─ Auth        │                         │ ├─ Presence    │
│ ├─ Boards      │                         │ ├─ Cursors     │
│ ├─ Cards       │                         │ ├─ Real-time   │
│ ├─ Comments    │                         │ │   Updates    │
│ └─ Permissions │                         │ └─ Reconnect   │
└───────┬────────┘                         └───────┬────────┘
        │                                           │
        └─────────────────────┬─────────────────────┘
                              │
                    ┌─────────▼─────────┐
                    │     Redis         │
                    │  ├─ Sessions      │
                    │  ├─ Cache         │
                    │  └─ Pub/Sub (WS)  │
                    └─────────┬─────────┘
                              │
                    ┌─────────▼─────────┐
                    │   PostgreSQL      │
                    │  ├─ Users         │
                    │  ├─ Workspaces    │
                    │  ├─ Boards        │
                    │  ├─ Lists         │
                    │  ├─ Cards         │
                    │  ├─ Subtasks      │
                    │  ├─ Comments      │
                    │  └─ Audit Logs    │
                    └───────────────────┘
```


### 4.3 Database Schema (PostgreSQL)

#### Core Tables

```sql
-- Users table
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    username VARCHAR(50) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    avatar_url VARCHAR(500),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);


-- Workspaces (teams)
CREATE TABLE workspaces (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    description TEXT,
    theme_color VARCHAR(7) DEFAULT '#6366F1', -- Indigo Cyber
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMP DEFAULT NOW()
);


-- Workspace members with roles
CREATE TABLE workspace_members (
    workspace_id UUID REFERENCES workspaces(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    role VARCHAR(20) NOT NULL CHECK (role IN ('admin', 'member', 'viewer')),
    joined_at TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (workspace_id, user_id)
);


-- Boards
CREATE TABLE boards (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workspace_id UUID REFERENCES workspaces(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    color_theme VARCHAR(7) DEFAULT '#6366F1',
    is_archived BOOLEAN DEFAULT FALSE,
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);


-- Lists (columns)
CREATE TABLE lists (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    board_id UUID REFERENCES boards(id) ON DELETE CASCADE,
    title VARCHAR(100) NOT NULL,
    position INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);


-- Cards
CREATE TABLE cards (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    list_id UUID REFERENCES lists(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    position INTEGER NOT NULL,
    due_date TIMESTAMP,
    progress_percentage INTEGER DEFAULT 0,
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);


-- Subtasks
CREATE TABLE subtasks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    card_id UUID REFERENCES cards(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    is_completed BOOLEAN DEFAULT FALSE,
    position INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);


-- Quick Notes (embedded in cards)
CREATE TABLE quick_notes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    card_id UUID REFERENCES cards(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);


-- Comments
CREATE TABLE comments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    card_id UUID REFERENCES cards(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id),
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);


-- Activity Log (audit trail)
CREATE TABLE activity_log (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    board_id UUID REFERENCES boards(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id),
    action VARCHAR(50) NOT NULL, -- 'card_moved', 'card_created', 'subtask_completed'
    entity_type VARCHAR(50), -- 'card', 'list', 'board'
    entity_id UUID,
    old_value JSONB,
    new_value JSONB,
    created_at TIMESTAMP DEFAULT NOW()
);


-- Indexes for performance
CREATE INDEX idx_cards_list_id ON cards(list_id);
CREATE INDEX idx_cards_position ON cards(list_id, position);
CREATE INDEX idx_subtasks_card_id ON subtasks(card_id);
CREATE INDEX idx_workspace_members_user ON workspace_members(user_id);
CREATE INDEX idx_activity_log_board ON activity_log(board_id, created_at DESC);
```


### 4.4 API Design (REST + WebSocket)

#### REST API Endpoints

| Method | Endpoint | Description | Auth |
| :-- | :-- | :-- | :-- |
| POST | `/api/auth/register` | User registration | Public |
| POST | `/api/auth/login` | User login | Public |
| GET | `/api/auth/me` | Current user profile | JWT |
| GET | `/api/workspaces` | List user's workspaces | JWT |
| POST | `/api/workspaces` | Create workspace | JWT |
| GET | `/api/workspaces/:id` | Get workspace details | JWT |
| POST | `/api/workspaces/:id/members` | Invite member | Admin |
| PUT | `/api/workspaces/:id/members/:userId` | Update role | Admin |
| GET | `/api/workspaces/:id/boards` | List boards | JWT |
| POST | `/api/workspaces/:id/boards` | Create board | JWT (Member+) |
| DELETE | `/api/boards/:id` | Delete board | Admin |
| GET | `/api/boards/:id` | Get board with all data | JWT |
| POST | `/api/boards/:id/lists` | Create list | JWT (Member+) |
| PUT | `/api/cards/:id` | Update card (move, edit) | JWT (Member+) |
| DELETE | `/api/cards/:id` | Delete card | JWT (Member+) |
| POST | `/api/cards/:id/subtasks` | Add subtask | JWT (Member+) |
| PUT | `/api/subtasks/:id` | Toggle subtask completion | JWT (Member+) |
| POST | `/api/cards/:id/notes` | Add quick note | JWT (Member+) |
| POST | `/api/cards/:id/comments` | Add comment | JWT (Member+) |
| POST | `/api/boards/:id/export` | Export board to JSON | JWT |
| POST | `/api/workspaces/:id/import` | Import boards from JSON | JWT (Admin) |

#### WebSocket Events

| Event | Direction | Payload | Description |
| :-- | :-- | :-- | :-- |
| `join_board` | Client→Server | `{ board_id }` | Join board channel |
| `leave_board` | Client→Server | `{ board_id }` | Leave board channel |
| `user_present` | Server→Client | `{ user_id, username, avatar }` | User came online |
| `user_absent` | Server→Client | `{ user_id }` | User went offline |
| `cursor_update` | Client→Server | `{ board_id, card_id, x, y }` | Cursor position |
| `cursor_broadcast` | Server→Client | `{ user_id, username, card_id, x, y }` | Show other's cursor |
| `card_moved` | Client→Server | `{ card_id, new_list_id, new_position }` | Card drag drop |
| `card_moved_broadcast` | Server→Client | `{ card_id, new_list_id, new_position, user_id }` | Show card move |
| `card_updated` | Client→Server | `{ card_id, field, value }` | Card field edit |
| `card_updated_broadcast` | Server→Client | `{ card_id, field, value, user_id }` | Sync card update |
| `subtask_toggled` | Client→Server | `{ subtask_id, is_completed }` | Toggle subtask |
| `subtask_broadcast` | Server→Client | `{ subtask_id, is_completed, user_id }` | Sync subtask |
| `comment_added` | Client→Server | `{ card_id, content }` | Post comment |
| `comment_broadcast` | Server→Client | `{ comment_id, card_id, content, user_id }` | Show comment |
| `typing_start` | Client→Server | `{ card_id }` | User typing in comments |
| `typing_broadcast` | Server→Client | `{ user_id, card_id }` | Show "is typing" |
| `offline_queue_sync` | Client→Server | `{ actions: [] }` | Sync offline actions |
| `sync_ack` | Server→Client | `{ synced_count, conflicts }` | Confirm sync |

### 4.5 Go Backend Structure

```
tnotes-teams-web/
├── cmd/
│   └── server/
│       └── main.go                 # Entry point
├── internal/
│   ├── config/
│   │   └── config.go               # Environment config
│   ├── database/
│   │   ├── postgres.go             # DB connection
│   │   └── migrations/             # SQL migrations
│   ├── models/
│   │   ├── user.go                 # User struct
│   │   ├── workspace.go            # Workspace struct
│   │   ├── board.go                # Board struct
│   │   ├── card.go                 # Card struct
│   │   └── subtask.go              # Subtask struct
│   ├── repository/
│   │   ├── user_repo.go            # User CRUD
│   │   ├── board_repo.go           # Board CRUD
│   │   └── card_repo.go            # Card CRUD
│   ├── service/
│   │   ├── auth_service.go         # JWT auth logic
│   │   ├── board_service.go        # Board business logic
│   │   └── sync_service.go         # Offline sync logic
│   ├── handler/
│   │   ├── auth_handler.go         # HTTP auth endpoints
│   │   ├── board_handler.go        # HTTP board endpoints
│   │   └── card_handler.go         # HTTP card endpoints
│   ├── websocket/
│   │   ├── hub.go                  # WebSocket hub (manages connections)
│   │   ├── client.go               # WebSocket client struct
│   │   └── events.go               # Event handling
│   ├── middleware/
│   │   ├── auth.go                 # JWT verification
│   │   ├── rbac.go                 # Role-based access control
│   │   ├── cors.go                 # CORS handling
│   │   └── rate_limit.go           # Rate limiting
│   └── utils/
│       ├── jwt.go                  # JWT helpers
│       ├── hash.go                 # Password hashing
│       └── response.go             # API response helpers
├── pkg/
│   └── api/
│       └── routes.go               # Route definitions
├── frontend/
│   ├── src/
│   │   ├── components/             # React components
│   │   ├── hooks/                  # Custom hooks
│   │   ├── stores/                 # Zustand stores
│   │   ├── services/               # API + WebSocket clients
│   │   └── utils/                  # Helpers
│   ├── package.json
│   └── tsconfig.json
├── Dockerfile
├── docker-compose.yml
├── Makefile
├── go.mod
└── README.md
```


***

## 5. User Interface Mockups

### 5.1 Key Screens

#### Screen 1: Board View (Primary Interface)

```
┌─────────────────────────────────────────────────────────────────────┐
│  🔵 TNotes Teams                                    [@user] [⚙️]   │
├─────────────────────────────────────────────────────────────────────┤
│  🏠 Workspace  >  📋 Project Alpha  >  🎯 Sprint Board             │
├────────────┬──────────────┬──────────────┬──────────────┬──────────┤
│  TO DO     │ IN PROGRESS  │    REVIEW    │     DONE     │ [+]      │
│  [12]      │    [5]       │    [3]       │    [18]      │ List     │
├────────────┼──────────────┼──────────────┼──────────────┼──────────┤
│ 📝 Card 1  │ 📝 Card 5    │ 📝 Card 9    │ 📝 Card 15   │          │
│   @sarah👤 │   @alex👤    │   @michael👤 │              │          │
│   [▬▬▬▬▬]  │   [▬▬▬▬▬▬▬] │   [▬▬▬▬▬▬▬▬] │  [▬▬▬▬▬▬▬▬▬] │          │
│   60%      │   80%        │   100%       │   100%       │          │
├────────────┼──────────────┼──────────────┼──────────────┼──────────┤
│ 📝 Card 2  │ 📝 Card 6    │ 📝 Card 10   │ 📝 Card 16   │          │
│   @alex👤  │            │              │              │          │
│   [▬▬▬▬▬▬] │   [▬▬▬▬▬▬▬] │              │              │          │
│   75%      │   85%        │              │              │          │
├────────────┼──────────────┼──────────────┼──────────────┼──────────┤
│ 📝 Card 3  │ 📝 Card 7    │              │ 📝 Card 17   │          │
│            │   @sarah👤   │              │              │          │
│   [▬▬▬▬▬]  │   [▬▬▬▬▬▬]  │              │              │          │
│   50%      │   70%        │              │              │          │
└────────────┴──────────────┴──────────────┴──────────────┴──────────┘


Legend:
👤 = Active co-presence cursor (highlighted card being edited)
[▬▬▬▬▬] = Linear progress bar (auto-computed from subtasks)
@sarah, @alex = Avatar labels showing who's actively editing
```


#### Screen 2: Card Detail Modal

```
┌─────────────────────────────────────────────────────────────────┐
│  📝 Implement WebSocket Reconnection                    [✕]     │
├─────────────────────────────────────────────────────────────────┤
│  Description:                                                   │
│  Implement exponential backoff WebSocket reconnection logic    │
│  with offline queue sync                                        │
│                                                                 │
│  📅 Due: June 15, 2026    🏷️ [Backend] [Go] [High Priority]    │
├─────────────────────────────────────────────────────────────────┤
│  PROGRESS: ████████████░░ 80% (4/5 subtasks)                   │
├─────────────────────────────────────────────────────────────────┤
│  SUBTASKS:                                                      │
│  ☑️ Design WebSocket hub architecture                          │
│  ☑️ Implement client reconnection logic                        │
│  ☑️ Add exponential backoff (1s, 2s, 4s, 8s)                   │
│  ☑️ Build offline queue with IndexedDB                         │
│  ☐ Write unit tests for reconnection                           │
├─────────────────────────────────────────────────────────────────┤
│  QUICK NOTES:                                                   │
│  ┌───────────────────────────────────────────────────────────┐ │
│  │ Remember to test with slow network (Chrome DevTools)     │ │
│  │ Use Redis pub/sub for WebSocket fan-out                  │ │
│  └───────────────────────────────────────────────────────────┘ │
│  [+] Add Note                                                   │
├─────────────────────────────────────────────────────────────────┤
│  COMMENTS:                                                      │
│  ┌─ @alex (2 min ago)                                          │
│  │ 👍 Good approach! Consider adding message deduplication    │
│  └───────────────────────────────────────────────────────────┐ │
│  │ Write a comment...                           [@] [📎] [💬] │
│  └───────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────────┘
```


#### Screen 3: Permission Popup (Viewer Attempting Edit)

```
┌─────────────────────────────────────────────────────────────┐
│  ⚠️ Permission Denied                                       │
├─────────────────────────────────────────────────────────────┤
│  You don't have permission to edit this card.               │
│                                                             │
│  Your role: **Viewer** (Read-only)                          │
│                                                             │
│  To edit cards, request **Member** or **Admin** access      │
│  from workspace administrators.                             │
│                                                             │
│                        [Cancel]  [Request Access]           │
└─────────────────────────────────────────────────────────────┘
```


***

## 6. Development Roadmap

### Phase 1: MVP (Weeks 1–4)

**Goal:** Core Kanban functionality with real-time collaboration


| Week | Deliverables |
| :-- | :-- |
| 1 | Go backend setup, PostgreSQL schema, auth (JWT), user/workspace CRUD |
| 2 | Board/list/card REST APIs, Gorm models, basic frontend (React + TypeScript) |
| 3 | WebSocket implementation (Gorilla WebSocket), real-time card moves, co-presence cursors |
| 4 | Material 3 UI (MUI), drag-and-drop (@dnd-kit), progress bar, MVP testing |

**MVP Success Criteria:**

- Create boards, lists, cards
- Drag-and-drop cards with real-time sync
- See other users' cursors and avatars
- JWT authentication working


### Phase 2: Offline-First \& Permissions (Weeks 5–7)

**Goal:** Offline support and role-based access control


| Week | Deliverables |
| :-- | :-- |
| 5 | Service Worker + IndexedDB (Dexie.js), offline queue, sync on reconnect |
| 6 | Permission system (Admin/Member/Viewer), RBAC middleware, UI restrictions |
| 7 | Quick Notes, comments, subtasks with progress bar, testing |

**Success Criteria:**

- Full offline mode (read/write without internet)
- Offline queue syncs instantly on reconnect
- Permission enforcement on frontend + backend


### Phase 3: Data Management \& Polish (Weeks 8–9)

**Goal:** Backup/export/import, search, visual polish


| Week | Deliverables |
| :-- | :-- |
| 8 | JSON export/import, CSV export, automatic backups |
| 9 | Full-text search, filters, dark mode, spring animations, bug fixes |

**Success Criteria:**

- Export/import boards as JSON
- Search across all card content
- Material 3 theme "Indigo Cyber" complete


### Phase 4: Testing \& Deployment (Weeks 10–12)

**Goal:** Production readiness


| Week | Deliverables |
| :-- | :-- |
| 10 | Unit tests (90%+ coverage), integration tests, WebSocket tests |
| 11 | E2E tests (Playwright), load testing, security audit |
| 12 | Docker deployment, CI/CD pipeline, documentation, launch |

**Success Criteria:**

- All tests passing (100% Green like Android version)
- 99.9% uptime in staging
- Deployed to production (Vercel + Render/AWS)

***

## 7. Testing Strategy

### 7.1 Test Pyramid

```
                    ┌─────────────┐
                    │  E2E Tests  │  ← Playwright (critical user flows)
                    │    10%      │
                   ─┴─────────────┴─
                 ┌───────────────────┐
                 │ Integration Tests │  ← API + DB + WebSockets
                 │      30%          │
                ─┴───────────────────┴─
              ┌─────────────────────────┐
              │    Unit Tests (Go)      │  ← Business logic, handlers
              │         60%             │
             ─┴─────────────────────────┴─
```


### 7.2 Test Coverage Targets

| Test Type | Framework | Coverage Target | Examples |
| :-- | :-- | :-- | :-- |
| **Unit** | Go `testing`, Vitest | 90%+ | Auth service, RBAC middleware, progress calculation |
| **Integration** | Go `testcontainers`, Playwright | 85%+ | Board CRUD with DB, WebSocket handshake |
| **E2E** | Playwright | Critical flows | Login → Create board → Add card → Drag-and-drop → Offline sync |
| **WebSocket** | Gorilla WebSocket tests | 100% reconnection | Reconnect after disconnect, offline queue sync |
| **Offline** | Playwright + Service Worker | All scenarios | Read/write offline, sync on reconnect, conflict resolution |
| **Screenshot** | Playwright visual regression | All key screens | Board view, card modal, permission popup |

### 7.3 Example Go Unit Test

```go
// internal/service/board_service_test.go
package service


import (
    "testing"
    "github.com/muliratendo/tnotes-teams-web/internal/models"
)


func TestCalculateProgress(t *testing.T) {
    tests := []struct {
        name     string
        subtasks []models.Subtask
        expected int
    }{
        {
            name: "all completed",
            subtasks: []models.Subtask{
                {Title: "Task 1", IsCompleted: true},
                {Title: "Task 2", IsCompleted: true},
            },
            expected: 100,
        },
        {
            name: "half completed",
            subtasks: []models.Subtask{
                {Title: "Task 1", IsCompleted: true},
                {Title: "Task 2", IsCompleted: false},
            },
            expected: 50,
        },
        {
            name: "no subtasks",
            subtasks: []models.Subtask{},
            expected: 0,
        },
    }


    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := CalculateProgress(tt.subtasks)
            if result != tt.expected {
                t.Errorf("expected %d, got %d", tt.expected, result)
            }
        })
    }
}
```


***

## 8. Success Metrics \& KPIs

| Metric | Target | Measurement |
| :-- | :-- | :-- |
| **User Adoption** | 100+ active users in first month | Google Analytics |
| **Real-time Latency** | < 100ms WebSocket message | Prometheus metrics |
| **Offline Sync Success** | 99%+ offline actions sync correctly | Error logging |
| **Test Coverage** | 90%+ backend, 85%+ frontend | `go test -cover` |
| **Uptime** | 99.9% | UptimeRobot |
| **Page Load Time** | < 2s on 3G | Lighthouse |
| **Bug Resolution** | < 48h for P0 bugs | GitHub Issues |


***

## 9. Risks \& Mitigations

| Risk | Impact | Probability | Mitigation |
| :-- | :-- | :-- | :-- |
| WebSocket scaling issues | High | Medium | Redis pub/sub, load testing early |
| Offline conflict resolution | High | Medium | Last-write-wins + conflict UI |
| Go WebSocket memory leaks | Medium | Low | Connection pooling, monitoring |
| Material 3 theme mismatch | Low | Low | Pixel-perfect design comparisons |
| PostgreSQL performance | Medium | Medium | Indexes, read replicas, query optimization |
| Frontend state complexity | Medium | Medium | Zustand for simple state, clear architecture |


***

## 10. Appendices

### Appendix A: Glossary

| Term | Definition |
| :-- | :-- |
| **Co-presence** | Real-time visibility of other users' activity (cursors, avatars) |
| **Offline-first** | App works fully without internet; syncs when online |
| **Spring animation** | Physics-based animation for smooth drag-and-drop |
| **Indigo Cyber** | Custom warm theme color scheme (\#6366F1 primary) |
| **RBAC** | Role-Based Access Control (Admin/Member/Viewer) |
| **IndexedDB** | Browser's local database for offline caching |

### Appendix B: Android ↔ Web Feature Mapping

| Android Feature | Technology | Web Equivalent | Technology |
| :-- | :-- | :-- | :-- |
| Room Database | SQLite + KSP | PostgreSQL + Gorm | Go ORM |
| Kotlin Coroutines | Async state | Goroutines + Channels | Go concurrency |
| Flow (reactive) | State streams | WebSocket events | Real-time push |
| Jetpack Compose | Material 3 UI | React + MUI | Material 3 Web |
| ViewModels | State managers | Zustand stores | React state |
| Team Update Engine | Simulated sync | WebSocket server | Gorilla WebSocket |
| Offline Queue | Local caching | IndexedDB + Service Worker | Dexie.js + Workbox |

### Appendix C: References

- Go WebSocket: [Gorilla WebSocket](https://github.com/gorilla/websocket)
- Material 3 Web: [MUI Material-UI](https://mui.com/)
- Drag-and-Drop: [@dnd-kit](https://dndKit.com/)
- IndexedDB: [Dexie.js](https://dexie.org/)
- Go ORM: [Gorm](https://gorm.io/)"
Regenerate this prd, using the information below;
"\# Project Concepts \& Architecture Overview
- Database Schema \& Entity Modeling: Managed via Ent, utilizing an AST-based code-generation engine for 100% type-safe schema modeling.
- Compile-Time Repository Layer: Built with sqlc to automatically parse and validate native PostgreSQL queries during the application build phase.
- Architectural Decoupling: Structurally organized around the clean Repository Pattern to completely isolate underlying PostgreSQL operations from core business logic.
- State Management \& Offline Resilience: Architected using a React + TypeScript state machine and Go background workers, utilizing asynchronous channels and reactive triggers to handle automated query retries and ensure seamless data persistence to PostgreSQL upon network reconnection.
- Creative Minimalist UI: Crafted a cohesive React + TypeScript visual workspace utilizing a themed design system mapped to dynamic PostgreSQL data nodes, styled with a custom, warm "Indigo Cyber" color palette.
- Fluid Frontend Experience: Built an edge-to-edge React + TypeScript user interface styled with custom dynamic themes, leveraging shadow elevations and Framer Motion spring physics to animate database payloads fetched from the PostgreSQL backend.


# Visual \& Functional Polish Implementation

- Real-Time Workspace Presence: Simulates real-time peer collaboration in the React + TypeScript frontend, showcasing live cursor trackers and teammate avatar indicators as active mutations occur within the PostgreSQL database.
- Asynchronous Update Engine: Engineered a background Go scheduler that pipes mocked teammate actions directly into the database layer, perfectly mimicking Pusher or WebSocket live-data synchronization events.
- Offline Synchronization Pipeline: Implemented a localized state queue that tracks offline client mutations and securely syncs pending payloads back to the PostgreSQL instance via Ent and sqlc immediately upon reconnection.
- Network Interruption Simulator: Features a front-end toggle switch to simulate offline states, allowing you to test client-side caching mechanism resilience and queue reconciliation logic on the fly.


# Team Permissions Contexts:

- Admin Permissions Layer: Grants full command over workspace mutations, board creation, and database-level deletions executed securely via the Go backend and Ent client.
- Member Interaction Layer: Enforces write permissions across the Go API to authorize data mutations, comment additions, and task updates inside the PostgreSQL database.
- Viewer Restrictions: Restricts UI components in the React + TypeScript frontend to a read-only state, utilizing informative popup snackbars to block unauthorized write operations.


# Rich Management \& Data Operations:

- Contextual Quick Notes: Integrates secure text fields directly inside React checklist cards to persist spontaneous micro-notes into the PostgreSQL database.
- Dynamic Progress Tracking: Features an automated progress tracker that computes real-time completion percentages based on subtask payloads fetched from the Go backend.
- Data Portability Pipeline: Enables users to generate and restore fully readable JSON payloads of boards for easy spreadsheet migration or full database backup recovery.
- Comprehensive Test Coverage: Backed by robust Go unit testing and React Testing Library suites to validate business logic, schema safety, and UI components

The TNotes Teams application will be rebuilt as a web platform utilizing a Go backend with Ent and sqlc for data modeling and query generation, a PostgreSQL database, and a React with TypeScript frontend."

Here is a regenerated PRD tailored to the revised architecture: **Go backend + PostgreSQL + Ent + sqlc + React + TypeScript**, with offline resilience, reactive state handling, and a warm Indigo Cyber UI. Ent is a type-safe entity framework for Go, sqlc generates Go code from SQL at build time, and Framer Motion supports spring-based UI transitions.[^1][^2][^3]

# Product Requirements Document

## TNotes Teams Web

**Version:** 1.0
**Date:** June 1, 2026
**Status:** Draft for Review

## 1. Executive Summary

### 1.1 Product Vision

TNotes Teams Web is a real-time collaborative Kanban platform for teams that need structured task management, simulated co-presence, offline-first editing, and strong permission controls. The web rebuild preserves the collaborative feel of the Android app while using a Go backend, PostgreSQL, Ent, sqlc, and a React + TypeScript frontend.[^2][^1]

### 1.2 Product Goals

- Enable teams to create and manage boards, lists, cards, subtasks, notes, and comments in one workspace.
- Simulate live collaboration with peer avatars, cursor presence, and background teammate activity.
- Support offline mutation, cached persistence, and automatic sync when the network returns.
- Enforce Admin, Member, and Viewer permissions consistently in the UI and backend.
- Provide export/import workflows for backups and spreadsheet migration.


## 2. Product Scope

### 2.1 In Scope

- Kanban board management, drag-and-drop, card detail views, comments, subtasks, quick notes, and progress tracking.
- Real-time presence, cursor tracking, simulated teammate updates, and reconnect handling.
- Offline queueing, local caching, and sync reconciliation.
- Role-based access control and audit trails.
- JSON export/import and optional CSV export for cards.
- Visual design system with Indigo Cyber theming, edge-to-edge layout, and spring animations. Framer Motion supports spring physics suitable for this motion style.[^3]


### 2.2 Out of Scope for v1

- Native mobile apps.
- Enterprise SSO beyond standard OAuth2/OpenID Connect.
- Multi-region active-active realtime infrastructure.
- Advanced CRDT collaboration beyond the agreed sync model.


## 3. Users and Roles

### 3.1 User Roles

- Admin: full workspace control, board creation/deletion, member management, backup/export/import, and data recovery.
- Member: create and edit cards, update subtasks, comment, move cards, and add quick notes.
- Viewer: read-only access with blocked write actions and clear permission feedback.


### 3.2 Primary Use Cases

- A product team runs sprint planning on a shared board.
- A distributed team edits cards while seeing live presence and cursor activity.
- A project lead exports a board as JSON before making major changes.
- A viewer reviews progress without accidentally mutating data.


## 4. Functional Requirements

### 4.1 Board Management

- FR-1.1: Users can create unlimited boards with names, descriptions, and themes.
- FR-1.2: Admins can delete boards with confirmation.
- FR-1.3: Boards can be archived and restored.
- FR-1.4: Board switching is available from a sidebar or workspace switcher.


### 4.2 List and Card Operations

- FR-2.1: Users can create unlimited lists per board.
- FR-2.2: Users can create cards with title, description, due date, labels, assignees, and metadata.
- FR-2.3: Cards can be dragged between lists with smooth spring motion.
- FR-2.4: Cards can be reordered within a list.
- FR-2.5: Members and Admins can delete cards.
- FR-2.6: Cards can be duplicated, including checklist items, notes, and comments where permitted.


### 4.3 Card Detail View

- FR-3.1: Cards open in a drawer or modal with full details.
- FR-3.2: Quick notes are embedded directly inside the card experience.
- FR-3.3: Progress is computed automatically from checklist completion.
- FR-3.4: Subtasks support checkbox completion and ordering.
- FR-3.5: Comments support @mentions and live timestamps.
- FR-3.6: Attachments support links in v1, file upload optionally in a later phase.
- FR-3.7: Activity log shows all material changes on the card.


### 4.4 Real-Time Collaboration

- FR-4.1: The app shows live peer avatars for online teammates.
- FR-4.2: Cursor labels identify who is actively editing which card.
- FR-4.3: Card movement is reflected in near real time for other users.
- FR-4.4: Checklist edits, note edits, and comment updates propagate live.
- FR-4.5: Typing indicators appear in comments.
- FR-4.6: Presence is simulated when no active external realtime provider is configured.


### 4.5 Team Update Engine

- FR-5.1: The backend runs a Go scheduler that emits mocked teammate actions into the system.
- FR-5.2: The update engine can imitate WebSocket/Pusher-style activity feeds.
- FR-5.3: Reconnection uses exponential backoff.
- FR-5.4: Conflicts resolve by server timestamp with deterministic merge rules.


### 4.6 Offline-First Behavior

- FR-6.1: Users can toggle live sync on or off from the UI.
- FR-6.2: The app supports read and write operations offline.
- FR-6.3: Local mutations are queued in the client state layer and persisted locally.
- FR-6.4: Queued changes sync automatically when connectivity returns.
- FR-6.5: A conflict resolution UI appears when the server detects divergent edits.
- FR-6.6: Static assets are cached for offline use.
- FR-6.7: Local state survives reloads and browser restarts.


### 4.7 Permissions

- FR-7.1: Permissions are assigned per workspace and optionally per board.
- FR-7.2: Admins can invite members by email and assign roles.
- FR-7.3: Role changes are recorded in an audit log.
- FR-7.4: Backend middleware enforces permissions on every mutating endpoint.
- FR-7.5: Viewers cannot edit, and the UI must show a clear snackbar when they attempt to do so.


### 4.8 Data Management

- FR-8.1: Users can export active boards as readable JSON.
- FR-8.2: Export output supports copy-paste and direct download.
- FR-8.3: Users can import boards from JSON backups.
- FR-8.4: CSV export is available for cards and task summaries.
- FR-8.5: Admins can back up a full workspace.
- FR-8.6: Backup jobs can be scheduled later as an admin setting.


### 4.9 Search and Filtering

- FR-9.1: Full-text search works across card titles, descriptions, notes, and comments.
- FR-9.2: Filters support labels, assignees, due dates, and completion status.
- FR-9.3: Saved filter presets are supported in a later release.


### 4.10 Visual and UI Requirements

- FR-10.1: The UI uses a Material 3-inspired design language.
- FR-10.2: The Indigo Cyber theme is the default visual style.
- FR-10.3: The app is edge-to-edge and uses the full viewport.
- FR-10.4: Dynamic dark/light themes are supported.
- FR-10.5: Cards and panels use shadow elevation for hierarchy.
- FR-10.6: Drag-and-drop uses spring animation for natural motion.
- FR-10.7: The layout is responsive, with desktop as the primary target.
- FR-10.8: Keyboard shortcuts support quick creation and deletion actions.


## 5. Non-Functional Requirements

### 5.1 Performance

- NFR-1: Initial contentful paint should be under 2 seconds on a throttled 3G profile.
- NFR-2: Realtime events should arrive within 100 ms under normal conditions.
- NFR-3: Drag interactions should remain at 60 fps on supported devices.
- NFR-4: The app should support boards with 500+ cards and 50+ lists.
- NFR-5: The system should support 100+ concurrent users per board.


### 5.2 Security

- NFR-6: Authentication should support JWT and optionally OAuth2.
- NFR-7: All traffic must use HTTPS.
- NFR-8: SQL injection must be prevented by parameterized queries through sqlc-generated code.
- NFR-9: XSS protection must include sanitization and CSP headers.
- NFR-10: CSRF protections must be applied where cookie auth is used.
- NFR-11: RBAC must be enforced in the backend.
- NFR-12: API rate limiting must be applied for abuse protection.


### 5.3 Reliability and Availability

- NFR-13: Service uptime target is 99.9%.
- NFR-14: Data durability target is 11 nines through PostgreSQL managed backups.
- NFR-15: Recovery time objective for failover is under 30 seconds.
- NFR-16: Retain backups for 30 days minimum.


### 5.4 Scalability

- NFR-17: Backend must scale horizontally.
- NFR-18: Realtime fan-out should be supported through Redis pub/sub or equivalent.
- NFR-19: Database read scaling should use read replicas where necessary.
- NFR-20: Static assets should be served via CDN.


### 5.5 Quality and Testing

- NFR-21: Go unit tests should cover critical business logic.
- NFR-22: Integration tests should cover API, DB, and sync behavior.
- NFR-23: End-to-end tests should cover the main user journeys.
- NFR-24: WebSocket and offline-reconnect scenarios must be tested.
- NFR-25: Screenshot/visual regression tests should cover key screens.
- NFR-26: Load testing should validate collaborative usage under concurrent load.


## 6. Technical Architecture

### 6.1 Backend Stack

- Language: Go 1.22+.
- Database: PostgreSQL.
- Schema modeling: Ent for type-safe entity definitions and code generation. Ent is designed for type-safe data modeling in Go.[^1]
- Query layer: sqlc for compile-time validated SQL and generated Go query code. sqlc generates code from schema and queries during build time.[^2]
- Pattern: Clean Repository Pattern to isolate persistence from business logic.
- Realtime: WebSocket service for presence, cursor updates, and collaboration events.
- Background jobs: Go workers and schedulers for mocked teammate actions, retries, and sync processing.


### 6.2 Frontend Stack

- Framework: React 18+ with TypeScript.
- State: State-machine-based client architecture with async channels and reactive triggers.
- Motion: Framer Motion for spring-based drag and UI transitions. Framer Motion’s spring transitions are appropriate for natural, physics-based movement.[^3]
- Offline storage: IndexedDB-based local cache and mutation queue.
- Service worker: Asset caching and offline support.
- Testing: React Testing Library and Playwright.


### 6.3 Data Flow

- User action is captured in the React state machine.
- The action is written to local cache and queued if offline.
- When online, the queue syncs to the Go API.
- The Go backend validates the request, applies business rules, and persists changes through Ent and sqlc-based repositories.
- Realtime events are broadcast to connected clients.
- The UI updates optimistic state and reconciles with server truth.


## 7. Data Modeling

### 7.1 Modeling Approach

- Ent defines canonical entities, relations, validations, and schema constraints.
- sqlc manages hand-written SQL for reads, joins, reporting, and optimized mutations.
- The repository layer shields the rest of the app from raw SQL and storage details.


### 7.2 Core Entities

- Users.
- Workspaces.
- Workspace members.
- Boards.
- Lists.
- Cards.
- Subtasks.
- Quick notes.
- Comments.
- Attachments.
- Activity logs.
- Sync operations.
- Export jobs.


### 7.3 Entity Rules

- A workspace contains many boards.
- A board contains many lists.
- A list contains many cards.
- A card contains many subtasks, notes, comments, and optional attachments.
- Every mutation must preserve ordering metadata for lists and cards.
- Audit trails must record actor, action, timestamp, before state, and after state.


## 8. Core UI Requirements

### 8.1 Board View

- Sidebar shows workspaces and boards.
- Board canvas shows lists horizontally with card stacks vertically.
- The top bar contains search, sync toggle, and user controls.
- Presence avatars appear near cards being edited.


### 8.2 Card Drawer

- Shows title, description, subtasks, notes, comments, attachments, progress, and activity log.
- Supports inline edits for members and admins.
- Displays read-only messaging for viewers.


### 8.3 Permission Feedback

- Unauthorized actions must not silently fail.
- The UI should explain why an action is blocked.
- Snackbars should be short, clear, and actionable.


### 8.4 Theme

- Indigo Cyber uses warm indigo primaries, dark surfaces, soft highlights, and elevated cards.
- Motion should feel smooth and subtle rather than flashy.
- The layout should feel dense enough for productivity but easy on the eyes.


## 9. API Design

### 9.1 REST Endpoints

- POST /api/auth/register
- POST /api/auth/login
- GET /api/auth/me
- GET /api/workspaces
- POST /api/workspaces
- GET /api/workspaces/:id
- POST /api/workspaces/:id/members
- PUT /api/workspaces/:id/members/:userId
- GET /api/workspaces/:id/boards
- POST /api/workspaces/:id/boards
- GET /api/boards/:id
- PATCH /api/boards/:id
- DELETE /api/boards/:id
- POST /api/boards/:id/lists
- PATCH /api/lists/:id
- DELETE /api/lists/:id
- POST /api/lists/:id/cards
- PATCH /api/cards/:id
- DELETE /api/cards/:id
- POST /api/cards/:id/subtasks
- PATCH /api/subtasks/:id
- POST /api/cards/:id/notes
- POST /api/cards/:id/comments
- GET /api/boards/:id/export
- POST /api/boards/import


### 9.2 Realtime Events

- join_board
- leave_board
- user_present
- user_absent
- cursor_update
- cursor_broadcast
- card_moved
- card_moved_broadcast
- card_updated
- card_updated_broadcast
- subtask_toggled
- subtask_broadcast
- comment_added
- comment_broadcast
- typing_start
- typing_broadcast
- offline_queue_sync
- sync_ack


## 10. Permission Rules

### 10.1 Admin

- Full control over workspace mutations.
- Can create, edit, delete boards and lists.
- Can delete workspace data and manage members.


### 10.2 Member

- Can create and edit cards.
- Can comment and complete subtasks.
- Can move cards and update notes.
- Cannot manage workspace ownership.


### 10.3 Viewer

- Can view boards and cards.
- Cannot mutate any write-protected field.
- Must receive explanatory feedback when attempting edits.


## 11. Offline Sync Strategy

### 11.1 Client Queue

- Store all offline mutations locally.
- Assign each queued operation a local idempotency key.
- Preserve action order for deterministic replay.
- Keep queued items across reloads.


### 11.2 Reconnection Flow

- Detect network restoration.
- Replay operations to the server in order.
- Resolve conflicts using server timestamps and merge rules.
- Refresh affected board slices after sync completes.


### 11.3 Conflict Handling

- Prefer last-write-wins for simple scalar fields.
- Preserve both values in the activity log for disputed edits.
- Prompt the user only for ambiguous collisions.


## 12. Testing Strategy

### 12.1 Go Backend Tests

- Ent schema tests.
- sqlc query validation tests.
- Repository tests against PostgreSQL test containers.
- Service-layer tests for permissions, sync, progress calculation, and export/import.


### 12.2 Frontend Tests

- React Testing Library for components and flows.
- State machine tests for offline and reconnect states.
- UI tests for permission feedback and progress rendering.


### 12.3 End-to-End Tests

- Login, create board, add list, add card, move card, go offline, queue changes, reconnect, verify sync.
- Viewer blocked from editing.
- Export JSON, import JSON, and validate restored board.


### 12.4 Visual Regression

- Board view.
- Card drawer.
- Permission snackbar.
- Offline indicator and sync badge.


## 13. Rollout Plan

### Phase 1: Foundation

- Set up Go backend, PostgreSQL, Ent models, sqlc queries, and auth.
- Build basic React shell and board navigation.


### Phase 2: Core Collaboration

- Add lists, cards, drag-and-drop, comments, subtasks, and progress.
- Add presence and realtime updates.


### Phase 3: Offline and Permissions

- Add offline queue, sync reconciliation, role enforcement, and read-only states.


### Phase 4: Export and Polish

- Add JSON export/import, CSV export, theming, motion polish, and accessibility improvements.


### Phase 5: Hardening

- Add load tests, monitoring, backups, and production deployment.


## 14. Success Metrics

- 100+ active users in the first month.
- Realtime message latency under 100 ms.
- Offline sync success above 99%.
- Backend test coverage above 90%.
- Frontend critical-flow test coverage complete.
- Page load under 2 seconds on slow network profiles.


## 15. Risks and Mitigations

- Realtime scaling risk: use Redis fan-out and shard by board.
- Offline conflict risk: use idempotent operations, timestamps, and a clear resolution UI.
- Schema complexity risk: keep Ent as the source of truth and restrict direct DB writes.
- Query drift risk: use sqlc to validate SQL at build time.
- UI complexity risk: isolate state machine logic and keep components small.


## 16. Implementation Notes

- Ent provides type-safe schema modeling and generation for the domain entities.[^1]
- sqlc should be used for query-heavy reads and reliable compile-time checked SQL paths.[^2]
- Framer Motion should handle drag, presence, and panel transitions with spring physics for a natural collaborative feel.[^3]

I can next turn this into one of three formats:

1. a polished **Markdown PRD document** ready for GitHub,
2. a **feature-by-feature engineering backlog** with epics and tasks,
3. a **system design + API spec** for implementation.
<span style="display:none">[^10][^4][^5][^6][^7][^8][^9]</span>

<div align="center">⁂</div>

[^1]: https://github.com/ent/ent

[^2]: https://docs.sqlc.dev/en/v1.21.0/tutorials/getting-started-postgresql.html

[^3]: https://refine.dev/blog/framer-motion/

[^4]: https://dev.to/jones_charles_ad50858dbc0/getting-started-with-ent-facebooks-open-source-orm-for-go-developers-1non

[^5]: https://pkg.go.dev/github.com/go-ent/ent

[^6]: https://encore.dev/docs/go/how-to/entgo-orm

[^7]: https://doc.ryansu.tech/guide/feature/ORM/ent.html

[^8]: https://www.framer.com/academy/lessons/framer-animations-transitions-and-easing

[^9]: https://pkg.go.dev/github.com/origadmin/toolkits/registry/orm/ent

[^10]: https://docs.sqlc.dev/en/v1.12.0/tutorials/getting-started-postgresql.html

