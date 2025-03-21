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

	currentTime := time.Now()
	dateStr := currentTime.Format("2006-01-02")
	hourStr := currentTime.Format("15h04")

	fileName := fmt.Sprintf("%s_report_%s_%s.xlsx", workspace, dateStr, hourStr)

	excelFileName := filepath.Join(outputDir, fileName)
	if err := f.SaveAs(excelFileName); err != nil {
		return err
	}
	logger.Info("Excel file saved: %s", excelFileName)
	return nil
}
