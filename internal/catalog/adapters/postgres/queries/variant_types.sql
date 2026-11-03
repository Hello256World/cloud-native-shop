-- name: InsertVariantType :one
INSERT INTO catalog.variant_types (name, data_type, display_order)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetVariantTypeByID :one
SELECT * FROM catalog.variant_types WHERE id = $1;

-- name: ListVariantTypes :many
SELECT * FROM catalog.variant_types ORDER BY display_order, id;

-- name: InsertVariantValue :one
INSERT INTO catalog.variant_values (variant_type_id, value, display_order)
VALUES ($1, $2, $3)
RETURNING *;

-- name: ListVariantValues :many
SELECT * FROM catalog.variant_values ORDER BY display_order, id;

-- name: ListVariantValuesByIDs :many
SELECT * FROM catalog.variant_values WHERE id = ANY($1::bigint[]);