# Research: Schlagen bei Mühlen

## R1: Definition „Stein liegt in einer Mühle“

**Decision**: Ein Feld mit gegnerischer Farbe gilt als **in einer geschlossenen
Mühle**, wenn es **Element mindestens einer** der 16 festen Linien ist, auf der
**alle drei Felder** die **gegnerische** Farbe tragen (identisch zu
`HasPlayerThreeStones`-Logik pro Linie, bezogen auf den Gegner).

**Rationale**: Konsistent mit bestehender Mühlen-Erkennung im Projekt.

## R2: Wann Mühlen-Steine geschützt sind

**Decision**: Vor `RemoveStone`: Wenn der Gegner **mindestens einen** Stein hat,
der **nicht** in einer geschlossenen Mühle liegt (`HasStoneOutsideMills`), dann
darf nur auf Felder geschlagen werden, die **nicht** `InClosedMill(enemy)` sind.
Sonst Fehler, Brett unverändert.

**Rationale**: FR-001.

## R3: Ausnahme

**Decision**: `HasStoneOutsideMills(enemy) == false` (jeder gegnerische Stein
liegt in mindestens einer vollen gegnerischen Mühle) → **jeder** gegnerische
Stein darf entfernt werden. FR-002.

**Rationale**: Verhindert Deadlock (z. B. nur noch Steine in überlappenden
Mühlen).

## R4: Ort der Logik

**Decision**: Hilfsfunktionen in **`BoardService`** (`IsFieldPartOfClosedMill`,
`EnemyHasStoneOutsideAnyMill` o. ä.); **`GameService.RemoveStone`** ruft sie vor
dem Entfernen auf.

**Rationale**: Brett/Mühlen-Geometrie bleibt in BoardService; GameService
orchestriert Zugrecht.

## R5: Fehlermeldung API

**Decision**: Explizite englische oder deutsche Fehlermeldung, z. B.
`cannot remove stone from mill`—für FR-004 nachvollziehbar.
