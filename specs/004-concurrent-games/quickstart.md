# Quickstart: Zwei Partien parallel

**Verifiziert:** Der Ablauf (zwei Partien, getrennte IDs, Züge nur in einer Partie) wird durch die Integrationstests in **`pkg/muehle/interfaces/gin_integration_test.go`** abgedeckt (`go test ./pkg/muehle/interfaces/...`). Optional manuell mit laufendem Server wie unten.

Server starten (`go run ./cmd/server` oder Binary). Terminal 1 und 2 simulieren zwei Gruppen.

## Partie A

```bash
A=$(curl -s -X POST http://localhost:8080/games | jq -r .id)
curl -s -X POST "http://localhost:8080/games/$A/players" -d "playerName=W1"
curl -s -X POST "http://localhost:8080/games/$A/players" -d "playerName=S1"
# … Züge mit gameId=$A
```

## Partie B (während A läuft)

```bash
B=$(curl -s -X POST http://localhost:8080/games | jq -r .id)
curl -s -X POST "http://localhost:8080/games/$B/players" -d "playerName=W2"
curl -s -X POST "http://localhost:8080/games/$B/players" -d "playerName=S2"
```

`GET http://localhost:8080/games/$A/board` und `…/$B/board` zeigen **unterschiedliche** Bretter nach unterschiedlichen Zügen.

## Unbekannte ID

```bash
curl -s -o /dev/null -w "%{http_code}" http://localhost:8080/games/00000000-0000-0000-0000-000000000000/board
# Erwartung: 404 (sofern diese UUID nicht zufällig existiert—besser zufällige UUID nutzen)
```
