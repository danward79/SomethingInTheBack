package decoders

import (
	"github.com/danward79/SomethingInTheBack/lib/decoder"
)

func init() {
	decoder.Register[16] = OutdoorClimate
}

//OutdoorClimate packet decoder
//struct {byte light; int humidity; int temperature; byte vcc; } payload;
func OutdoorClimate(msgData []byte) map[string]interface{} {
	m := make(map[string]interface{})

	if len(msgData) >= 6 {

		m["nodeid"] = int(msgData[0])
		m["light"] = int(msgData[1])
		m["humi"] = int(msgData[2]) + (256 * int(msgData[3]))
		m["temp"] = int(msgData[4]) + (256 * int(msgData[5]))
		m["battery"] = (int(msgData[6]) * 20) + 1000
	}
	return m
}
