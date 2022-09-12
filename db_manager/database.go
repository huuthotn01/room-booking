package db_manager

import (
	"context"
	"database/sql"
	"fmt"
	"go-app/db_queries"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const (
	username = "root"
	password = ""
	hostname = "localhost"
	db_name  = "room_booking"
)

func NewNullString(s string) sql.NullString {
	if len(s) == 0 {
		return sql.NullString{}
	}
	return sql.NullString{
		s, true,
	}
}

func Dsn(dbName string) string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", username, password, hostname, dbName)
}

func DbConnect() (*sql.DB, error) {
	db, err := sql.Open("mysql", Dsn(db_name))
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(20)
	db.SetConnMaxLifetime(time.Minute * 5)

	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}
	log.Printf("Connected to db %s successfully", db_name)
	return db, nil
}

func lockTable(db *sql.DB, table_name string, mode string) error {
	query := ""
	if mode == "READ" {
		query += "LOCK TABLES ? READ"
	} else {
		query += "LOCK TABLES ? WRITE"
	}
	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(table_name)
	if err != nil {
		return err
	}
	return nil
}

func unlockTable(db *sql.DB) error {
	query := "UNLOCK TABLES"
	_, err := db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

type User_Data struct {
	Email      string         `json:"email"`
	Fullname   string         `json:"fullname"`
	Birthdate  time.Time      `json:"Birthdate"`
	Mentor     sql.NullString `json:"Mentor"`
	Start_Date sql.NullTime   `json:"start_date"`
	End_Date   sql.NullTime   `json:"end_date"`
}

type Room_Data struct {
	Name    string `json:"name"`
	Descr   string `json:"description"`
	MaxCapa int    `json:"max_capa"`
}

type Booking_Info struct {
	Booker       string   `json:"booker"`
	Room         string   `json:"room"`
	Start_time   string   `json:"start_time"`
	End_time     string   `json:"end_time"`
	Tagged_Users []string `json:"tagged_users"`
}

func GetAllUser(db *sql.DB) ([]db_queries.User, error) {
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	query := db_queries.New(db)
	rows, err := query.GetAllUser(ctx)
	if err != nil {
		return nil, err
	}
	return rows, err
}

func InsertUser(db *sql.DB, usr User_Data) error {
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	query := db_queries.New(db)
	_, err := query.InsertUser(ctx, db_queries.InsertUserParams{
		Email:     usr.Email,
		Fullname:  usr.Fullname,
		Birthdate: usr.Birthdate,
		Mentor:    usr.Mentor,
		StartDate: usr.Start_Date,
		EndDate:   usr.End_Date,
	})

	if err != nil {
		return err
	}
	//no_rows, err := res.RowsAffected()
	//if err != nil {
	//	return err
	//}
	//log.Printf("Rows affected: %d", no_rows)
	//
	//log.Printf("User with Email %s created", usr.Email)
	//return nil
}

func GetRooms(db *sql.DB) ([]db_queries.Room, error) {
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	query := db_queries.New(db)
	rows, err := query.GetAllRoom(ctx)
	if err != nil {
		return nil, err
	}
	return rows, err
}

func GetBookingInfo(db *sql.DB, room string, start_time string, end_time string) ([]Booking_Info, error) {
	query := "SELECT * FROM booking WHERE 1=1"
	var filter []interface{}
	if room != "" {
		query += " AND room = ?"
		filter = append(filter, room)
	}
	if start_time != "" {
		query += " AND start_time >= ?"
		filter = append(filter, start_time)
	}
	if end_time != "" {
		query += " AND end_time <= ?"
		filter = append(filter, end_time)
	}
	query += " FOR UPDATE"
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return []Booking_Info{}, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, filter...)
	if err != nil {
		return []Booking_Info{}, err
	}
	defer rows.Close()

	var bookings = []Booking_Info{}
	for rows.Next() {
		var bks Booking_Info
		var id int
		if err := rows.Scan(&id, &bks.Booker, &bks.Room, &bks.Start_time, &bks.End_time); err != nil {
			return []Booking_Info{}, err
		}
		booking_detail_query := "SELECT tagged_user FROM userbooking WHERE id = ? FOR UPDATE"
		stm, err := db.PrepareContext(ctx, booking_detail_query)
		if err != nil {
			return []Booking_Info{}, err
		}
		defer stm.Close()
		r, e := stm.QueryContext(ctx, id)
		if e != nil {
			return []Booking_Info{}, e
		}
		defer r.Close()
		for r.Next() {
			var u string
			if e := r.Scan(&u); e != nil {
				return []Booking_Info{}, err
			}
			bks.Tagged_Users = append(bks.Tagged_Users, u)
		}
		bookings = append(bookings, bks)
	}
	if err := rows.Err(); err != nil {
		return []Booking_Info{}, err
	}
	return bookings, nil
}

func checkBookingTime(start_time string, end_time string) (string, string, error) {
	start, err := time.Parse("2006-01-02 15:04:05", start_time+":00")
	if err != nil {
		return start_time, end_time, err
	}
	end, err := time.Parse("2006-01-02 15:04:05", end_time+":00")
	if err != nil {
		return start_time, end_time, err
	}
	if start.Minute()%5 != 0 || end.Minute()%5 != 0 {
		return start_time, end_time, fmt.Errorf("Booking time must be divisible by 5, got %d and %d", start.Minute(), end.Minute())
	}
	if !start.Before(end) {
		return start_time, end_time, fmt.Errorf("Booking time not valid")
	}
	return start_time + ":00", end_time + ":00", nil
}

func removeDuplicate(str_slice []string) []string {
	all_keys := make(map[string]bool)
	list := []string{}
	for _, item := range str_slice {
		if _, value := all_keys[item]; !value {
			all_keys[item] = true
			list = append(list, item)
		}
	}
	return list
}

func BookRoom(db *sql.DB, booker string, room string, start_time string, end_time string, tagged_users ...string) error {
	start, end, err := checkBookingTime(start_time, end_time)
	if err != nil {
		return err
	}

	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()

	// Check booking time overlapping
	// Conditon: room=`room` and start_time <= `end` and end_time >= `start`
	query_overlap := "SELECT COUNT(*) FROM booking WHERE room=? AND start_time <= ? AND end_time >= ? FOR UPDATE"
	stmt, err := db.PrepareContext(ctx, query_overlap)
	if err != nil {
		return err
	}
	defer stmt.Close()

	num_row, err := stmt.QueryContext(ctx, room, end, start)
	if err != nil {
		return err
	}
	defer num_row.Close()
	var nrow int
	num_row.Next()
	if e := num_row.Scan(&nrow); e != nil {
		return e
	}

	if nrow > 0 {
		return errors.New("Overlapping booking")
	}

	// Begin Transaction
	tx, err := db.BeginTx(ctx, nil) // begin Transaction
	if err != nil {
		return err
	}
	defer tx.Rollback() // Rollback when necessary

	// Insert into table `booking`
	query_booking := "INSERT INTO booking (booker, room, start_time, end_time) VALUES (?, ?, ?, ?)"
	stmt, err = tx.PrepareContext(ctx, query_booking)
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := tx.ExecContext(ctx, query_booking, booker, room, start, end)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	// Insert into table `userbooking`
	query_userbks := "INSERT INTO userbooking VALUES "
	var userbks_inserts []string
	var params []interface{}
	userbks_inserts = append(userbks_inserts, "(?, ?)")
	params = append(params, id, booker)
	filtered_tagged_user := removeDuplicate(tagged_users)
	for _, v := range filtered_tagged_user {
		userbks_inserts = append(userbks_inserts, "(?, ?)")
		params = append(params, id, v)
	}
	query_tail := strings.Join(userbks_inserts, ",")
	query_userbks = query_userbks + query_tail
	stmt, err = tx.PrepareContext(ctx, query_userbks)
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err = stmt.ExecContext(ctx, params...)
	if err != nil {
		return err
	}
	no_rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	log.Printf("Num rows affected: %d", no_rows)

	// Commit if successful
	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}
