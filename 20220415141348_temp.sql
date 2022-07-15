-- +goose Up
-- +goose StatementBegin
CREATE TABLE temperatures
(
    id   int AUTO_INCREMENT NOT NULL,
    type VARCHAR(30),
    name VARCHAR(30),
    date DATETIME,
    temp float,
    PRIMARY KEY (id)
);
-- +goose StatementEnd
-- +goose StatementBegin
CREATE TABLE humidity
(
    id   int AUTO_INCREMENT NOT NULL,
    type VARCHAR(30),
    name VARCHAR(30),
    date DATETIME,
    temp float,
    PRIMARY KEY (id)
);
-- +goose StatementEnd
-- +goose StatementBegin
CREATE TABLE actions
(
    id     int AUTO_INCREMENT NOT NULL,
    name   VARCHAR(30),
    date   DATETIME,
    action int,
    PRIMARY KEY (id)
);
-- +goose StatementEnd
-- +goose StatementBegin
CREATE TABLE greenhouse
(
    id         int AUTO_INCREMENT NOT NULL,
    name       VARCHAR(30),
    lastonline DATETIME,
    lastdata   DATETIME,
    temp       float,
    humidity   float,
    online     bool,
    work       bool,
    PRIMARY KEY (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
Drop TABLE temperatures;
-- +goose StatementEnd
-- +goose StatementBegin
Drop TABLE actions;
-- +goose StatementEnd
-- +goose StatementBegin
Drop TABLE greenhouse;
-- +goose StatementEnd
-- +goose StatementBegin
Drop TABLE humidity;
-- +goose StatementEnd
