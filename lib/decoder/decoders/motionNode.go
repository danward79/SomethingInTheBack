package decoders

import (
	"bytes"
	"encoding/binary"

	"github.com/danward79/SomethingInTheBack/lib/decoder"
)

// MotionData Structure
type MotionData struct {
	Node  byte
	Light byte // light sensor: 0..255
	Moved byte // motion detector: 0..1
	//LoBat byte // supply voltage dropped under 3.1V: 0..1 //TODO: Add LoBat Handling.
}

func init() {
	decoder.Register[15] = Motion
	decoder.Register[19] = Motion
}

// Motion decoder
func Motion(msgData []byte) map[string]interface{} {
	var data MotionData
	m := make(map[string]interface{})

	if len(msgData) >= 3 {

		buf := bytes.NewReader(msgData)

		binary.Read(buf, binary.LittleEndian, &data)

		m["nodeid"] = int(data.Node)
		m["light"] = int((float64(data.Light) / 255) * 100)
		m["moved"] = int(data.Moved)
		//	m["lobat"] = int(data.LoBat)

	}
	return m
}
