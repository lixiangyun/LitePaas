package etcd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type Etcd struct {
	home   string
	status bool
	index  int
}

type EtcdClient struct {
	client []Etcd
}

func ReadFully(conn io.ReadCloser) ([]byte, error) {
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

	rsp, err = ReadFully(rspon.Body)
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
		node[i].home = "http://" + v
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

type VersionInfo struct {
	Server  string `json:"etcdserver"`
	Cluster string `json:"etcdcluster"`
}

func (c *EtcdClient) GetVersion() (ver VersionInfo, err error) {

	client := getClient(c)
	if client == nil {
		err = errors.New("not alive consul service")
		return
	}

	rsp, err := HttpRequest("GET", client.home+"/version", nil)
	if err != nil {
		return
	}

	fmt.Print(string(rsp))

	err = json.Unmarshal(rsp, &ver)
	if err != nil {
		return
	}

	return
}

type EtcdStats struct {
	CompareAndSwapFail    int `json:"compareAndSwapFail"`
	CompareAndSwapSuccess int `json:"compareAndSwapSuccess"`
	CreateFail            int `json:"createFail"`
	CreateSuccess         int `json:"createSuccess"`
	DeleteFail            int `json:"deleteFail"`
	DeleteSuccess         int `json:"deleteSuccess"`
	ExpireCount           int `json:"expireCount"`
	GetsFail              int `json:"getsFail"`
	GetsSuccess           int `json:"getsSuccess"`
	SetsFail              int `json:"setsFail"`
	SetsSuccess           int `json:"setsSuccess"`
	UpdateFail            int `json:"updateFail"`
	UpdateSuccess         int `json:"updateSuccess"`
	Watchers              int `json:"watchers"`
}

func (c *EtcdClient) GetStats() (stats EtcdStats, err error) {

	client := getClient(c)
	if client == nil {
		err = errors.New("not alive consul service")
		return
	}

	rsp, err := HttpRequest("GET", client.home+"/v2/stats/store", nil)
	if err != nil {
		return
	}

	err = json.Unmarshal(rsp, &stats)
	if err != nil {
		return
	}

	return
}
