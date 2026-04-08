# Launchdate Monorepo

This repository combines the original web, admin, and backend projects into a single workspace without preserving the original Git history.

## Structure

```text
launchdate/
  apps/
    web/
    admin/
  services/
    backend/
```

## Projects

- `apps/web`: user-facing frontend, built with Vite + React.
- `apps/admin`: admin frontend, built with Vite + React.
- `services/backend`: Go backend service.

## Prerequisites

- Node.js 20+
- npm
- Go 1.22+

## Install Dependencies

Install frontend dependencies separately:

```bash
cd apps/web && npm install
cd ../admin && npm install
```

Install backend dependencies:

```bash
cd services/backend && go mod download
```

## Run In Development

Start the web app:

```bash
cd apps/web
npm run dev
```

Start the admin app:

```bash
cd apps/admin
npm run dev
```

Start the backend:

```bash
cd services/backend
make run
```

## Build

Build the web app:

```bash
cd apps/web
npm run build
```

Build the admin app:

```bash
cd apps/admin
npm run build
```

Build the backend:

```bash
cd services/backend
make build
```

## Common Commands

Frontend lint:

```bash
cd apps/web && npm run lint
cd apps/admin && npm run lint
```

Backend checks:

```bash
cd services/backend && make test
cd services/backend && make lint
```

## Notes

- The original `.git` directories were not copied into this repository.
- Generated directories such as `node_modules` and `dist` were intentionally excluded during merge.
- If you need unified startup scripts later, they can be added at the repository root.