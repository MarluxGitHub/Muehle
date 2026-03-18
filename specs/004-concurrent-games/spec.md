# Feature Specification: Mehrere Spiele gleichzeitig

**Feature Branch**: `004-concurrent-games`  
**Created**: 2026-03-18  
**Status**: Draft  
**Input**: User description: "implementiere das mehrere Spiele gleichzeitig stattfinden können"

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Eigene Partie, eigener „Raum“ (Priority: P1)

Zwei oder mehr Gruppen wollen jeweils eine eigene Mühle-Partie spielen, ohne sich gegenseitig zu stören. Jede Gruppe startet oder betritt eine **eindeutige Partie**; alle laufenden Partien existieren **nebeneinander**.

**Why this priority**: Ohne getrennte Partien kann nur eine Gruppe spielen—das ist der Kern des Wunsches.

**Independent Test**: Zwei Mal „neue Partie anlegen“ (oder äquivalent); es existieren zwei unterscheidbare Partien, beide können zum Spielen genutzt werden.

**Acceptance Scenarios**:

1. **Given** keine oder bereits laufende Partien, **When** eine zweite Partie gestartet wird, **Then** bleibt die erste Partie bestehen und beide sind unterscheidbar adressierbar.
2. **Given** Partie A und Partie B sind aktiv, **When** in Partie A ein erlaubter Spielzug erfolgt, **Then** ändert sich der sichtbare Stand von Partie B nicht.

---

### User Story 2 - Paralleles Spielen ohne Vermischung (Priority: P2)

Spieler in Partie A und Spieler in Partie B können **gleichzeitig** ziehen oder Steine setzen; das System ordnet jede Aktion **nur der richtigen Partie** zu.

**Why this priority**: Gleichzeitigkeit ist nur dann sinnvoll, wenn Aktionen strikt pro Partie isoliert sind.

**Independent Test**: Zwei Clients (oder simulierte Aufrufe) an zwei verschiedene Partien; nacheinander oder parallel Aktionen ausführen und prüfen, dass Brett, Zugrecht und Spielende je Partie unabhängig sind.

**Acceptance Scenarios**:

1. **Given** Partie A und B mit je zwei Spielern im laufenden Spiel, **When** in A ein Zug ausgeführt wird und kurz darauf in B ein Zug, **Then** entspricht der Stand von A nur den A-Zügen und der von B nur den B-Zügen.
2. **Given** Partie A ist beendet (Sieg einer Seite), **When** in Partie B weiter gespielt wird, **Then** bleibt B spielbar und der Ausgang von A hat keinen Einfluss auf B.

---

### User Story 3 - Falsche oder unbekannte Partie (Priority: P3)

Wer eine Partie nicht (mehr) kennt oder eine falsche Referenz nutzt, erhält ein **klares, sicheres Verhalten** (kein Zugriff auf fremde Partie, keine stillen Änderungen an anderer Partie).

**Why this priority**: Schützt vor Verwechslung und versehentlicher Manipulation.

**Independent Test**: Aktion mit gültiger Form aber unbekannter Partie-Referenz; erwartetes Ablehnen ohne Seiteneffekte auf andere Partien.

**Acceptance Scenarios**:

1. **Given** existiert nur Partie A, **When** jemand mit Referenz einer nicht existierenden Partie eine Aktion auslöst, **Then** schlägt die Aktion fehl und Partie A bleibt unverändert.

---

### Edge Cases

- Sehr viele gleichzeitig angelegte Partien: System darf begrenzen oder langsamer werden, **darf aber keine Partien vermischen**.
- Dieselbe Person öffnet zwei Browser-Tabs mit zwei Partien: beide Partien bleiben getrennt, solange die Referenzen unterschiedlich sind.
- Partie-Referenz verloren: Nutzer kann nicht ohne gültige Referenz an einer Partie teilnehmen (kein Raten fremder Partien).

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: Das System MUST es erlauben, **mehr als eine** Mühle-Partie **gleichzeitig** zu führen (jede Partie eigener Lebenszyklus: anlegen, spielen, enden).
- **FR-002**: Jede Partie MUST eine **eindeutige, von außen verwendbare Referenz** haben, mit der alle spielrelevanten Aktionen dieser Partie ausgeführt werden.
- **FR-003**: Aktionen (z. B. beitreten, Zug ausführen, Stand abfragen) MUST **ausschließlich** die Partie betreffen, deren Referenz mitgesendet wurde.
- **FR-004**: Das System MUST **keinen** Spielstand, kein Zugrecht und keinen Spielausgang zwischen zwei verschiedenen Partien **übertragen oder vermischen**.
- **FR-005**: Beim Anlegen einer neuen Partie MUST **keine** bestehende Partie automatisch beendet oder zusammengeführt werden.
- **FR-006**: Bei ungültiger oder unbekannter Partie-Referenz MUST die angefragte Aktion **ablehnen** ohne andere Partien zu verändern.

### Key Entities

- **Partie (Spiel-Session)**: Eine laufende oder wartende Mühle-Partie; hat Referenz, zugehörige Spielerrollen, Brettzustand und Phase nach Mühle-Regeln.
- **Spieler in einer Partie**: Zu genau einer Partie zugeordnet; Aktionen gelten nur innerhalb dieser Partie.

## Assumptions & Dependencies

- **Annahme**: „Gleichzeitig“ bedeutet logisch parallele Sessions auf derselben Server-Instanz; Nutzer erwarten keine globale Begrenzung auf eine Partie pro Server.
- **Annahme**: Persistenz über Neustarts ist nicht Teil dieser Spezifikation, sofern nicht später ergänzt (Fokus: Korrektheit und Isolation bei parallelen Partien).
- **Abhängigkeit**: Bestehende Mühle-Spielregeln gelten **pro Partie** unverändert.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Zwei unabhängige Beobachter können **zwei gleichzeitig laufende Partien** unterscheiden und je Partie den korrekten Stand nachvollziehen, nachdem in beiden Partien jeweils mindestens ein legaler Zug erfolgt ist.
- **SC-002**: In einem Test mit **mindestens fünf gleichzeitig aktiven Partien** führt eine Folge erlaubter Züge in Partie 1 **nicht** zu einer Änderung des Bretts in Partien 2–5.
- **SC-003**: **100 %** der Versuche, mit einer **bekannt falschen** Partie-Referenz eine schreibende Aktion auszulösen, enden ohne erfolgreiche Ausführung und ohne Änderung an einer **anderen**, gültigen Partie.
- **SC-004**: Organisatoren berichten, dass sie **ohne Wartezeit** auf andere Gruppen eine **eigene** Partie starten können (qualitativ: paralleles Spielen ist im Alltag nutzbar).
