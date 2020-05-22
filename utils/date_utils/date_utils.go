package date_utils

import "time"

const (
	apiDateLayout = "02/01/2006 15:04:05"
	apiDbLayout   = "2006-01-02 15:04:05"
)

func GetNow() time.Time {
	loc, err := time.LoadLocation("America/Sao_Paulo")
	if err != nil {
		panic(err)
	}
	now := time.Now().In(loc)
	return now
}

func GetNowString() string {
	return GetNow().Format(apiDateLayout)
}

func GetNowDBFormat() string {
	return GetNow().Format(apiDbLayout)
}
