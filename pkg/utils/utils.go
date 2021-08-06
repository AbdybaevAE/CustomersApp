package utils

import (
	"io/ioutil"
	"math/rand"
	"text/template"
	"time"
)

const customerHashSize = 20
const birthDateLayout = "2006-01-02"

var dict = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

// format birthdate for input type = "date"
func FormatBirthDate(birthDate time.Time) string {
	return birthDate.Format(birthDateLayout)
}

// simple function to return random string hash(of size 20)
func RandomSizedString(size int) string {
	if size <= 0 {
		panic("size of string must be bigger than 0")
	}
	b := []rune{}
	for counter := 0; counter < size; counter++ {
		b = append(b, dict[rand.Intn(len(dict))])
	}
	return string(b)
}
func GenCustomerHash() string {
	return RandomSizedString(customerHashSize)
}

// it would be better to have singleton instance for less memory expenses.
// probably need to be refactored.
func LoadTemplates() *template.Template {
	var allFiles []string
	files, err := ioutil.ReadDir("./ui/html")
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		allFiles = append(allFiles, "./ui/html/"+file.Name())
	}
	templates, err := template.ParseFiles(allFiles...)
	if err != nil {
		panic(err)
	}
	return templates
}
