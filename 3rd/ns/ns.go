package ns

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
)

type Consul struct {
	home   string
	status bool
	index  int
}

type NsClient struct {
	cosl []Consul
}

type NsService struct {
	ID      string   `json:"ID"`
	Name    string   `json:"Service"`
	Address string   `json:"Address"`
	Port    int      `json:"Port"`
	Tags    []string `json:"Tags"`
}

type NsItem struct {
	Datacenter      string            `json:"Datacenter"`
	Node            string            `json:"Node"`
	Address         string            `json:"Address"`
	TaggedAddresses map[string]string `json:"TaggedAddresses"`
	NodeMeta        map[string]string `json:"NodeMeta"`
	Service         NsService         `json:"Service"`
}

type NsData struct {
	UUID            string            `json:"ID"`
	Datacenter      string            `json:"Datacenter"`
	Node            string            `json:"Node"`
	Address         string            `json:"Address"`
	TaggedAddresses map[string]string `json:"TaggedAddresses"`
	NodeMeta        map[string]string `json:"NodeMeta"`
	ServiceID       string            `json:"ServiceID"`
	ServiceName     string            `json:"ServiceName"`
	ServiceTags     []string          `json:"ServiceTags"`
	ServiceAddress  string            `json:"ServiceAddress"`
	Port            int               `json:"ServicePort"`
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

func NewNsClient(addr []string) (c *NsClient) {

	if len(addr) == 0 {
		addr = []string{"localhost:8500"}
	}

	consul := make([]Consul, len(addr))

	for i, v := range addr {
		consul[i].home = "http://" + v + "/v1"
		consul[i].status = true
		consul[i].index = i
	}

	return &NsClient{cosl: consul}
}

func getConsul(c *NsClient) *Consul {
	for _, v := range c.cosl {
		if v.status == true {
			return &v
		}
	}
	return nil
}

func setConsulStatus(index int, status bool, c *NsClient) {
	if index < len(c.cosl) {
		c.cosl[index].status = status
	}
}

func (c *NsClient) RegisterService(ns NsItem) error {

	consul := getConsul(c)
	if consul == nil {
		return errors.New("No alive consul service")
	}

	request, err := json.Marshal(ns)
	if err != nil {
		return err
	}

	rsp, err := ConsulRequest("PUT", consul.home+"/catalog/register", request)
	if err != nil {
		return err
	}

	if -1 == strings.Index(string(rsp), "true") {
		return errors.New("register failed ! " + string(rsp))
	}

	return nil
}

func (c *NsClient) GetService(servername string) ([]NsData, error) {

	consul := getConsul(c)
	if consul == nil {
		return nil, errors.New("No alive consul service")
	}

	rsp, err := ConsulRequest("GET", consul.home+"/catalog/service/"+servername, nil)
	if err != nil {
		return nil, err
	}

	if len(rsp) == 0 {
		return nil, errors.New("Not found server:" + servername)
	}

	data := make([]NsData, 0)

	err = json.Unmarshal(rsp, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

type DeNsItem struct {
	Datacenter string `json:"Datacenter"`
	Node       string `json:"Node"`
	ServiceID  string `json:"ServiceID"`
}

func (c *NsClient) DeRegisterService(ns DeNsItem) error {
	consul := getConsul(c)
	if consul == nil {
		return errors.New("No alive consul service")
	}

	request, err := json.Marshal(ns)
	if err != nil {
		return err
	}

	rsp, err := ConsulRequest("PUT", consul.home+"/catalog/deregister", request)
	if err != nil {
		return err
	}

	if -1 == strings.Index(string(rsp), "true") {
		return errors.New("register failed ! " + string(rsp))
	}

	return nil
}

func (c *NsClient) GetDataCenters() ([]string, error) {
	consul := getConsul(c)
	if consul == nil {
		return nil, errors.New("No alive consul service")
	}

	rsp, err := ConsulRequest("PUT", consul.home+"/catalog/datacenters", nil)
	if err != nil {
		return nil, err
	}

	if len(rsp) == 0 {
		return nil, errors.New("Not found datacenters")
	}

	dc := make([]string, 0)

	err = json.Unmarshal(rsp, &dc)
	if err != nil {
		return nil, err
	}

	if 0 == len(dc) {
		return nil, errors.New("Not found datacenters")
	}

	return dc, nil
}
