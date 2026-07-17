# Metis Scheduler

Metis Scheduler is een geavanceerde en snelle planning- en reken-engine geschreven in Go, gebaseerd op de **Critical Path Method (CPM)**. De engine is specifiek ontworpen om complexe projectnetwerken met parallelle paden, merge-points en dynamische taakstatussen door te rekenen.

Naast de traditionele CPM-berekeningen (`Early Start`, `Late Finish`, `Total Float`), beschikt deze engine over een unieke, intelligente `ReadyToStart` (Kan Starten) vlag die realtime bepaalt welke taken daadwerkelijk uitvoerbaar zijn op de werkvloer.

---

## 🚀 Features

*   **CPM Reken-engine:** Volledige forward en backward pass berekening voor netwerkplanningen.
*   **Kritieke Pad Analyse:** Directe identificatie van taken zonder speling (`Total Float = 0`).
*   **Smart "Ready to Start" Logica:** Filtert taken die direct kunnen starten op basis van de live status van hun specifieke voorgangers, onafhankelijk van de algehele projectvoortgang.
*   **Lokale PostgreSQL Integratie:** Robuuste opslag via de snelle en veilige `pgx` driver.
*   **Primavera Ready:** Het database-schema is functioneel al voorbereid op geavanceerde relatietypes (`SS`, `FF`, `SF`) en `lag_days` (wachttijden).
*   **Ingebouwde Swagger UI:** Interactieve API-documentatie direct bereikbaar via de browser.

---

## 🛠️ Tech Stack

*   **Language:** Go (v1.23+)
*   **Database:** PostgreSQL (Lokaal)
*   **Database Driver:** `jackc/pgx/v5`
*   **API Router:** Go Native `net/http` Multiplexer
*   **API Documentation:** OpenAPI 3.0 / Swagger UI

---

## ⚙️ Installatie & Setup

### 1. Database Voorbereiding
Zorg dat je een lokale PostgreSQL database hebt draaien genaamd `metis_scheduler`. Voer het volgende SQL-schema uit om de tabellen aan te maken:

```sql
CREATE TABLE IF NOT EXISTS projects (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS tasks (
    id VARCHAR(50) NOT NULL,
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    duration INT NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'NOT_STARTED',
    PRIMARY KEY (project_id, id)
);

CREATE TABLE IF NOT EXISTS task_dependencies (
    project_id UUID NOT NULL,
    task_id VARCHAR(50) NOT NULL,
    predecessor_id VARCHAR(50) NOT NULL,
    relation_type VARCHAR(10) NOT NULL DEFAULT 'FS', 
    lag_days INT NOT NULL DEFAULT 0,                 
    PRIMARY KEY (project_id, task_id, predecessor_id),
    FOREIGN KEY (project_id, task_id) REFERENCES tasks(project_id, id) ON DELETE CASCADE,
    FOREIGN KEY (project_id, predecessor_id) REFERENCES tasks(project_id, id) ON DELETE CASCADE
);
```

### 2. Configuratie
Pas in `internal/db/database.go` de database connectiestring aan met jouw lokale Postgres inloggegevens:

```go
connStr := "postgres://GEBRUIKER:WACHTWOORD@localhost:5432/metis_scheduler"
```

### 3. Server Starten
Zorg dat je in de hoofdmap van het project staat en voer uit:

```bash
# Dependencies downloaden en opschonen
go mod tidy

# Server opstarten
go run cmd/server/main.go
```

De server draait nu op: `http://localhost:8080`

---

## 🎯 API Endpoints

| Methode | Endpoint | Omschrijving |
| :--- | :--- | :--- |
| `GET` | `/swagger` | Open de interactieve Swagger UI documentatie |
| `POST` | `/projects/{id}/schedule` | Bereken de volledige CPM planning voor een project |
| `GET` | `/projects/{id}/critical-path` | Haal enkel de taken op die op het kritieke pad liggen |

---

## 📂 Project Structuur

```text
├── api/
│   └── openapi.yaml          # OpenAPI/Swagger specificatie
├── cmd/
│   └── server/
│       └── main.go           # Entrypoint van de applicatie
├── internal/
│   ├── cpm/
│   │   └── cpm.go            # De wiskundige CPM reken-engine
│   ├── db/
│   │   └── database.go       # PostgreSQL verbinding & database logica
│   └── server/
│       ├── handlers.go       # HTTP API Request handlers
│       └── server.go         # Server initialisatie
├── .gitignore
├── go.mod
└── README.md
```