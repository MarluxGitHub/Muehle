# Contract: Auslieferung OpenAPI & Swagger UI

| Ressource | URL (Dev) | Content-Type / Typ |
|-----------|-----------|-------------------|
| OpenAPI YAML | `GET /openapi.yaml` | `application/yaml` |
| Swagger UI | `GET /swagger/` (+ Assets z. B. `/swagger/swagger-ui.css`) | `text/html`, JS, CSS |

**Kopplung**: Swagger UI lädt die Spec von **`/openapi.yaml`** (relative zur gleichen Origin).

**Nicht Teil der Spiel-API**: reine Dokumentation; keine Auth.
