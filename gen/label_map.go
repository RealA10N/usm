package gen

import (
	"crypto/rand"
	"encoding/hex"

	"alon.kr/x/usm/core"
)

// LabelMap is a naive and default implementation of the LabelManger interface.
// It uses a simple go builtin map to store label information, using label names
// as keys.
type LabelMap map[string]*LabelInfo

func (m *LabelMap) GetLabel(name string) *LabelInfo {
	val, ok := (*m)[name]
	if !ok {
		return nil
	}
	return val
}

func (m *LabelMap) NewLabel(label *LabelInfo) core.ResultList {
	(*m)[label.Name] = label
	return core.ResultList{}
}

func generateRandomLabelName() (string, error) {
	b := make([]byte, 4)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return ".L" + hex.EncodeToString(b), nil
}

func (m *LabelMap) GenerateLabel() *LabelInfo {
	name, err := generateRandomLabelName()
	for err != nil || m.GetLabel(name) != nil {
		name, err = generateRandomLabelName()
	}

	return &LabelInfo{Name: name}
}

func NewLabelMap(*FileGenerationContext) LabelManager {
	return &LabelMap{}
}
