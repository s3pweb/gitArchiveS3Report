package excel

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode"

	"github.com/s3pweb/gitArchiveS3Report/config"
	"github.com/s3pweb/gitArchiveS3Report/utils/structs"
	"github.com/s3pweb/gitArchiveS3Report/utils/styles"
	"github.com/xuri/excelize/v2"
	"golang.org/x/text/unicode/norm"
)

func WriteBranchInfoToExcel(f *excelize.File, allBranches, mainBranches, developBranches []structs.BranchInfo, includeDevSheets bool) error {
	allBranchesSheet := "Branches"
	mainBranchesSheet := "Main Branches"
	developBranchesSheet := "Develop Branches"

	err := writeDataToSheet(f, allBranchesSheet, allBranches)
	if err != nil {
		return err
	}
	err = writeDataToSheet(f, mainBranchesSheet, mainBranches)
	if err != nil {
		return err
	}
	err = writeDataToSheet(f, developBranchesSheet, developBranches)
	if err != nil {
		return err
	}

	if includeDevSheets {
		err = createDeveloperSheets(f, allBranches)
		if err != nil {
			return err
		}
	}

	return nil
}

func writeDataToSheet(f *excelize.File, sheet string, branchesInfo []structs.BranchInfo) error {
	cfg := config.Get()
	columns := cfg.App.DefaultColumns
	columns = append(columns, cfg.App.TermsToSearch...)
	columns = append(columns, cfg.App.FilesToSearch...)

	sortBranchesByLastCommit(branchesInfo)
	nbrcolumn := 'A'

	// Create maps to store totals for terms and files
	termTotals := make(map[string]int)
	fileTotals := make(map[string]int)
	repoCount := countUniqueRepos(branchesInfo)

	// Initialize maps for terms and files
	for _, term := range cfg.App.TermsToSearch {
		termTotals[term] = 0
	}
	for _, file := range cfg.App.FilesToSearch {
		fileTotals[file] = 0
	}

	// Calculate total number of searchable items
	totalSearchableItems := len(cfg.App.TermsToSearch) + len(cfg.App.FilesToSearch)

	// Write column headers and data
	for _, column := range columns {
		row := 2
		for _, branchInfo := range branchesInfo {
			styles.SetOneHeader(f, sheet, strings.ToUpper(removeRegex(column)), nbrcolumn)

			// If column is "Count", calculate the total number of TRUE values
			if column == "Count" {
				trueCount := 0
				// Count only TRUE values in TermsToSearch
				for _, term := range cfg.App.TermsToSearch {
					if val, exists := branchInfo.TermsToSearch[term]; exists && val {
						trueCount++
					}
				}
				// Count only TRUE values in FilesToSearch
				for _, file := range cfg.App.FilesToSearch {
					if val, exists := branchInfo.FilesToSearch[file]; exists && val {
						trueCount++
					}
				}
				cell := fmt.Sprintf("%c%d", nbrcolumn, row)
				f.SetCellValue(sheet, cell, fmt.Sprintf("%d/%d", trueCount, totalSearchableItems))
				cellStyle, _ := styles.CreateCellStyle(f)
				f.SetCellStyle(sheet, cell, cell, cellStyle)
			} else {
				err := writeFieldToColumn(f, sheet, row, column, nbrcolumn, branchInfo)
				if err != nil {
					return err
				}
			}

			// Update totals for terms and files only when explicitly TRUE
			if val, exists := branchInfo.TermsToSearch[column]; exists && val {
				termTotals[column]++
			}
			if val, exists := branchInfo.FilesToSearch[column]; exists && val {
				fileTotals[column]++
			}
			row++
		}

		// Add totals row after all data
		cell := fmt.Sprintf("%c%d", nbrcolumn, row)
		cellStyle, _ := styles.CreateCellStyle(f)

		// Write totals for terms and files
		if termTotals[column] > 0 {
			percentage := float64(termTotals[column]) / float64(repoCount) * 100
			f.SetCellValue(sheet, cell, fmt.Sprintf("%d/%d (%.1f%%)", termTotals[column], repoCount, percentage))
			f.SetCellStyle(sheet, cell, cell, cellStyle)
		} else if fileTotals[column] > 0 {
			percentage := float64(fileTotals[column]) / float64(repoCount) * 100
			f.SetCellValue(sheet, cell, fmt.Sprintf("%d/%d (%.1f%%)", fileTotals[column], repoCount, percentage))
			f.SetCellStyle(sheet, cell, cell, cellStyle)
		} else if nbrcolumn == 'H' {
			f.SetCellValue(sheet, cell, "TOTAL")
			f.SetCellStyle(sheet, cell, cell, cellStyle)
		}

		f.SetRowHeight(sheet, row, 30)
		nbrcolumn++
	}
	return nil
}

func countUniqueRepos(branchesInfo []structs.BranchInfo) int {
	repos := make(map[string]bool)
	for _, info := range branchesInfo {
		repos[info.RepoName] = true
	}
	return len(repos)
}

func writeFieldToColumn(f *excelize.File, sheet string, row int, fieldName string, col rune, branchInfo interface{}) error {
	cellStyle, err := styles.CreateCellStyle(f)
	if err != nil {
		return err
	}
	falseStyle, err := styles.FalseCells(f)
	if err != nil {
		return err
	}

	v := reflect.ValueOf(branchInfo)
	fieldValue := v.FieldByName(fieldName)

	if !fieldValue.IsValid() {
		branchInfoTyped := branchInfo.(structs.BranchInfo)
		// Check specifically in FilesToSearch and TermsToSearch maps
		if val, exists := branchInfoTyped.FilesToSearch[fieldName]; exists {
			f.SetCellValue(sheet, fmt.Sprintf("%c%d", col, row), strings.ToUpper(fmt.Sprintf("%v", val)))
			if val {
				f.SetCellStyle(sheet, fmt.Sprintf("%c%d", col, row), fmt.Sprintf("%c%d", col, row), cellStyle)
			} else {
				f.SetCellStyle(sheet, fmt.Sprintf("%c%d", col, row), fmt.Sprintf("%c%d", col, row), falseStyle)
			}
			f.SetRowHeight(sheet, row, 30)
			return nil
		}
		if val, exists := branchInfoTyped.TermsToSearch[fieldName]; exists {
			f.SetCellValue(sheet, fmt.Sprintf("%c%d", col, row), strings.ToUpper(fmt.Sprintf("%v", val)))
			if val {
				f.SetCellStyle(sheet, fmt.Sprintf("%c%d", col, row), fmt.Sprintf("%c%d", col, row), cellStyle)
			} else {
				f.SetCellStyle(sheet, fmt.Sprintf("%c%d", col, row), fmt.Sprintf("%c%d", col, row), falseStyle)
			}
			f.SetRowHeight(sheet, row, 30)
			return nil
		}
		return fmt.Errorf("field %s not found in struct", fieldName)
	}

	cell := fmt.Sprintf("%c%d", col, row)
	if fieldName == "LastCommitDate" {
		f.SetCellValue(sheet, cell, fieldValue.Interface().(time.Time).Format("2006-01-02 15:04"))
	} else if fieldName == "LastDeveloperPercentage" || fieldName == "TopDeveloperPercentage" {
		f.SetCellValue(sheet, cell, fmt.Sprintf("%.2f%%", fieldValue.Float()))
	} else {
		f.SetCellValue(sheet, cell, fieldValue.Interface())
	}

	f.SetRowHeight(sheet, row, 30)
	f.SetCellStyle(sheet, cell, cell, cellStyle)
	return nil
}

func countTrueInMap(values map[string]bool) int {
	count := 0
	for _, value := range values {
		if value {
			count++
		}
	}
	return count
}

func sortBranchesByLastCommit(branches []structs.BranchInfo) {
	sort.Slice(branches, func(i, j int) bool {
		return branches[i].LastCommitDate.After(branches[j].LastCommitDate)
	})
}

func createDeveloperSheets(f *excelize.File, branchesInfo []structs.BranchInfo) error {
	developerSheets := make(map[string]bool)

	configColumns, err := getConfigColumn(".config")
	if err != nil {
		return err
	}
	termsColumns, err := getTermsColumn(".config")
	if err != nil {
		return err
	}
	filesColumns, err := getFilesColumn(".config")
	if err != nil {
		return err
	}
	columns := append(configColumns, append(termsColumns, filesColumns...)...)

	sortBranchesByLastCommit(branchesInfo)

	var developers []string
	for _, branchInfo := range branchesInfo {
		developers = append(developers, branchInfo.LastDeveloper, branchInfo.TopDeveloper)
	}
	for _, developer := range developers {
		nbrcolumn := 'A'
		for _, column := range columns {
			row := 2
			for _, branchInfo := range branchesInfo {
				if developer == "" {
					continue
				}
				developer = strings.ToLower(RemoveAccentsAndSpecialChars(developer))

				if !developerSheets[developer] {
					f.NewSheet(developer)
					developerSheets[developer] = true
					f.SetRowHeight(developer, 1, 40)
					for col := 'A'; col < 'Z'; col++ {
						f.SetColWidth(developer, string(col), string(col), 20)
					}
				}
				if strings.ToLower(RemoveAccentsAndSpecialChars(branchInfo.LastDeveloper)) == developer ||
					strings.ToLower(RemoveAccentsAndSpecialChars(branchInfo.TopDeveloper)) == developer {
					styles.SetOneHeader(f, developer, strings.ToUpper(removeRegex(column)), nbrcolumn)
					err := writeFieldToColumn(f, developer, row, column, nbrcolumn, branchInfo)
					if err != nil {
						return err
					}
					row++

				}
			}
			nbrcolumn++

		}
	}
	return nil
}

func RemoveAccentsAndSpecialChars(s string) string {
	t := norm.NFD.String(s)

	var b strings.Builder
	for _, r := range t {
		if unicode.IsLetter(r) || unicode.IsNumber(r) {
			b.WriteRune(r)
		}

	}
	re := regexp.MustCompile("[^a-zA-Z0-9]+")
	return re.ReplaceAllString(b.String(), "")

}

func getConfigColumn(config string) ([]string, error) {
	file, err := os.Open(config)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "DEFAULT_COLUMN=") {
			columns := strings.TrimPrefix(line, "DEFAULT_COLUMN=")
			return strings.Split(columns, ";"), nil
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return nil, fmt.Errorf("DEFAULT_COLUMN= line not found in the config file")
}

func getFilesColumn(config string) ([]string, error) {
	file, err := os.Open(config)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "FILES_TO_SEARCH=") {
			columns := strings.TrimPrefix(line, "FILES_TO_SEARCH=")
			return strings.Split(columns, ";"), nil
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return nil, fmt.Errorf("FILES_TO_SEARCH= line not found in the config file")
}

func getTermsColumn(config string) ([]string, error) {
	file, err := os.Open(config)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "TERMS_TO_SEARCH=") {
			columns := strings.TrimPrefix(line, "TERMS_TO_SEARCH=")
			return strings.Split(columns, ";"), nil
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return nil, fmt.Errorf("TERMS_TO_SEARCH= line not found in the config file")
}

func removeRegex(s string) string {
	re := regexp.MustCompile(`\(\?i\)|[^a-zA-Z0-9]`)
	return re.ReplaceAllString(s, "")
}

func countFixedColumns() int {
	cfg := config.Get()
	return len(cfg.App.DefaultColumns)
}
