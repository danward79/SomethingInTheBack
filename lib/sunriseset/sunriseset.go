// Package sunriseset provides sunrise and sunset events for the location provided.
package sunriseset

import (
	"fmt"
	"time"

	"github.com/keep94/sunrise"
)

type Loc struct {
	Location string
}

func New(loc string) *Loc {
	&Loc{Location: loc}
}

func (self *Loc) Start() {

	var s sunrise.Sunrise

	// Start time is June 1, 2013 PST
	location, _ := time.LoadLocation("Australia/Melbourne")
	startTime := time.Date(2014, 8, 7, 0, 0, 0, 0, location)

	// Coordinates of LA are 34.05N 118.25W
	s.Around(-37.88, 144.98, startTime)

	for s.Sunrise().Before(startTime) {
		s.AddDays(1)
	}

	formatStr := "Jan 2 15:04:05"
	for i := 0; i < 5; i++ {
		fmt.Printf("Sunrise: %s Sunset: %s\n", s.Sunrise().Format(formatStr), s.Sunset().Format(formatStr))
		s.AddDays(1)
	}

}
