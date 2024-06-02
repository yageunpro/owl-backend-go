-- +goose Up
-- +goose StatementBegin
CREATE SCHEMA IF NOT EXISTS location;

CREATE TABLE location.query
(
    id         uuid                      NOT NULL
        CONSTRAINT "PK_location_query_id"
            PRIMARY KEY,
    data       VARCHAR(512)              NOT NULL,
    created_at timestamptz DEFAULT NOW() NOT NULL,
    updated_at timestamptz DEFAULT NOW() NOT NULL,
    deleted_at timestamptz DEFAULT NULL
);

CREATE TABLE location.location
(
    id         uuid         NOT NULL
        CONSTRAINT "PK_location_location_id"
            PRIMARY KEY,
    query_id   uuid         NOT NULL
        CONSTRAINT "FK_location_query_id"
            REFERENCES location.query,
    title      VARCHAR(100) NOT NULL,
    address    VARCHAR(512) NOT NULL,
    category   VARCHAR(100) NOT NULL,
    map_x      INT          NOT NULL,
    map_y      INT          NOT NULL,
    deleted_at timestamptz  NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE location.query CASCADE;
DROP TABLE location.location CASCADE;
DROP SCHEMA location CASCADE;
-- +goose StatementEnd
