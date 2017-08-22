package etcd

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Etcd struct {
	home   string
	status bool
	index  int
}

type EtcdClient struct {
	client []Etcd
}

type KvInfo struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type GetKv struct {
	Key   string `json:"key"`
}

type Header struct {
	ClusterId string `json:"cluster_id"`
	MemberId  string `json:"member_id"`
	Revision  string `json:"revision"`
	RaftTerm  string `json:"raft_term"`
}

type KvPutRsp struct {
	H Header `json:"header"`
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

func HttpRequest(method string, url string, req []byte) (rsp []byte, err error) {
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

func NewClient(addr []string) (c *EtcdClient) {

	if len(addr) == 0 {
		addr = []string{"localhost:2379"}
	}

	node := make([]Etcd, len(addr))

	for i, v := range addr {
		node[i].home = "http://" + v + "/v3alpha/kv"
		node[i].status = true
		node[i].index = i
	}

	return &EtcdClient{client: node}
}

func getClient(c *EtcdClient) *Etcd {
	for _, v := range c.client {
		if v.status == true {
			return &v
		}
	}

	return nil
}

func setEtcdStatus(index int, status bool, c *EtcdClient) {
	if index < len(c.client) {
		c.client[index].status = status
	}
}

func (c *EtcdClient) NewKv(key, value string) error {

	consul := getClient(c)
	if consul == nil {
		return errors.New("No alive consul service")
	}

	var kv KvInfo
	kv.Key = base64.StdEncoding.EncodeToString([]byte(key))
	kv.Value = base64.StdEncoding.EncodeToString([]byte(value))

	buf, err := json.Marshal(kv)
	if err != nil {
		return err
	}

	//fmt.Println(consul)
	//fmt.Println(string(buf))

	rsp, err := ConsulRequest("POST", consul.home+"/put", buf)
	if err != nil {
		return err
	}

	//fmt.Println(string(rsp))

	var result KvPutRsp

	err = json.Unmarshal(rsp, &result)
	if err != nil {
		return err
	}

	fmt.Println(result)

	return nil
}

func (c *EtcdClient) GetKv(key string) ([]KvInfo, error) {

	consul := getClient(c)
	if consul == nil {
		return nil, errors.New("No alive consul service")
	}

	var kv GetKv
	kv.Key = base64.StdEncoding.EncodeToString([]byte(key))\

	buf, err := json.Marshal(kv)
	if err != nil {
		return err
	}

	rsp, err := ConsulRequest("POST", consul.home+"/range", buf)
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

func (c *EtcdClient) DelKv(key string) error {

	consul := getClient(c)
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
