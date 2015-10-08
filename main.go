package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/danward79/SomethingInTheBack/lib/decoder"
	_ "github.com/danward79/SomethingInTheBack/lib/decoder/decoders"
	"github.com/danward79/SomethingInTheBack/lib/logreplay"
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

	//Assemble input channels to be multiplexed
	var mapListChannels []<-chan map[string]interface{}
	chTime := make(chan interface{})
	var jeelink = rfm12b.Rfm12b{ChIn: make(chan interface{})}

	if config["portName"] != "" {
		jeelink := rfm12b.New(config["portName"], utils.Atoui(config["baud"]), config["logPathJeeLink"])
		chJeeLink := mapper.Map(decoder.ChannelDecode(jeelink.Open()))
		mapListChannels = append(mapListChannels, chJeeLink)
		chTime = timebroadcast.New(utils.Atoi(config["timeBroadcastPeriod"]))
	}

	if config["wemoIP"] != "" {
		wemos := wemodriver.New(config["wemoIP"], config["device"], utils.Atoi(config["timeout"]), config["logPathWemo"])
		mapListChannels = append(mapListChannels, wemos.Start())
	}

	if config["replayPath"] != "" {
		chJeeLink := mapper.Map(decoder.ChannelDecode(logreplay.Replay("./Logs/RFM12b/2015/20151001.txt")))
		mapListChannels = append(mapListChannels, chJeeLink)
	}

	//Declare a new client, Publish incomming data
	mqttClient := mqttservices.NewClient(config["mqttServerIP"])

	//Both the wemo and the Jeelink output onto a channel, which is multiplexed bellow with utils.FanInArray
	go mqttClient.PublishMap(utils.FanInArray(mapListChannels))

	//Timebroadcast and subscription, TODO: Need to work out how to manage this
	chSub := mqttClient.Subscribe("home/#")

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
