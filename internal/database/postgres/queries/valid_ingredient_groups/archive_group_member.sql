-- name: ArchiveValidIngredientGroupMember :exec

UPDATE valid_ingredient_group_members SET archived_at = NOW() WHERE id = $1 AND belongs_to_group = $2;