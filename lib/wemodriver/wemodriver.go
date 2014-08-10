package wemodriver

//Wemo driver to interface to wemo event subscriptions

import (
	"fmt"
	"github.com/danward79/SomethingInTheBack/lib/logger"
	"github.com/danward79/go.wemo"
	"log"
	"time"
)

//Wemos structure
type Wemos struct {
	listenerAddress string
	interfaceDevice string
	timeout         int
	loggerPath      string
}

//New create a new device
func New(a string, i string, t int, p string) *Wemos {
	return &Wemos{listenerAddress: a, interfaceDevice: i, timeout: t, loggerPath: p}
}

//Start listening
func (w *Wemos) Start() chan map[string]interface{} {

	logger := logger.New(w.loggerPath)

	api, _ := wemo.NewByInterface(w.interfaceDevice)

	devices, _ := api.DiscoverAll(3 * time.Second)

	subscriptions := make(map[string]*wemo.SubscriptionInfo)

	for _, device := range devices {
		_, err := device.ManageSubscription(w.listenerAddress, w.timeout, subscriptions)
		if err != 200 {
			log.Println("Wemo Initial Error Subscribing: ", err)
		}
	}

	chIn := make(chan wemo.SubscriptionEvent)
	chOut := make(chan map[string]interface{})

	go wemo.Listener(w.listenerAddress, chIn)

	go func(ch chan wemo.SubscriptionEvent, chOut chan map[string]interface{}) {

		for i := range ch {
			m := make(map[string]interface{})
			if _, ok := subscriptions[i.Sid]; ok {
				subscriptions[i.Sid].State = i.State

				// If Logging path is proved Log output to logger
				if w.loggerPath != "" {
					line := fmt.Sprintf("%s %t %s %s", subscriptions[i.Sid].Host, subscriptions[i.Sid].State, subscriptions[i.Sid].DeviceInfo.FriendlyName, i.Sid)
					logger.Log(line)
				}

				m["nodeid"] = subscriptions[i.Sid].Host
				m["state"] = subscriptions[i.Sid].State
				m["location"] = subscriptions[i.Sid].DeviceInfo.FriendlyName

			} else {
				log.Println("Wemo SID Does'nt exist, ", i.Sid)
			}
			chOut <- m
		}

	}(chIn, chOut)

	return chOut
}
