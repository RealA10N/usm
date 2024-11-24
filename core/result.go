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

type Result interface {
	GetType() ResultType
	GetMessage(ctx SourceContext) string
	GetLocation() *UnmanagedSourceView

	// Results may wrap other results, and this method returns the next result
	// in the chain. The first result should be logically the most significant,
	// upper level one, and results further down usually describe the problem
	// in more detail or suggest a quick fix.
	GetNext() Result
}

type ResultList = list.List[Result]

// MARK: GenericResult

type GenericResult struct {
	Type     ResultType
	Message  string
	Location *UnmanagedSourceView
	Next     Result
}

func (r GenericResult) GetType() ResultType {
	return r.Type
}

func (r GenericResult) GetMessage(SourceContext) string {
	return r.Message
}

func (r GenericResult) GetLocation() *UnmanagedSourceView {
	return r.Location
}

func (r GenericResult) GetNext() Result {
	return r.Next
}
