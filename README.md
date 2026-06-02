# TNotes Teams Web

Multi-user, collaborative Kanban board with realtime co-presence and offline-first sync.

## Quick start (Docker)

```bash
cd /home/tendo_mulira/projects/tnotes-teams
make docker-up
```

- Frontend: http://localhost:3000
- Backend: http://localhost:8080

Dev seed user:
- Email: `admin@tnotes.dev`
- Password: `password123`

## Local development (no Docker)

Prereqs:
- Go (per `go.mod`)
- Node.js (Vite warns on Node < 20.19 but still builds)
- PostgreSQL

```bash
export DATABASE_URL='postgres://tnotes:tnotes@localhost:5432/tnotes_teams?sslmode=disable'
go run ./cmd/server
```

In another terminal:

```bash
cd /home/tendo_mulira/projects/tnotes-teams/frontend
npm install
npm run dev
```

## Docs

- PRD: `/home/tendo_mulira/projects/tnotes-teams/TNotes Teams Web – Go-Based Multi-User Kanban Board PRD.md`
- Architecture: `/home/tendo_mulira/projects/tnotes-teams/docs/ARCHITECTURE.md`
- Feature checklist: `/home/tendo_mulira/projects/tnotes-teams/docs/FEATURES_CHECKLIST.md`

