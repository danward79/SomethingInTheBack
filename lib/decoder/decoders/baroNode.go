// Package decoders for the "BMP085demo.ino" sketch as in: http://github.com/jcw/jeelib/tree/master/examples/Ports/bmp085demo/bmp085demo.ino
// Decoder Registers as "Node-Bmp085".
//
// Tweaked version of TheDistractors Code. Added Light level and battery voltage packets.
package decoders

import (
	"bytes"
	"encoding/binary"

	"github.com/danward79/SomethingInTheBack/lib/decoder"
)

//BaroData Data structure
type BaroData struct {
	Node    uint8
	Light   byte
	Temp    uint16
	Press   uint32
	Battery uint16
}

func init() {
	decoder.Register[17] = BaroNode
}

//BaroNode carries out conversion
func BaroNode(msgData []byte) map[string]interface{} {
	m := make(map[string]interface{})

	if len(msgData) >= 8 {
		buf := bytes.NewReader(msgData)
		var data BaroData
		_ = binary.Read(buf, binary.LittleEndian, &data)

		m["nodeid"] = int(data.Node)
		m["light"] = int(data.Light)
		m["temp"] = int(data.Temp)
		m["pressure"] = int(data.Press)
		m["battery"] = int(data.Battery)
	}
	return m
}
