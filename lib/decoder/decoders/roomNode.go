package decoders

import (
	"bytes"
	"encoding/binary"
	//	"encoding/hex"
	"fmt"

	"github.com/danward79/SomethingInTheBack/lib/decoder"
)

//RoomData data structure from Arduino Jeenode Sketch
//	byte light;     // light sensor: 0..255
//	byte moved :1;  // motion detector: 0..1
//	byte humi  :7;  // humidity: 0..100
//	int temp   :10; // temperature: -500..+500 (tenths)
//	byte lobat :1;  // supply voltage dropped under 3.1V: 0..1

type RoomData struct {
	Node     byte
	Light    byte // light sensor: 0..255
	Moved    byte // motion detector: 0..1
	Humidity byte // humidity: 0..100
	Temp     int8 // temperature: -500..+500 (tenths)
	//LoBat    byte // supply voltage dropped under 3.1V: 0..1
}

func init() {
	decoder.Register[15] = Room
}

// Room decoder
func Room(msgData []byte) map[string]interface{} {
	var data RoomData
	m := make(map[string]interface{})

	if len(msgData) == 5 {

		buf := bytes.NewReader(msgData)

		e := binary.Read(buf, binary.LittleEndian, &data)

		fmt.Println(msgData, buf)
		fmt.Println(e, data)
		//fmt.Println("Node ID",data.Node, data.Light, data.Moved, data.Humidity, data.Temp, data.LoBat)

		m["nodeid"] = int(data.Node)
		m["temp"] = float64(data.Temp) / 100
		m["light"] = int((float64(data.Light) / 255) * 100)
		m["moved"] = int(data.Moved)
		m["humi"] = int((float64(data.Humidity) / 255) * 100)

	}
	return m
}
