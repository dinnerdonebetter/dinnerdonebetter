package main

type QueryType string

const (
	Exec       QueryType = ":exec"
	ExecResult QueryType = ":execresult"
	ExecRows   QueryType = ":execrows"
	ExecLastID QueryType = ":execlastid"
	Many       QueryType = ":many"
	One        QueryType = ":one"
	BatchExec  QueryType = ":batchexec"
	BatchMany  QueryType = ":batchmany"
	BatchOne   QueryType = ":batchone"
)

type QueryAnnotation struct {
	Name string
	Type QueryType
}

type Query struct {
	Content    string
	Annotation QueryAnnotation
}
