-- +goose Up
-- +goose StatementBegin

CREATE SCHEMA IF NOT EXISTS auth;

CREATE TABLE auth."user"
(
    id         uuid                      NOT NULL
        CONSTRAINT "PK_auth_user_id"
            PRIMARY KEY,
    email      VARCHAR(320)              NOT NULL
        CONSTRAINT "UQ_auth_user_email"
            UNIQUE,
    username   VARCHAR(100)              NOT NULL,
    created_at timestamptz DEFAULT NOW() NOT NULL,
    updated_at timestamptz DEFAULT NOW() NOT NULL,
    deleted_at timestamptz DEFAULT NULL
);

CREATE TABLE auth.password
(
    id            uuid                      NOT NULL
        CONSTRAINT "PK_auth_password_id"
            PRIMARY KEY
        CONSTRAINT "FK_auth_password_id_user_id"
            REFERENCES auth."user",
    password_hash VARCHAR(100)              NOT NULL,
    created_at    timestamptz DEFAULT NOW() NOT NULL,
    updated_at    timestamptz DEFAULT NOW() NOT NULL,
    deleted_at    timestamptz DEFAULT NULL
);

CREATE TABLE auth.oauth
(
    id            uuid                      NOT NULL
        CONSTRAINT "PK_auth_oauth_id"
            PRIMARY KEY
        CONSTRAINT "FK_auth_oauth_id_user_id"
            REFERENCES auth."user",
    open_id       VARCHAR(100)              NOT NULL,
    access_token  VARCHAR(512)              NOT NULL,
    refresh_token VARCHAR(512)              NULL,
    allow_sync    BOOLEAN                   NOT NULL,
    valid_until   timestamptz               NOT NULL,
    created_at    timestamptz DEFAULT NOW() NOT NULL,
    updated_at    timestamptz DEFAULT NOW() NOT NULL,
    deleted_at    timestamptz DEFAULT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE auth.oauth CASCADE;
DROP TABLE auth.password CASCADE;
DROP TABLE auth."user" CASCADE;
DROP SCHEMA auth CASCADE;
-- +goose StatementEnd
