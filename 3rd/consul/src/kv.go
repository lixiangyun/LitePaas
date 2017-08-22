package kv

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	//	"fmt"
	"io"
	"net/http"
	"strings"
)

type Consul struct {
	home   string
	status bool
	index  int
}

type KvClient struct {
	cosl []Consul
}

type KvInfo struct {
	Key   string
	Value string
}

type KvData struct {
	LockInfex   int    `json:"LockIndex"`
	Key         string `json:"Key"`
	Flags       int    `json:"Flags"`
	Value       string `json:"Value"`
	CreateIndex int    `json:"CreateIndex"`
	ModifyIndex int    `json:"ModifyIndex"`
}

func readFully(conn io.ReadCloser) ([]byte, error) {
	result := bytes.NewBuffer(nil)

	var buf [512]byte

	for {
		n, err := conn.Read(buf[0:])
		result.Write(buf[0:n])

		if err != nil {
			if err == io.EOF {
				break
			} else {
				return nil, err
			}
		}
	}

	return result.Bytes(), nil
}

func ConsulRequest(method string, url string, req []byte) (rsp []byte, err error) {
	body := bytes.NewBuffer(req)

	request, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	rspon, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}

	defer rspon.Body.Close()

	rsp, err = readFully(rspon.Body)
	if err != nil {
		return nil, err
	}

	return rsp, nil
}

func NewKvClient(addr []string) (c *KvClient) {

	if len(addr) == 0 {
		addr = []string{"localhost:8500"}
	}

	consul := make([]Consul, len(addr))

	for i, v := range addr {
		consul[i].home = "http://" + v + "/v1/kv/"
		consul[i].status = true
		consul[i].index = i
	}

	return &KvClient{cosl: consul}
}

func getConsul(c *KvClient) *Consul {
	for _, v := range c.cosl {
		if v.status == true {
			return &v
		}
	}

	return nil
}

func setConsulStatus(index int, status bool, c *KvClient) {
	if index < len(c.cosl) {
		c.cosl[index].status = status
	}
}

func (c *KvClient) NewKv(key, value string) error {

	consul := getConsul(c)
	if consul == nil {
		return errors.New("No alive consul service")
	}

	rsp, err := ConsulRequest("PUT", consul.home+key, []byte(value))
	if err != nil {
		return err
	}

	if -1 == strings.Index(string(rsp), "true") {
		return errors.New("New Key Fail: " + string(rsp))
	}

	return nil
}

func (c *KvClient) GetKv(key string) ([]KvInfo, error) {

	consul := getConsul(c)
	if consul == nil {
		return nil, errors.New("No alive consul service")
	}

	rsp, err := ConsulRequest("GET", consul.home+key, nil)
	if err != nil {
		return nil, err
	}

	if len(rsp) == 0 {
		return nil, errors.New("Not found key:" + key)
	}

	data := make([]KvData, 0)

	err = json.Unmarshal(rsp, &data)
	if err != nil {
		return nil, err
	}

	kv := make([]KvInfo, len(data))

	for i, v := range data {
		buf, err := base64.StdEncoding.DecodeString(v.Value)
		if err != nil {
			return nil, err
		}

		kv[i].Key = v.Key
		kv[i].Value = string(buf)
	}

	return kv, nil
}

func (c *KvClient) DelKv(key string) error {

	consul := getConsul(c)
	if consul == nil {
		return errors.New("No alive consul service")
	}

	rsp, err := ConsulRequest("DELETE", consul.home+key, nil)
	if err != nil {
		return err
	}

	if -1 == strings.Index(string(rsp), "true") {
		return errors.New("Delete Key Fail: " + string(rsp))
	}

	return nil
}
