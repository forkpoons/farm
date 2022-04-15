-- +goose Up
-- +goose StatementBegin
CREATE TABLE temperatures
(
    id   int AUTO_INCREMENT NOT NULL,
    date DATETIME,
    temp float,
    PRIMARY KEY (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
Drop TABLE temperatures;
-- +goose StatementEnd
