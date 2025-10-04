-- name: ListUsers :many
SELECT * FROM users ORDER BY id;

-- name: CreateUserClassic :one
WITH new_user AS (
	INSERT INTO users DEFAULT VALUES
	RETURNING *
), new_credentials AS (
	INSERT 
	INTO user_credentials (user_id, password_hash)
	SELECT
		id, crypt($1, gen_salt('bf'))
	FROM
		new_user
), new_contact AS (
	INSERT
	INTO
		contact_methods (user_id, type, value)
	SELECT
		id, $2, $3
	FROM
		new_user
)
SELECT
	u.*
FROM
	new_user u;

-- name: CreateUserOAuth :one
WITH new_user AS (
	INSERT INTO users DEFAULT VALUES
	RETURNING *
), new_oauth AS (
	INSERT INTO oauth_accounts (user_id, provider, provider_user_id)
	SELECT
		id, $1, $2
	FROM
		new_user
)
SELECT
	u.*
FROM
	new_user u;

-- name: FindUserByContactValue :one
SELECT
	u.*
FROM
	users u
	JOIN contact_methods cm ON u.id = cm.user_id
WHERE
	cm.type = $1
	AND cm.value = $2;

-- name: FindUserByContactMethod :one
SELECT
	u.*,
	cm.*
FROM
	users u
	JOIN contact_methods cm ON u.id = cm.user_id
	JOIN user_credentials uc ON u.id = uc.user_id
WHERE
	cm.type = $1
	AND cm.value = $2
	AND uc.password_hash = crypt($3, uc.password_hash);

-- name: VerifyContactMethod :one
UPDATE
	contact_methods
SET
	is_verified = TRUE
WHERE
	user_id = $1
	AND type = $2
	AND value = $3
RETURNING *;

-- name: FindUserByOAuth :one
SELECT
	u.*
FROM
	users u
	JOIN oauth_accounts oa ON u.id = oa.user_id
WHERE
	oa.provider = $1
	AND oa.provider_user_id = $2;

-- name: UpdateUserLoginTime :one
UPDATE
	users
SET
	last_login_at = now()
WHERE
	id = $1
RETURNING *;
