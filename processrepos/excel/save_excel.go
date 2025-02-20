package excel

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/s3pweb/gitArchiveS3Report/utils/logger"
	"github.com/xuri/excelize/v2"
)

func SaveExcelFile(f *excelize.File, workspacePath, outputDir string, logger *logger.Logger) error {
	workspace := filepath.Base(workspacePath)

	// Format current time
	currentTime := time.Now()
	fileName := fmt.Sprintf("%s_report_%s_%s.xlsx",
		workspace,
		currentTime.Format("2006-01-02"),
		currentTime.Format("15-04"))

	excelFileName := filepath.Join(outputDir, fileName)
	if err := f.SaveAs(excelFileName); err != nil {
		return err
	}
	logger.Info("Excel file saved: %s", excelFileName)
	return nil
}
