# Feature Specification: OpenAPI-Dokument & Auslieferungsroute

**Feature Branch**: `002-openapi-spec-route`  
**Created**: 2026-03-18  
**Status**: Draft  
**Input**: User description: "erstelle eine openapi specification zu den gin gonic routing und stelle dies unter einer route bereit"

## Clarifications

### Session 2026-03-18

- Q: Soll zusätzlich zur OpenAPI-Datei eine **Swagger/OpenAPI-UI** die API
  interaktiv anzeigen? → A: **Ja—Swagger-Ökosystem und eine UI** zur Anzeige
  (und idealerweise „Try it out“) der dokumentierten Endpunkte.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Maschinenlesbare API-Beschreibung abrufen (Priority: P1)

Als **Entwickler oder Integrationstool** möchte ich die Mühle-HTTP-API als
**OpenAPI-Dokument (3.x)** über eine **feste GET-Route** auf demselben Server
abrufen, damit ich Clients generieren, in Postman importieren oder die API
ohne Markdown lesen kann.

**Why this priority**: Ohne erreichbare Spec fehlt die zentrale technische
Referenz neben dem laufenden Dienst.

**Independent Test**: Mit einem HTTP-Client nur die Basis-URL und die
dokumentierte Spec-Route—Antwort ist gültiges OpenAPI (parsierbar) und listet
alle öffentlichen Spiel-Endpunkte.

**Acceptance Scenarios**:

1. **Given** der Server läuft, **When** ein Client die dokumentierte Route per
   GET aufruft, **Then** erhält er ein **OpenAPI-3.x-konformes** Dokument
   (JSON oder YAML gemäß Projektentscheid im Plan).
2. **Given** die Spec, **When** sie validiert wird (OpenAPI-Parser), **Then**
   treten **keine Schemafehler** auf, die das Dokument unbrauchbar machen.
3. **Given** die öffentlichen Spiel-Routen (`/games`, `/games/{gameId}/…`), **When**
   man die Spec mit der Router-Tabelle vergleicht, **Then** ist **jede**
   öffentliche Spiel-Operation mit Methode, Pfad und grober Request/Response-Beschreibung abgedeckt.

---

### User Story 2 - Swagger UI im Browser (Priority: P1)

Als **Entwickler** möchte ich eine **Swagger UI** (oder funktional gleichwertige
**OpenAPI-3-kompatible UI**) im **Browser** öffnen, damit ich alle Operationen
auf einen Blick sehe und optional **direkt gegen den Server ausprobieren** kann.

**Why this priority**: Entspricht der ausdrücklichen Produktvorgabe; ergänzt die
rohe Spec-Datei um nutzbare Visualisierung.

**Independent Test**: Nach Aufruf der dokumentierten **UI-URL** (ohne
Separatdownload) werden alle in der OpenAPI beschriebenen Pfade sichtbar; die UI
bezieht die Spec von der bekannten OpenAPI-Route (oder eingebetteter URL).

**Acceptance Scenarios**:

1. **Given** der Server läuft, **When** ein Nutzer die **UI-Route** im Browser
   öffnet, **Then** erscheint die **Swagger-typische** Darstellung (Gruppierung
   nach Pfaden, Methoden, Parameterbeschreibungen).
2. **Given** die UI, **When** sie lädt, **Then** verwendet sie dieselbe
   **OpenAPI-3.x-Quelle** wie FR-001/FR-002 (kein zweites widersprüchliches
   Schema).
3. **Given** CORS für die Spiel-API, **When** „Try it out“ von der UI aus
   genutzt wird, **Then** sind Anfragen an `localhost` **technisch möglich**
   (kein blockierendes CORS-Setup nur für die UI-Seite—Details im Plan).

---

### User Story 3 - Einheitliche Auffindbarkeit (Priority: P2)

Als **Nutzer der API-Doku** möchte ich, dass **OpenAPI-Route** und **Swagger-UI-Route**
in der **Benutzerdokumentation** (`docs/api.md` oder README) genannt werden.

**Why this priority**: Reduziert Support-Fragen.

**Independent Test**: Neuer Leser findet innerhalb von 1 Minute **Spec-URL** und
**UI-URL** aus der Repo-Doku (Stichworte OpenAPI, Swagger).

**Acceptance Scenarios**:

1. **Given** die Projekt-Doku, **When** jemand nach „OpenAPI“, „Swagger“ oder
   „API-Doku“ sucht, **Then** werden **beide** konkreten URLs genannt.

### Edge Cases

- Server im Produktionsmodus: **Spec-Route und UI-Route** bleiben erreichbar
  (Standard: **öffentlich**, keine Debug-only-Doku), sofern nicht explizit
  anders dokumentiert.
- Große Spec: eine YAML/JSON-Quelle; UI ist **Thin Client** auf dieselbe Quelle.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: Das System MUST ein **OpenAPI-3.x**-Dokument bereitstellen, das
  **alle** aktuell exponierten **Spiel-API-Endpunkte** beschreibt
  (Pfade, HTTP-Methoden, zentrale Parameter, typische Antwortcodes inkl. 404).
- **FR-002**: Die Spec MUST über **mindestens eine** dedizierte **GET-Route**
  ausgeliefert werden (Pfad im Plan/Doku); **Content-Type** passend zu YAML/JSON.
- **FR-003**: Die Spec-Antwort MUST den passenden **Content-Type** tragen
  (`application/json` bzw. `application/yaml` je nach Format).
- **FR-004**: **servers**-Eintrag in der OpenAPI MUST für typische lokale Nutzung
  sinnvoll vorbelegt sein (z. B. `http://localhost:8080`), ohne Geheimnisse.
- **FR-005**: Bei Änderung der HTTP-Routen MUST OpenAPI **mitgepflegt** werden;
  keine dauerhafte Divergenz zwischen Router und Spec.
- **FR-006**: Das System MUST **Swagger UI** (oder **OpenAPI-3-kompatible**
  gleichwertige UI) unter **mindestens einer dokumentierten Browser-Route**
  bereitstellen; die UI MUST die **gleiche** OpenAPI-Definition wie FR-001/002
  einbinden (per URL zur Spec oder eingebettet—Plan entscheidet).

### Key Entities

- **OpenAPI-Dokument**: Maschinenlesbare Schnittstellenbeschreibung.
- **Spec-Auslieferungsroute**: GET-Endpunkt für Rohspec.
- **Swagger UI**: Interaktive Darstellung derselben Spec im Browser.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: **100 %** der Spiel-API-Routen erscheinen als **Operation** in der
  OpenAPI-Datei (Review-Checkliste).
- **SC-002**: Abruf der Spec-Route liefert lokal **unter 2 Sekunden** vollständige
  Antwort.
- **SC-003**: Doku nennt **Spec-URL** und **Swagger-UI-URL**; Nutzer erreicht die
  UI **ohne** manuelles Herunterladen der YAML von Drittanbietern.
- **SC-004**: Im Browser zeigt die UI **alle** in OpenAPI definierten
  Operationen; manuelle Stichprobe **≥ 5** Operationen sichtbar und benannt wie
  in der Spec.

## Assumptions

- Umsetzung auf dem **bestehenden** Go/Gin-Server; **Swagger UI** typischerweise
  als **statische Assets** (embed) oder **CDN-Build** mit `url` zur lokalen
  `/openapi.yaml`—Details im Plan.
- **Spec- und UI-Routen** sind **öffentlich** lesbar; keine sensiblen Daten in
  der OpenAPI.
- **„Swagger“** meint hier die übliche Kombination **OpenAPI 3 + Swagger UI**,
  nicht veraltetes Swagger 2.0-Format.
