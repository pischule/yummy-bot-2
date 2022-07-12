package main

import (
	"strings"
	"time"
	"unicode/utf8"
)

type NameValidationError string

func (n NameValidationError) Error() string {
	return string(n)
}

func allCyrillicLetters(s string) bool {
	for _, r := range s {
		if !isCyrillicLetter(r) {
			return false
		}
	}
	return true
}

func isCyrillicLetter(r int32) bool {
	return r >= 0x0400 && r <= 0x04FF
}

func ValidateName(name string) (string, error) {
	trimmed := strings.TrimSpace(name)
	if trimmed == "" {
		return "", NameValidationError("укажите имя после команды /name")
	}

	if utf8.RuneCountInString(trimmed) <= 1 {
		return "", NameValidationError("имя должно быть длиннее 1 символа")
	}

	if !allCyrillicLetters(trimmed) {
		return "", NameValidationError("имя должно содержать только кириллические символы")
	}

	return trimmed, nil
}

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
