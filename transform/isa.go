package transform

import (
	"strings"

	"alon.kr/x/faststringmap"
	"alon.kr/x/usm/gen"
)

type InstructionSet struct {
	Name        string
	Aliases     []string
	Description string
	Extensions  []string

	GenerationContext gen.GenerationContext

	Transformations TransformationCollection
}

type InstructionSetCollection struct {
	InstructionSets      []*InstructionSet
	NameToInstructionSet faststringmap.Map[*InstructionSet]
}

func NewIsaCollection(
	instructionSets ...*InstructionSet,
) *InstructionSetCollection {
	entries := []faststringmap.MapEntry[*InstructionSet]{}
	for _, set := range instructionSets {
		entries = append(entries, faststringmap.MapEntry[*InstructionSet]{
			Key:   set.Name,
			Value: set,
		})

		for _, alias := range set.Aliases {
			entries = append(entries, faststringmap.MapEntry[*InstructionSet]{
				Key:   alias,
				Value: set,
			})
		}
	}

	return &InstructionSetCollection{
		InstructionSets:      instructionSets,
		NameToInstructionSet: faststringmap.NewMap(entries),
	}
}

func (c *InstructionSetCollection) TransformationNames() []string {
	names := []string{}
	for _, set := range c.InstructionSets {
		names = append(names, set.Transformations.Names()...)
	}
	return names
}

func (c *InstructionSetCollection) Traverse(
	start *InstructionSet,
	transformations []string,
) *InstructionSet {
	isa := start
	for _, transName := range transformations {
		trans, ok := isa.Transformations.NameToTransformation.LookupString(transName)
		if !ok {
			return nil
		}

		isaName := trans.Target
		isa, ok = c.NameToInstructionSet.LookupString(isaName)
		if !ok {
			return nil
		}
	}

	return isa
}

// Infer the instruction set from the filename extension.
// Returns nil if no matching instruction set is found.
// If multiple instruction sets match the filename, the one that matches with
// the longest prefix is returned.
func (c *InstructionSetCollection) FilenameToInstructionSet(
	filename string,
) *InstructionSet {
	// TODO: this can be implemented in linear time with, for example, a trie of
	// reversed prefixes, and traversal of the filename from the end. (faststringmap)

	longest := 0
	var longestSet *InstructionSet

	for _, set := range c.InstructionSets {
		for _, ext := range set.Extensions {
			if strings.HasSuffix(filename, ext) {
				if len(ext) > longest {
					longest = len(ext)
					longestSet = set
				}
			}
		}
	}

	return longestSet
}
