# Copyright (c) 2018 IBM
#
# SPDX-License-Identifier: Apache-2.0
#

# Power s390x settings

MACHINETYPE := s390-ccw-virtio
KERNELPARAMS :=
MACHINEACCELERATORS :=
# KERNELTYPE := uncompressed #This architecture must use an uncompressed kernel.
QEMUCMD := qemu-system-s390x
