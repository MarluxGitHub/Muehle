# Data Model: Schlagregel (Domain)

Keine neuen Persistenz-Entitäten.

## Regeln (fachlich)

| Prädikat | Bedeutung |
|----------|-----------|
| **ClosedMillLine** | Eine der 16 Standardlinien; alle 3 Felder = Farbe X |
| **FieldInEnemyMill** | Feldindex `i` hat Farbe Gegner **und** liegt auf mindestens einer ClosedMillLine für Gegner |
| **EnemyHasStoneOutsideMill** | ∃ Feld mit Gegnerfarbe, das **nicht** FieldInEnemyMill ist |

## RemoveStone (Erweiterung)

1. Bestehende Guards (State, Secret, aktueller Spieler, Feld = Gegner).
2. **Neu**: Wenn `EnemyHasStoneOutsideMill` **und** `FieldInEnemyMill(Zielfeld)` → **ablehnen**.
3. Wenn **nicht** `EnemyHasStoneOutsideMill` → Ziel wie bisher (jedes Gegnerfeld).

## Invarianten

- Kein Steinwechsel bei abgelehntem Remove.
- `Stones`-Zähler nur bei erfolgreichem Entfernen.
