package decoders

import (
	"bytes"
	"encoding/binary"

	"101/lib/decoder"
)

//LCDData data structure
type LCDData struct {
	Node  uint8
	Temp  uint16
	Light byte
}

func init() {
	decoder.Register[11] = EmonLcd
}

// EmonLcd decoder
func EmonLcd(msgData []byte) map[string]interface{} {
	var data LCDData
	m := make(map[string]interface{})

	if len(msgData) == 4 {
		buf := bytes.NewReader(msgData)

		_ = binary.Read(buf, binary.LittleEndian, &data)

		m["nodeid"] = int(data.Node)
		m["temp"] = int(data.Temp)
		m["light"] = int(255 - data.Light)

	}
	return m
}
