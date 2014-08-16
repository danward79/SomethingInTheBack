package main

import (
	"fmt"

	"github.com/danward79/SomethingInTheBack/lib/decoder"
	_ "github.com/danward79/SomethingInTheBack/lib/decoder/decoders"
	"github.com/danward79/SomethingInTheBack/lib/mapper"
	"github.com/danward79/SomethingInTheBack/lib/mqttservices"
	"github.com/danward79/SomethingInTheBack/lib/rfm12b"
	"github.com/danward79/SomethingInTheBack/lib/sunriseset"
	"github.com/danward79/SomethingInTheBack/lib/timebroadcast"
	"github.com/danward79/SomethingInTheBack/lib/wemodriver"
)

//TODO: Set up config file
//TODO: Set up command line parsing
const (
	//"/dev/ttyUSB0" rPi USB, "/dev/ttyAMA0" rPi Header, "/dev/tty.usbserial-A1014KGL" Mac
	portName            string = "/dev/tty.usbserial-A1014KGL" //Mac
	baud                uint32 = 57600
	logPathJeeLink      string = "./Logs/RFM12b/"
	wemoIP              string = "192.168.0.6:6767"
	device              string = "en0"
	timeout             int    = 600
	logPathWemo         string = "./Logs/Wemo/"
	mqttBrokerIP        string = ":1883" //"test.mosquitto.org:1883"
	timeBroadcastPeriod int    = 300
)

func main() {
	jeelink := rfm12b.New(portName, baud, logPathJeeLink)
	wemos := wemodriver.New(wemoIP, device, timeout, logPathWemo)
	melbourne := sunriseset.New(-37.81, 144.96)

	//Start mqtt Broker
	go mqttservices.NewBroker(mqttBrokerIP).Run()

	//Both the wemo and the Jeelink output onto a channel, which is multiplexed bellow with fanIn
	chJeeLink := mapper.Map(decoder.ChannelDecode(jeelink.Open()))

	//Declare a new client, Publish incomming data
	mqttClient := mqttservices.NewClient(mqttBrokerIP)

	//assemble input channels to be multiplexed
	var mapListChannels []<-chan map[string]interface{}
	mapListChannels = append(mapListChannels, wemos.Start())
	mapListChannels = append(mapListChannels, chJeeLink)
	mapListChannels = append(mapListChannels, melbourne.Start())
	go mqttClient.PublishMap(fanInArray(mapListChannels))

	//TODO: Need to work out how to manage this
	//Timebroadcast and subscription
	chSub := mqttClient.Subscribe("home/#")
	chTime := timebroadcast.New(timeBroadcastPeriod)

	for {
		select {
		case m := <-chSub:
			fmt.Printf("%s\t\t%s\n", m.TopicName, m.Payload)
		case m := <-chTime:
			jeelink.ChIn <- m
		}
	}

}

//TODO: Move to a seperate library?
//fanin Multiplex two channels to a single output, this code was pinched from a google presentation ;-)
func fanIn(input1 <-chan map[string]interface{}, input2 <-chan map[string]interface{}) chan map[string]interface{} {
	c := make(chan map[string]interface{})

	go func() {
		for {
			c <- <-input1
		}
	}()

	go func() {
		for {
			c <- <-input2
		}
	}()

	return c
}

//fanInArray is a version of fanIn which takes an array of chan map[string]interface{} making fanIn an expandable input multiplexer
func fanInArray(inputChannels []<-chan map[string]interface{}) chan map[string]interface{} {
	c := make(chan map[string]interface{})

	for i := range inputChannels {
		go func(chIn <-chan map[string]interface{}) {
			for {
				c <- <-chIn
			}
		}(inputChannels[i])
	}
	return c
}
