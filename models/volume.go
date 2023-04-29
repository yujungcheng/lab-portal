package models




 func CloneVolume() bool {
    // virsh vol-clone ....
	return false
 }


 func CreateVolume(diskPoolName, diskName, diskSize string) bool {
	_, err := RunCommand("virsh", "vol-create-as", "--format", "qcow2", "--prealloc-metadata",
		"--pool", diskPoolName, "--name", diskName, "--capacity", diskSize)
	if err != nil {
		return false
	} else {
		return true
	}
}