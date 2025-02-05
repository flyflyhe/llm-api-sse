package excelTool

import (
	"custom_insurance/pkg/logging"
	"fmt"
	"github.com/xuri/excelize/v2"
)

const Sheet1 = "Sheet1"

type ExcelGenerate struct {
	Headers   []interface{}
	Data      [][]interface{}
	ExcelFile *excelize.File
}

func NewExcelGenerate(headers []interface{}, data [][]interface{}) *ExcelGenerate {
	return &ExcelGenerate{Headers: headers, Data: data}
}

func (e *ExcelGenerate) Generate() error {
	f := excelize.NewFile()
	e.ExcelFile = f
	defer func() {
		if err := f.Close(); err != nil {
			logging.Logger.Sugar().Error(err)
		}
	}()

	//设置元素居中
	styleId, err := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
	})
	if err != nil {
		logging.Logger.Sugar().Error(err)
		return err
	}

	if err = f.SetColWidth(Sheet1, "A", "CZ", 24); err != nil {
		logging.Logger.Sugar().Error(err)

		return err
	}

	if err = f.SetRowStyle(Sheet1, 1, 10000, styleId); err != nil {
		logging.Logger.Sugar().Error(err)
		return err
	}
	// Create a new sheet.
	indexSheet, err := f.NewSheet(Sheet1)
	if err != nil {
		logging.Logger.Sugar().Error(err)
		return err
	}
	f.SetActiveSheet(indexSheet)
	index := 1
	if err := f.SetSheetRow(Sheet1, fmt.Sprintf("A%d", index), &e.Headers); err != nil {
		logging.Logger.Sugar().Error(err)
		return err
	}
	index++
	for _, v := range e.Data {
		v := v
		if err := f.SetSheetRow(Sheet1, fmt.Sprintf("A%d", index), &v); err != nil {
			logging.Logger.Sugar().Error(err)
			return err
		}
		index++
	}

	return nil
}
