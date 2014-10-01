//Package influxdb provides a method of interfacing to an influxdb database instance.
package influxdb

import (
	"fmt"
	proto "github.com/huin/mqtt"
	"net"
)

//Influxdb struct
type Influxdb struct {
	host string
	//	ChIn chan *proto.Publish
}

//New returns a new instance of an influxdb
func New(h string) *Influxdb {
	return &Influxdb{host: h} //, ChIn: make(chan *proto.Publish)}
}

//Start feeds channel to influxdb connection
func (db *Influxdb) Start(c chan *proto.Publish) {

	go func(cIn chan *proto.Publish) {
		for m := range cIn {
			fmt.Println("Inside ChIn", m)
			fmt.Printf("%s\t\t%s\n", m.TopicName, m.Payload)
		}
	}(c)
	return
}

func (db *Influxdb) connect() {
	host, err := net.ResolveUDPAddr("udp", db.host)
	if err != nil {
		return
	}

	con, err := net.DialUDP("udp", nil, host)
	if err != nil {
		return
	}
	defer con.Close()

}
