# Metis Scheduler — Project- & Agent-richtlijnen (GEMINI.md)

Dit bestand bevat de belangrijkste richtlijnen, standaarden en workflow-afspraken voor de ontwikkeling van de **Metis Scheduler**. Zowel menselijke ontwikkelaars als AI-assistenten dienen zich strikt aan deze regels te houden om codekwaliteit, consistentie en betrouwbaarheid te waarborgen.

---

## 🎯 Projectdoel & Visie

Metis Scheduler is een high-performance Turnaround (TA) / Plant Shutdown scheduling engine geschreven in Go. Het doel is om sneller, eenvoudiger en beter te zijn dan Primavera P6, specifiek gefocust op dynamische, automatische Critical Path Method (CPM) berekeningen op basis van proces-systemen en operations-vensters.

De scheduler functioneert als een standalone service die nauw samenwerkt met de hoofdapplicatie **METIS** (Node.js/TypeScript).

---

## 🏗️ Architectuur & Richtlijnen

### 1. OpenAPI-First
*   **Volgorde van werken:** Wijzigingen in API-endpoints of datastructuren moeten **eerst** worden doorgevoerd in `/api/openapi.yaml`. Pas daarna mag de server- of handlercode worden aangepast.
*   **Strikte scheiding:** De API is stateless en communiceert uitsluitend via JSON.

### 2. Go (Backend) & Database
*   **Geen zware ORM's:** We gebruiken de snelle en lichte PostgreSQL driver `jackc/pgx/v5`. SQL-queries worden expliciet geschreven en beheerd via migrations (`internal/db/migrations/`).
*   **Go-idomatische code:**
    *   Gebruik `go fmt` voor automatische opmaak.
    *   Hanteer expliciete foutafhandeling (`if err != nil`). Vermijd panics in request-handlers; geef in plaats daarvan nette HTTP-foutcodes terug.
    *   Gebruik altijd `context.Context` bij database-operaties en netwerkverzoeken.
*   **Pakketstructuur:**
    *   `cmd/server/`: Entrypoint van de server. Houd deze minimaal.
    *   `internal/cpm/`: De pure CPM-reken-engine (vrij van database- of HTTP-dependencies).
    *   `internal/db/`: Databaseconnectiviteit en schema-migraties.
    *   `internal/server/`: HTTP-handlers en routing-logica.

### 3. CPM Reken-engine Standstandaard (`internal/cpm/`)
*   **Forward & Backward Pass:** De forward pass berekent de vroege tijden. De backward pass berekent de late tijden via een recursieve successor map om te garanderen dat `LateFinish` altijd de exacte minimale `LateStart` van alle opvolgers pakt. Wijzigingen hierin moeten met uiterste precisie en volledige testdekking worden doorgevoerd.
*   **Ready To Start Logica:** Een taak is pas `ready_to_start` als deze de status `NOT_STARTED` heeft en al zijn directe voorgangers `COMPLETED` zijn.

---

## 🛠️ Ontwikkel- & Validatiestroom

### 1. Wijzigingen aanbrengen
Houd wijzigingen chirurgisch en gefocust op de specifieke taak. Voorkom onnodige refactoring van omliggende, niet-gerelateerde code.

### 2. Testen & Kwaliteitscontrole
Elke wijziging aan de CPM-engine of database-laag vereist validatie:
*   **Unit Tests:** Voer tests uit via:
    ```bash
    go test ./...
    ```
*   **Dependencies:** Zorg dat de dependencies up-to-date en schoon zijn na wijzigingen in imports:
    ```bash
    go mod tidy
    ```
*   **Lokaal Compileren:** Controleer altijd of de applicatie succesvol compileert alvorens af te ronden:
    ```bash
    go build -o /dev/null ./cmd/server/main.go
    ```

---

## 💡 Voorstel: Wat hoort er nog meer thuis in dit bestand?

Om dit bestand in de toekomst nog waardevoller te maken voor de samenwerking, adviseren we de volgende onderwerpen toe te voegen zodra de implementatie vordert:

1.  **Migratiebeleid & Schema-afspraken:**
    *   Hoe schrijven we nieuwe migrations in `internal/db/migrations/`? (Bijvoorbeeld: altijd een up- en down-migratie, naamgevingsconventie zoals `002_add_calendars.sql`).
2.  **Kalender- & Tijdzone-standaarden:**
    *   Hoe gaan we om met kalenders en shift-schema's? (Bijvoorbeeld: "Alle tijden in de CPM-engine worden berekend in relatieve uren/dagen ten opzichte van de projectstart, of we gebruiken UTC-timestamps").
3.  **METIS-integratiecontract:**
    *   Definieer hier de authenticatiemethode (bijv. API-keys, JWT) en de exacte payload-structuur voor de synchronisatie tussen METIS en de Scheduler.
4.  **Mocking & Integratietesten:**
    *   Richtlijnen voor het mocken van de database (`pgx`) in unit-tests om de CPM-engine en handlers onafhankelijk te kunnen testen zonder dat een actieve PostgreSQL database vereist is.
5.  **Foutafhandeling & Logging-standaard:**
    *   Welke logging-bibliotheek of structured-logging formaat (bijv. `slog` uit de Go standard library) we gebruiken voor productie-observability.
