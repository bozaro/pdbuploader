package main

import (
	"encoding/binary"
	"fmt"
	"os"
	"errors"
)

const mzHeaderSignature int16 = 0x5A4D
const peHeaderSignature int32 = 0x00004550
const imageDirectoryEntryDebug int32 = 6

type mzHeader struct {
	Signature int16      // 0x00-0x02 mzHeaderSignature
	Unused    [0x3A]byte // 0x02-0x3C
	PEOffset  int32      // 0x3C-0x40
}

type rvaAndSize struct {
	VirtualAddress int32
	VirtualSize    int32
}

type peHeader struct {
	Signature                   int32      // 0x00-0x04 peHeaderSignature
	Machine                     int16      // 0x04-0x06
	NumberOfSections            int16      // 0x06-0x08
	TimeDateStamp               int32      // 0x08-0x0C
	PointerToSymbolTable        int32      // 0x0C-0x10
	NumberOfSymbolTable         int32      // 0x10-0x14
	SizeOfOptionalHeader        int16      // 0x14-0x16
	Characteristics             int16      // 0x16-0x18
	StandadCOFFFields           [0x1C]byte // 0x18-0x34
	ImageBase                   int32      // 0x34-0x38
	SectionAlignment            int32      // 0x38-0x3C
	FileAlignment               int32      // 0x3C-0x40
	MajorOperatingSystemVersion int16      // 0x40-0x42
	MinorOperatingSystemVersion int16      // 0x40-0x42
	MajorImageVersion           int16      // 0x42-0x44
	MinorImageVersion           int16      // 0x44-0x46
	MajorSubsystemVersion       int16      // 0x46-0x48
	MinorSubsystemVersion       int16      // 0x48-0x4A
	Win32VersionValue           int32      // 0x4A-0x50
	SizeOfImage                 int32      // 0x50-0x54
	SizeOfHeaders               int32      // 0x54-0x58
	CheckSum                    int32      // 0x58-0x5C
	Subsystem                   int16      // 0x5C-0x5E
	DllCharacteristics          int16      // 0x5E-0x60
	SizeOfStackReserve          int32      // 0x60-0x64
	SizeOfStackCommit           int32      // 0x64-0x68
	SizeOfHeapReserve           int32      // 0x68-0x6C
	SizeOfHeapCommit            int32      // 0x6C-0x70
	LoaderFlags                 int32      // 0x70-0x74
	NumberOfRvaAndSizes         int32      // 0x74-0x78
}

type peSection struct {
	Name                 [0x08]byte
	VirtualSize          int32
	VirtualAddress       int32
	SizeOfRawData        int32
	PointerToRawData     int32
	PointerToRelocations int32
	PointerToLinenumbers int32
	NumberOfRelocations  int16
	NumberOfLinenumbers  int16
	Characteristics      int32
}

type peDebugDirectory struct {
	Characteristics  int32
	TimeDateStamp    int32
	MajorVersion     int16
	MinorVersion     int16
	Type             int32
	SizeOfData       int32
	AddressOfRawData int32
	PointerToRawData int32
}

type RSDSHeader struct {
	Signature     int32       // 0x00-0x04 0x53445352
	GUID          [0x10]byte  // 0x04-0x14
	TimeDateStamp int32       // 0x14-0x18
	PDBFileName   [0x104]byte // 0x18-0x11C
}

func ParseExe(file *os.File) (*DebugInfo, error) {
	var mz mzHeader
	// Read MZ DOS header
	if err := binary.Read(file, binary.LittleEndian, &mz); err != nil {
		return nil, err
	}
	if mz.Signature != mzHeaderSignature {
		return nil, errors.New("Invalid MZ header signature")
	}
	// Seek to PE header
	if _, err := file.Seek(int64(mz.PEOffset), 0); err != nil {
		return nil, err
	}
	// Read PE header
	var pe peHeader
	if err := binary.Read(file, binary.LittleEndian, &pe); err != nil {
		return nil, err
	}
	if pe.Signature != peHeaderSignature {
		return nil, errors.New("Invalid PE header signature")
	}

	fmt.Printf("PE timestamp: %08X\n", pe.TimeDateStamp)
	fmt.Printf("PE image size: %08X\n", pe.SizeOfImage)

	fmt.Printf("Sections count: %d\n", pe.NumberOfSections)
	fmt.Printf("Sections alignment: %d\n", pe.SectionAlignment)
	fmt.Printf("Size of headers: %d\n", pe.SizeOfHeaders)

	if pe.NumberOfRvaAndSizes < imageDirectoryEntryDebug {
		return nil, errors.New("Debug information not found in RVA table")
	}
	// Skip RVA entries before IMAGE_DIRECTORY_ENTRY_DEBUG
	for i := int32(0); i < imageDirectoryEntryDebug; i++ {
		var skip_rva rvaAndSize
		if err := binary.Read(file, binary.LittleEndian, &skip_rva); err != nil {
			return nil, err
		}
	}
	// Read RVA for IMAGE_DIRECTORY_ENTRY_DEBUG
	var debug_rva rvaAndSize
	if err := binary.Read(file, binary.LittleEndian, &debug_rva); err != nil {
		return nil, err
	}

	// Seek to PE sections after PE header
	if _, err := file.Seek(int64(mz.PEOffset) + int64(pe.SizeOfOptionalHeader) + 0x18, 0); err != nil {
		return nil, err
	}
	// Find file offset for IMAGE_DEBUG_DIRECTORY
	debug_dir_offest := int64(0)
	for i := int16(0); i < pe.NumberOfSections; i++ {
		var section peSection
		binary.Read(file, binary.LittleEndian, &section)
		if (section.VirtualAddress <= debug_rva.VirtualAddress) && (section.VirtualAddress + section.VirtualSize > debug_rva.VirtualAddress) {
			debug_dir_offest = int64(section.PointerToRawData + debug_rva.VirtualAddress - section.VirtualAddress)
			break
		}
	}
	if debug_dir_offest <= 0 {
		return nil, errors.New("Can't find offset for IMAGE_DEBUG_DIRECTORY")
	}

	// Seek to IMAGE_DEBUG_DIRECTORY
	if _, err := file.Seek(debug_dir_offest, 0); err != nil {
		return nil, err
	}

	rsds_offset := int64(0)
	var debug_dir peDebugDirectory
	fmt.Printf("IMAGE_DEBUG_DIRECTORY offset: %X (%d)\n", debug_dir_offest, binary.Size(&debug_dir))
	for i := 0; i < int(debug_rva.VirtualSize) / binary.Size(&debug_dir); i++ {
		binary.Read(file, binary.LittleEndian, &debug_dir)
		fmt.Printf("   %d: %d\n", i, debug_dir.Type)
		if debug_dir.Type == 2 {
			fmt.Printf("RSDS offset: %X\n", debug_dir.PointerToRawData)
			rsds_offset = int64(debug_dir.PointerToRawData)
			break
		}
	}

	var rsds RSDSHeader
	if rsds_offset > 0 {
		file.Seek(rsds_offset, 0)
		binary.Read(file, binary.LittleEndian, &rsds)
	}

	fmt.Printf("RSDS signature: %08X\n", rsds.Signature)
	fmt.Printf("RSDS timestamp: %08X\n", rsds.TimeDateStamp)

	return &DebugInfo{
		fmt.Sprintf("%X%x", pe.TimeDateStamp, pe.SizeOfImage),
		fmt.Sprintf("%s%d", guid_to_string(rsds.GUID), rsds.TimeDateStamp),
		CToGoString(rsds.PDBFileName[:]),
	}, nil
}
