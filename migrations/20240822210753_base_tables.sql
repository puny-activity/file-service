-- +goose Up
-- +goose StatementBegin
CREATE TYPE root_type AS ENUM ('LOCAL', 'MINIO');
CREATE TABLE roots
(
    id         UUID PRIMARY KEY,
    type       root_type NOT NULL,
    name       TEXT      NOT NULL,
    config     JSON      NOT NULL,
    scanned_at TIMESTAMPTZ
);

CREATE TABLE files
(
    id           UUID PRIMARY KEY,
    root_id      UUID     NOT NULL REFERENCES roots (id),
    path         TEXT     NOT NULL,
    name         TEXT     NOT NULL,
    content_type TEXT     NOT NULL,
    size         BIGINT   NOT NULL,
    metadata     JSON     NOT NULL,
    md5          CHAR(32) NOT NULL
);

CREATE TYPE action_type AS ENUM ('CREATE', 'UPDATE', 'DELETE');
CREATE TABLE file_history
(
    id           UUID PRIMARY KEY,
    file_id      UUID        NOT NULL,
    type         action_type NOT NULL,
    performed_at TIMESTAMPTZ NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE file_history;
DROP TYPE action_type;

DROP TABLE files;

DROP TABLE roots;
DROP TYPE root_type;
-- +goose StatementEnd
