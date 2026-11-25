package models

import "time"

// Zodiac wip
type Zodiac struct {
	date        time.Time
	Animal      string
	Element     string
	IsMasculine bool
	Color       string
}

func NewZodiac(date time.Time) *Zodiac {
	return &Zodiac{
		date:        date,
		Animal:      getAnimal(date),
		Element:     getElement(date),
		Color:       getColor(date),
		IsMasculine: isMasculine(date),
	}
}

func (z *Zodiac) GetDate() time.Time {
	return z.date
}

// getTsagaanSarDate returns the Mongolian Lunar New Year date.
// TODO: wip
func getTsagaanSarDate(date time.Time) time.Time {
	return date
}

var colors = [5]string{"blue", "red", "yellow", "white", "black"}

var elements = [5]string{"wood", "fire", "earth", "iron", "water"}

var animals = [12]string{
	"mouse",
	"cattle",
	"tiger",
	"rabbit",
	"dragon",
	"snake",
	"horse",
	"sheep",
	"monkey",
	"chicken",
	"dog",
	"pig",
}

// mod simple modulo operation, returns alway positive remainder
func mod(a, b int) int {
	remainder := a % b
	if remainder < 0 {
		remainder += b
	}

	return remainder
}

// IsMasculine calculates the year polarity
// Even numbers are mascuine, odd are female.
func isMasculine(date time.Time) bool {
	if mod(date.Year(), 2) == 0 {
		return true
	}

	return false
}

// getLunarYearOffset returns the offset that has to be substracted
// from year to find element, animal, and color etc.
func getLunarYearOffset(date time.Time) int {
	tsagaanSar := getTsagaanSarDate(date)
	// Check if date is before Mongolian Lunar New Year
	if date.YearDay() < tsagaanSar.YearDay() {
		// date belongs to the year before thus return 4
		return 4
	}

	return 3
}

func getAnimal(date time.Time) string {
	year := date.Year() - getLunarYearOffset(date)
	index := mod(year, 12)

	return animals[index]
}

func getElement(date time.Time) string {
	year := date.Year() - getLunarYearOffset(date)
	index := mod(year, 10) / 2

	return elements[index]
}

func getColor(date time.Time) string {
	year := date.Year() - getLunarYearOffset(date)
	index := mod(year, 10)

	return colors[index]
}
