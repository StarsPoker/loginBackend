package date_utils

import "time"

const (
	apiDateLayout = "02/01/2006 15:04:05"
	apiDbLayout   = "2006-01-02 15:04:05"
)

func GetNow() time.Time {
	now := time.Now().UTC()
	return now
}

func GetNowString() string {
	return GetNow().Format(apiDateLayout)
}

func GetNowDBFormat() string {
	return GetNow().Format(apiDbLayout)
}
