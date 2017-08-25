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
		err := kvs.SetKeyValue("/web", value)
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

		_, err := kvs.GetKeyValue("/web")
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

func main() {

	kvs = etcd.NewClient([]string{"localhost:2379"})

	banchmark(3000)
	banchmark2(3000)
}
