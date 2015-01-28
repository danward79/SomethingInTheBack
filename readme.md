SomethingInTheBack
==================

This repo is a further exploration into GO for me.

I was / am following the [Housemon](https://github.com/jcw/housemon) project closely. The work there is truly cool, but being a project of that nature the direction changes frequently and this meant that for me progress was too slow. So rather than get frustrated. I thought I would bugger off and learn a bit of GO. Being a practical person I wanted to learn with building a real thing. So.

This repo is an attempt of a backend for some JeeNode and Wemo devices I have. There is a large amount of inspiration from JCWs work. So don't be surprised when this seem familiar.

###Currently it has the following functionality.

- MQTT Broker which manages subscriptions and published data. Uses the MQTT Services library, which uses [Jeffallen's](https://github.com/jeffallen/mqtt) and [Huin's](https://github.com/huin/mqtt) libraries. Equally an alternative MQTT broker could be used such as [Mosquitto](http://mosquitto.org)
- Two-way Serial interface to JeeLink so that RF messages can be received which use the Arduino RFM12b driver from Jeelabs. Uses the serial library from [Chimera](https://github.com/chimera/rs232), outputs a channel of bytes.
- Wemo interface which subscribes to any discovered devices. So that state changes are received. This uses a version of [Savaki's library](https://github.com/savaki/go.wemo), to which I added a subscription service. Which is [here](https://github.com/danward79/go.wemo). I am waiting for my [pull](https://github.com/savaki/go.wemo/pull/1) request. This outputs a channel of map[string]interface{}.
- Logging device which can be passed to the Wemo and RFM12b drivers above to log the incoming data. This produces exactly the same output as JCW's logger. It is called with a simple string.
- RFM12b packets are decoded by the decoder and it's decoders. This outputs a channel of map[string]interface{}.
- These are routed through a mapper which adds location info. This outputs a channel of map[string]interface{}.
- Data is then published using the MQTT service. Which takes a channel of map[string]interface{}.
- Time packet transmission to keep displays up to date.
- Sunrise and Sunset MQTT events at -2, -1, 0, +1 hours. So an MQTT event at 2 hours, 1 hour before, at sunrise or sunset and one hour after.
- All config data is now stored in config.txt
- Added ability to replay a RFM12B log file. See below. This mimicks a RFM12b demo receiver.

```
chJeeLink := mapper.Map(decoder.ChannelDecode(logreplay.Replay("./Logs/RFM12b/2014/20140810.txt")))
```

**Using Channels** makes passing data around very easy. With a simple multiplexer and use of interfaces it is possible to push data in a similar manner to the publisher.

```
//Both the wemo and the Jeelink output onto a channel, which is multiplexed bellow with fanIn then published
chJeeLink := mapper.Map(decoder.ChannelDecode(jeelink.Start()))
chWemo := wemos.Start()

//Declare a new client, Publish incoming data
mqttClient := mqttservices.NewClient(mqttBrokerIP)
go mqttClient.PublishMap(fanIn(chWemo, chJeeLink))
```

###What's Next?

I don't know how far I will take this, but here is a list of things in no particular order that need doing.

- ~~Serial driver needs to be converted to be two way. The main reason for this is...~~
- ~~I need to make time transmissions to keep my display up to date. This is done elsewhere at the moment.~~
- ~~Time based task scheduler so that...~~
- ~~Sunrise and Set events can be sent~~
- ~~Add configuration file~~
- ~~As I am going away soon, a replay service for the sensor logs would be useful so I can fiddle on the plane! (Now I know a good reason JCW did it!)~~
- ~~Database needs choosing and a ...~~ I think I'll use influxdb
- ~~method of adding data decided on. Subscribe to all events and use MQTT or hook in earlier in the chain?~~ I'm going to subscribe via MQTT
- TODO: Method to allow command injection to Wemo, probably use format

```
/home/instruction/room/device/command value

e.g. /home/instruction/lounge/lamp/state false
```

- TODO: Weather forecast subscription, probably using Yahoo weather.
- TODO: Sonos subscription
- TODO: Method to allow command injection to Sonos
- TODO: How would rules be implemented?
