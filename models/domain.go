package models

import (
	"fmt"
	"log"
	"sort"
	"libvirt.org/go/libvirt"
	"libvirt.org/go/libvirtxml"
)

type Domain struct {
	Name       string
	ID         uint
	UUID       string
	State      int
	StateStr   string
	MaxMem     uint64
	MaxMemStr  string
	Memory     uint64
	MemoryStr  string
	Vcpu       uint
	CpuTime    uint64
	Disks      map[string]string  // disk device name and size
	Interfaces map[string]string  // MAC address and connected Network
}

func GetDomainStateStr(state libvirt.DomainState) string {
	var stateStr string = "No State"
	switch state {
	case 0:
		stateStr = "No State"
	case 1:
		stateStr = "Running"
	case 2:
		stateStr = "Blocked"
	case 3:
		stateStr = "Paused"
	case 4:
		stateStr = "Shutdown"
	case 5:
		stateStr = "Shutoff"
	case 6:
		stateStr = "Crashed"
	case 7:
		stateStr = "Pmsuspended"
	default:
		stateStr = "No State"
	}
	return stateStr
}

func GetFlag(flag string) libvirt.ConnectListAllDomainsFlags {
	var fg libvirt.ConnectListAllDomainsFlags
	switch flag {
	case "active":
		fg = libvirt.CONNECT_LIST_DOMAINS_ACTIVE
	case "inactive":
		fg = libvirt.CONNECT_LIST_DOMAINS_INACTIVE
	case "running":
		fg = libvirt.CONNECT_LIST_DOMAINS_RUNNING
	case "paused":
		fg = libvirt.CONNECT_LIST_DOMAINS_PAUSED
	case "shutoff":
		fg = libvirt.CONNECT_LIST_DOMAINS_SHUTOFF
	default:
		fg = libvirt.CONNECT_LIST_DOMAINS_PERSISTENT
	}
	return fg
}

/* ------------------------------------------------------------------------ */

func AllDomainsUUID() map[string]string {
	return nil
}

func GetAllDomains(flag string) []Domain {
	log.Println("Domain Model - get all domains")

	result := make([]Domain, 0)

	fg := GetFlag(flag)  // flag value: active, inactive, running, paused, shutoff, persistent
	domains, err := Conn.ListAllDomains(fg)
	if err != nil {
		log.Printf("Error: fail to get all domains")
	} else {
		for _, domain := range domains {
			d := new(Domain)
			info, _ := domain.GetInfo()
			d.Name, _ = domain.GetName()
			d.UUID, _  = domain.GetUUIDString()
			d.MaxMem = info.MaxMem
			d.MaxMemStr = ConvertSizeToString(info.MaxMem*1024, "GB")
			d.Memory = info.Memory
			d.MemoryStr = ConvertSizeToString(info.Memory*1024, "GB")
			d.Vcpu = info.NrVirtCpu
			d.CpuTime = info.CpuTime
			d.State = int(info.State)  // libvirt.DomainState
			d.StateStr = GetDomainStateStr(info.State)
			if info.State == 1 || info.State == 3 {
				d.ID, _ = domain.GetID()
			}

			domainxml, _ := domain.GetXMLDesc(0)
			domaincfg := &libvirtxml.Domain{}
			_ = domaincfg.Unmarshal(domainxml)

			var disks = map[string]string{}
			for _, disk := range domaincfg.Devices.Disks {
				if disk.Device == "disk" {
					blockInfo, err := domain.GetBlockInfo(disk.Target.Dev, 0)
					if err != nil {
						log.Printf("Error: fail to get disk %s capacity", disk.Target.Dev)
						disks[disk.Target.Dev] = "?"
					} else {
						disks[disk.Target.Dev] = ConvertSizeToString(blockInfo.Capacity, "GB")
					}
				}
			}
			d.Disks = disks

			var intfType, intfTypeNmae, intfTargetDev string
			var intfs = map[string]string{}  // todo: array of map to involve more details and keep item order.
			for _, intf := range domaincfg.Devices.Interfaces {
				if intf.Source.Network != nil {
					intfType = "N"
					intfTypeNmae = intf.Source.Network.Network
				} else if intf.Source.Bridge != nil {
					intfType = "B"
					intfTypeNmae = intf.Source.Bridge.Bridge
				}
				if intf.Target != nil {
					intfTargetDev = intf.Target.Dev
				} else {
					intfTargetDev = "?"
				}
				_ = intfTargetDev
				_ = intfType
				deviceID := intf.MAC.Address[9:]
				intfs[deviceID] = fmt.Sprintf("%s", intfTypeNmae)
			}
			d.Interfaces = intfs

			result = append(result, *d)
			log.Printf("Retrieve domain data (%s)", d.Name)
			domain.Free()
		}
	}
	// sort by name.
	sort.Slice(result, func(i, j int) bool {
		return result[i].Name < result[j].Name
	})
	return result
}