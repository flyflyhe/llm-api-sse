package excelTool

import (
	"custom_insurance/pkg/logging"
	"fmt"
	"github.com/xuri/excelize/v2"
)

type Header struct {
	MergeRow int32
	MergeCol int32
	Title    string
}

type ExcelGenerateV2 struct {
	Headers   [][]Header
	Data      [][]interface{}
	ExcelFile *excelize.File
}

func NewExcelGenerateV2(headers [][]Header, data [][]interface{}) *ExcelGenerateV2 {
	return &ExcelGenerateV2{Headers: headers, Data: data}
}

func (e *ExcelGenerateV2) Generate() error {
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
	//构建header

	row := int32(1)
	for _, header := range e.Headers {
		col := int32(1)
		for _, h := range header {
			fmt.Println(GetColNameByIndexV3(col, row), h.Title)
			if err = f.SetCellValue(Sheet1, GetColNameByIndexV3(col, row), h.Title); err != nil {
				return err
			}
			if h.MergeCol > 1 {
				col += int32(h.MergeCol)
			} else {
				col++
			}
		}
		row++
	}

	//合并
	row = int32(1)
	for _, header := range e.Headers {
		col := int32(1)
		for _, h := range header {
			if h.MergeCol > 1 || h.MergeRow > 1 {
				fmt.Println("合并行", GetColNameByIndexV3(col, row), "-", GetColNameByIndexV3(col+h.MergeCol-1, h.MergeRow))
				if err = f.MergeCell(Sheet1, GetColNameByIndexV3(col, row), GetColNameByIndexV3(col+h.MergeCol-1, h.MergeRow)); err != nil {
					logging.Logger.Sugar().Error(err)
				}
			}
			col += h.MergeCol
		}
		row++
	}

	index := len(e.Headers) + 1
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
