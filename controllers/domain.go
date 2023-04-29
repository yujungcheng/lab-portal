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
	DomainNames    map[string]string // store name and uuid only
}

type DomainClonePage struct {
	Templates    []mod.Template // original domain to be cloned
	StoragePools []mod.StoragePool
	Networks     []mod.Network

	// set cpu, ram, count, disk, interfaces
	Groups     map[string]string
	VCPUs      map[string]string
	RAMs       map[string]string
	Count      map[string]string
	Disks      map[string]string
	Interfaces map[string]string
}

type DomainCloneResult struct {
	CloneResult map[string]string
}

type DomainController struct {
	AllDomains  DomainList
	ClonePage   DomainClonePage
	CloneResult DomainCloneResult
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

	d.ClonePage = DomainClonePage{
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
	err = tpl.Execute(w, d.ClonePage)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", 500)
	}
}

func (d DomainController) Clone(w http.ResponseWriter, r *http.Request) {
	log.Println("Controller - clone domains")

	const maxDomainGroupCount = 2
	const maxInterfaceCount = 3
	const maxDiskCount = 3

	domainCloneResult := make(map[string]string)

	if err := r.ParseForm(); err != nil {
		log.Printf("Error: %s", err)
	} else {
		var ret bool
		var errStatus, strGroupID, group string
		var originalDomainName, newDomainNamePrefix, newDomainName string
		var storagePoolName, storagePoolPath string

		for groupID := 1; groupID <= maxDomainGroupCount; groupID++ {
			strGroupID = strconv.Itoa(groupID)
			group = "group" + strGroupID + "-"

			originalDomainName = r.PostFormValue(group + "original-domain-name")

			newDomainNamePrefix = r.PostFormValue(group + "prefix")
			if newDomainNamePrefix != "" {
				newDomainNamePrefix = newDomainNamePrefix + "-"
			}

			domainName := r.PostFormValue(group + "name")
			if domainName == "" {
				errStatus = "domain name empty"
				log.Printf("Error: %s", errStatus)
				domainCloneResult["group"+strGroupID] = errStatus
				continue
			}

			// initial storage pool
			newStoragePoolName := r.PostFormValue(group + "new-storage-pool")
			newStoragePoolPath := r.PostFormValue(group + "new-storage-pool-path")
			if newStoragePoolName != "" && newStoragePoolPath != "" {
				// create new storage pool
				ret = mod.CreateStoragePool(newStoragePoolName, newStoragePoolPath)
				if ret != true {
					errStatus = "fail to create new storage pool " + newStoragePoolName
					log.Printf("Error: %s", errStatus)
					domainCloneResult["group"+strGroupID] = errStatus
					continue
				}
				storagePoolName = newStoragePoolName
				storagePoolPath = newStoragePoolPath
			} else {
				storagePoolName = r.PostFormValue(group + "storage-pool")
				storagePool := mod.GetStoragePool(storagePoolName)
				storagePoolPath = storagePool.Path
			}

			// initial data disk naming
			diskNames := []string{"db", "dc", "dd", "de", "df", "dg"}
			diskTargetBus := r.PostFormValue(group + "disk-bus")
			diskTargetPrefix := "v"
			if diskTargetBus == "ide" {
				diskTargetPrefix = "h"
			} else if diskTargetBus == "sata" {
				diskTargetPrefix = "s"
			}
			diskDriverType := "qcow2" // or raw

			// set network interface driver
			intfDriver := r.PostFormValue("group1-nic-driver")

			domainCount := r.PostFormValue(group + "count")
			maxDomainCount, _ := strconv.Atoi(domainCount)
			for nameID := 1; nameID <= maxDomainCount; nameID++ {
				// set new domain name
				if maxDomainCount == 1 {
					newDomainName = newDomainNamePrefix + domainName
				} else {
					_nameID := strconv.Itoa(nameID)
					newDomainName = newDomainNamePrefix + domainName + "-" + _nameID
				}

				// clone domain
				newDomainDiskFile := storagePoolPath + "/" + newDomainName + ".qcow2"
				log.Printf("Cloning new domain %s", newDomainName)
				ret = mod.CloneDomain(originalDomainName, newDomainName, newDomainDiskFile)
				if ret != true {
					errStatus := "fail to clone domain " + newDomainName
					log.Printf("Error: %s", errStatus)
					domainCloneResult[newDomainName] = errStatus
					continue
				}
				log.Printf("Clone new domain %s successfully", newDomainName)

				// update description
				desc := "original_domain="+newDomainName
				ret = mod.SetDomainDesc(newDomainName, desc)
				if ret != true {
					errStatus := "fail to set description"
					log.Printf("Error: %s", errStatus)
					domainCloneResult[newDomainName] = errStatus
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
						domainCloneResult[newDomainName] = errStatus
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
						domainCloneResult[newDomainName] = errStatus
						continue
					}
				}

				// todo: set boot disk driver, bus

				// create data disk and attch to domain
				for diskNum := 1; diskNum <= maxDiskCount; diskNum++ {
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
							domainCloneResult[newDomainName] = errStatus
							continue
						}
						diskPath := storagePoolPath + "/" + diskName
						diskTarget := diskTargetPrefix + diskNames[diskNum-1]
						log.Printf("Attach data disk %s to domain %s", diskTarget, newDomainName)
						ret = mod.AttachDomainDisk(newDomainName, diskPath, diskTarget, diskTargetBus, diskDriverType)
						if ret != true {
							errStatus := "fail to attach data disk " + diskTarget
							log.Printf("Error: %s", errStatus)
							domainCloneResult[newDomainName] = errStatus
							continue
						}
					}
				}

				// attach network interface to domain
				log.Printf("Detach all network interface in %s", newDomainName)
				ret = mod.DetachDomainInterface(newDomainName, "")
				if ret != true {
					errStatus := "fail to detach all interfaces in " + newDomainName
					log.Printf("Error: %s", errStatus)
					domainCloneResult[newDomainName] = errStatus
					continue
				}
				for intfNum := 1; intfNum <= maxInterfaceCount; intfNum++ {
					strIntfNum := strconv.Itoa(intfNum)
					intfFormNmae := group + "nic" + strIntfNum
					intfNetwork := r.PostFormValue(intfFormNmae)
					if intfNetwork != "" {
						log.Printf("Attach %s network to %s", intfNetwork, newDomainName)
						ret = mod.AttachDomainInterface(newDomainName, intfDriver, intfNetwork)
						if ret != true {
							errStatus := "fail to attach interface to " + intfNetwork
							log.Printf("Error: %s", errStatus)
							domainCloneResult[newDomainName] = errStatus
							continue
						}
					}
				}

				domainCloneResult[newDomainName] = "created successfully"
			}
		}
	}

	// todo: reuse domain list struct to show result

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
	err = tpl.Execute(w, domainCloneResult)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", 500)
	}
}

func (d DomainController) Show(w http.ResponseWriter, r *http.Request) {
	domainUUID := r.URL.Query().Get("uuid")
	log.Printf("Controller - show domain %s", domainUUID)

	domain, err := mod.GetDomainByUUID(domainUUID)
    if err != nil {
		log.Printf("Error: %s", err)
	}
	tplFiles := []string{
		"templates/portal.tpl",
		"templates/base.tpl",
		"templates/domain_show_page.tpl",
	}
	tpl, err := template.ParseFiles(tplFiles...)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", 500)
	}
	err = tpl.Execute(w, domain)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", 500)
	}
}

func (d DomainController) Delete(w http.ResponseWriter, r *http.Request) {
	domainUUID := r.URL.Query().Get("uuid")
	log.Printf("Controller - delete domain %s", domainUUID)
}

func (d DomainController) Update(w http.ResponseWriter, r *http.Request) {
	domainUUID := r.URL.Query().Get("uuid")
	log.Printf("Controller - update domain %s", domainUUID)

	var ret bool

	domainName := r.PostFormValue("domainName")
	domainVcpu := r.PostFormValue("domainVcpu")
	domainMem := r.PostFormValue("domainMem")
	newDomainName := r.PostFormValue("newDomainName")
	newDomainVcpu := r.PostFormValue("newDomainVcpu")
	newDomainMem := r.PostFormValue("newDomainMem")
	
	if domainName != newDomainName {
		ret = mod.SetDomainName(domainName, newDomainName)
		if ret != true {
			log.Printf("Error: fail to update domain name to %s", newDomainName)
		}
	}
	if domainVcpu != newDomainVcpu {
		ret = mod.SetDomainvCPU(domainName, newDomainVcpu)
		if ret != true {
			log.Printf("Error: fail to update vcpu to %s", newDomainVcpu)
		}
	}
	if domainMem != newDomainMem {
		ret = mod.SetDomainMEM(domainName, newDomainMem)
		if ret != true {
			log.Printf("Error: fail to update memory to %s", newDomainMem)
		}
	}

	domain, err := mod.GetDomainByUUID(domainUUID)
    if err != nil {
		log.Printf("Error: %s", err)
	}
	tplFiles := []string{
		"templates/portal.tpl",
		"templates/base.tpl",
		"templates/domain_show_page.tpl",
	}
	tpl, err := template.ParseFiles(tplFiles...)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", 500)
	}
	err = tpl.Execute(w, domain)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", 500)
	}
}

func (d DomainController) Backup(w http.ResponseWriter, r *http.Request) {
	domainUUID := r.URL.Query().Get("uuid")
	log.Printf("Controller - backup domain %s", domainUUID)
}
