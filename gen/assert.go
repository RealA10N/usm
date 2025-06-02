package gen

import (
	"fmt"
	"math/big"
	"strconv"

	"alon.kr/x/list"
	"alon.kr/x/usm/core"
)

// MARK: Arguments

func AssertAtLeastArguments(
	info *InstructionInfo,
	atLeast int,
) core.ResultList {
	if len(info.Arguments) < atLeast {
		return list.FromSingle(core.Result{
			{
				Type:     core.ErrorResult,
				Message:  fmt.Sprintf("Expected at least %d arguments", atLeast),
				Location: info.Declaration,
			},
		})
	}

	return core.ResultList{}
}

func AssertAtMostArguments(
	info *InstructionInfo,
	atMost int,
) core.ResultList {
	if len(info.Arguments) > atMost {
		return list.FromSingle(core.Result{
			{
				Type:     core.ErrorResult,
				Message:  fmt.Sprintf("Expected at most %d arguments", atMost),
				Location: info.Declaration,
			},
		})
	}

	return core.ResultList{}
}

func AssertArgumentsBetween(
	info *InstructionInfo,
	atLeast int,
	atMost int,
) core.ResultList {
	if len(info.Arguments) < atLeast || len(info.Arguments) > atMost {
		return list.FromSingle(core.Result{
			{
				Type: core.ErrorResult,
				Message: fmt.Sprintf(
					"Expected between %d and %d arguments",
					atLeast,
					atMost,
				),
				Location: info.Declaration,
			},
		})
	}

	return core.ResultList{}
}

func AssertArgumentsExactly(
	info *InstructionInfo,
	count int,
) core.ResultList {
	if len(info.Arguments) != count {
		return list.FromSingle(core.Result{
			{
				Type:     core.ErrorResult,
				Message:  fmt.Sprintf("Expected %d arguments", count),
				Location: info.Declaration,
			},
		})
	}

	return core.ResultList{}
}

// MARK: Targets

func AssertTargetsExactly(
	info *InstructionInfo,
	count int,
) core.ResultList {
	if len(info.Targets) != count {
		return list.FromSingle(core.Result{
			{
				Type:     core.ErrorResult,
				Message:  fmt.Sprintf("Expected %d targets", count),
				Location: info.Declaration,
			},
		})
	}

	return core.ResultList{}
}

// MARK: Integers

func AssertBigIntInSet(
	view *core.UnmanagedSourceView,
	bigInt *big.Int,
	options []int64,
) (int64, core.ResultList) {
	var value int64
	isInvalid := !bigInt.IsInt64()

	if isInvalid {
		goto fail
	}

	value = bigInt.Int64()
	for _, option := range options {
		if value == option {
			return value, core.ResultList{}
		}
	}

fail:
	message := "Expected one of "
	message += "#" + strconv.FormatInt(options[0], 10)
	for _, option := range options[1:] {
		message += ", #" + strconv.FormatInt(option, 10)
	}

	return 0, list.FromSingle(core.Result{
		{
			Type:     core.ErrorResult,
			Message:  message,
			Location: view,
		},
	})
}
