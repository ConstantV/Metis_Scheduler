# Metis Scheduler — Project Plan

## Doel
Een scheduling tool vergelijkbaar met Primavera Project Planner: eenvoudiger, sneller, beter.
Gebouwd in Go. Geïntegreerd met METIS (Node.js/TypeScript/PostgreSQL).
Klant: PMXL (eigen gebruik).

---

## Architectuur

### Server (Go)
- **OpenAPI-first**: spec schrijven vóór implementatie
- **REST API**, stateless, JSON
- **PostgreSQL** als database (via `pgx` of `sqlc`, geen zware ORM)
- Database-laag zo ontworpen dat andere backends later haalbaar zijn

### Client (web)
- Volledig losgekoppeld van server — praat uitsluitend via de API
- MVP: lichte webstack voor testen (geen zware Gantt-library)
- Later: uitbreiden op basis van behoeften

### Integratie met METIS
- **METIS → Scheduler**: projecten, scope-items, resources, operaties komen uit METIS
- **Scheduler → METIS**: voortgang en executiestatus gaan terug naar METIS
- Integratiemechanisme: REST API (twee losse services die via API communiceren)

---

## MVP Scope

### Datamodel (kern)
- **Project** — naam, startdatum, kalender
- **WBS** — hiërarchische werkverdeling
- **Activiteit** — naam, duur, vroegst/laatste start/einde, % complete
- **Afhankelijkheid** — Finish-to-Start (FS), Start-to-Start (SS), Finish-to-Finish (FF)
- **Kalender** — werkdagen, feestdagen, uurroosters
- **Resource** — persoon of middel, capaciteit
- **Resource-toewijzing** — activiteit ↔ resource, belasting
- **Baseline** — snapshot van plan op een moment

### Berekeningen
- **CPM** (Critical Path Method) — vroegste/laatste start en einde per activiteit
- **Float** — total float en free float
- **Kritiek pad** markeren
- **Resource-belasting** per periode

### API (OpenAPI)
Minimale endpoints voor MVP:
- CRUD: projects, wbs, activities, dependencies, calendars, resources, assignments
- `POST /projects/{id}/schedule` — herbereken CPM
- `GET /projects/{id}/gantt` — Gantt-data voor frontend
- `GET /projects/{id}/critical-path` — kritiek pad
- `GET /projects/{id}/baseline` — baseline ophalen/vergelijken
- `PATCH /activities/{id}/progress` — voortgang bijwerken
- METIS-sync endpoints (import/export)

### Frontend (MVP / test)
- Lichte webstack — HTML + vanilla JS of HTMX
- Gantt-weergave: eenvoudige SVG of canvas rendering (zelf gebouwd, minimaal)
- Activiteitenlijst met % complete
- Kritiek pad visueel markeren

---

## Fasering

### Fase 1 — Fundament
- [ ] Go project initialiseren (`go mod init`)
- [ ] OpenAPI spec schrijven (kern-entiteiten)
- [ ] Database schema (PostgreSQL) — migrations
- [ ] CRUD API implementeren
- [ ] CPM-algoritme implementeren en testen

### Fase 2 — Planning features
- [ ] Kalender-logica (werkdagen, feestdagen)
- [ ] Resource-toewijzing en belastingsberekening
- [ ] Baselines
- [ ] Voortgang (% complete) bijwerken

### Fase 3 — Frontend
- [ ] Lichte Gantt-weergave (SVG/canvas)
- [ ] Activiteitenlijst
- [ ] Kritiek pad visualisatie

### Fase 4 — METIS-integratie
- [ ] Import vanuit METIS (projecten, scope-items, resources)
- [ ] Export voortgang naar METIS
- [ ] Authenticatie / API-keys

---

## Tech Stack

| Component | Keuze | Reden |
|-----------|-------|-------|
| Taal (server) | Go | Snel, eenvoudig, goede concurrency |
| API-stijl | REST / OpenAPI 3.x | Interoperabel, documenteerbaar |
| Database | PostgreSQL | Bewezen, flexibel, zelfde als METIS |
| DB-toegang | `pgx` of `sqlc` | Licht, performant |
| Frontend (MVP) | HTML + vanilla JS / HTMX | Minimaal, snel te testen |
| Gantt (MVP) | Zelf gebouwd (SVG) | Geen library-afhankelijkheid |

---

## Open vragen
- Authenticatie: JWT, API-key, of via METIS-sessie?
- Multiproject: meerdere projecten in één scheduler-instance of één per klant?
- Eenheden: werken we in dagen, uren, of beide?
- Kalender: één globale kalender of per project/resource?
