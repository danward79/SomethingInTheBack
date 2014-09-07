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
	"github.com/danward79/SomethingInTheBack/lib/utils"
	"github.com/danward79/SomethingInTheBack/lib/wemodriver"
)

//config stores config data read from the config file.
var config map[string]string

func init() {
	//Load the configuration data into the config map
	config = utils.ReadConfig("./config.cfg")
	//Start mqtt Broker
	go mqttservices.NewBroker(config["mqttBrokerIP"]).Run()
}

func main() {
	jeelink := rfm12b.New(config["portName"], utils.Atoui(config["baud"]), config["logPathJeeLink"])
	wemos := wemodriver.New(config["wemoIP"], config["device"], utils.Atoi(config["timeout"]), config["logPathWemo"])
	melbourne := sunriseset.New(-37.81, 144.96)

	//Both the wemo and the Jeelink output onto a channel, which is multiplexed below with fanIn
	chJeeLink := mapper.Map(decoder.ChannelDecode(jeelink.Open()))

	//Declare a new client, Publish incomming data
	mqttClient := mqttservices.NewClient(config["mqttBrokerIP"])

	//assemble input channels to be multiplexed
	var mapListChannels []<-chan map[string]interface{}
	mapListChannels = append(mapListChannels, wemos.Start())
	mapListChannels = append(mapListChannels, chJeeLink)
	mapListChannels = append(mapListChannels, melbourne.Start())
	go mqttClient.PublishMap(utils.FanInArray(mapListChannels))

	//Timebroadcast and subscription, TODO: Need to work out how to manage this
	chSub := mqttClient.Subscribe("home/#")
	chTime := timebroadcast.New(utils.Atoi(config["timeBroadcastPeriod"]))

	for {
		select {
		case m := <-chSub:
			fmt.Printf("%s\t\t%s\n", m.TopicName, m.Payload)
		case m := <-chTime:
			jeelink.ChIn <- m
		}
	}

}
