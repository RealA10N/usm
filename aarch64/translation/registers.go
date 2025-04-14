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

var validRegisterNames = []faststringmap.MapEntry[registers.GPRegister]{
	{Key: "%X0", Value: registers.X0},
	{Key: "%X1", Value: registers.X1},
	{Key: "%X2", Value: registers.X2},
	{Key: "%X3", Value: registers.X3},
	{Key: "%X4", Value: registers.X4},
	{Key: "%X5", Value: registers.X5},
	{Key: "%X6", Value: registers.X6},
	{Key: "%X7", Value: registers.X7},
	{Key: "%X8", Value: registers.X8},
	{Key: "%X9", Value: registers.X9},
	{Key: "%X10", Value: registers.X10},
	{Key: "%X11", Value: registers.X11},
	{Key: "%X12", Value: registers.X12},
	{Key: "%X13", Value: registers.X13},
	{Key: "%X14", Value: registers.X14},
	{Key: "%X15", Value: registers.X15},
	{Key: "%X16", Value: registers.X16},
	{Key: "%X17", Value: registers.X17},
	{Key: "%X18", Value: registers.X18},
	{Key: "%X19", Value: registers.X19},
	{Key: "%X20", Value: registers.X20},
	{Key: "%X21", Value: registers.X21},
	{Key: "%X22", Value: registers.X22},
	{Key: "%X23", Value: registers.X23},
	{Key: "%X24", Value: registers.X24},
	{Key: "%X25", Value: registers.X25},
	{Key: "%X26", Value: registers.X26},
	{Key: "%X27", Value: registers.X27},
	{Key: "%X28", Value: registers.X28},
	{Key: "%X29", Value: registers.X29},
	{Key: "%X30", Value: registers.X30},
	{Key: "%XZR", Value: registers.XZR},
}

// A mapping of USM register names to Aarch64 general purpose registers.
// This is used to convert USM register names to Aarch64 register.
//
// We use such mapping to avoid ambiguity in register names: this way, we
// explicitly DON'T allow register names like "X01" (prefixed with zeros),
// "x1" (lowercase 'x'), or "X31" (use "XZR" instead).
var registerNameToAarch64GPRegister = faststringmap.NewMap(validRegisterNames)

func RegisterNameToAarch64GPRegister(
	name string,
) (registers.GPRegister, bool) {
	return registerNameToAarch64GPRegister.LookupString(name)
}

// closestAarch64GPRegisterName finds the closest Aarch64 general-purpose
// register name to the given name using Levenshtein distance.
func closestAarch64GPRegisterName(name string) (string, int) {
	minDistance := math.MaxInt
	closestName := ""

	for _, entry := range validRegisterNames {
		distance := levenshtein.ComputeDistance(name, entry.Key)
		if distance < minDistance {
			minDistance = distance
			closestName = entry.Key
		}
	}

	return closestName, minDistance
}

func RegisterToAarch64GPRegister(
	register *gen.RegisterInfo,
) (registers.GPRegister, core.ResultList) {
	gpr, ok := RegisterNameToAarch64GPRegister(register.Name)

	if !ok {
		closestName, _ := closestAarch64GPRegisterName(register.Name)

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

	return gpr, core.ResultList{}
}
