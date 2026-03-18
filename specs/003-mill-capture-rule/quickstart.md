# Quickstart: Schlagregel manuell prüfen

1. Partie bis **RemovingStone** führen (Mühle schließen).
2. Gegner hat **Mühle + mindestens einen Stein außerhalb**: Versuch, Stein **aus der Mühle** zu entfernen → muss **fehlschlagen**.
3. Dann Stein **außerhalb** der Mühle entfernen → **Erfolg**.

Automatisiert: nach Implementierung **`go test ./pkg/muehle/domain/...`** mit
konstruierten Brettern.
