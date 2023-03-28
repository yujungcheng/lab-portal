package models

import (
	"log"
)

type Template struct {
	UUID string
	Name string
	StoragePoolName string
	StoragePoolUUID string
	StoragePoolPath string
	BootVolumeFile string
	DataVolume []int
}


func GetAllTemplates() []Template {
	result := make([]Template, 0)

	// get template from domain template group
	domainGroups := GetAllDomainsByGroup("persistent", "group")
	templateDomains := domainGroups["template"]
	for _, template := range templateDomains {
		log.Printf("+ Retriving template data (%s)", template.Name)
		if len(template.Disks) >= 1 {
			t := new(Template)
			t.UUID = template.UUID
			t.Name =  template.Name
			t.BootVolumeFile = template.Disks[0]["file"]
			result = append(result, *t)
		} else {
			log.Printf("  - domain %s does not have boot disk", template.Name)
			continue
		}
	}

	// todo: get template volume from template Storage pool

	return result
}