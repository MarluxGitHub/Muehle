# Contract: POST …/moves (action=remove)

**Domain-Verhalten** (HTTP unverändert: Form `fieldIndex`, `secretCode`).

| Situation | HTTP-Body / Status (heute: meist 500 mit `error`-Text) |
|-----------|--------------------------------------------------------|
| Ziel ist gegnerischer Stein **in Mühle**, Gegner hat **freien** Stein | **Fehler**, kein State-Change; Domain: `errors.New("cannot remove stone from mill")` |
| Ziel ist freier gegnerischer Stein | **200**, wie bisher |
| Alle Gegnersteine nur in Mühlen | **200** für beliebiges gegnerisches Feld |

Tests: **`go test`** gegen `GameService.RemoveStone` / BoardService-Hilfen.
