package mapper

//Mapper maps device messages to locations and preps for mqtt
//Map inputs mean it should be simple to make device agnostic.
//Thinking of RFM12b Jeenode devices, but also Wemo events and Sonos

//Register of all known devices(nodeid) to Locations and Type.
var Register = map[int]string{}

//Add devices with nodeids to Register at start. Not necessary for Wemo devices.
func init() {
	Register[11] = "Bedroom"
	Register[15] = "Lab"
	Register[16] = "Balcony"
	Register[17] = "Lounge"
	Register[19] = "LabBB"
}

//Map takes a map of strings and ensure that there is a location attribute.
func Map(chIn chan map[string]interface{}) chan map[string]interface{} {
	chOut := make(chan map[string]interface{})

	go func(chIn chan map[string]interface{}, chOut chan map[string]interface{}) {

		for c := range chIn {
			m := make(map[string]interface{})
			for k, v := range c {
				m[k] = v
			}
			if _, ok := m["location"]; !ok {
				if m["nodeid"] != nil {
					m["location"] = Register[m["nodeid"].(int)]
				}
			}
			chOut <- m
		}
	}(chIn, chOut)

	return chOut
}
