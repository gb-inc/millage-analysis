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
	//go:embed .\sql\bytwp-assmtincrease.sql
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

	// Parse WHSD rows into struct
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

	// Pull WHSD Excel template
	f1, err := excelize.OpenFile("./templ/WHSD_AssmtIncrease_.xlsx")
	if err != nil {
		log.Fatal(err)
		return
	}
	var schoolDist string = "WHSD"
	handleSchoolDist(f1, schoolDist, data)

	// Pull WESD Excel template
	f2, err := excelize.OpenFile("./templ/WESD_AssmtIncrease_.xlsx")
	if err != nil {
		log.Fatal(err)
		return
	}
	schoolDist = "WESD"
	handleSchoolDist(f2, schoolDist, data)

	// Pull WASD Excel template
	f3, err := excelize.OpenFile("./templ/WASD_AssmtIncrease_.xlsx")
	if err != nil {
		log.Fatal(err)
		return
	}
	schoolDist = "WASD"
	handleSchoolDist(f3, schoolDist, data)

	// Pull SQSD Excel template
	f4, err := excelize.OpenFile("./templ/SQSD_AssmtIncrease_.xlsx")
	if err != nil {
		log.Fatal(err)
		return
	}
	schoolDist = "SQSD"
	handleSchoolDist(f4, schoolDist, data)

	// Pull NPSD Excel template
	f5, err := excelize.OpenFile("./templ/NPSD_AssmtIncrease_.xlsx")
	if err != nil {
		log.Fatal(err)
		return
	}
	schoolDist = "NPSD"
	handleSchoolDist(f5, schoolDist, data)

	// Pull FCRSD Excel template
	f6, err := excelize.OpenFile("./templ/FCRSD_AssmtIncrease_.xlsx")
	if err != nil {
		log.Fatal(err)
		return
	}
	schoolDist = "FCRSD"
	handleSchoolDist(f6, schoolDist, data)

	// Set transaction success flag
	ok = true
}

func handleSchoolDist(f *excelize.File, schoolDist string, data []Row) {
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

	// Save file
	saveExcelFile(f, schoolDist)
}

func saveExcelFile(f *excelize.File, schoolDist string) {
	currTime := time.Now().Format("1504")
	var today string = time.Now().Format("2006-01-02")
	var savePath string = "C:/Users/Samuel/grandjean.net/FogBugz - Documents/"

	switch schoolDist {
	case "WHSD":
		savePath += "11466/daily/WHSD_AssmtIncrease_" + today + "_" + currTime + ".xlsx"
	case "WESD":
		savePath += "11528/daily/WESD_AssmtIncrease_" + today + "_" + currTime + ".xlsx"
	case "WASD":
		savePath += "11529/daily/WASD_AssmtIncrease_" + today + "_" + currTime + ".xlsx"
	case "SQSD":
		savePath += "11524/daily/SQSD_AssmtIncrease_" + today + "_" + currTime + ".xlsx"
	case "NPSD":
		savePath += "11530/daily/NPSD_AssmtIncrease_" + today + "_" + currTime + ".xlsx"
	case "FCRSD":
		savePath += "11521/daily/FCRSD_AssmtIncrease_" + today + "_" + currTime + ".xlsx"
	}

	if err := f.SaveAs(savePath); err != nil {
		log.Fatal(err)
	}
	fmt.Print(schoolDist + " saved at " + currTime + ".\n")
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
