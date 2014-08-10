package rfm12b

import (
	"bufio"
	"fmt"
	"github.com/chimera/rs232"
	"github.com/danward79/SomethingInTheBack/lib/logger"
	"strconv"
	"strings"
)

//Rfm12b2 configuration
type Rfm12b2 struct {
	portName   string
	baud       uint32
	loggerPath string
	chOut      chan []byte
	ChIn       chan []byte
	device     *rs232.Port
	logger     *logger.Logger
}

//New2 instance, portName, baud, LoggerPath
func New2(pN string, b uint32, p string) *Rfm12b2 {
	return &Rfm12b2{
		portName: pN, baud: b, loggerPath: p,
		chOut: make(chan []byte), ChIn: make(chan []byte),
		logger: logger.New(p),
	}
}

//open Starts an instance of RFM12b listening to the specified port and outputs on the channel provided.
func (r *Rfm12b2) read() {

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

		//If msgs are valid pass to channel
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
func (r *Rfm12b2) Open() chan []byte {
	go r.read()

	go func(chIn chan []byte) {
		for m := range chIn {
			r.write(m)
		}
	}(r.ChIn)

	return r.chOut
}

//Write packet to rfm12B
func (r *Rfm12b2) write(d []byte) {
	fmt.Println("Write:", string(d))
}
