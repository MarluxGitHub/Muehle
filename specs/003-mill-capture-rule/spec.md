# Feature Specification: Korrektes Schlagen bei Mühlen

**Feature Branch**: `003-mill-capture-rule`  
**Created**: 2026-03-18  
**Status**: Draft  
**Input**: User description: "aktuell kann ich aus einer Mühle einen stein entfernen beim schlagen das ist nicht erlaubt laut mühle regeln"

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Kein Stein aus gegnerischer Mühle (Priority: P1)

Als **Spieler**, der gerade eine Mühle geschlossen hat und schlagen darf, möchte
ich **keinen Stein des Gegners entfernen können**, der **Teil einer geschlossenen
Mühle** (drei in einer Reihe) des Gegners ist—solange der Gegner **mindestens
einen Stein hat, der in keiner Mühle liegt**.

**Why this priority**: Entspricht den gängigen Mühle-Regeln; aktuelles Verhalten
ist regelwidrig und frustriert faire Partien.

**Independent Test**: Auf einem Brett, auf dem der Gegner sowohl freie Steine als
auch Steine in Mühlen hat: nur ein Ziel auf freiem Stein ist erlaubt; Versuch,
einen Stein in einer gegnerischen Mühle zu wählen, schlägt fehl mit klarer
Rückmeldung.

**Acceptance Scenarios**:

1. **Given** der Gegner hat Steine außerhalb von Mühlen und in Mühlen, **When**
   der aktive Spieler nach Mühle schlägt und einen Stein **in einer gegnerischen
   Mühle** wählt, **Then** ist der Zug **ungültig** (kein Stein wird entfernt).
2. **Given** gleiche Situation, **When** der Spieler einen Stein wählt, der
   **nicht** in einer gegnerischen Mühle liegt, **Then** wird der Stein
   **regelkonform** entfernt und das Spiel geht weiter.
3. **Given** mehrere gegnerische Steine liegen außerhalb von Mühlen, **When**
   geschlagen wird, **Then** darf **nur** von diesen entfernt werden (nicht aus
   Mühlen).

---

### User Story 2 - Ausnahme: alle Steine in Mühlen (Priority: P2)

Als **Spieler** in einer seltenen Stellung, in der **jeder** Stein des Gegners
Teil mindestens einer **geschlossenen Mühle** ist, möchte ich **trotzdem schlagen
können**, damit das Spiel nicht stecken bleibt.

**Why this priority**: Standard-Ausnahme in der Turnier-/Volksregel; ohne sie
könnte Schlagen unmöglich werden.

**Independent Test**: Konstruierte oder gefundene Stellung „alle gegnerischen
Steine in Mühlen“: Entfernen eines **beliebigen** gegnerischen Steins ist
erlaubt.

**Acceptance Scenarios**:

1. **Given** **alle** verbleibenden Steine des Gegners liegen in geschlossenen
   Mühlen, **When** der Spieler nach Mühle schlägt, **Then** darf **ein**
   gegnerischer Stein **beliebig** entfernt werden (auch aus einer Mühle).

### Edge Cases

- Stein gehört zu **mehreren** überlappenden Mühlen: gilt weiterhin als „in
  Mühle“ für die Schlagregel.
- Nach dem Entfernen wechselt der Zug; keine doppelte Schlagpflicht.
- Ungültiger Schlagversuch: Spielzustand unverändert, klarer Hinweis (z. B.
  Fehlermeldung), wer am Zug ist.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: Befindet sich das Spiel in der Phase **„Stein des Gegners nach
  Mühle entfernen“**, MUST das System **verbieten**, einen gegnerischen Stein zu
  entfernen, der **Bestandteil einer geschlossenen Mühle** des Gegners ist—**so
 fern** der Gegner **mindestens einen Stein besitzt, der in keiner geschlossenen
  Mühle liegt**.
- **FR-002**: Erfüllt der Gegner die Bedingung **„alle seine Steine liegen in
  geschlossenen Mühlen“** (kein schlagbarer „freier“ Stein verfügbar), MUST das
  System **einen beliebigen** gegnerischen Stein zum Entfernen zulassen.
- **FR-003**: Die Regeln für **Setzen, Ziehen, Mühle bilden** und **Spielende**
  MUST **unverändert** bleiben, außer der hier beschriebenen **Schlaglogik**.
- **FR-004**: Jeder **unzulässige** Schlagversuch MUST **ohne** Änderung von Brett
  und Zugrecht **ablehnbar** sein mit **nachvollziehbarer** Rückmeldung an den
  Spieler.

### Key Entities

- **Geschlossene Mühle**: Drei eigene Steine des Gegners in einer geraden Linie
  auf dem Brett (laut gültigem Mühle-Brett).
- **Schlagbarer Stein**: Gegnerstein, der gemäß FR-001/FR-002 entfernt werden darf.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: In **100 %** der testbaren Stellungen mit „Gegner hat freie und
  Mühlen-Steine“ führt ein Versuch, einen **Mühlen-Stein** zu schlagen, zu
  **Ablehnung**; nur **freie** Steine sind wählbar.
- **SC-002**: In der dokumentierten **Ausnahmestellung** (alle Gegnersteine in
  Mühlen) ist **mindestens ein** Schlag möglich.
- **SC-003**: Bestehende **regelkonforme** Partien (ohne diesen Fehler) verhalten
  sich **wie bisher**, außer dass zuvor **falsche** Schläge nicht mehr möglich
  sind.

## Assumptions

- **Mühle-Definition** (drei in einer Reihe auf den definierten Linien)
entspricht der **bereits im Produkt** verwendeten Logik; diese Feature-Anforderung
betrifft **nur die Auswahl erlaubter Schlagziele**.
- Die **Ausnahme „alle Steine in Mühlen“** entspricht der **üblichen
Turnier-/Hausregel**; ohne sie wäre das Spiel in Randfällen nicht fortsetzbar.
