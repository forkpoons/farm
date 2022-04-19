-- +goose Up
-- +goose StatementBegin
CREATE TABLE actions
(
    id     int AUTO_INCREMENT NOT NULL,
    name VARCHAR(30),
    date   DATETIME,
    action int,
    PRIMARY KEY (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
Drop TABLE actions;
-- +goose StatementEnd
