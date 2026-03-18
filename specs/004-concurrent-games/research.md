# Research: 004-concurrent-games

## 1. Wie werden mehrere Partien technisch getrennt?

**Decision**: Eine **Registry** mappt eine **eindeutige ID** (UUID) auf genau eine **`Application`**-Instanz mit eigenem `GameService`/Brett.

**Rationale**: Erfüllt FR-002 (äußere Referenz), FR-003/004 (Aktionen nur auf gewählter Instanz), FR-005 (neues Spiel = neuer Map-Eintrag).

**Alternatives considered**:

- **Globales Singleton-Spiel**: verworfen (würde Spec verletzen).
- **Persistenz (DB/Redis)**: laut Spec-Annahme nicht Teil dieses Features; später für Neustart-Festigkeit.

## 2. Nebenläufigkeit (mehrere Requests gleichzeitig)

**Decision**: `GameRegistry` nutzt **`sync.RWMutex`**; jede Partie ist eigene `Application`. Parallele Requests auf **verschiedene** `gameId` blockieren sich nur kurz bei Map-Lookup; Schreibzugriffe auf **dieselbe** Partie laufen seriell über Mutex der Registry beim `Get` nicht pro App—hier: eine `Application` pro Request-Pfad; gleiche Partie parallel könnte Race geben.

**Rationale**: Für typische Zwei-Spieler-Nutzung ausreichend. Spec verlangt keine harte Serialisierung pro Partie; falls nötig, später feinere Locking pro `Application`.

**Alternatives considered**: Mutex pro `Application`—optional in Tasks, wenn Tests Races zeigen.

## 3. Unbekannte Partie-Referenz (FR-006)

**Decision**: `uuid.Parse` fehlgeschlagen oder ID nicht in Map → **404** JSON `{"error":"game not found"}`, kein State-Change.

**Rationale**: Bereits in `resolveGame` umgesetzt; entspricht US3.

## 4. Teststrategie für SC-002 / SC-003

**Decision**: **Integrationstests**: mehrere `POST /games`, je zwei Spieler, Züge nur in Spiel 1, Board-GET für Spiel 2–5 unverändert; `POST` mit Random-UUID → 404, danach gültiges Spiel unverändert.

**Rationale**: Deckt HTTP-Vertrag und Registry ab ohne Domain zu duplizieren.

## 5. Ist-Abgleich HTTP-Handler (T002 / Implementierung)

**Verifiziert**: In `pkg/muehle/interfaces/gin.go` nutzen `postGamePlayers`, `postGameMoves`, `getGameState`, `getGameCurrentPlayer`, `getGameBoard` jeweils `resolveGame` → `Registry.Get`. `POST /games` (`postGames`) ist die einzige Route ohne `gameId`, die eine Partie anlegt. Keine globale Partie außerhalb der Registry.
