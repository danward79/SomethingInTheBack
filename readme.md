SomethingInTheBack
==================

This repo is a further exploration into GO for me.

I was / am following the [Housemon](https://github.com/jcw/housemon) project closely. The work there is truely cool, but being a project of that nature the direction changes frequently and this meant that for me progress was too slow. So rather than get frustrated. I thought I would bugger off and learn a bit of GO. Being a practicle person I wanted to learn with building a real thing. So.

This repo is an attempt of a backend for some JeeNode and Wemo devices I have. There is a large amount of inspiration from JCWs work. So don't be surprised when this seem familiar.

###Currently it has the following functionality.

- MQTT Broker which manages subscriptions and published data. Uses the MQTT Services library, which uses [Jeffallen's](https://github.com/jeffallen/mqtt) and [Huin's](https://github.com/huin/mqtt) librarys
- Two-way Serial interface to JeeLink so that RF messages can be received which use the Arduino RFM12b driver from Jeelabs. Uses the serial library from [Chimera](https://github.com/chimera/rs232), outputs a channel of bytes.
- Wemo interface which subscribes to any discovered devices. So that state changes are received. This uses a version of [Savaki's library](https://github.com/savaki/go.wemo), to which I added a subscription service. Which is [here](https://github.com/danward79/go.wemo). I am waiting for my [pull](https://github.com/savaki/go.wemo/pull/1) request. This outputs a channel of map[string]interface{}.
- Logging device which can be passed to the Wemo and RFM12b drivers above to log the incomming data. This produces exactly the same output as JCW's logger. It is called with a simple string.
- RFM12b packets are decoded by the decoder and it's decoders. This outputs a channel of map[string]interface{}.
- These are routed through a mapper which adds location info. This outputs a channel of map[string]interface{}.
- Data is then published using the MQTT service. Which takes a channel of map[string]interface{}.
- Time packet transmission to keep displays up to date.

**Using Channels** makes passing data around very easy. With a simple multiplexer and use of interfaces it is possible to push data in a similar manner to the publisher.

```
//Both the wemo and the Jeelink output onto a channel, which is multiplexed bellow with fanIn then published
chJeeLink := mapper.Map(decoder.ChannelDecode(jeelink.Start()))
chWemo := wemos.Start()

//Declare a new client, Publish incomming data
mqttClient := mqttservices.NewClient(mqttBrokerIP)
go mqttClient.PublishMap(fanIn(chWemo, chJeeLink))
```

###What's Next?

I don't know how far I will take this, but here is a list of things in no particular order that need doing.

- ~~Serial driver needs to be converted to be two way. The main reason for this is...~~
- ~~I need to make time transmissions to keep my display upto date. This is done elsewhere at the moment.~~
- Database needs choosing and a method of adding data decided on. Do I subscribe to all events and use MQTT or do I hook in earlier in the chain?
- Time based task scheduler so that...
- Sunrise and Set events can be sent
- Weather forecast subscription, probably using Yahoo weather.
