package excel

import (
	"fmt"
	"github.com/tealeg/xlsx"
)

// Excel 数据
type Excel struct {
	Name string
	Out  string
}

// AddHeaderSingle 添加表头数据
func AddHeaderSingle(sheet *xlsx.Sheet, header []string, length []float64, color, fgColor string) {
	row := sheet.AddRow()
	row.SetHeightCM(1.3) //设置行的高度

	for _, h := range header {
		cell := row.AddCell()
		cell.Value = h
		cell.GetStyle().Alignment = xlsx.Alignment{Horizontal: "center", Vertical: "center", WrapText: true}
		cell.GetStyle().Font = xlsx.Font{Size: 14, Name: "Arial"}
		if color != "" {
			cell.GetStyle().Font = xlsx.Font{Size: 14, Name: "Arial", Color: color}
		}
		cell.GetStyle().Border = xlsx.Border{Left: "thin", Right: "thin", Top: "thin", Bottom: "thin"}

		if fgColor != "" {
			cell.GetStyle().Fill = xlsx.Fill{PatternType: xlsx.Solid_Cell_Fill, FgColor: fgColor}
		}
	}

	// 设置每列长度
	for i := 0; i < len(row.Sheet.Cols); i++ {
		if length == nil || len(length) != len(row.Sheet.Cols) {
			row.Sheet.Cols[i].Width = 20
		} else {
			row.Sheet.Cols[i].Width = length[i]
		}
	}
}

// AddData 添加数据
func AddData(sheet *xlsx.Sheet, data [][]string) {
	for _, r := range data {
		row := sheet.AddRow()
		row.SetHeightCM(1) //设置每行的高度

		var cell *xlsx.Cell
		for _, d := range r {
			cell = row.AddCell()
			cell.Value = d
			cell.GetStyle().Alignment = xlsx.Alignment{Horizontal: "left", Vertical: "center"}
			cell.GetStyle().Font = xlsx.Font{Size: 12, Name: "Arial"}
			cell.GetStyle().Border = xlsx.Border{Left: "thin", Right: "thin", Top: "thin", Bottom: "thin"}
		}
	}
}

// AddCommonCell 添加一个通用单元格
func AddCommonCell(row *xlsx.Row, value, fontColor string) {
	cell := row.AddCell()
	cell.Value = value
	cell.GetStyle().Alignment = xlsx.Alignment{Horizontal: "left", Vertical: "center", WrapText: true}
	cell.GetStyle().Font = xlsx.Font{Size: 12, Name: "Arial", Color: fontColor}
	cell.GetStyle().Border = xlsx.Border{Left: "thin", Right: "thin", Top: "thin", Bottom: "thin"}
}

// AddHMergeCell 添加一个合并行单元格
func AddHMergeCell(row *xlsx.Row, value, fontColor string, len int) {
	cell := row.AddCell()
	cell.Value = value
	// 用于合并单元格
	for i := 0; i < len; i++ {
		cellin := row.AddCell()
		cellin.GetStyle().Border = xlsx.Border{Left: "thin", Right: "thin", Top: "thin", Bottom: "thin"}
	}
	cell.HMerge = len
	cell.GetStyle().Alignment = xlsx.Alignment{Horizontal: "left", Vertical: "center", WrapText: true}
	cell.GetStyle().Font = xlsx.Font{Size: 12, Name: "Arial", Color: fontColor}
	cell.GetStyle().Border = xlsx.Border{Left: "thin", Right: "thin", Top: "thin", Bottom: "thin"}
}

// AddLinkCell 添加一个超链接单元格
func AddLinkCell(row *xlsx.Row, link, value, fontColor string) {
	cell := row.AddCell()
	cell.SetFormula(fmt.Sprintf("=HYPERLINK(\"#%s!A1\",\"%s\")", link, value))
	cell.GetStyle().Alignment = xlsx.Alignment{Horizontal: "left", Vertical: "center", WrapText: true}
	cell.GetStyle().Font = xlsx.Font{Size: 12, Name: "Arial", Color: fontColor, Underline: true}
	cell.GetStyle().Border = xlsx.Border{Left: "thin", Right: "thin", Top: "thin", Bottom: "thin"}
}

// AddBgCell 添加一个带背景的单元格
func AddBgCell(row *xlsx.Row, value, fontColor, fgColor string) {
	cell := row.AddCell()
	cell.Value = value

	style := xlsx.NewStyle()
	style.Alignment = xlsx.Alignment{Horizontal: "left", Vertical: "center", WrapText: true}
	style.Font = xlsx.Font{Size: 12, Name: "Arial", Color: fontColor}
	style.Border = xlsx.Border{Left: "thin", Right: "thin", Top: "thin", Bottom: "thin"}
	style.Fill = xlsx.Fill{PatternType: xlsx.Solid_Cell_Fill, FgColor: fgColor}

	cell.SetStyle(style)
}
