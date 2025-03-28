package util

import (
	"strconv"
	"time"
)

// Integer integer type constraints
type Integer interface {
	~uint | ~int | ~uint32 | ~int32 | ~uint64 | ~int64
}

// Float float type constraints
type Float interface {
	~float64 | ~float32
}

// Number number type constraints
type Number interface {
	Integer | Float
}

// Atot string to time helper
func Atot(a string, layout ...string) (time.Time, error) {
	l := "2006-01-02"
	if len(layout) > 0 {
		l = layout[0]
	}

	return time.Parse(l, a)
}

// Atof string to float
func Atof[T Float](s string) (T, error) {
	f, err := strconv.ParseFloat(s, 64)
	return T(f), err
}
