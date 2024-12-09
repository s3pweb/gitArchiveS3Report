package excel

import (
	"path/filepath"

	"github.com/s3pweb/gitArchiveS3Report/utils/logger"
	"github.com/xuri/excelize/v2"
)

func SaveExcelFile(f *excelize.File, basePath string, logger *logger.Logger) error {
	excelFileName := filepath.Join(basePath, "workspace_report.xlsx")
	logger.Info("Saving Excel file: %s", excelFileName)
	if err := f.SaveAs(excelFileName); err != nil {
		return err
	}
	logger.Info("Excel file saved: %s", excelFileName)
	return nil
}
