package main

import (
	"fmt"

	"github.com/danward79/SomethingInTheBack/lib/decoder"
	_ "github.com/danward79/SomethingInTheBack/lib/decoder/decoders"
	"github.com/danward79/SomethingInTheBack/lib/logreplay"
	"github.com/danward79/SomethingInTheBack/lib/mapper"
	"github.com/danward79/SomethingInTheBack/lib/mqttservices"
	"github.com/danward79/SomethingInTheBack/lib/utils"
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
	//Both the wemo and the Jeelink output onto a channel, which is multiplexed below with fanIn
	chJeeLink := mapper.Map(decoder.ChannelDecode(logreplay.Replay("./Logs/RFM12b/2014/20140810.txt")))

	//Declare a new client, Publish incomming data
	mqttClient := mqttservices.NewClient("localhost:1883")

	//Assemble input channels to be multiplexed
	var mapListChannels []<-chan map[string]interface{}
	mapListChannels = append(mapListChannels, chJeeLink)
	go mqttClient.PublishMap(utils.FanInArray(mapListChannels))

	//Timebroadcast and subscription
	chSub := mqttClient.Subscribe("home/#")

	for {
		select {
		case m := <-chSub:
			fmt.Printf("%s\t\t%s\n", m.TopicName, m.Payload)
		}
	}
}
