package custval

import (
	"testing"
	"time"
)

func TestIsValidBirthDate(t *testing.T) {
	type test struct {
		name string
		want bool
		args time.Time
	}
	tt := []test{
		{"normal age", true, time.Now().AddDate(-20, 0, 0)},
		{"too old", false, time.Now().AddDate(-70, 0, 0)},
		{"too young", false, time.Now().AddDate(-5, 0, 0)},
	}
	for _, tc := range tt {
		got := IsValidBirthDate(tc.args)
		if got != tc.want {
			t.Error("broken test ", tc.name)
		}
	}
}
