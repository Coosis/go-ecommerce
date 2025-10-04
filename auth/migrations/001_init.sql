-- +goose Up
CREATE TABLE users (
	id            UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
	created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
	last_login_at TIMESTAMPTZ
);

CREATE TABLE user_credentials (
	user_id       UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
	password_hash TEXT NOT NULL
);

CREATE TABLE contact_methods (
	id                 UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
	user_id            UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	type               VARCHAR(31) NOT NULL CHECK (type IN ('email','phone')),
	value              TEXT        NOT NULL UNIQUE,
	is_verified        BOOLEAN     NOT NULL DEFAULT FALSE,
	created_at         TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE oauth_accounts (
	id                   UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
	user_id              UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	provider             VARCHAR(31) NOT NULL CHECK (provider IN ('github' , 'google', 'facebook')),
	provider_user_id     TEXT        NOT NULL,
	-- access_token         TEXT, -- potential user info like avatar
	-- refresh_token        TEXT,
	-- token_expires_at     TIMESTAMPTZ,
	created_at           TIMESTAMPTZ NOT NULL DEFAULT now(),
	UNIQUE (provider, provider_user_id)
);

CREATE INDEX ON oauth_accounts(user_id);

CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- +goose Down
DROP TABLE IF EXISTS oauth_accounts;
DROP TABLE IF EXISTS contact_methods;
DROP TABLE IF EXISTS user_credentials;
