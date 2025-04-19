package aarch64translation

import (
	"fmt"
	"math"

	"alon.kr/x/aarch64codegen/immediates"
)

func int64toInt32(value int64) (int32, bool) {
	if value < math.MinInt32 || value > math.MaxInt32 {
		return 0, false
	}
	return int32(value), true
}

func Uint64DiffToOffset26Align4(
	dst, src uint64,
) (immediates.Offset26Align4, error) {
	diff := int64(dst) - int64(src)
	diff32, ok := int64toInt32(diff)
	if !ok {
		return 0, fmt.Errorf("offset %v does not fit in 26 bits", diff)
	}

	return immediates.NewOffset26Align4(diff32)
}
