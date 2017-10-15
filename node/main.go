package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type config struct {
	Nodename   string   `json:"nodename"`
	Address    string   `json:"address"`
	Port       int      `json:"port"`
	ServerAddr []string `json:"serveraddress"`
}

func Usage() {
	fmt.Println("Usage: <config.json>")
	os.Exit(1)
}

func ParseConfigfile(file string) *config {
	fd, err := os.Open(file)
	if err != nil {
		fmt.Println("open file failed!", err.Error())
		os.Exit(2)
	}

	buf := make([]byte, 0)

	for {
		var tmpbuf [120]byte
		cnt, err := fd.Read(tmpbuf[0:])
		if cnt > 0 {
			buf = append(buf, tmpbuf[0:cnt]...)
		}
		if err == io.EOF {
			break
		}
	}

	fd.Close()

	cfg := new(config)

	//fmt.Println("buf : ", string(buf))

	err = json.Unmarshal(buf, cfg)
	if err != nil {
		fmt.Println("json unmarshal failed!", err.Error())
		os.Exit(2)
	}

	return cfg
}

var configfile = "config.json"

var globalconfig *config

func main() {
	args := os.Args

	if len(args) > 3 {
		Usage()
	}

	if len(args) == 2 {
		configfile = args[1]
	}

	fmt.Println("load", configfile)

	globalconfig = ParseConfigfile(configfile)

	fmt.Println(globalconfig)
}
