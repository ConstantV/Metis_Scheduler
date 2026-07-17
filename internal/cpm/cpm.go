package cpm

type Activity struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Duration    int      `json:"duration"`
	DependsOn   []string `json:"depends_on"`
	EarlyStart  int      `json:"early_start"`
	EarlyFinish int      `json:"early_finish"`
	LateStart   int      `json:"late_start"`
	LateFinish  int      `json:"late_finish"`
	TotalFloat  int      `json:"total_float"`
	IsCritical  bool     `json:"is_critical"`

	// NIEUW: Status en de "Primavera-killer" vlag!
	Status       string `json:"status"`         // "NOT_STARTED", "IN_PROGRESS", "COMPLETED"
	ReadyToStart bool   `json:"ready_to_start"` // Onze nieuwe vlag
}

func CalculateSchedule(tasks []Activity) []Activity {
	if len(tasks) == 0 {
		return tasks
	}

	taskMap := make(map[string]*Activity)
	for i := range tasks {
		taskMap[tasks[i].ID] = &tasks[i]
	}

	// 1. Forward Pass
	projectDuration := 0
	for i := range tasks {
		calculateEarlyTimes(&tasks[i], taskMap)
		if tasks[i].EarlyFinish > projectDuration {
			projectDuration = tasks[i].EarlyFinish
		}
	}

	// 2. Backward Pass
	for i := range tasks {
		tasks[i].LateFinish = projectDuration
		tasks[i].LateStart = projectDuration - tasks[i].Duration
	}
	for i := range tasks {
		calculateLateTimes(&tasks[i], taskMap)
	}

	// 3. Float & Critical Path
	for i := range tasks {
		tasks[i].TotalFloat = tasks[i].LateStart - tasks[i].EarlyStart
		if tasks[i].TotalFloat == 0 {
			tasks[i].IsCritical = true
		}
	}

	// 4. NIEUW: BEREKEN "READY TO START"
	for i := range tasks {
		// Als de status leeg is, gaan we ervan uit dat hij nog niet gestart is
		if tasks[i].Status == "" {
			tasks[i].Status = "NOT_STARTED"
		}

		if tasks[i].Status != "NOT_STARTED" {
			tasks[i].ReadyToStart = false
			continue
		}
		// Als een taak al gestart of klaar is, hoeft hij niet meer te starten
		if tasks[i].Status != "NOT_STARTED" {
			tasks[i].ReadyToStart = false
			continue
		}

		allPredecessorsCompleted := true

		// Check alle voorgangers
		for _, depID := range tasks[i].DependsOn {
			depTask, exists := taskMap[depID]
			// Als de voorganger bestaat en NIET op COMPLETED staat, is deze taak geblokkeerd
			if exists && depTask.Status != "COMPLETED" {
				allPredecessorsCompleted = false
				break // We weten genoeg, stop deze sub-loop
			}
		}

		tasks[i].ReadyToStart = allPredecessorsCompleted
	}

	return tasks
}

// (De hulpfuncties calculateEarlyTimes en calculateLateTimes blijven ongewijzigd hieronder staan)
func calculateEarlyTimes(task *Activity, taskMap map[string]*Activity) {
	if len(task.DependsOn) == 0 {
		task.EarlyStart = 0
		task.EarlyFinish = task.Duration
		return
	}
	maxDependencyFinish := 0
	for _, depID := range task.DependsOn {
		depTask, exists := taskMap[depID]
		if exists {
			if depTask.EarlyFinish == 0 && depTask.Duration > 0 {
				calculateEarlyTimes(depTask, taskMap)
			}
			if depTask.EarlyFinish > maxDependencyFinish {
				maxDependencyFinish = depTask.EarlyFinish
			}
		}
	}
	task.EarlyStart = maxDependencyFinish
	task.EarlyFinish = task.EarlyStart + task.Duration
}

func calculateLateTimes(task *Activity, taskMap map[string]*Activity) {
	for _, depID := range task.DependsOn {
		depTask, exists := taskMap[depID]
		if exists {
			if task.LateStart < depTask.LateFinish {
				depTask.LateFinish = task.LateStart
				depTask.LateStart = depTask.LateFinish - depTask.Duration
			}
		}
	}
}
