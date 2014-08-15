// Package sunriseset provides sunrise and sunset events for the location provided.
package sunriseset

import (
	"fmt"
	"github.com/danward79/sunrise"
	"github.com/robfig/cron"
	"time"
)

//Loc stores data about the location of sunrise and set required.
type Loc struct {
	location  *sunrise.Location
	formatStr string
	cronSch   *cron.Cron
}

//New returns a new location
func New(latitude float64, longitude float64) *Loc {
	return &Loc{location: sunrise.NewLocation(latitude, longitude), formatStr: "Jan 2 15:04:05", cronSch: cron.New()}
}

//Start the process
func (l *Loc) Start() chan map[string]interface{} {
	chOut := make(chan map[string]interface{})

	//Is it before or after todays sunrise/sunset?
	l.location.Today()
	tSunrise := l.location.Sunrise()
	if time.Now().After(tSunrise) {
		tSunrise = l.nextSunrise()
	}

	tSunset := l.nextSunset()
	if time.Now().After(tSunset) {
		tSunset = l.nextSunset()
	}

	//Schedule cron
	l.sunriseSchedule(tSunrise)
	l.sunsetSchedule(tSunset)

	return chOut
}

func send(s string, ch chan map[string]interface{}) {
	m := make(map[string]interface{})
	m[s] = "true"
	ch <- m

}

func (l *Loc) sunriseSchedule(t time.Time) {
	chOut := make(chan map[string]interface{})
	l.cronSch.AddFunc(cronFormat(t), func() {
		send("sunrise", chOut)
		l.sunriseSchedule(l.nextSunrise())
	})
}

func (l *Loc) sunsetSchedule(t time.Time) {
	chOut := make(chan map[string]interface{})
	l.cronSch.AddFunc(cronFormat(t), func() {
		send("sunset", chOut)
		l.sunsetSchedule(l.nextSunset())
	})
}

func cronFormat(t time.Time) string {
	//Second, Minute, Hour, Dom, Month, Dow
	s := fmt.Sprintf("%d,%d,%d,%d,%d,*", t.Second(), t.Minute(), t.Hour(), t.Day(), int(t.Month()))
	return s
}

func (l *Loc) nextSunset() time.Time {
	l.location.Today()
	l.location.AddDays(1)
	return l.location.Sunset()
}

func (l *Loc) nextSunrise() time.Time {
	l.location.Today()
	l.location.AddDays(1)
	return l.location.Sunrise()
}

/*Process should be..

1. Create a new instance. Loc, Time Format

			melbourne := sunrise.NewLocation(-37.81, 144.96)

			formatStr := "Jan 2 15:04:05"

			cron.New

2. Start Calculate first sunset and sunrise times.

			calc sunset & sunrise
			schedule

3. Schedule those to be sent to mqtt

			cron.add....

4. On a scheduled broadcast recalc the next sunrise or sunset.

			fmt.Println(melbourne.Sunrise().Format(formatStr))
			fmt.Println(melbourne.Sunset().Format(formatStr))

*/
