-- DROP TABLE IF EXISTS events;

CREATE TABLE IF NOT EXISTS events (
    id VARCHAR(255) PRIMARY KEY,
    namespace VARCHAR(255),
    name VARCHAR(255),
    reason VARCHAR(255),
    message TEXT,
    type VARCHAR(50),
    involved_object VARCHAR(255),
    first_timestamp TIMESTAMP,
    last_timestamp TIMESTAMP,
    count INTEGER
);

CREATE TABLE IF NOT EXISTS watched_namespaces (
    namespace VARCHAR(255) PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);