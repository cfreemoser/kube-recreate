package util

import (
	"io"

	table "github.com/jedib0t/go-pretty/v6/table"
)

type Reporter struct {
	t table.Writer
}

func NewReporter(mirror io.Writer) *Reporter {
	r := &Reporter{
		t: table.NewWriter(),
	}

	r.t.SetOutputMirror(mirror)

	r.t.AppendHeader(table.Row{"Name", "Type", "Operation"})
	return r
}

func (r *Reporter) Append(name, resourceType, operation string) {
	r.t.AppendRow(table.Row{name, resourceType, operation})
}

func (r *Reporter) AddSeperator() {
	r.t.AppendSeparator()
}

func (r *Reporter) PrintReport() {
	r.t.Render()
}
