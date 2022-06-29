// Copyright 2016 - 2022 The excelize Authors. All rights reserved. Use of
// this source code is governed by a BSD-style license that can be found in
// the LICENSE file.
//
// Package excelize providing a set of functions that allow you to write to and
// read from XLAM / XLSM / XLSX / XLTM / XLTX files. Supports reading and
// writing spreadsheet documents generated by Microsoft Excel™ 2007 and later.
// Supports complex components by high compatibility, and provided streaming
// API for generating or reading data from a worksheet with huge amounts of
// data. This library needs Go version 1.15 or later.

package excelize

import (
	"math"
	"time"
)

const (
	nanosInADay    = float64((24 * time.Hour) / time.Nanosecond)
	dayNanoseconds = 24 * time.Hour
	maxDuration    = 290 * 364 * dayNanoseconds
	roundEpsilon   = 1e-9
)

var (
	daysInMonth           = []int{31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}
	excel1900Epoc         = time.Date(1899, time.December, 30, 0, 0, 0, 0, time.UTC)
	excel1904Epoc         = time.Date(1904, time.January, 1, 0, 0, 0, 0, time.UTC)
	excelMinTime1900      = time.Date(1899, time.December, 31, 0, 0, 0, 0, time.UTC)
	excelBuggyPeriodStart = time.Date(1900, time.March, 1, 0, 0, 0, 0, time.UTC).Add(-time.Nanosecond)
)

// timeToExcelTime provides a function to convert time to Excel time.
func timeToExcelTime(t time.Time) (float64, error) {
	// TODO in future this should probably also handle date1904 and like TimeFromExcelTime

	if t.Before(excelMinTime1900) {
		return 0.0, nil
	}

	tt := t
	diff := t.Sub(excelMinTime1900)
	result := float64(0)

	for diff >= maxDuration {
		result += float64(maxDuration / dayNanoseconds)
		tt = tt.Add(-maxDuration)
		diff = tt.Sub(excelMinTime1900)
	}

	rem := diff % dayNanoseconds
	result += float64(diff-rem)/float64(dayNanoseconds) + float64(rem)/float64(dayNanoseconds)

	// Excel dates after 28th February 1900 are actually one day out.
	// Excel behaves as though the date 29th February 1900 existed, which it didn't.
	// Microsoft intentionally included this bug in Excel so that it would remain compatible with the spreadsheet
	// program that had the majority market share at the time; Lotus 1-2-3.
	// https://www.myonlinetraininghub.com/excel-date-and-time
	if t.After(excelBuggyPeriodStart) {
		result += 1.0
	}
	return result, nil
}

// shiftJulianToNoon provides a function to process julian date to noon.
func shiftJulianToNoon(julianDays, julianFraction float64) (float64, float64) {
	switch {
	case -0.5 < julianFraction && julianFraction < 0.5:
		julianFraction += 0.5
	case julianFraction >= 0.5:
		julianDays++
		julianFraction -= 0.5
	case julianFraction <= -0.5:
		julianDays--
		julianFraction += 1.5
	}
	return julianDays, julianFraction
}

// fractionOfADay provides a function to return the integer values for hour,
// minutes, seconds and nanoseconds that comprised a given fraction of a day.
// values would round to 1 us.
func fractionOfADay(fraction float64) (hours, minutes, seconds, nanoseconds int) {
	const (
		c1us  = 1e3
		c1s   = 1e9
		c1day = 24 * 60 * 60 * c1s
	)

	frac := int64(c1day*fraction + c1us/2)
	nanoseconds = int((frac%c1s)/c1us) * c1us
	frac /= c1s
	seconds = int(frac % 60)
	frac /= 60
	minutes = int(frac % 60)
	hours = int(frac / 60)
	return
}

// julianDateToGregorianTime provides a function to convert julian date to
// gregorian time.
func julianDateToGregorianTime(part1, part2 float64) time.Time {
	part1I, part1F := math.Modf(part1)
	part2I, part2F := math.Modf(part2)
	julianDays := part1I + part2I
	julianFraction := part1F + part2F
	julianDays, julianFraction = shiftJulianToNoon(julianDays, julianFraction)
	day, month, year := doTheFliegelAndVanFlandernAlgorithm(int(julianDays))
	hours, minutes, seconds, nanoseconds := fractionOfADay(julianFraction)
	return time.Date(year, time.Month(month), day, hours, minutes, seconds, nanoseconds, time.UTC)
}

// doTheFliegelAndVanFlandernAlgorithm; By this point generations of
// programmers have repeated the algorithm sent to the editor of
// "Communications of the ACM" in 1968 (published in CACM, volume 11, number
// 10, October 1968, p.657). None of those programmers seems to have found it
// necessary to explain the constants or variable names set out by Henry F.
// Fliegel and Thomas C. Van Flandern.  Maybe one day I'll buy that jounal and
// expand an explanation here - that day is not today.
func doTheFliegelAndVanFlandernAlgorithm(jd int) (day, month, year int) {
	l := jd + 68569
	n := (4 * l) / 146097
	l = l - (146097*n+3)/4
	i := (4000 * (l + 1)) / 1461001
	l = l - (1461*i)/4 + 31
	j := (80 * l) / 2447
	d := l - (2447*j)/80
	l = j / 11
	m := j + 2 - (12 * l)
	y := 100*(n-49) + i + l
	return d, m, y
}

// timeFromExcelTime provides a function to convert an excelTime
// representation (stored as a floating point number) to a time.Time.
func timeFromExcelTime(excelTime float64, date1904 bool) time.Time {
	var date time.Time
	wholeDaysPart := int(excelTime)
	// Excel uses Julian dates prior to March 1st 1900, and Gregorian
	// thereafter.
	if wholeDaysPart <= 61 {
		const OFFSET1900 = 15018.0
		const OFFSET1904 = 16480.0
		const MJD0 float64 = 2400000.5
		var date time.Time
		if date1904 {
			date = julianDateToGregorianTime(MJD0, excelTime+OFFSET1904)
		} else {
			date = julianDateToGregorianTime(MJD0, excelTime+OFFSET1900)
		}
		return date
	}
	floatPart := excelTime - float64(wholeDaysPart) + roundEpsilon
	if date1904 {
		date = excel1904Epoc
	} else {
		date = excel1900Epoc
	}
	durationPart := time.Duration(nanosInADay * floatPart)
	return date.AddDate(0, 0, wholeDaysPart).Add(durationPart).Truncate(time.Second)
}

// ExcelDateToTime converts a float-based excel date representation to a time.Time.
func ExcelDateToTime(excelDate float64, use1904Format bool) (time.Time, error) {
	if excelDate < 0 {
		return time.Time{}, newInvalidExcelDateError(excelDate)
	}
	return timeFromExcelTime(excelDate, use1904Format), nil
}

// isLeapYear determine if leap year for a given year.
func isLeapYear(y int) bool {
	if y == y/400*400 {
		return true
	}
	if y == y/100*100 {
		return false
	}
	return y == y/4*4
}

// getDaysInMonth provides a function to get the days by a given year and
// month number.
func getDaysInMonth(y, m int) int {
	if m == 2 && isLeapYear(y) {
		return 29
	}
	return daysInMonth[m-1]
}

// validateDate provides a function to validate if a valid date by a given
// year, month, and day number.
func validateDate(y, m, d int) bool {
	if m < 1 || m > 12 {
		return false
	}
	if d < 1 {
		return false
	}
	return d <= getDaysInMonth(y, m)
}

// formatYear converts the given year number into a 4-digit format.
func formatYear(y int) int {
	if y < 1900 {
		if y < 30 {
			y += 2000
		} else {
			y += 1900
		}
	}
	return y
}
