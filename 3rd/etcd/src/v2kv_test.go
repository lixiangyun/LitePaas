package etcd

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestSetKeyValue01(t *testing.T) {

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

func TestCreateDir01(t *testing.T) {

	kvs := NewClient([]string{"localhost:2379"})

	kvs.DeleteDir("/test_dir", true)

	err := kvs.CreateDir("/test_dir")
	if err != nil {
		t.Errorf("%s", err.Error())
		return
	}

	err = kvs.DeleteDir("/test_dir", false)
	if err != nil {
		t.Errorf("%s", err.Error())
		return
	}
}

func TestListDir01(t *testing.T) {

	kvs := NewClient([]string{"localhost:2379"})

	err := kvs.CreateDir("/test_dir")
	if err != nil {
		t.Errorf("%s", err.Error())
		return
	}

	for i := 0; i < 10; i++ {

		key := fmt.Sprintf("/test_dir/test_key%d", i)
		value := fmt.Sprintf("test_%d", i)

		err := kvs.SetKeyValue(key, value)
		if err != nil {
			t.Errorf("%s", err.Error())
			return
		}
	}

	data, err := kvs.ListDir("/test_dir", false)
	if err != nil {
		t.Errorf("%s", err.Error())
		return
	}

	if len(data) != 10 {
		t.Errorf("get file list error (%d)", len(data))
	}

	for _, v := range data {

		var i int

		fmt.Sscanf(v.Key, "/test_dir/test_key%d", &i)

		key := fmt.Sprintf("/test_dir/test_key%d", i)
		value := fmt.Sprintf("test_%d", i)

		if v.IsDir == true || v.Key != key || v.Value != value {
			t.Errorf("get data failed! ", v.IsDir)
			t.Error(v)
		}
	}

	err = kvs.DeleteDir("/test_dir", true)
	if err != nil {
		t.Errorf("%s", err.Error())
		return
	}
}
