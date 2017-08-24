package etcd

import (
	"encoding/json"
	"errors"
	"fmt"
)

type Node struct {
	Key           string `json:"key"`
	Value         string `json:"value"`
	CreatedIndex  int    `json:"createdIndex"`
	ModifiedInfex int    `json:"modifiedIndex"`
}

type Rsponse struct {
	Action   string `json:"action"`
	NowNode  Node   `json:"node"`
	PrevNode Node   `json:"prevNode"`
}

func (c *EtcdClient) SetKeyValue(key, value string) error {

	client := getClient(c)
	if client == nil {
		return errors.New("not alive etcd service")
	}

	value = "value=" + value

	//fmt.Println(value)

	rsp, err := HttpRequest("PUT", client.home+"/v2/keys/"+key+"?"+value, nil)
	if err != nil {
		return err
	}

	//fmt.Println(string(rsp))

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

	client := getClient(c)
	if client == nil {
		return errors.New("not alive etcd service")
	}

	value = fmt.Sprintf("value=%s&ttl=%d", value, ttl)

	fmt.Println(value)

	rsp, err := HttpRequest("PUT", client.home+"/v2/keys/"+key+"?"+value, nil)
	if err != nil {
		return err
	}

	//fmt.Println(string(rsp))

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

	client := getClient(c)
	if client == nil {
		err = errors.New("not alive etcd service")
		return
	}

	rsp, err := HttpRequest("GET", client.home+"/v2/keys/"+key, nil)
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

	client := getClient(c)
	if client == nil {
		err = errors.New("not alive etcd service")
		return
	}

	rsp, err := HttpRequest("GET", client.home+"/v2/keys/"+key+"?wait=true", nil)
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

	client := getClient(c)
	if client == nil {
		err = errors.New("not alive etcd service")
		return
	}

	value := fmt.Sprintf("%s/v2/keys/%s?wait=true&waitIndex=%d", client.home, key, times)

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

	client := getClient(c)
	if client == nil {
		err = errors.New("not alive etcd service")
		return
	}

	rsp, err := HttpRequest("DELETE", client.home+"/v2/keys/"+key, nil)
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
