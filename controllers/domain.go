package controller

import (
	"html/template"
	mod "lab-portal/models"
	"log"
	"net/http"
	"strconv"
)

type DomainList struct {
	Domains        []mod.Domain
	DomainsByGroup map[string][]mod.Domain
}

type DomainClone struct {
	Templates    []mod.Template // original domain
	StoragePools []mod.StoragePool
	Networks     []mod.Network
}

type DomainController struct {
	AllDomains  DomainList
	CloneForm   DomainClone
	CloneResult map[string]string
}

func (d DomainController) List(w http.ResponseWriter, r *http.Request) {
	log.Println("Controller - list domains")

	status := r.URL.Query().Get("status")
	if status == "" {
		status = "persistent"
	}
	allDomains := mod.GetAllDomains(status)
	d.AllDomains = DomainList{
		Domains: allDomains,
	}

	tplFiles := []string{
		"templates/portal.tpl",
		"templates/base.tpl",
		"templates/domain_list.tpl",
	}
	tpl, err := template.ParseFiles(tplFiles...)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", 500)
	}
	err = tpl.Execute(w, d.AllDomains)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", 500)
	}

	//checked_at := time.Now()
	//checked_at.Format(time.RFC1123)
}

func (d DomainController) ListByGroup(w http.ResponseWriter, r *http.Request) {
	log.Println("Controller - list domains with grouping")

	var status, mode string
	var allDomainsByGroup map[string][]mod.Domain

	status = r.URL.Query().Get("status")
	if status == "" {
		status = "persistent"
	}

	mode = r.URL.Query().Get("mode")
	allDomainsByGroup = mod.GetAllDomainsByGroup(status, mode)
	d.AllDomains = DomainList{
		DomainsByGroup: allDomainsByGroup,
	}

	tplFiles := []string{
		"templates/portal.tpl",
		"templates/base.tpl",
		"templates/domain_list_by_group.tpl",
	}
	tpl, err := template.ParseFiles(tplFiles...)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", 500)
	}
	err = tpl.Execute(w, d.AllDomains)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", 500)
	}
}

func (d DomainController) GetClonePage(w http.ResponseWriter, r *http.Request) {
	log.Println("Controller - load domain clone page")

	storagePools := mod.GetAllStoragePools()
	networks := mod.GetAllNetworks()
	templates := mod.GetAllTemplates()

	d.CloneForm = DomainClone{
		StoragePools: storagePools,
		Networks:     networks,
		Templates:    templates,
	}

	tplFiles := []string{
		"templates/portal.tpl",
		"templates/base.tpl",
		"templates/domain_clone_page.tpl",
	}
	tpl, err := template.ParseFiles(tplFiles...)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", 500)
	}
	err = tpl.Execute(w, d.CloneForm)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", 500)
	}
}

func (d DomainController) Clone(w http.ResponseWriter, r *http.Request) {
	log.Println("Controller - clone domains")
	const maxGroupID = 1
	const maxInterfaceNum = 3
	const maxDiskNum = 3

	result := make(map[string]string)

	if err := r.ParseForm(); err != nil {
		log.Printf("Error: %s", err)
	} else {
		var newDomainName string
		for groupID := 1; groupID <= maxGroupID; groupID++ {
			strGroupID := strconv.Itoa(groupID)
			group := "group" + strGroupID + "-"

			original := r.PostFormValue("group1-original-domain")
			name := r.PostFormValue(group + "name")
			count := r.PostFormValue(group + "count")
			maxIDCount, _ := strconv.Atoi(count)

			prefix := r.PostFormValue(group + "prefix")
			if prefix != "" {
				prefix = prefix + "-"
			}

			storagePoolName := r.PostFormValue(group + "storage-pool")
			storagePool := mod.GetStoragePool(storagePoolName)
			storagePoolPath := storagePool.Path

			diskNames := []string{"db", "dc", "dd", "de", "df", "dg"}
			diskBus := r.PostFormValue(group + "disk-bus")
			diskTargetPrefix := "v"
			if diskBus == "ide" {
				diskTargetPrefix = "h"
			} else if diskBus == "sata" {
				diskTargetPrefix = "s"
			}

			intfDriver := r.PostFormValue("group1-nic-driver")

			var ret bool
			for nameID := 1; nameID <= maxIDCount; nameID++ {
				// set new domain name
				if maxIDCount == 1 {
					newDomainName = prefix + name
				} else {
					_nameID := strconv.Itoa(nameID)
					newDomainName = prefix + name + "-" + _nameID
				}

				// clone domain
				newDomainDiskFile := storagePoolPath + "/" + newDomainName + ".qcow2"
				log.Printf("Cloning new domain %s", newDomainName)
				ret = mod.CloneDomain(original, newDomainName, newDomainDiskFile)
				if ret != true {
					errStatus := "fail to clone domain " + newDomainName
					log.Printf("Error: %s", errStatus)
					result[newDomainName] = errStatus
					continue
				}

				// set vcpu
				vcpu := r.PostFormValue(group + "vcpu")
				if vcpu != "" {
					log.Printf("Set %s vcpu to domain %s", vcpu, newDomainName)
					ret = mod.SetDomainvCPU(newDomainName, vcpu)
					if ret != true {
						errStatus := "fail to set vcpu to " + vcpu
						log.Printf("Error: %s", errStatus)
						result[newDomainName] = errStatus
						continue
					}
				}

				// set ram size
				ram := r.PostFormValue(group + "ram")
				if ram != "" {
					log.Printf("Set %sG ram to domain %s", ram, newDomainName)
					ret = mod.SetDomainMEM(newDomainName, ram)
					if ret != true {
						errStatus := "fail to set ram to " + ram + "G"
						log.Printf("Error: %s", errStatus)
						result[newDomainName] = errStatus
						continue
					}
				}

				// attach interface to domain
				log.Printf("Detach all network interface in %s", newDomainName)
				ret = mod.DetachDomainInterface(newDomainName, "")
				if ret != true {
					errStatus := "fail to detach all interfaces in " + newDomainName
					log.Printf("Error: %s", errStatus)
					result[newDomainName] = errStatus
					continue
				}
				for intfNum := 1; intfNum <= maxInterfaceNum; intfNum++ {
					strIntfNum := strconv.Itoa(intfNum)
					intfFormNmae := group + "nic" + strIntfNum
					intfNetwork := r.PostFormValue(intfFormNmae)
					if intfNetwork != "" {
						log.Printf("Attach %s to %s", intfNetwork, newDomainName)
						ret = mod.AttachDomainInterface(newDomainName, intfDriver, intfNetwork)
						if ret != true {
							errStatus := "fail to attach interface to " + intfNetwork
							log.Printf("Error: %s", errStatus)
							result[newDomainName] = errStatus
							continue
						}
					}
				}

				// create data disk and attch to domain
				for diskNum := 1; diskNum <= maxDiskNum; diskNum++ {
					strDiskNum := strconv.Itoa(diskNum)
					diskFormName := group + "disk" + strDiskNum + "-size"
					diskSize := r.PostFormValue(diskFormName)
					if diskSize != "" {
						strDiskNum := strconv.Itoa(diskNum)
						diskName := newDomainName + ".data-disk" + strDiskNum + ".qcow2"
						log.Printf("Create %sGB data disk to domain %s",
							diskSize, newDomainName)
						ret = mod.CreateDomainDisk(storagePoolName, diskName, diskSize+"G")
						if ret != true {
							errStatus := "fail to create data disk " + diskName
							log.Printf("Error: %s", errStatus)
							result[newDomainName] = errStatus
							continue
						}
						diskPath := storagePoolPath + "/" + diskName
						diskTarget := diskTargetPrefix + diskNames[diskNum-1]
						log.Printf("Attach %s data disk to domain %s", diskTarget, newDomainName)
						ret = mod.AttachDomainDisk(newDomainName, diskPath, diskTarget)
						if ret != true {
							errStatus := "fail to attach data disk " + diskTarget
							log.Printf("Error: %s", errStatus)
							result[newDomainName] = errStatus
							continue
						}
					}
				}

				// update group in description

				log.Printf("Cloned new domain %s successfully", newDomainName)
				result[newDomainName] = "created successfully"
			}
		}
	}

	tplFiles := []string{
		"templates/portal.tpl",
		"templates/base.tpl",
		"templates/domain_clone_result.tpl",
	}
	tpl, err := template.ParseFiles(tplFiles...)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", 500)
	}
	err = tpl.Execute(w, result)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", 500)
	}
}

func (d DomainController) Delete(w http.ResponseWriter, r *http.Request) {
	domainUUID := r.URL.Query().Get("uuid")
	log.Printf("Delete domain %s", domainUUID)
}

func (d DomainController) Update(w http.ResponseWriter, r *http.Request) {
	domainUUID := r.URL.Query().Get("uuid")
	log.Printf("Update domain %s", domainUUID)
}

func (d DomainController) Backup(w http.ResponseWriter, r *http.Request) {
	domainUUID := r.URL.Query().Get("uuid")
	log.Printf("Backup domain %s", domainUUID)
}
