package decoders

import (
	"bytes"
	"encoding/binary"
	"github.com/danward79/SomethingInTheBack/lib/decoder"
)

func init() {
	decoder.Register[16] = OutdoorClimate
}

//OutdoorData data structure
//struct {byte light; int humidity; int temperature; byte vcc; } payload;
type OutdoorData struct {
	Node     uint8
	Light    byte
	Humidity uint16
	Temp     uint16
	Battery  byte
}

//OutdoorClimate packet decoder
func OutdoorClimate(msgData []byte) map[string]interface{} {
	m := make(map[string]interface{})
	var data OutdoorData

	if len(msgData) >= 6 {

		buf := bytes.NewReader(msgData)
		_ = binary.Read(buf, binary.LittleEndian, &data)

		m["nodeid"] = int(data.Node)
		m["light"] = int((float64(data.Light) / 255) * 100)
		m["humi"] = float64(data.Humidity) / 10
		m["temp"] = float64(data.Temp) / 10
		m["battery"] = (int(data.Battery) * 20) + 1000
	}
	return m
}
