package models

import (
	"libvirt.org/go/libvirt"
	"libvirt.org/go/libvirtxml"
	"log"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
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

// todo: use json string
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
	log.Println("  Get all domains")

	result := make([]Domain, 0)

	fg := GetFlag(flag) // flag value: active, inactive, running, paused, shutoff, persistent
	domains, err := Conn.ListAllDomains(fg)
	if err != nil {
		log.Printf("Error: fail to get all domains")
	} else {
		for _, domain := range domains {
			d := GetDomain(domain)
			result = append(result, d)
			domain.Free()
		}
	}
	// sort by name.
	sort.Slice(result, func(i, j int) bool {
		return result[i].Name < result[j].Name
	})
	return result
}

func GetDomainByName(domainName string) (Domain, error) {
	log.Printf("  Get domain by name (%s)", domainName)
	var d Domain
	domain, err := Conn.LookupDomainByName(domainName)
	if err != nil {
		log.Printf("Error: faile to lookup domain %s. %s", domainName, err)
		return d, err
	}
	d = GetDomain(*domain)
	return d, err
}

func GetDomainByUUID(domainUUID string) (Domain, error) {
	log.Printf("  Get domain by UUID (%s)", domainUUID)
	var d Domain
	domain, err := Conn.LookupDomainByUUIDString(domainUUID)
	if err != nil {
		log.Printf("Error: faile to lookup domain %s. %s", domainUUID, err)
		return d, err
	}
	d = GetDomain(*domain)
	return d, err
}

func GetDomain(domain libvirt.Domain) Domain {
	var d Domain
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

	log.Printf("  - Retriving domain data (%s)", d.Name)

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
				"name":        disk.Target.Dev,
				"file":        disk.Source.File.File,
				"capacity":    diskCapacity,
				"allocation":  diskAllocation,
				"physical":    diskPhysical,
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

	return d
}

func GetAllDomainsByGroup(flag, groupBy string) map[string][]Domain {
	log.Println("  Get all domains by group", groupBy)
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
			} else {
				log.Printf("Error: fail to get storage pool nmae for %s", domain.Name)
				groupName = "*unknown"
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
		log.Printf("  - Set '%s' as group name for %s", groupName, domain.Name)
		result[groupName] = append(result[groupName], domain)
	}
	return result
}

func CloneDomain(orgDomainName, newDomainName, newDomainDiskFile string) bool {
	if Debug == true {
		log.Println("  - virt-clone", "--quiet", "--original", orgDomainName, "--name", newDomainName, "--file", newDomainDiskFile)
	}
	args := []string{"--quiet", "--original", orgDomainName, "--name", newDomainName, "--file", newDomainDiskFile}
	_, err := RunCommand("virt-clone", args...)
	if err != nil {
		return false
	} else {
		// overwrite description, todo: move to a function
		_, err = RunVirsh("desc", newDomainName, "original domain", orgDomainName)
		if err != nil {
			return false
		}
		return true
	}
}

func SetDomainvCPU(domainName, vcpuNumber string) bool {
	if Debug == true {
		log.Println("  - virsh", "setvcpus", "--current", domainName, vcpuNumber)
	}
	_, err := RunVirsh("setvcpus", "--current", domainName, vcpuNumber)
	if err != nil {
		return false
	} else {
		return true
	}
}

func SetDomainMEM(domainName, ramSize string) bool {
	intSize, _ := strconv.Atoi(ramSize)
	kbSize := intSize * 1024 * 1024 // convert ramSize GB to KB
	strSize := strconv.Itoa(kbSize)
	if Debug == true {
		log.Println("  - virsh", "setmem", "--current", domainName, strSize)
	}
	_, err := RunVirsh("setmem", "--current", domainName, strSize)
	if err != nil {
		return false
	} else {
		return true
	}
}

// this should move to volume.go
func CreateDomainDisk(diskPoolName, diskName, diskSize string) bool {
	if Debug == true {
		log.Println("  - virsh", "vol-create-as", "--format", "qcow2", "--prealloc-metadata", "--pool", diskPoolName, "--name", diskName, "--capacity", diskSize)
	}
	_, err := RunVirsh("vol-create-as", "--format", "qcow2", "--prealloc-metadata",
		"--pool", diskPoolName, "--name", diskName, "--capacity", diskSize)
	if err != nil {
		return false
	} else {
		return true
	}
}

func AttachDomainDisk(domainName, diskPath, diskTarget, diskTargetBus, diskDriverType string) bool {

	if Debug == true {
		log.Println("  - virsh", "attach-disk", "--persistent", domainName, "--source", diskPath, "--target", diskTarget, "--subdriver", diskDriverType, "--targetbus", diskTargetBus)
	}
	_, err := RunVirsh("attach-disk", "--persistent", domainName,
		"--source", diskPath, "--target", diskTarget, "--subdriver", diskDriverType, "--targetbus", diskTargetBus)
	if err != nil {
		return false
	} else {
		return true
	}
}

func DetachDomainInterface(domainName, interfaceMac string) bool {
	var err error // err := error(nil)
	if interfaceMac != "" {
		if Debug == true {
			log.Println("  - virsh", "detach-interface", "--persistent", "--type", "network", "--domain", domainName, "--mac", interfaceMac)
		}
		_, err = RunVirsh("detach-interface", "--persistent", "--type", "network",
			"--domain", domainName, "--mac", interfaceMac)
	} else {
		if Debug == true {
			log.Println("  - virsh", "detach-interface", "--persistent", "--type", "network", "--domain", domainName)
		}
		_, err = RunVirsh("detach-interface", "--persistent", "--type", "network",
			"--domain", domainName)
	}
	if err != nil {
		return false
	} else {
		return true
	}
}

func AttachDomainInterface(domainName, intfDriver, networkName string) bool {
	if Debug == true {
		log.Println("  - virsh", "attach-interface", "--persistent", "--type", "network", "--domain", domainName, "--model", intfDriver, "--source", networkName)
	}
	_, err := RunVirsh("attach-interface", "--persistent", "--type", "network",
		"--domain", domainName, "--model", intfDriver, "--source", networkName)
	if err != nil {
		return false
	} else {
		return true
	}
}
