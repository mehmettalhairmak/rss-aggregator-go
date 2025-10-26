-- +goose Up

ALTER TABLE feeds ADD COLUMN description TEXT;
ALTER TABLE feeds ADD COLUMN logo_url TEXT;
ALTER TABLE feeds ADD COLUMN priority INTEGER NOT NULL DEFAULT 3;

-- +goose Down

ALTER TABLE feeds DROP COLUMN priority;
ALTER TABLE feeds DROP COLUMN logo_url;
ALTER TABLE feeds DROP COLUMN description;
