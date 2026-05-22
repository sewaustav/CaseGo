package profileService

import "time"

func isYesterday(t time.Time) bool {
	now := time.Now().In(t.Location())

	yesterday := now.AddDate(0, 0, -1)

	yYear, yMonth, yDay := yesterday.Date()
	tYear, tMonth, tDay := t.Date()

	return yYear == tYear && yMonth == tMonth && yDay == tDay
}

func isToday(t time.Time) bool {
	now := time.Now().In(t.Location())

	yYear, yMonth, yDay := now.Date()
	tYear, tMonth, tDay := t.Date()

	return yYear == tYear && yMonth == tMonth && yDay == tDay
}
