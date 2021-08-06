package custval

import "time"

// birthdate constraints
const (
	MinAge = 18
	MaxAge = 60
)

// custom validation for birthdate
func IsValidBirthDate(birthDate time.Time) bool {
	birthYear, birthMonth, birthDay := birthDate.Date()
	currYear, currMonth, currDay := time.Now().Date()
	age := currYear - birthYear
	if currMonth < birthMonth || currMonth == birthMonth && currDay < birthDay {
		age--
	}
	return age >= MinAge && age <= MaxAge
}

// compute available birthdate range for current time
func ComputeBirthDateRange() (time.Time, time.Time) {
	minDate := time.Now().AddDate(-MaxAge, 0, 0)
	maxDate := time.Now().AddDate(-MinAge, 0, 0)
	return minDate, maxDate
}
