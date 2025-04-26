package aarch64translation

import (
	"bytes"

	"alon.kr/x/macho/builder"
	"alon.kr/x/macho/header"
	"alon.kr/x/macho/load/build_version"
	"alon.kr/x/macho/load/nlist64"
	nlist64_builders "alon.kr/x/macho/load/nlist64/builders"
	"alon.kr/x/macho/load/section64"
	"alon.kr/x/macho/load/segment64"
	"alon.kr/x/macho/load/symtab"
	"alon.kr/x/macho/load/symtab/symbol"
	aarch64codegen "alon.kr/x/usm/aarch64/codegen"
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/gen"
)

func FileToMachoObject(file *gen.FileInfo) ([]byte, core.ResultList) {
	fileCtx := aarch64codegen.NewFileCodegenContext(file)

	data := bytes.Buffer{}
	results := fileCtx.Codegen(&data)
	if !results.IsEmpty() {
		return nil, results
	}

	symbols := []symbol.SymbolBuilder{}
	for _, function := range file.Functions {
		symbol := nlist64_builders.SectionNlist64Builder{
			Name:        "_" + function.Name[1:],
			Type:        nlist64.ExternalSymbol | nlist64.SectionSymbolType,
			Section:     1,
			Offset:      fileCtx.FunctionOffsets[function],
			Description: nlist64.ReferenceFlagUndefinedNonLazy,
		}

		symbols = append(symbols, symbol)
	}

	sectionBuilder := section64.Section64Builder{
		SectionName: [16]byte{'_', '_', 't', 'e', 'x', 't', 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		SegmentName: [16]byte{'_', '_', 'T', 'E', 'X', 'T', 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		Data:        data.Bytes(),
		Flags:       section64.AttrPureInstructions | section64.AttrSomeInstructions,
	}

	segmentBuilder := segment64.Segment64Builder{
		SegmentName:        [16]byte{'_', '_', 'T', 'E', 'X', 'T', 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		Sections:           []section64.Section64Builder{sectionBuilder},
		VirtualMemorySize:  16,
		MaxProtections:     segment64.AllowAllProtection,
		InitialProtections: segment64.AllowAllProtection,
	}

	symtabBuilder := symtab.SymtabBuilder{
		Symbols: symbols,
	}

	buildVersionBuilder := build_version.BuildVersionBuilder{
		Platform: build_version.PlatformMacOS,
	}

	headerBuilder := header.MachoHeaderBuilder{
		Magic:      header.Magic64Bit,
		CpuType:    header.Arm64CpuType,
		CpuSubType: header.AllArmProcessors,
		FileType:   header.Object,
	}

	machoBuilder := builder.MachoBuilder{
		Header: headerBuilder,
		Commands: []builder.CommandBuilder{
			segmentBuilder,
			symtabBuilder,
			buildVersionBuilder,
		},
	}

	buffer := new(bytes.Buffer)
	machoBuilder.WriteTo(buffer)
	return buffer.Bytes(), core.ResultList{}
}
