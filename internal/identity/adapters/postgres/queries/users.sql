-- name: InsertUser :one
INSERT INTO identity.users (full_name, email, phone, password_hash, birthday, gender, is_active)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetUserByID :one
SELECT * FROM identity.users
WHERE id = $1 AND deleted_at IS NULL;

-- name: GetUserByPhone :one
SELECT * FROM identity.users
WHERE phone = $1 AND deleted_at IS NULL;

-- name: AddUserRoleByName :exec
INSERT INTO identity.user_roles (user_id, role_id)
SELECT $1, id FROM identity.roles WHERE name = $2;

-- name: ListRoleNamesByUser :many
SELECT r.name
FROM identity.user_roles ur
JOIN identity.roles r ON r.id = ur.role_id
WHERE ur.user_id = $1
ORDER BY r.id;