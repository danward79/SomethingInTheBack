package utils

//FanIn Multiplex two channels to a single output, this code was pinched from a google presentation ;-)
func FanIn(input1 <-chan map[string]interface{}, input2 <-chan map[string]interface{}) chan map[string]interface{} {
	c := make(chan map[string]interface{})

	go func() {
		for {
			c <- <-input1
		}
	}()

	go func() {
		for {
			c <- <-input2
		}
	}()

	return c
}

//FanInArray is a version of fanIn which takes an array of chan map[string]interface{} making fanIn an expandable input multiplexer
func FanInArray(inputChannels []<-chan map[string]interface{}) chan map[string]interface{} {
	c := make(chan map[string]interface{})

	for i := range inputChannels {
		go func(chIn <-chan map[string]interface{}) {
			for {
				c <- <-chIn
			}
		}(inputChannels[i])
	}
	return c
}
