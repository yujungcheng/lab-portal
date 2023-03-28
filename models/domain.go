package models

import (
	"log"
	"sort"
	"strings"
	"strconv"
	"os/exec"
	"path/filepath"
	"libvirt.org/go/libvirt"
	"libvirt.org/go/libvirtxml"
)

type DomainGroup struct {  // unused
	GroupType string
	GroupName string
	Domains   []Domain
}

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
	Disks      []map[string]string // disk device name, sizes
	Interfaces []map[string]string // MAC address, connected Network/Bridge
	Metadata   map[string]string   // extra info from description
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

func parserDescription(desc string) map[string]string {
	result := map[string]string{
		"group":     "default", // default group
		"tag":       "",        // set tag
		"backup":    "",        // backup directory path of domain
		"account":   "",        // record username/password
		"ipaddress": "",        // network ip address
		"template":  "false",   // set as template or not
	}
	descLines := strings.Split(desc, "\n")
	for _, descLine := range descLines {
		for k := range result {
			if strings.Contains(descLine, k+"=") {
				tmp := strings.Split(descLine, "=")
				result[k] = tmp[1]
			}
		}
	}
	return result
}

/* ------------------------------------------------------------------------ */

func GetAllDomains(flag string) []Domain {
	log.Println("Domain Model - get all domains")

	result := make([]Domain, 0)

	fg := GetFlag(flag) // flag value: active, inactive, running, paused, shutoff, persistent
	domains, err := Conn.ListAllDomains(fg)
	if err != nil {
		log.Printf("Error: fail to get all domains")
	} else {
		for _, domain := range domains {
			d := new(Domain)
			info, _ := domain.GetInfo()
			d.Name, _ = domain.GetName()
			d.UUID, _ = domain.GetUUIDString()
			d.MaxMem = info.MaxMem
			d.MaxMemStr = ConvertSizeToString(info.MaxMem*1024, "GB")
			d.Memory = info.Memory
			d.MemoryStr = ConvertSizeToString(info.Memory*1024, "GB")
			d.Vcpu = info.NrVirtCpu
			d.CpuTime = info.CpuTime
			d.State = int(info.State) // libvirt.DomainState
			d.StateStr = GetDomainStateStr(info.State)
			if info.State == 1 || info.State == 3 {
				d.ID, _ = domain.GetID()
			}

			log.Printf("+ Retriving domain data (%s)", d.Name)

			domainxml, _ := domain.GetXMLDesc(0)
			domaincfg := &libvirtxml.Domain{}
			_ = domaincfg.Unmarshal(domainxml)

			var diskCapacity, diskAllocation, diskPhysical string
			var disks = []map[string]string{}
			for _, disk := range domaincfg.Devices.Disks {
				if disk.Device == "disk" {
					blockInfo, err := domain.GetBlockInfo(disk.Target.Dev, 0)
					if err != nil {
						log.Printf("Error: fail to get disk %s info", disk.Target.Dev)
					} else {
						diskCapacity = ConvertSizeToString(blockInfo.Capacity, "GB")
						diskAllocation = ConvertSizeToString(blockInfo.Allocation, "GB")
						diskPhysical = ConvertSizeToString(blockInfo.Physical, "GB")
					}
					// todo: get storage pool of disk
					d := map[string]string{
						"name":       disk.Target.Dev,
						"file":       disk.Source.File.File,
						"capacity":   diskCapacity,
						"allocation": diskAllocation,
						"physical":   diskPhysical,
						"storagePool": "",
					}
					disks = append(disks, d)
				} else if disk.Device == "cdrom" {
					log.Printf("  - disk %s is CDROM", disk.Target.Dev)
				}
			}
			d.Disks = disks

			var intfType, intfTypeNmae, intfTargetDev string
			var intfs = []map[string]string{}
			for _, intf := range domaincfg.Devices.Interfaces {
				if intf.Source.Network != nil {
					intfType = "Network"
					intfTypeNmae = intf.Source.Network.Network
				} else if intf.Source.Bridge != nil {
					intfType = "Bridge"
					intfTypeNmae = intf.Source.Bridge.Bridge
				}
				if intf.Target != nil {
					intfTargetDev = intf.Target.Dev
				}
				i := map[string]string{
					"mac":    intf.MAC.Address, //"mac": intf.MAC.Address[9:],
					"name":   intfTypeNmae,
					"type":   intfType,
					"target": intfTargetDev,
				}
				intfs = append(intfs, i)
			}
			d.Interfaces = intfs

			// parser metadata from desc
			d.Metadata = parserDescription(domaincfg.Description)
			result = append(result, *d)
			domain.Free()
		}
	}
	// sort by name.
	sort.Slice(result, func(i, j int) bool {
		return result[i].Name < result[j].Name
	})
	return result
}

func GetAllDomainsByGroup(flag, groupBy string) map[string][]Domain {
	var groupName string
	var result = make(map[string][]Domain)

	domains := GetAllDomains(flag)
	for _, domain := range domains {
		if groupBy == "group" {
			groupName = domain.Metadata["group"]
		} else if groupBy == "storage" {
			diskFileDir := filepath.Dir(domain.Disks[0]["file"])
			storagePool, err := Conn.LookupStoragePoolByTargetPath(diskFileDir)
			if err == nil && storagePool != nil {
				groupName, _ = storagePool.GetName()
				log.Printf("set storage pool %s as group name for %s", groupName, domain.Name)
			} else {
				log.Printf("Error: fail to get storage pool nmae for %s", domain.Name)
				groupName = "*unknown pool name"
			}
		} else if groupBy == "network" {
			if len(domain.Interfaces) >= 1 {
				intf := domain.Interfaces[0]
				groupName = intf["type"] + ":" + intf["name"]
			} else {
				log.Printf("Error: fail to get network/bridge nmae for %s", domain.Name)
				groupName = "default"
			}
			
		}
		log.Printf("+ Set %s as group name for %s", groupName, domain.Name)
		result[groupName] = append(result[groupName], domain)
	}
	return result
}

func CreateDomains(spec map[string]string) map[string][]Domain {
	var result = make(map[string][]Domain)
	var sourceDomainName, newDomainName, newDomainFile string

	log.Printf("+ Creating new domain, spec: %s", spec)
	count, _ := strconv.Atoi(spec["count"])
	for i := 1; i <= count; i++ {
		index := strconv.Itoa(i)
		if count == 1 {
			newDomainName = spec["prefix"]+"-"+spec["name"]
		} else {
			newDomainName = spec["prefix"]+"-"+spec["name"]+"-"+index
		}
		storagePool := GetStoragePool(spec["storagePool"])
		newDomainFile = storagePool.Path+"/"+newDomainName+".qcow2"
		sourceDomainName = spec["bootDiskDomain"]	

		log.Printf("  - Source Domain:%s, New Domain:%s, New Domain File:%s", sourceDomainName, newDomainName, newDomainFile)

		// clone domain 
		c := exec.Command("virt-clone", "--original", sourceDomainName, "--name", newDomainName, "--file", newDomainFile)
		_, err := c.Output()
		if err != nil {
			log.Printf("Error: fail to clone domain. %s", err)
			return result
		}

		// detach all network interface and then attach new interface
		/* 
		  assume only single interface attached, otherwise need use --mac option to remove all interface
		  and the interface is network type.
		*/
		c = exec.Command("virsh", "detach-interface", "--persistent", "--domain", newDomainName, "--type", "network")
		_, err = c.Output()
		if err != nil {
			log.Printf("Error: fail to remove interface. %s", err)
			return result
		}

		// so far support attach "network" type interface only.
		log.Printf("  - NIC Driver: %s", spec["nicDriver"])
		if spec["nic1"] != "" {
			log.Printf("  - NIC1: %s", spec["nic1"])
			c = exec.Command("virsh", "attach-interface", "--persistent", "--type", "network",
				"--domain", newDomainName, "--model", spec["nicDriver"], "--source", spec["nic1"])
			_, err = c.Output()
			if err != nil {
				log.Printf("Error: fail to add nic1. %s", err)
				return result
			}
		}
		if spec["nic2"] != "" {
			log.Printf("  - NIC2: %s", spec["nic2"])
			c = exec.Command("virsh", "attach-interface", "--persistent", "--type", "network",
				"--domain", newDomainName, "--model", spec["nicDriver"], "--source", spec["nic2"])
			_, err = c.Output()
			if err != nil {
				log.Printf("Error: fail to add nic2. %s", err)
				return result
			}
		}
		if spec["nic3"] != "" {
			log.Printf("  - NIC3: %s", spec["nic3"])
			c = exec.Command("virsh", "attach-interface", "--persistent", "--type", "network",
				"--domain", newDomainName,  "--model", spec["nicDriver"], "--source", spec["nic3"])
			_, err = c.Output()
			if err != nil {
				log.Printf("Error: fail to add nic3. %s", err)
				return result
			}
		}

		// create and attach data disk
		log.Printf("  - Disk Bus: %s", spec["diskBus"])
		if spec["disk2Size"] != "" {
			log.Printf("  - Disk2: %s", spec["disk2Size"])
		}
		if spec["disk3Size"] != "" {
			log.Printf("  - Disk3: %s", spec["disk3Size"])
		}
		if spec["disk4Size"] != "" {
			log.Printf("  - Disk4: %s", spec["disk4Size"])
		}
		//c = exec.Command()

	}

	return result
}