package logger

//The logger provides a method of recording strings in a text file. The files are stored in the path in folders for year and a file a day.
// Lines are formatted "L hh:mm:ss.sss - " followed by the string to be loggeed.

import (
	"fmt"
	"log"
	"os"
	"time"
)

// Logger path
type Logger struct {
	path string
}

//New Setup a new logger
func New(path string) *Logger {
	return &Logger{path: path}
}

//Log Call log each time an item needs logging.
func (l *Logger) Log(s string) {

	t := time.Now()
	year, month, day := t.Date()
	timeNow := t.Format("15:04:05.000")

	//Makesure the folder exists if not make it.
	filePath := fmt.Sprintf("%s%d/", l.path, year)
	err := os.MkdirAll(filePath, os.ModePerm)
	gotError(err)

	//Update the path to be a complete path and file name for today.
	filePath = filePath + fmt.Sprintf("%04d%02d%02d.txt", year, int(month), day)

	//Format the log line
	line := fmt.Sprintf("L %s - %s\n", timeNow, s)

	//Open the file if it does not exist create it.
	mode := os.O_WRONLY | os.O_APPEND | os.O_CREATE
	fs, err := os.OpenFile(filePath, mode, os.ModePerm)
	gotError(err)
	defer fs.Close()

	_, err = fs.WriteString(line)
	gotError(err)
}

func gotError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
