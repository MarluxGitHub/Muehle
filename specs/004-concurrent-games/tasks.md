---
description: "Tasks: Mehrere Spiele gleichzeitig (004-concurrent-games)"
---

# Tasks: Mehrere Spiele gleichzeitig

**Input**: `specs/004-concurrent-games/` (plan, spec, contracts, quickstart)  
**Pfad-Fokus**: `pkg/muehle/interfaces/gin_integration_test.go`, optional `docs/api.md`, `openapi.yaml`

## Format

`- [ ] Tnnn [P?] [USn?] Beschreibung mit Dateipfad`

---

## Phase 1: Setup

**Purpose**: Ausgangslage vor Änderungen

- [x] T001 Aus Repo-Root **`go test ./...`** ausführen—Baseline grün

---

## Phase 2: Foundational (Ist-Abgleich)

**Purpose**: Sicherstellen, dass kein Spielpfad ohne `gameId` existiert

- [x] T002 In **`pkg/muehle/interfaces/gin.go`** verifizieren: alle Handler unter **`/games/:gameId`** rufen **`resolveGame`** (bzw. **`Registry.Get`**) auf; **`POST /games`** ist die einzige Spiel-anlegende Route ohne `gameId`. Abweichungen beheben oder in **`specs/004-concurrent-games/research.md`** als Gap dokumentieren

**Checkpoint**: Kein direkter Zugriff auf eine einzelne globale Partie außerhalb der Registry.

---

## Phase 3: User Story 1 – Zwei Partien nebeneinander (Priority: P1) 🎯 MVP

**Goal**: Zwei unterschiedliche Partie-IDs; beide abfragbar (FR-001, FR-005)

**Independent Test**: Zwei `POST /games`, zwei `GET …/state` → je 200, unterschiedliche `gameId`

- [x] T003 [US1] In **`pkg/muehle/interfaces/gin_integration_test.go`**: Test **`TestConcurrentGames_TwoGamesCreatedBothAddressable`**—zweimal **`POST /games`**, IDs ungleich; für beide **`GET /games/{id}/state`** → **200** und plausibler Zustand (z. B. Warten auf Spieler)

---

## Phase 4: User Story 2 – Isolation bei Zügen (Priority: P2)

**Goal**: Züge in Partie A ändern Brett von B nicht (FR-003, FR-004, SC-002)

**Independent Test**: Mindestens fünf Partien; nach Zügen nur in Partie 1 bleiben Bretter 2–5 unverändert (leer bis auf deren eigene Züge—hier: nur Spiel 1 bekommt Züge)

- [x] T004 [US2] In **`pkg/muehle/interfaces/gin_integration_test.go`**: Test **`TestConcurrentGames_FiveGamesMoveInFirstOnlyOthersUnchanged`**—**5×** **`POST /games`**, je Partie 2 Spieler hinzufügen; in **nur** Partie 1 mindestens einen legalen **`place`**-Zug; **`GET /games/{id}/board`** für Partien 2–5: gleiche Belegung wie **vor** dem Zug in Partie 1 (z. B. alle Felder leer / `Color` 0)

---

## Phase 5: User Story 3 – Unbekannte Partie (Priority: P3)

**Goal**: Falsche `gameId` → Fehler ohne fremde Mutation (FR-006, SC-003)

**Independent Test**: `POST` Move mit gültigem UUID-Format, nicht in Registry → 404; bestehendes Spiel unverändert

- [x] T005 [US3] In **`pkg/muehle/interfaces/gin_integration_test.go`**: Test **`TestConcurrentGames_UnknownGameIdMoveReturns404AndDoesNotTouchExisting`**—eine echte Partie anlegen, **`GET /board`** Snapshot; dann **`POST /games/{random-uuid}/moves`** mit existierendem Pfad-Format aber ID die **nie** per **`POST /games`** erzeugt wurde → **404**; erneut **`GET`** Board der ersten Partie → identisch zum Snapshot

---

## Phase 6: Polish

- [x] T006 [P] **`docs/api.md`**: Absatz—mehrere Partien parallel über mehrfaches **`POST /games`**, jeweils eigene `gameId` in URLs
- [x] T007 [P] **`pkg/muehle/interfaces/openapi.yaml`**: Kurzbeschreibung bei **`POST /games`** oder Server-Description—mehrere gleichzeitige Partien pro Server-Instanz
- [x] T008 **`gofmt`** auf geänderte Go-Dateien
- [x] T009 **`go test ./...`** grün; optional **`go test -race ./pkg/muehle/interfaces/...`**
- [x] T010 Ablauf **`specs/004-concurrent-games/quickstart.md`** mit laufendem Server verifizieren (manuell oder skriptbar dokumentieren)

---

## Dependencies

```
T001 → T002 → T003 → T004 → T005 → T006,T007 (parallel) → T008 → T009 → T010
```

**T004/T005** können nach T003 parallel bearbeitet werden, wenn Hilfsfunktionen (`postForm`, `createGame`) im Testfile geteilt werden—sonst sequenziell.

---

## Parallel Opportunities

| Parallel | Tasks |
|----------|--------|
| Nach T005 | T006 + T007 (verschiedene Dateien) |
| Optional | T004 und T005 parallel (beide nur `gin_integration_test.go`)—Merge-Konflikt möglich → besser nacheinander oder eine Person |

---

## MVP

**T001–T003** liefern den minimalen Nachweis für US1 (zwei Partien). **T004–T005** decken SC-002/SC-003 ab.

---

## Notes

- Wenn T002 eine Lücke findet: zuerst **`gin.go`** reparieren, dann Tests.
- Keine Änderung an Mühle-Domain nötig, sofern Isolation über Registry gewährleistet ist.
