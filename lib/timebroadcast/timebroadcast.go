//Package timebroadcast provides a mechanism to send out time packets in the form of 116,h,m,s,s at a defined interval.
package timebroadcast

import (
	"fmt"
	"time"
)

//New needs a period interval and returns a channel to output on.
func New(p int) chan interface{} {
	chOut := make(chan interface{})

	go func(period time.Duration, chOut chan interface{}) {
		const layout = "15,04,00"

		var tmr *time.Timer
		tmr = time.AfterFunc(period, func() {
			tmr.Reset(period)

			chOut <- fmt.Sprintf("116," + time.Now().Format(layout) + ",s")
		})

	}(time.Duration(p)*time.Second, chOut)

	return chOut
}
