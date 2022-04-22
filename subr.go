package main

import (
	"unsafe"
)

func uuidToStr(uuid *uuid) string {
	return uuidToString(uuid)
}

var knownUuid = []struct {
	uuid uuid
	name string
}{
	{uuid{0x00000000, 0x0000, 0x0000, 0x00, 0x00, [UUID_NODE_LEN]uint8{0x00, 0x00, 0x00, 0x00, 0x00, 0x00}}, "UNUSED"},
	{uuid{0xc12a7328, 0xf81f, 0x11d2, 0xba, 0x4b, [UUID_NODE_LEN]uint8{0x00, 0xa0, 0xc9, 0x3e, 0xc9, 0x3b}}, "EFI"},
	{uuid{0x024dee41, 0x33e7, 0x11d3, 0x9d, 0x69, [UUID_NODE_LEN]uint8{0x00, 0x08, 0xc7, 0x81, 0xf3, 0x9f}}, "MBR"},
	{uuid{0x516e7cb4, 0x6ecf, 0x11d6, 0x8f, 0xf8, [UUID_NODE_LEN]uint8{0x00, 0x02, 0x2d, 0x09, 0x71, 0x2b}}, "FREEBSD"},
	{uuid{0x83bd6b9d, 0x7f41, 0x11dc, 0xbe, 0x0b, [UUID_NODE_LEN]uint8{0x00, 0x15, 0x60, 0xb8, 0x4f, 0x0f}}, "FREEBSD_BOOT"},
	{uuid{0x74ba7dd9, 0xa689, 0x11e1, 0xbd, 0x04, [UUID_NODE_LEN]uint8{0x00, 0xe0, 0x81, 0x28, 0x6a, 0xcf}}, "FREEBSD_NANDFS"},
	{uuid{0x516e7cb5, 0x6ecf, 0x11d6, 0x8f, 0xf8, [UUID_NODE_LEN]uint8{0x00, 0x02, 0x2d, 0x09, 0x71, 0x2b}}, "FREEBSD_SWAP"},
	{uuid{0x516e7cb6, 0x6ecf, 0x11d6, 0x8f, 0xf8, [UUID_NODE_LEN]uint8{0x00, 0x02, 0x2d, 0x09, 0x71, 0x2b}}, "FREEBSD_UFS"},
	{uuid{0x516e7cb8, 0x6ecf, 0x11d6, 0x8f, 0xf8, [UUID_NODE_LEN]uint8{0x00, 0x02, 0x2d, 0x09, 0x71, 0x2b}}, "FREEBSD_VINUM"},
	{uuid{0x516e7cba, 0x6ecf, 0x11d6, 0x8f, 0xf8, [UUID_NODE_LEN]uint8{0x00, 0x02, 0x2d, 0x09, 0x71, 0x2b}}, "FREEBSD_ZFS"},
	{uuid{0x9e1a2d38, 0xc612, 0x4316, 0xaa, 0x26, [UUID_NODE_LEN]uint8{0x8b, 0x49, 0x52, 0x1e, 0x5a, 0x8b}}, "PREP_BOOT"},
	{uuid{0xebd0a0a2, 0xb9e5, 0x4433, 0x87, 0xc0, [UUID_NODE_LEN]uint8{0x68, 0xb6, 0xb7, 0x26, 0x99, 0xc7}}, "MS_BASIC_DATA"},
	{uuid{0xaf9b60a0, 0x1431, 0x4f62, 0xbc, 0x68, [UUID_NODE_LEN]uint8{0x33, 0x11, 0x71, 0x4a, 0x69, 0xad}}, "MS_LDM_DATA"},
	{uuid{0x5808c8aa, 0x7e8f, 0x42e0, 0x85, 0xd2, [UUID_NODE_LEN]uint8{0xe1, 0xe9, 0x04, 0x34, 0xcf, 0xb3}}, "MS_LDM_METADATA"},
	{uuid{0xde94bba4, 0x06d1, 0x4d40, 0xa1, 0x6a, [UUID_NODE_LEN]uint8{0xbf, 0xd5, 0x01, 0x79, 0xd6, 0xac}}, "MS_RECOVERY"},
	{uuid{0xe3c9e316, 0x0b5c, 0x4db8, 0x81, 0x7d, [UUID_NODE_LEN]uint8{0xf9, 0x2d, 0xf0, 0x02, 0x15, 0xae}}, "MS_RESERVED"},
	{uuid{0xe75caf8f, 0xf680, 0x4cee, 0xaf, 0xa3, [UUID_NODE_LEN]uint8{0xb0, 0x01, 0xe5, 0x6e, 0xfc, 0x2d}}, "MS_SPACES"},
	{uuid{0x0fc63daf, 0x8483, 0x4772, 0x8e, 0x79, [UUID_NODE_LEN]uint8{0x3d, 0x69, 0xd8, 0x47, 0x7d, 0xe4}}, "LINUX_DATA"},
	{uuid{0xa19d880f, 0x05fc, 0x4d3b, 0xa0, 0x06, [UUID_NODE_LEN]uint8{0x74, 0x3f, 0x0f, 0x84, 0x91, 0x1e}}, "LINUX_RAID"},
	{uuid{0x0657fd6d, 0xa4ab, 0x43c4, 0x84, 0xe5, [UUID_NODE_LEN]uint8{0x09, 0x33, 0xc8, 0x4b, 0x4f, 0x4f}}, "LINUX_SWAP"},
	{uuid{0xe6d6d379, 0xf507, 0x44c2, 0xa2, 0x3c, [UUID_NODE_LEN]uint8{0x23, 0x8f, 0x2a, 0x3d, 0xf9, 0x28}}, "LINUX_LVM"},
	{uuid{0xaa31e02a, 0x400f, 0x11db, 0x95, 0x90, [UUID_NODE_LEN]uint8{0x00, 0x0c, 0x29, 0x11, 0xd1, 0xb8}}, "VMFS"},
	{uuid{0x9d275380, 0x40ad, 0x11db, 0xbf, 0x97, [UUID_NODE_LEN]uint8{0x00, 0x0c, 0x29, 0x11, 0xd1, 0xb8}}, "VMKDIAG"},
	{uuid{0x9198effc, 0x31c0, 0x11db, 0x8f, 0x78, [UUID_NODE_LEN]uint8{0x00, 0x0c, 0x29, 0x11, 0xd1, 0xb8}}, "VMRESERVED"},
	{uuid{0x381cfccc, 0x7288, 0x11e0, 0x92, 0xee, [UUID_NODE_LEN]uint8{0x00, 0x0c, 0x29, 0x11, 0xd0, 0xb2}}, "VMVSANHDR"},
	{uuid{0x426F6F74, 0x0000, 0x11aa, 0xaa, 0x11, [UUID_NODE_LEN]uint8{0x00, 0x30, 0x65, 0x43, 0xec, 0xac}}, "APPLE_BOOT"},
	{uuid{0x48465300, 0x0000, 0x11aa, 0xaa, 0x11, [UUID_NODE_LEN]uint8{0x00, 0x30, 0x65, 0x43, 0xec, 0xac}}, "APPLE_HFS"},
	{uuid{0x55465300, 0x0000, 0x11aa, 0xaa, 0x11, [UUID_NODE_LEN]uint8{0x00, 0x30, 0x65, 0x43, 0xec, 0xac}}, "APPLE_UFS"},
	{uuid{0x6a898cc3, 0x1dd2, 0x11b2, 0x99, 0xa6, [UUID_NODE_LEN]uint8{0x08, 0x00, 0x20, 0x73, 0x66, 0x31}}, "APPLE_ZFS"},
	{uuid{0x52414944, 0x0000, 0x11aa, 0xaa, 0x22, [UUID_NODE_LEN]uint8{0x00, 0x30, 0x65, 0x43, 0xec, 0xac}}, "APPLE_RAID"},
	{uuid{0x52414944, 0x5f4f, 0x11aa, 0xaa, 0x22, [UUID_NODE_LEN]uint8{0x00, 0x30, 0x65, 0x43, 0xec, 0xac}}, "APPLE_RAID_OFFLINE"},
	{uuid{0x4C616265, 0x6c00, 0x11aa, 0xaa, 0x11, [UUID_NODE_LEN]uint8{0x00, 0x30, 0x65, 0x43, 0xec, 0xac}}, "APPLE_LABEL"},
	{uuid{0x5265636f, 0x7665, 0x11AA, 0xaa, 0x11, [UUID_NODE_LEN]uint8{0x00, 0x30, 0x65, 0x43, 0xec, 0xac}}, "APPLE_TV_RECOVERY"},
	{uuid{0x53746f72, 0x6167, 0x11AA, 0xaa, 0x11, [UUID_NODE_LEN]uint8{0x00, 0x30, 0x65, 0x43, 0xec, 0xac}}, "APPLE_CORE_STORAGE"},
	{uuid{0x7c3457ef, 0x0000, 0x11aa, 0xaa, 0x11, [UUID_NODE_LEN]uint8{0x00, 0x30, 0x65, 0x43, 0xec, 0xac}}, "APPLE_APFS"},
	{uuid{0x49f48d5a, 0xb10e, 0x11dc, 0xb9, 0x9b, [UUID_NODE_LEN]uint8{0x00, 0x19, 0xd1, 0x87, 0x96, 0x48}}, "NETBSD_FFS"},
	{uuid{0x49f48d82, 0xb10e, 0x11dc, 0xb9, 0x9b, [UUID_NODE_LEN]uint8{0x00, 0x19, 0xd1, 0x87, 0x96, 0x48}}, "NETBSD_LFS"},
	{uuid{0x49f48d32, 0xb10e, 0x11dc, 0xB9, 0x9B, [UUID_NODE_LEN]uint8{0x00, 0x19, 0xd1, 0x87, 0x96, 0x48}}, "NETBSD_SWAP"},
	{uuid{0x49f48daa, 0xb10e, 0x11dc, 0xb9, 0x9b, [UUID_NODE_LEN]uint8{0x00, 0x19, 0xd1, 0x87, 0x96, 0x48}}, "NETBSD_RAID"},
	{uuid{0x2db519c4, 0xb10f, 0x11dc, 0xb9, 0x9b, [UUID_NODE_LEN]uint8{0x00, 0x19, 0xd1, 0x87, 0x96, 0x48}}, "NETBSD_CCD"},
	{uuid{0x2db519ec, 0xb10f, 0x11dc, 0xb9, 0x9b, [UUID_NODE_LEN]uint8{0x00, 0x19, 0xd1, 0x87, 0x96, 0x48}}, "NETBSD_CGD"},
	{uuid{0x9d087404, 0x1ca5, 0x11dc, 0x88, 0x17, [UUID_NODE_LEN]uint8{0x01, 0x30, 0x1b, 0xb8, 0xa9, 0xf5}}, "DRAGONFLY_LABEL32"},
	{uuid{0x9d58fdbd, 0x1ca5, 0x11dc, 0x88, 0x17, [UUID_NODE_LEN]uint8{0x01, 0x30, 0x1b, 0xb8, 0xa9, 0xf5}}, "DRAGONFLY_SWAP"},
	{uuid{0x9d94ce7c, 0x1ca5, 0x11dc, 0x88, 0x17, [UUID_NODE_LEN]uint8{0x01, 0x30, 0x1b, 0xb8, 0xa9, 0xf5}}, "DRAGONFLY_UFS1"},
	{uuid{0x9dd4478f, 0x1ca5, 0x11dc, 0x88, 0x17, [UUID_NODE_LEN]uint8{0x01, 0x30, 0x1b, 0xb8, 0xa9, 0xf5}}, "DRAGONFLY_VINUM"},
	{uuid{0xdbd5211b, 0x1ca5, 0x11dc, 0x88, 0x17, [UUID_NODE_LEN]uint8{0x01, 0x30, 0x1b, 0xb8, 0xa9, 0xf5}}, "DRAGONFLY_CCD"},
	{uuid{0x3d48ce54, 0x1d16, 0x11dc, 0x86, 0x96, [UUID_NODE_LEN]uint8{0x01, 0x30, 0x1b, 0xb8, 0xa9, 0xf5}}, "DRAGONFLY_LABEL64"},
	{uuid{0xbd215ab2, 0x1d16, 0x11dc, 0x86, 0x96, [UUID_NODE_LEN]uint8{0x01, 0x30, 0x1b, 0xb8, 0xa9, 0xf5}}, "DRAGONFLY_LEGACY"},
	{uuid{0x61dc63ac, 0x6e38, 0x11dc, 0x85, 0x13, [UUID_NODE_LEN]uint8{0x01, 0x30, 0x1b, 0xb8, 0xa9, 0xf5}}, "DRAGONFLY_HAMMER"},
	{uuid{0x5cbb9ad1, 0x862d, 0x11dc, 0xa9, 0x4d, [UUID_NODE_LEN]uint8{0x01, 0x30, 0x1b, 0xb8, 0xa9, 0xf5}}, "DRAGONFLY_HAMMER2"},
	{uuid{0xcab6e88e, 0xabf3, 0x4102, 0xa0, 0x7a, [UUID_NODE_LEN]uint8{0xd4, 0xbb, 0x9b, 0xe3, 0xc1, 0xd3}}, "CHROMEOS_FIRMWARE"},
	{uuid{0xfe3a2a5d, 0x4f32, 0x41a7, 0xb7, 0x25, [UUID_NODE_LEN]uint8{0xac, 0xcc, 0x32, 0x85, 0xa3, 0x09}}, "CHROMEOS_KERNEL"},
	{uuid{0x2e0a753d, 0x9e48, 0x43b0, 0x83, 0x37, [UUID_NODE_LEN]uint8{0xb1, 0x51, 0x92, 0xcb, 0x1b, 0x5e}}, "CHROMEOS_RESERVED"},
	{uuid{0x3cb8e202, 0x3b7e, 0x47dd, 0x8a, 0x3c, [UUID_NODE_LEN]uint8{0x7f, 0xf2, 0xa1, 0x3c, 0xfc, 0xec}}, "CHROMEOS_ROOT"},
	{uuid{0x824cc7a0, 0x36a8, 0x11e3, 0x89, 0x0a, [UUID_NODE_LEN]uint8{0x95, 0x25, 0x19, 0xad, 0x3f, 0x61}}, "OPENBSD_DATA"},
	{uuid{0x6a82cb45, 0x1dd2, 0x11b2, 0x99, 0xa6, [UUID_NODE_LEN]uint8{0x08, 0x00, 0x20, 0x73, 0x66, 0x31}}, "SOLARIS_BOOT"},
	{uuid{0x6a85cf4d, 0x1dd2, 0x11b2, 0x99, 0xa6, [UUID_NODE_LEN]uint8{0x08, 0x00, 0x20, 0x73, 0x66, 0x31}}, "SOLARIS_ROOT"},
	{uuid{0x6a87c46f, 0x1dd2, 0x11b2, 0x99, 0xa6, [UUID_NODE_LEN]uint8{0x08, 0x00, 0x20, 0x73, 0x66, 0x31}}, "SOLARIS_SWAP"},
	{uuid{0x6a8b642b, 0x1dd2, 0x11b2, 0x99, 0xa6, [UUID_NODE_LEN]uint8{0x08, 0x00, 0x20, 0x73, 0x66, 0x31}}, "SOLARIS_BACKUP"},
	{uuid{0x6a8ef2e9, 0x1dd2, 0x11b2, 0x99, 0xa6, [UUID_NODE_LEN]uint8{0x08, 0x00, 0x20, 0x73, 0x66, 0x31}}, "SOLARIS_VAR"},
	{uuid{0x6a90ba39, 0x1dd2, 0x11b2, 0x99, 0xa6, [UUID_NODE_LEN]uint8{0x08, 0x00, 0x20, 0x73, 0x66, 0x31}}, "SOLARIS_HOME"},
	{uuid{0x6a9283a5, 0x1dd2, 0x11b2, 0x99, 0xa6, [UUID_NODE_LEN]uint8{0x08, 0x00, 0x20, 0x73, 0x66, 0x31}}, "SOLARIS_ALTSEC"},
	{uuid{0x6a945a3b, 0x1dd2, 0x11b2, 0x99, 0xa6, [UUID_NODE_LEN]uint8{0x08, 0x00, 0x20, 0x73, 0x66, 0x31}}, "SOLARIS_RESERVED"},
	{uuid{0x21686148, 0x6449, 0x6e6f, 0x74, 0x4e, [UUID_NODE_LEN]uint8{0x65, 0x65, 0x64, 0x45, 0x46, 0x49}}, "BIOS_BOOT"},
}

func knownUuidToStr(uuid *uuid) string {
	assert(isLe())

	for _, x := range knownUuid {
		if *uuid == x.uuid {
			assert(len(x.name) <= 36)
			return x.name
		}
	}

	return ""
}

func isLe() bool {
	var x uint32 = 0x12345678

	ptr := unsafe.Pointer(&x)
	return *(*byte)(ptr) == 0x78
}

func assert(c bool) {
	if !c {
		panic("Assertion")
	}
}

func assertDs() {
	var a gptHdr
	assert(unsafe.Sizeof(a) == 92+4)

	var b gptEnt
	assert(unsafe.Sizeof(b) == 128)

	var c uuid
	assert(unsafe.Sizeof(c) == 16)
}