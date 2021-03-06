package log

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"runtime"
	"time"
)

func serialize(l level, ctx context.Context, data interface{}, restOfData ...interface{}) []byte {
	mergedData := append([]interface{}{data}, restOfData...)

	pc, file, line, _ := runtime.Caller(3)
	fname := "UNKNOWN"
	if fn := runtime.FuncForPC(pc); fn != nil {
		fname = fn.Name()
	}

	m := message{
		Time:      time.Now(),
		Level:     l,
		Message:   fmt.Sprintf("%v", data),
		RequestID: GetRequestId(ctx),
		Source: source{
			file,
			line,
			fname,
		},
	}

	for _, d := range mergedData {
		if err, ok := d.(error); ok {
			m.Data = append(m.Data, newDataItem(err))
			continue
		}

		m.Data = append(m.Data, newDataItem(d))
	}

	b, err := json.Marshal(m)
	if err != nil {
		m.Data = []dataItem{newDataItem("Unable to serialize log message"), newDataItem(err)}
		b, err = json.Marshal(m)
		if err != nil {
			panic("unable to marshal secondary log message")
		}
	}

	return append(b, []byte("\n")...)
}

func newDataItem(i interface{}) dataItem {
	typ := reflect.TypeOf(i).PkgPath()
	if typ != "" {
		typ = typ + "."
	}

	f := dataItem{
		typ + reflect.TypeOf(i).Name(),
		data{i},
		fmt.Sprintf("%#v", i),
		map[string]data{},
	}

	for name, convert := range converters {
		fct, err := convert(i)
		if err != nil {
			fct = "UNABLE TO CONVERT: " + err.Error()
		}

		if fct != nil {
			f.Converted[name] = data{fct}
		}
	}

	return f
}

type dataItem struct {
	Type      string
	Marshaled data
	Printed   string
	Converted map[string]data
}

type data struct {
	i interface{}
}

func (f data) MarshalJSON() ([]byte, error) {
	b, err := json.Marshal(f.i)
	switch err.(type) {
	case *json.UnsupportedTypeError, *json.UnsupportedValueError:
		b, err = json.Marshal(struct{ SerializationError string }{"serialization error: " + err.Error()})
	}

	return b, err
}

type message struct {
	Time      time.Time
	Level     level
	RequestID string
	Message   string
	Source    source
	Data      []dataItem
}

type source struct {
	FilePath     string
	LineNumber   int
	FunctionName string
}
