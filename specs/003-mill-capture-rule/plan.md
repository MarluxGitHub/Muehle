# Implementation Plan: Korrektes Schlagen bei Mühlen

**Branch**: `003-mill-capture-rule` | **Date**: 2026-03-18 | **Spec**: [spec.md](./spec.md)

## Summary

Die Domain erweitert **`BoardService`** um Hilfen zur Mühlen-Zugehörigkeit eines
Feldes und ob der Gegner noch „freie“ Steine hat. **`GameService.RemoveStone`**
lehnt das Entfernen eines Steins aus einer **geschlossenen gegnerischen Mühle**
ab, **sofern** der Gegner mindestens einen Stein **außerhalb** jeder Mühle hat.
Liegen **alle** Gegnersteine in Mühlen, bleibt das Entfernen **beliebig** möglich.
Siehe [research.md](./research.md).

## Technical Context

**Language/Version**: Go 1.25+ (`marluxgithub/muehle`)  
**Primary Dependencies**: keine neuen  
**Storage**: N/A  
**Testing**: **`go test`** in `pkg/muehle/domain/services` (Tabelle der 16 Linien
wie in `boardService.go` wiederverwenden oder extrahieren)  
**Target Platform**: Domain + bestehende HTTP-API  
**Project Type**: Backend-Domain-Fix  
**Performance**: O(24 × Linien) pro Remove—vernachlässigbar  
**Constraints**: HTTP-Schema unverändert; nur Fehlertexte ggf. neu

## Constitution Check

| Gate | Status |
|------|--------|
| Domain-First | Pass – entspricht offizieller Mühle-Regel |
| SpecKit | Pass |
| Testbare Kernlogik | Pass – **neue Unit-Tests** für Remove + Mill-Ausnahme **Pflicht** |
| Schichtentrennung | Pass – Logik in Services, nicht in Gin |
| Einfachheit | Pass – kleine Hilfsfunktionen, keine neue Abstraktionsschicht |

## Project Structure

```text
pkg/muehle/domain/services/
├── boardService.go    # ggf. Mill-Linien-Konstante + IsFieldInClosedMill, EnemyHasStoneOutsideMill
├── gameService.go     # RemoveStone: neue Guards
└── *_test.go          # neue Fälle
```

## Complexity Tracking

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| — | — | — |

## Phase 0 & 1 Outputs

| Artefakt | Pfad |
|----------|------|
| Research | [research.md](./research.md) |
| Data model | [data-model.md](./data-model.md) |
| Contract | [contracts/remove-stone-behavior.md](./contracts/remove-stone-behavior.md) |
| Quickstart | [quickstart.md](./quickstart.md) |

## Implementation Notes

1. **Mill lines**: Array `[][]int` wie in `HasPlayerThreeStones`—optional als
   Package-Level-Konstante teilen, um Duplikat zu vermeiden.
2. **`IsFieldPartOfClosedMill(board, index, color)`**: Feld `index` muss `color`
   haben; ∃ Zeile in Mill-Lines mit index ∈ Zeile und alle drei Felder `color`.
3. **`EnemyHasStoneOutsideMill(board, enemyColor)`**: ∀ Felder mit `enemyColor`,
   mindestens eines mit `!IsFieldPartOfClosedMill`.
4. **RemoveStone**: nach Prüfung „Feld gehört Gegner“ → wenn
   `EnemyHasStoneOutsideMill && IsFieldPartOfClosedMill(..., enemy)` →
   `return errors.New("...")`.
5. **Tests**: (a) Gegner: 3 in Mühle + 1 isoliert → Remove aus Mühle fehl, Remove
   isoliert OK. (b) Nur 3 in einer Mühle → Remove aus Mühle OK. (c) Überlappende
   Mühlen optional.
