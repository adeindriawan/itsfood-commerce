package utils

import (
	"time"
	"strconv"
)

func ConvertDateToPhrase(date time.Time, withTime bool) string {
	var months = [12]string{"Januari", "Februari", "Maret", "April", "Mei", "Juni", "Juli", "Agustus", "September", "Oktober", "November", "Desember"}
	var readableDate = ""
	var m int = int(date.Month())
	dateYear := strconv.Itoa(date.Year())
	dateMonth := months[m-1]
	dateDay := strconv.Itoa(date.Day())

	readableDate += dateDay + " " + dateMonth + " " + dateYear

	if withTime {
		prefixHour := ""
		prefixMinute := ""

		if date.Hour() < 10 {
			prefixHour += "0"
		}

		if date.Minute() < 10 {
			prefixMinute += "0"
		}

		dateHour := strconv.Itoa(date.Hour())
		dateMinute := strconv.Itoa(date.Minute())

		readableDate += " " + prefixHour + dateHour + ":" + prefixMinute + dateMinute
	}

	return readableDate
}