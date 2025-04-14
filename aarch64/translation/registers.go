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

var x0toX30RegisterNames = []string{"%X0", "%X1", "%X2", "%X3", "%X4", "%X5", "%X6", "%X7", "%X8", "%X9", "%X10", "%X11", "%X12", "%X13", "%X14", "%X15", "%X16", "%X17", "%X18", "%X19", "%X20", "%X21", "%X22", "%X23", "%X24", "%X25", "%X26", "%X27", "%X28", "%X29", "%X30"}
var validGPRegisterNames = append(append([]string{}, x0toX30RegisterNames...), "%XZR")
var validGPorSPRegisterNames = append(append([]string{}, x0toX30RegisterNames...), "%SP")

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

var registerNameToAarch64GPRegister = newStringMapFromKeys[registers.GPRegister](validGPRegisterNames)
var registerNameToAarch64GPorSPRegister = newStringMapFromKeys[registers.GPorSPRegister](validGPorSPRegisterNames)

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

	for _, option := range validGPRegisterNames {
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
		closestName, _ := closestLevenshteinDistance(name, validGPRegisterNames)

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
		closestName, _ := closestLevenshteinDistance(name, validGPorSPRegisterNames)

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
