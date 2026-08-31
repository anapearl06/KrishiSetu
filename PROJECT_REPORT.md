# KrishiSetu — Full Project Report

Generated: 2026-08-29

---

## 1. Overview

KrishiSetu is an AI-powered agricultural marketplace (SIH 2026, Problem Statement 26033) connecting farmers/FPOs directly to consumers and bulk buyers, eliminating middlemen.

The project has **two parts**:

| Part | Location | Stack | Status |
|------|----------|-------|--------|
| Backend API | `D:\Coding\Projects\SIH'26\KrishiSetu` | Go 1.25, Gin, GORM, PostgreSQL | Deployed on Render (https://krishisetu-api-tiau.onrender.com) |
| Frontend design system / prototypes | `C:\Users\HP\OneDrive\Desktop\HTML\stitch_markdown_design_system\stitch_markdown_design_system` | HTML, Tailwind CSS, Material Symbols, Inter | Static prototypes only (not wired to API) |

---

## 2. Repository File Tree

```
KrishiSetu\
├── .env                      # Local dev environment (untracked in git, present locally)
├── .env.example              # Template for local env config
├── .gitignore
├── docker-compose.yml        # Local PostgreSQL 17 container
├── go.mod                    # Module definition, Go 1.25.5, Gin/GORM/deps
├── go.sum                    # Dependency checksums
├── README.md                 # Full product / architecture documentation
├── PROJECT_REPORT.md         # This file
├── cmd\
│   └── server\
│       └── main.go           # Entry point: env load, DB connect, router start
├── internal\
│   ├── auth\                 # Auth primitives (JWT + bcrypt)
│   │   ├── jwt.go
│   │   └── password.go
│   ├── buyer\                # Buyer domain
│   │   ├── handler.go
│   │   ├── model.go
│   │   ├── repository.go
│   │   └── service.go
│   ├── config\config.go      # Env-based configuration
│   ├── database\database.go  # GORM/Postgres connection
│   ├── farmer\               # Farmer domain
│   │   ├── handler.go
│   │   ├── model.go
│   │   ├── repository.go
│   │   └── service.go
│   ├── listing\model.go      # CropListing entity (schema only, no CRUD yet)
│   ├── middleware\auth.go    # JWT bearer-auth middleware
│   ├── migrations\           # SQL migrations 000001..000003 (up/down)
│   │   ├── 000001_create_farmers.up/down.sql
│   │   ├── 000002_create_buyers.up/down.sql
│   │   └── 000003_create_crop_listing.up/down.sql
│   └── server\
│       ├── response\response.go  # Standard {success,data}/{success,error} envelope
│       └── server.go             # Gin router + route wiring
```

---

## 3. Backend Files (detailed)

### 3.1 `cmd/server/main.go`
Loads `.env` (non-fatal if missing), loads config, connects DB, builds router, starts Gin on `:PORT` (default 8080). Fails with `log.Fatal` on DB or server errors.

### 3.2 `internal/config/config.go`
Reads env vars: `APP_ENV`, `PORT`/`APP_PORT`, `DATABASE_URL`, `JWT_SECRET`, `JWT_EXPIRATION_HOURS` (default 24h). Note: local `.env` uses discrete `DB_HOST/PORT/USER/...` vars but the code only reads a single `DATABASE_URL` — so local `.env` is not directly consumable unless `DATABASE_URL` is provided.

### 3.3 `internal/database/database.go`
`Connect(cfg)` → `gorm.Open(postgres.Open(cfg.DatabaseURL))`. Errors are wrapped; requires `DATABASE_URL` set.

### 3.4 `internal/auth/password.go`
bcrypt helpers: `HashPassword` (default cost), `ComparePassword`, and boolean `CheckPassword` used by services.

### 3.5 `internal/auth/jwt.go`
`Claims{UserID, Role, RegisteredClaims}`. `GenerateToken(userID, role, secret, expiry)` signs HS256 JWT with `iat`/`exp`. `ValidateToken` parses/verifies and rejects non-HS256 signing methods.

### 3.6 `internal/middleware/auth.go`
`JWTAuth(secret)` middleware: requires `Authorization: Bearer <token>`, validates JWT, injects `user_id` and `role` into context (keys `UserIDKey`, `RoleKey`). Responds 401 for missing/invalid/expired tokens.

### 3.7 `internal/server/response/response.go`
Standardizes API output:
- Success: `{ "success": true, "data": ... }`
- Error: `{ "success": false, "error": { "code", "message" } }`

### 3.8 `internal/server/server.go`
Wires the Gin router:
- `GET /health` → `{status: ok}`
- `/api/v1/farmers` → `POST /register`, `POST /login`, `GET /me` (auth), `PUT /me` (auth)
- `/api/v1/buyers` → `POST /register`, `POST /login`, `GET /me` (auth), `PUT /me` (auth)

The `listing` package is imported nowhere yet — no listing routes exist.

### 3.9 Farmer domain
- **`model.go`**: `Farmer{ID, Name, Phone (unique), PasswordHash (json:"-"), State, District, CreatedAt, UpdatedAt}`
- **`repository.go`**: interface `Create / FindByPhone / FindByID / Update`; Postgres impl via GORM (`Save` for update).
- **`service.go`**: `Register` (validates name/phone/password/state/district, hashes password), `Login` (validates password → JWT, role `"farmer"`), `GetProfile`, `UpdateProfile` (name/state/district only).
- **`handler.go`**: HTTP adapters. `Me` only echoes `user_id`/`role` and does NOT fetch the farmer row (inconsistent with buyer's `Me`). Register returns 500 on validation errors (`INTERNAL_ERROR`).

### 3.10 Buyer domain
- **`model.go`**: `Buyer{ID, Name, Phone (unique), PasswordHash (json:"-"), BusinessName, BusinessType, State, District, CreatedAt, UpdatedAt}`
- **`repository.go`**: same interface shape but ctx-aware (`WithContext(ctx)`).
- **`service.go`**: ctx-aware, trims/validates inputs, blocks duplicate phone on register, JWT role `"buyer"`.
- **`handler.go`**: `Me` verifies role == "buyer" and returns the full buyer profile (more complete than farmer). Register/Update return finer-grained 400s.

### 3.11 `internal/listing/model.go`
`CropListing{ID, FarmerID, CropName, Quantity, Unit, ExpectedPrice, QualityGrade, State, District, HarvestDate, Status, CreatedAt, UpdatedAt}`. **Schema/entity only** — no repository, service, handler, or routes yet.

### 3.12 Migrations
| Migration | Up | Down |
|-----------|----|----|
| `000001_create_farmers` | `farmers` (id BIGSERIAL PK, name, phone UNIQUE, password_hash, state, district, timestamps) | DROP farmers |
| `000002_create_buyers` | `buyers` (id, name, phone UNIQUE, password_hash, business_name, business_type, state, district, timestamps) | DROP buyers |
| `000003_create_crop_listing` | `crop_listings` + FK → farmers(id) ON DELETE CASCADE + 5 indexes (farmer_id, crop_name, state, district, status) | DROP crop_listings |

Note: no table migrations or auto-migration calls in code; buyer table timestamps use `TIMESTAMP` while others use `TIMESTAMPTZ` (minor inconsistency).

### 3.13 `docker-compose.yml`
Single service: `postgres:17-alpine`, user/pass/db `krishisetu`, port 5432, named volume for persistence.

### 3.14 `.env` / `.env.example`
Dev values: `APP_PORT=8080`, DB connection fields, `JWT_SECRET=change-this...`. **Note**: `.env` is present locally and committed-adjacent; the code expects `DATABASE_URL`/`JWT_SECRET`, so the file layout does not match what the app reads.

### 3.15 `go.mod`
Dependencies: Gin v1.12.0, GORM + Postgres driver, pgx v5, golang-jwt/v5, bcrypt (x/crypto), godotenv, sonic JSON, quic-go (indirect). Go 1.25.5. All non-std libs except godotenv are marked `// indirect` (would be cleaned up by `go mod tidy`).

### 3.16 `README.md`
Extensive product doc: problem statement, solution, feature roadmap (agri-intel, matching algorithm, ML/price intelligence, logistics & payment scope, MVP scope), architecture diagram (Flutter → Go/Gin → Postgres/Redis/FastAPI), and entity list. Several features described here do **not** exist in code yet (Redis, FastAPI/ML, offers, orders, payments, matching).

---

## 4. Frontend Design System (all files)

### 4.1 Design-token docs
| File | Purpose |
|------|---------|
| `design.md` | High-level tokens: Primary #2E7D32, Secondary #8D6E63, Bg #FAFAF5, Inter, 4–8px corners, 44px touch targets, ~14-screen MVP rule, component reuse across roles, drawers over pages, states in components. |
| `agricore_design_system/DESIGN.md` | Expanded Material-3-style token set (YAML frontmatter): full color ramp (surface-*, primary-*, tertiary-*, error-*), typography scale (headline-lg 44px → caption 12px), radii, spacing, touch targets. |

### 4.2 Screens (13 screens, each = `code.html` + `screen.png`)
All use a shared pattern: Tailwind CDN + inline `tailwind.config` palette generated from the DESIGN.md tokens, Inter font, Material Symbols Outlined icons, responsive layout (desktop sidebar → mobile bottom-nav).

| Folder | Title | Role | Components |
|--------|-------|------|-----------|
| `krishisetu_home` | Empowering the Ground Level | Public | Top nav, hero, bento feature grid, footer |
| `registration_krishisetu` | Registration | Public | Signup form (role/tabs) |
| `login_krishisetu` | Login | Public | Login form, phone + password |
| `farmer_dashboard` | Farmer Dashboard | Farmer | Sidebar (Dashboard/My Produce/Offers/Orders/Profile), stat cards, recent listings |
| `my_produce_krishisetu` | My Produce | Farmer | Listing cards + create-listing action |
| `my_offers_krishisetu` | My Offers | Farmer | Incoming offer list |
| `my_orders_krishisetu` | My Orders | Farmer | Order list with status chips |
| `buyer_dashboard_krishisetu` | Buyer Dashboard | Buyer | Sidebar (Dashboard/Marketplace/Prices/Offers/Orders/Profile), stats |
| `browse_produce_marketplace` | Browse Produce | Buyer | Marketplace grid, filters, search (543 lines, largest screen) |
| `market_intelligence_prices_trends` | Market Prices | Buyer | Price tables, trends, mandi comparison |
| `my_offers_sent_krishisetu` | My Offers | Buyer | Outgoing offer list (reused component, inverse direction) |
| `my_orders_buyer_krishisetu` | My Orders | Buyer | Buyer-side order list |
| `profile_settings_krishisetu` | Profile & Settings | Shared | Profile edit form, settings |

### 4.3 Assets
| Folder | Purpose |
|--------|---------|
| `krishisetu_logo/screen.png` | Logo asset (referenced as image in nav/dashboard sidebars) |
| `high_quality_..._crop_field/screen.png` | Hero/landing crop-field photography asset |

### 4.4 Prototype characteristics
- All data is mock/static, hard-coded inline.
- No JS logic beyond Tailwind config; forms do not submit.
- Assets are hot-linked from `lh3.googleusercontent.com` (external URL dependency).
- Dark-mode classes (`dark:bg-...`) are present on most screens but only activated by the `dark` class on `<html>`.
- Note: buyer/farmer dashboards and lists rely on the same list components per screen — consistent with the "reuse across roles" MVP rule.

---

## 5. API Surface (backend) ↔ Screens (frontend) mapping

| Endpoint (backend) | Screen(s) it powers |
|--------------------|---------------------|
| `POST /api/v1/farmers/register` | `registration_krishisetu` |
| `POST /api/v1/buyers/register` | `registration_krishisetu` |
| `POST /api/v1/farmers/login` / `buyers/login` | `login_krishisetu` |
| `GET/PUT /api/v1/farmers/me` | `farmer_dashboard`, `profile_settings_krishisetu` |
| `GET/PUT /api/v1/buyers/me` | `buyer_dashboard_krishisetu`, `profile_settings_krishisetu` |
| *(not implemented)* listings CRUD | `my_produce_krishisetu`, `browse_produce_marketplace` |
| *(not implemented)* offers, orders, market prices | offer/order screens, `market_intelligence_prices_trends` |

The `CropListing` entity and migration exist but there is **no listing API**, and offers/orders/prices have neither backend endpoints nor tables — the frontend prototypes currently outpace the backend.

---

## 6. Observations & Gaps

1. **Listing module incomplete**: `internal/listing/model.go` + migration 000003 exist, but no repository/service/handler/routes.
2. **Farmer vs buyer asymmetry**: farmer `Me` doesn't fetch the profile; buyer `Me` does. Farmer register maps service errors to 500 `INTERNAL_ERROR` instead of 400.
3. **Error mapping**: service returning `err != nil` from `FindByPhone` on login is exposed as generic 401 even for not-found.
4. **Config mismatch**: code reads `DATABASE_URL` only, but `.env` defines discrete `DB_*` vars. `JWT_EXPIRY_HOURS` in `.env.example` vs `JWT_EXPIRATION_HOURS` read by config.
5. **`go.mod` cleanliness**: third-party modules listed as `// indirect` — needs `go mod tidy`.
6. **No tests, no CI, no lint config** despite README promising Go testing, Testify, Swagger.
7. **No auto-migration**: migrations run externally; GORM models won't auto-create tables.
8. **Recent commit added listing + buyer code but no wiring to router** for listing.
9. **Frontend is static HTML**; not convertible directly — will need Flutter implementation per README.
10. **`.env` committed** to the working tree (contains dev-only creds, harmless but should stay out of git).

---

## 7. Deployment

- **Backend**: Render at `https://krishisetu-api-tiau.onrender.com` — verified live (`GET /health` → `{"success":true,"data":{"status":"ok"}}`).
- **Database**: PostgreSQL via `DATABASE_URL` env var configured in Render (provider not visible in repo; local dev uses docker-compose Postgres 17).