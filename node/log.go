package main

import (
	"log"
	"os"
	"time"
)

var glogfile *os.File

func LogDirSet(dir string) error {
	tm := time.Now()
	logfilename := tm.Format("20060102150405") + ".log"
	logfilepath := dir + logfilename

	fd, err := os.OpenFile(logfilepath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}

	log.SetOutput(fd)
	return nil
}
