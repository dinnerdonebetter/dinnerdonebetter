INSERT INTO
	webhooks (
	id,
	name,
	content_type,
	url,
	method,
	belongs_to_household
)
VALUES
	($1, $2, $3, $4, $5, $6);