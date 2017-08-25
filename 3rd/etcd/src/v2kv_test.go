package etcd

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestSetKeyValue(t *testing.T) {

	kvs := NewClient([]string{"localhost:2379"})

	defer kvs.DelKeyValue("/message")

	err := kvs.SetKeyValue("/message", "helloworld")
	if err != nil {
		t.Errorf("%s", err.Error())
		return
	}

	data, err := kvs.GetKeyValue("/message")
	if err != nil {
		t.Errorf("%s", err.Error())
		return
	}

	t.Log("1", data)

	if data.IsDir == true || data.Key != "/message" || data.Value != "helloworld" {
		t.Errorf("get data failed! ", data.IsDir)
		t.Error(data)
	}

	t.Log("1", data)

	if data.IsDir == true || data.Key != "/message" {
		t.Errorf("get data failed! ", data.IsDir)
		t.Error(data)
	}
}

func TestSetKeyValueByTTL(t *testing.T) {

	kvs := NewClient([]string{"localhost:2379"})

	defer kvs.DelKeyValue("/msg_ttl")

	err := kvs.SetKeyValueByTTL("/msg_ttl", "1234567890", 1)
	if err != nil {
		t.Errorf("%s", err.Error())
		return
	}

	data, err := kvs.GetKeyValue("/msg_ttl")
	if err != nil {
		t.Errorf("%s", err.Error())
		return
	}

	t.Log("1", data)

	if data.IsDir == true || data.Key != "/msg_ttl" || data.Value != "1234567890" {
		t.Errorf("get data failed! ", data.IsDir)
		t.Error(data)
	}

	time.Sleep(2 * time.Second)

	data, err = kvs.GetKeyValue("/msg_ttl")
	if err == nil {
		t.Errorf("find the key after ttl!", data)
	}

	if strings.IndexAny(err.Error(), "errorCode") > 0 {
		t.Errorf("can not return error!", err.Error())
	}

	t.Log(err.Error())
}

func testSetKV(times int, key string, t *testing.T, c *EtcdClient) {

	time.Sleep(500 * time.Millisecond)

	for i := 0; i < times; i++ {
		value := fmt.Sprintf("%d", i)
		err := c.SetKeyValue(key, value)
		if err != nil {
			t.Errorf("%s", err.Error())
			return
		}
	}
}

func TestWatchKeyValue01(t *testing.T) {
	kvs := NewClient([]string{"localhost:2379"})

	defer kvs.DelKeyValue("/watch_test")

	err := kvs.SetKeyValue("/watch_test", "test_value")
	if err != nil {
		t.Errorf("%s", err.Error())
		return
	}

	go testSetKV(1, "/watch_test", t, kvs)

	data, err := kvs.WatchKeyValue("/watch_test")
	if err != nil {
		t.Errorf("%s", err.Error())
		return
	}

	t.Log("1", data)

	if data.IsDir == true || data.Key != "/watch_test" || data.Value != "0" {
		t.Errorf("get data failed! ", data.IsDir)
		t.Error(data)
	}
}

/*
func TestWatchKeyValueByTimes01(t *testing.T) {

	kvs := NewClient([]string{"localhost:2379"})

	defer kvs.DelKeyValue("/watch_test")

	err := kvs.SetKeyValue("/watch_test", "test2")
	if err != nil {
		t.Errorf("%s", err.Error())
		return
	}

	go testSetKV(10, "/watch_test", t, kvs)

	data, err := kvs.WatchKeyValueByTimes("/watch_test", 3205)
	if err != nil {
		t.Errorf("%s", err.Error())
		return
	}

	t.Log("1", data)

	if data.IsDir == true || data.Key != "/watch_test" || data.Value != "4" {
		t.Errorf("get data failed! ", data.IsDir)
		t.Error(data)
	}
}*/
