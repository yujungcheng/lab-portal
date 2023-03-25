package models

import (
	"log"
	"sort"
	"libvirt.org/go/libvirt"
	"libvirt.org/go/libvirtxml"
)

type Interface struct {
	Name string
	UUID string
	MacAddress string
	Owner string
	OwnerUUID string

	IpAddress string
	IpNetmask string

	Parameters *libvirt.NetworkPortParameters
}

type Network struct {
	Name string
	UUID string
	Bridge string
	IsActive bool	// State
	IsPersistent bool
	AutoStart bool
    ForwardMode string
	DhcpRange []string
	DhcpHosts []string    
}

/* ------------------------------------------------------------------------ */

func GetAllNetworks() []Network {
	result := make([]Network, 0)

	fg := libvirt.CONNECT_LIST_NETWORKS_PERSISTENT
	networks, err := Conn.ListAllNetworks(fg)
	if err != nil {
		log.Println("Error: fail to get all networks")
	} else {
		for _, network := range networks {
			n := new(Network)
			n.Name, _ = network.GetName()
			n.Bridge, _ = network.GetBridgeName()
			n.UUID, _ = network.GetUUIDString()
			n.AutoStart, _ = network.GetAutostart()
			n.IsActive, _ = network.IsActive()
			n.IsPersistent, _ = network.IsPersistent()
			
			log.Printf("+ Retriving network data (%s)", n.Name)

			networkxml, _ := network.GetXMLDesc(0)
			networkcfg := &libvirtxml.Network{}
			_ = networkcfg.Unmarshal(networkxml)
			if networkcfg.Forward != nil {
				n.ForwardMode = networkcfg.Forward.Mode
			}
			for _, ip := range networkcfg.IPs {
				//address := ip.Address
				//netmask := ip.Netmask
				if ip.DHCP != nil {
					for _, _dhcpRange := range ip.DHCP.Ranges {
						_range := _dhcpRange.Start+"-"+_dhcpRange.End
						n.DhcpRange = append(n.DhcpRange, _range)
					}
					for _, _dhcpHost := range ip.DHCP.Hosts {
						_host := _dhcpHost.Name+" "+_dhcpHost.MAC+" "+_dhcpHost.IP
						n.DhcpHosts = append(n.DhcpHosts, _host)
					}
				}
				//log.Printf("DHCP: %s %s %s %s", address, netmask, dhcpRange, dhcpHosts)
			}
			result = append(result, *n)	
		}
	}

	// sort by name.
	sort.Slice(result, func(i, j int) bool {
		return result[i].Name < result[j].Name
	})
	return result
}