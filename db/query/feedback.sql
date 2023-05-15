-- name: CreateFeedback :one
INSERT INTO feedback (
    user_id,
    food_id,
    edit_by
) VALUES (
    $1, $2, $3
)
RETURNING *;

-- name: GetFeedback :many
SELECT *
FROM feedback
LIMIT $1 OFFSET $2;

-- name: GetFeedbackByDateRange :many
SELECT *
FROM feedback
WHERE created_at >= $1 AND created_at <= $2
LIMIT $3 OFFSET $4;

-- name: GetFeedbackByStatus :many
SELECT *
FROM feedback;
