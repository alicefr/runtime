// Copyright (c) 2018 IBM
//
// SPDX-License-Identifier: Apache-2.0
//

package virtcontainers

import (
	"fmt"
	govmmQemu "github.com/intel/govmm/qemu"
	"github.com/kata-containers/runtime/virtcontainers/device/config"
	"github.com/sirupsen/logrus"
)

type qemuS390x struct {
	// inherit from qemuArchBase, overwrite methods if needed
	qemuArchBase
}

const defaultQemuPath = "/usr/bin/qemu-system-s390x"

const defaultQemuMachineType = QemuCCWVirtio

const defaultQemuMachineOptions = "accel=kvm"

const defaultPCBridgeBus = "pci.0"

const VirtioSerialS390x = "virtio-serial-ccw"

var qemuPaths = map[string]string{
	QemuCCWVirtio: defaultQemuPath,
}

var kernelRootParams = []Param{}

// TODO Verify needed parameters
var kernelParams = []Param{
	{"console", "ttysclp0"},
}

var supportedQemuMachines = []govmmQemu.Machine{
	{
		Type:    QemuCCWVirtio,
		Options: defaultQemuMachineOptions,
	},
}

// MaxQemuVCPUs returns the maximum number of vCPUs supported
func MaxQemuVCPUs() uint32 {
	return uint32(128) // TODO 255
}

func newQemuArch(config HypervisorConfig) qemuArch {
	machineType := config.HypervisorMachineType
	if machineType == "" {
		machineType = defaultQemuMachineType
	}

	q := &qemuS390x{
		qemuArchBase{
			machineType:           machineType,
			qemuPaths:             qemuPaths,
			supportedQemuMachines: supportedQemuMachines,
			kernelParamsNonDebug:  kernelParamsNonDebug,
			kernelParamsDebug:     kernelParamsDebug,
			kernelParams:          kernelParams,
		},
	}

	if config.ImagePath != "" {
		q.kernelParams = append(q.kernelParams, kernelRootParams...)
		q.kernelParamsNonDebug = append(q.kernelParamsNonDebug, kernelParamsSystemdNonDebug...)
		q.kernelParamsDebug = append(q.kernelParamsDebug, kernelParamsSystemdDebug...)
	}

	return q
}

// appendBridges appends to devices the given bridges
func (q *qemuS390x) appendBridges(devices []govmmQemu.Device, bridges []Bridge) []govmmQemu.Device {
	return genericAppendBridges(devices, bridges, q.machineType)
}

func (q *qemuS390x) appendConsole(devices []govmmQemu.Device, path string) []govmmQemu.Device {
	serial := govmmQemu.SerialDevice{
		Driver:        VirtioSerialS390x,
		ID:            "serial0",
		DisableModern: q.nestedRun,
	}

	devices = append(devices, serial)

	var console govmmQemu.CharDevice

	console = govmmQemu.CharDevice{
		Driver:   govmmQemu.Console,
		Backend:  govmmQemu.Socket,
		DeviceID: "console0",
		ID:       "charconsole0",
		Path:     path,
	}

	devices = append(devices, console)

	return devices
}

// appendVhostUserDevice throws an error if vhost devices are tried to be used.
// See issue https://github.com/kata-containers/runtime/issues/659
func (q *qemuS390x) appendVhostUserDevice(devices []govmmQemu.Device, attr config.VhostUserDeviceAttrs) []govmmQemu.Device {
	logrus.Fatalln("No vhost-user devices supported on s390x")
	return nil
}

// supportGuestMemoryHotplug return false for s390x architecture. The pc-dimm backend device for s390x
// is not support. PC-DIMM is not listed in the devices supported by qemu-system-s390x -device help
func (q *qemuS390x) supportGuestMemoryHotplug() bool {
	return false
}

func (q *qemuS390x) appendVSockPCI(devices []govmmQemu.Device, vsock kataVSOCK) []govmmQemu.Device {
	devices = append(devices,
		govmmQemu.VSOCKDevice{
			ID:        fmt.Sprintf("vsock-%d", vsock.contextID),
			ContextID: vsock.contextID,
			VHostFD:   vsock.vhostFd,
		},
	)

	return devices
}
