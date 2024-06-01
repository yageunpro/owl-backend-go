-- +goose NO TRANSACTION
-- +goose Up

CREATE SCHEMA IF NOT EXISTS calendar;

CREATE TYPE calendar.appointment_status AS ENUM ('DRAFT', 'CONFIRM', 'DONE', 'CANCEL');

CREATE TABLE calendar.appointment
(
    id           uuid                         NOT NULL
        CONSTRAINT "PK_calendar_appointment_id"
            PRIMARY KEY,
    title        VARCHAR(100)                 NOT NULL,
    description  VARCHAR(200)                 NOT NULL,
    status       calendar.appointment_status  NOT NULL,
    keywords     VARCHAR(100)[] DEFAULT '{}'  NOT NULL,
    start_time   timestamptz                  NULL,
    end_time     timestamptz                  NULL,
    confirm_time timestamptz                  NULL,
    created_at   timestamptz    DEFAULT NOW() NOT NULL,
    updated_at   timestamptz    DEFAULT NOW() NOT NULL,
    deleted_at   timestamptz    DEFAULT NULL
);

CREATE TABLE calendar.location
(
    id         uuid                      NOT NULL
        CONSTRAINT "PK_calendar_location_id"
            PRIMARY KEY,
    title      VARCHAR(100)              NOT NULL,
    address    VARCHAR(400)              NOT NULL,
    category   VARCHAR(100)              NOT NULL,
    pos_x      INT                       NOT NULL,
    pos_y      INT                       NOT NULL,
    created_at timestamptz DEFAULT NOW() NOT NULL,
    updated_at timestamptz DEFAULT NOW() NOT NULL,
    deleted_at timestamptz DEFAULT NULL
);

CREATE TABLE calendar.appointment_location_map
(
    appointment_id uuid NOT NULL
        CONSTRAINT "PK_calendar_map_appointment_id"
            PRIMARY KEY
        CONSTRAINT "FK_calendar_map_appointment_id"
            REFERENCES calendar.appointment,
    location_id    uuid NOT NULL
        CONSTRAINT "FK_calendar_map_location_id"
            REFERENCES calendar.location
);

CREATE TABLE calendar.schedule
(
    id         uuid                      NOT NULL
        CONSTRAINT "PK_calendar_schedule_id"
            PRIMARY KEY,
    title      VARCHAR(100)              NOT NULL,
    start_time timestamptz               NOT NULL,
    end_time   timestamptz               NOT NULL,
    created_at timestamptz DEFAULT NOW() NOT NULL,
    updated_at timestamptz DEFAULT NOW() NOT NULL,
    deleted_at timestamptz DEFAULT NULL
);
