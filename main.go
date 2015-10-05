package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/danward79/SomethingInTheBack/lib/decoder"
	_ "github.com/danward79/SomethingInTheBack/lib/decoder/decoders"
	"github.com/danward79/SomethingInTheBack/lib/mapper"
	"github.com/danward79/SomethingInTheBack/lib/mqttservices"
	"github.com/danward79/SomethingInTheBack/lib/rfm12b"
	"github.com/danward79/SomethingInTheBack/lib/timebroadcast"
	"github.com/danward79/SomethingInTheBack/lib/utils"
	"github.com/danward79/SomethingInTheBack/lib/wemodriver"
)

//config stores config data read from the config file.
var config map[string]string

func init() {
	//Load the configuration data into the config map
	file := flag.String("c", "", "path to config")
	flag.Parse()

	if *file == "" {
		log.Fatal("Need to specifiy config file")
	}

	config = utils.ReadConfig(*file)

}

func main() {
	jeelink := rfm12b.New(config["portName"], utils.Atoui(config["baud"]), config["logPathJeeLink"])
	wemos := wemodriver.New(config["wemoIP"], config["device"], utils.Atoi(config["timeout"]), config["logPathWemo"])

	//Both the wemo and the Jeelink output onto a channel, which is multiplexed bellow with utils.FanInArray
	chJeeLink := mapper.Map(decoder.ChannelDecode(jeelink.Open()))

	//Declare a new client, Publish incomming data
	mqttClient := mqttservices.NewClient(config["mqttServerIP"])

	//Assemble input channels to be multiplexed
	var mapListChannels []<-chan map[string]interface{}
	mapListChannels = append(mapListChannels, wemos.Start())
	mapListChannels = append(mapListChannels, chJeeLink)
	go mqttClient.PublishMap(utils.FanInArray(mapListChannels))

	//Timebroadcast and subscription, TODO: Need to work out how to manage this
	chSub := mqttClient.Subscribe("home/#")
	chTime := timebroadcast.New(utils.Atoi(config["timeBroadcastPeriod"]))

	for {
		select {
		case m := <-chSub:
			fmt.Printf("%s\t\t%s\n", m.TopicName, m.Payload)
		case m := <-chTime:
			fmt.Println("***Time Broadcast***")
			jeelink.ChIn <- m
		}
	}

}
