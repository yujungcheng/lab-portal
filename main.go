package main

import (
	"log"
	"net/http"
	ctl "lab-portal/controllers"
	mod "lab-portal/models"
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
	mux.HandleFunc("/domains", domainCtl.List)

	
	l := http.ListenAndServe(":3000", mux)  // listen on port 3000
	log.Printf("Close Lab-Portal. %s", l)
}