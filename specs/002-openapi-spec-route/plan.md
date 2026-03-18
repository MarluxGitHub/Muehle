# Implementation Plan: OpenAPI-Dokument, Route & Swagger UI

**Branch**: `002-openapi-spec-route` | **Date**: 2026-03-18 | **Spec**: [spec.md](./spec.md)  
**Input**: Feature specification (inkl. Clarify: Swagger UI)

## Summary

1. **OpenAPI 3.0.3** als **`openapi.yaml`**, per **`go:embed`** in
   `pkg/muehle/interfaces`, Auslieferung **`GET /openapi.yaml`**
   (`application/yaml`).
2. **Swagger UI**: eingebettete **swagger-ui-dist**-Dateien unter
   `pkg/muehle/interfaces/swaggerui/`, Auslieferung per **`gin.StaticFS`**
   unter **`/swagger`**. **`index.html`** konfiguriert **`url: "/openapi.yaml"`**.
3. Nutzer öffnet **`http://localhost:8080/swagger/index.html`** für interaktive
   Doku (Try it out, gleiche Origin wie API).

Siehe [research.md](./research.md) (R6), [contracts/swagger-serving.md](./contracts/swagger-serving.md).

## Technical Context

**Language/Version**: Go 1.25+ (`marluxgithub/muehle`)  
**Primary Dependencies**: `github.com/gin-gonic/gin`; **keine** zusätzliche
OpenAPI-Go-Library; Swagger UI als **reine Static Files**  
**Storage**: N/A  
**Testing**: `go test`; manuell Browser + curl  
**Target Platform**: HTTP `:8080`  
**Project Type**: HTTP-API + eingebettete Spec + eingebettete Swagger UI  
**Performance Goals**: Statische Auslieferung < 2 s  
**Constraints**: Spec + UI öffentlich; keine Secrets in YAML  
**Scale/Scope**: OpenAPI mit allen Spiel-Operationen + Meta `/openapi.yaml` +
Swagger-UI-Bundle (~1–2 MB im Repo)

## Constitution Check

| Gate | Status | Notiz |
|------|--------|-------|
| Domain-First | Pass | Keine Spiellogik-Änderung. |
| SpecKit | Pass | Spec → Plan → Tasks. |
| Testbare Kernlogik | Pass | GameService unberührt. |
| Schichtentrennung | Pass | Nur HTTP/Doku-Schicht. |
| Einfachheit | Pass | Embed + StaticFS; kein Swagger-Codegen-Zwang. |

**Post-Phase-1**: Statische UI erhöht Repo-Größe—akzeptierter Tradeoff (R7).

## Project Structure

### Documentation

```text
specs/002-openapi-spec-route/
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/openapi.yaml
├── contracts/swagger-serving.md
└── tasks.md   # nach Plan-Update neu generieren oder ergänzen
```

### Source Code (Ziel)

```text
pkg/muehle/interfaces/
├── gin.go                 # GET /openapi.yaml, StaticFS /swagger
├── openapi.yaml           # go:embed
└── swaggerui/             # go:embed (dist + index.html)
    ├── index.html         # SwaggerUIBundle url: /openapi.yaml
    ├── swagger-ui-bundle.js
    ├── swagger-ui.css
    └── … (weitere swagger-ui-dist Dateien)
```

## Complexity Tracking

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| Große statische UI-Dateien | FR-006 Swagger UI im Browser | CDN abgelehnt (Offline/Repro); reine YAML ohne UI verletzt Spec |

## Phase 0 & 1 Outputs

| Artefakt | Pfad |
|----------|------|
| Research | [research.md](./research.md) |
| Data model | [data-model.md](./data-model.md) |
| OpenAPI | [contracts/openapi.yaml](./contracts/openapi.yaml) |
| Serving contract | [contracts/swagger-serving.md](./contracts/swagger-serving.md) |
| Quickstart | [quickstart.md](./quickstart.md) |

## Implementation Notes (für Tasks)

1. **`npm pack swagger-ui-dist`** oder Kopie aus `node_modules/swagger-ui-dist/dist`
   nach `swaggerui/`; **eigenes `index.html`** (Minimalbeispiel Swagger 5).
2. **`embed.FS`** für Unterbaum `swaggerui`; **keine** Pfade mit führendem `/`
   in embed-Pattern.
3. Optional **`GET /swagger`** → **302** auf `/swagger/index.html`.
4. **`docs/api.md`**: beide URLs (`/openapi.yaml`, `/swagger/index.html`).
5. Nach Router-Änderungen: **openapi.yaml** mitpflegen.
