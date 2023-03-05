package controller

import (
	"log"
	"net/http"
	"html/template"
	mod "lab-portal/models"
)

type DomainList struct {
	Domains []mod.Domain
}

type DomainController struct {
	AllDomains DomainList 
}

// list all domain
func (d DomainController) List(w http.ResponseWriter, r *http.Request) {
	log.Println("Controller - list domains")

	allDomains:= mod.GetAllDomains("persistent")
	d.AllDomains = DomainList {
		Domains: allDomains,
	}
	tplFiles := []string {
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

	log.Println("Controller - list domains completed")
}

