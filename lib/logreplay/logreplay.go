package logreplay

/*
logreplay is used for replaying log files into the system,
which have previously been logged and stored in log txt files.

Log lines take the format below.

L 00:04:00.735 - OK 16 23 62 3 81 0 77
L 00:04:04.748 - OK 11 45 7 255
L 00:04:09.416 - OK 17 1 178 0 194 146 1 0 160 10

We need to remove the content that preceeds the OK then feed the lines to the
mapper, from which point they are handled in the same manner as RFM12B demo lines.
*/

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

//Replay reads log files and outputs a chan of bytes in the same format as the RFM12B demo
func Replay(filePath string) chan []byte {
	fmt.Println("Replaying", filePath)

	chOut := make(chan []byte)

	go read(filePath, chOut)

	return chOut
}

func read(filePath string, chOut chan []byte) {
	//Open the file
	mode := os.O_RDONLY
	fs, err := os.OpenFile(filePath, mode, os.ModePerm)
	gotError(err)
	defer fs.Close()

	//Read file lines using scanner, then split the lines and only output the byte array
	scanner := bufio.NewScanner(fs)

	for scanner.Scan() {
		var out []byte

		line := scanner.Text()
		oa := strings.Split(line[17:], ` `)

		//If msgs are valid pass to channel -minus the "OK"
		if oa[0] == "OK" {
			for i := 1; i < len(oa); i++ {
				v, err := strconv.ParseInt(oa[i], 10, 16)
				if err == nil {
					out = append(out, byte(v))
				}
			}
			chOut <- out
		}
		//Wait an arbitary time to make data a little more real.
		time.Sleep(10 * time.Second)
	}

}

//gotError helps handle errors
func gotError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
