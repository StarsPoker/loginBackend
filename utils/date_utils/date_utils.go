package date_utils

import (
	"fmt"
	"time"

	"github.com/go-shadow/moment"
	"github.com/panamafrancis/tizzy"
)

const (
	apiDateLayout = "02/01/2006 15:04:05"
	apiDbLayout   = "2006-01-02 15:04:05"
	apiBBLayout   = "2006-01-02T15:04:05"
)

func GetNow() time.Time {
	loc, err := tizzy.LoadLocation("America/Sao_Paulo")
	if err != nil {
		fmt.Println(err)
	}

	now := time.Now().In(loc)
	return now
}

func GetLocalDate(dateString string) string {
	loc, err := time.LoadLocation("America/Sao_Paulo")
	if err != nil {
		fmt.Println(err)
	}
	date, err := time.Parse(apiDbLayout, dateString)
	localDate := date.In(loc).Format(apiDbLayout)
	return localDate
}

func GetTimeFromNow(timeToAdd time.Duration) time.Time {
	now := GetNow()
	return now.Add(timeToAdd)
}

func GetNowMoment() int64 {
	return moment.New().ValueOf()
}

func GetNowString() string {
	return GetNow().Format(apiDateLayout)
}

func GetNowDBFormat() string {
	return GetNow().Format(apiDbLayout)
}

func GetUnixBbFormat(timestamp int64) string {
	unixTimeUTC := time.Unix(timestamp, 0)
	return unixTimeUTC.Format(apiBBLayout)
}

func GetUnixDbFormat(timestamp int64) string {
	unixTimeUTC := time.Unix(timestamp, 0)
	return unixTimeUTC.Format(apiDbLayout)
}

func GetTimeFromNowDbFormat(timeToAdd time.Duration) string {
	return GetTimeFromNow(timeToAdd).Format(apiDbLayout)
}

func GetTimeFromNowDbBB(timeToAdd time.Duration) string {
	return GetTimeFromNow(timeToAdd).Format(apiBBLayout)
}
