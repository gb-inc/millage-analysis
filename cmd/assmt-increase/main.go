package main

import (
	_ "embed"
	"fmt"
	"log"
	"time"

	"github.com/gb-inc/millage-analysis/utils"
	"github.com/xuri/excelize/v2"
)

var (
	//go:embed bytwp-assmtincrease.sql
	valuationchangeSql string
)

type Row struct {
	TownShipBorough string
	OldLandAssmt    float64
	OldImprAssmt    float64
	NewLandAssmt    float64
	NewImprAssmt    float64
	LandDiff        float64 // NewLandAssmt - OldLandAssmt
	ImprDiff        float64 // NewImprAssmt - OldImprAssmt

}

func main() {
	db, err := utils.NewDB("localhost", "1433", "TaxDB_Dev")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Start transaction
	tx, err := db.Begin()
	if err != nil {
		_ = db.Close()
		log.Fatal(err)
	}
	var ok bool
	defer utils.HandleTxFunc(tx, &ok)

	// Query database
	rows, err := tx.Query(valuationchangeSql)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Parse rows into struct
	var data []Row
	for rows.Next() {
		var r Row
		if err := rows.Scan(&r.TownShipBorough, &r.OldLandAssmt, &r.OldImprAssmt, &r.NewLandAssmt, &r.NewImprAssmt, &r.LandDiff, &r.ImprDiff); err != nil {
			log.Fatal(err)
		}
		data = append(data, r)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	// Generate Excel report
	f := excelize.NewFile()
	var today string = time.Now().Format("2006-01-02")
	sheetName := "ValuationChanges_" + today
	index, err := f.NewSheet(sheetName)
	if err != nil {
		log.Fatal(err)
	}
	f.SetActiveSheet(index)
	popHeaderCells(f, sheetName)
	popDataCells(f, sheetName, data)
	utils.AutoFitColumns(f, sheetName, "A", "G")

	if err := f.SaveAs("./" + sheetName + ".xlsx"); err != nil {
		log.Fatal(err)
	}

	// Set transaction success flag
	ok = true
}

func popDataCells(f *excelize.File, sheet string, data []Row) {
	for i, r := range data {
		f.SetCellValue(sheet, fmt.Sprintf("A%d", i+2), r.TownShipBorough)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", i+2), r.OldLandAssmt)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", i+2), r.OldImprAssmt)
		f.SetCellValue(sheet, fmt.Sprintf("D%d", i+2), r.NewLandAssmt)
		f.SetCellValue(sheet, fmt.Sprintf("E%d", i+2), r.NewImprAssmt)
		f.SetCellValue(sheet, fmt.Sprintf("F%d", i+2), r.LandDiff)
		f.SetCellValue(sheet, fmt.Sprintf("G%d", i+2), r.ImprDiff)

		// Set number format of columns B to G to "$#,##0.00"
		styleID, err := f.NewStyle(&excelize.Style{
			NumFmt: 3,
		})
		if err != nil {
			log.Fatal(err)
		}
		for j := 1; j <= 6; j++ {
			f.SetCellStyle(sheet, fmt.Sprintf("%c%d", 'B'+j-1, i+2), fmt.Sprintf("%c%d", 'B'+j-1, i+2), styleID)
		}
	}

	fmtDataCells(f, sheet)
}

func fmtDataCells(f *excelize.File, sheet string) {
	numRows, err := f.GetRows(sheet)
	if err != nil {
		log.Fatal(err)
	}
	for i := 2; i <= len(numRows); i++ {
		utils.FmtDataCell(f, "#E2EFDA", sheet, []string{fmt.Sprintf("D%d", i), fmt.Sprintf("E%d", i)}) // Green - indicates new valuation
		utils.FmtDataCell(f, "#FCE4D6", sheet, []string{fmt.Sprintf("F%d", i), fmt.Sprintf("G%d", i)}) // Yellow - indicates change in valuation
	}

	err = f.AutoFilter(sheet, "A1:G1", []excelize.AutoFilterOptions{})
	if err != nil {
		fmt.Println(err)
	}
}

func popHeaderCells(f *excelize.File, sheet string) {
	f.SetCellValue(sheet, "A1", "Township Borough")
	f.SetCellValue(sheet, "B1", "Old Land Assmt")
	f.SetCellValue(sheet, "C1", "Old Impr Assmt")
	f.SetCellValue(sheet, "D1", "New Land Assmt")
	f.SetCellValue(sheet, "E1", "New Impr Assmt")
	f.SetCellValue(sheet, "F1", "Land Diff")
	f.SetCellValue(sheet, "G1", "Impr Diff")
	fmtHeaderCells(f, sheet)
}
func fmtHeaderCells(f *excelize.File, sheet string) {

	var newvalcells []string = []string{"D1", "E1"}
	var chgvalcells []string = []string{"F1", "G1"}
	var othercells []string = []string{"A1", "B1", "C1"}

	newvalID, err := f.NewStyle(&excelize.Style{ // Bold Header Green - indicates new valuation
		Fill: excelize.Fill{
			Type:    "pattern",
			Pattern: 1,
			Color:   []string{"#E2EFDA"},
		},
		Font: &excelize.Font{
			Bold: true,
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
	})

	if err != nil {
		fmt.Println(err)
	}

	chgvalID, err := f.NewStyle(&excelize.Style{ // Bold Header Yellow - indicates change in valuation
		Fill: excelize.Fill{
			Type:    "pattern",
			Pattern: 1,
			Color:   []string{"#FCE4D6"},
		},
		Font: &excelize.Font{
			Bold: true,
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
	})

	if err != nil {
		fmt.Println(err)
	}

	//Apply style to column
	for _, cell := range newvalcells {
		f.SetCellStyle(sheet, cell, cell, newvalID)
	}
	for _, cell := range chgvalcells {
		f.SetCellStyle(sheet, cell, cell, chgvalID)
	}

	utils.BoldCells(f, sheet, othercells)

}
