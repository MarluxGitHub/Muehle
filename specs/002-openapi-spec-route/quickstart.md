# Quickstart: OpenAPI & Swagger UI

Server starten:

```bash
go run ./cmd/server
```

## Rohspec

```bash
curl -s http://localhost:8080/openapi.yaml | head -20
```

## Swagger UI im Browser

Öffnen:

```text
http://localhost:8080/swagger/
```

oder: `http://localhost:8080/swagger` → Redirect.

Erwartung: Alle Operationen sichtbar; **Try it out** gegen `http://localhost:8080` möglich (gleiche Origin).

## Doku

Siehe `docs/api.md`—dort **OpenAPI-URL** und **Swagger-URL** aufführen.
