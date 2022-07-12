package main

import (
	"time"
)

func GetMinskHour() int {
	return time.Now().In(time.FixedZone("Europe/Minsk", 3*60*60)).Hour()
}

var weekdayToRussian = map[time.Weekday]string{
	time.Monday:    "понедельник",
	time.Tuesday:   "вторник",
	time.Wednesday: "среду",
	time.Thursday:  "четверг",
	time.Friday:    "пятницу",
	time.Saturday:  "субботу",
	time.Sunday:    "воскресенье",
}

func Today() time.Time {
	return time.Now().Truncate(24 * time.Hour)
}
