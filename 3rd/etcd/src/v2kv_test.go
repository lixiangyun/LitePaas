package etcd

import (
	"testing"
)

func TestSetKeyValue(t *testing.T) {

	kvs := NewClient([]string{"localhost:2379"})

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

	data, err = kvs.DelKeyValue("/message")
	if err != nil {
		t.Errorf("%s", err.Error())
		return
	}

	t.Log("1", data)

	if data.IsDir == true || data.Key != "/message" {
		t.Errorf("get data failed! ", data.IsDir)
		t.Error(data)
	}
}
