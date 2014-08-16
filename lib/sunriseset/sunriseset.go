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

	l.location.Today()
	tSunset := l.location.Sunset()
	if time.Now().After(tSunset) {
		tSunset = l.nextSunset()
	}

	//Schedule cron
	l.scheduleNext(tSunrise, true, chOut)
	l.scheduleNext(tSunset, false, chOut)
	l.cronSch.Start()
	//TODO: Need to work out how to remove completed jobs or reschedule existing jobs.
	l.test(time.Now(), chOut)

	return chOut
}

//send a sunrise or sunset event to the output channel
func send(s string, t string, ch chan map[string]interface{}) {
	m := make(map[string]interface{})
	m["location"] = s
	m["state"] = t
	ch <- m
}

//test to make sure it is all working
func (l *Loc) test(t time.Time, ch chan map[string]interface{}) {

	t = t.Add(10 * time.Second)
	//Second, Minute, Hour, Dom, Month, Dow
	l.cronSch.AddFunc(cronFormat(t), func() {

		send("sunrise", fmt.Sprintf("%d", -1), ch)
		l.test(time.Now().Add(10*time.Second), ch)

	})

	for k, v := range l.cronSch.Entries() {
		fmt.Println(k, v)
	}

}

//schedule the next sunrise or sunset, set sunrise true for Sunrise
func (l *Loc) scheduleNext(t time.Time, sunrise bool, ch chan map[string]interface{}) {
	for i := -1; i <= 1; i++ {
		l.cronSch.AddFunc(cronFormat(t.Add(time.Duration(i)*time.Hour)), func() {
			if sunrise {
				send("sunrise", fmt.Sprintf("%d", i), ch)
				l.scheduleNext(l.nextSunrise(), true, ch)
			} else {
				send("sunset", fmt.Sprintf("%d", i), ch)
				l.scheduleNext(l.nextSunset(), false, ch)
			}
		})
	}
}

//cronFormat converts a time.Time to a cron schedule string
func cronFormat(t time.Time) string {
	//Second, Minute, Hour, Dom, Month, Dow
	s := fmt.Sprintf("%d %d %d %d %d *", t.Second(), t.Minute(), t.Hour(), t.Day(), int(t.Month()))
	return s
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
