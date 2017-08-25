package etcd

import (
	"encoding/json"
	"errors"
	"fmt"
)

type Node struct {
	IsDir bool   `json:"dir"`
	Key   string `json:"key"`
	Value string `json:"value"`
	Nodes []Node `json:"nodes"`
}

type Rsponse struct {
	Action   string `json:"action"`
	NowNode  Node   `json:"node"`
	PrevNode Node   `json:"prevNode"`
	Nodes    []Node `json:"nodes"`
}

func (c *EtcdClient) SetKeyValue(key, value string) error {

	if key[0] != '/' {
		return errors.New("input invaild " + key)
	}

	client := getClient(c)
	if client == nil {
		return errors.New("not alive etcd service")
	}

	value = "value=" + value

	rsp, err := HttpRequest("PUT", client.home+"/v2/keys"+key+"?"+value, nil)
	if err != nil {
		return err
	}

	var result Rsponse

	err = json.Unmarshal(rsp, &result)
	if err != nil {
		return err
	}

	if result.Action != "set" {
		return errors.New("rsponse error! " + string(rsp))
	}

	return nil
}

func (c *EtcdClient) SetKeyValueByTTL(key, value string, ttl int) error {

	if key[0] != '/' {
		return errors.New("input invaild " + key)
	}

	client := getClient(c)
	if client == nil {
		return errors.New("not alive etcd service")
	}

	value = fmt.Sprintf("value=%s&ttl=%d", value, ttl)

	fmt.Println(value)

	rsp, err := HttpRequest("PUT", client.home+"/v2/keys"+key+"?"+value, nil)
	if err != nil {
		return err
	}

	var result Rsponse

	err = json.Unmarshal(rsp, &result)
	if err != nil {
		return err
	}

	if result.Action != "set" {
		return errors.New("rsponse error! " + string(rsp))
	}

	return nil
}

func (c *EtcdClient) GetKeyValue(key string) (node Node, err error) {

	if key[0] != '/' {
		err = errors.New("input invaild " + key)
		return
	}

	client := getClient(c)
	if client == nil {
		err = errors.New("not alive etcd service")
		return
	}

	rsp, err := HttpRequest("GET", client.home+"/v2/keys"+key, nil)
	if err != nil {
		return
	}

	var result Rsponse

	err = json.Unmarshal(rsp, &result)
	if err != nil {
		return
	}

	if result.Action != "get" {
		err = errors.New("rsponse error! " + string(rsp))
		return
	}

	node = result.NowNode
	return
}

func (c *EtcdClient) WatchKeyValue(key string) (node Node, err error) {

	if key[0] != '/' {
		err = errors.New("input invaild " + key)
		return
	}

	client := getClient(c)
	if client == nil {
		err = errors.New("not alive etcd service")
		return
	}

	rsp, err := HttpRequest("GET", client.home+"/v2/keys"+key+"?wait=true", nil)
	if err != nil {
		return
	}

	var result Rsponse

	err = json.Unmarshal(rsp, &result)
	if err != nil {
		return
	}

	if result.Action != "set" {
		err = errors.New("rsponse error! " + string(rsp))
		return
	}

	node = result.NowNode
	return
}

func (c *EtcdClient) WatchKeyValueByTimes(key string, times int) (node Node, err error) {

	if key[0] != '/' {
		err = errors.New("input invaild " + key)
		return
	}

	client := getClient(c)
	if client == nil {
		err = errors.New("not alive etcd service")
		return
	}

	value := fmt.Sprintf("%s/v2/keys%s?wait=true&waitIndex=%d", client.home, key, times)

	rsp, err := HttpRequest("GET", value, nil)
	if err != nil {
		return
	}

	var result Rsponse

	err = json.Unmarshal(rsp, &result)
	if err != nil {
		return
	}

	if result.Action != "set" {
		err = errors.New("rsponse error! " + string(rsp))
		return
	}

	node = result.NowNode
	return
}

func (c *EtcdClient) DelKeyValue(key string) (node Node, err error) {

	if key[0] != '/' {
		err = errors.New("input invaild " + key)
		return
	}

	client := getClient(c)
	if client == nil {
		err = errors.New("not alive etcd service")
		return
	}

	rsp, err := HttpRequest("DELETE", client.home+"/v2/keys"+key, nil)
	if err != nil {
		return
	}

	var result Rsponse

	err = json.Unmarshal(rsp, &result)
	if err != nil {
		return
	}

	if result.Action != "delete" {
		err = errors.New("rsponse error! " + string(rsp))
		return
	}

	node = result.NowNode
	return
}

func (c *EtcdClient) CreateDir(dir string) error {

	if dir[0] != '/' {
		return errors.New("input invaild " + dir)
	}

	client := getClient(c)
	if client == nil {
		return errors.New("not alive etcd service")
	}

	rsp, err := HttpRequest("PUT", client.home+"/v2/keys"+dir+"?dir=true", nil)
	if err != nil {
		return err
	}

	var result Rsponse

	err = json.Unmarshal(rsp, &result)
	if err != nil {
		return err
	}

	if result.Action != "set" {
		return errors.New("rsponse error! " + string(rsp))
	}

	return nil
}

func (c *EtcdClient) ListDir(dir string, recursive bool) error {

	if dir[0] != '/' {
		return errors.New("input invaild " + dir)
	}

	client := getClient(c)
	if client == nil {
		return errors.New("not alive etcd service")
	}

	if recursive == true {
		dir = fmt.Sprintf("%s?recursive=true", dir)
	}

	rsp, err := HttpRequest("DELETE", client.home+"/v2/keys"+dir, nil)
	if err != nil {
		return err
	}

	var result Rsponse

	err = json.Unmarshal(rsp, &result)
	if err != nil {
		return err
	}

	if result.Action != "delete" {
		return errors.New("rsponse error! " + string(rsp))
	}

	return nil
}

func (c *EtcdClient) DeleteDir(dir string, recursive bool) error {

	if dir[0] != '/' {
		return errors.New("input invaild " + dir)
	}

	client := getClient(c)
	if client == nil {
		return errors.New("not alive etcd service")
	}

	if recursive == true {
		dir = fmt.Sprintf("%s?recursive=true", dir)
	} else {
		dir = fmt.Sprintf("%s?dir=true", dir)
	}

	rsp, err := HttpRequest("DELETE", client.home+"/v2/keys"+dir, nil)
	if err != nil {
		return err
	}

	//fmt.Println(string(rsp))

	var result Rsponse

	err = json.Unmarshal(rsp, &result)
	if err != nil {
		return err
	}

	if result.Action != "delete" {
		return errors.New("rsponse error! " + string(rsp))
	}

	return nil
}
