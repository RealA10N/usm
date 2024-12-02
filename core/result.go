package core

import "alon.kr/x/list"

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
