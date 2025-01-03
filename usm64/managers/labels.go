package managers

import (
	"crypto/rand"
	"encoding/hex"

	"alon.kr/x/usm/core"
	"alon.kr/x/usm/gen"
)

type LabelMap map[string]*gen.LabelInfo

func (m *LabelMap) GetLabel(name string) *gen.LabelInfo {
	val, ok := (*m)[name]
	if !ok {
		return nil
	}
	return val
}

func (m *LabelMap) NewLabel(label *gen.LabelInfo) core.Result {
	(*m)[label.Name] = label
	return nil
}

func generateRandomLabelName() (string, error) {
	b := make([]byte, 4)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return ".L" + hex.EncodeToString(b), nil
}

func (m *LabelMap) GenerateLabel(block *gen.BasicBlockInfo) *gen.LabelInfo {
	name, err := generateRandomLabelName()
	for err != nil || m.GetLabel(name) != nil {
		name, err = generateRandomLabelName()
	}

	label := &gen.LabelInfo{
		Name:       name,
		BasicBlock: block,
	}

	return label
}

func NewLabelManager() gen.LabelManager {
	return &LabelMap{}
}
