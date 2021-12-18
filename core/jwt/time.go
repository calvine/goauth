package jwt

import (
	"strconv"
	"time"
)

// Time is an alias of time.Time and is useful for unix time stamps coming in from jtws
type Time time.Time

func NewTime() Time {
	return Time(time.Now().UTC())
}

func NewTimeLocal() Time {
	return Time(time.Now())
}

func FromDuration(d time.Duration) Time {
	return Time(time.Now().Add(d).UTC())
}

func FromDurationLocal(d time.Duration) Time {
	return Time(time.Now().Add(d))
}

func (t Time) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatInt(time.Time(t).Unix(), 10)), nil
}

func (t *Time) UnmarshalJSON(s []byte) error {
	timestampString := string(s)
	q, err := strconv.ParseInt(timestampString, 10, 64)
	if err != nil {
		return err
	}
	*(*time.Time)(t) = time.Unix(q, 0)
	return nil
}

func (t Time) Unix() int64 {
	return time.Time(t).Unix()
}

func (t Time) TimeLocal() time.Time {
	return time.Time(t)
}

func (t Time) StringLocal(layout string) string {
	return t.TimeLocal().Format(layout)
}

func (t Time) Time() time.Time {
	return time.Time(t).UTC()
}

func (t Time) String(layout string) string {
	return t.Time().Format(layout)
}

func (t Time) IsZero() bool {
	return time.Time(t).IsZero()
}

func (t Time) IsInPast() bool {
	now := time.Now()
	return time.Time(t).Before(now)
}

func (t Time) IsInFuture() bool {
	now := time.Now()
	return time.Time(t).After(now)
}
