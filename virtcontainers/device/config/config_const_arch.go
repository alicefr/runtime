// +build !s390x

package config


const (
	// DeviceVFIO is the VFIO device type
	DeviceVFIO DeviceType = "vfio"

	// DeviceBlock is the block device type
	DeviceBlock DeviceType = "block"

	// DeviceGeneric is a generic device type
	DeviceGeneric DeviceType = "generic"

	//VhostUserSCSI - SCSI based vhost-user type
	VhostUserSCSI = "vhost-user-scsi-pci"

	//VhostUserNet - net based vhost-user type
	VhostUserNet = "virtio-net-pci"

	//VhostUserBlk represents a block vhostuser device type
	VhostUserBlk = "vhost-user-blk-pci"
)
