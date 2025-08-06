-- +goose Up
ALTER TABLE users
DROP COLUMN IF EXISTS hashed_password;
ALTER TABLE users
ADD COLUMN hashed_password TEXT NOT NULL DEFAULT 'unset';
