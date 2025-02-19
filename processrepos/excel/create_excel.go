package excel

import (
	"fmt"

	"github.com/s3pweb/gitArchiveS3Report/utils/structs"
	"github.com/s3pweb/gitArchiveS3Report/utils/styles"
	"github.com/xuri/excelize/v2"
)

func CreateExcelFile(Branches []structs.BranchInfo) (*excelize.File, error) {
	f := excelize.NewFile()

	// Create sheets
	allBranchesSheet := "Branches"
	mainBranchesSheet := "Main Branches"
	developBranchesSheet := "Develop Branches"

	// Add sheets
	index := f.NewSheet(allBranchesSheet)
	f.NewSheet(mainBranchesSheet)
	f.NewSheet(developBranchesSheet)

	// Set active sheet
	f.SetActiveSheet(index)

	fixedColumns := countFixedColumns()
	maxFilesToSearch := 0
	maxTermsToSearch := 0

	for _, branch := range Branches {
		if len(branch.FilesToSearch) > maxFilesToSearch {
			maxFilesToSearch = len(branch.FilesToSearch)
		}
		if len(branch.TermsToSearch) > maxTermsToSearch {
			maxTermsToSearch = len(branch.TermsToSearch)
		}
	}

	totalColumns := fixedColumns + maxFilesToSearch + maxTermsToSearch

	// Create header style
	headerStyle, err := styles.CreateHeaderStyle(f)
	if err != nil {
		return nil, err
	}

	// Apply header style and set row height for each sheet
	sheets := []string{allBranchesSheet, mainBranchesSheet, developBranchesSheet}
	for _, sheet := range sheets {
		// Apply header style
		f.SetCellStyle(sheet, "A1", fmt.Sprintf("%c1", 'A'+totalColumns-1), headerStyle)
		f.SetRowHeight(sheet, 1, 40)

		// Set column widths
		for col := 'A'; col < 'A'+rune(totalColumns); col++ {
			f.SetColWidth(sheet, string(col), string(col), 20)
		}

		// Freeze the first row
		if err := f.SetPanes(sheet, `{
			"freeze": true,
			"split": false,
			"x_split": 0,
			"y_split": 1,
			"top_left_cell": "A2",
			"active_pane": "bottomLeft"
		}`); err != nil {
			return nil, err
		}
	}
	f.DeleteSheet("Sheet1")

	return f, nil
}
