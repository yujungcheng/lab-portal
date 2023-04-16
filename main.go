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
	mod.SetDebug(true)

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
	mux.HandleFunc("/domains/clone-page", domainCtl.GetClonePage)
	mux.HandleFunc("/domains/clone", domainCtl.Clone)

	/*
		mux.HandleFunc("/domains/delete-page", domainCtl.GetDeletePage)
		mux.HandleFunc("/domains/delete", domainCtl.Delete)
		mux.HandleFunc("/domains/update-page", domainCtl.GetUpdatePage)
		mux.HandleFunc("/domains/update", domainCtl.Update)

		mux.HandleFunc("/storagePool/list", poolCtl.List)
		mux.HandleFunc("/network/list", networkCtl.List)
	*/

	l := http.ListenAndServe(":3600", mux) // listen on port 3000
	log.Printf("Close Lab-Portal. %s", l)
}
