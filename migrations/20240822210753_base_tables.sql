-- +goose Up
-- +goose StatementBegin
CREATE TYPE root_type AS ENUM ('LOCAL', 'MINIO');
CREATE TABLE roots
(
    id         UUID PRIMARY KEY,
    type       root_type NOT NULL,
    name       TEXT      NOT NULL,
    config     JSON      NOT NULL,
    scanned_at TIMESTAMP
);

CREATE TABLE files
(
    id           UUID PRIMARY KEY,
    root_id      UUID     NOT NULL REFERENCES roots (id),
    path         TEXT     NOT NULL,
    name         TEXT     NOT NULL,
    content_type TEXT     NOT NULL,
    size         BIGINT   NOT NULL,
    metadata     JSON,
    md5          CHAR(32) NOT NULL
);

CREATE TYPE action_type AS ENUM ('CREATE', 'UPDATE', 'DELETE');
CREATE TABLE file_history
(
    id           UUID PRIMARY KEY,
    file_id      UUID        NOT NULL,
    action_type  action_type NOT NULL,
    performed_at TIMESTAMP   NOT NULL
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
