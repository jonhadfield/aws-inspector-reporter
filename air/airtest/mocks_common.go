package airtest

import "time"

func ptrToInt64(in int64) *int64 {
	return &in
}

func ptrToFloat64(in float64) *float64 {
	return &in
}

func ptrToStr(in string) *string {
	return &in
}

func ptrToTime(in time.Time) *time.Time {
	return &in
}

func ptrToBool(in bool) *bool {
	return &in
}
