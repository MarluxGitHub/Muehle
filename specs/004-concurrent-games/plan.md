# Implementation Plan: Mehrere Spiele gleichzeitig

**Branch**: `004-concurrent-games` | **Date**: 2026-03-18 | **Spec**: [spec.md](./spec.md)  
**Input**: Feature specification from `/specs/004-concurrent-games/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command.

## Summary

Mehrere Mühle-Partien sollen **parallel und isoliert** laufen (FR-001–006). Im Ist-Stand existiert bereits eine **`GameRegistry`** (UUID → `Application`) und die HTTP-Schicht routet alle Spielaktionen über **`/games/:gameId`**. Die **Spiellogik pro Partie** ist unverändert; Schwerpunkt der Umsetzung/Verifikation liegt auf **Nachweis der Isolation** (Tests) und **Dokumentation des Vertrags**. Keine neue Persistenzschicht nötig (Spec schließt Neustart-Persistenz aus).

## Technical Context

**Language/Version**: Go 1.22+ (Modul `marluxgithub/muehle`)  
**Primary Dependencies**: Gin, google/uuid  
**Storage**: In-Memory `map[string]*Application` (pro Server-Prozess)  
**Testing**: `go test ./...`; Integrationstests mit `httptest` + Gin TestMode  
**Target Platform**: Linux/macOS Server (`:8080`)  
**Project Type**: Web-Service (REST)  
**Performance Goals**: Kein hartes Ziel in Spec; viele parallele Partien dürfen linear mit Speicher wachsen  
**Constraints**: Keine Vermischung von Spielzuständen; unbekannte `gameId` → 404, keine fremden Mutationen  
**Scale/Scope**: „Viele“ gleichzeitige Partien auf einer Instanz; keine Mandantenfähigkeit über Servergrenzen hinaus

## Constitution Check

*GATE: passed. Re-checked after Phase 1 design.*

| Prinzip | Erfüllung |
|---------|-----------|
| **Domain-First** | Mühle-Regeln unverändert pro `Application`; dieses Feature betrifft nur Session-Isolation. |
| **SpecKit** | Spec → Plan → Tasks; Anforderungen auf Registry + API abbildbar. |
| **Testbare Kernlogik** | Isolation per HTTP-Integrationstests (und optional `GameRegistry`-Concurrency); Domain bleibt getrennt testbar. |
| **Schichtentrennung** | `interfaces` → `application`/`GameRegistry` → Domain-Services unverändert. |
| **Einfachheit** | Kein zweites Speichersystem; Verifikation statt großer Refactorings. |

## Project Structure

### Documentation (this feature)

```text
specs/004-concurrent-games/
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
│   └── concurrent-games-behavior.md
└── tasks.md              # /speckit.tasks
```

### Source Code (repository root)

```text
cmd/server/main.go
pkg/muehle/application/
  game_registry.go      # CreateGame, Get (UUID)
  application.go        # pro Partie: GameService-Fassade
pkg/muehle/interfaces/
  gin.go                # POST /games, /games/:gameId/...
  gin_integration_test.go
pkg/muehle/domain/services/
  …                     # unverändert für dieses Feature
```

**Structure Decision**: Monorepo-Go-Service; alle spielbezogenen Routen hängen an `gameId`. Feature 004 erweitert primär **Tests + Docs**, nicht zwingend Produktionscode.

## Complexity Tracking

Keine Verfassungsverletzungen; Tabelle nicht nötig.
