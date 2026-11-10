-- name: InsertAttributeDefinition :one
INSERT INTO catalog.attribute_definitions (code, name, value_type, unit, is_filterable, display_order)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: ListAttributeDefinitions :many
SELECT * FROM catalog.attribute_definitions ORDER BY display_order, id;

-- name: ListAttributeDefinitionsByIDs :many
SELECT * FROM catalog.attribute_definitions WHERE id = ANY(sqlc.arg(ids)::bigint[]);

-- name: DeleteProductAttributes :exec
DELETE FROM catalog.product_attributes WHERE product_id = $1;

-- name: InsertProductAttribute :exec
INSERT INTO catalog.product_attributes (product_id, attribute_id, value_text, value_number)
VALUES ($1, $2, $3, $4);

-- name: ListProductAttributes :many
SELECT pa.attribute_id, ad.code, ad.name, ad.value_type, ad.unit, ad.display_order,
       pa.value_text, pa.value_number
FROM catalog.product_attributes pa
JOIN catalog.attribute_definitions ad ON ad.id = pa.attribute_id
WHERE pa.product_id = $1
ORDER BY ad.display_order, ad.id;