package server

import (
	"encoding/json"
	"net/http"

	"metis-scheduler/internal/cpm"
)

// HandleSchedule haalt data uit de DB, berekent de CPM, slaat het op en geeft antwoord
func (s *Server) HandleSchedule(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	projectID := r.PathValue("id")
	if projectID == "" {
		http.Error(w, "Project ID is verplicht", http.StatusBadRequest)
		return
	}

	// 1. Haal alle taken op voor dit project uit de DB
	rows, err := s.db.Query(ctx,
		"SELECT id, name, duration, status FROM tasks WHERE project_id = $1",
		projectID,
	)
	if err != nil {
		http.Error(w, "Fout bij ophalen taken: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var tasks []cpm.Activity
	taskMap := make(map[string]*cpm.Activity)

	for rows.Next() {
		var t cpm.Activity
		if err := rows.Scan(&t.ID, &t.Name, &t.Duration, &t.Status); err != nil {
			http.Error(w, "Fout bij scannen taak: "+err.Error(), http.StatusInternalServerError)
			return
		}
		t.DependsOn = []string{} // Initialiseer lege slice
		tasks = append(tasks, t)
	}

	// Zet in een map voor snelle referentie bij het koppelen van dependencies
	for i := range tasks {
		taskMap[tasks[i].ID] = &tasks[i]
	}

	// 2. Haal de relaties (dependencies) op uit de DB en koppel ze aan de taken
	depRows, err := s.db.Query(ctx,
		"SELECT task_id, predecessor_id FROM task_dependencies WHERE project_id = $1",
		projectID,
	)
	if err != nil {
		http.Error(w, "Fout bij ophalen relaties: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer depRows.Close()

	for depRows.Next() {
		var taskID, predID string
		if err := depRows.Scan(&taskID, &predID); err != nil {
			http.Error(w, "Fout bij scannen relatie: "+err.Error(), http.StatusInternalServerError)
			return
		}
		if task, exists := taskMap[taskID]; exists {
			task.DependsOn = append(task.DependsOn, predID)
		}
	}

	if len(tasks) == 0 {
		http.Error(w, "Geen taken gevonden in de database voor dit project UUID", http.StatusNotFound)
		return
	}

	// 3. Voer de CPM berekening uit via de engine
	calculatedTasks := cpm.CalculateSchedule(tasks)

	// 4. Sla de berekende resultaten op in de database
	for _, t := range calculatedTasks {
		_, err := s.db.Exec(ctx, `
			UPDATE tasks 
			SET early_start = $1, early_finish = $2, late_start = $3, late_finish = $4, total_float = $5, ready_to_start = $6
			WHERE project_id = $7 AND id = $8`,
			t.EarlyStart, t.EarlyFinish, t.LateStart, t.LateFinish, t.TotalFloat, t.ReadyToStart,
			projectID, t.ID,
		)
		if err != nil {
			http.Error(w, "Fout bij opslaan rekenresultaten: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// 5. Geef de berekende planning terug als JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(calculatedTasks)
}

// HandleCriticalPath haalt direct de kritieke taken op uit de database
func (s *Server) HandleCriticalPath(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	projectID := r.PathValue("id")
	if projectID == "" {
		http.Error(w, "Project ID is verplicht", http.StatusBadRequest)
		return
	}

	// Taken op het kritieke pad hebben een TotalFloat van 0
	rows, err := s.db.Query(ctx, `
		SELECT id, name, duration, status, early_start, early_finish, late_start, late_finish, total_float, ready_to_start 
		FROM tasks 
		WHERE project_id = $1 AND total_float = 0`,
		projectID,
	)
	if err != nil {
		http.Error(w, "Fout bij ophalen kritieke pad: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var criticalTasks []cpm.Activity
	for rows.Next() {
		var t cpm.Activity
		err := rows.Scan(
			&t.ID, &t.Name, &t.Duration, &t.Status,
			&t.EarlyStart, &t.EarlyFinish, &t.LateStart, &t.LateFinish,
			&t.TotalFloat, &t.ReadyToStart,
		)
		if err != nil {
			http.Error(w, "Fout bij scannen kritieke taak: "+err.Error(), http.StatusInternalServerError)
			return
		}
		criticalTasks = append(criticalTasks, t)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(criticalTasks)
}
