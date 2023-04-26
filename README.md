# VM portal for home lab

A web portal to simplify VM management for my home lab.

OS: ubuntu22.04
Go version: go1.18.1


#### Initial project dir and install go libvirt
```
$ mkdir lab-portal; cd lab-portal
$ go mod init lab-portal
$ go mod tidy
```

libvirt and libvirtxml (published on 2023)
```
$ go get libvirt.org/go/libvirt
$ go get libvirt.org/go/libvirtxml
$ go install libvirt.org/go/libvirt@latest
$ go install libvirt.org/go/libvirtxml@latest
```

libvirt-go and libvirt-go-xml (published on 2021, replaced by libvirt and libvirtxml)
```
$ go get libvirt.org/libvirt-go
$ go install libvirt.org/libvirt-go

$ go get libvirt.org/libvirt-go-xml
$ go install libvirt.org/libvirt-go-xml
```

To format the code
```
$ gofmt -s -w .

OR

$ go fmt ...
```

#### Run
Run command `go run main.go` and then open web browser on port 3000.


#### Reference
```
https://go.dev/doc/gopath_code  
https://go.dev/ref/mod  
https://libvirt.org/go/libvirt.html  

https://pkg.go.dev/github.com/libvirt/libvirt-go
https://pkg.go.dev/libvirt.org/libvirt-go-xml

https://pkg.go.dev/libvirt.org/go/libvirt
https://pkg.go.dev/libvirt.org/go/libvirtxml

https://www.cyberciti.biz/faq/how-to-add-disk-image-to-kvm-virtual-machine-with-virsh-command/
```

```
https://access.redhat.com/documentation/zh-tw/red_hat_enterprise_linux/6/html/virtualization_administration_guide/sect-virtualization-storage_pools-creating-local_directories-virsh#doc-wrapper
https://access.redhat.com/documentation/en-us/red_hat_enterprise_linux/7/html/virtualization_deployment_and_administration_guide/sect-managing_guest_virtual_machines_with_virsh-storage_pool_commands
https://access.redhat.com/documentation/en-us/red_hat_enterprise_linux/7/html/virtualization_deployment_and_administration_guide/sect-managing_guest_virtual_machines_with_virsh-storage_pool_commands#sect-virsh-pool-list
```


#### KB
Before using libvirtxml, I use "gopkg.in/xmlpath.v2" to parser XML.
```
Use "gopkg.in/xmlpath.v2" package to parser xml
func ParserXML(xml string, xpath string) []string {
	path := xmlpath.MustCompile(xpath)
	root, err := xmlpath.Parse(strings.NewReader(xml))
	var result []string
	if err == nil {
		values := path.Iter(root)
		for values.Next() {
			node := values.Node()
			result = append(result, node.String())
		}
		/* return first match as string
		  if value, ok := path.String(root); ok {
			return value
		  }
		*/
	}
	return result
}
```