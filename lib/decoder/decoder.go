package decoder

// Provides an interface to rfm12b serial driver. Uses NodeId to identify which decoder to use to decode the RFM12B messages.

//Register of all known decoders.
var Register = map[byte]func([]byte) map[string]interface{}{}

//Decode function uses Nodeid to find correct decoder in registry returns a map of key, values
func decode(chIn chan []byte, chOut chan map[string]interface{}) {
	var m map[string]interface{}

	for v := range chIn {
		if len(v) > 0 {
			if d, ok := Register[v[0]]; ok {
				m = d(v)
			}
			chOut <- m
		}
	}
}

//ChannelDecode function uses Nodeid to find correct decoder in registry returns a map of key, values to the returned channel
func ChannelDecode(chIn chan []byte) chan map[string]interface{} {
	chOut := make(chan map[string]interface{})
	go decode(chIn, chOut)
	return chOut
}
