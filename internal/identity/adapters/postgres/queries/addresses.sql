-- name: InsertAddress :one
INSERT INTO identity.addresses (user_id, receiver_name, phone, address, postal_code, no, unit)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: ListAddressesByUser :many
SELECT * FROM identity.addresses
WHERE user_id = $1 AND deleted_at IS NULL
ORDER BY id;

-- name: UpdateAddress :one
UPDATE identity.addresses
SET receiver_name = $3, phone = $4, address = $5, postal_code = $6, no = $7, unit = $8
WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL
RETURNING *;

-- name: DeleteAddress :execrows
UPDATE identity.addresses SET deleted_at = now()
WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL;