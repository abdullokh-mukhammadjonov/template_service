CREATE TABLE IF NOT EXISTS file_track
(
    id uuid primary key,
    entity_id VARCHAR NOT NULL,
    property_id VARCHAR NOT NULL,
    to_delete boolean NOT NULL,
    file_name VARCHAR NOT NULL,
    bucket_name VARCHAR NOT NULL,
    created_at TIMESTAMP DEFAULT now()
)