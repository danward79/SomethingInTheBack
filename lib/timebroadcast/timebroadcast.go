//Package timebroadcast provides a mechanism to send out time packets in the form of 116,h,m,s,s at a defined interval.
package timebroadcast

import (
	"fmt"
	"log"
	"time"
)

//New needs a period interval and returns a channel to output on.
func New(p int) chan interface{} {
	chOut := make(chan interface{})

	go func(period time.Duration, chOut chan interface{}) {
		const layout = "15,04,00"

		t := time.NewTicker(period)
		for _ = range t.C {
<<<<<<< Updated upstream
			chOut <- fmt.Sprint("116," + time.Now().Format(layout) + ",s")
			log.Println("Time Broadcast")
		}
=======
			log.Println("***Time Broadcast")
			chOut <- fmt.Sprint("116," + time.Now().Format(layout) + ",s")
		}

		/*
			var tmr *time.Timer
			tmr = time.AfterFunc(period, func() {
				tmr.Reset(period)
>>>>>>> Stashed changes

				chOut <- fmt.Sprint("116," + time.Now().Format(layout) + ",s")
			})
		*/
	}(time.Duration(p)*time.Second, chOut)

	return chOut
}
