---
description: "Tasks: OpenAPI + Swagger UI (002-openapi-spec-route)"
---

# Tasks: OpenAPI-Dokument, Route & Swagger UI

**Input**: `specs/002-openapi-spec-route/` (plan, contracts, quickstart)  
**Prerequisites**: plan.md, spec.md, `contracts/openapi.yaml`, `contracts/swagger-serving.md`

## Format

`- [ ] Tnnn [P?] [USn?] Beschreibung mit Dateipfad`

---

## Phase 1: Setup

**Purpose**: Contract und Tooling vorbereiten

- [x] T001 Inhalt von **`specs/002-openapi-spec-route/contracts/openapi.yaml`** nach **`pkg/muehle/interfaces/openapi.yaml`** kopieren/synchronisieren (kanonische Embed-Quelle im Paket)

---

## Phase 2: Foundational – Swagger UI Assets

**Purpose**: Statische Dateien vor Gin-Registrierung bereitstellen

- [x] T002 **`swagger-ui-dist`** bereitstellen: z. B. `npm install swagger-ui-dist@5` in temp-Verzeichnis oder `npx` und **`dist/`**-Inhalt nach **`pkg/muehle/interfaces/swaggerui/`** kopieren (alle benötigten JS/CSS; siehe upstream `swagger-ui-dist` README)
- [x] T003 **`swagger-initializer.js`** mit **`url: '/openapi.yaml'`**; **`index.html`** aus Dist (petstore-URL ersetzt)

**Checkpoint**: Ordner `swaggerui/` enthält Dist + funktionsfähiges `index.html`.

---

## Phase 3: User Story 1 – OpenAPI-Route (Priority: P1)

**Goal**: `GET /openapi.yaml` → gültiges YAML

**Independent Test**: `curl -s http://localhost:8080/openapi.yaml | head -5`

- [x] T004 [US1] **`pkg/muehle/interfaces/openapi_embed.go`**: **`//go:embed openapi.yaml`**
- [x] T005 [US1] **`gin.go`**: **`GET /openapi.yaml`**, **`application/yaml`**

---

## Phase 4: User Story 2 – Swagger UI (Priority: P1)

**Goal**: Browser **`/swagger/`** zeigt interaktive Doku

**Independent Test**: `http://localhost:8080/swagger/` (Redirect von `/swagger`)

- [x] T006 [US2] **`//go:embed swaggerui/*`** in **`openapi_embed.go`**
- [x] T007 [US2] **`StaticFS("/swagger", …)`** + **`fs.Sub(..., "swaggerui")`**
- [x] T008 [US2] **`GET /swagger`** → **302** **`/swagger/`**

---

## Phase 5: User Story 3 – Auffindbarkeit (Priority: P2)

- [x] T009 [US3] **`docs/api.md`** – OpenAPI + Swagger-URLs
- [x] T010 [P] [US3] README.md – _nicht vorhanden, übersprungen_

---

## Phase 6: Polish

- [x] T011 `gofmt`
- [x] T012 **`go test ./...`**
- [x] T013 **quickstart.md** – Pfade an **`/swagger/`** angepasst
- [x] T014 **`TestOpenAPIAndSwaggerRoutes`** in **`gin_integration_test.go`**

---

## Dependencies

(siehe ursprünglicher Plan—umgesetzt)

---

## Notes

- Swagger UI **~1.7 MB** unter **`pkg/muehle/interfaces/swaggerui/`** (git-tracked).
