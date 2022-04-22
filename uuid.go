package main

import (
	"fmt"
)

const UUID_NODE_LEN = 6

type uuid struct {
	Time_low                  uint32
	Time_mid                  uint16
	Time_hi_and_version       uint16
	Clock_seq_hi_and_reserved uint8
	Clock_seq_low             uint8
	Node                      [UUID_NODE_LEN]uint8
}

func uuidToString(u *uuid) string {
	return fmt.Sprintf("%08x-%04x-%04x-%02x%02x-%02x%02x%02x%02x%02x%02x",
		u.Time_low, u.Time_mid, u.Time_hi_and_version,
		u.Clock_seq_hi_and_reserved, u.Clock_seq_low,
		u.Node[0], u.Node[1], u.Node[2], u.Node[3], u.Node[4], u.Node[5])
}
