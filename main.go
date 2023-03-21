package main

import (
	ctl "lab-portal/controllers"
	mod "lab-portal/models"
	"log"
	"net/http"
)

func main() {

	conn := mod.GetLibvirtConnect()
	defer conn.Close()

	mod.SetLibvirtConnect(conn)
	mod.SetStartTime()
	mod.SetProcessID()

	hostname := mod.GetHostname()
	log.Printf("Start Lab-Portal on %s", hostname)

	/* Start HTTP service
	Initialize controller and define path */
	mux := http.NewServeMux()

	var domainCtl = ctl.DomainController{}
	mux.HandleFunc("/", domainCtl.List)
	mux.HandleFunc("/domains", domainCtl.List)
	mux.HandleFunc("/domains/list", domainCtl.List)
	mux.HandleFunc("/domains/list-by-group", domainCtl.ListByGroup)

	/*
		mux.HandleFunc("/domains/create-page", domainCtl.Create)
		mux.HandleFunc("/domains/delete-page", domainCtl.Create)
		mux.HandleFunc("/domains/update-page", domainCtl.Create)
		mux.HandleFunc("/domains/create", domainCtl.Create)
		mux.HandleFunc("/domains/delete", domainCtl.Create)
		mux.HandleFunc("/domains/update", domainCtl.Create)
	*/

	l := http.ListenAndServe(":3000", mux) // listen on port 3000
	log.Printf("Close Lab-Portal. %s", l)
}
