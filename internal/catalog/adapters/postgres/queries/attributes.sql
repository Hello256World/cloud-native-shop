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

-- name: GetAttributeDefinitionByCode :one
SELECT * FROM catalog.attribute_definitions WHERE code = $1;

-- name: UpdateAttributeDefinition :one
UPDATE catalog.attribute_definitions
SET name = $2, unit = $3, is_filterable = $4, display_order = $5
WHERE id = $1
RETURNING *;

-- name: ListAttributeFacets :many
SELECT ad.id AS attribute_id, ad.code, ad.name, ad.value_type, ad.unit, ad.display_order,
       pa.value_text, pa.value_number, COUNT(*)::bigint AS product_count
FROM catalog.product_attributes pa
JOIN catalog.attribute_definitions ad ON ad.id = pa.attribute_id AND ad.is_filterable
JOIN catalog.products p ON p.id = pa.product_id
JOIN catalog.categories c ON c.id = p.category_id AND c.is_active AND c.deleted_at IS NULL
WHERE p.deleted_at IS NULL
  AND p.is_active
  AND p.min_price IS NOT NULL
  AND (sqlc.arg(category_id)::bigint = 0 OR p.category_id = sqlc.arg(category_id)::bigint)
GROUP BY ad.id, ad.code, ad.name, ad.value_type, ad.unit, ad.display_order,
         pa.value_text, pa.value_number
ORDER BY ad.display_order, ad.id, pa.value_number, pa.value_text;