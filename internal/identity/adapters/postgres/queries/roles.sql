-- name: ListRoles :many
SELECT id, name FROM identity.roles ORDER BY id;