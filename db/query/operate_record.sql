-- name: CreateOperateRecord :one
INSERT INTO operate_record (
    user_id,
    food_id,
    before,
    after,
    operate_category
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING *;

-- name: GetOperateRecords :many
SELECT *
FROM operate_record
LIMIT $1 OFFSET $2;

-- name: GetOperateRecordsByUserID :many
SELECT *
FROM operate_record
WHERE user_id = $1
LIMIT $2 OFFSET $3;

-- name: GetOperateRecordsByDateRange :many
SELECT *
FROM operate_record
WHERE update_at >= $1 AND update_at <= $2
LIMIT $3 OFFSET $4;
