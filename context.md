# Metis Industrial Scheduler - Project Context

## Doel
Een high-performance Turnaround (TA) / Plant Shutdown scheduling engine gebouwd in Go en PostgreSQL. Het lost de beperkingen van Primavera P6 op door middel van dynamische, automatische CPM-berekeningen op basis van proces-systemen en operations-vensters.

## Huidige Database Architectuur (PostgreSQL)
- `projects`: Basis projectinformatie.
- `standard_tasks`: Catalogus / Takenboek met nummerreeksen (`PRE_TA`, `SHUTDOWN`, `MAINTENANCE`, `STARTUP`).
- `process_systems`: Fysieke fabriekssystemen die fasen doorlopen.
- `system_phase_windows`: Bepaalt per systeem wanneer fasen starten/eindigen (`start_hour`, `end_hour`).
- `scope_items`: Koppelt scope-nummers en equipment aan processystemen.
- `tasks`: De CPM-taken met kolommen als `early_start`, `early_finish`, `late_start`, `late_finish`, `total_float`, `is_critical`, en `ready_to_start`.
- `task_dependencies`: Netwerklogica tussen taken.

## CPM Reken-engine Logic (`internal/cpm/cpm.go`)
- **Forward Pass:** Recursieve berekening van vroege tijden op basis van `DependsOn`.
- **Backward Pass:** Volledig recursieve berekening via een *Successor Map*. Garandeert dat `LateFinish` altijd de exacte minimale `LateStart` van alle opvolgers pakt. Voorkomt overschrijf-bugs bij meerdere opvolgers.
- **Ready To Start:** Een taak wordt pas `true` als status `NOT_STARTED` is én alle directe voorgangers op `COMPLETED` staan.

## Volgende Ontwikkelstap
De Go structs in `cpm.go` en de database database-queries (`repository`) uitbreiden zodat ze de nieuwe tabellen (`scope_items`, `system_phase_windows`, `standard_tasks`) ondersteunen. De CPM-engine moet bij de start van de Forward Pass het `start_hour` van het bijbehorende planningsvenster inladen als de minimale starttijd voor die specifieke taak.