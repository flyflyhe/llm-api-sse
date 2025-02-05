package tool

import (
	"testing"
)

// TestGetMonthListV2 测试 GetMonthListV2 函数
func TestGetMonthListV2(t *testing.T) {
	// 定义测试用例
	tests := []struct {
		start   string
		end     string
		layout  string
		want    []string
		wantErr bool
	}{
		{
			start:   "2023-01-01",
			end:     "2023-03-01",
			layout:  "2006-01-02",
			want:    []string{"2023-01-01", "2023-02-01", "2023-03-01"},
			wantErr: false,
		},
		// 可以添加更多测试用例
	}

	for _, tt := range tests {
		t.Run(tt.start+"_"+tt.end+"_"+tt.layout, func(t *testing.T) {
			got, err := GetMonthListV2(tt.start, tt.end, tt.layout)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetMonthListV2() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !isEqual(got, tt.want) {
				t.Errorf("GetMonthListV2() got = %v, want %v", got, tt.want)
			}
		})
	}
}

// isEqual 判断两个字符串切片是否相等
func isEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
