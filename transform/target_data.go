package transform

import (
	"bytes"
	"io"

	"alon.kr/x/usm/gen"
)

type TargetData struct {
	// A target instance that describes current target type.
	*Target

	// A representation of the target data as code.
	// If nil, an artifact representation is used instead, and must be provided.
	Code *gen.FileInfo

	// A representation of the target data as a binary artifact.
	// If not nil, this is the preferred representation of the target data,
	// and the code representation is ignored.
	Artifact *bytes.Buffer
}

func NewTargetData(target *Target, code *gen.FileInfo) *TargetData {
	return &TargetData{
		Target: target,
		Code:   code,
	}
}

func (d *TargetData) WriteTo(writer io.Writer) (int64, error) {
	if d.Artifact != nil {
		return d.Artifact.WriteTo(writer)
	}

	str := d.Code.String()
	n, err := writer.Write([]byte(str))
	return int64(n), err
}
