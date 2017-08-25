package etcd

import (
	"testing"
)

func TestGetVersion01(t *testing.T) {

	kvs := NewClient([]string{"localhost:2379"})

	ver, err := kvs.GetVersion()
	if err != nil {
		t.Errorf("%s", err.Error())
		t.Error(ver)
		return
	}

	t.Log(ver)

	if ver.Server == "" {
		t.Errorf("get version info failed! (%s,%s)", ver.Server, ver.Cluster)
	}
}

func TestGetStats01(t *testing.T) {
	kvs := NewClient([]string{"localhost:2379"})

	stats, err := kvs.GetStats()
	if err != nil {
		t.Errorf("%s", err.Error())
		t.Error(stats)
		return
	}

	t.Log(stats)
}
