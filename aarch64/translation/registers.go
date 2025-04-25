package aarch64translation

import (
	"fmt"
	"math"

	"alon.kr/x/aarch64codegen/registers"
	"alon.kr/x/faststringmap"
	"alon.kr/x/list"
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/gen"
	"github.com/agnivade/levenshtein"
)

var spRegisterName = "%sp"
var xzrRegisterName = "%xzr"
var x0toX30RegisterNames = []string{"%x0", "%x1", "%x2", "%x3", "%x4", "%x5", "%x6", "%x7", "%x8", "%x9", "%x10", "%x11", "%x12", "%x13", "%x14", "%x15", "%x16", "%x17", "%x18", "%x19", "%x20", "%x21", "%x22", "%x23", "%x24", "%x25", "%x26", "%x27", "%x28", "%x29", "%x30"}

var ValidGPRegisterNames = append(append([]string{}, x0toX30RegisterNames...), xzrRegisterName)
var ValidGPorSPRegisterNames = append(append([]string{}, x0toX30RegisterNames...), spRegisterName)

var AllRegisterNames = append(append([]string{}, x0toX30RegisterNames...), spRegisterName, xzrRegisterName)

func newStringMapFromKeys[T ~uint8](keys []string) faststringmap.Map[T] {
	entries := make([]faststringmap.MapEntry[T], len(keys))
	for i, key := range keys {
		entries[i] = faststringmap.MapEntry[T]{Key: key, Value: T(i)}
	}

	return faststringmap.NewMap(entries)
}

// A mapping of USM register names to Aarch64 general purpose registers.
// This is used to convert USM register names to Aarch64 register.
//
// We use such mapping to avoid ambiguity in register names: this way, we
// explicitly DON'T allow register names like "X01" (prefixed with zeros),
// "x1" (lowercase 'x'), or "X31" (use "XZR" instead).

var registerNameToAarch64GPRegister = newStringMapFromKeys[registers.GPRegister](ValidGPRegisterNames)
var registerNameToAarch64GPorSPRegister = newStringMapFromKeys[registers.GPorSPRegister](ValidGPorSPRegisterNames)

func RegisterNameToAarch64GPRegister(
	name string,
) (registers.GPRegister, bool) {
	return registerNameToAarch64GPRegister.LookupString(name)
}

func RegisterNameToAarch64GPorSPRegister(
	name string,
) (registers.GPorSPRegister, bool) {
	return registerNameToAarch64GPorSPRegister.LookupString(name)
}

func closestLevenshteinDistance(name string, options []string) (string, int) {
	minDistance := math.MaxInt
	closestName := ""

	for _, option := range ValidGPRegisterNames {
		distance := levenshtein.ComputeDistance(name, option)
		if distance < minDistance {
			minDistance = distance
			closestName = option
		}
	}

	return closestName, minDistance
}

func RegisterToAarch64GPRegister(
	register *gen.RegisterInfo,
) (registers.GPRegister, core.ResultList) {
	name := register.Name
	reg, ok := RegisterNameToAarch64GPRegister(name)

	if !ok {
		// TODO: add a more sophisticated way to find the closest name
		// for example, if user wrote "X31", suggest "XZR" as an alternative.
		closestName, _ := closestLevenshteinDistance(name, ValidGPRegisterNames)

		return 0, list.FromSingle(core.Result{
			{
				Type:     core.ErrorResult,
				Message:  "Expected Aarch64 General Purpose register",
				Location: &register.Declaration,
			},
			{
				Type:    core.HintResult,
				Message: fmt.Sprintf("Did you mean \"%s\"?", closestName),
			},
		})
	}

	return reg, core.ResultList{}
}

func RegisterToAarch64GPOrSPRegister(
	register *gen.RegisterInfo,
) (registers.GPorSPRegister, core.ResultList) {
	name := register.Name
	reg, ok := RegisterNameToAarch64GPorSPRegister(name)

	if !ok {
		closestName, _ := closestLevenshteinDistance(name, ValidGPorSPRegisterNames)

		return 0, list.FromSingle(core.Result{
			{
				Type:     core.ErrorResult,
				Message:  "Expected Aarch64 General Purpose or SP register",
				Location: &register.Declaration,
			},
			{
				Type:    core.HintResult,
				Message: fmt.Sprintf("Did you mean \"%s\"?", closestName),
			},
		})
	}

	return reg, core.ResultList{}
}
