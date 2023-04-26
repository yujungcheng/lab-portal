package models

import (
	"libvirt.org/go/libvirt"
	"libvirt.org/go/libvirtxml"
	"log"
	"sort"
)

type StoragePool struct {
	Name         string
	UUID         string
	IsActive     bool
	NumOfVolumes int

	State      string
	Capacity   uint64
	Allocation uint64
	Available  uint64

	Path string
}

func GetStoragePoolStateStr(state libvirt.StoragePoolState) string {
	var stateStr string = "Unknown"
	switch state {
	case 0:
		stateStr = "Inactive"
	case 1:
		stateStr = "Building"
	case 2:
		stateStr = "Running"
	case 3:
		stateStr = "Degraded"
	case 4:
		stateStr = "Inaccessible"
	}
	return stateStr
}

/* ------------------------------------------------------------------------ */

func GetAllStoragePools() []StoragePool {
	log.Println("  Get all storage pools")
	result := make([]StoragePool, 0)

	pools, err := Conn.ListAllStoragePools(0)
	if err != nil {
		log.Printf("Error: fail to get all storage pools")
	} else {
		for _, pool := range pools {
			s := new(StoragePool)
			s.Name, _ = pool.GetName()
			s.UUID, _ = pool.GetUUIDString()
			s.IsActive, _ = pool.IsActive()
			s.NumOfVolumes, _ = pool.NumOfStorageVolumes()

			info, _ := pool.GetInfo()
			s.State = GetStoragePoolStateStr(info.State)
			s.Capacity = info.Capacity
			s.Allocation = info.Allocation
			s.Available = info.Available

			log.Printf("  - Retriving storage pool data (%s)", s.Name)

			storagePoolxml, _ := pool.GetXMLDesc(0)
			storagePoolcfg := &libvirtxml.StoragePool{}
			_ = storagePoolcfg.Unmarshal(storagePoolxml)
			s.Path = storagePoolcfg.Target.Path

			result = append(result, *s)
		}
	}

	// sort by name.
	sort.Slice(result, func(i, j int) bool {
		return result[i].Name < result[j].Name
	})

	return result
}

func GetStoragePool(poolName string) StoragePool {
	log.Println("  Get storage pool", poolName)
	pool, _ := Conn.LookupStoragePoolByName(poolName)
	s := new(StoragePool)
	s.Name, _ = pool.GetName()
	s.UUID, _ = pool.GetUUIDString()
	s.IsActive, _ = pool.IsActive()
	s.NumOfVolumes, _ = pool.NumOfStorageVolumes()

	info, _ := pool.GetInfo()
	s.State = GetStoragePoolStateStr(info.State)
	s.Capacity = info.Capacity
	s.Allocation = info.Allocation
	s.Available = info.Available

	log.Printf("  - Retriving storage pool data (%s)", s.Name)

	storagePoolxml, _ := pool.GetXMLDesc(0)
	storagePoolcfg := &libvirtxml.StoragePool{}
	_ = storagePoolcfg.Unmarshal(storagePoolxml)
	s.Path = storagePoolcfg.Target.Path
	return *s
}

func CreateStoragePool(poolName, poolPath string) bool {
	_, err := RunCommand("virsh", "pool-create-as", "--type", "dir",
		"--name", poolName, "--target", poolPath)
	if err != nil {
		return false
	} else {
		return true
	}
}

func DeleteStoragePool(poolName string) bool {
	if Debug == true {
	}

	return false
}
