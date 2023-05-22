package main

import (
	_ "embed"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gb-inc/millage-analysis/utils"
	"github.com/xuri/excelize/v2"
)

var (
	//go:embed .\sql\bytwp-assmtincrease.sql
	newconstructionSql string

	//go:embed .\sql\total-taxable.sql
	totaltaxableSql string
)

type txblRow struct {
	TotalTxbl float64 // Necessary to calculate millages
}

type Row struct {
	TownShipBorough string
	DistrictName    string
	NewImprAssmt    float64
	OldImprAssmt    float64
	ImprDiff        float64 // NewImprAssmt - OldImprAssmt
}

func main() {
	db, err := utils.NewDB("10.0.7.24", "1433", "TaxDB")
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

	// Define "WHERE" conditions for each SD
	sdWhereMap := map[string]string{
		"WHSD":  "('030','070','150','200','230','020','090','110','130','010','050','170','210','270')",
		"WESD":  "('040','061','120','220','240','260','280')",
		"WASD":  "('080','100','180','190','273')",
		"SQSD":  "('250')",
		"NPSD":  "('140')",
		"FCRSD": "('062','160')",
	}

	// Query database
	for sdName, sdWhere := range sdWhereMap {
		// Create a new variable to store the modified query
		modifiedSql := strings.Replace(newconstructionSql, "GROUP BY", "WHERE bpic.TownShipBorough IN "+sdWhere+" GROUP BY", 1)

		// Get total taxable amount for each SD
		modifiedtxblSql := totaltaxableSql + " AND P.TownShipBorough IN " + sdWhere
		rowsTxbl, err := tx.Query(modifiedtxblSql)
		if err != nil {
			log.Fatal(err)
		}
		defer rowsTxbl.Close()

		// Get assessment increase for each SD
		rows, err := tx.Query(modifiedSql)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		// Parse rows into struct
		var txblData []txblRow
		for rowsTxbl.Next() {
			var r txblRow
			if err := rowsTxbl.Scan(&r.TotalTxbl); err != nil {
				log.Fatal(err)
			}
			txblData = append(txblData, r)
		}
		if err := rowsTxbl.Err(); err != nil {
			log.Fatal(err)
		}

		var data []Row
		for rows.Next() {
			var r Row
			if err := rows.Scan(&r.TownShipBorough, &r.DistrictName, &r.NewImprAssmt, &r.OldImprAssmt, &r.ImprDiff); err != nil {
				log.Fatal(err)
			}
			data = append(data, r)
		}
		if err := rows.Err(); err != nil {
			log.Fatal(err)
		}

		// Pull Excel template
		templatePath := fmt.Sprintf("./templ/%s_AssmtIncrease_.xlsx", sdName)
		f, err := excelize.OpenFile(templatePath)
		if err != nil {
			log.Fatal(err)
		}

		// Populate data and cover sheets
		handleSchoolDist(f, sdName, data, txblData)

		// Save file
		saveExcelFile(f, sdName)
	}

	// Set transaction success flag
	ok = true
}

func handleSchoolDist(f *excelize.File, schoolDist string, data []Row, txblData []txblRow) {
	/* Populate "Data" sheet */
	sheetIndex := 1 // "Data" sheet
	sheetName := f.GetSheetMap()[sheetIndex]
	f.SetActiveSheet(sheetIndex)
	popDataCells(f, sheetName, data)

	/* Populate "Cover" sheet */
	sheetIndex = 2 // "Cover" sheet
	sheetName = f.GetSheetMap()[sheetIndex]
	f.SetActiveSheet(sheetIndex)
	popHeaderCells(f, sheetName, txblData)

	// Save file
	saveExcelFile(f, schoolDist)
}

func saveExcelFile(f *excelize.File, schoolDist string) {
	currTime := time.Now().Format("1504")
	var today string = time.Now().Format("2006-01-02")
	var savePath string = "C:/Users/Samuel/grandjean.net/FogBugz - Documents/"

	switch schoolDist {
	case "WHSD":
		savePath += "11520/daily/WHSD_AssmtIncrease_" + today + "_" + currTime + ".xlsx"
	case "WESD":
		savePath += "11520/daily/WESD_AssmtIncrease_" + today + "_" + currTime + ".xlsx"
	case "WASD":
		savePath += "11520/daily/WASD_AssmtIncrease_" + today + "_" + currTime + ".xlsx"
	case "SQSD":
		savePath += "11520/daily/SQSD_AssmtIncrease_" + today + "_" + currTime + ".xlsx"
	case "NPSD":
		savePath += "11520/daily/NPSD_AssmtIncrease_" + today + "_" + currTime + ".xlsx"
	case "FCRSD":
		savePath += "11520/daily/FCRSD_AssmtIncrease_" + today + "_" + currTime + ".xlsx"
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

func popHeaderCells(f *excelize.File, sheet string, txblData []txblRow) {
	switch sheet {
	case "Cover":
		var today string = time.Now().Format("01/02/2006")
		f.SetCellValue(sheet, "A2", "Parcels w/ New Construction as of "+today+":")
		f.SetCellValue(sheet, "A3", "Parcels w/ New Construction as of 01/01/2023:")
		for i, r := range txblData {
			f.SetCellValue(sheet, fmt.Sprintf("B%d", i+8), r.TotalTxbl)
		}
	}
}
