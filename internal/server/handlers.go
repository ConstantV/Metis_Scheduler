package server

import (
	"encoding/json"
	"net/http"

	"metis-scheduler/internal/cpm"
)

type Server struct{}

func NewServer() *Server {
	return &Server{}
}

// Zorg ervoor dat hier de 'Status' velden exact zo bij staan!
func getMockProject() []cpm.Activity {
	return []cpm.Activity{
		// Hoofd-startpunt (Al voltooid)
		{ID: "task-1", Name: "Ontwerp & Vergunning", Duration: 5, Status: "COMPLETED", DependsOn: []string{}},

		// Fundering is nu in uitvoering (Niet voltooid!)
		{ID: "task-2", Name: "Grondwerk & Fundering", Duration: 4, Status: "IN_PROGRESS", DependsOn: []string{"task-1"}},

		// Parallel Pad A: Wacht op fundering
		{ID: "task-3", Name: "Ruwbouw Kring A", Duration: 6, Status: "NOT_STARTED", DependsOn: []string{"task-2"}},

		// Parallel Pad B: Wacht op fundering (Duurt langer, dus kritiek!)
		{ID: "task-4", Name: "Ruwbouw Kring B", Duration: 8, Status: "NOT_STARTED", DependsOn: []string{"task-2"}},

		// Parallel Pad C: Onafhankelijk van de fundering. De aanvraag is al gedaan!
		{ID: "task-5", Name: "Nutsnet Aanvragen", Duration: 2, Status: "COMPLETED", DependsOn: []string{"task-1"}},

		// Deze taak wacht op de aanvraag (task-5). Omdat task-5 COMPLETED is, MOET deze vlag op TRUE springen!
		{ID: "task-6", Name: "Nutsleidingen Leggen", Duration: 3, Status: "NOT_STARTED", DependsOn: []string{"task-5"}},

		// Merge Point 1: Wacht tot BEIDE ruwbouw-kringen (task-3 en task-4) klaar zijn
		{ID: "task-7", Name: "Afwerking & Interieur", Duration: 5, Status: "NOT_STARTED", DependsOn: []string{"task-3", "task-4"}},

		// Eindstation: Wacht op de afwerking en de nutsleidingen
		{ID: "task-8", Name: "Oplevering & Inspectie", Duration: 2, Status: "NOT_STARTED", DependsOn: []string{"task-7", "task-6"}},
	}
}

// HandleSchedule berekent de volledige planning
func (s *Server) HandleSchedule(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	tasks := getMockProject()
	calculatedSchedule := cpm.CalculateSchedule(tasks)

	json.NewEncoder(w).Encode(calculatedSchedule)
}

// HandleCriticalPath filtert alleen de taken die kritiek zijn eruit!
func (s *Server) HandleCriticalPath(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	tasks := getMockProject()
	calculatedSchedule := cpm.CalculateSchedule(tasks)

	var criticalTasks []cpm.Activity
	for _, task := range calculatedSchedule {
		if task.IsCritical {
			criticalTasks = append(criticalTasks, task)
		}
	}

	json.NewEncoder(w).Encode(criticalTasks)
}
