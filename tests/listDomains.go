package main

import (
	"os"
	"fmt"
	"libvirt.org/go/libvirt"
	"libvirt.org/go/libvirtxml"
)

func main() {

	conn, err := libvirt.NewConnect("qemu:///system")
	if err != nil {
		fmt.Println("Fail to create libvirt connection.")
	}
	defer conn.Close()

	doms, err := conn.ListAllDomains(libvirt.CONNECT_LIST_DOMAINS_ACTIVE)
	if err != nil {
		fmt.Println("Fail to get active domains.")
	}
	fmt.Printf("%d active domains:\n", len(doms))
	IterateDomains(doms)
	
	doms, err = conn.ListAllDomains(libvirt.CONNECT_LIST_DOMAINS_INACTIVE)
	if err != nil {
		fmt.Println("Fail to get inactive domains.")
	}
	fmt.Printf("%d inactive domains:\n", len(doms))
	IterateDomains(doms)

	os.Exit(0)
}

func IterateDomains(doms []libvirt.Domain) {
	var intfType string
	var intfTypeNmae string
	var intfTargetDev string
	var count int
	for i, dom := range doms {
		count = i + 1
		name, err := dom.GetName()
		if err == nil {
			domxml, _ := dom.GetXMLDesc(0)
			domcfg := &libvirtxml.Domain{}
			_ = domcfg.Unmarshal(domxml)
			fmt.Printf("  %02d - %s | %s | %d | %d\n", 
			        count,
					domcfg.UUID, 
					name, 
					domcfg.VCPU.Value, 
					domcfg.Memory.Value)
			for _, disk := range domcfg.Devices.Disks {
				fmt.Printf("       %s | %s | %s | %s\n", 
					disk.Device, 
					disk.Target.Dev,
					disk.Target.Bus,
					//disk.Driver.Name, 
					//disk.Driver.Type,
					disk.Source.File.File)
			}
			for _, intf := range domcfg.Devices.Interfaces {
				if intf.Source.Network != nil {
					intfType = "network"
					intfTypeNmae = intf.Source.Network.Network
				} else if intf.Source.Bridge != nil {
					intfType = "bridge"
					intfTypeNmae = intf.Source.Bridge.Bridge
				}
				if intf.Target != nil {
					intfTargetDev = intf.Target.Dev
				}
				fmt.Printf("       %s | %s | %s | %s | %s\n",
					intfType,
					intfTypeNmae,
					intf.MAC.Address,
					intf.Model.Type,
					intfTargetDev)
			}
		}
		dom.Free()
	}
	fmt.Println("")
}