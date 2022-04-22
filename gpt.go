package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"unsafe"
)

type gptHdr struct {
	Hdr_sig       [8]uint8
	Hdr_revision  uint32
	Hdr_size      uint32
	Hdr_crc_self  uint32
	Reserved      uint32
	Hdr_lba_self  uint64
	Hdr_lba_alt   uint64
	Hdr_lba_start uint64
	Hdr_lba_end   uint64
	Hdr_uuid      uuid
	Hdr_lba_table uint64
	Hdr_entries   uint32
	Hdr_entsz     uint32
	Hdr_crc_table uint32
	Padding       uint32
}

type gptEnt struct {
	Ent_type      uuid
	Ent_uuid      uuid
	Ent_lba_start uint64
	Ent_lba_end   uint64
	Ent_attr      uint64
	Ent_name      [36]uint16
}

func tryKnownUuidToStr(uuid *uuid) string {
	var s string

	if dumpOptSymbol {
		s = knownUuidToStr(uuid)
		if s == "" {
			s = uuidToStr(uuid)
		}
	} else {
		s = uuidToStr(uuid)
	}

	return s
}

func allocBuffer() []byte {
	const UNIT_SIZE = 512

	buf := make([]byte, UNIT_SIZE)

	assert(len(buf) == UNIT_SIZE)
	assert(len(buf)%512 == 0)

	return buf
}

func dumpHeader(fp *os.File, hdr_lba uint64) (*gptHdr, error) {
	buf := allocBuffer()

	hdr_offset := hdr_lba * uint64(len(buf))

	ret, err := fp.ReadAt(buf, int64(hdr_offset))
	if err != nil {
		return nil, err
	} else if ret != len(buf) {
		return nil, errors.New("failed to read")
	}

	hdr := gptHdr{}
	n := int(unsafe.Sizeof(hdr))
	r := bytes.NewReader(buf[:n])
	binary.Read(r, binary.LittleEndian, &hdr)

	l1 := []byte("EFI PART")
	l2 := []byte(hdr.Hdr_sig[:])
	if !bytes.Equal(l1, l2) {
		return nil, errors.New("not GPT")
	}

	fmt.Printf("sig      = \"%c%c%c%c%c%c%c%c\"\n",
		hdr.Hdr_sig[0],
		hdr.Hdr_sig[1],
		hdr.Hdr_sig[2],
		hdr.Hdr_sig[3],
		hdr.Hdr_sig[4],
		hdr.Hdr_sig[5],
		hdr.Hdr_sig[6],
		hdr.Hdr_sig[7])

	p := make([]byte, 4)
	binary.LittleEndian.PutUint32(p, hdr.Hdr_revision)
	fmt.Printf("revision = %02x %02x %02x %02x\n",
		p[0], p[1], p[2], p[3])

	fmt.Printf("size     = %d\n", hdr.Hdr_size)
	fmt.Printf("crc_self = 0x%x\n", hdr.Hdr_crc_self)
	fmt.Printf("lba_self = 0x%016x\n", hdr.Hdr_lba_self)
	fmt.Printf("lba_alt  = 0x%016x\n", hdr.Hdr_lba_alt)
	fmt.Printf("lba_start= 0x%016x\n", hdr.Hdr_lba_start)
	fmt.Printf("lba_end  = 0x%016x\n", hdr.Hdr_lba_end)

	fmt.Printf("uuid     = %s\n", tryKnownUuidToStr(&hdr.Hdr_uuid))

	fmt.Printf("lba_table= 0x%016x\n", hdr.Hdr_lba_table)
	fmt.Printf("entries  = %d\n", hdr.Hdr_entries)
	fmt.Printf("entsz    = %d\n", hdr.Hdr_entsz)
	fmt.Printf("crc_table= 0x%x\n", hdr.Hdr_crc_table)

	// XXX
	if hdr.Hdr_entries > 512 {
		return nil, errors.New("likely corrupted entries")
	}

	return &hdr, nil
}

func dumpEntries(fp *os.File, hdr *gptHdr) error {
	buf := allocBuffer()

	lba_table_size := hdr.Hdr_entsz * hdr.Hdr_entries
	lba_table_sectors := lba_table_size / uint32(len(buf))
	total := 0

	fmt.Printf("%-3s %-36s %-36s %-16s %-16s %-16s %s\n",
		"#", "type", "uniq", "lba_start", "lba_end", "attr", "name")

	for i := 0; i < int(lba_table_sectors); i++ {
		offset := (hdr.Hdr_lba_table + uint64(i)) * uint64(len(buf))
		ret, err := fp.ReadAt(buf, int64(offset))
		if err != nil {
			return err
		} else if ret != len(buf) {
			return errors.New("failed to read")
		}

		sector_entries := uint32(len(buf)) / hdr.Hdr_entsz
		entry_offset := 0

		for j := 0; j < int(sector_entries); j++ {
			p := gptEnt{}
			n := int(unsafe.Sizeof(p))
			r := bytes.NewReader(buf[entry_offset : entry_offset+n])
			binary.Read(r, binary.LittleEndian, &p)
			entry_offset += n

			empty := gptEnt{}
			if !dumpOptVerbose && p == empty {
				total++
				continue
			}

			name := make([]byte, 36)
			nlen := 0
			for k := 0; k < len(name); k++ {
				name[k] = byte(p.Ent_name[k] & 0xFF) // XXX ascii
				if name[k] == 0 {
					nlen = k
					break
				}
			}

			fmt.Printf("%-3d %-36s %-36s %016x %016x %016x %s\n",
				i*int(sector_entries)+j,
				tryKnownUuidToStr(&p.Ent_type),
				tryKnownUuidToStr(&p.Ent_uuid),
				p.Ent_lba_start,
				p.Ent_lba_end,
				p.Ent_attr,
				string(name[:nlen]))
			total++
		}
	}
	assert(total == int(hdr.Hdr_entries))

	return nil
}

func dumpGpt(fp *os.File) error {
	var hdr1, hdr2 *gptHdr
	var err error

	// primary header
	fmt.Println("primary header")
	hdr1, err = dumpHeader(fp, 1)
	if err != nil {
		return err
	}

	// secondary header
	if !dumpOptNoalt {
		fmt.Println("")
		fmt.Println("secondary header")
		hdr2, err = dumpHeader(fp, hdr1.Hdr_lba_alt)
		if err != nil {
			return err
		}
	}

	// primary entries
	fmt.Println("")
	fmt.Println("primary entries")
	err = dumpEntries(fp, hdr1)
	if err != nil {
		return err
	}

	// secondary entries
	if !dumpOptNoalt {
		fmt.Println("")
		fmt.Println("secondary entries")
		err := dumpEntries(fp, hdr2)
		if err != nil {
			return err
		}
	}

	return nil
}
