package excelTool

import (
	"custom_insurance/configs"
	"custom_insurance/pkg/logging"
	"github.com/shopspring/decimal"
	"log"
	"testing"
)

func init() {
	configs.InitConfigAuto("config.yaml")
	logging.InitLogger()
}

func TestGetColNameByIndex(t *testing.T) {
	log.Println(GetColNameByIndex(0, 1))
	log.Println(GetColNameByIndex(2, 1))
	log.Println(GetColNameByIndex(25, 1))
	log.Println(GetColNameByIndex(26, 1))
	log.Println(GetColNameByIndex(27, 1))
	log.Println(GetColNameByIndex(28, 1))
	log.Println(GetColNameByIndex(29, 1))
	log.Println(GetColNameByIndex(30, 1))
	log.Println(GetColNameByIndex(51, 1))
	log.Println(GetColNameByIndex(52, 1))
}

func TestGetRows(t *testing.T) {
	if copyMoney, err := decimal.NewFromString("0"); err != nil {
		t.Error(err)
	} else {
		t.Log(copyMoney)
	}
}
