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
	HdrSig      [8]uint8
	HdrRevision uint32
	HdrSize     uint32
	HdrCrcSelf  uint32
	Reserved    uint32
	HdrLbaSelf  uint64
	HdrLbaAlt   uint64
	HdrLbaStart uint64
	HdrLbaEnd   uint64
	HdrUuid     uuid
	HdrLbaTable uint64
	HdrEntries  uint32
	HdrEntsz    uint32
	HdrCrcTable uint32
	Padding     uint32
}

type gptEnt struct {
	EntType     uuid
	EntUuid     uuid
	EntLbaStart uint64
	EntLbaEnd   uint64
	EntAttr     uint64
	EntName     [36]uint16
}

func tryKnownUuidToStr(uuid *uuid) string {
	if optSymbol {
		if s := knownUuidToStr(uuid); len(s) != 0 {
			return s
		}
	}

	return uuidToStr(uuid)
}

func allocBuffer() []byte {
	const UNIT_SIZE = 512

	buf := make([]byte, UNIT_SIZE)

	assert(len(buf) == UNIT_SIZE)
	assert(len(buf)%512 == 0)

	return buf
}

func dumpHeader(fp *os.File, hdrLba uint64) (*gptHdr, error) {
	buf := allocBuffer()

	hdrOffset := hdrLba * uint64(len(buf))

	if ret, err := fp.ReadAt(buf, int64(hdrOffset)); err != nil {
		return nil, err
	} else if ret != len(buf) {
		return nil, errors.New("failed to read")
	}

	hdr := gptHdr{}
	n := int(unsafe.Sizeof(hdr))
	r := bytes.NewReader(buf[:n])
	if err := binary.Read(r, binary.LittleEndian, &hdr); err != nil {
		return nil, err
	}

	l1 := []byte("EFI PART")
	l2 := []byte(hdr.HdrSig[:])
	if !bytes.Equal(l1, l2) {
		return nil, errors.New("not GPT")
	}

	fmt.Printf("sig      = \"%c%c%c%c%c%c%c%c\"\n",
		hdr.HdrSig[0],
		hdr.HdrSig[1],
		hdr.HdrSig[2],
		hdr.HdrSig[3],
		hdr.HdrSig[4],
		hdr.HdrSig[5],
		hdr.HdrSig[6],
		hdr.HdrSig[7])

	p := make([]byte, 4)
	binary.LittleEndian.PutUint32(p, hdr.HdrRevision)
	fmt.Printf("revision = %02x %02x %02x %02x\n", p[0], p[1], p[2], p[3])

	fmt.Printf("size     = %d\n", hdr.HdrSize)
	fmt.Printf("crc_self = 0x%x\n", hdr.HdrCrcSelf)
	fmt.Printf("lba_self = 0x%016x\n", hdr.HdrLbaSelf)
	fmt.Printf("lba_alt  = 0x%016x\n", hdr.HdrLbaAlt)
	fmt.Printf("lba_start= 0x%016x\n", hdr.HdrLbaStart)
	fmt.Printf("lba_end  = 0x%016x\n", hdr.HdrLbaEnd)

	fmt.Printf("uuid     = %s\n", tryKnownUuidToStr(&hdr.HdrUuid))

	fmt.Printf("lba_table= 0x%016x\n", hdr.HdrLbaTable)
	fmt.Printf("entries  = %d\n", hdr.HdrEntries)
	fmt.Printf("entsz    = %d\n", hdr.HdrEntsz)
	fmt.Printf("crc_table= 0x%x\n", hdr.HdrCrcTable)

	// XXX
	if hdr.HdrEntries > 512 {
		return nil, errors.New("likely corrupted entries")
	}

	return &hdr, nil
}

func dumpEntries(fp *os.File, hdr *gptHdr) error {
	buf := allocBuffer()

	lbaTableSize := hdr.HdrEntsz * hdr.HdrEntries
	lbaTableSectors := lbaTableSize / uint32(len(buf))
	total := 0

	fmt.Printf("%-3s %-36s %-36s %-16s %-16s %-16s %s\n",
		"#", "type", "uniq", "lba_start", "lba_end", "attr", "name")

	for i := 0; i < int(lbaTableSectors); i++ {
		offset := (hdr.HdrLbaTable + uint64(i)) * uint64(len(buf))
		if ret, err := fp.ReadAt(buf, int64(offset)); err != nil {
			return err
		} else if ret != len(buf) {
			return errors.New("failed to read")
		}

		sectorEntries := uint32(len(buf)) / hdr.HdrEntsz
		entryOffset := 0

		for j := 0; j < int(sectorEntries); j++ {
			p := gptEnt{}
			n := int(unsafe.Sizeof(p))
			r := bytes.NewReader(buf[entryOffset : entryOffset+n])
			if err := binary.Read(r, binary.LittleEndian, &p); err != nil {
				return err
			}

			entryOffset += n

			empty := gptEnt{}
			if !optVerbose && p == empty {
				total++
				continue
			}

			name := make([]byte, 36)
			nlen := 0
			for k := 0; k < len(name); k++ {
				name[k] = byte(p.EntName[k] & 0xFF) // XXX ascii
				if name[k] == 0 {
					nlen = k
					break
				}
			}

			fmt.Printf("%-3d %-36s %-36s %016x %016x %016x %s\n",
				i*int(sectorEntries)+j,
				tryKnownUuidToStr(&p.EntType),
				tryKnownUuidToStr(&p.EntUuid),
				p.EntLbaStart,
				p.EntLbaEnd,
				p.EntAttr,
				string(name[:nlen]))
			total++
		}
	}
	assert(total == int(hdr.HdrEntries))

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
	if !optNoAlt {
		fmt.Println("")
		fmt.Println("secondary header")
		hdr2, err = dumpHeader(fp, hdr1.HdrLbaAlt)
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
	if !optNoAlt {
		fmt.Println("")
		fmt.Println("secondary entries")
		err := dumpEntries(fp, hdr2)
		if err != nil {
			return err
		}
	}

	return nil
}
