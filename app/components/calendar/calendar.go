package calendar

import (
	"strings"
	"time"
)

type Weeks []Week
type Week [8]int

func (w Week) Start(wd time.Weekday, ptr time.Time) Week {
	_, weekNum := ptr.ISOWeek()
	ints := [8]int{weekNum}

	weekday := ptr.Weekday()
	shift := (7+weekday-wd)%7 + 1
	day := 1

	for i := shift; i < 8; i++ {
		ints[i] = day
		day++
	}

	return ints
}

func (w Week) Fill(year int, month time.Month) Weeks {
	weeks := Weeks{w}
	end := time.Date(year, month, 1, 0, 0, 0, 0, time.Local).AddDate(0, 1, 0).AddDate(0, 0, -1)

	for day := w[7] + 1; day <= end.Day(); day += 7 {
		_, isoWeek := time.Date(year, month, day, 0, 0, 0, 0, time.Local).ISOWeek()
		week := Week{isoWeek}

		for i := 0; i < 7; i++ {
			d := day + i
			if d > end.Day() {
				break
			}

			week[i+1] = d
		}

		weeks = append(weeks, week)
	}

	return weeks
}

func (w Week) FillWeek(ptr time.Time) Week {
	_, w[0] = ptr.ISOWeek()

	for i := 1; i <= 7; i++ {
		w[i] = ptr.Day()

		if _, wn := ptr.ISOWeek(); wn > w[0] {
			w[0] = wn
		}

		ptr = ptr.AddDate(0, 0, 1)
	}

	return w
}

func (w Weeks) FigureOut(year int, month time.Month, wd time.Weekday) Weeks {
	ptr := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)

	return Week{}.Start(wd, ptr).Fill(year, month)
}

type YearMonth struct {
	year  int
	month time.Month
}

func NewYearMonth(year int, month time.Month) YearMonth {
	return YearMonth{
		year:  year,
		month: month,
	}
}

func (m YearMonth) Calendar(wd time.Weekday) Calendar {
	return Calendar{
		month: m.month,
		wd:    wd,
		weeks: MonthWeeklies(wd, m.year, m.month),
	}
}

type Calendar struct {
	wd    time.Weekday
	weeks Weeklies
	month time.Month
}

func (c Calendar) WeekLayout(weekNum bool) string {
	line := strings.Repeat("c", 7)
	if !weekNum {
		return line
	}

	return "c|" + line
}

func (c Calendar) WeekHeader(weekNum bool) string {
	names := append(make([]string, 0, 8), "W")

	for i := time.Sunday; i < 7; i++ {
		names = append(names, ((c.wd + i) % 7).String()[:1])
	}

	if !weekNum {
		names = names[1:]
	}

	return strings.Join(names, " & ")
}

func (c Calendar) WeekHeaderFull(weekNum bool) string {
	names := make([]string, 0, 7)

	for i := time.Sunday; i < 7; i++ {
		names = append(names, "\\hfil{}"+((c.wd+i)%7).String())
	}

	out := strings.Join(names, " & ")
	if weekNum {
		out = "& " + out
	}

	return out
}

func (c Calendar) WeekHeaderLen(weekNum bool) int {
	if weekNum {
		return 8
	}

	return 7
}

func (c Calendar) MonthName() time.Month {
	return c.month
}

func (c Calendar) Matrix() Weeklies {
	return c.weeks
}
