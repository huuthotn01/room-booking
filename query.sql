-- name: GetAllUser :many
SELECT * FROM user FOR UPDATE;

-- name: GetAllRoom :many
SELECT * FROM room FOR UPDATE;

-- name: InsertUser :execresult
INSERT INTO user(email, fullname, birthdate, mentor, start_date, end_date) VALUES(?, ?, ?, ?, ?, ?);