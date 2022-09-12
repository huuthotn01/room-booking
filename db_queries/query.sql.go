// Code generated by sqlc. DO NOT EDIT.
// source: query.sql

package db_queries

import (
	"context"
	"database/sql"
	"time"
)

const getAllRoom = `-- name: GetAllRoom :many
SELECT name, description, max_capacity FROM room FOR UPDATE
`

func (q *Queries) GetAllRoom(ctx context.Context) ([]Room, error) {
	rows, err := q.db.QueryContext(ctx, getAllRoom)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Room
	for rows.Next() {
		var i Room
		if err := rows.Scan(&i.Name, &i.Description, &i.MaxCapacity); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getAllUser = `-- name: GetAllUser :many
SELECT email, fullname, birthdate, mentor, start_date, end_date FROM user FOR UPDATE
`

func (q *Queries) GetAllUser(ctx context.Context) ([]User, error) {
	rows, err := q.db.QueryContext(ctx, getAllUser)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []User
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.Email,
			&i.Fullname,
			&i.Birthdate,
			&i.Mentor,
			&i.StartDate,
			&i.EndDate,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const insertUser = `-- name: InsertUser :execresult
INSERT INTO user(email, fullname, birthdate, mentor, start_date, end_date) VALUES(?, ?, ?, ?, ?, ?)
`

type InsertUserParams struct {
	Email     string         `json:"email"`
	Fullname  string         `json:"fullname"`
	Birthdate time.Time      `json:"birthdate"`
	Mentor    sql.NullString `json:"mentor"`
	StartDate sql.NullTime   `json:"start_date"`
	EndDate   sql.NullTime   `json:"end_date"`
}

func (q *Queries) InsertUser(ctx context.Context, arg InsertUserParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, insertUser,
		arg.Email,
		arg.Fullname,
		arg.Birthdate,
		arg.Mentor,
		arg.StartDate,
		arg.EndDate,
	)
}
