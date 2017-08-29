package integrationkit

import (
	"strings"
)

// Platform describes an OS and CPU architecture.
type Platform struct {
	OS   OS   `json:"os"`
	Arch Arch `json:"arch"`
}

// Arch type is used to represent CPU architecture.
type Arch string

// List of supported CPU architectures.
const (
	ArchAMD64   Arch = "amd64"
	ArchARM          = "arm"
	ArchARM64        = "arm64"
	ArchPPC64LE      = "ppc64le"
	ArchS390X        = "s390x"
)

// OS type is used to represent operating system.
type OS string

// List of supported OSes.
const (
	OSLinux   OS = "linux"
	OSWindows    = "windows"
)

// NormalizeArch normalizes the architecture.
func NormalizeArch(arch string) Arch {
	arch = strings.ToLower(arch)
	var res Arch
	switch arch {
	case "x86_64", "x86-64":
		res = ArchAMD64
	case "aarch64":
		res = ArchARM64
	case "armhf", "armel":
		res = ArchARM
	}
	if res != "" {
		return res
	}
	return Arch(arch)
}

// NormalizeOS normalizes the operating system.
func NormalizeOS(os string) OS {
	os = strings.ToLower(os)
	return OS(os)
}
