package utils

import (
	"fmt"

	"github.com/xuri/excelize/v2"
)

func AutoFitColumns(f *excelize.File, sheet, startCol, endCol string) {
	f.SetColWidth(sheet, startCol, endCol, 15)
}

// Typical formatting I use for data (non-header) cells
// Have the option of color fill type
func FmtDataCell(f *excelize.File, fillColor, sheet string, cells []string) {
	//Color must be valid Hexadecimal
	styleID, err := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{
			Type:    "pattern",
			Pattern: 1,
			Color:   []string{fillColor},
		},
		Font: &excelize.Font{
			Color: "#FFFFFF",
		},
		Border: []excelize.Border{
			{
				Type:  "left",
				Color: "#B2B2B2",
				Style: 1,
			},
			{
				Type:  "top",
				Color: "#B2B2B2",
				Style: 1,
			},
			{
				Type:  "right",
				Color: "#B2B2B2",
				Style: 1,
			},
			{
				Type:  "bottom",
				Color: "#B2B2B2",
				Style: 1,
			},
		},
		NumFmt: 3,
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	//Apply style to column
	for _, cell := range cells {
		f.SetCellStyle(sheet, cell, cell, styleID)
	}
}

func BoldCells(f *excelize.File, sheet string, cells []string) {
	styleID, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
		},
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	// Apply style to the header cells
	for _, cell := range cells {
		f.SetCellStyle(sheet, cell, cell, styleID)
	}
}
