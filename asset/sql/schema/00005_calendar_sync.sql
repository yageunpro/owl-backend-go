-- +goose Up
-- +goose StatementBegin
CREATE TABLE calendar.sync
(
    id         uuid                      NOT NULL
        CONSTRAINT "PK_calendar_sync_id"
            PRIMARY KEY,
    user_id    uuid                      NOT NULL,
    sync_token VARCHAR(200)              NOT NULL,
    created_at timestamptz DEFAULT NOW() NOT NULL,
    updated_at timestamptz DEFAULT NOW() NOT NULL,
    deleted_at timestamptz DEFAULT NULL
);

ALTER TABLE calendar.schedule
    ADD google_calc_id VARCHAR(200)
        CONSTRAINT "UQ_google_calc_id" UNIQUE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE calendar.schedule
    DROP google_calc_id;
DROP TABLE calendar.sync;
-- +goose StatementEnd
