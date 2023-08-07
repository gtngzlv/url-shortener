-- +goose Up
CREATE TABLE IF NOT EXISTS URL_STORAGE(
    ID SERIAL NOT NULL PRIMARY KEY,
    short TEXT NOT NULL,
    long TEXT NOT NULL,
    userID TEXT NOT NULL,
    is_deleted BIT DEFAULT 0::bit NOT NULL);
CREATE UNIQUE INDEX long_id ON URL_STORAGE USING btree(long);

-- +goose Down
-- +goose StatementBegin
DROP TABLE URL_STORAGE;
-- +goose StatementEnd
