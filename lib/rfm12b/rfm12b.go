package rfm12b

//RFM12b provides a serial driver for handling Jeelabs RFM12b demo sketch.

import (
	"bufio"
	"errors"
	"log"
	"strconv"
	"strings"

	"github.com/chimera/rs232"
	"github.com/danward79/SomethingInTheBack/lib/logger"
)

//Rfm12b configuration data
type Rfm12b struct {
	portName   string
	baud       uint32
	loggerPath string
	chOut      chan []byte
	ChIn       chan interface{}
	device     *rs232.Port
	logger     *logger.Logger
}

//New instance, portName, baud, LoggerPath
func New(pN string, b uint32, p string) *Rfm12b {
	return &Rfm12b{
		portName: pN, baud: b, loggerPath: p,
		chOut: make(chan []byte), ChIn: make(chan interface{}),
		logger: logger.New(p),
	}
}

//open Starts an instance of RFM12b listening to the specified port and outputs on the channel provided.
func (r *Rfm12b) read() {

	dev, err := rs232.Open(r.portName, rs232.Options{BitRate: r.baud, DataBits: 8, StopBits: 1})
	gotError(err)
	r.device = dev
	defer dev.Close()

	lineScanner := bufio.NewScanner(r.device)
	for lineScanner.Scan() {

		var out []byte
		line := lineScanner.Text()
		oa := strings.Split(line, ` `)

		// If Logging path is proved Log output to logger
		if r.loggerPath != "" {
			r.logger.Log(line)
		}

		//FIXME: RFM2PI Compatibility
		//If msgs are valid pass to channel -minus the "OK"
		if oa[0] == "OK" {
			for i := 1; i < len(oa); i++ {
				v, err := strconv.ParseInt(oa[i], 10, 16)
				if err == nil {
					out = append(out, byte(v))
				}
			}
			r.chOut <- out
		}
	}
	gotError(lineScanner.Err())
}

//Open rfm12b driver
func (r *Rfm12b) Open() chan []byte {
	go r.read()

	go func(chIn chan interface{}) {
		for m := range chIn {
			r.write(m)
		}
	}(r.ChIn)

	return r.chOut
}

//Write packet to rfm12B
func (r *Rfm12b) write(d interface{}) error {

	switch value := d.(type) {
	case string:
		_, err := r.device.Write([]byte(value + "/n"))
		return err
		//TODO: Handle other types
		//case byte:
		//r.device.Write(value)
	}

	return errors.New("Invalid data type")
}

//Generic Function to catch errors
func gotError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
