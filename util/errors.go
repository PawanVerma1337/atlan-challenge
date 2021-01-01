package util

import (
	"log"
)

// CheckError utility function to check error
// and log them to terminal/logfile.
func CheckError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
