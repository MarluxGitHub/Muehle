# Research: OpenAPI, Auslieferungsroute & Swagger UI

## R1: Format (JSON vs. YAML)

**Decision**: **YAML** als **einzige** Quelldatei (`openapi.yaml`), Auslieferung per
`GET /openapi.yaml` mit `Content-Type: application/yaml`.

**Rationale**: Lesbar in Reviews; Postman/OpenAPI-Importer akzeptieren YAML.

**Alternatives**: Nur JSON—weniger lesbar im Repo.

## R2: Einbindung Spec in Go/Gin

**Decision**: **`go:embed`** für `openapi.yaml`; dedizierter Handler.

**Rationale**: Kein Dateisystem-Pfad zur Laufzeit; ein Binary.

**Alternatives**: `swaggo/swag` Codegen—nicht nötig für statische Spec.

## R3: Kanonische Spec-Datei

**Decision**: **`pkg/muehle/interfaces/openapi.yaml`** (Embed); Contract-Spiegel
unter `specs/002-openapi-spec-route/contracts/openapi.yaml` synchron halten.

## R4: servers-URL

**Decision**: `http://localhost:8080` in OpenAPI `servers`.

## R5: OpenAPI-Route in paths

**Decision**: `GET /openapi.yaml` in OpenAPI `paths` dokumentieren (Meta).

## R6: Swagger UI (FR-006)

**Decision**: **Offizielle Swagger UI Dist** (statische Dateien: `swagger-ui.css`,
`swagger-ui-bundle.js`, `swagger-ui-standalone-preset.js`) plus **minimales
`index.html`**, das `SwaggerUIBundle({ url: "/openapi.yaml", … })` setzt. Alles
per **`go:embed`** unter `pkg/muehle/interfaces/swaggerui/` (Unterordner).

**Gin**: `router.StaticFS("/swagger", http.FS(swaggerUIFS))`; Nutzer öffnet
**`http://localhost:8080/swagger/index.html`** (oder Redirect von `GET /swagger`
→ `/swagger/index.html`).

**Rationale**: Keine zusätzliche Go-Dependency (`gin-swagger`/`swag` oft an
generierte Spec gebunden); volle Kontrolle; **gleiche Origin** wie API →
**Try it out** ohne CORS-Zusatzprobleme (Spec und XHR gehen gegen `:8080`).

**Alternatives**:

- **CDN** für Swagger UI—offline/build-repro schlechter, CSP.
- **swaggo/files** + handlers—mehr Magic, Versionskopplung.

## R7: Asset-Größe

**Decision**: Swagger-UI-Dist (~1–2 MB) **versioniert im Repo** unter
`swaggerui/` oder per **dokumentiertem einmaligen Copy** aus npm
`swagger-ui-dist` (Tasks). Akzeptierter Tradeoff für Einfachheit.
