package models

import (
	"log"
	"sort"
	"libvirt.org/go/libvirt"
	//"libvirt.org/go/libvirtxml"
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
			d.MaxMemStr = ConvertSizeToString(info.MaxMem*1024, "MB")
			d.Memory = info.Memory
			d.MemoryStr = ConvertSizeToString(info.Memory*1024, "MB")
			d.Vcpu = info.NrVirtCpu
			d.CpuTime = info.CpuTime
			d.State = int(info.State)  // libvirt.DomainState
			d.StateStr = GetDomainStateStr(info.State)
			if info.State == 1 || info.State == 3 {
				id, err := domain.GetID()
				if err == nil {
					d.ID = id
				} else {
					log.Printf("Error: fail to get domain ID")
				}
			}
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