package main

import (
	"io"
	"log"
	"os"
)

var (
	// LogW warning logger
	LogW *log.Logger
	// LogI info logger
	LogI *log.Logger
	// LogE error logger
	LogE *log.Logger
)

func init() {
	// Setp logger
	file, _ := os.OpenFile("./logs.log", os.O_APPEND|os.O_CREATE, 0666)
	mw := io.MultiWriter(os.Stdout, file)

	LogW = log.New(mw, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	LogI = log.New(mw, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	LogE = log.New(mw, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}
