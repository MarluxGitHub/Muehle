# Data Model: Mehrere Spiele gleichzeitig

## Registry (Anwendungsebene)

| Konzept | Bedeutung |
|---------|-----------|
| **Partie-Referenz (`gameId`)** | String, kanonisch UUID (RFC 4122); von Clients in allen URLs unter `/games/{gameId}/…` verwendet. |
| **Aktive Partie** | Eintrag in der Registry: Referenz → genau eine gekapselte Spiel-Session (`Application`). |

## Partie (Spiel-Session)

Wie in Spec **Key Entities**; technisch pro Eintrag:

- Eigener **Brettzustand**, **Phase** (Setzen/Ziehen/…), **Spielerliste** mit Geheimcodes.
- **Keine** gemeinsamen veränderlichen Strukturen mit anderen Partien.

## Validierung

- `gameId` muss gültiges UUID-Format haben und in der Registry existieren, sonst keine Spieloperation.
- Anlegen einer Partie erzeugt neue UUID und leeren Spielzustand (Warten auf Spieler).

## Beziehungen

```text
GameRegistry
  └── 0..N Partie(gameId) → Application → GameService → Board / Players / State
```

Keine Relationen zwischen Partien untereinander.
