# Data Model: OpenAPI + Swagger UI

Keine **Domain**-Entitäten.

## Artefakte

| Artefakt | Rolle |
|----------|--------|
| **openapi.yaml** | OpenAPI-3.x, Quelle für Embed + Swagger UI `url`. |
| **swaggerui/** | Eingebettete statische Swagger-UI-Dist + `index.html`. |

## HTTP-Routen (Doku-Schicht)

| Route | Methode | Zweck |
|-------|---------|--------|
| `/openapi.yaml` | GET | Rohspec, `application/yaml` |
| `/swagger/*` | GET | Swagger UI (HTML/JS/CSS) |

## Validierung

- OpenAPI 3.0.x parsierbar; UI lädt Spec ohne Cross-Origin-Fehler (gleicher Host/Port).
