create extension if not exists "uuid-ossp";

CREATE TABLE IF NOT EXISTS action_logs
(
    entity_id    VARCHAR,
    action_id    VARCHAR NOT NULL,
    status       VARCHAR(100),
    indicator VARCHAR(100),
    action_type VARCHAR,
    error_string VARCHAR,
    updated_at TIMESTAMP,
    entity_number VARCHAR,
    req_body VARCHAR,
    order_number INTEGER,
    id serial INTEGER
);

CREATE TABLE IF NOT EXISTS qaror_action_logs
(
    entity_number    VARCHAR,
    request_url    VARCHAR NOT NULL,
    status          VARCHAR(100),
    created_at       TIMESTAMP DEFAULT now(),
    response   VARCHAR,
    integration VARCHAR(100)
);

CREATE TABLE IF NOT EXISTS orders
(
    "entity_id" VARCHAR,
    "order_id" VARCHAR,
    "law_type" VARCHAR(12),
    "language" VARCHAR,
    "soato" INTEGER,
    "name" VARCHAR,
    "price" DOUBLE PRECISION,
    "address" VARCHAR,
    "category" VARCHAR,
    "closed" INTEGER,
    "prop_set" INTEGER,
    "additional_info" VARCHAR,
    "lat" DOUBLE PRECISION,
    "lng" DOUBLE PRECISION,
    "created_at" TIMESTAMP DEFAULT now()
);

CREATE TABLE IF NOT EXISTS order_documents
(
    "order_id" VARCHAR,
    "document_id" INTEGER,
    "language" VARCHAR,
    "show_on_front" VARCHAR,
    "description" VARCHAR,
    "file" VARCHAR,
    "type" VARCHAR,
    "created_at" TIMESTAMP DEFAULT now()
);

CREATE TABLE IF NOT EXISTS order_images
(
    "order_id" INTEGER,
    "image_id" INTEGER,
    "language" VARCHAR,
    "is_main" VARCHAR,
    "image_position" VARCHAR,
    "description" VARCHAR,
    "file" VARCHAR,
    "created_at" TIMESTAMP DEFAULT now()
);

CCREATE TABLE IF NOT EXISTS function_error_logs (
	id serial primary key,
	target VARCHAR,
  	target_id VARCHAR,
    other_target VARCHAR,
    endpoint_url VARCHAR,
    http_method VARCHAR,
    handler_function VARCHAR,
    call_stack VARCHAR,
    req_body VARCHAR,
    res_body VARCHAR,
    indicator VARCHAR,
    collection VARCHAR,
    service VARCHAR,
    created_at timestamp
);

CREATE TABLE IF NOT EXISTS function_info_logs (
	id serial primary key,
	target VARCHAR,
  	target_id VARCHAR,
    other_target VARCHAR,
    endpoint_url VARCHAR,
    http_method VARCHAR,
    handler_function VARCHAR,
    call_stack VARCHAR,
    req_body VARCHAR,
    res_body VARCHAR,
    indicator VARCHAR,
    collection VARCHAR,
    service VARCHAR,
    created_at timestamp
);

CREATE TABLE IF NOT EXISTS push_from_auction (
	id serial primary key,
	order_id INTEGER,
  	category_id INTEGER,
    protocol INTEGER,
    url VARCHAR,
    is_finished BOOLEAN,
    success VARCHAR,
    req_body VARCHAR,
    error VARCHAR,
    entity_id VARCHAR,
    entity_number VARCHAR,
    old_status VARCHAR,
    new_status VARCHAR,
    indicator VARCHAR,
    auction_status_code INTEGER,
    auction_status_name VARCHAR,
    created_at timestamp
);

CREATE TABLE IF NOT EXISTS cronjob_logs (
	id serial primary key,
    entity_id VARCHAR,
    current_status_id VARCHAR,
    next_status_id VARCHAR,
    delayed_days INTEGER,
    organizations VARCHAR,
    update_parallel_action_request VARCHAR,
    setmap VARCHAR,
    action_history_id VARCHAR,
    created_at timestamp
);

CREATE TABLE IF NOT EXISTS cronjob_error_logs (
	id serial primary key,
    error VARCHAR,
    indicator VARCHAR,
    additional VARCHAR,
    created_at timestamp
);