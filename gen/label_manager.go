package gen

import "alon.kr/x/usm/core"

type LabelManager[InstT BaseInstruction] interface {
	GetLabel(name string) *LabelInfo[InstT]
	NewLabel(info *LabelInfo[InstT]) core.Result
}
