// Copyright (c) 2017 Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0
//

package types

import "fmt"

// Type represents a type of bus and bridge.
type Type string

type bridge interface {
	// AddDevice on success adds the device ID to the bridge and returns the address where the device was added.
	AddDevice(ID string) (uint32, error)
	// RemoveDevice removes the device ID from the bridge.
	RemoveDevice(ID string) error
}

type Bridge struct {
	// Address contains information about devices plugged and its address in the bridge
	Address map[uint32]string

	// ID is used to identify the bridge in the hypervisor
	ID string

	// Addr is the slot of the bridge
	Addr int

	// Type is the type of the bridge (pci, pcie, etc)
	Type Type

	// BridgeMaxCapacity is the max capacity of the bridge
	BridgeMaxCapacity uint32
}

func (b *Bridge) AddDevice(ID string) (uint32, error) {
	var addr uint32

	// looking for the first available address
	for i := uint32(1); i <= b.BridgeMaxCapacity; i++ {
		if _, ok := b.Address[i]; !ok {
			addr = i
			break
		}
	}

	if addr == 0 {
		return 0, fmt.Errorf("Unable to hot plug device on bridge: there are not empty slots")
	}

	// save address and device
	b.Address[addr] = ID
	return addr, nil
}

func (b *Bridge) RemoveDevice(ID string) error {
	// check if the device was hot plugged in the bridge
	for addr, devID := range b.Address {
		if devID == ID {
			// free address to re-use the same slot with other devices
			delete(b.Address, addr)
			return nil
		}
	}

	return fmt.Errorf("Unable to hot unplug device %s: not present on bridge", ID)
}

const pciBridgeMaxCapacity = 30

const (
	// PCI represents a PCI bus and bridge
	PCI Type = "pci"

	// PCIE represents a PCIe bus and bridge
	PCIE Type = "pcie"
)

const CCWBridgeMaxCapacity = 0xffff

const (
	CCW Type = "ccw"
)

// AddressFormatCCW returns the address format for the device number. The channel subsystem-ID 0xfe is reserved to the virtual channel and the address format is in the form fe.n.dddd, where n is subchannel set ID and ddd the device number. More details at https://www.ibm.com/support/knowledgecenter/en/linuxonibm/com.ibm.linux.z.ldva/ldva_t_configuringSCSIdevices.html
func (b *Bridge) AddressFormatCCW(addr string) (string, error) {
	if b.Type != CCW {
		return "", fmt.Errorf("Wrong bridge type")
	}
	return fmt.Sprintf("fe.%x.%s", b.Addr, addr), nil
}

// AddressFormatCCWForVirtServer returns the address format for the virtual server. The address format is in the form of 0.n.dddd
func (b *Bridge) AddressFormatCCWForVirtServer(addr string) (string, error) {
	if b.Type != CCW {
		return "", fmt.Errorf("Wrong bridge type")
	}
	return fmt.Sprintf("0.%x.%s", b.Addr, addr), nil
}
