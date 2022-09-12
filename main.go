package main

import (
	"encoding/json"
	"fmt"
	db_manager "go-app/db_manager"
	"log"
	"net/http"
)

func UserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" { // get users info
		db, err := db_manager.DbConnect()
		if err != nil {
			fmt.Fprintf(w, "%s", err)
			return
		}
		defer db.Close()
		user_data, err := db_manager.GetAllUser(db)
		if err != nil {
			fmt.Fprintf(w, "%s", err)
			return
		}

		enc := json.NewEncoder(w)
		w.Header().Set("Content-Type", "application/json")
		err = enc.Encode(user_data)
		if err != nil {
			fmt.Fprintf(w, "%s", err)
			return
		}
	} else {
		http.Error(w, "Method not supported", http.StatusNotFound)
		return
	}
}

func RoomHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" { // get rooms info
		db, err := db_manager.DbConnect()
		if err != nil {
			fmt.Fprintf(w, "%s", err)
			return
		}
		defer db.Close()
		room_data, err := db_manager.GetRooms(db)
		if err != nil {
			fmt.Fprintf(w, "%s", err)
			return
		}

		enc := json.NewEncoder(w)
		w.Header().Set("Content-Type", "application/json")
		err = enc.Encode(room_data)
		if err != nil {
			fmt.Fprintf(w, "%s", err)
			return
		}
	} else {
		http.Error(w, "Method not supported", http.StatusNotFound)
		return
	}
}

//func RoomBookingHandler(w http.ResponseWriter, r *http.Request) {
//	if r.Method == "GET" { // get room booking info
//		//queries := r.URL.Query()
//		//room, start_time, end_time := queries.Get("room"), queries.Get("start_time"), queries.Get("end_time")
//		//fmt.Fprintf(w, "%s, %s, %s", room, start_time, end_time)
//		var user_filter db_manager.Booking_Info
//		dec := json.NewDecoder(r.Body)
//		err := dec.Decode(&user_filter)
//		if err != nil {
//			if !errors.Is(err, io.EOF) {
//				fmt.Fprintf(w, "JSON %s", err)
//				return
//			}
//		}
//
//		db, err := db_manager.DbConnect()
//		if err != nil {
//			fmt.Fprintf(w, "%s", err)
//			return
//		}
//		defer db.Close()
//		booking_data, err := db_manager.GetBookingInfo(db, user_filter.Room, user_filter.Start_time, user_filter.End_time)
//		if err != nil {
//			fmt.Fprintf(w, "%s", err)
//			return
//		}
//
//		enc := json.NewEncoder(w)
//		w.Header().Set("Content-Type", "application/json")
//		err = enc.Encode(booking_data)
//		if err != nil {
//			fmt.Fprintf(w, "%s", err)
//			return
//		}
//	} else if r.Method == "POST" { // submit room booking
//		decoder := json.NewDecoder(r.Body)
//		var booking_data db_manager.Booking_Info
//		err := decoder.Decode(&booking_data)
//		if err != nil {
//			fmt.Fprintf(w, "%s", err)
//			return
//		}
//
//		// Check input format with regex
//		booker_regex, _ := regexp.Compile("^[a-zA-Z0-9]+@teko.vn$")
//		room_regex, _ := regexp.Compile("[0-9][0-9][0-9][0-9]")
//		if !booker_regex.MatchString(booking_data.Booker) || !room_regex.MatchString(booking_data.Room) {
//			fmt.Fprintf(w, "Booker and Room wrong format")
//			return
//		}
//
//		db, err := db_manager.DbConnect()
//		if err != nil {
//			fmt.Fprintf(w, "%s", err)
//			return
//		}
//		defer db.Close()
//		err = db_manager.BookRoom(db, booking_data.Booker, booking_data.Room, booking_data.Start_time, booking_data.End_time, booking_data.Tagged_Users...)
//		if err != nil {
//			fmt.Fprintf(w, "%s", err)
//			return
//		}
//		fmt.Fprintf(w, "Room booking successfully!")
//		return
//	} else {
//		http.Error(w, "Method not supported", http.StatusNotFound)
//		return
//	}
//}

func main() {
	http.HandleFunc("/users", UserHandler)
	http.HandleFunc("/rooms", RoomHandler)
	//http.HandleFunc("/room_booking", RoomBookingHandler)

	if err := http.ListenAndServe(":4000", nil); err != nil {
		log.Fatal(err)
	}
}
