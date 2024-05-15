package main

import (
	"testing"
	"time"

	"github.com/blue-davinci/gosnip/internal/assert"
)

func TestHumanDate(t *testing.T) {
	tests := []struct {
		name     string
		testtime time.Time
		want     string
	}{
		{
			name:     "Normal UTC time",
			testtime: time.Date(2022, 3, 17, 10, 15, 0, 0, time.UTC),
			want:     "17 Mar 2022 at 10:15",
		},
		{
			name:     "Empty time supplied",
			testtime: time.Time{},
			want:     "",
		},
		{
			name:     "Different timezone: CET",
			testtime: time.Date(2022, 3, 17, 10, 15, 0, 0, time.FixedZone("CET", 1*60*60)),
			want:     "17 Mar 2022 at 09:15",
		},
	}
	// Loop over the test cases.
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := humanDate(test.testtime)
			assert.Equal(t, got, test.want)
		})
	}

}
