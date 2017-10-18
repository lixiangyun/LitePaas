package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
)

type MasterConfig struct {
	MasterName string   `json:"mastername"`
	Address    string   `json:"address"`
	Port       int      `json:"port"`
	MasterAddr []string `json:"othermasteraddr"`
	LogDir     string   `json:"logdir"`
}

func LoadCfg(file string) (*MasterConfig, error) {
	fd, err := os.Open(file)
	if err != nil {
		return nil, errors.New("open " + file + " failed!")
	}
	defer fd.Close()
	buf := make([]byte, 0)

	for {
		var tmpbuf [128]byte
		cnt, err := fd.Read(tmpbuf[0:])
		if cnt > 0 {
			buf = append(buf, tmpbuf[0:cnt]...)
		}
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return nil, err
			}
		}
	}
	cfg := new(MasterConfig)
	err = json.Unmarshal(buf, cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func SaveCfg(cfg *MasterConfig, file string) error {
	var fd *os.File
	var err error

	for {
		fd, err = os.Create(file)
		if err != nil {
			if os.IsExist(err) {
				err = os.Remove(file)
				if nil != err {
					return err
				}
			} else {
				return err
			}
		} else {
			break
		}
	}
	defer fd.Close()

	buf, err := json.Marshal(cfg)
	if err != nil {
		return err
	}
	begin := 0
	end := len(buf)

	for {
		cnt, err := fd.Write(buf[begin:end])
		if cnt > 0 {
			begin += cnt
		}
		if err != nil {
			return err
		}
		if begin >= end {
			break
		}
	}
	return nil
}

func ShowCfg(cfg *MasterConfig) string {

	str := fmt.Sprintf("Master     : %s \r\n", cfg.MasterName)
	str += fmt.Sprintf("Addr       : %s:%d \r\n", cfg.Address, cfg.Port)

	for idx, v := range cfg.MasterAddr {
		if idx == 0 {
			str += fmt.Sprintf("MasterAddr : %s \r\n", v)
		} else {
			str += fmt.Sprintf("             %s \r\n", v)
		}
	}
	str += fmt.Sprintf("LogDir     : %s \r\n", cfg.LogDir)
	return str
}
