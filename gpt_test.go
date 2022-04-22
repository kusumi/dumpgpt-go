package main

import (
	"testing"
)

func Test_tryKnownUuidToStr(t *testing.T) {
	u := &uuid{0x516e7cb4, 0x6ecf, 0x11d6, 0x8f, 0xf8, [UUID_NODE_LEN]uint8{0x00, 0x02, 0x2d, 0x09, 0x71, 0x2b}}
	s := knownUuidToStr(u)
	if s != "FREEBSD" {
		t.Error(s)
	}

	u = &uuid{0x416e7cb4, 0x6ecf, 0x11d6, 0x8f, 0xf8, [UUID_NODE_LEN]uint8{0x00, 0x02, 0x2d, 0x09, 0x71, 0x2b}}
	s = tryKnownUuidToStr(u)
	if s != "416e7cb4-6ecf-11d6-8ff8-00022d09712b" {
		t.Error(s)
	}
}

func Test_allocBuffer(t *testing.T) {
	buf := allocBuffer()
	var i interface{} = buf

	switch i.(type) {
	case []byte:
	default:
		t.Error(buf)
	}

	if len(buf) < 512 {
		t.Error(len(buf))
	}

	if len(buf)%512 != 0 {
		t.Error(len(buf))
	}

	for i := 0; i < len(buf); i++ {
		if buf[i] != 0 {
			t.Error(i, buf[i])
		}
	}
}
