-- name: CreateAuditLog :exec

INSERT INTO audit_log (
	id,
	resource_type,
	relevant_id,
	event_type,
	changes,
	belongs_to_user,
	belongs_to_household
) VALUES (
	sqlc.arg(id),
	sqlc.arg(resource_type),
	sqlc.arg(relevant_id),
	sqlc.arg(event_type),
	sqlc.arg(changes),
	sqlc.arg(belongs_to_user),
	sqlc.arg(belongs_to_household)
);

-- name: GetAuditLog :one

SELECT
	audit_log.id as audit_log_id,
	audit_log.resource_type as audit_log_resource_type,
	audit_log.relevant_id as audit_log_relevant_id,
	audit_log.event_type as audit_log_event_type,
	audit_log.changes as audit_log_changes,
	audit_log.belongs_to_user as audit_log_belongs_to_user,
	audit_log.belongs_to_household as audit_log_belongs_to_household,
	audit_log.created_at as audit_log_created_at
FROM audit_log
WHERE audit_log.id = sqlc.arg(id);

-- name: GetAuditLogsForUser :many

SELECT
	audit_log.id,
	audit_log.resource_type,
	audit_log.relevant_id,
	audit_log.event_type,
	audit_log.changes,
	audit_log.belongs_to_user,
	audit_log.belongs_to_household,
	audit_log.created_at,
	(
		SELECT COUNT(audit_log.id)
		FROM audit_log
		WHERE audit_log.created_at > COALESCE(sqlc.narg(created_after), (SELECT NOW() - '999 years'::INTERVAL))
			AND audit_log.created_at < COALESCE(sqlc.narg(created_before), (SELECT NOW() + '999 years'::INTERVAL))
			AND audit_log.belongs_to_user = sqlc.arg(belongs_to_user)
	) AS filtered_count,
	(
		SELECT COUNT(audit_log.id)
		FROM audit_log
		WHERE
			audit_log.belongs_to_user = sqlc.arg(belongs_to_user)
	) AS total_count
FROM audit_log
WHERE audit_log.created_at > COALESCE(sqlc.narg(created_after), (SELECT NOW() - '999 years'::INTERVAL))
	AND audit_log.created_at < COALESCE(sqlc.narg(created_before), (SELECT NOW() + '999 years'::INTERVAL))
	AND audit_log.belongs_to_user = sqlc.arg(belongs_to_user)
LIMIT sqlc.narg(query_limit)
OFFSET sqlc.narg(query_offset);

-- name: GetAuditLogsForUserAndResourceType :many

SELECT
	audit_log.id,
	audit_log.resource_type,
	audit_log.relevant_id,
	audit_log.event_type,
	audit_log.changes,
	audit_log.belongs_to_user,
	audit_log.belongs_to_household,
	audit_log.created_at,
	(
		SELECT COUNT(audit_log.id)
		FROM audit_log
		WHERE audit_log.created_at > COALESCE(sqlc.narg(created_after), (SELECT NOW() - '999 years'::INTERVAL))
			AND audit_log.created_at < COALESCE(sqlc.narg(created_before), (SELECT NOW() + '999 years'::INTERVAL))
			AND audit_log.belongs_to_user = sqlc.arg(belongs_to_user)
			AND audit_log.resource_type = sqlc.arg(resource_type)
	) AS filtered_count,
	(
		SELECT COUNT(audit_log.id)
		FROM audit_log
		WHERE
			audit_log.belongs_to_user = sqlc.arg(belongs_to_user)
			AND audit_log.resource_type = sqlc.arg(resource_type)
	) AS total_count
FROM audit_log
WHERE audit_log.created_at > COALESCE(sqlc.narg(created_after), (SELECT NOW() - '999 years'::INTERVAL))
	AND audit_log.created_at < COALESCE(sqlc.narg(created_before), (SELECT NOW() + '999 years'::INTERVAL))
	AND audit_log.belongs_to_user = sqlc.arg(belongs_to_user)
	AND audit_log.resource_type = sqlc.arg(resource_type)
LIMIT sqlc.narg(query_limit)
OFFSET sqlc.narg(query_offset);

-- name: GetAuditLogsForHousehold :many

SELECT
	audit_log.id,
	audit_log.resource_type,
	audit_log.relevant_id,
	audit_log.event_type,
	audit_log.changes,
	audit_log.belongs_to_user,
	audit_log.belongs_to_household,
	audit_log.created_at,
	(
		SELECT COUNT(audit_log.id)
		FROM audit_log
		WHERE audit_log.created_at > COALESCE(sqlc.narg(created_after), (SELECT NOW() - '999 years'::INTERVAL))
			AND audit_log.created_at < COALESCE(sqlc.narg(created_before), (SELECT NOW() + '999 years'::INTERVAL))
			AND audit_log.belongs_to_household = sqlc.arg(belongs_to_household)
	) AS filtered_count,
	(
		SELECT COUNT(audit_log.id)
		FROM audit_log
		WHERE
			audit_log.belongs_to_household = sqlc.arg(belongs_to_household)
	) AS total_count
FROM audit_log
WHERE audit_log.created_at > COALESCE(sqlc.narg(created_after), (SELECT NOW() - '999 years'::INTERVAL))
	AND audit_log.created_at < COALESCE(sqlc.narg(created_before), (SELECT NOW() + '999 years'::INTERVAL))
	AND audit_log.belongs_to_household = sqlc.arg(belongs_to_household)
LIMIT sqlc.narg(query_limit)
OFFSET sqlc.narg(query_offset);

-- name: GetAuditLogsForHouseholdAndResourceType :many

SELECT
	audit_log.id,
	audit_log.resource_type,
	audit_log.relevant_id,
	audit_log.event_type,
	audit_log.changes,
	audit_log.belongs_to_user,
	audit_log.belongs_to_household,
	audit_log.created_at,
	(
		SELECT COUNT(audit_log.id)
		FROM audit_log
		WHERE audit_log.created_at > COALESCE(sqlc.narg(created_after), (SELECT NOW() - '999 years'::INTERVAL))
			AND audit_log.created_at < COALESCE(sqlc.narg(created_before), (SELECT NOW() + '999 years'::INTERVAL))
			AND audit_log.belongs_to_household = sqlc.arg(belongs_to_household)
			AND audit_log.resource_type = sqlc.arg(resource_type)
	) AS filtered_count,
	(
		SELECT COUNT(audit_log.id)
		FROM audit_log
		WHERE
			audit_log.belongs_to_household = sqlc.arg(belongs_to_household)
			AND audit_log.resource_type = sqlc.arg(resource_type)
	) AS total_count
FROM audit_log
WHERE audit_log.created_at > COALESCE(sqlc.narg(created_after), (SELECT NOW() - '999 years'::INTERVAL))
	AND audit_log.created_at < COALESCE(sqlc.narg(created_before), (SELECT NOW() + '999 years'::INTERVAL))
	AND audit_log.belongs_to_household = sqlc.arg(belongs_to_household)
	AND audit_log.resource_type = sqlc.arg(resource_type)
LIMIT sqlc.narg(query_limit)
OFFSET sqlc.narg(query_offset);
