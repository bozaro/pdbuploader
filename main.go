// https://rsdn.ru/article/baseserv/pe_coff.xml
// http://www.godevtool.com/Other/pdb.htm
package main

import (
	"encoding/binary"
	"fmt"
	"os"
)

type DebugInfo struct {
	CodeId  string
	DebugId string
}

type MZHeader struct {
	Signature int16      // 0x00-0x02 0x5A4D
	Unused    [0x3A]byte // 0x02-0x3C
	PEOffset  int32      // 0x3C-0x40
}

type RvaAndSize struct {
	VirtualAddress int32
	VirtualSize    int32
}

type PEHeader struct {
	Signature                   int32      // 0x00-0x04 0x50450000
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

type PESection struct {
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

type PEDebugDirectory struct {
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
	Signature     int32      // 0x00-0x04 0x53445352
	GUID          [0x10]byte // 0x04-0x14
	TimeDateStamp int32      // 0x14-0x18
}

type PDBHeader struct {
	Description          [0x20]byte // 0x00-0x20
	PageSize             int32      // 0x20-0x24
	Unknown1             int32      // 0x24-0x28 0x00000002
	NumberOfPages        int32      // 0x28-0x2C
	StreamDirectoryBytes int32      // 0x2C-0x30
	Unknown2             int32      // 0x30-0x34 0x00000000
	StreamDirectoryPage  int32      // 0x34-0x38
}

type PDBAuthStream struct {
	Signature     int32      // 0x00-0x04 0x01312E94
	Unknown       int32      // 0x04-0x08
	TimeDateStamp int32      // 0x08-0x0C
	GUID          [0x10]byte // 0x0C-0x1C
}

func guid_to_string(guid [0x10]byte) string {
	return fmt.Sprintf("%02X%02X%02X%02X%02X%02X%02X%02X%16X",
		guid[3], guid[2], guid[1], guid[0],
		guid[5], guid[4],
		guid[7], guid[6],
		guid[8:])
}

func read_pdb_debug_id(file *os.File) string {
	var pdb PDBHeader
	binary.Read(file, binary.LittleEndian, &pdb)
	fmt.Printf("PDB signature: %s\n", pdb.Description)
	stream_dirs_offset := int64(pdb.StreamDirectoryPage) * int64(pdb.PageSize)
	fmt.Printf("Stream directory list offset: %08X (%d)\n", stream_dirs_offset, pdb.NumberOfPages)

	file.Seek(stream_dirs_offset, 0)
	var stream_dir_page int32
	binary.Read(file, binary.LittleEndian, &stream_dir_page)

	stream_dir := int64(stream_dir_page) * int64(pdb.PageSize)
	fmt.Printf("Stream directory offset: %08X\n", stream_dir)
	file.Seek(stream_dir, 0)

	var stream_count int32
	binary.Read(file, binary.LittleEndian, &stream_count)
	fmt.Printf("Streams count: %08X\n", stream_count)

	var size int32
	file.Seek(int64(stream_count)*int64(binary.Size(&size)), 1)
	var offset int32
	file.Seek(int64(1)*int64(binary.Size(&offset)), 1)

	var stream_offset int32
	binary.Read(file, binary.LittleEndian, &stream_offset)
	file.Seek(int64(stream_offset)*int64(pdb.PageSize), 0)

	var auth PDBAuthStream
	binary.Read(file, binary.LittleEndian, &auth)

	return fmt.Sprintf("%s%d", guid_to_string(auth.GUID), auth.TimeDateStamp)
}

func read_exe_debug_info(file *os.File) DebugInfo {
	var mz MZHeader
	var pe PEHeader
	binary.Read(file, binary.LittleEndian, &mz)

	fmt.Printf("MZ signature: %04X\n", mz.Signature)
	fmt.Printf("PE offset: %08X\n", mz.PEOffset)

	file.Seek(int64(mz.PEOffset), 0)
	binary.Read(file, binary.LittleEndian, &pe)

	fmt.Printf("PE signature: %08X\n", pe.Signature)
	fmt.Printf("PE timestamp: %08X\n", pe.TimeDateStamp)
	fmt.Printf("PE image size: %08X\n", pe.SizeOfImage)

	fmt.Printf("Sections count: %d\n", pe.NumberOfSections)
	fmt.Printf("Sections alignment: %d\n", pe.SectionAlignment)
	fmt.Printf("Size of headers: %d\n", pe.SizeOfHeaders)

	var debug_rva RvaAndSize
	if pe.NumberOfRvaAndSizes < 7 {
		// todo: opss...
		fmt.Println("OPS....")
	}
	for i := 0; i < 7; i++ {
		binary.Read(file, binary.LittleEndian, &debug_rva)
	}

	file.Seek(int64(mz.PEOffset)+int64(pe.SizeOfOptionalHeader)+0x18, 0)

	fmt.Printf("Section offset: %08X\n", int64(mz.PEOffset)+int64(pe.SizeOfOptionalHeader)+0x18)
	debug_dir_offest := int64(0)
	for i := int16(0); i < pe.NumberOfSections; i++ {
		var section PESection
		binary.Read(file, binary.LittleEndian, &section)
		if (section.VirtualAddress <= debug_rva.VirtualAddress) && (section.VirtualAddress+section.VirtualSize > debug_rva.VirtualAddress) {
			debug_dir_offest = int64(section.PointerToRawData + debug_rva.VirtualAddress - section.VirtualAddress)
			break
		}
	}

	rsds_offset := int64(0)
	if debug_dir_offest > 0 {
		file.Seek(debug_dir_offest, 0)
		var debug_dir PEDebugDirectory
		fmt.Printf("IMAGE_DEBUG_DIRECTORY offset: %X (%d)\n", debug_dir_offest, binary.Size(&debug_dir))
		for i := 0; i < int(debug_rva.VirtualSize)/binary.Size(&debug_dir); i++ {
			binary.Read(file, binary.LittleEndian, &debug_dir)
			fmt.Printf("   %d: %d\n", i, debug_dir.Type)
			if debug_dir.Type == 2 {
				fmt.Printf("RSDS offset: %X\n", debug_dir.PointerToRawData)
				rsds_offset = int64(debug_dir.PointerToRawData)
				break
			}
		}
	}

	var rsds RSDSHeader
	if rsds_offset > 0 {
		file.Seek(rsds_offset, 0)
		binary.Read(file, binary.LittleEndian, &rsds)
	}

	fmt.Printf("RSDS signature: %08X\n", rsds.Signature)
	fmt.Printf("RSDS timestamp: %08X\n", rsds.TimeDateStamp)

	return DebugInfo{
		fmt.Sprintf("%X%x", pe.TimeDateStamp, pe.SizeOfImage),
		fmt.Sprintf("%s%d", guid_to_string(rsds.GUID), rsds.TimeDateStamp),
	}
}

func main() {
	file, _ := os.Open("sample/hello.exe")
	info := read_exe_debug_info(file)
	file, _ = os.Open("sample/hello.pdb")
	debug_id := read_pdb_debug_id(file)
	fmt.Println("EXE")
	fmt.Printf("  Code ID: %s\n", info.CodeId)
	fmt.Printf("  Debug ID: %s\n", info.DebugId)
	fmt.Println("PDB")
	fmt.Printf("  Debug ID: %s\n", debug_id)
}
