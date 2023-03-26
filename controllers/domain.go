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
	reqData := r.URL.Query()
	_ = reqData
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
