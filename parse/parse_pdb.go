package parse

import (
	"encoding/binary"
	"errors"
	"os"
)

const pdbAuthStreamSignature int32 = 0x01312E94

type pdbHeader struct {
	Description          [0x20]byte // 0x00-0x20
	PageSize             int32      // 0x20-0x24
	Unknown1             int32      // 0x24-0x28 0x00000002
	NumberOfPages        int32      // 0x28-0x2C
	StreamDirectoryBytes int32      // 0x2C-0x30
	Unknown2             int32      // 0x30-0x34 0x00000000
	StreamDirectoryPage  int32      // 0x34-0x38
}

type pdbAuthStream struct {
	Signature     int32      // 0x00-0x04 pdbAuthStreamSignature
	Unknown       int32      // 0x04-0x08
	TimeDateStamp int32      // 0x08-0x0C
	Guid          [0x10]byte // 0x0C-0x1C
}

func ParsePdb(file *os.File) (*DebugId, error) {
	var pdb pdbHeader
	// Read PDB header
	if err := binary.Read(file, binary.LittleEndian, &pdb); err != nil {
		return nil, err
	}
	pageSize := int64(pdb.PageSize)
	// Seek to PDB stream directory list
	if _, err := file.Seek(int64(pdb.StreamDirectoryPage)*pageSize, 0); err != nil {
		return nil, err
	}
	// Read first stream directory offset
	var stream_dir_page int32
	if err := binary.Read(file, binary.LittleEndian, &stream_dir_page); err != nil {
		return nil, err
	}
	// Seek to first stream directory
	if _, err := file.Seek(int64(stream_dir_page)*pageSize, 0); err != nil {
		return nil, err
	}
	// Read stream count (first DWORD in stream directory)
	var stream_count int32
	if err := binary.Read(file, binary.LittleEndian, &stream_count); err != nil {
		return nil, err
	}
	if stream_count < 2 {
		return nil, errors.New("Unexpected PDB stream count (need at least 2 streams)")
	}
	// Skip all stream sizes
	var size int32
	if _, err := file.Seek(int64(stream_count)*int64(binary.Size(&size)), 1); err != nil {
		return nil, err
	}
	// Read stream offsets
	stream_offset := make([]int32, stream_count)
	if err := binary.Read(file, binary.LittleEndian, &stream_offset); err != nil {
		return nil, err
	}
	for _, offset := range stream_offset {
		// Seek to second stream with PDB debug information
		if _, err := file.Seek(int64(offset)*int64(pdb.PageSize), 0); err != nil {
			return nil, err
		}
		// Read PDB debug information
		var auth pdbAuthStream
		if err := binary.Read(file, binary.LittleEndian, &auth); err != nil {
			return nil, err
		}
		if auth.Signature == pdbAuthStreamSignature {
			return &DebugId{
				Guid{
					auth.Guid,
				},
				int(auth.TimeDateStamp),
			}, nil
		}
	}
	return nil, errors.New("Can't find PDB auth stream")
}
