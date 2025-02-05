package logging

import (
	"custom_insurance/configs"
	"testing"
)

func init() {
	configs.InitConfigAuto("config.yaml")
	InitLogger()
}

func TestInitLogger(t *testing.T) {
	zapLog.Warn("hhhh ")
}
