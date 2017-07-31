package main

import (
	"LitePaas/3rd/kv"
	"fmt"
)

func main() {

	kvs := kv.NewKvClient(nil)

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

}
