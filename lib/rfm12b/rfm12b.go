package rfm12b

//RFM12b provides a serial driver for handling Jeelabs RFM12b demo sketch.

import (
	"bufio"

	"log"
	"strconv"
	"strings"

	"github.com/chimera/rs232"
	"github.com/danward79/SomethingInTheBack/lib/logger"
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
func (r *Rfm12b) open(chOut chan []byte) {

	logger := logger.New(r.loggerPath)

	device, err := rs232.Open(r.portName, rs232.Options{BitRate: r.baud, DataBits: 8, StopBits: 1})
	gotError(err)
	defer device.Close()

	lineScanner := bufio.NewScanner(device)
	for lineScanner.Scan() {

		var out []byte
		line := lineScanner.Text()
		oa := strings.Split(line, ` `)

		// If Logging path is proved Log output to logger
		if r.loggerPath != "" {
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
			chOut <- out
		}
	}
	gotError(lineScanner.Err())
}

//Start rfm12b driver
func (r *Rfm12b) Start() chan []byte {
	chOut := make(chan []byte)
	go r.open(chOut)
	return chOut
}

//Generic Function to catch errors
func gotError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
