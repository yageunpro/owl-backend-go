-- +goose Up
-- +goose StatementBegin
CREATE SCHEMA IF NOT EXISTS appointment;

CREATE TYPE appointment.status AS ENUM ('DRAFT', 'CONFIRM', 'DONE', 'CANCEL', 'DELETE');

CREATE TABLE appointment.appointment
(
    id           uuid                               NOT NULL
        CONSTRAINT "PK_appointment_appointment_id"
            PRIMARY KEY,
    organizer_id uuid                               NOT NULL,
    title        VARCHAR(100)                       NOT NULL,
    description  VARCHAR(200)                       NOT NULL,
    category     VARCHAR(100)[]                     NOT NULL,
    status       appointment.status DEFAULT 'DRAFT' NOT NULL,
    location_id  uuid                               NULL,
    deadline     timestamptz                        NOT NULL,
    confirm_time timestamptz                        NULL,
    created_at   timestamptz        DEFAULT NOW()   NOT NULL,
    updated_at   timestamptz        DEFAULT NOW()   NOT NULL,
    deleted_at   timestamptz        DEFAULT NULL
);

CREATE TABLE appointment.participant
(
    id             uuid                      NOT NULL
        CONSTRAINT "PK_appointment_participant_id"
            PRIMARY KEY,
    appointment_id uuid                      NOT NULL
        CONSTRAINT "FK_appointment_participant_appointment_id"
            REFERENCES appointment.appointment,
    user_id        uuid                      NOT NULL,
    created_at     timestamptz DEFAULT NOW() NOT NULL,
    updated_at     timestamptz DEFAULT NOW() NOT NULL,
    deleted_at     timestamptz DEFAULT NULL,
    CONSTRAINT "UQ_appointment_id_user_id"
        UNIQUE (appointment_id, user_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE appointment.participant CASCADE;
DROP TABLE appointment.appointment CASCADE;
DROP TYPE appointment.status CASCADE;
DROP SCHEMA appointment CASCADE;
-- +goose StatementEnd
