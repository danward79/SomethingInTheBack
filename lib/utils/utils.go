package utils

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

//Atoi Helper to convert string to int.
func Atoi(s string) int {
	i, e := strconv.Atoi(s)
	if e != nil {
		return 0
	}
	return i
}

//Atoui helper to convert string to uint32
func Atoui(s string) uint32 {
	i, e := strconv.ParseUint(s, 10, 32)
	if e != nil {
		return 0
	}
	return uint32(i)
}

//ReadConfig takes a path to a configuration file and returns a map of configuration parameters
func ReadConfig(path string) map[string]string {

	configMap := make(map[string]string)

	f, err := os.Open(path)
	defer f.Close()
	if err != nil {
		return nil
	}

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {

		line := scanner.Text()

		if !strings.HasPrefix(line, "//") {
			fields := strings.SplitN(scanner.Text(), "=", 2)

			configMap[strings.TrimSpace(fields[0])] = strings.TrimSpace(fields[1])
		}

	}

	return configMap
}

//GotError checks for an error
func GotError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
