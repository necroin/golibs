package document

import (
	"fmt"
	"strconv"
	"time"

	"github.com/necroin/golibs/utils"
)

type Key struct {
	Name  string
	Value any
}

// String returns string representation of value.
func (key *Key) String() string {
	if key.Value == nil {
		return ""
	}
	return fmt.Sprintf("%v", key.Value)
}

// Bool returns bool type value.
func (key *Key) Bool() (bool, error) {
	return utils.ParseBool(key.String())
}

// Int returns int type value.
func (key *Key) Int() (int, error) {
	v, err := strconv.ParseInt(key.String(), 0, 64)
	return int(v), err
}

// Int64 returns int64 type value.
func (key *Key) Int64() (int64, error) {
	return strconv.ParseInt(key.String(), 0, 64)
}

// Uint returns uint type valued.
func (key *Key) Uint() (uint, error) {
	u, e := strconv.ParseUint(key.String(), 0, 64)
	return uint(u), e
}

// Uint64 returns uint64 type value.
func (key *Key) Uint64() (uint64, error) {
	return strconv.ParseUint(key.String(), 0, 64)
}

// Float64 returns float64 type value.
func (key *Key) Float64() (float64, error) {
	return strconv.ParseFloat(key.String(), 64)
}

// Duration returns time.Duration type value.
func (key *Key) Duration() (time.Duration, error) {
	return time.ParseDuration(key.String())
}

// TimeFormat parses with given format and returns time.Time type value.
func (key *Key) TimeFormat(format string) (time.Time, error) {
	return time.Parse(format, key.String())
}
