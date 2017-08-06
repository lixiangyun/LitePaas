package main

import (
	"LitePaas/3rd/kv"
	"fmt"
	"time"
)

func benchmarktest_NewKv(kvs *kv.KvClient, sec int) {

	t1 := time.Now()

	for i := 0; i < 10000; i++ {

		key := fmt.Sprintf("key/%02x", i)

		err := kvs.NewKv(key, "helloworld")
		if err != nil {
			fmt.Println(err.Error())
		}

	}

	t2 := time.Now()

	t2.Sub(t1)

	fmt.Println("time = ", t2.Second())
	fmt.Println("speed = ", 10000/t2.Second(), " ps")
}

func main() {

	addr := []string{"192.168.0.107:8500"}

	kvs := kv.NewKvClient(addr)

	err := kvs.NewKv("web", "helloworld")
	if err != nil {
		fmt.Println(err.Error())
	}

	data, err := kvs.GetKv("web")
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println("Items:", len(data))
	fmt.Println("KV: ", data)

	err = kvs.DelKv("web")
	if err != nil {
		fmt.Println(err.Error())
	}

	benchmarktest_NewKv(kvs, 0)

}
