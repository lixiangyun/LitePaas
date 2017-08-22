package main

import (
	"LitePaas/3rd/etcd/src"
	"fmt"
)

func main() {

	addr := []string{"localhost:2379"}

	kvs := etcd.NewClient(addr)

	err := kvs.NewKv("web", "helloworld")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	data, err := kvs.GetKv("web")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("Items:", len(data))
	fmt.Println("KV: ", data)

	err = kvs.DelKv("web")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
