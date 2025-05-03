package transform

import (
	"alon.kr/x/faststringmap"
	"alon.kr/x/usm/core"
)

type DoTransform func(*TargetData) (*TargetData, core.ResultList)

type Transformation struct {
	Names       []string
	Description string

	// The name of the target of this transformation.
	TargetName string

	Transform DoTransform
}

type TransformationCollection struct {
	Transformations      []*Transformation
	NameToTransformation faststringmap.Map[*Transformation]
}

func (c *TransformationCollection) Names() []string {
	names := []string{}
	for _, t := range c.Transformations {
		names = append(names, t.Names...)
	}
	return names
}

func NewTransformationCollection(transformations ...*Transformation) *TransformationCollection {
	entries := []faststringmap.MapEntry[*Transformation]{}
	for _, t := range transformations {
		for _, name := range t.Names {
			entries = append(entries, faststringmap.MapEntry[*Transformation]{
				Key:   name,
				Value: t,
			})
		}
	}

	return &TransformationCollection{
		Transformations:      transformations,
		NameToTransformation: faststringmap.NewMap(entries),
	}
}
