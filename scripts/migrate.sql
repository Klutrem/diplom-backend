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

CREATE TABLE IF NOT EXISTS telegram_alerts (
    id SERIAL PRIMARY KEY,
    bot_token VARCHAR(255) NOT NULL,
    chat_id VARCHAR(255) NOT NULL,
    thread_id INTEGER,
    alert_type VARCHAR(50) NOT NULL CHECK (alert_type IN ('all', 'normal', 'warning')),
    namespace VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(namespace, chat_id, thread_id)
);