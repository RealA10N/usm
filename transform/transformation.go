package transform

import (
	"alon.kr/x/faststringmap"
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/gen"
)

type Transformation struct {
	Name        string
	Aliases     []string
	Description string

	// The name of the target ISA of this transformation.
	Target string

	// Transform takes a file in the input ISA, and returns a file in the
	// output ISA.
	//
	// The transformation may modify, and possibly invalidate the input
	// structure.
	Transform func(*gen.FileInfo) (*gen.FileInfo, core.ResultList)
}

type TransformationCollection struct {
	Transformations      []*Transformation
	NameToTransformation faststringmap.Map[*Transformation]
}

func (c *TransformationCollection) Names() []string {
	names := []string{}
	for _, t := range c.Transformations {
		names = append(names, t.Name)
		names = append(names, t.Aliases...)
	}
	return names
}

func NewTransformationCollection(transformations ...*Transformation) *TransformationCollection {
	entries := []faststringmap.MapEntry[*Transformation]{}
	for _, t := range transformations {
		entries = append(entries, faststringmap.MapEntry[*Transformation]{
			Key:   t.Name,
			Value: t,
		})
		for _, alias := range t.Aliases {
			entries = append(entries, faststringmap.MapEntry[*Transformation]{
				Key:   alias,
				Value: t,
			})
		}
	}

	return &TransformationCollection{
		Transformations:      transformations,
		NameToTransformation: faststringmap.NewMap(entries),
	}
}
