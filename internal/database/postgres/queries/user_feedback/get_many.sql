-- name: GetUserFeedback :many
SELECT
	user_feedback.id,
	user_feedback.prompt,
	user_feedback.feedback,
	user_feedback.rating,
	user_feedback.context,
	user_feedback.by_user,
	user_feedback.created_at,
	(
	 SELECT
		COUNT(user_feedback.id)
	 FROM
		user_feedback
	 WHERE
		user_feedback.created_at > COALESCE($1, (SELECT NOW() - interval '999 years'))
	 AND user_feedback.created_at < COALESCE($2, (SELECT NOW() + interval '999 years'))
	) as filtered_count,
	(
	 SELECT
		COUNT(user_feedback.id)
	 FROM
		user_feedback
	) as total_count
FROM
	user_feedback
WHERE
	user_feedback.created_at > COALESCE($1, (SELECT NOW() - interval '999 years'))
	AND user_feedback.created_at < COALESCE($2, (SELECT NOW() + interval '999 years'))
GROUP BY
	user_feedback.id
ORDER BY
	user_feedback.id
	LIMIT $3;
