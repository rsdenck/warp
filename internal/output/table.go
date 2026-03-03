package output

import (
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
)

// NewTable creates a new table with standard styling
func NewTable(title string) table.Writer {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleRounded)
	
	if title != "" {
		t.SetTitle(title)
	}
	
	return t
}

// RenderSimpleTable renders a simple table with headers and rows
func RenderSimpleTable(title string, headers []interface{}, rows [][]interface{}) {
	t := NewTable(title)
	
	if len(headers) > 0 {
		t.AppendHeader(table.Row(headers))
	}
	
	for _, row := range rows {
		t.AppendRow(table.Row(row))
	}
	
	t.Render()
}
