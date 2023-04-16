package models

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"time"
	"strconv"
	"strings"
	"gopkg.in/xmlpath.v2"
	"libvirt.org/go/libvirt"
)

/* variables for model
------------------------------------------------------------- */
var StartTime time.Time
var ProcessID int
var Conn libvirt.Connect
var Debug bool

func SetStartTime() { StartTime = time.Now() }

func SetProcessID() { ProcessID = os.Getpid() }

func SetLibvirtConnect(conn libvirt.Connect) { Conn = conn }

func SetDebug(enabled bool) { Debug = enabled }

/* functions
------------------------------------------------------------- */
func GetLibvirtConnect() libvirt.Connect {
	Conn, err := libvirt.NewConnect("qemu:///system")
	if err != nil {
		log.Fatalf("Fatal: Fail to create libvirt connection.", err)
	}
	return *Conn
}

func GetElapsedTime(startTime time.Time) time.Duration {
	return time.Since(StartTime)
}

func GetHostname() string {
	var conn = GetLibvirtConnect()
	defer conn.Close()

	hostname, err := conn.GetHostname()
	if err != nil {
		log.Fatalf("Fatal: fail to get hostname")
	}
	return hostname
}

func GetLibvirtVersion() string {
	var conn = GetLibvirtConnect()
	defer conn.Close()

	libvirtVersion, err := conn.GetLibVersion()
	if err != nil {
		log.Fatalf("Fatal: fail to get libvirt version")
	}
	major := libvirtVersion / 1000000
	minor := (libvirtVersion - (major * 1000000)) / 1000
	release := (libvirtVersion - (major * 1000000) - (minor * 1000))
	return fmt.Sprintf("%d.%d.%d", major, minor, release)
	//return strconv.FormatUint(uint64(libvirtVersion), 10)
}

// unused soon
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

func ConvertSizeToString(size uint64, unit string) string {
	var newSize uint64
	switch unit {
	case "KB":
		newSize = size / 1024
	case "MB":
		newSize = size / 1048576
	case "GB":
		newSize = size / 1073741824
	default:
		newSize = size / 1024
	}
	return strconv.FormatUint(newSize, 10) + unit
}

func ConvertNetmaskToNumber(mask string) int {
	ip := net.ParseIP(mask)
	sz, _ := net.IPMask(ip.To4()).Size()
	return sz
}

func RunVirsh(args... string) (string, error) {
	c := exec.Command("virsh", args...)
	out, err := c.Output()
	/*
	if err != nil {
		log.Printf("    RunVirsh Error: %s", err)
	} else {
		log.Printf("    RunVirsh Output: %s", out)
	}
	*/	
	return string(out), err
}

func RunCommand(cmd string, args... string) (string, error) {
	c := exec.Command(cmd, args...)
	out, err := c.Output()
	/*
	if err != nil {
		log.Printf("    RunCommand Error: %s", err)
	} else {
		log.Printf("    RunCommand Output: %s", out)
	}
	*/
	return string(out), err
}