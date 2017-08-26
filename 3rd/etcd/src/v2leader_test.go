package etcd

import (
	"testing"
	"time"
)

func TestVoteLeader01(t *testing.T) {

	kvs := NewClient([]string{"localhost:2379"})

	err := kvs.VoteLeader()
	if err != nil {
		t.Errorf("%s", err.Error())
		return
	}
}

func TestKeepAlive01(t *testing.T) {

	kvs := NewClient([]string{"localhost:2379"})

	stats, err := kvs.KeepAlive()
	if err != nil {
		t.Errorf("%s", err.Error())
		return
	}

	t.Log(stats)

	stats2, err := kvs.KeepAlive()
	if err != nil {
		t.Errorf("%s", err.Error())
		return
	}

	t.Log(stats2)

	var t1, t2 time.Time

	t1.UnmarshalText([]byte(stats.StartTime))
	t2.UnmarshalText([]byte(stats2.StartTime))

	t.Log(t1, t2)

	if t1 != t2 {
		t.Error("diff time ", t1, t2)
	}
}
