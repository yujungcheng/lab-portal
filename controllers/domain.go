package controller

import (
	"html/template"
	mod "lab-portal/models"
	"log"
	"net/http"
)

type DomainList struct {
	Domains        []mod.Domain
	DomainsByGroup map[string][]mod.Domain
}

type DomainCreate struct {
	StoragePools []mod.StoragePool
	Networks     []mod.Network
	Templates    []mod.Template
}

type DomainController struct {
	AllDomains DomainList
	CreateForm DomainCreate
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

func (d DomainController) GetCreatePage(w http.ResponseWriter, r *http.Request) {
	log.Println("Controller - load domain create page")

	storagePools := mod.GetAllStoragePools()
	networks := mod.GetAllNetworks()
	templates := mod.GetAllTemplates()

	d.CreateForm = DomainCreate{
		StoragePools: storagePools,
		Networks: networks,
		Templates: templates,
	}

	tplFiles := []string{
		"templates/portal.tpl",
		"templates/base.tpl",
		"templates/domain_create_page.tpl",
	}
	tpl, err := template.ParseFiles(tplFiles...)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", 500)
	}
	err = tpl.Execute(w, d.CreateForm)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", 500)
	}
}

func (d DomainController) Create(w http.ResponseWriter, r *http.Request) {
	log.Println("Controller - create domains")
	if err := r.ParseForm(); err != nil {
		log.Printf("Error: %s", err)
	} else {
		group1 := map[string]string{}
		group1["prefix"] = r.PostFormValue("group1-prefix")
		group1["name"] = r.PostFormValue("group1-name")
		group1["vcpu"] = r.PostFormValue("group1-vcpu")
		group1["ram"] = r.PostFormValue("group1-ram")
		group1["count"] = r.PostFormValue("group1-count")

		group1["diskBus"] = r.PostFormValue("group1-disk-bus")
		group1["storagePool"] = r.PostFormValue("group1-storage-pool")
		group1["bootDiskDomain"] = r.PostFormValue("group1-boot-disk-domain")
		group1["disk2Size"] = r.PostFormValue("group1-disk2-size")
		group1["disk3Size"] = r.PostFormValue("group1-disk3-size")
		group1["disk4Size"] = r.PostFormValue("group1-disk4-size")

		group1["nicDriver"] = r.PostFormValue("group1-nic-driver")
		group1["nic1"] = r.PostFormValue("group1-nic1")
		group1["nic2"] = r.PostFormValue("group1-nic2")
		group1["nic3"] = r.PostFormValue("group1-nic3")

		// clone domain
		_ = mod.CreateDomains(group1)

		//todo: update vCPU and RAM

		//todo: create data volume and attch to domain

		//todo: attach interface and attach to domain
	}

	//d.GetCreatePage(w, r)
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
