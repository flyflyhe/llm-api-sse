package tool

import (
	"custom_insurance/pkg/logging"
	"encoding/json"
	"github.com/jinzhu/now"
	"github.com/shopspring/decimal"
	"strings"
	"time"
)

func Ternary(condition bool, trueVal, falseVal interface{}) interface{} {
	if condition {
		return trueVal
	}
	return falseVal
}

func TernaryFunc(condition bool, f1 func() interface{}, f2 func() interface{}) interface{} {
	if condition {
		return f1()
	} else {
		return f2()
	}
}

func SafeGetStrVal(s []string, i int) string {
	if i < len(s) {
		return s[i]
	}

	return ""
}

func FenToYuan[T int | int64 | int32 | uint | uint32 | uint64](total T) decimal.Decimal {
	f := decimal.NewFromInt(int64(total))

	return f.Div(decimal.NewFromInt(100))
}

func YuanToFen(decimal2 decimal.Decimal) int64 {
	return decimal2.Mul(decimal.NewFromInt(100)).BigInt().Int64()
}

func ToJson(v interface{}) string {
	j, err := json.Marshal(v)
	if err != nil {
		logging.Logger.Sugar().Error("json encode :", err.Error())
	}

	return string(j)
}

func GetLastMonth(month, layout string) (string, error) {
	parse, err := time.Parse(layout, month)
	if err != nil {
		return "", err
	}
	return now.With(parse).BeginningOfMonth().AddDate(0, 0, -1).Format(layout), nil
}

func GetNextMonth(month, layout string) (string, error) {
	parse, err := time.Parse(layout, month)
	if err != nil {
		return "", err
	}
	return now.With(parse).EndOfMonth().AddDate(0, 0, 1).Format(layout), nil
}

func GetMonthList(start, end, layout string) ([]string, error) {
	result := make([]string, 0)
	c := start
	for strings.Compare(c, end) < 0 {
		logging.Logger.Sugar().Info("月份", c)
		result = append(result, c)
		cTime, err := time.Parse(layout, c)
		if err != nil {
			return nil, err
		}
		c = now.With(cTime).EndOfMonth().AddDate(0, 0, 2).Format(layout)
	}

	return result, nil
}

func GetMonthListV2(start, end, layout string) ([]string, error) {
	result := make([]string, 0)
	c := start
	for strings.Compare(c, end) <= 0 {
		logging.Logger.Sugar().Info("月份", c)
		result = append(result, c)
		cTime, err := time.Parse(layout, c)
		if err != nil {
			return nil, err
		}
		c = now.With(cTime).EndOfMonth().AddDate(0, 0, 2).Format(layout)
	}

	return result, nil
}

func AsyncTask(f func() error) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				logging.Logger.Sugar().Error(err)
			}
		}()

		if err := f(); err != nil {
			logging.Logger.Sugar().Error(err)
		}
	}()
}
