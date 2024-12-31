package utils

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"
)

// NewTimeNowUTC In case of doesn't have any date in query param, we need use a constructor to create datetime with actual datetime
// because if we don't make this, when insert on database will be +3 hours.
// The default when insert a datetime in database is use UTC, when we use
// time.Now the time created haven't location,
// so when insert the database set to default location(that is UTC) and when make this the time change.
// We need create the datetime with this constructor to create the time with location UTC.
// ------Example------
// Created without UTC -> 2022-07-07 14:09:00
// Inserted on database -> 2022-07-07 17:09:00
// Because when is converted to UTC is added 3 hours.
//
// Created with UTC -> 2022-07-07 14:09:09
// Inserted on database -> 2022-07-07 14:09:00
// Because the database doesn't convert the time
// -------------------
func NewTimeNowUTC() *time.Time {
	now := time.Now()
	nowUTC := time.Date(
		now.Year(),
		now.Month(),
		now.Day(),
		now.Hour(),
		now.Minute(),
		now.Second(),
		0,
		time.UTC,
	)

	return &nowUTC
}

func NewDateTime(convert *time.Time) *DateTime {
	if convert != nil {
		x := DateTime(*convert)
		return &x
	}
	return nil
}

func NewDate(convert *time.Time) *Date {
	if convert != nil {
		x := Date(*convert)
		return &x
	}
	return nil
}

func DateTimeByString(convert string) *DateTime {
	layout := "2006-01-02 15:04:05"

	t, err := time.Parse(layout, convert)
	if err == nil {
		d := DateTime(t)
		return &d
	}

	return nil
}

func DateByString(convert string, layout string) *Date {
	t, err := time.Parse(layout, convert)
	if err != nil {
		log.Println("DateByString error", err)
		return nil
	}
	d := Date(t)
	return &d

}

func DateByStrings(convert string, layouts []string) *Date {
	var t time.Time
	var err error

	for _, layout := range layouts {
		t, err = time.Parse(layout, convert)
		if err == nil {
			d := Date(t)
			return &d
		}
	}

	log.Println("DateByString error", err)
	return nil
}

type DateTime time.Time
type Time time.Time

func (d *DateTime) IsNilOrZero() bool {
	if d == nil {
		return true
	}
	if d.Time() == nil {
		return true
	}
	return d.Time().IsZero()
}

// CorrectDateTime return nil if d.Time() == nil or d.Time is zero
// else return d.Time()
func (d *DateTime) CorrectDateTime() *time.Time {
	if d.Time() == nil || d.Time().IsZero() {
		return nil
	}

	return d.Time()
}

func (d DateTime) MarshalJSON() ([]byte, error) {
	date := time.Time(d).Format("2006-01-02 15:04:05")
	return []byte(fmt.Sprintf("%q", date)), nil
}

func (d *DateTime) UnmarshalJSON(b []byte) error {
	date := strings.Trim(string(b), `"`)
	if date == "null" || date == "" {
		return nil
	} else {
		tmp, err := time.Parse("2006-01-02 15:04:05", date)
		if err == nil {
			*d = DateTime(tmp)
		}
		return err
	}
}

func (d *Time) UnmarshalJSON(b []byte) error {
	date := strings.Trim(string(b), `"`)
	if date == "null" || date == "" {
		return nil
	} else {
		tmp, err := time.Parse("15:04", date)
		if err == nil {
			*d = Time(tmp)
		}
		return err
	}
}

func (d Time) MarshalJSON() ([]byte, error) {
	date := time.Time(d).Format("15:04")
	return []byte(fmt.Sprintf("%q", date)), nil
}

func (d *DateTime) Time() *time.Time {
	if d == nil {
		return nil
	}
	t := time.Time(*d)
	return &t
}

func (d *Time) Time() *time.Time {
	if d == nil {
		return nil
	}
	t := time.Time(*d)
	return &t
}

func (d *Date) CorrectDate() *time.Time {
	if d.Time() == nil || d.Time().IsZero() {
		return nil
	}

	return d.Time()
}

// </editor-fold>

type Date time.Time

func (d *Date) IsNilOrZero() bool {
	if d == nil {
		return true
	}
	if d.Time() == nil {
		return true
	}
	return d.Time().IsZero()
}

func (d Date) MarshalJSON() ([]byte, error) {
	df := time.Time(d).Format("2006-01-02")
	return []byte(fmt.Sprintf("%q", df)), nil
}

func (d *Date) UnmarshalJSON(b []byte) error {
	date := strings.Trim(string(b), `"`)
	if date == "null" || date == "" {
		return nil
	} else {
		tmp, err := time.Parse("2006-01-02", date)
		if err == nil {
			*d = Date(tmp)
		} else {
			return errors.New(fmt.Sprintf("Data %s inválida.", date))
		}

		return err
	}
}

func (d *Date) Time() *time.Time {
	if d == nil {
		return nil
	}
	t := time.Time(*d)
	return &t
}

func NewDuration(convert *time.Time) *Duration {
	if convert != nil {
		x := Duration(*convert)
		return &x
	}
	return nil
}

type Duration time.Time

func (d Duration) MarshalJSON() ([]byte, error) {
	df := time.Time(d).Format("15:04:05")
	return []byte(fmt.Sprintf("%q", df)), nil
}

func (d *Duration) UnmarshalJSON(b []byte) error {
	duration := strings.Trim(string(b), `"`)
	if duration == "null" || duration == "" {
		return nil
	} else {
		tmp, err := time.Parse("15:04:05", duration)
		if err == nil {
			*d = Duration(tmp)
		}
		return err
	}
}

func (d *Duration) Time() *time.Time {
	if d == nil {
		return nil
	}
	t := time.Time(*d)
	return &t
}

func GetDaysBetweenInclusive(date1 time.Time, date2 time.Time) []time.Time {
	var startDay time.Time
	var finalDay time.Time

	if date1.IsZero() || date2.IsZero() {
		return nil
	}

	if date1.Equal(date2) {
		return []time.Time{date1}
	}

	if date1.Before(date2) {
		startDay = date1
		finalDay = date2
	} else {
		startDay = date2
		finalDay = date1
	}

	var datesBetween []time.Time
	datesBetween = append(datesBetween, startDay)
	for startDay.Before(finalDay) {
		startDay = startDay.Add(time.Duration(time.Hour) * 24)
		datesBetween = append(datesBetween, startDay)
	}

	return datesBetween
}

// ClearTime return a new instance of time.Time, keeping only the date related fields
func ClearTime(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

// DateBefore check whether t1 is before t2, ignoring the time
// E.g:
//
//	t1 := time.Parse("20060402150405", "20230101235509") // 2023-01-01-23 55 09
//	t2 := time.Parse("20060402150405", "20250208235509") // 2025-02-08-23 55 09
//	fmt.Print(DateBefore(t1, t2))
//
// >> true
func DateBefore(t1 time.Time, t2 time.Time) bool {
	t1Cleaned := ClearTime(t1)
	t2Cleaned := ClearTime(t2)

	return t1Cleaned.Before(t2Cleaned)
}

// DateAfter check whether t1 is after t2, ignoring the time
// E.g:
//
//	t1 := time.Parse("20060402150405", "20230101235509") // 2023-01-01-23 55 09
//	t2 := time.Parse("20060402150405", "20250208235509") // 2025-02-08-23 55 09
//	fmt.Print(DateBefore(t1, t2))
//
// >> false
func DateAfter(t1 time.Time, t2 time.Time) bool {
	t1Cleaned := ClearTime(t1)
	t2Cleaned := ClearTime(t2)

	return t1Cleaned.After(t2Cleaned)
}

// DateEquals check whether t1 is equal to t2, ignoring the time
// E.g:
//
//	t1 := time.Parse("20060402150405", "20230101235509") // 2023-01-01-23 55 09
//	t2 := time.Parse("20060402150405", "20250208235509") // 2025-02-08-23 55 09
//	fmt.Print(DateBefore(t1, t2))
//
// >> false
func DateEquals(t1 time.Time, t2 time.Time) bool {
	t1Cleaned := ClearTime(t1)
	t2Cleaned := ClearTime(t2)

	return t1Cleaned.Equal(t2Cleaned)
}
