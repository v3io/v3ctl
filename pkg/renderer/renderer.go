package renderer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/nuclio/errors"
	"gopkg.in/yaml.v3"
)

type Renderer struct {
	output io.Writer
}

func NewRenderer(output io.Writer) *Renderer {
	return &Renderer{
		output: output,
	}
}

func (r *Renderer) RenderTable(header []string, records [][]string) {
	tw := table.NewWriter()
	tw.SetOutputMirror(r.output)
	tw.SetStyle(table.Style{
		Name: "Nuclio",
		Box: table.BoxStyle{
			MiddleVertical: "|",
			PaddingLeft:    " ",
			PaddingRight:   " ",
		},
		Options: table.Options{
			DoNotColorBordersAndSeparators: true,
			DrawBorder:                     false,
			SeparateColumns:                true,
			SeparateFooter:                 false,
			SeparateHeader:                 false,
			SeparateRows:                   false,
		},
		Color:  table.ColorOptionsDefault,
		Format: table.FormatOptionsDefault,
		HTML:   table.DefaultHTMLOptions,
		Title:  table.TitleOptionsDefault,
	})
	tw.AppendHeader(r.rowStringToTableRow(header), table.RowConfig{})
	tw.AppendRows(r.rowsStringToTableRows(records), table.RowConfig{})
	tw.Render()
}

func (r *Renderer) RenderYAML(items interface{}) error {
	body, err := yaml.Marshal(items)
	if err != nil {
		return errors.Wrap(err, "Failed to render YAML")
	}

	fmt.Fprintln(r.output, string(body)) // nolint: errcheck

	return nil
}

func (r *Renderer) RenderJSON(items interface{}) error {
	body, err := json.Marshal(items)
	if err != nil {
		return errors.Wrap(err, "Failed to render JSON")
	}

	var pbody bytes.Buffer
	if err := json.Indent(&pbody, body, "", "\t"); err != nil {
		return errors.Wrap(err, "Failed to indent JSON")
	}

	fmt.Fprintln(r.output, pbody.String()) // nolint: errcheck

	return nil
}

func (r *Renderer) rowsStringToTableRows(rows [][]string) []table.Row {
	tableRows := make([]table.Row, len(rows))
	for rowIndex, rowValue := range rows {
		tableRows[rowIndex] = r.rowStringToTableRow(rowValue)
	}
	return tableRows
}

func (r *Renderer) rowStringToTableRow(row []string) table.Row {
	tableRow := make(table.Row, len(row))
	for cellIndex, cellValue := range row {
		tableRow[cellIndex] = cellValue
	}
	return tableRow
}
