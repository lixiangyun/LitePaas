package main

import (
	"LitePaas/3rd/ns"
	"fmt"
)

func main() {

	addr := []string{"127.0.0.1:8500"}

	nss := ns.NewNsClient(addr)

	var hw ns.NsItem
	var service ns.NsService

	hw.Datacenter = "dc1"
	hw.Node = "abc"
	hw.NodeMeta = make(map[string]string, 0)
	hw.TaggedAddresses = make(map[string]string, 0)
	hw.Address = "192.168.0.1"

	service.Address = "192.168.0.100"
	service.ID = "001"
	service.Name = "helloworld"
	service.Port = 1010
	service.Tags = make([]string, 0)

	hw.Service = service

	err := nss.RegisterService(hw)
	if err != nil {
		fmt.Println(err.Error())
	}

	service.Address = "192.168.0.99"
	service.ID = "002"
	service.Name = "helloworld"
	service.Port = 1010
	service.Tags = make([]string, 0)

	hw.Service = service

	err = nss.RegisterService(hw)
	if err != nil {
		fmt.Println(err.Error())
	}

	data, err := nss.GetService("helloworld")
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println("Items:", len(data))
	fmt.Println("server: ", data)

}
