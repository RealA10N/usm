package core

import (
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
	titles map[ResultType]string
}

func NewResultStringer() ResultStringer {
	return ResultStringer{
		titles: map[ResultType]string{
			InternalErrorResult: color.New(color.Bold, color.BgRed, color.FgWhite).Sprint("panic:"),
			ErrorResult:         color.New(color.Bold, color.FgRed).Sprint("error:"),
			WarningResult:       color.New(color.Bold, color.FgYellow).Sprint("warning:"),
			HintResult:          color.New(color.Bold, color.FgCyan).Sprint("note:"),
		},
	}
}

func (w *ResultStringer) StringResultDetails(details ResultDetails) string {
	title := w.titles[details.Type]
	return title + " " + details.Message
}

func (w *ResultStringer) StringResult(result Result) string {
	s := ""
	for _, details := range result {
		s += w.StringResultDetails(details) + "\n"
	}
	return s
}
