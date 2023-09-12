package main

type QueryType string

const (
	ExecType       QueryType = ":exec"
	ExecResultType QueryType = ":execresult"
	ExecRowsType   QueryType = ":execrows"
	ExecLastIDType QueryType = ":execlastid"
	ManyType       QueryType = ":many"
	OneType        QueryType = ":one"
	BatchExecType  QueryType = ":batchexec"
	BatchManyType  QueryType = ":batchmany"
	BatchOneType   QueryType = ":batchone"
)

type QueryAnnotation struct {
	Name string
	Type QueryType
}

type Query struct {
	Content    string
	Annotation QueryAnnotation
}
