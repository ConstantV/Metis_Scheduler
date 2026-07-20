-- Migration 001: Add support for standard tasks, process systems, system phase windows, and scope items

CREATE TABLE IF NOT EXISTS standard_tasks (
    id VARCHAR(50) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    duration INT NOT NULL,
    phase VARCHAR(50) NOT NULL -- PRE_TA, SHUTDOWN, MAINTENANCE, STARTUP
);

CREATE TABLE IF NOT EXISTS process_systems (
    id VARCHAR(50) PRIMARY KEY,
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS system_phase_windows (
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    system_id VARCHAR(50) NOT NULL REFERENCES process_systems(id) ON DELETE CASCADE,
    phase VARCHAR(50) NOT NULL, -- PRE_TA, SHUTDOWN, MAINTENANCE, STARTUP
    start_hour INT NOT NULL,
    end_hour INT NOT NULL,
    PRIMARY KEY (project_id, system_id, phase)
);

CREATE TABLE IF NOT EXISTS scope_items (
    id VARCHAR(50) PRIMARY KEY,
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    system_id VARCHAR(50) NOT NULL REFERENCES process_systems(id) ON DELETE CASCADE
);

-- Alter tasks to add scope_item_id and standard_task_id columns if they don't exist
ALTER TABLE tasks ADD COLUMN IF NOT EXISTS scope_item_id VARCHAR(50) REFERENCES scope_items(id) ON DELETE SET NULL;
ALTER TABLE tasks ADD COLUMN IF NOT EXISTS standard_task_id VARCHAR(50) REFERENCES standard_tasks(id) ON DELETE SET NULL;
