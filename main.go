package main

import (
	"flag"
	"log"
	"sync"

	"github.com/danward79/SomethingInTheBack/lib/decoder"
	_ "github.com/danward79/SomethingInTheBack/lib/decoder/decoders"
	"github.com/danward79/SomethingInTheBack/lib/logreplay"
	"github.com/danward79/SomethingInTheBack/lib/mapper"
	"github.com/danward79/SomethingInTheBack/lib/rfm12b"
	"github.com/danward79/SomethingInTheBack/lib/timebroadcast"
	"github.com/danward79/SomethingInTheBack/lib/utils"
	"github.com/danward79/SomethingInTheBack/lib/wemodriver"
	"github.com/danward79/mqttservices"
	proto "github.com/huin/mqtt"
)

//config stores config data read from the config file.
var config map[string]string

func init() {
	log.Println("SomeThingInTheBack - Started")

	//Load the configuration data into the config map
	file := flag.String("c", "", "path to config")
	flag.Parse()

	if *file == "" {
		log.Fatal("Need to specifiy config file")
	}

	config = utils.ReadConfig(*file)
}

func main() {

	//Waitgroup for go routines
	var wg sync.WaitGroup

	//Assemble input channels to be multiplexed
	var mapListChannels []<-chan map[string]interface{}
	chTime := make(chan interface{})
	var jeelink = &rfm12b.Rfm12b{ChIn: make(chan interface{})}

	if config["portName"] != "" {
		jeelink = rfm12b.New(config["portName"], utils.Atoui(config["baud"]), config["logPathJeeLink"])
		chJeeLink := mapper.Map(decoder.ChannelDecode(jeelink.Open()))
		mapListChannels = append(mapListChannels, chJeeLink)
		chTime = timebroadcast.New(utils.Atoi(config["timeBroadcastPeriod"]))

		broadcastTime(chTime, jeelink.ChIn, &wg)
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

	//Subscribe and log events
	subscriptions("home/#", mqttClient, &wg)

	//Wait for go routines to return
	wg.Wait()
}

//broadcastTime - pass time to the serial device
func broadcastTime(chTime chan interface{}, chLink chan interface{}, wg *sync.WaitGroup) {
	wg.Add(1)

	go func() {
		defer wg.Done()

		for {
			chLink <- <-chTime
		}
	}()
}

//subscriptions - Log subscriptions
func subscriptions(t string, c *mqttservices.MqttClient, wg *sync.WaitGroup) {
	wg.Add(1)

	go func() {
		defer wg.Done()

		chSub := c.Subscribe([]proto.TopicQos{{
			Topic: t,
			Qos:   proto.QosAtMostOnce,
		}})

		for {
			m := <-chSub
			log.Printf("%s\t\t%s\n", m.TopicName, m.Payload)
		}
	}()
}
