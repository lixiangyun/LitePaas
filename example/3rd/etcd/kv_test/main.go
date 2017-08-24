package main

import (
	"LitePaas/3rd/etcd/src"
	"fmt"
	"time"
)

var kvs *etcd.EtcdClient

func banchmark(num int) {

	t1 := time.Now()

	for i := 0; i < num; i++ {
		value := fmt.Sprint("helloworld_%d", i)
		err := kvs.SetKeyValue("web", value)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}

	t2 := time.Now()

	subtime := t2.Sub(t1)

	fmt.Println("Time: ", subtime.Seconds())
	fmt.Println("Speed: ", float64(num)/subtime.Seconds())
}

func banchmark2(num int) {

	t1 := time.Now()

	for i := 0; i < num; i++ {

		_, err := kvs.GetKeyValue("web")
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}

	t2 := time.Now()

	subtime := t2.Sub(t1)

	fmt.Println("Time: ", subtime.Seconds())
	fmt.Println("Speed: ", float64(num)/subtime.Seconds())
}

func test1(times int, key string) {

	for i := 0; i < times; i++ {
		err := kvs.SetKeyValue(key, "helloworld")
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}
}

func main() {

	addr := []string{"localhost:2379"}

	kvs = etcd.NewClient(addr)

	err := kvs.SetKeyValue("message", "helloworld")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	data, err := kvs.GetKeyValue("message")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("KV: ", data)

	data, err = kvs.DelKeyValue("message")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("KV: ", data)

	err = kvs.SetKeyValueByTTL("message", "helloworld", 2)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	data, err = kvs.GetKeyValue("message")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("KV: ", data)

	time.Sleep(3 * time.Second)

	data, err = kvs.GetKeyValue("message")

	fmt.Println(err.Error())

	fmt.Println("KV: ", data)

	go test1(1, "wait")

	data, err = kvs.WatchKeyValue("wait")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("KV: ", data)

	go test1(10, "waitindex")

	data, err = kvs.WatchKeyValueByTimes("waitindex", 5)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("KV: ", data)

	//banchmark(3000)
	//banchmark2(3000)
}
