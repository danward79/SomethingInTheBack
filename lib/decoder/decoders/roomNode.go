package decoders

import (
	"bytes"
	"encoding/binary"

	"github.com/danward79/SomethingInTheBack/lib/decoder"
)

// RoomData Structure
type RoomData struct {
	Node     byte
	Light    byte // light sensor: 0..255
	Moved    byte // motion detector: 0..1
	Humidity byte // humidity: 0..100
	Temp     int8 // temperature: -500..+500 (tenths)
	//LoBat    byte // supply voltage dropped under 3.1V: 0..1 //TODO: Add LoBat Handling.
}

func init() {
	decoder.Register[99] = Room
}

// Room decoder
func Room(msgData []byte) map[string]interface{} {
	var data RoomData
	m := make(map[string]interface{})

	if len(msgData) == 5 {

		buf := bytes.NewReader(msgData)

		binary.Read(buf, binary.LittleEndian, &data)

		m["nodeid"] = int(data.Node)
		m["temp"] = float64(data.Temp) / 100
		m["light"] = int((float64(data.Light) / 255) * 100)
		m["moved"] = int(data.Moved)
		m["humi"] = int((float64(data.Humidity) / 255) * 100)

	}
	return m
}
