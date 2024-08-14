package main

import (
	"fmt"
)

const UUID_NODE_LEN = 6

type uuid struct {
	TimeLow               uint32
	TimeMid               uint16
	TimeHiAndVersion      uint16
	ClockSeqHiAndReserved uint8
	ClockSeqLow           uint8
	Node                  [UUID_NODE_LEN]uint8
}

func uuidToString(u *uuid) string {
	return fmt.Sprintf("%08x-%04x-%04x-%02x%02x-%02x%02x%02x%02x%02x%02x",
		u.TimeLow, u.TimeMid, u.TimeHiAndVersion,
		u.ClockSeqHiAndReserved, u.ClockSeqLow,
		u.Node[0], u.Node[1], u.Node[2], u.Node[3], u.Node[4], u.Node[5])
}
