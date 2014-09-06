package main

import (
	"bufio"
	"fmt"
	"github.com/danward79/SomethingInTheBack/lib/decoder"
	_ "github.com/danward79/SomethingInTheBack/lib/decoder/decoders"
	"github.com/danward79/SomethingInTheBack/lib/mapper"
	"github.com/danward79/SomethingInTheBack/lib/mqttservices"
	"github.com/danward79/SomethingInTheBack/lib/rfm12b"
	"github.com/danward79/SomethingInTheBack/lib/sunriseset"
	"github.com/danward79/SomethingInTheBack/lib/timebroadcast"
	"github.com/danward79/SomethingInTheBack/lib/wemodriver"
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	baud uint32 = 57600
)

var config map[string]string

func init() {
	config = readConfig("./config.txt")
	//Start mqtt Broker
	go mqttservices.NewBroker(config["mqttBrokerIP"]).Run()

}

func main() {
	jeelink := rfm12b.New(config["portName"], baud, config["logPathJeeLink"])
	wemos := wemodriver.New(config["wemoIP"], config["device"], Atoi(config["timeout"]), config["logPathWemo"])
	melbourne := sunriseset.New(-37.81, 144.96)

	//Both the wemo and the Jeelink output onto a channel, which is multiplexed bellow with fanIn
	chJeeLink := mapper.Map(decoder.ChannelDecode(jeelink.Open()))

	//Declare a new client, Publish incomming data
	mqttClient := mqttservices.NewClient(config["mqttBrokerIP"])

	//assemble input channels to be multiplexed
	var mapListChannels []<-chan map[string]interface{}
	mapListChannels = append(mapListChannels, wemos.Start())
	mapListChannels = append(mapListChannels, chJeeLink)
	mapListChannels = append(mapListChannels, melbourne.Start())
	go mqttClient.PublishMap(fanInArray(mapListChannels))

	//Timebroadcast and subscription, TODO: Need to work out how to manage this
	chSub := mqttClient.Subscribe("home/#")
	chTime := timebroadcast.New(Atoi(config["timeBroadcastPeriod"]))

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

//Atoi Helper to convert string to int.
func Atoi(s string) int {
	i, e := strconv.Atoi(s)
	if e != nil {
		return 0
	}
	return i
}

//ReadConfig takes a path to a configuration file and returns a map of configuration parameters
func readConfig(path string) map[string]string {

	configMap := make(map[string]string)

	f, err := os.Open(path)
	defer f.Close()
	if err != nil {
		return nil
	}

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {

		line := scanner.Text()

		if !strings.HasPrefix(line, "//") {
			fields := strings.SplitN(scanner.Text(), "=", 2)

			configMap[strings.TrimSpace(fields[0])] = strings.TrimSpace(fields[1])
		}

	}

	return configMap
}

//gotError checks for an error
func gotError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
