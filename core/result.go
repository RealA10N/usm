package core

import (
	"fmt"
	"sort"

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
	Titles     map[ResultType]string
	Filepath   string
	LineStarts []UsmUint
}

func newTitleStrings() map[ResultType]string {
	return map[ResultType]string{
		InternalErrorResult: color.New(color.Bold, color.BgRed, color.FgWhite).Sprint("panic:"),
		ErrorResult:         color.New(color.Bold, color.FgRed).Sprint("error:"),
		WarningResult:       color.New(color.Bold, color.FgYellow).Sprint("warning:"),
		HintResult:          color.New(color.Bold, color.FgCyan).Sprint("note:"),
	}
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
		Titles:     newTitleStrings(),
		LineStarts: calculateLineStartsFromSource(ctx),
		Filepath:   Filepath,
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

func (w *ResultStringer) StringResultDetails(details ResultDetails) string {
	location := w.Filepath
	if details.Location != nil {
		row, col := w.viewToLocation(*details.Location)
		location += fmt.Sprintf(":%d:%d", row+1, col+1)
	}

	title := w.Titles[details.Type]
	message := details.Message
	return location + ": " + title + " " + message
}

func (w *ResultStringer) StringResult(result Result) string {
	s := ""
	for _, details := range result {
		s += w.StringResultDetails(details) + "\n"
	}
	return s
}
