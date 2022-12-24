package util

import (
	"time"
)

func GetSpecificHourTimeOfTheDay(hour int) time.Time {
	loc := time.FixedZone("UTC", 0)

	now := time.Now()
	yyyy, mm, dd := now.Date()
	specificTime := time.Date(yyyy, mm, dd, hour, 0, 0, 0, loc)
	return specificTime
}

func GetTimeDifferenceFromNow(workStartTime, inTime time.Time) int {

	diff := inTime.Sub(workStartTime)

	return int(diff.Seconds())
}
