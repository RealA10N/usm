package gen

import "alon.kr/x/usm/core"

type LabelManager interface {
	GetLabel(name string) *LabelInfo
	NewLabel(info LabelInfo) core.Result
}
