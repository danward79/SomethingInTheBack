package rfm12b

//RFM12b provides a serial driver for handling Jeelabs RFM12b demo sketch.

import (
	"bufio"
	"fmt"
	"log"
	"strconv"
	"strings"

	"101/lib/logger"
	"github.com/chimera/rs232"
)

//Rfm12b configuration
type Rfm12b struct {
	portName   string
	baud       uint32
	loggerPath string
}

//New instance, providing a port string such as "/dev/ttyUSB0", baud rate 57600 and if logging is required a path to store the logs in. Logging uses the an instance of logger. If logging is not required pass an empty string ""
func New(portName string, baud uint32, path string) *Rfm12b {
	return &Rfm12b{portName: portName, baud: baud, loggerPath: path}
}

//Starts an instance of RFM12b listening to the specified port and outputs on the channel provided.
func (self *Rfm12b) open(cs chan []byte) {

	logger := logger.New(self.loggerPath)

	device, err := rs232.Open(self.portName, rs232.Options{BitRate: self.baud, DataBits: 8, StopBits: 1})
	gotError(err)
	defer device.Close()

	lineScanner := bufio.NewScanner(device)
	for lineScanner.Scan() {

		var out []byte
		line := lineScanner.Text()
		oa := strings.Split(line, ` `)

		// If Logging path is proved Log output to logger
		if self.loggerPath != "" {
			logger.Log(line)
		}

		//If msgs are valid pass to channel
		if oa[0] == "OK" {
			for i := 1; i < len(oa); i++ {
				v, err := strconv.ParseInt(oa[i], 10, 16)
				if err == nil {
					out = append(out, byte(v))
				}
			}
			cs <- out
		}
	}
	gotError(lineScanner.Err())
}

//Start logger
func (self *Rfm12b) Start() chan []byte {
	chOut := make(chan []byte)
	go self.open(chOut)
	return chOut
}

//Generic Function to catch errors
func gotError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

//Write packet to rfm12B
func (self *Rfm12b) Write() {
	device, err := rs232.Open(self.portName, rs232.Options{BitRate: self.baud, DataBits: 8, StopBits: 1})
	gotError(err)
	defer device.Close()
	fmt.Println(device)
}
