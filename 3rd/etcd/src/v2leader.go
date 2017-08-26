package etcd

import (
	"encoding/json"
	"errors"
)

type Count struct {
	Success int `json:"success"`
	Fail    int `json:"fail"`
}

type Latency struct {
	Average      float64 `json:"average"`
	Current      float64 `json:"current"`
	Maximum      float64 `json:"maximum"`
	Minimum      float64 `json:"minimum"`
	StdDeviation float64 `json:"standardDeviation"`
}

type Leader struct {
	LeaderName string `json:"leader"`
}

func (c *EtcdClient) VoteLeader() error {

	client := getClient(c)
	if client == nil {
		return errors.New("not alive etcd service")
	}

	rsp, err := HttpRequest("GET", client.home+"/v2/stats/leader", nil)
	if err != nil {
		return err
	}

	var result Leader

	err = json.Unmarshal(rsp, &result)
	if err != nil {
		return err
	}

	if result.LeaderName == "" {
		return errors.New("rsponse error! " + string(rsp))
	}

	return nil
}

type LeaderInfo struct {
	Leader    string `json:"leader"`
	StartTime string `json:"startTime"`
	UpTime    string `json:"uptime"`
}

type SelfStats struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Role      string `json:"state"`
	StartTime string `json:"startTime"`

	RecvAppRequestCnt int `json:"recvAppendRequestCnt"`
	SendAppRequestCnt int `json:"sendAppendRequestCnt"`

	Leader LeaderInfo `json:"leaderInfo"`
}

func (c *EtcdClient) KeepAlive() (SelfStats, error) {

	var result SelfStats

	client := getClient(c)
	if client == nil {
		return result, errors.New("not alive etcd service")
	}

	rsp, err := HttpRequest("GET", client.home+"/v2/stats/self", nil)
	if err != nil {
		return result, err
	}

	err = json.Unmarshal(rsp, &result)
	if err != nil {
		return result, err
	}

	if result.Id == "" {
		return result, errors.New("rsponse error! " + string(rsp))
	}

	return result, nil
}
