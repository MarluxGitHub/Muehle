---
description: "Tasks: Korrektes Schlagen bei Mühlen (003-mill-capture-rule)"
---

# Tasks: Korrektes Schlagen bei Mühlen

**Input**: `specs/003-mill-capture-rule/` (plan, research, contracts)  
**Pfad-Fokus**: `pkg/muehle/domain/services/`

## Format

`- [ ] Tnnn [P?] [USn?] Beschreibung mit Dateipfad`

---

## Phase 1: Setup

- [x] T001 Aus Repo-Root **`go test ./pkg/muehle/domain/...`** ausführen—Baseline grün vor Änderungen

---

## Phase 2: Foundational (Mill-Geometrie)

**Purpose**: Eine Quelle für die 16 Mühlen-Linien, keine doppelte Logik

- [x] T002 In **`pkg/muehle/domain/services/boardService.go`**: Mühlen-Linien als **gemeinsame Konstante** oder Hilfsfunktion extrahieren (bestehende `millCombination` aus **`HasPlayerThreeStones`** wiederverwenden); **`HasPlayerThreeStones`** auf dieselbe Definition umstellen, Verhalten unverändert lassen

**Checkpoint**: `go test` weiterhin grün.

---

## Phase 3: User Story 1 – Kein Schlag aus Mühle (Priority: P1)

**Goal**: Remove aus Mühle verboten, wenn Gegner „freien“ Stein hat

**Independent Test**: Unit-Tests + manuell laut quickstart

- [x] T003 [P] [US1] In **`boardService.go`**: **`IsFieldPartOfClosedMill(board entities.Board, fieldIndex int, color entities.Color) bool`** — Feld `fieldIndex` trägt `color` **und** liegt auf mindestens einer Linie, deren **alle drei** Felder `color` sind
- [x] T004 [US1] In **`boardService.go`**: **`EnemyHasStoneOutsideMill(board, enemyColor) bool`** — **true**, wenn **mindestens ein** Feld mit `enemyColor` existiert mit **`!IsFieldPartOfClosedMill(...)`**
- [x] T005 [P] [US1] **`boardService_test.go`** (neu oder erweitern): Fälle für `IsFieldPartOfClosedMill` und `EnemyHasStoneOutsideMill` (mind. eine volle Mühle + isolierter Stein)
- [x] T006 [US1] In **`gameService.go`** **`RemoveStone`**: nach Validierung „Feld = Gegner“, **vor** State-Mutation: wenn **`EnemyHasStoneOutsideMill`** **und** **`IsFieldPartOfClosedMill(..., enemyColor, fieldIndex)`** → **`return errors.New("cannot remove stone from mill")`** (oder gleichwertiger klarer Text)
- [x] T007 [P] [US1] **`gameService_test.go`** (neu): Partie in **`GameStateRemovingStone`** konstruieren—Gegner mit Mühle + Stein außerhalb; **`RemoveStone`** auf Mühlenfeld → **Fehler**; auf freies Feld → **nil**, Brett konsistent

---

## Phase 4: User Story 2 – Ausnahme alle in Mühlen (Priority: P2)

- [x] T008 [US2] In **`gameService_test.go`**: Stellung, in der **jeder** Gegnerstein in mindestens einer geschlossenen Mühle liegt (z. B. 6 Steine, zwei überlappende Mühlen o. ä.); **`RemoveStone`** auf beliebiges Gegnerfeld → **erfolgreich**

---

## Phase 5: Polish

- [x] T009 **`gofmt`** auf geänderte Go-Dateien
- [x] T010 **`go test ./...`** grün
- [x] T011 Optional: **`contracts/remove-stone-behavior.md`** mit tatsächlichem Fehlerstring abgleichen

---

## Dependencies

```
T001 → T002 → T003,T004 → T005 → T006 → T007 → T008 → T009–T011
```

**T003/T005** parallel nach T002 möglich.

---

## MVP

**T002–T007** decken US1; **T008** verifiziert US2.

---

## Notes

- Keine Änderung an **`pkg/muehle/interfaces/gin.go`** nötig, sofern Fehler wie bisher als 500/JSON durchgereicht werden.
