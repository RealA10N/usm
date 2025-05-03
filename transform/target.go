package transform

import (
	"fmt"
	"strings"

	"alon.kr/x/faststringmap"
	"alon.kr/x/list"
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/gen"
)

type Target struct {
	// The name(s) of the target.
	// The convention is to use lowercase letters, with '-' as a separator, if
	// required.
	// The first name in the list is the default name, which is used and
	// displayed in the CLI. Other names in the list are valid names when the
	// user inputs them. At least one name must be provided.
	Names []string

	// A short description of the target.
	// This is used to enhance the user experience in the CLI, and possibly,
	// in other places.
	Description string

	// The list of file extensions (suffixes) that represent this target.
	// The first extension in the list is the default extension, which is
	// used as the default extension for output (generated) files.
	// At least one extension must be provided.
	// Other extensions in the list are used to automatically detect the target
	// type from the input file name.
	Extensions []string

	// An optional generation context for this target.
	// This is used to parse and generate the target code from a USM source.
	// If nil, assumes that the target cannot be generated from a USM source.
	GenerationContext *gen.GenerationContext

	Transformations TransformationCollection
}

type TargetCollection struct {
	Targets      []*Target
	NameToTarget faststringmap.Map[*Target]
}

func NewTargetCollection(
	targets ...*Target,
) *TargetCollection {
	entries := []faststringmap.MapEntry[*Target]{}
	for _, t := range targets {
		for _, name := range t.Names {
			entries = append(entries, faststringmap.MapEntry[*Target]{
				Key:   name,
				Value: t,
			})
		}
	}

	return &TargetCollection{
		Targets:      targets,
		NameToTarget: faststringmap.NewMap(entries),
	}
}

func (c *TargetCollection) TransformationNames() []string {
	names := []string{}
	for _, t := range c.Targets {
		names = append(names, t.Transformations.Names()...)
	}
	return names
}

// Traverse the target collection according to the provided transformation names.
// This does not actually perform the transformations, but checks for the
// legality of the transformation chain.
func (c *TargetCollection) Traverse(
	start *Target,
	transformationNames []string,
) (
	transformations []*Transformation,
	end *Target,
	results core.ResultList,
) {
	end = start
	for _, transName := range transformationNames {
		transformation, ok := end.Transformations.NameToTransformation.LookupString(transName)
		if !ok {
			return nil, nil, list.FromSingle(core.Result{
				{
					Type: core.ErrorResult,
					Message: fmt.Sprintf(
						"Transformation \"%s\" not supported on target \"%s\"",
						transName,
						end.Names[0],
					),
				},
			})
		}

		transformations = append(transformations, transformation)

		end, ok = c.NameToTarget.LookupString(transformation.TargetName)
		if !ok {
			return nil, nil, list.FromSingle(core.Result{
				{
					Type: core.InternalErrorResult,
					Message: fmt.Sprintf(
						"Target \"%s\" does not exist",
						transformation.TargetName,
					),
				},
			})
		}
	}

	return
}

func (c *TargetCollection) Transform(
	data *TargetData,
	transformationNames []string,
) (end *TargetData, results core.ResultList) {
	transformations, _, results := c.Traverse(data.Target, transformationNames)
	if !results.IsEmpty() {
		return nil, results
	}

	for _, transformation := range transformations {
		data, results = transformation.Transform(data)
		if !results.IsEmpty() {
			return nil, results
		}

		targetName := transformation.TargetName
		target, ok := c.NameToTarget.LookupString(targetName)
		if !ok {
			return nil, list.FromSingle(core.Result{
				{
					Type: core.InternalErrorResult,
					Message: fmt.Sprintf(
						"Target \"%s\" does not exist",
						targetName,
					),
				},
			})
		}

		data.Target = target
	}

	return data, core.ResultList{}
}

// Infer the target from the filename extension.
// Returns nil if no matching target is found.
// If multiple targets match the filename, the one that matches with the longest
// prefix is returned.
func (c *TargetCollection) FilenameToTarget(
	filename string,
) (longestTarget *Target, longestExt string) {
	// TODO: this can be implemented in linear time with, for example, a trie of
	// reversed prefixes, and traversal of the filename from the end. (faststringmap)

	for _, t := range c.Targets {
		for _, ext := range t.Extensions {
			if strings.HasSuffix(filename, ext) {
				if len(ext) > len(longestExt) {
					longestExt = ext
					longestTarget = t
				}
			}
		}
	}

	return longestTarget, longestExt
}
