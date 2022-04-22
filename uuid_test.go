package main

import (
	"testing"
)

func Test_uuidToString(t *testing.T) {
	u := &uuid{}
	s := uuidToString(u)
	if s != "00000000-0000-0000-0000-000000000000" {
		t.Error(s)
	}

	u = &uuid{0x516e7cb4, 0x6ecf, 0x11d6, 0x8f, 0xf8, [UUID_NODE_LEN]uint8{0x00, 0x02, 0x2d, 0x09, 0x71, 0x2b}}
	s = uuidToStr(u)
	if s != "516e7cb4-6ecf-11d6-8ff8-00022d09712b" {
		t.Error(s)
	}
}
