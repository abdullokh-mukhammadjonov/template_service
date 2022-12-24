CREATE TABLE IF NOT EXISTS uzkad_action_logs
(
    entity_number    VARCHAR,
    request_url    VARCHAR NOT NULL,
    status          VARCHAR(100),
    created_at       TIMESTAMP DEFAULT now(),
    response   VARCHAR,
    integration VARCHAR(100)
);
