-- name: ListEnabledConfigs :many
SELECT id, name, pattern, version, enabled, priority, metadata_config, payload_config
FROM pathfinder_config
WHERE enabled = true
ORDER BY priority DESC;

-- name: GetConfigByName :one
SELECT id, name, pattern, version, enabled, priority, metadata_config, payload_config
FROM pathfinder_config
WHERE name = $1
LIMIT 1;

-- name: InsertConfig :one
INSERT INTO pathfinder_config (name, pattern, version, enabled, priority, metadata_config, payload_config)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id, name, pattern, version, enabled, priority, metadata_config, payload_config;
