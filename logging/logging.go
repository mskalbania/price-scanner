package logging

import (
	"io"
	"log"
	"os"
)

var L *log.Logger

func init() {
	file, err := os.OpenFile("log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Error creating/opening log file - %v", err)
	}
	L = log.New(io.MultiWriter(os.Stdout, file), "", log.Ldate|log.Ltime|log.Lshortfile)
}
