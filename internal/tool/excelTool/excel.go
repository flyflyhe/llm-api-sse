package excelTool

import (
	"custom_insurance/pkg/logging"
	"github.com/xuri/excelize/v2"
	"strconv"
)

// 1:A 2:B

func GetColNameByIndex(i, j int32) string {
	res, err := excelize.ColumnNumberToName(int(i) + 1)
	if err != nil {
		logging.Logger.Sugar().Error(err)
	}

	return res + strconv.Itoa(int(j))
}

func GetColNameByIndexV2(i, j int32) string {
	if i < 26 {
		return string('A'+i) + strconv.Itoa(int(j))
	}

	a := (i - 26) / 26

	var b int32
	if i > 26 {
		b = (i - 26) % 26
	}

	return string('A'+a) + string('A'+b) + strconv.Itoa(int(j))
}

func GetColNameByIndexV3(i, j int32) string {
	res, err := excelize.ColumnNumberToName(int(i))
	if err != nil {
		logging.Logger.Sugar().Error(err)
	}

	return res + strconv.Itoa(int(j))
}

func GetIndexByAxis(axis string) (int, int, error) {
	col, rowInt, err := excelize.SplitCellName(axis)
	if err != nil {
		return 0, 0, err
	}

	colInt, err := excelize.ColumnNameToNumber(col)
	if err != nil {
		return 0, 0, err
	}

	return rowInt - 1, colInt - 1, nil
}

func GetRows(f *excelize.File, sheet string, opts ...excelize.Options) ([][]string, error) {
	rows, err := f.Rows(sheet)
	if err != nil {
		return nil, err
	}
	results, cur, max := make([][]string, 0, 64), 0, 0
	for rows.Next() {
		cur++
		if cur == 5000 {
			break
		}
		row, err := rows.Columns(opts...)
		if err != nil {
			break
		}
		results = append(results, row)
		if len(row) > 0 {
			max = cur
		}
	}
	return results[:max], rows.Close()
}
