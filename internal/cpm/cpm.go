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

	// Status en de "Primavera-killer" vlag!
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
		// Reset kritieke vlag voor een schone berekening
		tasks[i].IsCritical = false
	}

	// 1. Forward Pass (Bereken vroege tijden)
	projectDuration := 0
	for i := range tasks {
		calculateEarlyTimes(&tasks[i], taskMap)
		if tasks[i].EarlyFinish > projectDuration {
			projectDuration = tasks[i].EarlyFinish
		}
	}

	// 2. Backward Pass (Nu volledig recursief via opvolgers!)
	// We bouwen een kaart: welke taken volgen er op deze taak?
	successors := make(map[string][]string)
	for _, t := range tasks {
		for _, depID := range t.DependsOn {
			successors[depID] = append(successors[depID], t.ID)
		}
	}

	lateComputed := make(map[string]bool)
	var computeLate func(id string)

	computeLate = func(id string) {
		// Als deze taak al berekend is in de recursie, sla hem over
		if lateComputed[id] {
			return
		}
		t := taskMap[id]
		succIDs := successors[id]

		if len(succIDs) == 0 {
			// Als er geen opvolgers zijn, is het een eindtaak.
			// LateFinish is dan gelijk aan de totale projectduur.
			t.LateFinish = projectDuration
		} else {
			// Als er wel opvolgers zijn, is de LateFinish de ALLERKLEINSTE LateStart van zijn opvolgers!
			minLateStart := projectDuration
			for _, succID := range succIDs {
				computeLate(succID) // Zorg dat de opvolger éérst berekend is!
				if taskMap[succID].LateStart < minLateStart {
					minLateStart = taskMap[succID].LateStart
				}
			}
			t.LateFinish = minLateStart
		}
		t.LateStart = t.LateFinish - t.Duration
		lateComputed[id] = true
	}

	// Voer de backward pass uit voor alle taken
	for i := range tasks {
		computeLate(tasks[i].ID)
	}

	// 3. Float & Critical Path
	for i := range tasks {
		tasks[i].TotalFloat = tasks[i].LateStart - tasks[i].EarlyStart
		if tasks[i].TotalFloat == 0 {
			tasks[i].IsCritical = true
		}
	}

	// 4. BEREKEN "READY TO START"
	for i := range tasks {
		if tasks[i].Status == "" {
			tasks[i].Status = "NOT_STARTED"
		}

		if tasks[i].Status != "NOT_STARTED" {
			tasks[i].ReadyToStart = false
			continue
		}

		allPredecessorsCompleted := true

		// Check alle voorgangers
		for _, depID := range tasks[i].DependsOn {
			depTask, exists := taskMap[depID]
			if exists && depTask.Status != "COMPLETED" {
				allPredecessorsCompleted = false
				break
			}
		}

		tasks[i].ReadyToStart = allPredecessorsCompleted
	}

	return tasks
}

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
