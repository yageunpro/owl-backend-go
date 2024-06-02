-- +goose Up
-- +goose StatementBegin
CREATE SCHEMA IF NOT EXISTS calendar;

CREATE TABLE calendar.schedule
(
    id         uuid                      NOT NULL
        CONSTRAINT "PK_calendar_schedule_id"
            PRIMARY KEY,
    user_id    uuid                      NOT NULL,
    title      VARCHAR(100)              NOT NULL,
    period     tstzrange                 NOT NULL,
    created_at timestamptz DEFAULT NOW() NOT NULL,
    updated_at timestamptz DEFAULT NOW() NOT NULL,
    deleted_at timestamptz DEFAULT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE calendar.schedule CASCADE;
DROP SCHEMA calendar CASCADE;
-- +goose StatementEnd
