# HTTP API (Mühle)

- **Dev-Base-URL:** `http://localhost:8080`

## OpenAPI

- **Rohspezifikation (YAML):** `GET http://localhost:8080/openapi.yaml`  
  Content-Type: `application/yaml` — für Codegeneratoren, Postman (Import via URL).

## Swagger UI

- **Interaktive Doku (Browser):** `http://localhost:8080/swagger/`  
  Kurz: `http://localhost:8080/swagger` → Redirect auf `/swagger/`.  
  Lädt dieselbe OpenAPI wie oben; **Try it out** gegen denselben Host möglich.

## Weitere Doku

- Legacy-Mapping: [`specs/001-rest-routing-cleanup/contracts/http-api.md`](../specs/001-rest-routing-cleanup/contracts/http-api.md)  
- Ablauf Spiel anlegen: [`specs/001-rest-routing-cleanup/quickstart.md`](../specs/001-rest-routing-cleanup/quickstart.md)

Kurz: `POST /games` → leeres Spiel → Spieler unter `/games/{id}/players` → Züge unter `/games/{id}/moves`.
