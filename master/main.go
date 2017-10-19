package main

import (
	"fmt"
	"log"
	"os"
)

func Usage() {
	fmt.Println("Usage: <config.json>")
}

func main() {
	args := os.Args

	if len(args) > 3 {
		Usage()
		return
	}

	var cfgname string

	if len(args) == 2 {
		cfgname = args[1]
	} else {
		cfgname = "config.json"
	}

	fmt.Println("load : ", cfgname)

	cfg, err := LoadCfg(cfgname)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	err = SaveCfg(cfg, "save.json")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	cfg2, err := LoadCfg("save.json")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(ShowCfg(cfg))
	fmt.Println(ShowCfg(cfg2))

	err = LogDirSet(cfg.LogDir)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	log.Println("helloworld!")
}
