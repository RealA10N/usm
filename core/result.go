package core

import (
	"fmt"
	"sort"
	"strings"

	"alon.kr/x/list"
	"github.com/fatih/color"
)

// MARK: Result

type ResultType uint8

const (
	ErrorResult = ResultType(iota)
	InternalErrorResult
	WarningResult
	HintResult
)

type Result []ResultDetails

type ResultDetails struct {
	Type     ResultType
	Message  string
	Location *UnmanagedSourceView
}

type ResultList = list.List[Result]

type ResultStringer struct {
	SourceContext
	Titles             map[ResultType]string
	SourceErrorPointer string
	Filepath           string
	LineStarts         []UsmUint
}

func calculateLineStartsFromSource(ctx SourceContext) []UsmUint {
	starts := []UsmUint{0}
	for i, c := range ctx {
		if c == '\n' {
			starts = append(starts, UsmUint(i+1))
		}
	}
	return starts
}

func NewResultStringer(ctx SourceContext, Filepath string) ResultStringer {
	return ResultStringer{
		SourceContext: ctx,
		Titles: map[ResultType]string{
			InternalErrorResult: color.New(color.Bold, color.BgRed, color.FgWhite).Sprint(" panic "),
			ErrorResult:         color.New(color.Bold, color.FgRed).Sprint("error:"),
			WarningResult:       color.New(color.Bold, color.FgYellow).Sprint("warning:"),
			HintResult:          color.New(color.Bold, color.FgCyan).Sprint("note:"),
		},
		SourceErrorPointer: color.New(color.Bold, color.FgMagenta).Sprint("^"),
		LineStarts:         calculateLineStartsFromSource(ctx),
		Filepath:           Filepath,
	}
}

func (w *ResultStringer) viewToLocation(
	view UnmanagedSourceView,
) (line UsmUint, col UsmUint) {
	start := view.Start
	row := sort.Search(len(w.LineStarts), func(i int) bool {
		return start < w.LineStarts[i]
	}) - 1
	col = start - w.LineStarts[row]
	return UsmUint(row), UsmUint(col)
}

func (w *ResultStringer) getLine(row UsmUint) string {
	lineStart := w.LineStarts[row]
	var lineEnd UsmUint
	if row >= UsmUint(len(w.LineStarts)-1) {
		lineEnd = UsmUint(len(w.SourceContext))
	} else {
		lineEnd = w.LineStarts[row+1] - 1
	}
	return string(w.SourceContext[lineStart:lineEnd])
}

func (w *ResultStringer) stringLineContext(row, col UsmUint) string {
	firstLinePad := fmt.Sprintf("%*d", 5, row+1)
	border := " | "
	line := w.getLine(row)

	secondLinePad := strings.Repeat(" ", len(firstLinePad))
	pointerLine := strings.Repeat(" ", int(col)) + w.SourceErrorPointer

	firstLine := firstLinePad + border + line
	secondLine := secondLinePad + border + pointerLine
	return firstLine + "\n" + secondLine + "\n"
}

func (w *ResultStringer) StringResultDetails(details ResultDetails) string {
	locationPrefix := ""
	locationSuffix := ""
	if details.Location != nil {
		row, col := w.viewToLocation(*details.Location)
		locationPrefix = fmt.Sprintf("%s:%d:%d: ", w.Filepath, row+1, col+1)
		locationSuffix = w.stringLineContext(row, col)
	}

	title := w.Titles[details.Type]
	message := details.Message
	return locationPrefix + title + " " + message + "\n" + locationSuffix
}

func (w *ResultStringer) StringResult(result Result) string {
	s := ""
	for _, details := range result {
		s += w.StringResultDetails(details)
	}
	return s
}
