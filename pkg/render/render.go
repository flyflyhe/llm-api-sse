package render

import (
	"bytes"
	"encoding/json"
	hjson "github.com/cloudwego/hertz/pkg/common/json"
	"github.com/cloudwego/hertz/pkg/protocol"
	jsontime "github.com/liamylian/jsontime/v2/v2"
	"time"
)

var cJson = jsontime.ConfigWithCustomTimeFormat

// JSONMarshaler customize json.Marshal as you like
type JSONMarshaler func(v interface{}) ([]byte, error)

var jsonMarshalFunc JSONMarshaler

func init() {
	timeZoneShanghai, _ := time.LoadLocation("Asia/Shanghai")
	jsontime.AddTimeFormatAlias("sql_datetime", "2006-01-02 15:04:05")
	jsontime.AddLocaleAlias("shanghai", timeZoneShanghai)
	ResetJSONMarshal(hjson.Marshal)
}

func ResetJSONMarshal(fn JSONMarshaler) {
	jsonMarshalFunc = fn
}

func ResetStdJSONMarshal() {
	ResetJSONMarshal(json.Marshal)
}

// JSONRender JSON contains the given interface object.
type JSONRender struct {
	Data interface{}
}

var jsonContentType = "application/json; charset=utf-8"

// Render (JSON) writes data with custom ContentType.
func (r JSONRender) Render(resp *protocol.Response) error {
	writeContentType(resp, jsonContentType)
	jsonBytes, err := cJson.Marshal(r.Data)
	if err != nil {
		return err
	}

	resp.AppendBody(jsonBytes)
	return nil
}

// WriteContentType (JSON) writes JSON ContentType.
func (r JSONRender) WriteContentType(resp *protocol.Response) {
	writeContentType(resp, jsonContentType)
}

// PureJSON contains the given interface object.
type PureJSON struct {
	Data interface{}
}

// Render (JSON) writes data with custom ContentType.
func (r PureJSON) Render(resp *protocol.Response) (err error) {
	r.WriteContentType(resp)
	buffer := new(bytes.Buffer)
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	err = encoder.Encode(r.Data)
	if err != nil {
		return
	}
	resp.AppendBody(buffer.Bytes())
	return
}

// WriteContentType (JSON) writes JSON ContentType.
func (r PureJSON) WriteContentType(resp *protocol.Response) {
	writeContentType(resp, jsonContentType)
}

// IndentedJSON contains the given interface object.
type IndentedJSON struct {
	Data interface{}
}

// Render (IndentedJSON) marshals the given interface object and writes it with custom ContentType.
func (r IndentedJSON) Render(resp *protocol.Response) (err error) {
	writeContentType(resp, jsonContentType)
	jsonBytes, err := jsonMarshalFunc(r.Data)
	if err != nil {
		return err
	}
	var buf bytes.Buffer
	err = json.Indent(&buf, jsonBytes, "", "    ")
	if err != nil {
		return err
	}
	resp.AppendBody(buf.Bytes())
	return nil
}

// WriteContentType (JSON) writes JSON ContentType.
func (r IndentedJSON) WriteContentType(resp *protocol.Response) {
	writeContentType(resp, jsonContentType)
}

func writeContentType(resp *protocol.Response, value string) {
	resp.Header.SetContentType(value)
}
