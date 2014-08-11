package main

//TODO: Move to a seperate library?
//fanin Multiplex two channels to a single output, this code was pinched from a google presentation ;-)
func fanIn(input1 <-chan map[string]interface{}, input2 <-chan map[string]interface{}) chan map[string]interface{} {
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
