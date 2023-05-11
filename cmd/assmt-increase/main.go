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
	newconstructionSql string
)

type Row struct {
	TownShipBorough string
	DistrictName    string
	OldImprAssmt    float64
	NewImprAssmt    float64
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
	rows, err := tx.Query(newconstructionSql)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Parse rows into struct
	var data []Row
	for rows.Next() {
		var r Row
		if err := rows.Scan(&r.TownShipBorough, &r.DistrictName, &r.OldImprAssmt, &r.NewImprAssmt, &r.ImprDiff); err != nil {
			log.Fatal(err)
		}
		data = append(data, r)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	// Pull SD Excel template
	f, err := excelize.OpenFile("./templ/WHSD_AssmtIncrease_.xlsx")
	if err != nil {
		log.Fatal(err)
		return
	}

	var today string = time.Now().Format("2006-01-02")

	/* Populate "Data" sheet */
	sheetIndex := 1 // "Data" sheet
	sheetName := f.GetSheetMap()[sheetIndex]
	f.SetActiveSheet(sheetIndex)
	popDataCells(f, sheetName, data)

	/* Populate "Cover" sheet */
	sheetIndex = 2 // "Cover" sheet
	sheetName = f.GetSheetMap()[sheetIndex]
	f.SetActiveSheet(sheetIndex)
	popHeaderCells(f, sheetName)

	currTime := time.Now().Format("1504")
	if err := f.SaveAs("C:/Users/Samuel/grandjean.net/FogBugz - Documents/11466/daily/WHSD_AssmtIncrease_" + today + "_" + currTime + ".xlsx"); err != nil {
		log.Fatal(err)
	}

	// Set transaction success flag
	ok = true
}

func popDataCells(f *excelize.File, sheet string, data []Row) {
	for i, r := range data {
		f.SetCellValue(sheet, fmt.Sprintf("A%d", i+5), r.TownShipBorough)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", i+5), r.DistrictName)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", i+5), r.NewImprAssmt)
		f.SetCellValue(sheet, fmt.Sprintf("D%d", i+5), r.OldImprAssmt)
		f.SetCellValue(sheet, fmt.Sprintf("E%d", i+5), r.ImprDiff)

		// Set number format of columns C to E to "$#,##0.00"
		styleID, err := f.NewStyle(&excelize.Style{
			NumFmt: 3,
		})
		if err != nil {
			log.Fatal(err)
		}
		for j := 3; j <= 5; j++ {
			f.SetCellStyle(sheet, fmt.Sprintf("%c%d", 'A'+j-1, i+5), fmt.Sprintf("%c%d", 'A'+j-1, i+5), styleID)
		}
	}

	fmtDataCells(f, sheet)
}

func fmtDataCells(f *excelize.File, sheet string) {
	numRows, err := f.GetRows(sheet)
	if err != nil {
		log.Fatal(err)
	}
	for i := 5; i <= len(numRows)-1; i++ {
		utils.FmtDataCell(f, "#4472C4", sheet, []string{fmt.Sprintf("C%d", i)}) // Green - indicates new valuation
		utils.FmtDataCell(f, "#70AD47", sheet, []string{fmt.Sprintf("E%d", i)}) // Yellow - indicates change in valuation
	}
}

func popHeaderCells(f *excelize.File, sheet string) {
	switch sheet {
	case "Cover":
		var today string = time.Now().Format("01/02/2006")
		f.SetCellValue(sheet, "A2", "Parcels w/ New Construction as of "+today+":")
		f.SetCellValue(sheet, "A3", "Parcels w/ New Construction as of 01/01/2023:")
	}
}
