# Contract: Mehrere parallele Partien (HTTP)

Kanonische Routen wie `specs/001-rest-routing-cleanup/contracts/http-api.md`; hier nur **mehrspielbezogenes** Verhalten.

| Anforderung | Verhalten |
|-------------|-----------|
| FR-001 / neue Partie | `POST /games` → `201`, Body enthält eindeutige `id`. Jeder Aufruf legt **zusätzliche** Partie an; bestehende bleiben. |
| FR-002 | `id` ist die `gameId` für alle folgenden Pfade. |
| FR-003 | Nur Routen unter `/games/{gameId}/…` betreffen diese Partie. |
| FR-006 | Ungültiges UUID-Format oder unbekannte `gameId`: **`404`**, Body z. B. `{"error":"game not found"}`. Kein erfolgreiches Ausführen von Moves/Players auf fremder Partie. |
| Parallele Partien | Zwei gültige IDs A und B: Aktionen mit Pfad `…/A/…` ändern nur A; `GET …/B/board` bleibt unabhängig von Zügen in A. |

**Tests**: siehe `pkg/muehle/interfaces` (Integration) für Isolation und 404.
