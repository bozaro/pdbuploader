package parse

// https://msdn.microsoft.com/en-us/windows/hardware/gg463119.aspx

import (
	"encoding/binary"
	"errors"
	"os"
)

const (
	mzHeaderSignature           int16 = 0x5A4D
	peHeaderSignature           int32 = 0x00004550
	rsdsHeaderSignature         int32 = 0x53445352
	IMAGE_DIRECTORY_ENTRY_DEBUG int32 = 6
	IMAGE_DEBUG_TYPE_CODEVIEW   int32 = 2

	COFF_PE32_32 int16 = 0x010B
	COFF_PE32_64 int16 = 0x020B
)

type mzHeader struct {
	Signature int16      // 0x00-0x02 mzHeaderSignature
	Unused    [0x3A]byte // 0x02-0x3C
	PEOffset  int32      // 0x3C-0x40
}

type rvaAndSize struct {
	VirtualAddress int32
	VirtualSize    int32
}

type coffHeader struct {
	Machine              uint16 // 0x04-0x06
	NumberOfSections     int16  // 0x06-0x08
	TimeDateStamp        int32  // 0x08-0x0C
	PointerToSymbolTable int32  // 0x0C-0x10
	NumberOfSymbolTable  int32  // 0x10-0x14
	SizeOfOptionalHeader int16  // 0x14-0x16
	Characteristics      int16  // 0x16-0x18

	// Optional Header Standard Fields
	Magic                   int16 // 0x18-0x1A
	MajorLinkerVersion      int8  // 0x1A-0x1B
	MinorLinkerVersion      int8  // 0x1B-0x1C
	SizeOfCode              int32 // 0x1C-0x20
	SizeOfInitializedData   int32 // 0x20-0x24
	SizeOfUninitializedData int32 // 0x24-0x28
	AddressOfEntryPoint     int32 // 0x28-0x2C
	BaseOfCode              int32 // 0x2C-0x30
}

type peHeader32 struct {
	BaseOfData                  int32 // 0x30-0x34
	ImageBase                   int32 // 0x34-0x38
	SectionAlignment            int32 // 0x38-0x3C
	FileAlignment               int32 // 0x3C-0x40
	MajorOperatingSystemVersion int16 // 0x40-0x42
	MinorOperatingSystemVersion int16 // 0x40-0x42
	MajorImageVersion           int16 // 0x42-0x44
	MinorImageVersion           int16 // 0x44-0x46
	MajorSubsystemVersion       int16 // 0x46-0x48
	MinorSubsystemVersion       int16 // 0x48-0x4A
	Win32VersionValue           int32 // 0x4A-0x50
	SizeOfImage                 int32 // 0x50-0x54
	SizeOfHeaders               int32 // 0x54-0x58
	CheckSum                    int32 // 0x58-0x5C
	Subsystem                   int16 // 0x5C-0x5E
	DllCharacteristics          int16 // 0x5E-0x60
	SizeOfStackReserve          int32 // 0x60-0x64
	SizeOfStackCommit           int32 // 0x64-0x68
	SizeOfHeapReserve           int32 // 0x68-0x6C
	SizeOfHeapCommit            int32 // 0x6C-0x70
	LoaderFlags                 int32 // 0x70-0x74
	NumberOfRvaAndSizes         int32 // 0x74-0x78
}

type peHeader64 struct {
	ImageBase                   int64 // 0x30-0x38
	SectionAlignment            int32 // 0x38-0x3C
	FileAlignment               int32 // 0x3C-0x40
	MajorOperatingSystemVersion int16 // 0x40-0x42
	MinorOperatingSystemVersion int16 // 0x40-0x42
	MajorImageVersion           int16 // 0x42-0x44
	MinorImageVersion           int16 // 0x44-0x46
	MajorSubsystemVersion       int16 // 0x46-0x48
	MinorSubsystemVersion       int16 // 0x48-0x4A
	Win32VersionValue           int32 // 0x4A-0x50
	SizeOfImage                 int32 // 0x50-0x54
	SizeOfHeaders               int32 // 0x54-0x58
	CheckSum                    int32 // 0x58-0x5C
	Subsystem                   int16 // 0x5C-0x5E
	DllCharacteristics          int16 // 0x5E-0x60
	SizeOfStackReserve          int64 // 0x60-0x68
	SizeOfStackCommit           int64 // 0x68-0x70
	SizeOfHeapReserve           int64 // 0x70-0x78
	SizeOfHeapCommit            int64 // 0x78-0x80
	LoaderFlags                 int32 // 0x80-0x84
	NumberOfRvaAndSizes         int32 // 0x84-0x88
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

type rsdsHeader struct {
	Signature     int32       // 0x00-0x04 rsdsHeaderSignature
	Guid          [0x10]byte  // 0x04-0x14
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
	// Read PE signature
	var peSignature int32
	if err := binary.Read(file, binary.LittleEndian, &peSignature); err != nil {
		return nil, err
	}
	if peSignature != peHeaderSignature {
		return nil, errors.New("Invalid PE header signature")
	}
	// Read COFF header
	var coff coffHeader
	if err := binary.Read(file, binary.LittleEndian, &coff); err != nil {
		return nil, err
	}
	// Read platform specific COFF header
	var sizeOfImage int32
	switch coff.Magic {
	case COFF_PE32_32:
		var pe peHeader32
		if err := binary.Read(file, binary.LittleEndian, &pe); err != nil {
			return nil, err
		}
		if pe.NumberOfRvaAndSizes < IMAGE_DIRECTORY_ENTRY_DEBUG {
			return nil, errors.New("Debug information not found in RVA table")
		}
		sizeOfImage = pe.SizeOfImage
	case COFF_PE32_64:
		var pe peHeader64
		if err := binary.Read(file, binary.LittleEndian, &pe); err != nil {
			return nil, err
		}
		if pe.NumberOfRvaAndSizes < IMAGE_DIRECTORY_ENTRY_DEBUG {
			return nil, errors.New("Debug information not found in RVA table")
		}
		sizeOfImage = pe.SizeOfImage
	default:
		return nil, errors.New("Unsupported COFF-header format")
	}
	// Skip RVA entries before IMAGE_DIRECTORY_ENTRY_DEBUG
	for i := int32(0); i < IMAGE_DIRECTORY_ENTRY_DEBUG; i++ {
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
	if _, err := file.Seek(int64(mz.PEOffset)+int64(coff.SizeOfOptionalHeader)+0x18, 0); err != nil {
		return nil, err
	}
	// Find file offset for IMAGE_DEBUG_DIRECTORY
	debug_dir_offest := int64(0)
	for i := int16(0); i < coff.NumberOfSections; i++ {
		var section peSection
		binary.Read(file, binary.LittleEndian, &section)
		if (section.VirtualAddress <= debug_rva.VirtualAddress) && (section.VirtualAddress+section.VirtualSize > debug_rva.VirtualAddress) {
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

	// Search IMAGE_DEBUG_TYPE_CODEVIEW offset
	rsds_offset := int64(0)
	var debug_dir peDebugDirectory
	for i := 0; i < int(debug_rva.VirtualSize)/binary.Size(&debug_dir); i++ {
		if err := binary.Read(file, binary.LittleEndian, &debug_dir); err != nil {
			return nil, err
		}
		if debug_dir.Type == IMAGE_DEBUG_TYPE_CODEVIEW {
			rsds_offset = int64(debug_dir.PointerToRawData)
			break
		}
	}
	if rsds_offset <= 0 {
		return nil, errors.New("Can't find offset for CV_INFO_PDB20 debug information")
	}

	// CV_INFO_PDB20 debug information
	var rsds rsdsHeader
	if _, err := file.Seek(rsds_offset, 0); err != nil {
		return nil, err
	}
	if err := binary.Read(file, binary.LittleEndian, &rsds); err != nil {
		return nil, err
	}
	if rsds.Signature != rsdsHeaderSignature {
		return nil, errors.New("Invalid PE header signature")
	}

	return &DebugInfo{
		CodeId{
			coff.TimeDateStamp,
			sizeOfImage,
		},
		DebugId{
			Guid{
				rsds.Guid,
			},
			int(rsds.TimeDateStamp),
		},
		CToGoString(rsds.PDBFileName[:]),
	}, nil
}
