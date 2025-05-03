package aarch64translation

import (
	"bytes"
	"fmt"

	"alon.kr/x/list"
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
	"alon.kr/x/usm/transform"
)

func ToMachoObject(
	data *transform.TargetData,
) (*transform.TargetData, core.ResultList) {
	file := data.Code
	fileCtx := aarch64codegen.NewFileCodegenContext(file)

	codeBuffer := bytes.Buffer{}
	results := fileCtx.Codegen(&codeBuffer)
	if !results.IsEmpty() {
		return nil, results
	}

	symbols := []symbol.SymbolBuilder{}
	for _, function := range file.Functions {
		symbol := nlist64_builders.SectionNlist64Builder{
			Name:        "_" + function.Name[1:],
			Type:        nlist64.ExternalSymbol,
			Description: nlist64.ReferenceFlagUndefinedNonLazy,
		}

		if function.IsDefined() {
			symbol.Type |= nlist64.SectionSymbolType
			symbol.Section = 1
			symbol.Offset = fileCtx.FunctionOffsets[function]
		}

		symbols = append(symbols, symbol)
	}

	sectionBuilder := section64.Section64Builder{
		SectionName: [16]byte{'_', '_', 't', 'e', 'x', 't', 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		SegmentName: [16]byte{'_', '_', 'T', 'E', 'X', 'T', 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		Data:        codeBuffer.Bytes(),
		Flags:       section64.AttrPureInstructions | section64.AttrSomeInstructions,
		Relocations: fileCtx.Relocations,
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

	machoBuffer := new(bytes.Buffer)
	_, err := machoBuilder.WriteTo(machoBuffer)
	if err != nil {
		return nil, list.FromSingle(core.Result{
			{
				Type:    core.ErrorResult,
				Message: fmt.Sprintf("Failed to generate Mach-O object file: %v", err),
			},
		})
	}

	data.Artifact = machoBuffer
	return data, core.ResultList{}
}
