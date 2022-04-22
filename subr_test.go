package main

import (
	"testing"
	"unsafe"
)

func Test_uuidToStr(t *testing.T) {
	u := &uuid{}
	s := uuidToStr(u)
	if s != "00000000-0000-0000-0000-000000000000" {
		t.Error(s)
	}

	u = &uuid{0x516e7cb4, 0x6ecf, 0x11d6, 0x8f, 0xf8, [UUID_NODE_LEN]uint8{0x00, 0x02, 0x2d, 0x09, 0x71, 0x2b}}
	s = uuidToStr(u)
	if s != "516e7cb4-6ecf-11d6-8ff8-00022d09712b" {
		t.Error(s)
	}
}

func Test_knownUuid(t *testing.T) {
	n := 0
	for i, a := range knownUuid {
		for j, b := range knownUuid {
			if i == j {
				n++
				continue
			}
			if a.uuid == b.uuid {
				t.Error(a.uuid, b.uuid)
			}
			if a.name == b.name {
				t.Error(a.name, b.name)
			}
			if a.name == "" {
				t.Error(a.name)
			}
			if b.name == "" {
				t.Error(b.name)
			}
		}
	}

	if n != len(knownUuid) {
		t.Error(n)
	}
}

func Test_knownUuidToStr(t *testing.T) {
	u := &uuid{}
	s := knownUuidToStr(u)
	if s != "UNUSED" {
		t.Error(s)
	}

	u = &uuid{0x516e7cb4, 0x6ecf, 0x11d6, 0x8f, 0xf8, [UUID_NODE_LEN]uint8{0x00, 0x02, 0x2d, 0x09, 0x71, 0x2b}}
	s = knownUuidToStr(u)
	if s != "FREEBSD" {
		t.Error(s)
	}

	u = &uuid{0x416e7cb4, 0x6ecf, 0x11d6, 0x8f, 0xf8, [UUID_NODE_LEN]uint8{0x00, 0x02, 0x2d, 0x09, 0x71, 0x2b}}
	s = knownUuidToStr(u)
	if s != "" {
		t.Error(s)
	}
}

func Test_isLe(t *testing.T) {
	var x uint16 = 0x1234
	ptr := unsafe.Pointer(&x)
	b := *(*byte)(ptr)

	if b == 0x34 {
		if !isLe() {
			t.Error("!le")
		}
	} else if b == 0x12 {
		if isLe() {
			t.Error("le")
		}
	} else {
		t.Error("unknown")
	}
}

func Test_Ds(t *testing.T) {
	var a gptHdr
	if unsafe.Sizeof(a) != 92+4 {
		t.Error(a)
	}

	var b gptEnt
	if unsafe.Sizeof(b) != 128 {
		t.Error(b)
	}

	var c uuid
	if unsafe.Sizeof(c) != 16 {
		t.Error(c)
	}
}
