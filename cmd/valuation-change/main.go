package main

import (
	_ "embed"
	"fmt"
	"log"

	"github.com/orklepee/millage-analysis/utils"
	"github.com/qax-os/excelize"
)

var (
	//go:embed valuationchange.sql
	valuationchangeSql string
)

type Row struct {
	TownShipBorough string
	OldLandAssmt int
	OldImprAssmt int
	NewLandAssmt int
	NewImprAssmt int
	LandDiff     int // NewLandAssmt - OldLandAssmt
	ImprDiff     int // NewImprAssmt - OldImprAssmt

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
	sheetName := "ValuationChanges"
	index := f.NewSheet(sheetName)
	f.SetActiveSheet(index)
	f.SetCellValue(sheetName, "A1", "Township Borough")
	f.SetCellValue(sheetName, "B1", "Old Land Assmt")
	f.SetCellValue(sheetName, "C1", "Old Impr Assmt")
	f.SetCellValue(sheetName, "D1", "New Land Assmt")
	f.SetCellValue(sheetName, "E1", "New Impr Assmt")
	f.SetCellValue(sheetName, "F1", "Land Diff")
	f.SetCellValue(sheetName, "G1", "Impr Diff")
	for i, r := range data {
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", i+2), r.TownShipBorough)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", i+2), r.OldLandAssmt)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", i+2), r.OldImprAssmt)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", i+2), r.NewLandAssmt)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", i+2), r.NewImprAssmt)
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", i+2), r.LandDiff)
		f.SetCellValue(sheetName, fmt.Sprintf("G%d", i+2), r.ImprDiff)
	}
	if err := f.SaveAs("./valuationchange.xlsx"); err != nil {
		log.Fatal(err)
	}

	// Set transaction success flag
	ok = true
}
