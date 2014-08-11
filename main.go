package main

import (
	"fmt"

	"github.com/danward79/SomethingInTheBack/lib/decoder"
	_ "github.com/danward79/SomethingInTheBack/lib/decoder/decoders"
	"github.com/danward79/SomethingInTheBack/lib/mapper"
	"github.com/danward79/SomethingInTheBack/lib/mqttservices"
	"github.com/danward79/SomethingInTheBack/lib/rfm12b"
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

	//Start mqtt Broker
	go mqttservices.NewBroker(mqttBrokerIP).Run()

	//Both the wemo and the Jeelink output onto a channel, which is multiplexed bellow with fanIn
	chJeeLink := mapper.Map(decoder.ChannelDecode(jeelink.Open()))

	//Declare a new client, Publish incomming data
	mqttClient := mqttservices.NewClient(mqttBrokerIP)
	go mqttClient.PublishMap(fanIn(wemos.Start(), chJeeLink))

	//Timebroadcast
	go func() {
		for t := range timebroadcast.New(timeBroadcastPeriod) {
			jeelink.ChIn <- t
		}
	}()

	//Subscribe to all "home" topics
	for m := range mqttClient.Subscribe("home/#") {
		fmt.Printf("%s\t\t%s\n", m.TopicName, m.Payload)
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
