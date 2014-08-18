// Package sunriseset provides sunrise and sunset events for the location provided.
package sunriseset

import (
	"fmt"
	"github.com/danward79/cron"
	"github.com/danward79/sunrise"
	"time"
)

//Loc stores data about the location of sunrise and set required.
type Loc struct {
	location *sunrise.Location
	cronSch  *cron.Cron
}

//New returns a new location
func New(latitude float64, longitude float64) *Loc {
	return &Loc{location: sunrise.NewLocation(latitude, longitude), cronSch: cron.New()}
}

//Start the process
func (l *Loc) Start() chan map[string]interface{} {
	chOut := make(chan map[string]interface{})
	l.cronSch.Start()

	//Is it before or after todays sunrise/sunset?
	l.location.Today()
	tSunrise := l.location.Sunrise()
	tSunset := l.location.Sunset()
	if time.Now().After(tSunrise) {
		tSunrise = l.nextSunrise()
	}
	if time.Now().After(tSunset) {
		tSunset = l.nextSunset()
	}

	//Schedule cron
	l.scheduleNext(tSunrise, true, chOut)
	l.scheduleNext(tSunset, false, chOut)

	return chOut
}

//send a sunrise or sunset event to the output channel
func send(s string, t string, ch chan map[string]interface{}) {
	m := make(map[string]interface{})
	m["location"] = s
	m["state"] = t
	ch <- m
}

//schedule the next sunrise or sunset, set sunrise true for Sunrise
func (l *Loc) scheduleNext(t time.Time, rise bool, ch chan map[string]interface{}) {
	for i := -2; i <= 1; i++ {
		l.schedule(t, i, rise, ch)
	}
}

func (l *Loc) schedule(t time.Time, i int, rise bool, ch chan map[string]interface{}) {
	func(t time.Time, i int, rise bool, ch chan map[string]interface{}) {
		var e *cron.Entry
		e, _ = l.cronSch.AddFunc(cronFormat(t.Add(time.Duration(i)*time.Hour)), func() {
			l.cronSch.Remove(e)
			if rise {
				send("sunrise", fmt.Sprintf("%d", i), ch)
				l.schedule(l.nextSunrise(), i, true, ch)
			} else {
				send("sunset", fmt.Sprintf("%d", i), ch)
				l.schedule(l.nextSunset(), i, false, ch)
			}
		})
	}(t, i, rise, ch)
}

//cronFormat converts a time.Time to a cron schedule string Second, Minute, Hour, Dom, Month, Dow
func cronFormat(t time.Time) string {
	return fmt.Sprintf("%d %d %d %d %d *", t.Second(), t.Minute(), t.Hour(), t.Day(), int(t.Month()))
}

//nextSunset returns the time of the next sunset
func (l *Loc) nextSunset() time.Time {
	l.location.Today()
	l.location.AddDays(1)
	return l.location.Sunset()
}

//nextSunrise returns the time of the next sunrise
func (l *Loc) nextSunrise() time.Time {
	l.location.Today()
	l.location.AddDays(1)
	return l.location.Sunrise()
}
