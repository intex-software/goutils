## 1. Arbeitsweise: Review vs. Editing

- **Review-Modus:** Ein "Review" oder "Audit" bedeutet ausdrücklich **nur Analyse und Bericht**.
- **Änderungsverbot:** Es dürfen niemals eigenmächtig Änderungen am Code vorgenommen werden, wenn nur nach einem Review gefragt wurde.
- **Vorschläge:** Korrekturen müssen als Empfehlung/Vorschlag im Text präsentiert werden. Änderungen erfolgen erst nach expliziter Bestätigung durch den User.

---

## 2. Architektur & Utility-Designprinzipien (Go)

- **Single Responsibility:** Jedes Package in `goutils` kapselt genau eine Domäne oder Funktionalität (z. B. `csv`, `secret`, `ssl`, `krypta`, `obfuscate`).
- **Zero Surprises & Explizites Error-Handling:**
  - Fehler sind Werte. Kein `panic` in Bibliotheksfunktionen (außer explizite `Must...`-Initialisierer).
  - Nutze konsequentes Error-Wrapping mit `%w` (`fmt.Errorf("package: operation failed: %w", err)`).
- **Interface-Design & Komposition:**
  - „Accept interfaces, return structs.“
  - Halte Interfaces so klein und interoperabel wie möglich (z. B. `io.Reader`, `io.Writer`).
- **Control Flow & Guard Clauses:**
  - Flache Hierarchien, Early Returns. **Vermeide `else`** nach Fehlerrückgaben.
- **Concurrency & Thread-Safety:**
  - Bibliotheksfunktionen, die Shared State verwalten, MÜSSEN threadsicher sein (Mutex, Atomic, Read-only Structures).
  - Keine versteckten globalen Zustände.

---

## 3. Code-Style & Namenskonventionen

- **Idiomatisches Go:** Folge den offiziellen [Effective Go](https://go.dev/doc/effective_go) Richtlinien.
- **Benennung:**
  - `PascalCase` für exportierte Typen, Konstanten und Funktionen.
  - `camelCase` für interne Implementierungsdetails.
  - Kurze, prägnante Bezeichner in kleinen Scopes (`r` für Reader, `err` für Error).
- **Performance & Allokationen:**
  - In Utility-Kernroutinen Allokationen minimieren (z. B. `sync.Pool`, Re-Use von Buffern, Vermeidung unnötiger Casts).

---

## 4. Testing-Standards

- **Table Driven Tests:** Verwende tabellengesteuerte Tests für alle Utility-Funktionen.
- **Test-Abdeckung:**
  - Happy Path, Edge Cases (nil Pointers, leere Slices/Strings, Grenzwerte).
  - Sad Paths: Fehlerfälle gezielt testen und Fehlertypen prüfen (`errors.Is`, `errors.As`).
- **Isolation:** Keine Abhängigkeiten zu externen Netzwerken oder realen Dateisystemen ohne saubere Test-Fixtures/Mocks.

---

## 5. Dokumentation & Technical Writing

- **Sprache:** Präzise, entwicklernah und auf Deutsch.
- **Godoc:** Vollständige `godoc`-Kommentare für alle exportierten Packages, Typen und Funktionen.
- **Beispiele:** Erstelle bei nichttrivialen Utility-Funktionen ausführbare Tests/Beispiele (`Example...`).
