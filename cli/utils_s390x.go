// +build s390x

package main

func archConvertStatFs(cgroupFsType int) uint32 {
        return uint32(cgroupFsType)
}
