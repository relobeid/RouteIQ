-- vehicles
CREATE TABLE IF NOT EXISTS vehicles (
    id UUID PRIMARY KEY,
    position_x INTEGER NOT NULL,
    position_y INTEGER NOT NULL,
    speed DECIMAL(5,2),
    destination_x INTEGER,
    destination_y INTEGER,
    last_updated TIMESTAMP DEFAULT NOW()
);

-- intersections
CREATE TABLE IF NOT EXISTS intersections (
    id SERIAL PRIMARY KEY,
    position_x INTEGER NOT NULL,
    position_y INTEGER NOT NULL,
    light_state VARCHAR(10) CHECK (light_state IN ('red', 'yellow', 'green')),
    last_changed TIMESTAMP DEFAULT NOW()
);

-- incidents
CREATE TABLE IF NOT EXISTS incidents (
    id UUID PRIMARY KEY,
    incident_type VARCHAR(20),
    position_x INTEGER NOT NULL,
    position_y INTEGER NOT NULL,
    severity INTEGER CHECK (severity BETWEEN 1 AND 5),
    created_at TIMESTAMP DEFAULT NOW(),
    resolved_at TIMESTAMP
);
