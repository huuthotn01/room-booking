package db_manager

import (
	"testing"
)

func TestDSN(t *testing.T) {
	got := Dsn(db_name)
	want := "root:@tcp(localhost)/room_booking"

	if got != want {
		t.Errorf("Expected %s, got %s", want, got)
	}
}

func TestDbConnect(t *testing.T) {
	db, err := DbConnect()
	if err != nil {
		t.Errorf("Expected successful connection, got %s", err)
	}
	defer db.Close()
}

func TestInsertUser(t *testing.T) {
	db, _ := DbConnect()
	defer db.Close()

	var usr = User_Data{
		"teko@teko.vn",
		"Lê Văn Teko",

		NewNullString(""),
		NewNullString(""),
		NewNullString(""),
	}
	err := InsertUser(db, usr)
	if err != nil {
		t.Errorf("Expected successful query, got %s", err)
	}

	var usr_concurrent = []User_Data{
		{
			"lock@teko.vn",
			"Phạm Thị Lock",
			"2002-02-20",
			NewNullString(""),
			NewNullString(""),
			NewNullString(""),
		},
		{
			"lock@teko.vn",
			"Nguyễn Lê Mutex",
			"2002-02-20",
			NewNullString(""),
			NewNullString(""),
			NewNullString(""),
		},
	}
	var wg sync.WaitGroup
	
	for i := 0; i <= 1; i++ {
		wg.Add(1)
		i := i
		go func() {
			defer wg.Done()
			err := InsertUser(db, usr_concurrent[i])
			if err == nil {
				t.Errorf("Expected key constraint violated, got successful")
			}
			fmt.Printf("%s", err)
		}()
	}
	wg.Wait()
}

func TestGetAllUser(t *testing.T) {
	db, _ := DbConnect()
	defer db.Close()
	_, err := GetAllUser(db)
	//want := []User_Data{
	//	{"a@teko.vn", "Trần Văn Test", "2022-06-27", sql.NullString{"", false}, sql.NullString{"", false}, sql.NullString{"", false}},
	//	{"abc@teko.vn", "Trần Văn ABC", "2020-03-24", sql.NullString{"", false}, sql.NullString{"", false}, sql.NullString{"", false}},
	//	{"b@teko.vn", "Nguyễn Thị Test B", "2020-12-21", sql.NullString{"", false}, sql.NullString{"", false}, sql.NullString{"", false}},
	//	{"teko@teko.vn", "Lê Văn Teko", "2002-02-20", sql.NullString{"", false}, sql.NullString{"", false}, sql.NullString{"", false}},
	//	{"test@teko.vn", "Trần Văn Test", "2022-06-27", sql.NullString{"abc@teko.vn", true}, sql.NullString{"2022-06-27", true}, sql.NullString{"", false}},
	//}
	if err != nil {
		t.Errorf("Expected successful query, got %s", err)
		//} else {
		//	for i, d := range want {
		//		if !(got[i].Email == d.Email && got[i].Fullname == d.Fullname && got[i].Birthdate == d.Birthdate && got[i].Mentor == d.Mentor && got[i].Start_Date == d.Start_Date && got[i].End_Date == d.End_Date) {
		//			t.Errorf("Want %v, got %v", d, got[i])
		//		}
		//	}
		//}
	}
}

func TestGetBookingInfo(t *testing.T) {
	db, _ := DbConnect()
	defer db.Close()
	want := []Booking_Info{
		{"test@teko.vn", "0802", "2022-06-27 14:15:00", "2022-06-27 14:30:00", []string{"a@teko.vn", "b@teko.vn", "test@teko.vn"}},
		{"test@teko.vn", "1006", "2022-06-27 14:20:00", "2022-06-27 15:00:00", []string{"test@teko.vn"}},
		{"test@teko.vn", "1006", "2022-06-27 15:20:00", "2022-06-27 15:40:00", []string{"a@teko.vn", "b@teko.vn", "test@teko.vn"}},
		{"a@teko.vn", "0802", "2022-06-28 13:00:00", "2022-06-28 15:00:00", []string{"a@teko.vn"}},
		{"a@teko.vn", "0802", "2022-06-28 16:00:00", "2022-06-28 16:30:00", []string{"a@teko.vn", "b@teko.vn", "test@teko.vn"}},
		{"a@teko.vn", "0802", "2022-06-28 17:00:00", "2022-06-28 17:30:00", []string{"a@teko.vn", "b@teko.vn", "teko@teko.vn", "test@teko.vn"}},
	}
	got, err := GetBookingInfo(db, "", "", "")
	if err != nil {
		t.Errorf("Expected successful query, got %s", err)
		return
	} else if len(got) != len(want) {
		t.Errorf("Num rows retrieved incorrect.")
		return
	} else {
		for i, d := range want {
			if !(got[i].Booker == d.Booker && got[i].Room == d.Room && got[i].Start_time == d.Start_time && got[i].End_time == d.End_time && testEq(got[i].Tagged_Users, d.Tagged_Users)) {
				t.Errorf("Want %v, got %v", d, got[i])
				return
			}
		}
	}

	want = []Booking_Info{
		{"test@teko.vn", "0802", "2022-06-27 14:15:00", "2022-06-27 14:30:00", []string{"a@teko.vn", "b@teko.vn", "test@teko.vn"}},
		{"a@teko.vn", "0802", "2022-06-28 13:00:00", "2022-06-28 15:00:00", []string{"a@teko.vn"}},
		{"a@teko.vn", "0802", "2022-06-28 16:00:00", "2022-06-28 16:30:00", []string{"a@teko.vn", "b@teko.vn", "test@teko.vn"}},
		{"a@teko.vn", "0802", "2022-06-28 17:00:00", "2022-06-28 17:30:00", []string{"a@teko.vn", "b@teko.vn", "teko@teko.vn", "test@teko.vn"}},
	}
	got, err = GetBookingInfo(db, "0802", "", "")
	if err != nil {
		t.Errorf("Expected successful query, got %s", err)
		return
	} else if len(got) != len(want) {
		t.Errorf("Num rows retrieved incorrect.")
		return
	} else {
		for i, d := range want {
			if !(got[i].Booker == d.Booker && got[i].Room == d.Room && got[i].Start_time == d.Start_time && got[i].End_time == d.End_time && testEq(got[i].Tagged_Users, d.Tagged_Users)) {
				t.Errorf("Want %v, got %v", d, got[i])
				return
			}
		}
	}
}

func testEq(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestGetRoom(t *testing.T) {
	db, _ := DbConnect()
	defer db.Close()
	_, err := GetRooms(db)
	if err != nil {
		t.Errorf("Expected successful query, got %s", err)
	}
}

func TestRemoveDuplicate(t *testing.T) {
	input := []string{"a", "a", "a", "a"}
	want := []string{"a"}
	got := removeDuplicate(input)
	if !testEq(want, got) {
		t.Errorf("Expected %s, got %s", want, got)
	}
}

func TestBooking(t *testing.T) {
	db, _ := DbConnect()
	defer db.Close()

	err := BookRoom(db, "test@teko.vn", "0802", "2022-06-27 14:15", "2022-06-27 14:30", "a@teko.vn", "b@teko.vn")
	if err != nil {
		t.Errorf("Expected successful booking, got %s", err)
	}

	err = BookRoom(db, "test@teko.vn", "1006", "2022-06-27 15:15", "2022-06-27 14:33", "a@teko.vn", "b@teko.vn")
	if err == nil {
		t.Errorf("Expected Booking time must be divisible by 5, got 15 and 33, got successful")
	} else if err.Error() != "Booking time must be divisible by 5, got 15 and 33" {
		t.Errorf("Expected Booking time must be divisible by 5, got 15 and 33, got %s", err)
	}

	err = BookRoom(db, "test@teko.vn", "1006", "2022-06-27 15:15", "2022-06-27 14:30", "a@teko.vn", "b@teko.vn")
	if err == nil {
		t.Errorf("Expected Booking time not valid, got successful")
	} else if err.Error() != "Booking time not valid" {
		t.Errorf("Expected Booking time not valid, got %s", err)
	}

	err = BookRoom(db, "test@teko.vn", "0802", "2022-06-27 14:20", "2022-06-27 15:00")
	if err == nil {
		t.Errorf("Expected Overlapping booking, got successful")
	} else if err.Error() != "Overlapping booking" {
		t.Errorf("Expected Overlapping booking, got %s", err)
	}

	err = BookRoom(db, "test@teko.vn", "1006", "2022-06-27 14:20", "2022-06-27 15:00")
	if err != nil {
		t.Errorf("Expected successful booking, got %s", err)
	}

	err = BookRoom(db, "test@teko.vn", "1006", "2022-06-27 15:20", "2022-06-27 15:40", "a@teko.vn", "b@teko.vn", "b@teko.vn", "a@teko.vn")
	if err != nil {
		t.Errorf("Expected successful booking, got %s", err)
	}
}
